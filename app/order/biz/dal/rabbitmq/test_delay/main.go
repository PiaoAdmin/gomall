package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	rabbitMQURL   = "amqp://admin:123456@localhost:5672/"
	delayExchange = "order_delay_exchange"
	delayQueue    = "order_delay_queue"
	cancelQueue   = "order_cancel_queue"
	delayTTL      = 30 * time.Second // 测试用 30 秒
)

type OrderCancelMessage struct {
	OrderID   string `json:"order_id"`
	UserID    uint64 `json:"user_id"`
	CreatedAt int64  `json:"created_at"`
	Reason    string `json:"reason"`
}

func main() {
	fmt.Println("=== RabbitMQ 延迟队列测试 ===")
	fmt.Printf("延迟时间: %v\n\n", delayTTL)

	// 连接 RabbitMQ
	conn, err := amqp.Dial(rabbitMQURL)
	if err != nil {
		log.Fatalf("连接 RabbitMQ 失败: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("创建 Channel 失败: %v", err)
	}
	defer ch.Close()

	// 声明延迟交换机
	err = ch.ExchangeDeclare(delayExchange, "direct", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("声明交换机失败: %v", err)
	}
	fmt.Println("✓ 交换机已声明")

	// 声明取消队列（接收过期消息）
	_, err = ch.QueueDeclare(cancelQueue, true, false, false, false, nil)
	if err != nil {
		log.Fatalf("声明取消队列失败: %v", err)
	}
	err = ch.QueueBind(cancelQueue, "order.cancel", delayExchange, false, nil)
	if err != nil {
		log.Fatalf("绑定取消队列失败: %v", err)
	}
	fmt.Println("✓ 取消队列已声明并绑定")

	// 声明延迟队列（带 TTL 和死信转发）
	_, err = ch.QueueDeclare(
		delayQueue,
		true,
		false,
		false,
		false,
		amqp.Table{
			"x-message-ttl":             int32(delayTTL.Milliseconds()),
			"x-dead-letter-exchange":    delayExchange,
			"x-dead-letter-routing-key": "order.cancel",
		},
	)
	if err != nil {
		log.Fatalf("声明延迟队列失败: %v", err)
	}
	err = ch.QueueBind(delayQueue, "order.delay", delayExchange, false, nil)
	if err != nil {
		log.Fatalf("绑定延迟队列失败: %v", err)
	}
	fmt.Printf("✓ 延迟队列已声明 (TTL=%v)\n\n", delayTTL)

	// 发送测试消息
	testOrderID := fmt.Sprintf("TEST_ORDER_%d", time.Now().UnixNano())
	msg := OrderCancelMessage{
		OrderID:   testOrderID,
		UserID:    12345,
		CreatedAt: time.Now().Unix(),
		Reason:    "timeout",
	}

	body, _ := json.Marshal(msg)
	err = ch.PublishWithContext(
		context.Background(),
		delayExchange,
		"order.delay",
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Body:         body,
		},
	)
	if err != nil {
		log.Fatalf("发送消息失败: %v", err)
	}

	sendTime := time.Now()
	fmt.Printf("✓ 消息已发送到延迟队列\n")
	fmt.Printf("  订单 ID: %s\n", testOrderID)
	fmt.Printf("  发送时间: %s\n", sendTime.Format("15:04:05"))
	fmt.Printf("  预计到达取消队列: %s\n\n", sendTime.Add(delayTTL).Format("15:04:05"))

	// 启动消费者监听取消队列
	fmt.Println("正在监听取消队列...")
	msgs, err := ch.Consume(cancelQueue, "", false, false, false, false, nil)
	if err != nil {
		log.Fatalf("启动消费者失败: %v", err)
	}

	// 设置超时
	timeout := time.After(delayTTL + 10*time.Second)

	for {
		select {
		case msg := <-msgs:
			receiveTime := time.Now()
			var cancelMsg OrderCancelMessage
			json.Unmarshal(msg.Body, &cancelMsg)

			delay := receiveTime.Sub(sendTime)
			fmt.Printf("\n✓✓✓ 收到取消消息! ✓✓✓\n")
			fmt.Printf("  订单 ID: %s\n", cancelMsg.OrderID)
			fmt.Printf("  接收时间: %s\n", receiveTime.Format("15:04:05"))
			fmt.Printf("  实际延迟: %v\n", delay)
			fmt.Printf("  取消原因: %s\n", cancelMsg.Reason)

			msg.Ack(false)
			fmt.Println("\n=== 测试成功! 延迟队列工作正常 ===")
			return

		case <-timeout:
			fmt.Println("\n✗ 超时! 未收到消息")
			return
		}
	}
}
