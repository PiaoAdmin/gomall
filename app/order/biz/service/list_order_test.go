package service

import (
	"context"
	"testing"

	"github.com/PiaoAdmin/gomall/app/order/biz/dal"
	order "github.com/PiaoAdmin/gomall/rpc_gen/kitex_gen/order"
	"github.com/joho/godotenv"
)

func TestListOrder_Run(t *testing.T) {
	_ = godotenv.Load("../../.env")
	dal.Init()
	ctx := context.Background()
	s := NewListOrderService(ctx)
	// init req and assert value

	req := &order.ListOrderReq{
		UserId: 1892469484459921408,
	}
	resp, err := s.Run(req)
	t.Logf("err: %v", err)
	t.Logf("resp: %v", resp)

	// todo: edit your unit test

}
