/*
 * @Author: liaosijie
 * @Date: 2021-02-16 10:13:51
 * @Last Modified by: liaosijie
 * @Last Modified time: 2025-02-16 17:53:58
 */

package order

import (
	"context"
	"database/sql"
	"gomall/biz/model/order"
	"gomall/utils"
)
 
 type OrderQuery struct {
     db *sql.DB
 }
 
 func NewOrderQuery(db *sql.DB) *OrderQuery {
     return &OrderQuery{db: db}
 }
 
 // PlaceOrder 创建订单
 func (q *OrderQuery) PlaceOrder(ctx context.Context, o *order.Order) error {
     tx, err := q.db.Begin()
     if err != nil {
         return err
     }
     defer tx.Rollback()
 
     o.ID = utils.CreateId(1) // 确保每个节点ID唯一
     result, err := tx.ExecContext(ctx, `
         INSERT INTO orders (id, user_id, user_currency, address_id, email, created_at, status)
         VALUES (?, ?, ?, ?, ?, ?, ?)
     `, o.ID, o.UserID, o.UserCurrency, o.AddressID, o.Email, o.CreatedAt, o.Status)
     if err != nil {
         return err
     }
 
     for _, item := range o.OrderItems {
         item.ID = utils.CreateId(1) // 确保每个节点ID唯一
         _, err := tx.ExecContext(ctx, `
             INSERT INTO order_items (id, order_id, product_id, quantity, price, created_at, updated_at, deleted_at)
             VALUES (?, ?, ?, ?, ?, ?, ?, ?)
         `, item.ID, o.ID, item.ProductID, item.Quantity, item.Price, item.CreatedAt, item.UpdatedAt, item.DeletedAt)
         if err != nil {
             return err
         }
     }
 
     return tx.Commit()
 }
 
 // ListOrder 查询订单列表
 func (q *OrderQuery) ListOrder(ctx context.Context, userId int64) ([]*order.Order, error) {
     rows, err := q.db.QueryContext(ctx, `
         SELECT id, user_id, user_currency, address_id, email, created_at, updated_at, deleted_at, status
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
         err := rows.Scan(&o.ID, &o.UserID, &o.UserCurrency, &o.AddressID, &o.Email, &o.CreatedAt, &o.UpdatedAt, &o.DeletedAt, &o.Status)
         if err != nil {
             return nil, err
         }
 
         // 获取订单项
         items, err := q.GetOrderItems(ctx, o.ID)
         if err != nil {
             return nil, err
         }
         o.OrderItems = items
 
         orders = append(orders, &o)
     }
 
     return orders, nil
 }
 
 func (q *OrderQuery) GetOrderItems(ctx context.Context, orderID int64) ([]*order.OrderItem, error) {
     rows, err := q.db.QueryContext(ctx, `
         SELECT id, order_id, product_id, quantity, price, created_at, updated_at, deleted_at
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
         err := rows.Scan(&item.ID, &item.OrderID, &item.ProductID, &item.Quantity, &item.Price, &item.CreatedAt, &item.UpdatedAt, &item.DeletedAt)
         if err != nil {
             return nil, err
         }
         items = append(items, &item)
     }
 
     return items, nil
 }
 
 // MarkOrderPaid 标记订单为已支付
 func (q *OrderQuery) MarkOrderPaid(ctx context.Context, orderId int64) error {
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