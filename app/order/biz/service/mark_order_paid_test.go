package service

import (
	"context"
<<<<<<< HEAD
	"testing"

	// order "douyin-gomall/gomall/rpc_gen/kitex_gen/order"
	order "github.com/PiaoAdmin/gomall/rpc_gen/kitex_gen/order"
=======
	order "github.com/PiaoAdmin/gomall/rpc_gen/kitex_gen/order"
	"testing"
>>>>>>> b6e73c27fce12b01552c5334097a847176b8f26a
)

func TestMarkOrderPaid_Run(t *testing.T) {
	ctx := context.Background()
	s := NewMarkOrderPaidService(ctx)
	// init req and assert value

	req := &order.MarkOrderPaidReq{}
	resp, err := s.Run(req)
	t.Logf("err: %v", err)
	t.Logf("resp: %v", resp)

	// todo: edit your unit test

}
