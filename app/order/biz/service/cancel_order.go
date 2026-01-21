package service

import (
	"context"

	"github.com/PiaoAdmin/pmall/app/order/biz/dal/mysql"
	"github.com/PiaoAdmin/pmall/app/order/biz/model"
	"github.com/PiaoAdmin/pmall/common/errs"
	order "github.com/PiaoAdmin/pmall/rpc_gen/order"
)

type CancelOrderService struct {
	ctx context.Context
}

func NewCancelOrderService(ctx context.Context) *CancelOrderService {
	return &CancelOrderService{ctx: ctx}
}

func (s *CancelOrderService) Run(req *order.CancelOrderReq) (*order.CancelOrderResp, error) {
	if req == nil || req.OrderId == "" {
		return nil, errs.New(errs.ErrParam.Code, "order_id empty")
	}

	var ord model.Order
	if err := mysql.DB.Where("order_id = ?", req.OrderId).First(&ord).Error; err != nil {
		return nil, errs.New(errs.ErrRecordNotFound.Code, err.Error())
	}

	if err := mysql.DB.Model(&model.Order{}).Where("order_id = ?", req.OrderId).Update("status", model.OrderStateCanceled).Error; err != nil {
		return nil, errs.New(errs.ErrInternal.Code, "cancel order failed: "+err.Error())
	}

	return &order.CancelOrderResp{Success: true}, nil
}
