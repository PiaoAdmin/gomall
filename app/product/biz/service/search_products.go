package service

import (
	"context"

	"github.com/PiaoAdmin/pmall/app/product/biz/dal/mysql"
	"github.com/PiaoAdmin/pmall/app/product/biz/model"
	"github.com/PiaoAdmin/pmall/app/product/biz/utils"
	"github.com/PiaoAdmin/pmall/common/errs"
	product "github.com/PiaoAdmin/pmall/rpc_gen/product"
)

type SearchProductsService struct {
	ctx context.Context
}

func NewSearchProductsService(ctx context.Context) *SearchProductsService {
	return &SearchProductsService{ctx: ctx}
}

// Run C端商品搜索，返回可购买的SKU列表
func (s *SearchProductsService) Run(req *product.SearchProductsRequest) (*product.SearchProductsResponse, error) {
	page := int(req.Page)
	pageSize := int(req.PageSize)

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	var minPrice, maxPrice float64
	var err error
	if req.MinPrice != "" {
		minPrice, err = utils.PriceConvert(req.MinPrice)
		if err != nil {
			return nil, errs.New(errs.ErrParam.Code, "invalid min_price")
		}
	}
	if req.MaxPrice != "" {
		maxPrice, err = utils.PriceConvert(req.MaxPrice)
		if err != nil {
			return nil, errs.New(errs.ErrParam.Code, "invalid max_price")
		}
	}

	skus, total, err := model.SearchSKUs(
		s.ctx,
		mysql.DB,
		page,
		pageSize,
		req.Keyword,
		req.CategoryId,
		req.BrandId,
		minPrice,
		maxPrice,
		req.SortType,
	)
	if err != nil {
		return nil, errs.New(errs.ErrInternal.Code, "search products failed: "+err.Error())
	}

	spuIDs := make([]uint64, 0, len(skus))
	spuIDSet := make(map[uint64]bool)
	for _, sku := range skus {
		if !spuIDSet[sku.SpuID] {
			spuIDs = append(spuIDs, sku.SpuID)
			spuIDSet[sku.SpuID] = true
		}
	}

	spuMap := make(map[uint64]*model.ProductSPU)
	if len(spuIDs) > 0 {
		spus, err := model.GetProductsByIds(s.ctx, mysql.DB, spuIDs)
		if err != nil {
			return nil, errs.New(errs.ErrInternal.Code, "get products failed: "+err.Error())
		}
		for _, spu := range spus {
			spuMap[spu.ID] = spu
		}
	}

	brandIDs := make([]uint64, 0)
	categoryIDs := make([]uint64, 0)
	brandIDSet := make(map[uint64]bool)
	categoryIDSet := make(map[uint64]bool)

	for _, spu := range spuMap {
		if !brandIDSet[spu.BrandID] {
			brandIDs = append(brandIDs, spu.BrandID)
			brandIDSet[spu.BrandID] = true
		}
		if !categoryIDSet[spu.CategoryID] {
			categoryIDs = append(categoryIDs, spu.CategoryID)
			categoryIDSet[spu.CategoryID] = true
		}
	}

	brandMap := make(map[uint64]*model.ProductBrand)
	if len(brandIDs) > 0 {
		brands, err := model.GetBrandsByIds(s.ctx, mysql.DB, brandIDs)
		if err != nil {
			return nil, errs.New(errs.ErrInternal.Code, "get brands failed: "+err.Error())
		}
		for _, brand := range brands {
			brandMap[brand.ID] = brand
		}
	}

	categoryMap := make(map[uint64]*model.ProductCategory)
	if len(categoryIDs) > 0 {
		categories, err := model.GetCategoriesByIds(s.ctx, mysql.DB, categoryIDs)
		if err != nil {
			return nil, errs.New(errs.ErrInternal.Code, "get categories failed: "+err.Error())
		}
		for _, category := range categories {
			categoryMap[category.ID] = category
		}
	}

	items := make([]*product.SearchProductItem, 0, len(skus))
	for _, sku := range skus {
		spu := spuMap[sku.SpuID]
		if spu == nil {
			continue
		}

		brand := brandMap[spu.BrandID]
		category := categoryMap[spu.CategoryID]

		brandName := ""
		if brand != nil {
			brandName = brand.Name
		}

		categoryName := ""
		if category != nil {
			categoryName = category.Name
		}

		items = append(items, &product.SearchProductItem{
			Sku: &product.ProductSKU{
				Id:          sku.ID,
				SpuId:       sku.SpuID,
				SkuCode:     sku.SkuCode,
				Name:        sku.Name,
				SubTitle:    sku.SubTitle,
				MainImage:   sku.MainImage,
				Price:       utils.PriceToString(sku.Price),
				MarketPrice: utils.PriceToString(sku.MarketPrice),
				Stock:       int32(sku.Stock),
				SkuSpecData: sku.SkuSpecData,
				Version:     int32(sku.Version),
			},
			SpuId:        spu.ID,
			SpuName:      spu.Name,
			BrandId:      spu.BrandID,
			BrandName:    brandName,
			CategoryId:   spu.CategoryID,
			CategoryName: categoryName,
			SpuSaleCount: int32(spu.SaleCount),
		})
	}

	return &product.SearchProductsResponse{
		List:  items,
		Total: total,
	}, nil
}
