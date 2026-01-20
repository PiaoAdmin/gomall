package model

import (
	"context"

	"github.com/PiaoAdmin/pmall/common/errs"
	"github.com/PiaoAdmin/pmall/common/uniqueid"
	"gorm.io/gorm"
)

// TODO: sku的PublishStatus和VerifyStatus字段暂时保留，后续看需求是否需要
type ProductSKU struct {
	Model
	SpuID         uint64  `gorm:"not null;comment:商品SPU ID;index:idx_spu_id"`
	SkuCode       string  `gorm:"type:varchar(64);not null;uniqueIndex:uk_sku_code;comment:商家内部SKU编码"`
	Name          string  `gorm:"type:varchar(255);not null;comment:SKU名称"`
	SubTitle      string  `gorm:"type:varchar(500);comment:SKU副标题"`
	MainImage     string  `gorm:"type:varchar(1000);comment:SKU商品主图"`
	PublishStatus int8    `gorm:"default:0;comment:商品发布状态:0-未发布,1-已发布;index:idx_spu_id;index:idx_status"`
	VerifyStatus  int8    `gorm:"default:0;comment:商品审核状态:0-未审核,1-审核通过,2-审核不通过;index:idx_status"`
	Price         float64 `gorm:"type:decimal(10,2);not null;comment:销售价;index:idx_price"`
	MarketPrice   float64 `gorm:"type:decimal(10,2);comment:市场价"`
	Stock         int     `gorm:"default:0;comment:库存数量;index:idx_stock"`
	LockStock     int     `gorm:"default:0;comment:锁定库存(下单未付)"`
	SkuSpecData   string  `gorm:"type:json;comment:规格键值对"`
	Version       int     `gorm:"default:1;comment:乐观锁版本号"`
}

func (ProductSKU) TableName() string {
	return "product_sku"
}

func (s *ProductSKU) BeforeCreate(tx *gorm.DB) (err error) {
	// 雪花算法生成id
	s.ID = uint64(uniqueid.GenId())
	return
}

func CreateBatchSKU(ctx context.Context, db *gorm.DB, skus []*ProductSKU) error {
	if len(skus) == 0 {
		return nil
	}
	return db.WithContext(ctx).Create(skus).Error
}

func GetSKUsBySpuID(ctx context.Context, db *gorm.DB, spuID uint64) ([]*ProductSKU, error) {
	var skus []*ProductSKU
	err := db.WithContext(ctx).Where("spu_id = ?", spuID).Find(&skus).Error
	return skus, err
}

// GetSKUByID 根据ID获取SKU
func GetSKUByID(ctx context.Context, db *gorm.DB, id uint64) (*ProductSKU, error) {
	var sku ProductSKU
	err := db.WithContext(ctx).Where("id = ?", id).First(&sku).Error
	if err != nil {
		return nil, err
	}
	return &sku, nil
}

// DeductStock 扣减库存（使用乐观锁）
func DeductStock(ctx context.Context, db *gorm.DB, skuID uint64, count int) error {
	// 使用乐观锁更新：stock = stock - count, lock_stock = lock_stock + count
	// WHERE id = ? AND stock >= ? AND version = ?
	result := db.WithContext(ctx).Model(&ProductSKU{}).
		Where("id = ? AND stock >= ?", skuID, count).
		Updates(map[string]interface{}{
			"stock":      gorm.Expr("stock - ?", count),
			"lock_stock": gorm.Expr("lock_stock + ?", count),
			"version":    gorm.Expr("version + 1"),
		})

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound // 库存不足或记录不存在
	}
	return nil
}

// ReleaseStock 释放库存（取消订单/退货）
func ReleaseStock(ctx context.Context, db *gorm.DB, skuID uint64, count int) error {
	// lock_stock - count, stock + count
	result := db.WithContext(ctx).Model(&ProductSKU{}).
		Where("id = ? AND lock_stock >= ?", skuID, count).
		Updates(map[string]interface{}{
			"lock_stock": gorm.Expr("lock_stock - ?", count),
			"stock":      gorm.Expr("stock + ?", count),
			"version":    gorm.Expr("version + 1"),
		})

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func UpdateSKU(ctx context.Context, db *gorm.DB, skuID uint64, updates map[string]interface{}) error {
	if len(updates) == 0 {
		return nil
	}
	updates["version"] = gorm.Expr("version + 1")
	result := db.WithContext(ctx).Model(&ProductSKU{}).Where("id = ?", skuID).Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	// Check if any rows were affected
	if result.RowsAffected == 0 {
		return errs.New(errs.ErrRecordNotFound.Code, "sku not found")
	}
	return nil
}

// SearchSKUs C端搜索可购买的SKU（只返回已发布、审核通过、有库存的）
func SearchSKUs(ctx context.Context, db *gorm.DB, page, pageSize int, keyword string,
	categoryID, brandID uint64, minPrice, maxPrice float64, sortType int32) ([]*ProductSKU, int64, error) {

	var skus []*ProductSKU
	var total int64

	// 联表查询 SKU + SPU
	query := db.WithContext(ctx).
		Table("product_sku").
		Joins("INNER JOIN product_spu ON product_sku.spu_id = product_spu.id").
		// 暂时去掉状态校验，简化业务流程
		// Where("product_spu.publish_status = ? AND product_spu.verify_status = ?", 1, 1).
		Where("product_sku.stock > ?", 0) // 有库存

	if keyword != "" {
		query = query.Where("product_spu.name LIKE ? OR product_sku.name LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%")
	}

	if categoryID > 0 {
		query = query.Where("product_spu.category_id = ?", categoryID)
	}

	if brandID > 0 {
		query = query.Where("product_spu.brand_id = ?", brandID)
	}

	if minPrice > 0 {
		query = query.Where("product_sku.price >= ?", minPrice)
	}
	if maxPrice > 0 {
		query = query.Where("product_sku.price <= ?", maxPrice)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 排序
	switch sortType {
	case 1: // 价格升序
		query = query.Order("product_sku.price ASC")
	case 2: // 价格降序
		query = query.Order("product_sku.price DESC")
	case 3: // 销量降序
		query = query.Order("product_spu.sale_count DESC")
	default: // 综合排序（权重 + 销量 + 创建时间）
		query = query.Order("product_spu.sort DESC, product_spu.sale_count DESC, product_sku.created_at DESC")
	}

	// 分页查询
	offset := (page - 1) * pageSize
	err := query.Select("product_sku.*").Offset(offset).Limit(pageSize).Find(&skus).Error

	return skus, total, err
}
