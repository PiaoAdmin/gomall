package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/PiaoAdmin/pmall/app/order/biz/model"
	"github.com/cloudwego/kitex/pkg/klog"
	amqp "github.com/rabbitmq/amqp091-go"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// 测试配置
var testConfig = struct {
	RabbitMQURL string
	MySQLDSN    string
	OrderQueue  string
	Exchange    string
	WorkerCount int
	BatchSize   int // 每批次发送的订单数
	TotalOrders int // 总订单数
}{
	RabbitMQURL: "amqp://admin:123456@localhost:5672/",
	MySQLDSN:    "root:123456@tcp(piaohost:3306)/p_order?charset=utf8mb4&parseTime=True&loc=Local",
	OrderQueue:  "order_test_queue",
	Exchange:    "order_test_exchange",
	WorkerCount: 5,
	BatchSize:   100,
	TotalOrders: 1000,
}

func init() {
	// 初始化日志
	klog.SetLevel(klog.LevelInfo)
}

// setupTestDB 初始化测试数据库
func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(mysql.Open(testConfig.MySQLDSN), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// 自动迁移表结构
	db.AutoMigrate(&model.Order{}, &model.OrderItem{})

	// 清空测试数据
	db.Exec("DELETE FROM order_items WHERE order_id LIKE 'test_%'")
	db.Exec("DELETE FROM orders WHERE order_id LIKE 'test_%'")

	return db
}

// setupTestRabbitMQ 初始化测试 RabbitMQ
func setupTestRabbitMQ(t *testing.T) (*amqp.Connection, *amqp.Channel) {
	conn, err := amqp.Dial(testConfig.RabbitMQURL)
	if err != nil {
		t.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		t.Fatalf("Failed to create channel: %v", err)
	}

	// 声明交换机和队列
	ch.ExchangeDeclare(testConfig.Exchange, "direct", true, false, false, false, nil)
	ch.QueueDeclare(testConfig.OrderQueue, true, false, false, false, nil)
	ch.QueueBind(testConfig.OrderQueue, "order.create", testConfig.Exchange, false, nil)

	// 清空队列
	ch.QueuePurge(testConfig.OrderQueue, false)

	return conn, ch
}

// TestMQProducerPerformance 测试消息队列生产者性能
func TestMQProducerPerformance(t *testing.T) {
	conn, ch := setupTestRabbitMQ(t)
	defer conn.Close()
	defer ch.Close()

	orderCount := testConfig.TotalOrders

	t.Run("单线程发送", func(t *testing.T) {
		start := time.Now()

		for i := 0; i < orderCount; i++ {
			msg := createTestOrderMessage(fmt.Sprintf("test_single_%d", i), uint64(i))
			body, _ := json.Marshal(msg)

			ch.PublishWithContext(
				context.Background(),
				testConfig.Exchange,
				"order.create",
				false, false,
				amqp.Publishing{
					ContentType:  "application/json",
					DeliveryMode: amqp.Persistent,
					Body:         body,
				},
			)
		}

		elapsed := time.Since(start)
		qps := float64(orderCount) / elapsed.Seconds()
		t.Logf("✓ 单线程发送 %d 条消息耗时: %v, QPS: %.2f", orderCount, elapsed, qps)
	})

	ch.QueuePurge(testConfig.OrderQueue, false)

	t.Run("并发发送", func(t *testing.T) {
		start := time.Now()
		var wg sync.WaitGroup
		var successCount int64
		concurrency := 10

		for i := 0; i < concurrency; i++ {
			wg.Add(1)
			go func(workerID int) {
				defer wg.Done()
				localCh, _ := conn.Channel()
				defer localCh.Close()

				for j := 0; j < orderCount/concurrency; j++ {
					msg := createTestOrderMessage(fmt.Sprintf("test_concurrent_%d_%d", workerID, j), uint64(workerID*1000+j))
					body, _ := json.Marshal(msg)

					err := localCh.PublishWithContext(
						context.Background(),
						testConfig.Exchange,
						"order.create",
						false, false,
						amqp.Publishing{
							ContentType:  "application/json",
							DeliveryMode: amqp.Persistent,
							Body:         body,
						},
					)
					if err == nil {
						atomic.AddInt64(&successCount, 1)
					}
				}
			}(i)
		}

		wg.Wait()
		elapsed := time.Since(start)
		qps := float64(successCount) / elapsed.Seconds()
		t.Logf("✓ 并发发送 %d 条消息耗时: %v, QPS: %.2f", successCount, elapsed, qps)
	})
}

