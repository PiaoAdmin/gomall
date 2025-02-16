/*
 * @Author: liaosijie
 * @Date: 2025-02-16 10:15:18
 * @Last Modified by:   liaosijie
 * @Last Modified time: 2025-02-16 13:15:18
 */

package order

import (
	"context"
	"database/sql"
	"gomall/biz/dal/queries/order"
	"gomall/biz/model/order"
	"testing"
)

func TestPlaceOrder(t *testing.T) {
    // 初始化测试环境
    db, err := sql.Open("mysql", "root:jinitaimei114514@tcp(39.103.237.155:10112)/gomall")
    if err != nil {
        t.Fatalf("failed to connect to database: %v", err)
    }
    defer db.Close()

    query := order.NewOrderQuery(db)
    service := NewOrderService(query)

    // 创建一个订单
    o := &order.Order{
        UserID:      1,
        UserCurrency: "USD",
        // ... 其他字段
    }

    // 调用PlaceOrder方法
    err = service.PlaceOrder(context.Background(), o)
    if err != nil {
        t.Fatalf("failed to place order: %v", err)
    }

    // 验证订单是否正确创建
    placedOrder, err := query.GetOrderByID(context.Background(), o.OrderID)
    if err != nil {
        t.Fatalf("failed to get order: %v", err)
    }

    if placedOrder.UserID != o.UserID {
        t.Errorf("expected user ID %d, got %d", o.UserID, placedOrder.UserID)
    }
    // ... 其他字段验证
}