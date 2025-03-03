/**
 * @Author: ZhangHaoChen
 * @Date:   2/20/25 AM11:18
 */

package model

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

type Category struct {
	Base
	Name     string     `gorm:"column:name" json:"Name"`
	ParentId int64      `gorm:"column:parent_id" json:"ParentId"`
	Status   int        `gorm:"column:status" json:"Status"`
	Products []*Product `gorm:"many2many:prod_category" json:"Products"`
}

func (c Category) TableName() string {
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
	err = c.db.Joins("JOIN prod_category ON prod.id = prod_category.product_id").
		Where("prod_category.category_id IN ?", categoryIds).
		Find(&products).Error
	return
}

// CreateCategory 创建一个Category
func CreateCategory(DB *gorm.DB, category *Category) error {
	return DB.Create(&category).Error
}

// AssociateProductWithCategory 根据productId和categoryId，将product添加到对应的category中
func AssociateProductWithCategory(DB *gorm.DB, productId int64, categoryId int64) error {
	// TODO: NND 怎么都不对
	return DB.Transaction(func(tx *gorm.DB) error {
		var product Product
		if err := tx.Preload("Categories").First(&product, productId).Error; err != nil {
			return err
		}

		for _, c := range product.Categories {
			if c.ID == categoryId {
				return nil
			}
		}
		fmt.Printf("categoryId: %+v\n", categoryId)
		var category Category
		if err := tx.First(&category, categoryId).Error; err != nil {
			return err
		}
		fmt.Printf("category: %+v\n", category)
		return tx.Model(&product).Association("Categories").Append(&category)
	})
}
