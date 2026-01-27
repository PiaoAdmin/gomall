package main

import (
	"context"

	"github.com/PiaoAdmin/pmall/app/payment/biz/service"
	payment "github.com/PiaoAdmin/pmall/rpc_gen/payment"
)

// PaymentServiceImpl implements the last service interface defined in the IDL.
type PaymentServiceImpl struct{}

// Pay implements the PaymentServiceImpl interface.
func (s *PaymentServiceImpl) Pay(ctx context.Context, req *payment.PayRequest) (resp *payment.PayResponse, err error) {
	return service.NewPayService(ctx).Run(req)
}
