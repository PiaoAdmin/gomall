/**
 * @Author: ZhangHaoChen
 * @Date:   2/20/25 AM11:25
 */

package model

type Sku struct {
	Base
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
