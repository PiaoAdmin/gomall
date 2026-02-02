package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/redis/go-redis/v9"
)

// 缓存key前缀和过期时间常量
const (
	// 商品详情缓存
	ProductDetailKeyPrefix = "product:detail:"
	ProductDetailExpire    = 30 * time.Minute

	// 热门商品列表缓存
	HotProductsKey    = "product:hot:list"
	HotProductsExpire = 10 * time.Minute
	HotProductsLimit  = 100 // 缓存前100个热门商品

	// 商品列表缓存
	ProductListKeyPrefix = "product:list:"
	ProductListExpire    = 5 * time.Minute

	// 分类商品缓存
	CategoryProductsKeyPrefix = "product:category:"
	CategoryProductsExpire    = 10 * time.Minute
)

// CachedProductDetail 缓存的商品详情结构
type CachedProductDetail struct {
	SPU      *CachedSPU      `json:"spu"`
	SKUs     []*CachedSKU    `json:"skus"`
	Category *CachedCategory `json:"category,omitempty"`
	Brand    *CachedBrand    `json:"brand,omitempty"`
	Detail   *CachedDetail   `json:"detail,omitempty"`
}

type CachedSPU struct {
	ID            uint64  `json:"id"`
	BrandID       uint64  `json:"brand_id"`
	CategoryID    uint64  `json:"category_id"`
	Name          string  `json:"name"`
	SubTitle      string  `json:"sub_title"`
	MainImage     string  `json:"main_image"`
	PublishStatus int8    `json:"publish_status"`
	VerifyStatus  int8    `json:"verify_status"`
	LowPrice      float64 `json:"low_price"`
	HighPrice     float64 `json:"high_price"`
	SaleCount     int     `json:"sale_count"`
	Sort          int     `json:"sort"`
	ServiceBits   int64   `json:"service_bits"`
	Version       int     `json:"version"`
}

type CachedSKU struct {
	ID          uint64  `json:"id"`
	SpuID       uint64  `json:"spu_id"`
	SkuCode     string  `json:"sku_code"`
	Name        string  `json:"name"`
	SubTitle    string  `json:"sub_title"`
	MainImage   string  `json:"main_image"`
	Price       float64 `json:"price"`
	MarketPrice float64 `json:"market_price"`
	Stock       int     `json:"stock"`
	LockStock   int     `json:"lock_stock"`
	SkuSpecData string  `json:"sku_spec_data"`
	Version     int     `json:"version"`
}

type CachedCategory struct {
	ID       uint64 `json:"id"`
	ParentID uint64 `json:"parent_id"`
	Name     string `json:"name"`
	Level    int    `json:"level"`
	Icon     string `json:"icon"`
	Unit     string `json:"unit"`
	Sort     int    `json:"sort"`
}

type CachedBrand struct {
	ID          uint64 `json:"id"`
	Name        string `json:"name"`
	FirstLetter string `json:"first_letter"`
	Logo        string `json:"logo"`
	Sort        int    `json:"sort"`
}

type CachedDetail struct {
	Description   string   `json:"description"`
	Images        []string `json:"images"`
	Videos        []string `json:"videos"`
	MarketTagJSON string   `json:"market_tag_json"`
	TechTagJSON   string   `json:"tech_tag_json"`
	FaqJSON       string   `json:"faq_json"`
}

// GetProductDetailFromCache 从缓存获取商品详情
func GetProductDetailFromCache(ctx context.Context, productID uint64) (*CachedProductDetail, error) {
	key := fmt.Sprintf("%s%d", ProductDetailKeyPrefix, productID)
	data, err := RedisClient.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // 缓存未命中
		}
		return nil, err
	}

	var cached CachedProductDetail
	if err := json.Unmarshal(data, &cached); err != nil {
		klog.Warnf("Failed to unmarshal product detail cache: %v", err)
		return nil, err
	}
	return &cached, nil
}

// SetProductDetailCache 设置商品详情缓存
func SetProductDetailCache(ctx context.Context, productID uint64, detail *CachedProductDetail) error {
	key := fmt.Sprintf("%s%d", ProductDetailKeyPrefix, productID)
	data, err := json.Marshal(detail)
	if err != nil {
		return err
	}
	return RedisClient.Set(ctx, key, data, ProductDetailExpire).Err()
}

// DeleteProductDetailCache 删除商品详情缓存
func DeleteProductDetailCache(ctx context.Context, productID uint64) error {
	key := fmt.Sprintf("%s%d", ProductDetailKeyPrefix, productID)
	return RedisClient.Del(ctx, key).Err()
}

// BatchDeleteProductDetailCache 批量删除商品详情缓存
func BatchDeleteProductDetailCache(ctx context.Context, productIDs []uint64) error {
	if len(productIDs) == 0 {
		return nil
	}
	keys := make([]string, len(productIDs))
	for i, id := range productIDs {
		keys[i] = fmt.Sprintf("%s%d", ProductDetailKeyPrefix, id)
	}
	return RedisClient.Del(ctx, keys...).Err()
}

// HotProductInfo 热门商品简要信息（用于排行榜展示）
type HotProductInfo struct {
	ID        uint64  `json:"id"`
	Name      string  `json:"name"`
	SubTitle  string  `json:"sub_title"`
	MainImage string  `json:"main_image"`
	LowPrice  float64 `json:"low_price"`
	SaleCount int     `json:"sale_count"`
}

