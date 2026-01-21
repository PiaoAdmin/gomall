package service

import (
	"context"
	"strconv"

	"github.com/PiaoAdmin/pmall/app/cart/biz/model"
	"github.com/PiaoAdmin/pmall/app/cart/biz/rpc"
	cart "github.com/PiaoAdmin/pmall/rpc_gen/cart"
	"github.com/PiaoAdmin/pmall/rpc_gen/product"
	"github.com/cloudwego/kitex/pkg/klog"
)

type GetCartDetailsService struct {
	ctx context.Context
}

func NewGetCartDetailsService(ctx context.Context) *GetCartDetailsService {
	return &GetCartDetailsService{ctx: ctx}
}

func (s *GetCartDetailsService) Run(req *cart.GetCartDetailsRequest) (*cart.GetCartDetailsResponse, error) {
	klog.CtxInfof(s.ctx, "GetCartDetails: userID=%d", req.UserId)

	cartItems, err := model.GetCartItems(s.ctx, req.UserId)
	if err != nil {
		klog.CtxErrorf(s.ctx, "GetCartItems failed: %v", err)
		return nil, err
	}

	if len(cartItems) == 0 {
		return &cart.GetCartDetailsResponse{
			Items:         []*cart.CartItem{},
			TotalQuantity: 0,
			TotalAmount:   "0.00",
		}, nil
	}

	// 获取所有 SKU IDs
	skuIDs := make([]uint64, 0, len(cartItems))
	for skuID := range cartItems {
		skuIDs = append(skuIDs, skuID)
	}

	// 批量查询商品信息
	skuDetails, err := s.getSkuDetails(skuIDs)
	if err != nil {
		klog.CtxWarnf(s.ctx, "GetSkuDetails failed: %v", err)
		// 即使查询失败，也返回基础购物车数据
		skuDetails = make(map[uint64]*product.ProductSKU)
	}

	// 组装返回数据
	var totalQuantity int32
	var totalAmount float64
	items := make([]*cart.CartItem, 0, len(cartItems))

	for skuID, quantity := range cartItems {
		totalQuantity += quantity

		item := &cart.CartItem{
			UserId:   req.UserId,
			SkuId:    skuID,
			Quantity: quantity,
		}

		// 填充 SKU 详情
		if sku, ok := skuDetails[skuID]; ok {
			item.SkuName = sku.Name
			item.SkuImage = sku.MainImage
			item.Price = sku.Price
			item.MarketPrice = sku.MarketPrice
			item.Stock = sku.Stock
			item.SpuId = sku.SpuId
			item.SkuSpecData = sku.SkuSpecData

			// 计算总金额
			if price, err := strconv.ParseFloat(sku.Price, 64); err == nil {
				totalAmount += price * float64(quantity)
			}
		}

		items = append(items, item)
	}

	return &cart.GetCartDetailsResponse{
		Items:         items,
		TotalQuantity: totalQuantity,
		TotalAmount:   strconv.FormatFloat(totalAmount, 'f', 2, 64),
	}, nil
}

func (s *GetCartDetailsService) getSkuDetails(skuIDs []uint64) (map[uint64]*product.ProductSKU, error) {
	if len(skuIDs) == 0 {
		return make(map[uint64]*product.ProductSKU), nil
	}

	req := &product.GetSkusByIdsRequest{
		SkuIds: skuIDs,
	}

	resp, err := rpc.ProductClient.GetSkusByIds(s.ctx, req)
	if err != nil {
		klog.CtxErrorf(s.ctx, "GetSkusByIds RPC call failed: %v", err)
		return nil, err
	}

	return resp.Skus, nil
}
