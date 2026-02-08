package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/cloudwego/kitex/pkg/klog"
	amqp "github.com/rabbitmq/amqp091-go"
)

// OrderCancelMessage 订单取消消息结构
type OrderCancelMessage struct {
	OrderID   string `json:"order_id"`
	UserID    uint64 `json:"user_id"`
	CreatedAt int64  `json:"created_at"` // 订单创建时间
	Reason    string `json:"reason"`     // 取消原因
}

// PublishOrderCancelDelay 发布延迟取消订单消息
// 订单创建后调用此方法，消息将在 30 分钟后被消费
func PublishOrderCancelDelay(ctx context.Context, orderID string, userID uint64) error {
	msg := &OrderCancelMessage{
		OrderID:   orderID,
		UserID:    userID,
		CreatedAt: time.Now().Unix(),
		Reason:    "timeout", // 超时未支付
	}

	body, err := json.Marshal(msg)
	if err != nil {
		klog.CtxErrorf(ctx, "Failed to marshal cancel message: %v", err)
		return err
	}

	publishing := amqp.Publishing{
		ContentType:  "application/json",
		DeliveryMode: amqp.Persistent,
		Timestamp:    time.Now(),
		MessageId:    orderID + "_cancel",
		Body:         body,
	}

	pubCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// 发送到延迟队列，消息将在队列 TTL 到期后自动转发到取消队列
	err = Channel.PublishWithContext(
		pubCtx,
		OrderDelayExchange, // 延迟交换机
		"order.delay",      // 延迟路由键
		false,
		false,
		publishing,
	)

	if err != nil {
		klog.CtxErrorf(ctx, "Failed to publish cancel delay message: %v", err)
		return err
	}

	klog.CtxInfof(ctx, "Order cancel delay message published: order_id=%s, will cancel after %v",
		orderID, OrderCancelDelayTime)
	return nil
}

// PublishOrderCancelDelayWithCustomTTL 发布自定义延迟时间的取消消息（用于测试）
func PublishOrderCancelDelayWithCustomTTL(ctx context.Context, orderID string, userID uint64, delay time.Duration) error {
	msg := &OrderCancelMessage{
		OrderID:   orderID,
		UserID:    userID,
		CreatedAt: time.Now().Unix(),
		Reason:    "timeout",
	}

	body, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	// 使用消息级别的 TTL（会覆盖队列级别的 TTL，取较小值）
	publishing := amqp.Publishing{
		ContentType:  "application/json",
		DeliveryMode: amqp.Persistent,
		Timestamp:    time.Now(),
		MessageId:    orderID + "_cancel",
		Expiration:   formatTTL(delay), // 消息级别 TTL
		Body:         body,
	}

	pubCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err = Channel.PublishWithContext(
		pubCtx,
		OrderDelayExchange,
		"order.delay",
		false,
		false,
		publishing,
	)

	if err != nil {
		klog.CtxErrorf(ctx, "Failed to publish cancel delay message: %v", err)
		return err
	}

	klog.CtxInfof(ctx, "Order cancel delay message published: order_id=%s, custom delay=%v", orderID, delay)
	return nil
}

// formatTTL 格式化 TTL 为字符串（毫秒）
func formatTTL(d time.Duration) string {
	return fmt.Sprintf("%d", d.Milliseconds())
}
