package service

import (
	"context"

	"github.com/PiaoAdmin/pmall/app/product/biz/dal/mysql"
	"github.com/PiaoAdmin/pmall/app/product/biz/dal/redis"
	"github.com/PiaoAdmin/pmall/app/product/biz/model"
	"github.com/PiaoAdmin/pmall/app/product/biz/utils"
	"github.com/PiaoAdmin/pmall/common/errs"
	product "github.com/PiaoAdmin/pmall/rpc_gen/product"
	"github.com/cloudwego/kitex/pkg/klog"
)

type GetProductDetailService struct {
	ctx context.Context
}

func NewGetProductDetailService(ctx context.Context) *GetProductDetailService {
	return &GetProductDetailService{ctx: ctx}
}

func (s *GetProductDetailService) Run(req *product.GetProductDetailRequest) (*product.GetProductDetailResponse, error) {
	if req.Id == 0 {
		return nil, errs.New(errs.ErrParam.Code, "product id is required")
	}

	// 优先从缓存获取商品详情
	cached, err := redis.GetProductDetailFromCache(s.ctx, req.Id)
	if err != nil {
		klog.Warnf("Failed to get product detail from cache: %v", err)
	}
	if cached != nil {
		klog.Debugf("Product %d cache hit", req.Id)
		return s.convertCachedToResponse(cached), nil
	}

	// 缓存未命中，从数据库获取
	klog.Debugf("Product %d cache miss, fetching from database", req.Id)
	spu, err := model.GetSPUByID(s.ctx, mysql.DB, req.Id)
	if err != nil {
		return nil, errs.New(errs.ErrRecordNotFound.Code, "product not found")
	}

	skus, err := model.GetSKUsBySpuID(s.ctx, mysql.DB, req.Id)
	if err != nil {
		return nil, errs.New(errs.ErrInternal.Code, "get skus failed: "+err.Error())
	}

	var category *model.ProductCategory
	var protoCategory *product.Category
	if spu.CategoryID > 0 {
		category, err = model.GetCategoryByID(s.ctx, mysql.DB, spu.CategoryID)
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

	var brand *model.ProductBrand
	var protoBrand *product.Brand
	if spu.BrandID > 0 {
		brand, err = model.GetBrandByID(s.ctx, mysql.DB, spu.BrandID)
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

	var detail *model.ProductDetail
	var protoDetail *product.ProductDetail
	detail, err = model.GetProductDetailBySpuID(s.ctx, req.Id)
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

	// 异步设置缓存
	go func() {
		cachedDetail := s.buildCacheData(spu, skus, category, brand, detail)
		if err := redis.SetProductDetailCache(s.ctx, req.Id, cachedDetail); err != nil {
			klog.Warnf("Failed to set product detail cache: %v", err)
		}
	}()

	return &product.GetProductDetailResponse{
		Spu:      protoSPU,
		Skus:     protoSKUs,
		Category: protoCategory,
		Brand:    protoBrand,
		Detail:   protoDetail,
	}, nil
}

// convertCachedToResponse 将缓存数据转换为响应
func (s *GetProductDetailService) convertCachedToResponse(cached *redis.CachedProductDetail) *product.GetProductDetailResponse {
	var protoSPU *product.ProductSPU
	if cached.SPU != nil {
		protoSPU = &product.ProductSPU{
			Id:            cached.SPU.ID,
			BrandId:       cached.SPU.BrandID,
			CategoryId:    cached.SPU.CategoryID,
			Name:          cached.SPU.Name,
			SubTitle:      cached.SPU.SubTitle,
			MainImage:     cached.SPU.MainImage,
			PublishStatus: int32(cached.SPU.PublishStatus),
			VerifyStatus:  int32(cached.SPU.VerifyStatus),
			SaleCount:     int32(cached.SPU.SaleCount),
			Sort:          int32(cached.SPU.Sort),
			ServiceBits:   cached.SPU.ServiceBits,
			Version:       int32(cached.SPU.Version),
		}
	}

	protoSKUs := make([]*product.ProductSKU, 0, len(cached.SKUs))
	for _, sku := range cached.SKUs {
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

	var protoCategory *product.Category
	if cached.Category != nil {
		protoCategory = &product.Category{
			Id:       cached.Category.ID,
			ParentId: cached.Category.ParentID,
			Name:     cached.Category.Name,
			Level:    int32(cached.Category.Level),
			Icon:     cached.Category.Icon,
			Unit:     cached.Category.Unit,
			Sort:     int32(cached.Category.Sort),
		}
	}

	var protoBrand *product.Brand
	if cached.Brand != nil {
		protoBrand = &product.Brand{
			Id:          cached.Brand.ID,
			Name:        cached.Brand.Name,
			FirstLetter: cached.Brand.FirstLetter,
			Logo:        cached.Brand.Logo,
			Sort:        int32(cached.Brand.Sort),
		}
	}

	var protoDetail *product.ProductDetail
	if cached.Detail != nil {
		protoDetail = &product.ProductDetail{
			Description:   cached.Detail.Description,
			Images:        cached.Detail.Images,
			Videos:        cached.Detail.Videos,
			MarketTagJson: cached.Detail.MarketTagJSON,
			TechTagJson:   cached.Detail.TechTagJSON,
			FaqJson:       cached.Detail.FaqJSON,
		}
	}

	return &product.GetProductDetailResponse{
		Spu:      protoSPU,
		Skus:     protoSKUs,
		Category: protoCategory,
		Brand:    protoBrand,
		Detail:   protoDetail,
	}
}

// buildCacheData 构建缓存数据
func (s *GetProductDetailService) buildCacheData(spu *model.ProductSPU, skus []*model.ProductSKU, category *model.ProductCategory, brand *model.ProductBrand, detail *model.ProductDetail) *redis.CachedProductDetail {
	cached := &redis.CachedProductDetail{}

	if spu != nil {
		cached.SPU = &redis.CachedSPU{
			ID:            spu.ID,
			BrandID:       spu.BrandID,
			CategoryID:    spu.CategoryID,
			Name:          spu.Name,
			SubTitle:      spu.SubTitle,
			MainImage:     spu.MainImage,
			PublishStatus: spu.PublishStatus,
			VerifyStatus:  spu.VerifyStatus,
			LowPrice:      spu.LowPrice,
			HighPrice:     spu.HighPrice,
			SaleCount:     spu.SaleCount,
			Sort:          spu.Sort,
			ServiceBits:   spu.ServiceBits,
			Version:       spu.Version,
		}
	}

	cached.SKUs = make([]*redis.CachedSKU, 0, len(skus))
	for _, sku := range skus {
		cached.SKUs = append(cached.SKUs, &redis.CachedSKU{
			ID:          sku.ID,
			SpuID:       sku.SpuID,
			SkuCode:     sku.SkuCode,
			Name:        sku.Name,
			SubTitle:    sku.SubTitle,
			MainImage:   sku.MainImage,
			Price:       sku.Price,
			MarketPrice: sku.MarketPrice,
			Stock:       sku.Stock,
			LockStock:   sku.LockStock,
			SkuSpecData: sku.SkuSpecData,
			Version:     sku.Version,
		})
	}

	if category != nil {
		cached.Category = &redis.CachedCategory{
			ID:       category.ID,
			ParentID: category.ParentID,
			Name:     category.Name,
			Level:    category.Level,
			Icon:     category.Icon,
			Unit:     category.Unit,
			Sort:     category.Sort,
		}
	}

	if brand != nil {
		cached.Brand = &redis.CachedBrand{
			ID:          brand.ID,
			Name:        brand.Name,
			FirstLetter: brand.FirstLetter,
			Logo:        brand.Logo,
			Sort:        brand.Sort,
		}
	}

	if detail != nil {
		cached.Detail = &redis.CachedDetail{
			Description:   detail.Description,
			Images:        detail.Images,
			Videos:        detail.Videos,
			MarketTagJSON: detail.MarketTagJSON,
			TechTagJSON:   detail.TechTagJSON,
			FaqJSON:       detail.FaqJSON,
		}
	}

	return cached
}
