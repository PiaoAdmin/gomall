package service

import (
	"context"
	"time"

	"github.com/PiaoAdmin/pmall/app/product/biz/dal/mysql"
	"github.com/PiaoAdmin/pmall/app/product/biz/dal/redis"
	"github.com/PiaoAdmin/pmall/app/product/biz/model"
	"github.com/cloudwego/kitex/pkg/klog"
)

// CacheWarmUpService 缓存预热服务
type CacheWarmUpService struct {
	ctx context.Context
}

func NewCacheWarmUpService(ctx context.Context) *CacheWarmUpService {
	return &CacheWarmUpService{ctx: ctx}
}

// WarmUpHotProducts 预热热门商品缓存
func (s *CacheWarmUpService) WarmUpHotProducts() error {
	klog.Info("Starting hot products cache warm-up...")

	// 从数据库获取销量最高的商品
	spus, _, err := model.ListProductsBySaleCount(s.ctx, mysql.DB, redis.HotProductsLimit)
	if err != nil {
		klog.Errorf("Failed to get hot products from database: %v", err)
		return err
	}

	if len(spus) == 0 {
		klog.Warn("No products found for cache warm-up")
		return nil
	}

	// 更新热门商品ZSet
	products := make([]struct {
		ID        uint64
		SaleCount int
	}, len(spus))
	for i, spu := range spus {
		products[i] = struct {
			ID        uint64
			SaleCount int
		}{
			ID:        spu.ID,
			SaleCount: spu.SaleCount,
		}
	}
	if err := redis.RefreshHotProductsZSet(s.ctx, products); err != nil {
		klog.Errorf("Failed to refresh hot products ZSet: %v", err)
		return err
	}

	// 更新热门商品列表缓存
	hotProducts := make([]*redis.HotProductInfo, 0, len(spus))
	for _, spu := range spus {
		hotProducts = append(hotProducts, &redis.HotProductInfo{
			ID:        spu.ID,
			Name:      spu.Name,
			SubTitle:  spu.SubTitle,
			MainImage: spu.MainImage,
			LowPrice:  spu.LowPrice,
			SaleCount: spu.SaleCount,
		})
	}
	if err := redis.SetHotProductsCache(s.ctx, hotProducts); err != nil {
		klog.Errorf("Failed to set hot products cache: %v", err)
		return err
	}

	klog.Infof("Hot products cache warm-up completed, cached %d products", len(spus))
	return nil
}

// WarmUpTopProductDetails 预热热门商品详情缓存
func (s *CacheWarmUpService) WarmUpTopProductDetails(limit int) error {
	klog.Infof("Starting top %d product details cache warm-up...", limit)

	// 获取热门商品ID
	hotIDs, err := redis.GetTopHotProductIDs(s.ctx, limit)
	if err != nil {
		klog.Errorf("Failed to get hot product IDs: %v", err)
		return err
	}

	if len(hotIDs) == 0 {
		klog.Warn("No hot product IDs found for cache warm-up")
		return nil
	}

	// 逐个预热商品详情
	warmedCount := 0
	for _, productID := range hotIDs {
		// 检查缓存是否已存在
		cached, _ := redis.GetProductDetailFromCache(s.ctx, productID)
		if cached != nil {
			continue // 已缓存，跳过
		}

		// 获取商品详情并缓存
		spu, err := model.GetSPUByID(s.ctx, mysql.DB, productID)
		if err != nil {
			klog.Warnf("Failed to get SPU %d: %v", productID, err)
			continue
		}

		skus, err := model.GetSKUsBySpuID(s.ctx, mysql.DB, productID)
		if err != nil {
			klog.Warnf("Failed to get SKUs for SPU %d: %v", productID, err)
			continue
		}

		var category *model.ProductCategory
		if spu.CategoryID > 0 {
			category, _ = model.GetCategoryByID(s.ctx, mysql.DB, spu.CategoryID)
		}

		var brand *model.ProductBrand
		if spu.BrandID > 0 {
			brand, _ = model.GetBrandByID(s.ctx, mysql.DB, spu.BrandID)
		}

		detail, _ := model.GetProductDetailBySpuID(s.ctx, productID)

		cachedDetail := buildCacheData(spu, skus, category, brand, detail)
		if err := redis.SetProductDetailCache(s.ctx, productID, cachedDetail); err != nil {
			klog.Warnf("Failed to cache product %d: %v", productID, err)
			continue
		}
		warmedCount++
	}

	klog.Infof("Product details cache warm-up completed, cached %d products", warmedCount)
	return nil
}

// buildCacheData 构建缓存数据
func buildCacheData(spu *model.ProductSPU, skus []*model.ProductSKU, category *model.ProductCategory, brand *model.ProductBrand, detail *model.ProductDetail) *redis.CachedProductDetail {
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

// StartCacheRefreshTask 启动定时刷新缓存任务
func StartCacheRefreshTask(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Minute) // 每5分钟刷新一次热门商品缓存
	go func() {
		for {
			select {
			case <-ctx.Done():
				ticker.Stop()
				klog.Info("Cache refresh task stopped")
				return
			case <-ticker.C:
				service := NewCacheWarmUpService(ctx)
				if err := service.WarmUpHotProducts(); err != nil {
					klog.Errorf("Periodic hot products cache refresh failed: %v", err)
				}
			}
		}
	}()
	klog.Info("Cache refresh task started")
}
