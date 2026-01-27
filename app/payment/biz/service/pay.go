package service

import (
	"context"
	"fmt"

	"github.com/PiaoAdmin/pmall/common/errs"
	"github.com/PiaoAdmin/pmall/common/uniqueid"
	payment "github.com/PiaoAdmin/pmall/rpc_gen/payment"
)

type PayService struct {
	ctx context.Context
}

func NewPayService(ctx context.Context) *PayService {
	return &PayService{ctx: ctx}
}

func (s *PayService) Run(req *payment.PayRequest) (*payment.PayResponse, error) {
	if req == nil || req.OrderId == "" || req.UserId == 0 {
		return nil, errs.New(errs.ErrParam.Code, "invalid request")
	}

	tradeNo := fmt.Sprintf("%d", uniqueid.GenId())
	return &payment.PayResponse{
		Success: true,
		TradeNo: tradeNo,
	}, nil
}
