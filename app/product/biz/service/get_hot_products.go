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

type GetHotProductsService struct {
	ctx context.Context
}

func NewGetHotProductsService(ctx context.Context) *GetHotProductsService {
	return &GetHotProductsService{ctx: ctx}
}

// Run 获取热门商品列表
func (s *GetHotProductsService) Run(req *product.GetHotProductsRequest) (*product.GetHotProductsResponse, error) {
	limit := int(req.Limit)
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	// 优先从缓存获取热门商品列表
	cachedProducts, err := redis.GetHotProductsFromCache(s.ctx)
	if err != nil {
		klog.Warnf("Failed to get hot products from cache: %v", err)
	}
	if len(cachedProducts) > 0 {
		klog.Debugf("Hot products cache hit, count=%d", len(cachedProducts))
		// 根据limit截取
		if len(cachedProducts) > limit {
			cachedProducts = cachedProducts[:limit]
		}
		return s.convertCachedToResponse(cachedProducts), nil
	}

	// 尝试从ZSet获取热门商品ID
	hotIDs, err := redis.GetTopHotProductIDs(s.ctx, limit)
	if err != nil {
		klog.Warnf("Failed to get hot product IDs from ZSet: %v", err)
	}

	var spus []*model.ProductSPU

	if len(hotIDs) > 0 {
		// 从数据库批量获取商品信息
		spus, err = model.GetProductsByIds(s.ctx, mysql.DB, hotIDs)
		if err != nil {
			return nil, errs.New(errs.ErrInternal.Code, "get hot products failed: "+err.Error())
		}
	} else {
		// ZSet为空，从数据库按销量排序获取
		klog.Info("Hot products ZSet is empty, fetching from database by sale_count")
		spus, _, err = model.ListProductsBySaleCount(s.ctx, mysql.DB, limit)
		if err != nil {
			return nil, errs.New(errs.ErrInternal.Code, "get hot products failed: "+err.Error())
		}

		// 异步初始化热门商品ZSet
		go func() {
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
				klog.Warnf("Failed to refresh hot products ZSet: %v", err)
			}
		}()
	}

	// 构建响应
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

	// 异步设置缓存
	go func() {
		if err := redis.SetHotProductsCache(s.ctx, hotProducts); err != nil {
			klog.Warnf("Failed to set hot products cache: %v", err)
		}
	}()

	return s.convertCachedToResponse(hotProducts), nil
}

// convertCachedToResponse 将缓存数据转换为响应
func (s *GetHotProductsService) convertCachedToResponse(products []*redis.HotProductInfo) *product.GetHotProductsResponse {
	protoProducts := make([]*product.HotProductInfo, 0, len(products))
	for _, p := range products {
		protoProducts = append(protoProducts, &product.HotProductInfo{
			Id:        p.ID,
			Name:      p.Name,
			SubTitle:  p.SubTitle,
			MainImage: p.MainImage,
			LowPrice:  p.LowPrice,
			SaleCount: int32(p.SaleCount),
		})
	}
	return &product.GetHotProductsResponse{
		Products: protoProducts,
	}
}

// RefreshHotProductsCache 刷新热门商品缓存（可由定时任务调用）
func RefreshHotProductsCache(ctx context.Context) error {
	// 从数据库获取销量最高的商品
	spus, _, err := model.ListProductsBySaleCount(ctx, mysql.DB, redis.HotProductsLimit)
	if err != nil {
		return err
	}

	// 更新ZSet
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
	if err := redis.RefreshHotProductsZSet(ctx, products); err != nil {
		return err
	}

	// 更新缓存列表
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

	return redis.SetHotProductsCache(ctx, hotProducts)
}
