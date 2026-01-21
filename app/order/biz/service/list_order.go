package service

import (
	"context"
	"fmt"

	"github.com/PiaoAdmin/pmall/app/order/biz/dal/mysql"
	"github.com/PiaoAdmin/pmall/app/order/biz/model"
	"github.com/PiaoAdmin/pmall/common/errs"
	order "github.com/PiaoAdmin/pmall/rpc_gen/order"
)

type ListOrderService struct {
	ctx context.Context
}

func NewListOrderService(ctx context.Context) *ListOrderService {
	return &ListOrderService{ctx: ctx}
}

func (s *ListOrderService) Run(req *order.ListOrderReq) (*order.ListOrderResp, error) {
	if req == nil || req.UserId == 0 {
		return nil, errs.New(errs.ErrParam.Code, "user_id empty")
	}

	var orders []model.Order
	if err := mysql.DB.Preload("Items").Where("user_id = ?", req.UserId).Order("created_at desc").Find(&orders).Error; err != nil {
		return nil, errs.New(errs.ErrInternal.Code, "list orders failed: "+err.Error())
	}

	protoOrders := make([]*order.Order, 0, len(orders))
	for _, o := range orders {
		po := &order.Order{
			OrderId:   o.OrderId,
			UserId:    o.UserId,
			Email:     o.Email,
			Status:    o.Status,
			CreatedAt: int32(o.CreatedAt.Unix()),
		}
		items := make([]*order.CartItem, 0, len(o.Items))
		for _, it := range o.Items {
			items = append(items, &order.CartItem{
				SkuId:    it.SkuId,
				Quantity: it.Quantity,
				SkuName:  it.SkuName,
				Price:    fmt.Sprintf("%.2f", it.Price),
			})
		}
		po.Items = items
		protoOrders = append(protoOrders, po)
	}

	return &order.ListOrderResp{Orders: protoOrders}, nil
}
