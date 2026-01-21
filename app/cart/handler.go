package main

import (
	"context"

	"github.com/PiaoAdmin/pmall/app/cart/biz/service"
	cart "github.com/PiaoAdmin/pmall/rpc_gen/cart"
)

// CartServiceImpl implements the last service interface defined in the IDL.
type CartServiceImpl struct{}

// AddToCart implements the CartServiceImpl interface.
func (s *CartServiceImpl) AddToCart(ctx context.Context, req *cart.AddToCartRequest) (resp *cart.AddToCartResponse, err error) {
	return service.NewAddToCartService(ctx).Run(req)
}

// RemoveFromCart implements the CartServiceImpl interface.
func (s *CartServiceImpl) RemoveFromCart(ctx context.Context, req *cart.RemoveFromCartRequest) (resp *cart.RemoveFromCartResponse, err error) {
	return service.NewRemoveFromCartService(ctx).Run(req)
}

// GetCartDetails implements the CartServiceImpl interface.
func (s *CartServiceImpl) GetCartDetails(ctx context.Context, req *cart.GetCartDetailsRequest) (resp *cart.GetCartDetailsResponse, err error) {
	return service.NewGetCartDetailsService(ctx).Run(req)
}

// ClearCart implements the CartServiceImpl interface.
func (s *CartServiceImpl) ClearCart(ctx context.Context, req *cart.ClearCartRequest) (resp *cart.ClearCartResponse, err error) {
	return service.NewClearCartService(ctx).Run(req)
}
