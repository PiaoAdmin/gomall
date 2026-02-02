package service

import (
	"context"

	"github.com/PiaoAdmin/pmall/app/product/biz/dal/mysql"
	"github.com/PiaoAdmin/pmall/app/product/biz/model"
	"github.com/PiaoAdmin/pmall/app/product/biz/utils"
	"github.com/PiaoAdmin/pmall/common/errs"
	product "github.com/PiaoAdmin/pmall/rpc_gen/product"
)

type CreateProductService struct {
	ctx context.Context
}

func NewCreateProductService(ctx context.Context) *CreateProductService {
	return &CreateProductService{ctx: ctx}
}

func (s *CreateProductService) Run(req *product.CreateProductRequest) (resp *product.CreateProductResponse, err error) {
	if req.Spu == nil || req.Skus == nil || req.Detail == nil {
		return nil, errs.New(errs.ErrParam.Code, "invalid request")
	}
	newSPU := &model.ProductSPU{
		Name:          req.Spu.Name,
		CategoryID:    req.Spu.CategoryId,
		BrandID:       req.Spu.BrandId,
		SubTitle:      req.Spu.SubTitle,
		MainImage:     req.Spu.MainImage,
		PublishStatus: 1,
		VerifyStatus:  1,
		ServiceBits:   req.Spu.ServiceBits,
	}
	newDetail := &model.ProductDetail{
		Description:   req.Detail.Description,
		Images:        req.Detail.Images,
		Videos:        req.Detail.Videos,
		MarketTagJSON: req.Detail.MarketTagJson,
		TechTagJSON:   req.Detail.TechTagJson,
		FaqJSON:       req.Detail.FaqJson,
	}
	newSKUs := make([]*model.ProductSKU, 0, len(req.Skus))
	for _, sku := range req.Skus {
		price, err := utils.PriceConvert(sku.Price)
		if err != nil {
			return nil, errs.New(errs.ErrParam.Code, "invalid sku price")
		}
		if price < 0 {
			return nil, errs.New(errs.ErrParam.Code, "sku price cannot be negative")
		}
		markerPrice, err := utils.PriceConvert(sku.MarketPrice)
		if err != nil {
			return nil, errs.New(errs.ErrParam.Code, "invalid sku market price")
		}
		if markerPrice < 0 {
			return nil, errs.New(errs.ErrParam.Code, "sku market price cannot be negative")
		}
		newSKU := &model.ProductSKU{
			SkuCode:       sku.SkuCode,
			Name:          sku.Name,
			SubTitle:      sku.SubTitle,
			MainImage:     sku.MainImage,
			Price:         price,
			MarketPrice:   markerPrice,
			Stock:         int(sku.Stock),
			LockStock:     0,
			SkuSpecData:   sku.SkuSpecData,
			Version:       1,
			PublishStatus: 1,
			VerifyStatus:  1,
		}
		newSKUs = append(newSKUs, newSKU)
	}
	spuid, err := model.CreateProductWithTransaction(s.ctx, mysql.DB, newSPU, newSKUs, newDetail)
	return &product.CreateProductResponse{SpuId: spuid}, err
}
