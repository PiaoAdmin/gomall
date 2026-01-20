package model

import (
	"context"

	"github.com/PiaoAdmin/pmall/common/uniqueid"
	"gorm.io/gorm"
)

type ProductCategory struct {
	Model
	Name       string `gorm:"type:varchar(64);not null;uniqueIndex;comment:分类名称"`
	ParentID   uint64 `gorm:"default:0;comment:父级分类ID"`
	Level      int    `gorm:"default:0;comment:分类级别"`
	Sort       int    `gorm:"default:0;comment:排序"`
	ShowStatus bool   `gorm:"default:true;comment:是否显示"`
	Icon       string `gorm:"type:varchar(255);comment:图标URL"`
	Unit       string `gorm:"type:varchar(16);comment:计量单位"`
}

func (ProductCategory) TableName() string {
	return "product_category"
}

// TODO: 目前所有分类都是直接写入库中，没有用到雪花id
func (c *ProductCategory) BeforeCreate(tx *gorm.DB) (err error) {
	c.ID = uint64(uniqueid.GenId())
	return
}

func GetCategoriesByParentID(ctx context.Context, db *gorm.DB, parentID uint64) ([]*ProductCategory, error) {
	var categories []*ProductCategory
	err := db.WithContext(ctx).Where("parent_id = ? AND show_status = ?", parentID, true).
		Order("sort ASC, id ASC").Find(&categories).Error
	return categories, err
}

func GetAllCategories(ctx context.Context, db *gorm.DB) ([]*ProductCategory, error) {
	var categories []*ProductCategory
	err := db.WithContext(ctx).Where("show_status = ?", true).
		Order("level ASC, sort ASC, id ASC").Find(&categories).Error
	return categories, err
}

func GetCategoryByID(ctx context.Context, db *gorm.DB, id uint64) (*ProductCategory, error) {
	var category ProductCategory
	err := db.WithContext(ctx).Where("id = ?", id).First(&category).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

// GetCategoriesByIds 批量获取分类
func GetCategoriesByIds(ctx context.Context, db *gorm.DB, ids []uint64) ([]*ProductCategory, error) {
	if len(ids) == 0 {
		return []*ProductCategory{}, nil
	}
	var categories []*ProductCategory
	err := db.WithContext(ctx).Where("id IN ?", ids).Find(&categories).Error
	return categories, err
}
