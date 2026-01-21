package service

import (
	"context"

	apiOrder "github.com/PiaoAdmin/pmall/app/api/biz/model/api/order"
	"github.com/PiaoAdmin/pmall/app/api/rpc"
	orderrpc "github.com/PiaoAdmin/pmall/rpc_gen/order"
	"github.com/cloudwego/hertz/pkg/app"
)

type CancelOrderService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewCancelOrderService(ctx context.Context, c *app.RequestContext) *CancelOrderService {
	return &CancelOrderService{RequestContext: c, Context: ctx}
}

func (s *CancelOrderService) Run(req *apiOrder.CancelOrderReq) (resp *apiOrder.CancelOrderResp, err error) {
	rpcReq := &orderrpc.CancelOrderReq{OrderId: req.OrderId}
	rpcResp, err := rpc.OrderClient.CancelOrder(s.Context, rpcReq)
	if err != nil {
		return nil, err
	}
	return &apiOrder.CancelOrderResp{Success: rpcResp.Success}, nil
}
