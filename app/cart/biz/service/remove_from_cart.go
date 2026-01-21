package service

import (
	"context"

	"github.com/PiaoAdmin/pmall/app/cart/biz/model"
	"github.com/PiaoAdmin/pmall/common/errs"
	cart "github.com/PiaoAdmin/pmall/rpc_gen/cart"
	"github.com/cloudwego/kitex/pkg/klog"
)

type RemoveFromCartService struct {
	ctx context.Context
}

func NewRemoveFromCartService(ctx context.Context) *RemoveFromCartService {
	return &RemoveFromCartService{ctx: ctx}
}

func (s *RemoveFromCartService) Run(req *cart.RemoveFromCartRequest) (*cart.RemoveFromCartResponse, error) {
	klog.CtxInfof(s.ctx, "RemoveFromCart: userID=%d, cartItemIDs=%v", req.UserId, req.CartItemIds)
	if len(req.CartItemIds) == 0 {
		return nil, errs.New(errs.ErrParam.Code, "cart_item_ids is required")
	}
	err := model.RemoveFromCart(s.ctx, req.UserId, req.CartItemIds)
	if err != nil {
		klog.CtxErrorf(s.ctx, "RemoveFromCart failed: %v", err)
		return nil, err
	}

	return &cart.RemoveFromCartResponse{
		Success: true,
	}, nil
}
