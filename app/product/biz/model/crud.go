package model

import (
	"context"

	"gorm.io/gorm"
)

// TODO: 可优化为消息队列异步处理MongoDB插入，提升性能和解耦
func CreateProductWithTransaction(ctx context.Context, db *gorm.DB,
	spu *ProductSPU, skus []*ProductSKU, detail *ProductDetail) (spuid uint64, err error) {
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err = tx.Error; err != nil {
		return
	}

	spuid, err = CreateSPU(ctx, tx, spu)
	if err != nil {
		tx.Rollback()
		return
	}

	for _, sku := range skus {
		sku.SpuID = spuid
	}
	err = CreateBatchSKU(ctx, tx, skus)
	if err != nil {
		tx.Rollback()
		return
	}

	if err = tx.Commit().Error; err != nil {
		return
	}

	// TODO: 优化方案 - 使用消息队列异步处理，解耦MySQL和MongoDB操作
	detail.SpuID = spuid
	err = CreateProductDetail(ctx, detail)
	if err != nil {
		// MongoDB插入失败，记录日志，不回滚MySQL
		// 可以通过后台任务或消息队列重试
		return spuid, err
	}

	return
}
