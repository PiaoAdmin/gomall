/*
 * @Author: liaosijie
 * @Date: 2025-02-18 16:47:51
 * @Last Modified by: liaosijie
 * @Last Modified time: 2025-02-18 17:05:14
 */

package model

import (
	"time"

	"gorm.io/gorm"
)

type SnowflakeBase struct {
    ID        int64 `gorm:"primarykey;autoIncrement:false"`
    CreatedAt time.Time
    UpdatedAt time.Time
    DeletedAt gorm.DeletedAt `gorm:"index"`
}

type Consignee struct {
	Email	string
	StreetAddress	string
	City	string
	State	string
	Country	string
	Zipcode	string
}

type Order struct {
    SnowflakeBase
    Orderid     string `gorm:"type:varchar(100);uniqueIndex"`
    Userid      uint   `gorm:"type:varchar(11)"`
    UserCurrency string `gorm:"type:varchar(10)"`
    Consignee   Consignee `gorm:"embedded"`
    OrderItem   []OrderItem `gorm:"foreignKey:OrderIdRefer;references:ID"`
}

func (Order) TableName() string {
    return "order"
}