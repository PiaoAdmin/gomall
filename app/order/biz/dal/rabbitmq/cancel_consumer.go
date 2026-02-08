package rabbitmq

import (
	"context"
	"encoding/json"
	"sync"
	"sync/atomic"
	"time"

	"github.com/PiaoAdmin/pmall/app/order/biz/dal/mysql"
	"github.com/PiaoAdmin/pmall/app/order/biz/model"
	"github.com/PiaoAdmin/pmall/app/order/biz/rpc"
	"github.com/PiaoAdmin/pmall/app/order/conf"
	"github.com/PiaoAdmin/pmall/rpc_gen/product"
	"github.com/cloudwego/kitex/pkg/klog"
	amqp "github.com/rabbitmq/amqp091-go"
)

// CancelConsumer 订单取消消费者
type CancelConsumer struct {
	running  int32
	stopChan chan struct{}
	wg       sync.WaitGroup
}

var cancelConsumer *CancelConsumer

// StartCancelConsumer 启动订单取消消费者
func StartCancelConsumer(ctx context.Context) {
	cancelConsumer = &CancelConsumer{
		stopChan: make(chan struct{}),
	}
	cancelConsumer.Start(ctx)
}

// StopCancelConsumer 停止取消消费者
func StopCancelConsumer() {
	if cancelConsumer != nil {
		cancelConsumer.Stop()
	}
}

// Start 启动消费者
func (c *CancelConsumer) Start(ctx context.Context) {
	if !atomic.CompareAndSwapInt32(&c.running, 0, 1) {
		klog.Warn("Cancel consumer already running")
		return
	}

	cfg := conf.GetConf().RabbitMQ
	workerCount := cfg.WorkerCount
	if workerCount <= 0 {
		workerCount = 2 // 取消消费者使用较少的工作线程
	}

	for i := 0; i < workerCount; i++ {
		c.wg.Add(1)
		go c.consume(ctx, i)
	}

	klog.Infof("Order cancel consumer started with %d workers", workerCount)
}

// Stop 停止消费者
func (c *CancelConsumer) Stop() {
	if !atomic.CompareAndSwapInt32(&c.running, 1, 0) {
		return
	}
	close(c.stopChan)
	c.wg.Wait()
	klog.Info("Order cancel consumer stopped")
}

// consume 消费取消消息
func (c *CancelConsumer) consume(ctx context.Context, workerID int) {
	defer c.wg.Done()

	cfg := conf.GetConf().RabbitMQ

	ch, err := Connection.Channel()
	if err != nil {
		klog.Errorf("Cancel worker %d: Failed to create channel: %v", workerID, err)
		return
	}
	defer ch.Close()

	if err := ch.Qos(cfg.PrefetchCount, 0, false); err != nil {
		klog.Errorf("Cancel worker %d: Failed to set QoS: %v", workerID, err)
		return
	}

	msgs, err := ch.Consume(
		OrderCancelQueue,
		"cancel-consumer-"+string(rune(workerID)),
		false, // auto-ack
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		klog.Errorf("Cancel worker %d: Failed to register consumer: %v", workerID, err)
		return
	}

	klog.Infof("Cancel worker %d: Started consuming from queue %s", workerID, OrderCancelQueue)

	for {
		select {
		case <-c.stopChan:
			klog.Infof("Cancel worker %d: Received stop signal", workerID)
			return
		case <-ctx.Done():
			klog.Infof("Cancel worker %d: Context cancelled", workerID)
			return
		case msg, ok := <-msgs:
			if !ok {
				klog.Warnf("Cancel worker %d: Message channel closed", workerID)
				return
			}
			c.handleCancelMessage(ctx, workerID, msg)
		}
	}
}

// handleCancelMessage 处理取消订单消息
func (c *CancelConsumer) handleCancelMessage(ctx context.Context, workerID int, msg amqp.Delivery) {
	startTime := time.Now()

	var cancelMsg OrderCancelMessage
	if err := json.Unmarshal(msg.Body, &cancelMsg); err != nil {
		klog.Errorf("Cancel worker %d: Failed to unmarshal message: %v", workerID, err)
		msg.Reject(false)
		return
	}

	klog.Infof("Cancel worker %d: Processing cancel for order: %s, created_at=%d",
		workerID, cancelMsg.OrderID, cancelMsg.CreatedAt)

	// 执行取消订单逻辑
	err := cancelOrderIfUnpaid(ctx, cancelMsg.OrderID)
	if err != nil {
		klog.Errorf("Cancel worker %d: Failed to cancel order %s: %v", workerID, cancelMsg.OrderID, err)
		// 取消失败也确认消息，避免无限重试
		// 可以记录日志或发送告警
		msg.Ack(false)
		return
	}

	msg.Ack(false)
	elapsed := time.Since(startTime)
	klog.Infof("Cancel worker %d: Order %s cancel processed in %v", workerID, cancelMsg.OrderID, elapsed)
}

// cancelOrderIfUnpaid 检查并取消未支付的订单
func cancelOrderIfUnpaid(ctx context.Context, orderID string) error {
	// 1. 查询订单状态
	var order model.Order
	if err := mysql.DB.Preload("Items").Where("order_id = ?", orderID).First(&order).Error; err != nil {
		klog.Infof("Order %s not found, may already be deleted", orderID)
		return nil // 订单不存在，视为已处理
	}

	// 2. 检查订单状态
	switch order.Status {
	case model.OrderStatePaid:
		// 已支付，不需要取消
		klog.Infof("Order %s already paid, skip cancel", orderID)
		return nil
	case model.OrderStateCanceled:
		// 已取消，不需要再次取消
		klog.Infof("Order %s already canceled", orderID)
		return nil
	case model.OrderStatePlaced:
		// 待支付状态，执行取消
		klog.Infof("Order %s is unpaid after timeout, canceling...", orderID)
	default:
		klog.Warnf("Order %s has unknown status: %s", orderID, order.Status)
		return nil
	}

	// 3. 释放库存
	if len(order.Items) > 0 {
		releaseItems := make([]*product.SkuDeductItem, 0, len(order.Items))
		for _, item := range order.Items {
			if item.SkuId == 0 || item.Quantity <= 0 {
				continue
			}
			releaseItems = append(releaseItems, &product.SkuDeductItem{
				SkuId: item.SkuId,
				Count: item.Quantity,
			})
		}

		if len(releaseItems) > 0 {
			_, err := rpc.ProductClient.ReleaseStock(ctx, &product.ReleaseStockRequest{
				OrderSn: orderID,
				Items:   releaseItems,
			})
			if err != nil {
				klog.Errorf("Failed to release stock for order %s: %v", orderID, err)
				// 继续执行取消，库存问题后续处理
			} else {
				klog.Infof("Stock released for order %s", orderID)
			}
		}
	}

	// 4. 更新订单状态为已取消
	if err := mysql.DB.Model(&model.Order{}).
		Where("order_id = ? AND status = ?", orderID, model.OrderStatePlaced).
		Update("status", model.OrderStateCanceled).Error; err != nil {
		klog.Errorf("Failed to update order status: %v", err)
		return err
	}

	klog.Infof("Order %s canceled successfully due to timeout", orderID)
	return nil
}
