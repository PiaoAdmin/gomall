package service

import (
	"context"

	"github.com/PiaoAdmin/pmall/app/order/biz/dal/mysql"
	"github.com/PiaoAdmin/pmall/app/order/biz/model"
	"github.com/PiaoAdmin/pmall/app/order/biz/rpc"
	"github.com/PiaoAdmin/pmall/common/errs"
	order "github.com/PiaoAdmin/pmall/rpc_gen/order"
	"github.com/PiaoAdmin/pmall/rpc_gen/product"
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
	if err := mysql.DB.Preload("Items").Where("order_id = ?", req.OrderId).First(&ord).Error; err != nil {
		return nil, errs.New(errs.ErrRecordNotFound.Code, err.Error())
	}

	if ord.Status == model.OrderStateCanceled {
		return &order.CancelOrderResp{Success: true}, nil
	}
	if ord.Status == model.OrderStatePaid {
		return nil, errs.New(errs.ErrParam.Code, "order already paid")
	}

	if len(ord.Items) > 0 {
		releaseItems := make([]*product.SkuDeductItem, 0, len(ord.Items))
		for _, it := range ord.Items {
			if it.SkuId == 0 || it.Quantity <= 0 {
				continue
			}
			releaseItems = append(releaseItems, &product.SkuDeductItem{
				SkuId: it.SkuId,
				Count: it.Quantity,
			})
		}
		if len(releaseItems) > 0 {
			if _, err := rpc.ProductClient.ReleaseStock(s.ctx, &product.ReleaseStockRequest{
				OrderSn: req.OrderId,
				Items:   releaseItems,
			}); err != nil {
				return nil, errs.New(errs.ErrInternal.Code, "release stock failed: "+err.Error())
			}
		}
	}

	if err := mysql.DB.Model(&model.Order{}).Where("order_id = ?", req.OrderId).Update("status", model.OrderStateCanceled).Error; err != nil {
		return nil, errs.New(errs.ErrInternal.Code, "cancel order failed: "+err.Error())
	}

	return &order.CancelOrderResp{Success: true}, nil
}
