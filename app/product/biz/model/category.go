/**
 * @Author: ZhangHaoChen
 * @Date:   2/20/25 AM11:18
 */

package model

import (
	"context"
	"gorm.io/gorm"
)

type Category struct {
	CustomModel
	Name     string `gorm:"column:name" json:"Name"`
	ParentId int64  `gorm:"column:parent_id" json:"ParentId"`
	Status   int    `gorm:"column:status" json:"Status"`
}

func (c Category) tableName() string {
	return "category"
}

type CategoryQuery struct {
	ctx context.Context
	db  *gorm.DB
}

func NewCategoryQuery(ctx context.Context, db *gorm.DB) *CategoryQuery {
	return &CategoryQuery{ctx: ctx, db: db}
}
func (c CategoryQuery) GetProductsByCategoryName(categoryName string) (products []*Product, err error) {
	var categories []*Category
	err = c.db.Where("name = ?", categoryName).Find(&categories).Error
	if err != nil {
		return nil, err
	}
	if len(categories) == 0 {
		return nil, nil
	}
	var categoryIds []int64
	for _, category := range categories {
		categoryIds = append(categoryIds, category.ID)
	}
	err = c.db.Where("category_id in ?", categoryIds).Find(&products).Error
	return
}
