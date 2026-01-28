package service

import (
	"context"
	"fmt"

	"github.com/PiaoAdmin/pmall/common/errs"
	"github.com/PiaoAdmin/pmall/common/uniqueid"
	payment "github.com/PiaoAdmin/pmall/rpc_gen/payment"
)

type PayService struct {
	ctx context.Context
}

func NewPayService(ctx context.Context) *PayService {
	return &PayService{ctx: ctx}
}

func (s *PayService) Run(req *payment.PayRequest) (*payment.PayResponse, error) {
	if req == nil || req.OrderId == "" || req.UserId == 0 {
		return nil, errs.New(errs.ErrParam.Code, "invalid request")
	}
	if req.Amount == "" {
		return nil, errs.New(errs.ErrParam.Code, "amount is required")
	}

	if req.CreditCard == "" {
		return nil, errs.New(errs.ErrParam.Code, "credit card is required")
	}
	if !isValidCreditCard(req.CreditCard) {
		return nil, errs.New(errs.ErrParam.Code, "invalid credit card")
	}

	tradeNo := fmt.Sprintf("%d", uniqueid.GenId())
	return &payment.PayResponse{
		Success: true,
		TradeNo: tradeNo,
	}, nil
}

func isValidCreditCard(card string) bool {
	// Strip spaces and dashes.
	clean := make([]byte, 0, len(card))
	for i := 0; i < len(card); i++ {
		ch := card[i]
		if ch >= '0' && ch <= '9' {
			clean = append(clean, ch)
		} else if ch == ' ' || ch == '-' {
			continue
		} else {
			return false
		}
	}
	if len(clean) < 13 || len(clean) > 19 {
		return false
	}
	// Luhn check.
	sum := 0
	double := false
	for i := len(clean) - 1; i >= 0; i-- {
		d := int(clean[i] - '0')
		if double {
			d *= 2
			if d > 9 {
				d -= 9
			}
		}
		sum += d
		double = !double
	}
	return sum%10 == 0
}
