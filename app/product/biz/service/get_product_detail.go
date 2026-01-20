package service

import (
	"context"

	"github.com/PiaoAdmin/pmall/app/product/biz/dal/mysql"
	"github.com/PiaoAdmin/pmall/app/product/biz/model"
	"github.com/PiaoAdmin/pmall/app/product/biz/utils"
	"github.com/PiaoAdmin/pmall/common/errs"
	product "github.com/PiaoAdmin/pmall/rpc_gen/product"
)

type GetProductDetailService struct {
	ctx context.Context
}

func NewGetProductDetailService(ctx context.Context) *GetProductDetailService {
	return &GetProductDetailService{ctx: ctx}
}

// TODO: cache result
func (s *GetProductDetailService) Run(req *product.GetProductDetailRequest) (*product.GetProductDetailResponse, error) {
	if req.Id == 0 {
		return nil, errs.New(errs.ErrParam.Code, "product id is required")
	}

	spu, err := model.GetSPUByID(s.ctx, mysql.DB, req.Id)
	if err != nil {
		return nil, errs.New(errs.ErrRecordNotFound.Code, "product not found")
	}

	skus, err := model.GetSKUsBySpuID(s.ctx, mysql.DB, req.Id)
	if err != nil {
		return nil, errs.New(errs.ErrInternal.Code, "get skus failed: "+err.Error())
	}

	var protoCategory *product.Category
	if spu.CategoryID > 0 {
		category, err := model.GetCategoryByID(s.ctx, mysql.DB, spu.CategoryID)
		if err == nil && category != nil {
			protoCategory = &product.Category{
				Id:       category.ID,
				ParentId: category.ParentID,
				Name:     category.Name,
				Level:    int32(category.Level),
				Icon:     category.Icon,
				Unit:     category.Unit,
				Sort:     int32(category.Sort),
			}
		}
	}

	var protoBrand *product.Brand
	if spu.BrandID > 0 {
		brand, err := model.GetBrandByID(s.ctx, mysql.DB, spu.BrandID)
		if err == nil && brand != nil {
			protoBrand = &product.Brand{
				Id:          brand.ID,
				Name:        brand.Name,
				FirstLetter: brand.FirstLetter,
				Logo:        brand.Logo,
				Sort:        int32(brand.Sort),
			}
		}
	}

	var protoDetail *product.ProductDetail
	detail, err := model.GetProductDetailBySpuID(s.ctx, req.Id)
	if err == nil && detail != nil {
		protoDetail = &product.ProductDetail{
			Description:   detail.Description,
			Images:        detail.Images,
			Videos:        detail.Videos,
			MarketTagJson: detail.MarketTagJSON,
			TechTagJson:   detail.TechTagJSON,
			FaqJson:       detail.FaqJSON,
		}
	}

	protoSPU := &product.ProductSPU{
		Id:            spu.ID,
		BrandId:       spu.BrandID,
		CategoryId:    spu.CategoryID,
		Name:          spu.Name,
		SubTitle:      spu.SubTitle,
		MainImage:     spu.MainImage,
		PublishStatus: int32(spu.PublishStatus),
		VerifyStatus:  int32(spu.VerifyStatus),
		SaleCount:     int32(spu.SaleCount),
		Sort:          int32(spu.Sort),
		ServiceBits:   spu.ServiceBits,
		Version:       int32(spu.Version),
	}

	protoSKUs := make([]*product.ProductSKU, 0, len(skus))
	for _, sku := range skus {
		protoSKUs = append(protoSKUs, &product.ProductSKU{
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
		})
	}

	return &product.GetProductDetailResponse{
		Spu:      protoSPU,
		Skus:     protoSKUs,
		Category: protoCategory,
		Brand:    protoBrand,
		Detail:   protoDetail,
	}, nil
}
