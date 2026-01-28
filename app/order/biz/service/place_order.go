package service

import (
	"context"
	"fmt"
	"strconv"

	"github.com/PiaoAdmin/pmall/app/order/biz/dal/mysql"
	"github.com/PiaoAdmin/pmall/app/order/biz/model"
	"github.com/PiaoAdmin/pmall/app/order/biz/rpc"
	"github.com/PiaoAdmin/pmall/common/errs"
	"github.com/PiaoAdmin/pmall/common/uniqueid"
	order "github.com/PiaoAdmin/pmall/rpc_gen/order"
	"github.com/PiaoAdmin/pmall/rpc_gen/product"
	"gorm.io/gorm"
)

type PlaceOrderService struct {
	ctx context.Context
}

func NewPlaceOrderService(ctx context.Context) *PlaceOrderService {
	return &PlaceOrderService{ctx: ctx}
}

func (s *PlaceOrderService) Run(req *order.PlaceOrderReq) (*order.PlaceOrderResp, error) {
	if req == nil || req.UserId == 0 {
		return nil, errs.New(errs.ErrParam.Code, "invalid request")
	}
	if len(req.Items) == 0 {
		return nil, errs.New(errs.ErrParam.Code, "items empty")
	}

	newOrderId := fmt.Sprintf("%d", uniqueid.GenId())

	deductItems := make([]*product.SkuDeductItem, 0, len(req.Items))
	for _, it := range req.Items {
		if it == nil || it.GetSkuId() == 0 || it.GetQuantity() <= 0 {
			return nil, errs.New(errs.ErrParam.Code, "invalid sku_id or quantity")
		}
		deductItems = append(deductItems, &product.SkuDeductItem{
			SkuId: it.GetSkuId(),
			Count: it.GetQuantity(),
		})
	}

	if _, err := rpc.ProductClient.DeductStock(s.ctx, &product.DeductStockRequest{
		OrderSn: newOrderId,
		Items:   deductItems,
	}); err != nil {
		return nil, errs.New(errs.ErrInternal.Code, "deduct stock failed: "+err.Error())
	}

	o := &model.Order{
		OrderId: newOrderId,
		UserId:  req.UserId,
		Email:   req.Email,
		Status:  model.OrderStatePlaced,
	}
	if req.ShippingAddress != nil {
		o.ShippingAddress = model.Address{
			Name:          req.ShippingAddress.GetName(),
			StreetAddress: req.ShippingAddress.GetStreetAddress(),
			City:          req.ShippingAddress.GetCity(),
			ZipCode:       req.ShippingAddress.GetZipCode(),
		}
	}

	err := mysql.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(o).Error; err != nil {
			return err
		}
		for _, it := range req.Items {
			price := 0.0
			if it != nil && it.GetPrice() != "" {
				if p, perr := strconv.ParseFloat(it.GetPrice(), 64); perr == nil {
					price = p
				}
			}
			qty := it.GetQuantity()
			oi := &model.OrderItem{
				OrderId:  newOrderId,
				SkuId:    it.GetSkuId(),
				SkuName:  it.GetSkuName(),
				Price:    price,
				Quantity: qty,
			}
			if err := tx.Create(oi).Error; err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		if _, relErr := rpc.ProductClient.ReleaseStock(s.ctx, &product.ReleaseStockRequest{
			OrderSn: newOrderId,
			Items:   deductItems,
		}); relErr != nil {
			return nil, errs.New(errs.ErrInternal.Code, "place order failed and release stock failed: "+relErr.Error())
		}
		if e, ok := err.(*errs.Error); ok {
			return nil, e
		}
		return nil, errs.New(errs.ErrInternal.Code, "place order failed: "+err.Error())
	}
	// TODO: 使用消息队列定时取消订单
	return &order.PlaceOrderResp{
		Order: &order.OrderResult{OrderId: newOrderId},
	}, nil
}
