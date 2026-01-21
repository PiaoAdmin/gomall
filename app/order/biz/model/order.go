package model

import (
	"strconv"

	"gorm.io/gorm"
)

const (
	OrderStatePlaced   string = "placed"
	OrderStatePaid     string = "paid"
	OrderStateCanceled string = "canceled"
)

type Address struct {
	Name          string `gorm:"column:shipping_name;type:varchar(64);not null;default:''"`
	StreetAddress string `gorm:"column:shipping_street_address;type:varchar(255);not null;default:''"`
	City          string `gorm:"column:shipping_city;type:varchar(64);not null;default:''"`
	ZipCode       int32  `gorm:"column:shipping_zip_code;type:int;not null;default:0"`
}

type Order struct {
	Model
	OrderId         string      `gorm:"primaryKey;column:order_id;type:varchar(64);not null"`
	UserId          uint64      `gorm:"column:user_id;type:bigint unsigned;index;not null"`
	Email           string      `gorm:"column:email;type:varchar(255);not null;default:''"`
	ShippingAddress Address     `gorm:"embedded"`
	Status          string      `gorm:"column:status;type:varchar(32);not null;default:''"`
	Items           []OrderItem `gorm:"foreignKey:OrderId;references:OrderId"`
}

func (Order) TableName() string {
	return "orders"
}

func (o *Order) BeforeCreate(tx *gorm.DB) error {
	id, _ := strconv.ParseInt(o.OrderId, 10, 64)
	o.ID = uint64(id)
	return nil
}
