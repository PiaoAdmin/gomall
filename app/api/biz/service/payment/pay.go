package payment

import (
	"context"

	apiPayment "github.com/PiaoAdmin/pmall/app/api/biz/model/api/payment"
	"github.com/PiaoAdmin/pmall/app/api/md/jwt"
	"github.com/PiaoAdmin/pmall/app/api/rpc"
	paymentrpc "github.com/PiaoAdmin/pmall/rpc_gen/payment"
	"github.com/cloudwego/hertz/pkg/app"
)

type PayService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewPayService(ctx context.Context, c *app.RequestContext) *PayService {
	return &PayService{RequestContext: c, Context: ctx}
}

func (s *PayService) Run(req *apiPayment.PayReq) (*apiPayment.PayResp, error) {
	claims := jwt.ExtractClaims(s.Context, s.RequestContext)
	userID := uint64(claims[jwt.JwtMiddleware.IdentityKey].(float64))

	rpcResp, err := rpc.PaymentClient.Pay(s.Context, &paymentrpc.PayRequest{
		OrderId:    req.OrderId,
		UserId:     userID,
		Amount:     req.Amount,
		CreditCard: req.CreditCard,
	})
	if err != nil {
		return nil, err
	}

	return &apiPayment.PayResp{
		Success: rpcResp.Success,
		TradeNo: rpcResp.TradeNo,
	}, nil
}
