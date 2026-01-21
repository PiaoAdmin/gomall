package service

import (
	"context"

	apiCart "github.com/PiaoAdmin/pmall/app/api/biz/model/api/cart"
	"github.com/PiaoAdmin/pmall/app/api/md/jwt"
	"github.com/PiaoAdmin/pmall/app/api/rpc"
	"github.com/PiaoAdmin/pmall/rpc_gen/cart"
	"github.com/cloudwego/hertz/pkg/app"
)

type AddToCartService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewAddToCartService(ctx context.Context, c *app.RequestContext) *AddToCartService {
	return &AddToCartService{
		RequestContext: c,
		Context:        ctx,
	}
}

func (s *AddToCartService) Run(req *apiCart.AddToCartReq) (resp *apiCart.AddToCartResp, err error) {
	claims := jwt.ExtractClaims(s.Context, s.RequestContext)
	userID := uint64(claims[jwt.JwtMiddleware.IdentityKey].(float64))

	// 调用购物车 RPC 服务
	rpcResp, err := rpc.CartClient.AddToCart(s.Context, &cart.AddToCartRequest{
		UserId:   userID,
		SkuId:    req.SkuId,
		Quantity: req.Quantity,
	})
	if err != nil {
		return nil, err
	}

	// 转换为 API 响应
	resp = &apiCart.AddToCartResp{
		CartItem: convertCartItem(rpcResp.CartItem),
	}

	return resp, nil
}

// convertCartItem 转换 RPC CartItem 为 API CartItem
func convertCartItem(item *cart.CartItem) *apiCart.CartItem {
	if item == nil {
		return nil
	}

	return &apiCart.CartItem{
		SkuId:       item.SkuId,
		Quantity:    item.Quantity,
		SkuName:     item.SkuName,
		SkuImage:    item.SkuImage,
		Price:       item.Price,
		MarketPrice: item.MarketPrice,
		Stock:       item.Stock,
		SpuId:       item.SpuId,
		SpuName:     item.SpuName,
		SkuSpecData: item.SkuSpecData,
	}
}
