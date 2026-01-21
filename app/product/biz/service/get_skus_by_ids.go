package service

import (
	"context"

	"github.com/PiaoAdmin/pmall/app/product/biz/dal/mysql"
	"github.com/PiaoAdmin/pmall/app/product/biz/model"
	"github.com/PiaoAdmin/pmall/app/product/biz/utils"
	product "github.com/PiaoAdmin/pmall/rpc_gen/product"
	"github.com/cloudwego/kitex/pkg/klog"
)

type GetSkusByIdsService struct {
	ctx context.Context
}

func NewGetSkusByIdsService(ctx context.Context) *GetSkusByIdsService {
	return &GetSkusByIdsService{ctx: ctx}
}

func (s *GetSkusByIdsService) Run(req *product.GetSkusByIdsRequest) (*product.GetSkusByIdsResponse, error) {
	klog.CtxInfof(s.ctx, "GetSkusByIds: skuIDs=%v", req.SkuIds)

	if len(req.SkuIds) == 0 {
		return &product.GetSkusByIdsResponse{
			Skus: make(map[uint64]*product.ProductSKU),
		}, nil
	}

	// 从数据库批量查询 SKU
	var skus []*model.ProductSKU
	err := mysql.DB.Where("id IN ?", req.SkuIds).Find(&skus).Error
	if err != nil {
		klog.CtxErrorf(s.ctx, "GetSkusByIds query failed: %v", err)
		return nil, err
	}

	// 转换为 map 结构
	skuMap := make(map[uint64]*product.ProductSKU)
	for _, sku := range skus {
		skuMap[sku.ID] = &product.ProductSKU{
			Id:          sku.ID,
			SpuId:       sku.SpuID,
			SkuCode:     sku.SkuCode,
			Name:        sku.Name,
			SubTitle:    sku.SubTitle,
			MainImage:   sku.MainImage,
			Price:       utils.PriceToString(sku.Price),
			MarketPrice: utils.PriceToString(sku.MarketPrice),
			Stock:       int32(sku.Stock),
			LockStock:   int32(sku.LockStock),
			SkuSpecData: sku.SkuSpecData,
			Version:     int32(sku.Version),
		}
	}

	klog.CtxInfof(s.ctx, "GetSkusByIds found %d SKUs", len(skuMap))
	return &product.GetSkusByIdsResponse{
		Skus: skuMap,
	}, nil
}
