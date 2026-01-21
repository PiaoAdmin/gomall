package service

import (
	"context"

	apiCart "github.com/PiaoAdmin/pmall/app/api/biz/model/api/cart"
	"github.com/PiaoAdmin/pmall/app/api/md/jwt"
	"github.com/PiaoAdmin/pmall/app/api/rpc"
	"github.com/PiaoAdmin/pmall/rpc_gen/cart"
	"github.com/cloudwego/hertz/pkg/app"
)

type GetCartDetailsService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewGetCartDetailsService(ctx context.Context, c *app.RequestContext) *GetCartDetailsService {
	return &GetCartDetailsService{
		RequestContext: c,
		Context:        ctx,
	}
}

func (s *GetCartDetailsService) Run(req *apiCart.GetCartDetailsReq) (resp *apiCart.GetCartDetailsResp, err error) {
	// 从 JWT 获取用户 ID
	claims := jwt.ExtractClaims(s.Context, s.RequestContext)
	userID := uint64(claims[jwt.JwtMiddleware.IdentityKey].(float64))

	// 调用购物车 RPC 服务
	rpcResp, err := rpc.CartClient.GetCartDetails(s.Context, &cart.GetCartDetailsRequest{
		UserId: userID,
	})
	if err != nil {
		return nil, err
	}

	// 转换购物车项列表
	items := make([]*apiCart.CartItem, 0, len(rpcResp.Items))
	for _, item := range rpcResp.Items {
		items = append(items, convertCartItem(item))
	}

	resp = &apiCart.GetCartDetailsResp{
		Items:         items,
		TotalQuantity: rpcResp.TotalQuantity,
		TotalAmount:   rpcResp.TotalAmount,
	}

	return resp, nil
}
