package main

import (
	"context"

	"github.com/PiaoAdmin/pmall/app/checkout/biz/service"
	checkout "github.com/PiaoAdmin/pmall/rpc_gen/checkout"
)

// CheckoutServiceImpl implements the last service interface defined in the IDL.
type CheckoutServiceImpl struct{}

// Checkout implements the CheckoutServiceImpl interface.
func (s *CheckoutServiceImpl) Checkout(ctx context.Context, req *checkout.CheckoutRequest) (resp *checkout.CheckoutResponse, err error) {
	return service.NewCheckoutService(ctx).Run(req)
}
