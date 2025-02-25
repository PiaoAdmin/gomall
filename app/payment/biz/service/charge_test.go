package service

import (
	"context"
	"testing"

	"github.com/PiaoAdmin/gomall/app/payment/biz/dal"
	payment "github.com/PiaoAdmin/gomall/rpc_gen/kitex_gen/payment"
	"github.com/joho/godotenv"
)

func TestCharge_Run(t *testing.T) {
	_ = godotenv.Load("../../.env")
	dal.Init()
	ctx := context.Background()
	s := NewChargeService(ctx)
	// init req and assert value

	req := &payment.ChargeReq{
		OrderId: 1894334578726805504,
		Amount:  100.0,
		CreditCard: &payment.CreditCardInfo{
			CreditCardNumber:          "4242424242424242",
			CreditCardCvv:             12,
			CreditCardExpirationYear:  2026,
			CreditCardExpirationMonth: 02,
		},
	}
	resp, err := s.Run(req)
	t.Logf("err: %v", err)
	t.Logf("resp: %v", resp)

	// todo: edit your unit test

}
