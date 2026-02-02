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

type DeductStockService struct {
	ctx context.Context
}

func NewDeductStockService(ctx context.Context) *DeductStockService {
	return &DeductStockService{ctx: ctx}
}

// only order service call this interface
func (s *DeductStockService) Run(req *product.DeductStockRequest) (*product.DeductStockResponse, error) {
	if req.OrderSn == "" {
		return nil, errs.New(errs.ErrParam.Code, "order_sn is required")
	}
	if len(req.Items) == 0 {
		return nil, errs.New(errs.ErrParam.Code, "items is empty")
	}

	// 用于记录需要更新缓存的商品ID和销量增量
	spuSaleUpdates := make(map[uint64]int)

	err := mysql.DB.Transaction(func(tx *gorm.DB) error {
		for _, item := range req.Items {
			if item.SkuId == 0 || item.Count <= 0 {
				return errs.New(errs.ErrParam.Code, "invalid sku_id or count")
			}

			err := model.DeductStock(s.ctx, tx, item.SkuId, int(item.Count))
			if err != nil {
				if err == gorm.ErrRecordNotFound {
					return errs.New(errs.ErrInternal.Code, "stock not enough or sku not found")
				}
				return err
			}
			// TODO: 不是主要业务不应该影响下单
			sku, err := model.GetSKUByID(s.ctx, tx, item.SkuId)
			if err != nil {
				return err
			}
			if sku == nil {
				return errs.New(errs.ErrRecordNotFound.Code, "sku not found")
			}
			if err := model.AddSaleCount(s.ctx, tx, sku.SpuID, int(item.Count)); err != nil {
				return err
			}
			// 记录销量增量
			spuSaleUpdates[sku.SpuID] += int(item.Count)
		}
		return nil
	})

	if err != nil {
		if e, ok := err.(*errs.Error); ok {
			return nil, e
		}
		return nil, errs.New(errs.ErrInternal.Code, "deduct stock failed: "+err.Error())
	}

	// 异步更新热门商品排行榜和缓存
	go func() {
		for spuID, increment := range spuSaleUpdates {
			// 更新热门商品排行榜中的销量分数
			if err := redis.IncrementProductSaleCount(s.ctx, spuID, increment); err != nil {
				klog.Warnf("Failed to update hot product score for SPU %d: %v", spuID, err)
			}
			// 删除商品详情缓存（库存变化）
			if err := redis.DeleteProductDetailCache(s.ctx, spuID); err != nil {
				klog.Warnf("Failed to delete product detail cache for SPU %d: %v", spuID, err)
			}
		}
		// 清除热门商品列表缓存（销量变化可能影响排序）
		if err := redis.DeleteHotProductsCache(s.ctx); err != nil {
			klog.Warnf("Failed to delete hot products cache: %v", err)
		}
	}()

	return &product.DeductStockResponse{
		Success: true,
	}, nil
}
