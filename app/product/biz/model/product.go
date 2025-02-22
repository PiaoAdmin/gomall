/**
 * @Author: ZhangHaoChen
 * @Date:   2/19/25 PM2:36
 */

package model

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type Product struct {
	Base
	ProdName        string     `gorm:"column:prod_name" json:"ProdName"`               //type:string       comment:商品名称                                   version:2025-01-19 17:09
	ShopId          int64      `gorm:"column:shop_id" json:"ShopId"`                   //type:BIGINT       comment:商铺ID                                     version:2025-01-19 17:09
	OriPrice        float64    `gorm:"column:ori_price" json:"OriPrice"`               //type:*float64     comment:商品原价                                   version:2025-01-19 17:09
	Price           float64    `gorm:"column:price" json:"Price"`                      //type:*float64     comment:商品现价                                   version:2025-01-19 17:09
	Brief           string     `gorm:"column:brief" json:"Brief"`                      //type:string       comment:商品概述                                   version:2025-01-19 17:09
	Content         string     `gorm:"column:content" json:"Content"`                  //type:string       comment:商品详情                                   version:2025-01-19 17:09
	MainImage       string     `gorm:"column:main_image" json:"MainImage"`             //type:string       comment:商品主图;在首页展示的图片                  version:2025-01-19 17:09
	SecondaryImages string     `gorm:"column:secondary_images" json:"SecondaryImages"` //type:string       comment:商品辅图;点开商品后的辅图列表，用','隔开   version:2025-01-19 17:09
	Status          int        `gorm:"column:status" json:"Status"`                    //type:*int         comment:商品状态;-1:删除 0:下架 1:上架             version:2025-01-19 17:09
	SoldNum         int        `gorm:"column:sold_num" json:"SoldNum"`                 //type:*int         comment:销量                                       version:2025-01-19 17:09
	TotalStock      int        `gorm:"column:total_stocks" json:"TotalStock"`          //type:*int         comment:库存                                       version:2025-01-19 17:09
	ListingTime     time.Time  `gorm:"column:listing_time" json:"ListingTime"`         //type:*time.Time   comment:商品上架时间                               version:2025-01-19 17:09
	Categories      []Category `gorm:"many2many:prod_category" json:"Categories"`
}

// TableName 表名:prod，商品表。
// 说明:
func (Product) TableName() string {
	return "prod"
}

func CreateProduct(DB *gorm.DB, product *Product) error {
	return DB.Create(&product).Error
}

type ProductQuery struct {
	ctx context.Context
	db  *gorm.DB
}

func (p ProductQuery) GetById(productId int64) (product Product, err error) {
	err = p.db.WithContext(p.ctx).Model(&Product{}).First(&product, productId).Error
	return
}

func (p ProductQuery) SearchProducts(q string) (products []*Product, err error) {
	err = p.db.WithContext(p.ctx).Model(&Product{}).Find(&products, "prod_name LIKE ? or brief like ?", "%"+q+"%", "%"+q+"%").Error
	return
}

func NewProductQuery(ctx context.Context, db *gorm.DB) *ProductQuery {
	return &ProductQuery{
		ctx: ctx,
		db:  db,
	}
}
