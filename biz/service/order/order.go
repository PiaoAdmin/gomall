/*
 * @Author: liaosijie
 * @Date: 2025-02-16 10:15:31
 * @Last Modified by: liaosijie
 * @Last Modified time: 2025-02-16 17:52:32
 */

package order

import (
	"context"
	"gomall/biz/dal/queries/order"
)
 
 type OrderService struct {
     query *order.OrderQuery
 }
 
 func NewOrderService(query *order.OrderQuery) *OrderService {
     return &OrderService{query: query}
 }
 
 // PlaceOrder 创建订单
 func (s *OrderService) PlaceOrder(ctx context.Context, o *order.Order) error {
     return s.query.PlaceOrder(ctx, o)
 }
 
 // ListOrder 查询订单列表
 func (s *OrderService) ListOrder(ctx context.Context, userId int64) ([]*order.Order, error) {
     return s.query.ListOrder(ctx, userId)
 }
 
 // MarkOrderPaid 标记订单为已支付
 func (s *OrderService) MarkOrderPaid(ctx context.Context, orderId int64) error {
     return s.query.MarkOrderPaid(ctx, orderId)
 }