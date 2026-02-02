package service

import (
	"context"

	"github.com/PiaoAdmin/pmall/app/product/biz/dal/mysql"
	"github.com/PiaoAdmin/pmall/app/product/biz/dal/redis"
	"github.com/PiaoAdmin/pmall/app/product/biz/model"
	"github.com/PiaoAdmin/pmall/common/errs"
	product "github.com/PiaoAdmin/pmall/rpc_gen/product"
	"github.com/cloudwego/kitex/pkg/klog"
)

type ListProductsService struct {
	ctx context.Context
}

func NewListProductsService(ctx context.Context) *ListProductsService {
	return &ListProductsService{ctx: ctx}
}

func (s *ListProductsService) Run(req *product.ListProductsRequest) (*product.ListProductsResponse, error) {
	page := int(req.Page)
	pageSize := int(req.PageSize)

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100 // 限制最大页面大小
	}

	// 尝试从缓存获取商品列表
	cachedList, err := redis.GetProductListFromCache(s.ctx, page, pageSize, req.Keyword, req.CategoryId, req.BrandId)
	if err != nil {
		klog.Warnf("Failed to get product list from cache: %v", err)
	}
	if cachedList != nil {
		klog.Debugf("Product list cache hit: page=%d, pageSize=%d", page, pageSize)
		return s.convertCachedListToResponse(cachedList), nil
	}

	// 缓存未命中，从数据库获取
	klog.Debugf("Product list cache miss, fetching from database")
	spus, total, err := model.ListProducts(
		s.ctx,
		mysql.DB,
		page,
		pageSize,
		req.Keyword,
		req.CategoryId,
		req.BrandId,
		req.PublishStatus,
		req.VerifyStatus,
	)
	if err != nil {
		return nil, errs.New(errs.ErrInternal.Code, "list products failed: "+err.Error())
	}

	protoSPUs := make([]*product.ProductSPU, 0, len(spus))
	for _, spu := range spus {
		protoSPUs = append(protoSPUs, &product.ProductSPU{
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
		})
	}

	// 异步设置缓存
	go func() {
		cachedSPUs := make([]*redis.CachedSPU, 0, len(spus))
		for _, spu := range spus {
			cachedSPUs = append(cachedSPUs, &redis.CachedSPU{
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
			})
		}
		cacheData := &redis.CachedProductList{
			Products: cachedSPUs,
			Total:    total,
		}
		if err := redis.SetProductListCache(s.ctx, page, pageSize, req.Keyword, req.CategoryId, req.BrandId, cacheData); err != nil {
			klog.Warnf("Failed to set product list cache: %v", err)
		}
	}()

	return &product.ListProductsResponse{
		List:  protoSPUs,
		Total: total,
	}, nil
}

// convertCachedListToResponse 将缓存数据转换为响应
func (s *ListProductsService) convertCachedListToResponse(cached *redis.CachedProductList) *product.ListProductsResponse {
	protoSPUs := make([]*product.ProductSPU, 0, len(cached.Products))
	for _, spu := range cached.Products {
		protoSPUs = append(protoSPUs, &product.ProductSPU{
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
		})
	}
	return &product.ListProductsResponse{
		List:  protoSPUs,
		Total: cached.Total,
	}
}
