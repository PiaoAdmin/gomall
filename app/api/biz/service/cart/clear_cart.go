package service

import (
	"context"

	apiCart "github.com/PiaoAdmin/pmall/app/api/biz/model/api/cart"
	"github.com/PiaoAdmin/pmall/app/api/md/jwt"
	"github.com/PiaoAdmin/pmall/app/api/rpc"
	"github.com/PiaoAdmin/pmall/rpc_gen/cart"
	"github.com/cloudwego/hertz/pkg/app"
)

type ClearCartService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewClearCartService(ctx context.Context, c *app.RequestContext) *ClearCartService {
	return &ClearCartService{
		RequestContext: c,
		Context:        ctx,
	}
}

func (s *ClearCartService) Run(req *apiCart.ClearCartReq) (resp *apiCart.ClearCartResp, err error) {
	// 从 JWT 获取用户 ID
	claims := jwt.ExtractClaims(s.Context, s.RequestContext)
	userID := uint64(claims[jwt.JwtMiddleware.IdentityKey].(float64))

	// 调用购物车 RPC 服务
	rpcResp, err := rpc.CartClient.ClearCart(s.Context, &cart.ClearCartRequest{
		UserId: userID,
	})
	if err != nil {
		return nil, err
	}

	resp = &apiCart.ClearCartResp{
		Success: rpcResp.Success,
	}

	return resp, nil
}
