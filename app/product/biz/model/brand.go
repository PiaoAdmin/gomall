package model

import (
	"context"

	"github.com/PiaoAdmin/pmall/common/uniqueid"
	"gorm.io/gorm"
)

type ProductBrand struct {
	Model
	Name        string `gorm:"type:varchar(64);not null;uniqueIndex;comment:品牌名称"`
	FirstLetter string `gorm:"type:char(3);comment:品牌首字母"`
	Logo        string `gorm:"type:varchar(500);comment:品牌LOGO URL"`
	Sort        int    `gorm:"default:0"`
	ShowStatus  bool   `gorm:"default:true;comment:是否显示"`
}

func (ProductBrand) TableName() string {
	return "product_brand"
}

func (b *ProductBrand) BeforeCreate(tx *gorm.DB) (err error) {
	// 雪花算法生成id
	b.ID = uint64(uniqueid.GenId())
	return
}

func ListBrands(ctx context.Context, db *gorm.DB, page, pageSize int) ([]*ProductBrand, int64, error) {
	var brands []*ProductBrand
	var total int64

	query := db.WithContext(ctx).Model(&ProductBrand{}).Where("show_status = ?", true)

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	err := query.Order("sort ASC, id ASC").Offset(offset).Limit(pageSize).Find(&brands).Error
	return brands, total, err
}

func GetBrandByID(ctx context.Context, db *gorm.DB, id uint64) (*ProductBrand, error) {
	var brand ProductBrand
	err := db.WithContext(ctx).Where("id = ?", id).First(&brand).Error
	if err != nil {
		return nil, err
	}
	return &brand, nil
}

// GetBrandsByIds 批量获取品牌
func GetBrandsByIds(ctx context.Context, db *gorm.DB, ids []uint64) ([]*ProductBrand, error) {
	if len(ids) == 0 {
		return []*ProductBrand{}, nil
	}
	var brands []*ProductBrand
	err := db.WithContext(ctx).Where("id IN ?", ids).Find(&brands).Error
	return brands, err
}
