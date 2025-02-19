/*
 * @Author: liaosijie
 * @Date: 2025-02-18 17:09:37
 * @Last Modified by: liaosijie
 * @Last Modified time: 2025-02-18 23:37:37
 */

package service

import (
	"context"
	// "douyin-gomall/gomall/app/order/biz/dal/mysql"
	// "douyin-gomall/gomall/app/order/biz/model"
	// order "douyin-gomall/gomall/rpc_gen/kitex_gen/order"
	"github.com/PiaoAdmin/gomall/app/order/biz/dal/mysql"
	"github.com/PiaoAdmin/gomall/app/order/biz/model"
	order "github.com/PiaoAdmin/gomall/rpc_gen/kitex_gen/order"

	"github.com/cloudwego/kitex/pkg/kerrors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PlaceOrderService struct {
	ctx context.Context
}

// NewPlaceOrderService new PlaceOrderService
func NewPlaceOrderService(ctx context.Context) *PlaceOrderService {
	return &PlaceOrderService{ctx: ctx}
}

// Run create note info
func (s *PlaceOrderService) Run(req *order.PlaceOrderReq) (resp *order.PlaceOrderResp, err error) {
	// Finish your business logic.
	if len(req.Items)==0{
		err = kerrors.NewBizStatusError(50001,"items is empty")
		return
	}
	err = mysql.DB.Transaction(func(tx *gorm.DB) error {
		orderId, _ := uuid.NewUUID()

		o := &model.Order{
			OrderId:     orderId.String(),
			UserId:      req.UserId,
			UserCurrency: req.Currency,
			Consignee: model.Consignee{
				Email: req.Email,
			},
		}
		if req.Address!= nil {
			a := req.Address
			o.Consignee.StreetAddress = a.StreetAddress
			o.Consignee.City = a.City
			o.Consignee.State = a.State
			o.Consignee.Country = a.Country
		}
		if err := tx.Create(&o).Error; err!= nil {
			return err
		}
		var items []model.OrderItem
		for _, v := range req.Items {
			items = append(items, model.OrderItem{
				OrderIdRefer:     orderId.String(),
				ProductID:   v.Item.ProductId,
				Quantity:    v.Item.Quantity,
				Cost:        v.Item.Cost,
			})
		}

		if err := tx.Create(items).Error; err!= nil {
			return err
		}

		resp = &order.PlaceOrderResp{
			Order: &order.OrderResult{
				OrderId: orderId.String(),
			},
		}

		return nil
	})

	return
}
