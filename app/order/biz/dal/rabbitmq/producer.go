package rabbitmq

import (
	"context"
	"encoding/json"
	"time"

	"github.com/PiaoAdmin/pmall/app/order/conf"
	"github.com/cloudwego/kitex/pkg/klog"
	amqp "github.com/rabbitmq/amqp091-go"
)

// OrderMessage 订单消息结构
type OrderMessage struct {
	OrderID   string             `json:"order_id"`
	UserID    uint64             `json:"user_id"`
	Email     string             `json:"email"`
	Address   OrderAddress       `json:"address"`
	Items     []OrderMessageItem `json:"items"`
	CreatedAt int64              `json:"created_at"`
	Retry     int                `json:"retry"` // 重试次数
}

type OrderAddress struct {
	Name          string `json:"name"`
	StreetAddress string `json:"street_address"`
	City          string `json:"city"`
	ZipCode       int32  `json:"zip_code"`
}

type OrderMessageItem struct {
	SkuID    uint64  `json:"sku_id"`
	SkuName  string  `json:"sku_name"`
	Price    float64 `json:"price"`
	Quantity int32   `json:"quantity"`
}

// PublishOrderMessage 发布订单消息到 RabbitMQ
func PublishOrderMessage(ctx context.Context, msg *OrderMessage) error {
	cfg := conf.GetConf().RabbitMQ

	body, err := json.Marshal(msg)
	if err != nil {
		klog.CtxErrorf(ctx, "Failed to marshal order message: %v", err)
		return err
	}

	// 设置消息属性
	publishing := amqp.Publishing{
		ContentType:  "application/json",
		DeliveryMode: amqp.Persistent, // 持久化消息
		Timestamp:    time.Now(),
		MessageId:    msg.OrderID, // 使用订单ID作为消息ID，便于去重
		Body:         body,
	}

	// 使用带超时的上下文发布消息
	pubCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err = Channel.PublishWithContext(
		pubCtx,
		cfg.OrderExchange, // 交换机
		"order.create",    // 路由键
		false,             // mandatory
		false,             // immediate
		publishing,
	)

	if err != nil {
		klog.CtxErrorf(ctx, "Failed to publish order message: %v", err)
		return err
	}

	klog.CtxInfof(ctx, "Order message published successfully: order_id=%s", msg.OrderID)
	return nil
}

// PublishOrderMessageWithConfirm 发布订单消息（带发布确认）
func PublishOrderMessageWithConfirm(ctx context.Context, msg *OrderMessage) error {
	cfg := conf.GetConf().RabbitMQ

	// 开启发布确认模式
	if err := Channel.Confirm(false); err != nil {
		klog.CtxWarnf(ctx, "Failed to enable confirm mode: %v", err)
		// 降级为普通发布
		return PublishOrderMessage(ctx, msg)
	}

	body, err := json.Marshal(msg)
	if err != nil {
		klog.CtxErrorf(ctx, "Failed to marshal order message: %v", err)
		return err
	}

	confirms := Channel.NotifyPublish(make(chan amqp.Confirmation, 1))

	publishing := amqp.Publishing{
		ContentType:  "application/json",
		DeliveryMode: amqp.Persistent,
		Timestamp:    time.Now(),
		MessageId:    msg.OrderID,
		Body:         body,
	}

	pubCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err = Channel.PublishWithContext(
		pubCtx,
		cfg.OrderExchange,
		"order.create",
		false,
		false,
		publishing,
	)

	if err != nil {
		klog.CtxErrorf(ctx, "Failed to publish order message: %v", err)
		return err
	}

	// 等待确认
	select {
	case confirm := <-confirms:
		if !confirm.Ack {
			klog.CtxErrorf(ctx, "Message was not confirmed: order_id=%s", msg.OrderID)
			return ErrMessageNotConfirmed
		}
		klog.CtxInfof(ctx, "Order message confirmed: order_id=%s", msg.OrderID)
	case <-time.After(5 * time.Second):
		klog.CtxErrorf(ctx, "Message confirmation timeout: order_id=%s", msg.OrderID)
		return ErrConfirmTimeout
	}

	return nil
}
