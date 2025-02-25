<<<<<<< HEAD
/*
 * @Author: liaosijie
 * @Date: 2025-02-18 16:51:32
 * @Last Modified by:   liaosijie
 * @Last Modified time: 2025-02-18 16:51:32
 */

=======
>>>>>>> b6e73c27fce12b01552c5334097a847176b8f26a
package main

import (
	"context"
<<<<<<< HEAD
	// order "douyin-gomall/gomall/rpc_gen/kitex_gen/order"
	// "douyin-gomall/gomall/app/order/biz/service"
=======
>>>>>>> b6e73c27fce12b01552c5334097a847176b8f26a
	"github.com/PiaoAdmin/gomall/app/order/biz/service"
	order "github.com/PiaoAdmin/gomall/rpc_gen/kitex_gen/order"
)

// OrderServiceImpl implements the last service interface defined in the IDL.
type OrderServiceImpl struct{}

// PlaceOrder implements the OrderServiceImpl interface.
func (s *OrderServiceImpl) PlaceOrder(ctx context.Context, req *order.PlaceOrderReq) (resp *order.PlaceOrderResp, err error) {
	resp, err = service.NewPlaceOrderService(ctx).Run(req)

	return resp, err
}

// ListOrder implements the OrderServiceImpl interface.
func (s *OrderServiceImpl) ListOrder(ctx context.Context, req *order.ListOrderReq) (resp *order.ListOrderResp, err error) {
	resp, err = service.NewListOrderService(ctx).Run(req)

	return resp, err
}

// MarkOrderPaid implements the OrderServiceImpl interface.
func (s *OrderServiceImpl) MarkOrderPaid(ctx context.Context, req *order.MarkOrderPaidReq) (resp *order.MarkOrderPaidResp, err error) {
	resp, err = service.NewMarkOrderPaidService(ctx).Run(req)

	return resp, err
}