// TestSyncVsAsyncDBWrite 对比同步和异步写入数据库性能
func TestSyncVsAsyncDBWrite(t *testing.T) {
	db := setupTestDB(t)

	orderCount := 500

	// 同步写入测试
	t.Run("同步写入数据库", func(t *testing.T) {
		start := time.Now()
		var successCount int64

		for i := 0; i < orderCount; i++ {
			orderID := fmt.Sprintf("test_sync_%d_%d", time.Now().UnixNano(), i)
			order := &model.Order{
				OrderId: orderID,
				UserId:  uint64(i),
				Email:   fmt.Sprintf("test%d@example.com", i),
				Status:  model.OrderStatePlaced,
				ShippingAddress: model.Address{
					Name:          "Test User",
					StreetAddress: "123 Test St",
					City:          "Test City",
					ZipCode:       12345,
				},
			}

			err := db.Transaction(func(tx *gorm.DB) error {
				if err := tx.Create(order).Error; err != nil {
					return err
				}
				// 创建3个订单项
				for j := 0; j < 3; j++ {
					item := &model.OrderItem{
						OrderId:  orderID,
						SkuId:    uint64(j + 1),
						SkuName:  fmt.Sprintf("SKU-%d", j+1),
						Price:    99.99,
						Quantity: int32(j + 1),
					}
					if err := tx.Create(item).Error; err != nil {
						return err
					}
				}
				return nil
			})
			if err == nil {
				atomic.AddInt64(&successCount, 1)
			}
		}

		elapsed := time.Since(start)
		qps := float64(successCount) / elapsed.Seconds()
		avgLatency := elapsed.Milliseconds() / int64(orderCount)
		t.Logf("✓ 同步写入 %d 条订单耗时: %v, QPS: %.2f, 平均延迟: %dms", successCount, elapsed, qps, avgLatency)
	})

	// 模拟异步写入（批量写入）
	t.Run("异步批量写入数据库", func(t *testing.T) {
		start := time.Now()
		batchSize := 50

		orders := make([]*model.Order, 0, batchSize)
		items := make([]*model.OrderItem, 0, batchSize*3)

		for i := 0; i < orderCount; i++ {
			orderID := fmt.Sprintf("test_async_%d_%d", time.Now().UnixNano(), i)
			order := &model.Order{
				OrderId: orderID,
				UserId:  uint64(i),
				Email:   fmt.Sprintf("test%d@example.com", i),
				Status:  model.OrderStatePlaced,
				ShippingAddress: model.Address{
					Name:          "Test User",
					StreetAddress: "123 Test St",
					City:          "Test City",
					ZipCode:       12345,
				},
			}
			orders = append(orders, order)

			for j := 0; j < 3; j++ {
				items = append(items, &model.OrderItem{
					OrderId:  orderID,
					SkuId:    uint64(j + 1),
					SkuName:  fmt.Sprintf("SKU-%d", j+1),
					Price:    99.99,
					Quantity: int32(j + 1),
				})
			}

			// 批量写入
			if len(orders) >= batchSize {
				db.Transaction(func(tx *gorm.DB) error {
					tx.CreateInBatches(orders, batchSize)
					tx.CreateInBatches(items, batchSize*3)
					return nil
				})
				orders = orders[:0]
				items = items[:0]
			}
		}

		// 处理剩余
		if len(orders) > 0 {
			db.Transaction(func(tx *gorm.DB) error {
				tx.CreateInBatches(orders, len(orders))
				tx.CreateInBatches(items, len(items))
				return nil
			})
		}

		elapsed := time.Since(start)
		qps := float64(orderCount) / elapsed.Seconds()
		avgLatency := elapsed.Milliseconds() / int64(orderCount)
		t.Logf("✓ 异步批量写入 %d 条订单耗时: %v, QPS: %.2f, 平均延迟: %dms", orderCount, elapsed, qps, avgLatency)
	})
}

