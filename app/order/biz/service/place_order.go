/*
 * @Author: liaosijie
 * @Date: 2025-02-18 17:09:37
 * @Last Modified by: liaosijie
 * @Last Modified time: 2025-02-18 17:13:15
 */

package service

import (
	"context"
	order "douyin-gomall/gomall/rpc_gen/kitex_gen/order"

	"github.com/cloudwego/kitex/pkg/kerrors"
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
	}
	

	return
}
