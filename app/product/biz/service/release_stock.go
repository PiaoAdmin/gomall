package service

import (
	"context"

	"github.com/PiaoAdmin/pmall/app/product/biz/dal/mysql"
	"github.com/PiaoAdmin/pmall/app/product/biz/dal/redis"
	"github.com/PiaoAdmin/pmall/app/product/biz/model"
	"github.com/PiaoAdmin/pmall/common/errs"
	product "github.com/PiaoAdmin/pmall/rpc_gen/product"
	"github.com/cloudwego/kitex/pkg/klog"
	"gorm.io/gorm"
)

type ReleaseStockService struct {
	ctx context.Context
}

func NewReleaseStockService(ctx context.Context) *ReleaseStockService {
	return &ReleaseStockService{ctx: ctx}
}

// only order service call this interface
func (s *ReleaseStockService) Run(req *product.ReleaseStockRequest) (*product.ReleaseStockResponse, error) {
	if req.OrderSn == "" {
		return nil, errs.New(errs.ErrParam.Code, "order_sn is required")
	}
	if len(req.Items) == 0 {
		return nil, errs.New(errs.ErrParam.Code, "items is empty")
	}

	// 用于记录需要更新缓存的商品ID和销量减少量
	spuSaleUpdates := make(map[uint64]int)

	err := mysql.DB.Transaction(func(tx *gorm.DB) error {
		for _, item := range req.Items {
			if item.SkuId == 0 || item.Count <= 0 {
				return errs.New(errs.ErrParam.Code, "invalid sku_id or count")
			}

			err := model.ReleaseStock(s.ctx, tx, item.SkuId, int(item.Count))
			if err != nil {
				if err == gorm.ErrRecordNotFound {
					return errs.New(errs.ErrInternal.Code, "lock stock not enough or sku not found")
				}
				return err
			}
			sku, err := model.GetSKUByID(s.ctx, tx, item.SkuId)
			if err != nil {
				return err
			}
			if sku == nil {
				return errs.New(errs.ErrRecordNotFound.Code, "sku not found")
			}
			if err := model.DecreaseSaleCount(s.ctx, tx, sku.SpuID, int(item.Count)); err != nil {
				if err == gorm.ErrRecordNotFound {
					return errs.New(errs.ErrInternal.Code, "sale count not enough or spu not found")
				}
				return err
			}
			// 记录销量减少量
			spuSaleUpdates[sku.SpuID] += int(item.Count)
		}
		return nil
	})

	if err != nil {
		if e, ok := err.(*errs.Error); ok {
			return nil, e
		}
		return nil, errs.New(errs.ErrInternal.Code, "release stock failed: "+err.Error())
	}

	// 异步更新热门商品排行榜和缓存
	go func() {
		for spuID, decrement := range spuSaleUpdates {
			// 减少热门商品排行榜中的销量分数
			if err := redis.IncrementProductSaleCount(s.ctx, spuID, -decrement); err != nil {
				klog.Warnf("Failed to update hot product score for SPU %d: %v", spuID, err)
			}
			// 删除商品详情缓存（库存变化）
			if err := redis.DeleteProductDetailCache(s.ctx, spuID); err != nil {
				klog.Warnf("Failed to delete product detail cache for SPU %d: %v", spuID, err)
			}
		}
		// 清除热门商品列表缓存
		if err := redis.DeleteHotProductsCache(s.ctx); err != nil {
			klog.Warnf("Failed to delete hot products cache: %v", err)
		}
	}()

	return &product.ReleaseStockResponse{
		Success: true,
	}, nil
}
