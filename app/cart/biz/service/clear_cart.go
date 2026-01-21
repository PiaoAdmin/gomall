package service

import (
	"context"

	"github.com/PiaoAdmin/pmall/app/cart/biz/model"
	cart "github.com/PiaoAdmin/pmall/rpc_gen/cart"
	"github.com/cloudwego/kitex/pkg/klog"
)

type ClearCartService struct {
	ctx context.Context
}

func NewClearCartService(ctx context.Context) *ClearCartService {
	return &ClearCartService{ctx: ctx}
}

func (s *ClearCartService) Run(req *cart.ClearCartRequest) (*cart.ClearCartResponse, error) {
	klog.CtxInfof(s.ctx, "ClearCart: userID=%d", req.UserId)

	err := model.ClearCart(s.ctx, req.UserId)
	if err != nil {
		klog.CtxErrorf(s.ctx, "ClearCart failed: %v", err)
		return nil, err
	}

	return &cart.ClearCartResponse{
		Success: true,
	}, nil
}