// TestDatabasePressure 数据库压力测试
func TestDatabasePressure(t *testing.T) {
	db := setupTestDB(t)

	testCases := []struct {
		name        string
		concurrency int
		ordersPerGo int
	}{
		{"低并发(5)", 5, 100},
		{"中并发(20)", 20, 50},
		{"高并发(50)", 50, 20},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			start := time.Now()
			var wg sync.WaitGroup
			var successCount, failCount int64

			for i := 0; i < tc.concurrency; i++ {
				wg.Add(1)
				go func(workerID int) {
					defer wg.Done()

					for j := 0; j < tc.ordersPerGo; j++ {
						orderID := fmt.Sprintf("test_pressure_%d_%d_%d", time.Now().UnixNano(), workerID, j)
						order := &model.Order{
							OrderId: orderID,
							UserId:  uint64(workerID*1000 + j),
							Email:   fmt.Sprintf("test%d_%d@example.com", workerID, j),
							Status:  model.OrderStatePlaced,
						}

						err := db.Transaction(func(tx *gorm.DB) error {
							if err := tx.Create(order).Error; err != nil {
								return err
							}
							item := &model.OrderItem{
								OrderId:  orderID,
								SkuId:    1,
								SkuName:  "Test SKU",
								Price:    99.99,
								Quantity: 1,
							}
							return tx.Create(item).Error
						})

						if err != nil {
							atomic.AddInt64(&failCount, 1)
						} else {
							atomic.AddInt64(&successCount, 1)
						}
					}
				}(i)
			}

			wg.Wait()
			elapsed := time.Since(start)
			total := successCount + failCount
			qps := float64(total) / elapsed.Seconds()
			successRate := float64(successCount) / float64(total) * 100

			t.Logf("✓ %s: 总请求=%d, 成功=%d, 失败=%d, 耗时=%v, QPS=%.2f, 成功率=%.1f%%",
				tc.name, total, successCount, failCount, elapsed, qps, successRate)
		})
	}
}

