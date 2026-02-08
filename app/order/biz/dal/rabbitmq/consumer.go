package rabbitmq

import (
	"context"
	"encoding/json"
	"sync"
	"sync/atomic"
	"time"

	"github.com/PiaoAdmin/pmall/app/order/biz/dal/mysql"
	"github.com/PiaoAdmin/pmall/app/order/biz/model"
	"github.com/PiaoAdmin/pmall/app/order/conf"
	"github.com/cloudwego/kitex/pkg/klog"
	amqp "github.com/rabbitmq/amqp091-go"
	"gorm.io/gorm"
)

const (
	MaxRetryCount = 3 // 最大重试次数
)

// Consumer 订单消息消费者
type Consumer struct {
	running  int32
	stopChan chan struct{}
	wg       sync.WaitGroup
}

// 全局消费者实例
var consumer *Consumer

// StartConsumer 启动订单消息消费者
func StartConsumer(ctx context.Context) {
	consumer = &Consumer{
		stopChan: make(chan struct{}),
	}
	consumer.Start(ctx)
}

// StopConsumer 停止消费者
func StopConsumer() {
	if consumer != nil {
		consumer.Stop()
	}
}

// Start 启动消费者
func (c *Consumer) Start(ctx context.Context) {
	if !atomic.CompareAndSwapInt32(&c.running, 0, 1) {
		klog.Warn("Consumer already running")
		return
	}

	cfg := conf.GetConf().RabbitMQ
	workerCount := cfg.WorkerCount
	if workerCount <= 0 {
		workerCount = 3
	}

	// 启动多个工作协程
	for i := 0; i < workerCount; i++ {
		c.wg.Add(1)
		go c.consume(ctx, i)
	}

	klog.Infof("Order consumer started with %d workers", workerCount)
}

// Stop 停止消费者
func (c *Consumer) Stop() {
	if !atomic.CompareAndSwapInt32(&c.running, 1, 0) {
		return
	}
	close(c.stopChan)
	c.wg.Wait()
	klog.Info("Order consumer stopped")
}

// consume 消费消息的工作协程
func (c *Consumer) consume(ctx context.Context, workerID int) {
	defer c.wg.Done()

	cfg := conf.GetConf().RabbitMQ

	// 创建独立的 channel
	ch, err := Connection.Channel()
	if err != nil {
		klog.Errorf("Worker %d: Failed to create channel: %v", workerID, err)
		return
	}
	defer ch.Close()

	// 设置 QoS
	if err := ch.Qos(cfg.PrefetchCount, 0, false); err != nil {
		klog.Errorf("Worker %d: Failed to set QoS: %v", workerID, err)
		return
	}

	// 开始消费
	msgs, err := ch.Consume(
		cfg.OrderQueue, // 队列名
		"order-consumer-"+string(rune(workerID)), // 消费者标签
		false, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		klog.Errorf("Worker %d: Failed to register consumer: %v", workerID, err)
		return
	}

	klog.Infof("Worker %d: Started consuming from queue %s", workerID, cfg.OrderQueue)

	for {
		select {
		case <-c.stopChan:
			klog.Infof("Worker %d: Received stop signal", workerID)
			return
		case <-ctx.Done():
			klog.Infof("Worker %d: Context cancelled", workerID)
			return
		case msg, ok := <-msgs:
			if !ok {
				klog.Warnf("Worker %d: Message channel closed", workerID)
				return
			}
			c.handleMessage(ctx, workerID, msg)
		}
	}
}

// handleMessage 处理单条消息
func (c *Consumer) handleMessage(ctx context.Context, workerID int, msg amqp.Delivery) {
	startTime := time.Now()

	var orderMsg OrderMessage
	if err := json.Unmarshal(msg.Body, &orderMsg); err != nil {
		klog.Errorf("Worker %d: Failed to unmarshal message: %v", workerID, err)
		// 解析失败，直接拒绝不重试
		msg.Reject(false)
		return
	}

	klog.Infof("Worker %d: Processing order: order_id=%s, user_id=%d", workerID, orderMsg.OrderID, orderMsg.UserID)

	// 处理订单写入数据库
	err := processOrderToDB(ctx, &orderMsg)
	if err != nil {
		klog.Errorf("Worker %d: Failed to process order %s: %v", workerID, orderMsg.OrderID, err)

		// 检查重试次数
		if orderMsg.Retry < MaxRetryCount {
			// 重新发送消息，增加重试计数
			orderMsg.Retry++
			if pubErr := PublishOrderMessage(ctx, &orderMsg); pubErr != nil {
				klog.Errorf("Worker %d: Failed to republish order %s: %v", workerID, orderMsg.OrderID, pubErr)
			}
		} else {
			klog.Errorf("Worker %d: Order %s exceeded max retries, sending to DLQ", workerID, orderMsg.OrderID)
		}
		// 拒绝消息，不重新入队（会进入死信队列）
		msg.Reject(false)
		return
	}

	// 成功处理，确认消息
	if err := msg.Ack(false); err != nil {
		klog.Errorf("Worker %d: Failed to ack message: %v", workerID, err)
	}

	elapsed := time.Since(startTime)
	klog.Infof("Worker %d: Order %s processed successfully in %v", workerID, orderMsg.OrderID, elapsed)
}

// processOrderToDB 将订单数据写入数据库
func processOrderToDB(ctx context.Context, msg *OrderMessage) error {
	// 先检查订单是否已存在（幂等性检查）
	var existingOrder model.Order
	if err := mysql.DB.Where("order_id = ?", msg.OrderID).First(&existingOrder).Error; err == nil {
		klog.Infof("Order %s already exists, skipping", msg.OrderID)
		return nil // 订单已存在，视为成功
	} else if err != gorm.ErrRecordNotFound {
		return err
	}

	// 构建订单模型
	order := &model.Order{
		OrderId: msg.OrderID,
		UserId:  msg.UserID,
		Email:   msg.Email,
		Status:  model.OrderStatePlaced,
		ShippingAddress: model.Address{
			Name:          msg.Address.Name,
			StreetAddress: msg.Address.StreetAddress,
			City:          msg.Address.City,
			ZipCode:       msg.Address.ZipCode,
		},
	}

	// 使用事务写入订单和订单项
	return mysql.DB.Transaction(func(tx *gorm.DB) error {
		// 创建订单
		if err := tx.Create(order).Error; err != nil {
			return err
		}

		// 创建订单项
		for _, item := range msg.Items {
			orderItem := &model.OrderItem{
				OrderId:  msg.OrderID,
				SkuId:    item.SkuID,
				SkuName:  item.SkuName,
				Price:    item.Price,
				Quantity: item.Quantity,
			}
			if err := tx.Create(orderItem).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

// GetConsumerStats 获取消费者统计信息
func GetConsumerStats() map[string]interface{} {
	if consumer == nil {
		return map[string]interface{}{
			"running": false,
		}
	}
	return map[string]interface{}{
		"running": atomic.LoadInt32(&consumer.running) == 1,
	}
}
