# RabbitMQ 异步下单实现说明

## 概述

本次改进将订单服务的下单操作从**同步写入数据库**改为**异步消息队列处理**，使用 RabbitMQ 作为消息中间件。

## 架构图

```
┌─────────────┐     ┌──────────────┐     ┌─────────────┐     ┌──────────┐
│   API层     │────>│  Order 服务  │────>│  RabbitMQ   │────>│ Consumer │
│ (Checkout)  │     │ (Producer)   │     │   Queue     │     │ (Worker) │
└─────────────┘     └──────────────┘     └─────────────┘     └────┬─────┘
                           │                                      │
                           │ 扣减库存(同步)                        │ 写入DB(异步)
                           ▼                                      ▼
                    ┌──────────────┐                       ┌──────────┐
                    │ Product 服务 │                       │  MySQL   │
                    └──────────────┘                       └──────────┘
```

## 实现原理

### 1. 下单流程变化

**改进前（同步模式）：**
```
请求 -> 扣减库存 -> 写入订单表 -> 写入订单项表 -> 返回响应
        (RPC)       (事务)        (事务)
```
- 整个流程串行执行
- 数据库写入是阻塞操作
- 高并发下数据库成为瓶颈

**改进后（异步模式）：**
```
请求 -> 扣减库存 -> 发送MQ消息 -> 返回响应
        (RPC)       (非阻塞)      (立即)
                        │
                        ▼ (异步)
                   消费者处理 -> 写入订单表 -> 写入订单项表
                                  (事务)        (事务)
```
- 库存扣减保持同步，保证准确性
- 数据库写入变为异步，不阻塞请求
- 快速返回响应，提升用户体验

### 2. 核心组件

| 组件 | 文件 | 职责 |
|------|------|------|
| 初始化 | `rabbitmq/init.go` | 连接管理、队列声明、交换机配置 |
| 生产者 | `rabbitmq/producer.go` | 消息序列化、发布到队列 |
| 消费者 | `rabbitmq/consumer.go` | 消息消费、数据库写入、重试机制 |
| 消息结构 | `rabbitmq/producer.go` | OrderMessage 定义 |

### 3. 消息队列配置

```yaml
rabbitmq:
  url: "amqp://admin:123456@localhost:5672/"
  order_queue: "order_create_queue"      # 订单创建队列
  order_exchange: "order_exchange"       # 订单交换机
  prefetch_count: 10                     # 预取消息数
  worker_count: 5                        # 消费者工作线程数
```

### 4. 可靠性保障

1. **消息持久化**：消息设置 `DeliveryMode: amqp.Persistent`
2. **手动确认**：消费成功后手动 ACK，失败时 Reject
3. **死信队列**：超过重试次数的消息进入 DLQ
4. **幂等性检查**：消费者处理前检查订单是否已存在
5. **重试机制**：最多重试 3 次

## 为什么这么做

### 1. 削峰填谷

```
高峰期请求:    ████████████████████  (10000/s)
                    │
                    ▼ (MQ缓冲)
数据库处理:    ████████████          (平稳 500/s)
```

消息队列起到缓冲作用，将瞬时高峰请求平滑分散处理。

### 2. 解耦服务

- **API响应时间**：不再依赖数据库写入速度
- **独立扩展**：可以单独增加消费者数量
- **故障隔离**：数据库临时故障不影响下单

### 3. 提升吞吐量

**性能测试结果：**

| 指标 | 同步模式 | 异步模式 | 提升 |
|------|---------|---------|------|
| 消息发送 QPS | - | 220,258/s | - |
| 批量写入 QPS | ~50/s | 1,301/s | **26x** |
| 请求延迟 | 13ms | <1ms | **13x** |

### 4. 失败重试

同步模式下写入失败需要立即返回错误；异步模式下可以：
- 自动重试 3 次
- 失败消息进入死信队列，后续人工处理
- 不影响用户体验

## 好处总结

| 优势 | 说明 |
|------|------|
| ⚡ **低延迟** | 请求快速返回，不阻塞等待数据库 |
| 📈 **高吞吐** | 消息发送速度远高于数据库写入 |
| 🔄 **削峰填谷** | 平滑处理流量高峰 |
| 🔌 **解耦** | 服务间松耦合，独立部署和扩展 |
| 🛡️ **容错** | 失败重试、死信队列保证可靠性 |
| 🔍 **可追溯** | 消息可持久化，便于问题排查 |

## 压力测试

### 运行测试

```bash
# RabbitMQ 生产者性能测试
cd app/order
go test -v -run TestMQProducerPerformance ./biz/dal/rabbitmq/...

# 同步 vs 异步数据库写入对比
go test -v -run TestSyncVsAsyncDBWrite ./biz/dal/rabbitmq/...

# 数据库压力测试（不同并发级别）
go test -v -run TestDatabasePressure ./biz/dal/rabbitmq/...

# 基准测试
go test -bench=BenchmarkMQPublish -benchmem ./biz/dal/rabbitmq/...
```

### 测试结果

