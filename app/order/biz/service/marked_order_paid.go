package service

import (
	"context"

	"github.com/PiaoAdmin/pmall/app/order/biz/dal/mysql"
	"github.com/PiaoAdmin/pmall/app/order/biz/model"
	"github.com/PiaoAdmin/pmall/common/errs"
	order "github.com/PiaoAdmin/pmall/rpc_gen/order"
)

type MarkOrderPaidService struct {
	ctx context.Context
}

func NewMarkOrderPaidService(ctx context.Context) *MarkOrderPaidService {
	return &MarkOrderPaidService{ctx: ctx}
}

func (s *MarkOrderPaidService) Run(req *order.MarkOrderPaidReq) (*order.MarkOrderPaidResp, error) {
	if req == nil || req.OrderId == "" {
		return nil, errs.New(errs.ErrParam.Code, "order_id empty")
	}

	var ord model.Order
	if err := mysql.DB.Where("order_id = ?", req.OrderId).First(&ord).Error; err != nil {
		return nil, errs.New(errs.ErrRecordNotFound.Code, err.Error())
	}

	if err := mysql.DB.Model(&model.Order{}).Where("order_id = ?", req.OrderId).Update("status", model.OrderStatePaid).Error; err != nil {
		return nil, errs.New(errs.ErrInternal.Code, "mark paid failed: "+err.Error())
	}

	return &order.MarkOrderPaidResp{Success: true}, nil
}
