package model

import (
	"context"

	"github.com/PiaoAdmin/pmall/common/uniqueid"
	"gorm.io/gorm"
)

type ProductSPU struct {
	Model
	BrandID       uint64  `gorm:"not null;comment:品牌ID;index:idx_cat_brand"`
	CategoryID    uint64  `gorm:"not null;comment:分类ID;index:idx_cat_brand"`
	Name          string  `gorm:"type:varchar(255);not null;comment:商品名称"`
	SubTitle      string  `gorm:"type:varchar(500);comment:商品副标题"`
	MainImage     string  `gorm:"type:varchar(1000);comment:商品主图"`
	PublishStatus int8    `gorm:"default:0;comment:商品发布状态:0-未发布,1-已发布;index:idx_status"`
	VerifyStatus  int8    `gorm:"default:0;comment:商品审核状态:0-未审核,1-审核通过,2-审核不通过;index:idx_status"`
	LowPrice      float64 `gorm:"type:decimal(10,2);comment:最低售价;index:idx_price_range"`
	HighPrice     float64 `gorm:"type:decimal(10,2);comment:最高售价;index:idx_price_range"`
	SaleCount     int     `gorm:"default:0;comment:销量;index:idx_sale_count"`
	Sort          int     `gorm:"default:0;comment:排序权重;index:idx_sort"`
	ServiceBits   int64   `gorm:"default:0;comment:商品服务:用二进制位存储,每一位代表一种服务"`
	Version       int     `gorm:"default:1;comment:乐观锁版本号"`
}

func (ProductSPU) TableName() string {
	return "product_spu"
}

func (s *ProductSPU) BeforeCreate(tx *gorm.DB) (err error) {
	// 雪花算法生成id
	s.ID = uint64(uniqueid.GenId())
	return
}

func CreateSPU(ctx context.Context, db *gorm.DB, spu *ProductSPU) (spuid uint64, err error) {
	err = gorm.G[ProductSPU](db).Create(ctx, spu)
	if err == nil {
		spuid = spu.ID
	}
	return
}

func UpdateSPU(ctx context.Context, db *gorm.DB, id uint64, updates map[string]interface{}) (effectRows int64, err error) {
	if len(updates) == 0 {
		return 0, nil
	}
	updates["version"] = gorm.Expr("version + 1")
	result := db.WithContext(ctx).Model(&ProductSPU{}).Where("id = ?", id).Updates(updates)
	return result.RowsAffected, result.Error
}

func GetSPUByID(ctx context.Context, db *gorm.DB, spuid uint64) (*ProductSPU, error) {
	spu, err := gorm.G[ProductSPU](db).Where("id = ?", spuid).First(ctx)
	if err != nil {
		return nil, err
	}
	return &spu, nil
}

func BatchUpdateStatus(ctx context.Context, db *gorm.DB, ids []uint64, updates map[string]interface{}) error {
	if len(ids) == 0 || len(updates) == 0 {
		return nil
	}
	return db.WithContext(ctx).Model(&ProductSPU{}).Where("id IN ?", ids).Updates(updates).Error
}

func ListProducts(ctx context.Context, db *gorm.DB, page, pageSize int, keyword string, categoryID, brandID uint64, publishStatus, verifyStatus int32) ([]*ProductSPU, int64, error) {
	var spus []*ProductSPU
	var total int64

	query := db.WithContext(ctx).Model(&ProductSPU{})

	if keyword != "" {
		query = query.Where("name LIKE ? OR sub_title LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	if categoryID > 0 {
		query = query.Where("category_id = ?", categoryID)
	}

	if brandID > 0 {
		query = query.Where("brand_id = ?", brandID)
	}

	// 暂时去掉状态校验，简化业务流程
	// if publishStatus >= 0 {
	// 	query = query.Where("publish_status = ?", publishStatus)
	// }

	// if verifyStatus >= 0 {
	// 	query = query.Where("verify_status = ?", verifyStatus)
	// }

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Order("sort DESC, created_at DESC").Offset(offset).Limit(pageSize).Find(&spus).Error
	return spus, total, err
}

func GetProductsByIds(ctx context.Context, db *gorm.DB, ids []uint64) ([]*ProductSPU, error) {
	if len(ids) == 0 {
		return []*ProductSPU{}, nil
	}
	var spus []*ProductSPU
	err := db.WithContext(ctx).Where("id IN ?", ids).Find(&spus).Error
	return spus, err
}

// ListProductsBySaleCount 按销量排序获取商品列表（用于热门商品）
func ListProductsBySaleCount(ctx context.Context, db *gorm.DB, limit int) ([]*ProductSPU, int64, error) {
	var spus []*ProductSPU
	var total int64

	query := db.WithContext(ctx).Model(&ProductSPU{})

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Order("sale_count DESC, sort DESC").Limit(limit).Find(&spus).Error
	return spus, total, err
}

func AddSaleCount(ctx context.Context, db *gorm.DB, spuid uint64, count int) error {
	updates := map[string]interface{}{
		"sale_count": gorm.Expr("sale_count + ?", count),
	}
	result := db.WithContext(ctx).Model(&ProductSPU{}).Where("id = ?", spuid).Updates(updates)
	return result.Error
}

func DecreaseSaleCount(ctx context.Context, db *gorm.DB, spuid uint64, count int) error {
	if count <= 0 {
		return nil
	}
	updates := map[string]interface{}{
		"sale_count": gorm.Expr("GREATEST(sale_count - ?, 0)", count),
	}
	result := db.WithContext(ctx).Model(&ProductSPU{}).Where("id = ?", spuid).Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
