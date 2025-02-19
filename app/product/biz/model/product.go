/**
 * @Author: ZhangHaoChen
 * @Date:   2/19/25 PM2:36
 */

package model

import (
	"gorm.io/gorm"
	"time"
)

type Product struct {
	CustomModel
	ProdName        string    `gorm:"column:prod_name" json:"ProdName"`               //type:string       comment:商品名称                                   version:2025-01-19 17:09
	ShopId          int64     `gorm:"column:shop_id" json:"ShopId"`                   //type:BIGINT       comment:商铺ID                                     version:2025-01-19 17:09
	OriPrice        float64   `gorm:"column:ori_price" json:"OriPrice"`               //type:*float64     comment:商品原价                                   version:2025-01-19 17:09
	Price           float64   `gorm:"column:price" json:"Price"`                      //type:*float64     comment:商品现价                                   version:2025-01-19 17:09
	Brief           string    `gorm:"column:brief" json:"Brief"`                      //type:string       comment:商品概述                                   version:2025-01-19 17:09
	Content         string    `gorm:"column:content" json:"Content"`                  //type:string       comment:商品详情                                   version:2025-01-19 17:09
	MainImage       string    `gorm:"column:main_image" json:"MainImage"`             //type:string       comment:商品主图;在首页展示的图片                  version:2025-01-19 17:09
	SecondaryImages string    `gorm:"column:secondary_images" json:"SecondaryImages"` //type:string       comment:商品辅图;点开商品后的辅图列表，用','隔开   version:2025-01-19 17:09
	Status          int       `gorm:"column:status" json:"Status"`                    //type:*int         comment:商品状态;-1:删除 0:下架 1:上架             version:2025-01-19 17:09
	SoldNum         int       `gorm:"column:sold_num" json:"SoldNum"`                 //type:*int         comment:销量                                       version:2025-01-19 17:09
	TotalStocks     int       `gorm:"column:total_stocks" json:"TotalStocks"`         //type:*int         comment:库存                                       version:2025-01-19 17:09
	ListingTime     time.Time `gorm:"column:listing_time" json:"ListingTime"`         //type:*time.Time   comment:商品上架时间                               version:2025-01-19 17:09
}
type Sku struct {
	CustomModel
	ProdId   int64    `gorm:"column:prod_id" json:"ProdId"`     //type:BIGINT     comment:sku对应商品的id    version:2025-01-19 17:12
	OriPrice *float64 `gorm:"column:ori_price" json:"OriPrice"` //type:*float64   comment:原价               version:2025-01-19 17:12
	Price    *float64 `gorm:"column:price" json:"Price"`        //type:*float64   comment:现价               version:2025-01-19 17:12
	Stocks   *int     `gorm:"column:stocks" json:"Stocks"`      //type:*int       comment:库存               version:2025-01-19 17:12
	Img      string   `gorm:"column:img" json:"Img"`            //type:string     comment:sku的图片          version:2025-01-19 17:12
	SkuName  string   `gorm:"column:sku_name" json:"SkuName"`   //type:string     comment:sku名称            version:2025-01-19 17:12
	Status   *int     `gorm:"column:status" json:"Status"`      //type:*int       comment:sku状态            version:2025-01-19 17:12
}

// TableName 表名:sku，商品sku表。
// 说明:
func (Sku) TableName() string {
	return "sku"
}

// TableName 表名:prod，商品表。
// 说明:
func (Product) TableName() string {
	return "prod"
}

func CreateProduct(DB *gorm.DB, product *Product) error {
	return DB.Create(&product).Error
}
