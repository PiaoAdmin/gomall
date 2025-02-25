package service

import (
	"context"
	"testing"

	"github.com/PiaoAdmin/gomall/app/order/biz/dal"
	"github.com/PiaoAdmin/gomall/rpc_gen/kitex_gen/cart"
	order "github.com/PiaoAdmin/gomall/rpc_gen/kitex_gen/order"
	"github.com/joho/godotenv"
)

func TestPlaceOrder_Run(t *testing.T) {
	_ = godotenv.Load("../../.env")
	dal.Init()
	ctx := context.Background()
	s := NewPlaceOrderService(ctx)
	// init req and assert value

	req := &order.PlaceOrderReq{
		UserId: 1892469484459921408,
		Email:  "test@example.com",
		Address: &order.Address{
			StreetAddress: "Test Street",
			City:          "Test City",
			State:         "Test State",
			Country:       "Test Country",
			ZipCode:       123456,
		},
		OrderItems: []*order.OrderItem{
			{
				Item: &cart.CartItem{
					ProductId: 987654321,
					Quantity:  1,
				},
				Cost: 100.0,
			},
		},
	}
	resp, err := s.Run(req)
	t.Logf("err: %v", err)
	t.Logf("resp: %v", resp)

	// todo: edit your unit test

}
