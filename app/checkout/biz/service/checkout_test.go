package service

import (
	"context"
	"testing"

	"github.com/PiaoAdmin/gomall/app/checkout/biz/dal"
	"github.com/PiaoAdmin/gomall/app/checkout/infra/rpc"
	checkout "github.com/PiaoAdmin/gomall/rpc_gen/kitex_gen/checkout"
	"github.com/PiaoAdmin/gomall/rpc_gen/kitex_gen/payment"
	"github.com/joho/godotenv"
)

func TestCheckout_Run(t *testing.T) {
	_ = godotenv.Load("../../.env")
	dal.Init()
	rpc.InitClient()
	ctx := context.Background()
	s := NewCheckoutService(ctx)
	// init req and assert value

	req := &checkout.CheckoutReq{

		CreditCard: &payment.CreditCardInfo{
			CreditCardNumber:          "1234567890123456",
			CreditCardCvv:             123,
			CreditCardExpirationYear:  2025,
			CreditCardExpirationMonth: 12,
		},
		Firstname: "Test",
		Lastname:  "Test",
		UserId:    1892469484459921408,
		Email:     "test@example.com",
		Address: &checkout.Address{
			StreetAddress: "Test Street",
			City:          "Test City",
			State:         "Test State",
			Country:       "Test Country",
			ZipCode:       123456,
		},
	}
	resp, err := s.Run(req)
	t.Logf("err: %v", err)
	t.Logf("resp: %v", resp)

	// todo: edit your unit test

}
