package service

import (
	"context"

	apiCart "github.com/PiaoAdmin/pmall/app/api/biz/model/api/cart"
	"github.com/PiaoAdmin/pmall/app/api/md/jwt"
	"github.com/PiaoAdmin/pmall/app/api/rpc"
	"github.com/PiaoAdmin/pmall/common/errs"
	"github.com/PiaoAdmin/pmall/rpc_gen/cart"
	"github.com/cloudwego/hertz/pkg/app"
)

type RemoveFromCartService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewRemoveFromCartService(ctx context.Context, c *app.RequestContext) *RemoveFromCartService {
	return &RemoveFromCartService{
		RequestContext: c,
		Context:        ctx,
	}
}

func (s *RemoveFromCartService) Run(req *apiCart.RemoveFromCartReq) (resp *apiCart.RemoveFromCartResp, err error) {
	// 从 JWT 获取用户 ID
	claims := jwt.ExtractClaims(s.Context, s.RequestContext)
	userID := uint64(claims[jwt.JwtMiddleware.IdentityKey].(float64))
	if len(req.SkuIds) == 0 {
		return nil, errs.New(errs.ErrParam.Code, "len sku_ids is 0")
	}
	// 调用购物车 RPC 服务
	rpcResp, err := rpc.CartClient.RemoveFromCart(s.Context, &cart.RemoveFromCartRequest{
		UserId:      userID,
		CartItemIds: req.SkuIds, // 这里 CartItemIds 实际是 SKU IDs
	})
	if err != nil {
		return nil, err
	}

	resp = &apiCart.RemoveFromCartResp{
		Success: rpcResp.Success,
	}

	return resp, nil
}
