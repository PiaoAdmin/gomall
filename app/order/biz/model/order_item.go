package model

type OrderItem struct {
	ID       uint64  `gorm:"primaryKey;autoIncrement"`
	OrderId  string  `gorm:"column:order_id;type:varchar(64);not null;index"`
	SkuId    uint64  `gorm:"column:sku_id;type:bigint unsigned;not null"`
	SkuName  string  `gorm:"column:sku_name;type:varchar(255);not null;default:''"`
	Price    float64 `gorm:"column:price;type:decimal(10,2);not null;default:0.00"`
	Quantity int32   `gorm:"column:quantity;type:int;not null;default:1"`
}

func (OrderItem) TableName() string {
	return "order_items"
}
