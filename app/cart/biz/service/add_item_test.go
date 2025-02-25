package service

import (
	"context"
	"testing"

	"github.com/PiaoAdmin/gomall/app/cart/biz/dal"
	"github.com/PiaoAdmin/gomall/app/cart/rpc"
	cart "github.com/PiaoAdmin/gomall/rpc_gen/kitex_gen/cart"
	"github.com/joho/godotenv"
)

func TestAddItem_Run(t *testing.T) {
	_ = godotenv.Load("../../.env")
	rpc.InitClient()
	dal.Init()
	ctx := context.Background()
	s := NewAddItemService(ctx)
	// init req and assert value

	req := &cart.AddItemReq{
		UserId: 1892469484459921408,
		Item: &cart.CartItem{
			ProductId: 1894233990104858624,
			Quantity:  10,
		},
	}
	resp, err := s.Run(req)
	t.Logf("err: %v", err)
	t.Logf("resp: %v", resp)

	// todo: edit your unit test

}
