/*
 * @Author: liaosijie
 * @Date: 2025-02-18 22:28:50
 * @Last Modified by: liaosijie
 * @Last Modified time: 2025-02-18 23:30:55
 */

package service

import (
	"context"

	"github.com/PiaoAdmin/gomall/app/cart"

	// "douyin-gomall/gomall/app/order/biz/dal/mysql"
	// "douyin-gomall/gomall/app/order/biz/model"
	"github.com/PiaoAdmin/gomall/app/order/biz/dal/mysql"
	"github.com/PiaoAdmin/gomall/app/order/biz/model"

	// order "douyin-gomall/gomall/rpc_gen/kitex_gen/order"
	order "github.com/PiaoAdmin/gomall/rpc_gen/kitex_gen/order"
	"github.com/cloudwego/kitex/pkg/kerrors"
)

type ListOrderService struct {
	ctx context.Context
}

// NewListOrderService new ListOrderService
func NewListOrderService(ctx context.Context) *ListOrderService {
	return &ListOrderService{ctx: ctx}
}

// Run create note info
func (s *ListOrderService) Run(req *order.ListOrderReq) (resp *order.ListOrderResp, err error) {
	// Finish your business logic.
	model.ListOrder(s.ctx,mysql.DB, req.UserId)
	if err != nil {
		return nil, kerrors.NewBizStatusError(500001,err.Error())
	}

	var orders []*order.Order
	for _, v := range orders {
		var items []*order.OrderItem
		for _, oi := range v.OrderItems {
			items = append(items, &order.OrderItem{
				Item:&cart.CartItem{
					ProductId:oi.ProductId,
					Quantity:oi.Quantity,
				},
				Cost:oi.Cost,
			})
		}
		orders = append(orders, &order.Order{
			OrderId:		v.OrderId,
			UserId:			v.UserId,
			UserCurrency:	v.UserCurrency,
			Email:			v.Consignee.Email,
			Address:		&order.Address{
				StreetAddress: v.Consignee.StreetAddress,
				Country:	   v.Consignee.Country,
				City:          v.Consignee.City,
				State:         v.Consignee.State,
				ZipCode:       v.Consignee.ZipCode,
			},
			Items: items,
		})
	}
	resp = &order.ListOrderResp{
		Orders: orders,
	}
	return
}
