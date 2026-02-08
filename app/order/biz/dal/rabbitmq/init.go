package rabbitmq

import (
	"sync"
	"time"

	"github.com/PiaoAdmin/pmall/app/order/conf"
	"github.com/cloudwego/kitex/pkg/klog"
	amqp "github.com/rabbitmq/amqp091-go"
)

// 延迟队列相关常量
const (
	// TODO:订单超时取消延迟时间（测试用：30秒，生产环境改回 30 * time.Minute）
	OrderCancelDelayTime = 30 * time.Second
	// 延迟队列名称
	OrderDelayQueue = "order_delay_queue"
	// 取消订单队列名称
	OrderCancelQueue = "order_cancel_queue"
	// 延迟交换机名称
	OrderDelayExchange = "order_delay_exchange"
)

var (
	Connection *amqp.Connection
	Channel    *amqp.Channel
	once       sync.Once
)

// Init 初始化 RabbitMQ 连接
func Init() {
	once.Do(func() {
		cfg := conf.GetConf().RabbitMQ
		var err error

		// 带重试的连接逻辑
		for i := 0; i < 5; i++ {
			Connection, err = amqp.Dial(cfg.URL)
			if err == nil {
				break
			}
			klog.Warnf("Failed to connect to RabbitMQ, retrying in 2s... (attempt %d/5): %v", i+1, err)
			time.Sleep(2 * time.Second)
		}
		if err != nil {
			klog.Errorf("Failed to connect to RabbitMQ after 5 attempts: %v", err)
			panic(err)
		}

		Channel, err = Connection.Channel()
		if err != nil {
			klog.Errorf("Failed to create RabbitMQ channel: %v", err)
			panic(err)
		}

		// 设置 QoS，限制未确认消息数量
		if err := Channel.Qos(cfg.PrefetchCount, 0, false); err != nil {
			klog.Errorf("Failed to set QoS: %v", err)
			panic(err)
		}

		// ========== 订单创建队列 ==========
		// 声明交换机
		if err := Channel.ExchangeDeclare(
			cfg.OrderExchange, // 交换机名称
			"direct",          // 类型
			true,              // 持久化
			false,             // 自动删除
			false,             // 内部
			false,             // 不等待
			nil,               // 参数
		); err != nil {
			klog.Errorf("Failed to declare exchange: %v", err)
			panic(err)
		}

		// 声明队列
		_, err = Channel.QueueDeclare(
			cfg.OrderQueue, // 队列名称
			true,           // 持久化
			false,          // 自动删除
			false,          // 排他
			false,          // 不等待
			amqp.Table{
				"x-dead-letter-exchange":    cfg.OrderExchange + ".dlx",
				"x-dead-letter-routing-key": "order.dead",
			},
		)
		if err != nil {
			klog.Errorf("Failed to declare queue: %v", err)
			panic(err)
		}

		// 绑定队列到交换机
		if err := Channel.QueueBind(
			cfg.OrderQueue,    // 队列名称
			"order.create",    // 路由键
			cfg.OrderExchange, // 交换机名称
			false,             // 不等待
			nil,               // 参数
		); err != nil {
			klog.Errorf("Failed to bind queue: %v", err)
			panic(err)
		}

		// 声明死信交换机和队列
		if err := Channel.ExchangeDeclare(
			cfg.OrderExchange+".dlx",
			"direct",
			true,
			false,
			false,
			false,
			nil,
		); err != nil {
			klog.Warnf("Failed to declare DLX exchange: %v", err)
		}

		_, err = Channel.QueueDeclare(
			cfg.OrderQueue+".dlq",
			true,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			klog.Warnf("Failed to declare DLQ: %v", err)
		}

		if err := Channel.QueueBind(
			cfg.OrderQueue+".dlq",
			"order.dead",
			cfg.OrderExchange+".dlx",
			false,
			nil,
		); err != nil {
			klog.Warnf("Failed to bind DLQ: %v", err)
		}

		// ========== 订单延迟取消队列（30分钟超时） ==========
		initDelayQueue(cfg)

		klog.Info("Successfully connected to RabbitMQ")
	})
}

// initDelayQueue 初始化延迟队列（用于订单超时取消）
func initDelayQueue(cfg conf.RabbitMQ) {
	// 1. 声明延迟交换机（用于接收延迟消息）
	if err := Channel.ExchangeDeclare(
		OrderDelayExchange, // 延迟交换机
		"direct",
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		klog.Errorf("Failed to declare delay exchange: %v", err)
		return
	}

	// 2. 声明延迟队列（消息在此队列等待，TTL到期后转发到取消队列）
	// 设置 x-message-ttl 使消息在队列中等待指定时间
	// 设置 x-dead-letter-exchange 指定消息过期后转发的交换机
	_, err := Channel.QueueDeclare(
		OrderDelayQueue,
		true,
		false,
		false,
		false,
		amqp.Table{
			"x-message-ttl":             int32(OrderCancelDelayTime.Milliseconds()), // 30分钟 TTL
			"x-dead-letter-exchange":    OrderDelayExchange,                         // 过期后转发到同一个交换机
			"x-dead-letter-routing-key": "order.cancel",                             // 使用取消路由键
		},
	)
	if err != nil {
		klog.Errorf("Failed to declare delay queue: %v", err)
		return
	}

	// 3. 绑定延迟队列到交换机（接收延迟消息）
	if err := Channel.QueueBind(
		OrderDelayQueue,
		"order.delay", // 延迟消息路由键
		OrderDelayExchange,
		false,
		nil,
	); err != nil {
		klog.Errorf("Failed to bind delay queue: %v", err)
		return
	}

	// 4. 声明取消订单队列（接收延迟队列过期的消息）
	_, err = Channel.QueueDeclare(
		OrderCancelQueue,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		klog.Errorf("Failed to declare cancel queue: %v", err)
		return
	}

	// 5. 绑定取消队列到交换机
	if err := Channel.QueueBind(
		OrderCancelQueue,
		"order.cancel", // 取消消息路由键
		OrderDelayExchange,
		false,
		nil,
	); err != nil {
		klog.Errorf("Failed to bind cancel queue: %v", err)
		return
	}

	klog.Infof("Delay queue initialized: messages will expire after %v", OrderCancelDelayTime)
}

// Close 关闭 RabbitMQ 连接
func Close() {
	if Channel != nil {
		Channel.Close()
	}
	if Connection != nil {
		Connection.Close()
	}
	klog.Info("RabbitMQ connection closed")
}

// IsConnected 检查连接是否正常
func IsConnected() bool {
	return Connection != nil && !Connection.IsClosed()
}
