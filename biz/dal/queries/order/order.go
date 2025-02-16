/*
 * @Author: liaosijie
 * @Date: 2021-02-16 10:13:51
 * @Last Modified by: liaosijie
 * @Last Modified time: 2025-02-16 13:14:35
 */

package order

import (
	"context"
	"database/sql"
	"gomall/biz/model/order"
)

type OrderQuery struct {
    db *sql.DB
}

func NewOrderQuery(db *sql.DB) *OrderQuery {
    return &OrderQuery{db: db}
}

// PlaceOrder 创建订单
func (q *OrderQuery) PlaceOrder(ctx context.Context, o *order.Order) error {
    // 开始事务
    tx, err := q.db.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback()

    // 插入订单
    result, err := tx.ExecContext(ctx, `
        INSERT INTO orders (user_id, user_currency, address_id, email, status)
        VALUES (?, ?, ?, ?, ?)
    `, o.UserID, o.UserCurrency, o.AddressID, o.Email, o.Status)
    if err != nil {
        return err
    }

    orderID, err := result.LastInsertId()
    if err != nil {
        return err
    }

    // 插入订单项
    for _, item := range o.OrderItems {
        _, err := tx.ExecContext(ctx, `
            INSERT INTO order_items (order_id, product_id, quantity, price)
            VALUES (?, ?, ?, ?)
        `, orderID, item.ProductID, item.Quantity, item.Price)
        if err != nil {
            return err
        }
    }

    // 提交事务
    return tx.Commit()
}

// ListOrder 查询订单列表
func (q *OrderQuery) ListOrder(ctx context.Context, userId int) ([]*order.Order, error) {
    rows, err := q.db.QueryContext(ctx, `
        SELECT id, user_id, user_currency, address_id, email, created_at, status
        FROM orders
        WHERE user_id = ?
    `, userId)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var orders []*order.Order
    for rows.Next() {
        var o order.Order
        err := rows.Scan(&o.ID, &o.UserID, &o.UserCurrency, &o.AddressID, &o.Email, &o.CreatedAt, &o.Status)
        if err != nil {
            return nil, err
        }

        // 获取订单项
        items, err := q.getOrderItems(ctx, o.ID)
        if err != nil {
            return nil, err
        }
        o.OrderItems = items

        orders = append(orders, &o)
    }

    return orders, nil
}

func (q *OrderQuery) getOrderItems(ctx context.Context, orderID string) ([]*order.OrderItem, error) {
    rows, err := q.db.QueryContext(ctx, `
        SELECT id, order_id, product_id, quantity, price
        FROM order_items
        WHERE order_id = ?
    `, orderID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var items []*order.OrderItem
    for rows.Next() {
        var item order.OrderItem
        err := rows.Scan(&item.ID, &item.OrderID, &item.ProductID, &item.Quantity, &item.Price)
        if err != nil {
            return nil, err
        }
        items = append(items, &item)
    }

    return items, nil
}


// MarkOrderPaid 标记订单为已支付
func (q *OrderQuery) MarkOrderPaid(ctx context.Context, orderId string) error {
    _, err := q.db.ExecContext(ctx, `
        UPDATE orders
        SET status = 'paid'
        WHERE id = ?
    `, orderId)
    if err != nil {
        return err
    }

    return nil
}
