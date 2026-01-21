package service

import (
	"context"

	"github.com/PiaoAdmin/pmall/app/cart/biz/model"
	"github.com/PiaoAdmin/pmall/app/cart/biz/rpc"
	cart "github.com/PiaoAdmin/pmall/rpc_gen/cart"
	"github.com/PiaoAdmin/pmall/rpc_gen/product"
	"github.com/cloudwego/kitex/pkg/klog"
)

type AddToCartService struct {
	ctx context.Context
}

func NewAddToCartService(ctx context.Context) *AddToCartService {
	return &AddToCartService{ctx: ctx}
}

func (s *AddToCartService) Run(req *cart.AddToCartRequest) (*cart.AddToCartResponse, error) {
	klog.CtxInfof(s.ctx, "AddToCart: userID=%d, skuID=%d, quantity=%d", req.UserId, req.SkuId, req.Quantity)

	quantity := req.Quantity
	if quantity <= 0 {
		quantity = 1
	}

	err := model.AddToCart(s.ctx, req.UserId, req.SkuId, quantity)
	if err != nil {
		klog.CtxErrorf(s.ctx, "AddToCart failed: %v", err)
		return nil, err
	}

	totalQuantity, err := model.GetCartItemQuantity(s.ctx, req.UserId, req.SkuId)
	if err != nil {
		klog.CtxErrorf(s.ctx, "GetCartItemQuantity failed: %v", err)
		return nil, err
	}

	cartItem := &cart.CartItem{
		UserId:   req.UserId,
		SkuId:    req.SkuId,
		Quantity: totalQuantity,
	}

	skuDetail, err := s.getSkuDetail(req.SkuId)
	if err != nil {
		klog.CtxWarnf(s.ctx, "GetSkuDetail failed: %v, return basic info only", err)
	} else if skuDetail != nil {
		cartItem.SkuName = skuDetail.Name
		cartItem.SkuImage = skuDetail.MainImage
		cartItem.Price = skuDetail.Price
		cartItem.MarketPrice = skuDetail.MarketPrice
		cartItem.Stock = skuDetail.Stock
		cartItem.SpuId = skuDetail.SpuId
		cartItem.SkuSpecData = skuDetail.SkuSpecData
	}

	return &cart.AddToCartResponse{
		Success:  true,
		CartItem: cartItem,
	}, nil
}

func (s *AddToCartService) getSkuDetail(skuID uint64) (*product.ProductSKU, error) {
	req := &product.GetSkusByIdsRequest{
		SkuIds: []uint64{skuID},
	}

	resp, err := rpc.ProductClient.GetSkusByIds(s.ctx, req)
	if err != nil {
		klog.CtxErrorf(s.ctx, "GetSkusByIds RPC call failed: %v", err)
		return nil, err
	}

	if sku, ok := resp.Skus[skuID]; ok {
		return sku, nil
	}

	return nil, nil
}