```
=== TestMQProducerPerformance ===
✓ 单线程发送 1000 条消息耗时: 4.54ms, QPS: 220,258
✓ 并发发送 1000 条消息耗时: 10.15ms, QPS: 98,527

=== BenchmarkMQPublish ===
BenchmarkMQPublish-28    113019    13976 ns/op    912 B/op    29 allocs/op
```

## 文件变更清单

```
app/order/
├── biz/
│   └── dal/
│       ├── init.go                    # 修改：添加 RabbitMQ 初始化
│       └── rabbitmq/                  # 新增目录
│           ├── init.go                # RabbitMQ 连接初始化 + 延迟队列
│           ├── producer.go            # 消息生产者（订单创建）
│           ├── consumer.go            # 消息消费者（订单写入DB）
│           ├── delay_producer.go      # 延迟消息生产者（订单取消）
│           ├── cancel_consumer.go     # 取消消费者（超时取消处理）
│           ├── errors.go              # 错误定义
│           ├── pressure_test.go       # 压力测试
│           └── README.md              # 说明文档
│   └── service/
│       └── place_order.go             # 修改：异步发消息 + 发送延迟取消
├── conf/
│   ├── conf.go                        # 修改：添加 RabbitMQ 配置结构
│   └── test/
│       └── conf.yaml                  # 修改：添加 RabbitMQ 配置
└── main.go                            # 修改：启动消费者 + 取消消费者
```

## 启动说明

### 1. 启动 RabbitMQ

```bash
docker run -d --name pmall-rabbitmq \
  -p 5672:5672 -p 15672:15672 \
  -e RABBITMQ_DEFAULT_USER=admin \
  -e RABBITMQ_DEFAULT_PASS=123456 \
  rabbitmq:3-management
```

访问管理界面：http://localhost:15672 (admin/123456)

### 2. 启动服务

```bash
cd scripts
./start_services.sh
```

### 3. 验证

1. 查看 RabbitMQ 管理界面的队列状态
2. 观察 `order_create_queue` 的消息流转
3. 检查数据库中的订单记录

## 延迟队列 - 订单超时自动取消

### 功能说明

当用户创建订单后，如果 **30 分钟内未支付**，系统会自动取消订单并释放库存。

### 实现原理

使用 RabbitMQ 的 **TTL + 死信队列** 模式实现延迟消息：

```
┌─────────────────┐   30分钟后过期   ┌─────────────────┐
│ order_delay_    │ ───────────────> │ order_cancel_   │
│ queue (TTL)     │    (死信转发)     │ queue           │
└─────────────────┘                  └────────┬────────┘
        ▲                                     │
        │                                     ▼
   订单创建时发送                         Cancel Consumer
   延迟取消消息                           检查订单状态
                                              │
                                    ┌─────────┴─────────┐
                                    ▼                   ▼
                               未支付则取消          已支付则忽略
                               释放库存
```

### 核心组件

| 组件 | 文件 | 职责 |
|------|------|------|
| 队列初始化 | `init.go` | 声明延迟队列、取消队列、交换机 |
| 延迟生产者 | `delay_producer.go` | 发送延迟取消消息 |
| 取消消费者 | `cancel_consumer.go` | 消费过期消息，执行订单取消 |

### 队列配置

```yaml
order_delay_queue:    # 延迟队列
  x-message-ttl: 1800000     # 30分钟 (毫秒)
  x-dead-letter-exchange: order_delay_exchange
  x-dead-letter-routing-key: order.cancel

order_cancel_queue:   # 取消处理队列
  # 接收从延迟队列转发的过期消息
```

### 工作流程

1. **订单创建成功** → 发送延迟消息到 `order_delay_queue`
2. **消息等待 30 分钟** → 消息在延迟队列中等待过期
3. **消息过期** → 自动转发到 `order_cancel_queue`（通过死信交换机）
4. **取消消费者处理**：
   - 查询订单当前状态
   - 如果是 `placed`（未支付）→ 更新为 `canceled` + 释放库存
   - 如果是 `paid`（已支付）→ 忽略，记录日志

### 消息结构

```go
type OrderCancelMessage struct {
    OrderID   string    `json:"order_id"`
    UserID    int64     `json:"user_id"`
    CreatedAt time.Time `json:"created_at"`
    Reason    string    `json:"reason"`
}
```

### 可靠性保障

1. **幂等性**：取消前检查订单状态，避免重复取消
2. **事务性**：订单状态更新和库存释放在同一事务中
3. **可追溯**：记录取消原因和时间
4. **手动确认**：处理成功后才 ACK 消息

## 注意事项

1. **库存扣减仍为同步**：保证库存准确性
2. **订单号立即返回**：用户可以用订单号查询
3. **最终一致性**：订单数据延迟写入（通常 < 100ms）
4. **消息确认**：消费成功才确认，保证不丢失
5. **幂等处理**：重复消费不会创建重复订单
6. **超时取消**：未支付订单 30 分钟后自动取消

## 监控建议

1. **队列积压监控**：队列消息数超过阈值告警
2. **消费延迟监控**：消息处理时间超过阈值告警
3. **死信队列监控**：DLQ 有消息需要人工处理
4. **消费者健康检查**：确保消费者在线