// TestEndToEndAsync 端到端异步测试
func TestEndToEndAsync(t *testing.T) {
	// 跳过如果没有配置好环境
	if testing.Short() {
		t.Skip("Skipping end-to-end test in short mode")
	}

	conn, ch := setupTestRabbitMQ(t)
	defer conn.Close()
	defer ch.Close()

	db := setupTestDB(t)

	orderCount := 100
	start := time.Now()

	// 1. 发送消息到队列
	t.Log("开始发送订单消息...")
	for i := 0; i < orderCount; i++ {
		msg := createTestOrderMessage(fmt.Sprintf("test_e2e_%d_%d", time.Now().UnixNano(), i), uint64(i))
		body, _ := json.Marshal(msg)

		ch.PublishWithContext(
			context.Background(),
			testConfig.Exchange,
			"order.create",
			false, false,
			amqp.Publishing{
				ContentType:  "application/json",
				DeliveryMode: amqp.Persistent,
				Body:         body,
			},
		)
	}
	sendElapsed := time.Since(start)
	t.Logf("✓ 发送 %d 条消息耗时: %v", orderCount, sendElapsed)

	// 2. 模拟消费者处理
	consumeStart := time.Now()
	var processedCount int64

	msgs, err := ch.Consume(testConfig.OrderQueue, "", false, false, false, false, nil)
	if err != nil {
		t.Fatalf("Failed to consume: %v", err)
	}

	done := make(chan bool)
	go func() {
		for msg := range msgs {
			var orderMsg OrderMessage
			if err := json.Unmarshal(msg.Body, &orderMsg); err != nil {
				msg.Reject(false)
				continue
			}

			// 写入数据库
			order := &model.Order{
				OrderId: orderMsg.OrderID,
				UserId:  orderMsg.UserID,
				Email:   orderMsg.Email,
				Status:  model.OrderStatePlaced,
			}

			err := db.Transaction(func(tx *gorm.DB) error {
				if err := tx.Create(order).Error; err != nil {
					return err
				}
				for _, item := range orderMsg.Items {
					oi := &model.OrderItem{
						OrderId:  orderMsg.OrderID,
						SkuId:    item.SkuID,
						SkuName:  item.SkuName,
						Price:    item.Price,
						Quantity: item.Quantity,
					}
					if err := tx.Create(oi).Error; err != nil {
						return err
					}
				}
				return nil
			})

			if err == nil {
				msg.Ack(false)
				atomic.AddInt64(&processedCount, 1)
			} else {
				msg.Reject(false)
			}

			if atomic.LoadInt64(&processedCount) >= int64(orderCount) {
				done <- true
				return
			}
		}
	}()

	select {
	case <-done:
	case <-time.After(30 * time.Second):
		t.Log("Timeout waiting for messages")
	}

	consumeElapsed := time.Since(consumeStart)
	totalElapsed := time.Since(start)

	t.Logf("✓ 消费处理 %d 条消息耗时: %v", processedCount, consumeElapsed)
	t.Logf("✓ 端到端总耗时: %v, 平均每条: %v", totalElapsed, totalElapsed/time.Duration(orderCount))
}

// createTestOrderMessage 创建测试订单消息
func createTestOrderMessage(orderID string, userID uint64) *OrderMessage {
	return &OrderMessage{
		OrderID:   orderID,
		UserID:    userID,
		Email:     fmt.Sprintf("user%d@test.com", userID),
		CreatedAt: time.Now().Unix(),
		Address: OrderAddress{
			Name:          "Test User",
			StreetAddress: "123 Test Street",
			City:          "Test City",
			ZipCode:       12345,
		},
		Items: []OrderMessageItem{
			{SkuID: 1, SkuName: "Test Product 1", Price: 99.99, Quantity: 2},
			{SkuID: 2, SkuName: "Test Product 2", Price: 49.99, Quantity: 1},
		},
	}
}

// BenchmarkMQPublish 基准测试：消息发布
func BenchmarkMQPublish(b *testing.B) {
	conn, err := amqp.Dial(testConfig.RabbitMQURL)
	if err != nil {
		b.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		b.Fatalf("Failed to create channel: %v", err)
	}
	defer ch.Close()

	ch.ExchangeDeclare(testConfig.Exchange, "direct", true, false, false, false, nil)

	msg := createTestOrderMessage("bench_order", 1)
	body, _ := json.Marshal(msg)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ch.PublishWithContext(
			context.Background(),
			testConfig.Exchange,
			"order.create",
			false, false,
			amqp.Publishing{
				ContentType:  "application/json",
				DeliveryMode: amqp.Persistent,
				Body:         body,
			},
		)
	}
}

// BenchmarkDBWrite 基准测试：数据库写入
func BenchmarkDBWrite(b *testing.B) {
	db, err := gorm.Open(mysql.Open(testConfig.MySQLDSN), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		b.Fatalf("Failed to connect: %v", err)
	}

	db.AutoMigrate(&model.Order{}, &model.OrderItem{})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		orderID := fmt.Sprintf("bench_%d_%d", time.Now().UnixNano(), i)
		order := &model.Order{
			OrderId: orderID,
			UserId:  uint64(i),
			Email:   "bench@test.com",
			Status:  model.OrderStatePlaced,
		}
		db.Create(order)
	}
}