// GetHotProductsFromCache 从缓存获取热门商品列表
func GetHotProductsFromCache(ctx context.Context) ([]*HotProductInfo, error) {
	data, err := RedisClient.Get(ctx, HotProductsKey).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	var products []*HotProductInfo
	if err := json.Unmarshal(data, &products); err != nil {
		klog.Warnf("Failed to unmarshal hot products cache: %v", err)
		return nil, err
	}
	return products, nil
}

// SetHotProductsCache 设置热门商品列表缓存
func SetHotProductsCache(ctx context.Context, products []*HotProductInfo) error {
	data, err := json.Marshal(products)
	if err != nil {
		return err
	}
	return RedisClient.Set(ctx, HotProductsKey, data, HotProductsExpire).Err()
}

// DeleteHotProductsCache 删除热门商品列表缓存
func DeleteHotProductsCache(ctx context.Context) error {
	return RedisClient.Del(ctx, HotProductsKey).Err()
}

// 使用 ZSet 维护热门商品排行榜（按销量）
const HotProductsZSetKey = "product:hot:zset"

// UpdateHotProductScore 更新热门商品评分（使用销量作为分数）
func UpdateHotProductScore(ctx context.Context, productID uint64, saleCount int) error {
	return RedisClient.ZAdd(ctx, HotProductsZSetKey, redis.Z{
		Score:  float64(saleCount),
		Member: strconv.FormatUint(productID, 10),
	}).Err()
}

// GetTopHotProductIDs 获取销量最高的商品ID列表
func GetTopHotProductIDs(ctx context.Context, limit int) ([]uint64, error) {
	// 按分数从高到低获取
	result, err := RedisClient.ZRevRange(ctx, HotProductsZSetKey, 0, int64(limit-1)).Result()
	if err != nil {
		return nil, err
	}

	ids := make([]uint64, 0, len(result))
	for _, s := range result {
		id, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			continue
		}
		ids = append(ids, id)
	}
	return ids, nil
}

// IncrementProductSaleCount 增加商品销量（原子操作）
func IncrementProductSaleCount(ctx context.Context, productID uint64, increment int) error {
	return RedisClient.ZIncrBy(ctx, HotProductsZSetKey, float64(increment), strconv.FormatUint(productID, 10)).Err()
}

// 商品列表缓存（用于首页、分类页等场景）

// ProductListCacheKey 生成商品列表缓存key
func ProductListCacheKey(page, pageSize int, keyword string, categoryID, brandID uint64) string {
	return fmt.Sprintf("%s%d:%d:%s:%d:%d", ProductListKeyPrefix, page, pageSize, keyword, categoryID, brandID)
}

// CachedProductList 缓存的商品列表
type CachedProductList struct {
	Products []*CachedSPU `json:"products"`
	Total    int64        `json:"total"`
}

// GetProductListFromCache 从缓存获取商品列表
func GetProductListFromCache(ctx context.Context, page, pageSize int, keyword string, categoryID, brandID uint64) (*CachedProductList, error) {
	key := ProductListCacheKey(page, pageSize, keyword, categoryID, brandID)
	data, err := RedisClient.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	var cached CachedProductList
	if err := json.Unmarshal(data, &cached); err != nil {
		klog.Warnf("Failed to unmarshal product list cache: %v", err)
		return nil, err
	}
	return &cached, nil
}

// SetProductListCache 设置商品列表缓存
func SetProductListCache(ctx context.Context, page, pageSize int, keyword string, categoryID, brandID uint64, list *CachedProductList) error {
	key := ProductListCacheKey(page, pageSize, keyword, categoryID, brandID)
	data, err := json.Marshal(list)
	if err != nil {
		return err
	}
	return RedisClient.Set(ctx, key, data, ProductListExpire).Err()
}

// InvalidateProductListCache 使商品列表缓存失效（当商品更新时调用）
func InvalidateProductListCache(ctx context.Context) error {
	// 使用 SCAN 删除所有商品列表缓存
	var cursor uint64
	pattern := ProductListKeyPrefix + "*"
	for {
		keys, newCursor, err := RedisClient.Scan(ctx, cursor, pattern, 100).Result()
		if err != nil {
			return err
		}
		if len(keys) > 0 {
			if err := RedisClient.Del(ctx, keys...).Err(); err != nil {
				return err
			}
		}
		cursor = newCursor
		if cursor == 0 {
			break
		}
	}
	return nil
}

// 缓存预热：批量设置商品详情缓存
func WarmUpProductDetailCache(ctx context.Context, details map[uint64]*CachedProductDetail) error {
	if len(details) == 0 {
		return nil
	}

	pipe := RedisClient.Pipeline()
	for productID, detail := range details {
		key := fmt.Sprintf("%s%d", ProductDetailKeyPrefix, productID)
		data, err := json.Marshal(detail)
		if err != nil {
			klog.Warnf("Failed to marshal product detail for warm up: %v", err)
			continue
		}
		pipe.Set(ctx, key, data, ProductDetailExpire)
	}
	_, err := pipe.Exec(ctx)
	return err
}

// RefreshHotProductsZSet 刷新热门商品 ZSet（从数据库同步数据）
func RefreshHotProductsZSet(ctx context.Context, products []struct {
	ID        uint64
	SaleCount int
}) error {
	if len(products) == 0 {
		return nil
	}

	pipe := RedisClient.Pipeline()
	for _, p := range products {
		pipe.ZAdd(ctx, HotProductsZSetKey, redis.Z{
			Score:  float64(p.SaleCount),
			Member: strconv.FormatUint(p.ID, 10),
		})
	}
	_, err := pipe.Exec(ctx)
	return err
}
