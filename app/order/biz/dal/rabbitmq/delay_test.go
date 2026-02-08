package rabbitmq

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/cloudwego/kitex/pkg/klog"
)

func init() {
	// 设置测试环境
	os.Setenv("GO_ENV", "test")
}

// TestDelayQueueFlow 测试延迟队列完整流程
// 1. 初始化连接
// 2. 发送延迟消息
// 3. 等待消息过期
// 4. 验证取消消费者收到消息
func TestDelayQueueFlow(t *testing.T) {
	Init()

	ctx := context.Background()

	// 测试订单 ID
	testOrderID := fmt.Sprintf("TEST_ORDER_%d", time.Now().UnixNano())
	testUserID := uint64(12345)

	t.Logf("=== 延迟队列测试开始 ===")
	t.Logf("测试订单 ID: %s", testOrderID)
	t.Logf("当前 TTL 设置: %v", OrderCancelDelayTime)

	// 发送延迟取消消息
	err := PublishOrderCancelDelay(ctx, testOrderID, testUserID)
	if err != nil {
		t.Fatalf("发送延迟消息失败: %v", err)
	}
	t.Logf("✓ 延迟取消消息已发送，将在 %v 后过期", OrderCancelDelayTime)

	// 启动取消消费者
	t.Log("启动取消消费者...")
	StartCancelConsumer(ctx)

	// 等待消息过期并被处理
	waitTime := OrderCancelDelayTime + 5*time.Second
	t.Logf("等待 %v 让消息过期并被消费...", waitTime)

	// 创建一个 channel 用于接收消费结果
	done := make(chan struct{})
	go func() {
		time.Sleep(waitTime)
		close(done)
	}()

	<-done
	t.Log("✓ 等待完成，检查日志确认消息是否被处理")

	// 停止消费者
	StopCancelConsumer()
	t.Log("✓ 消费者已停止")

	t.Log("=== 延迟队列测试完成 ===")
}

// TestDelayQueueManual 手动测试 - 发送消息后等待观察
func TestDelayQueueManual(t *testing.T) {
	Init()

	ctx := context.Background()

	t.Logf("当前延迟时间: %v", OrderCancelDelayTime)

	// 发送多个测试消息
	for i := 1; i <= 3; i++ {
		orderID := fmt.Sprintf("DELAY_TEST_%d_%d", time.Now().Unix(), i)
		err := PublishOrderCancelDelay(ctx, orderID, uint64(i))
		if err != nil {
			t.Errorf("发送消息 %d 失败: %v", i, err)
		} else {
			t.Logf("✓ 消息 %d 已发送: %s", i, orderID)
		}
	}

	t.Log("\n请在 RabbitMQ 管理界面观察:")
	t.Log("1. order_delay_queue 中应该有 3 条消息")
	t.Logf("2. %v 后消息将转移到 order_cancel_queue", OrderCancelDelayTime)
	t.Log("3. 访问 http://localhost:15672 (admin/123456)")
}

// TestCancelConsumerOnly 单独测试取消消费者
func TestCancelConsumerOnly(t *testing.T) {
	Init()

	ctx := context.Background()

	t.Log("启动取消消费者，等待消费消息...")
	t.Log("如果 order_cancel_queue 中有消息，将被处理")

	StartCancelConsumer(ctx)

	// 监听一段时间
	klog.Info("消费者运行中，按 Ctrl+C 停止或等待 60 秒自动结束...")

	select {
	case <-time.After(60 * time.Second):
		t.Log("60 秒已到，停止消费者")
	}

	StopCancelConsumer()
	t.Log("✓ 测试完成")
}
