package service

import (
	"context"
<<<<<<< HEAD
	// order "douyin-gomall/gomall/rpc_gen/kitex_gen/order"
=======
>>>>>>> b6e73c27fce12b01552c5334097a847176b8f26a
	order "github.com/PiaoAdmin/gomall/rpc_gen/kitex_gen/order"
)

type MarkOrderPaidService struct {
	ctx context.Context
<<<<<<< HEAD
}

// NewMarkOrderPaidService new MarkOrderPaidService
=======
} // NewMarkOrderPaidService new MarkOrderPaidService
>>>>>>> b6e73c27fce12b01552c5334097a847176b8f26a
func NewMarkOrderPaidService(ctx context.Context) *MarkOrderPaidService {
	return &MarkOrderPaidService{ctx: ctx}
}

// Run create note info
func (s *MarkOrderPaidService) Run(req *order.MarkOrderPaidReq) (resp *order.MarkOrderPaidResp, err error) {
	// Finish your business logic.

	return
}
