/*
 * @Author: liaosijie
 * @Date: 2025-02-18 16:47:51
 * @Last Modified by: liaosijie
 * @Last Modified time: 2025-02-18 23:32:44
 */

package model

import (
	"context"
	"log"

	"github.com/bwmarrin/snowflake"
	"gorm.io/gorm"
)

type Consignee struct {
	Email         string
	StreetAddress string
	City          string
	State         string
	Country       string
	ZipCode       int32
}

type OrderState string

const (
	OrderStatePlaced   OrderState = "placed"
	OrderStatePaid     OrderState = "paid"
	OrderStateCanceled OrderState = "canceled"
)

type Order struct {
	Base
	OrderId      int64 `gorm:"uniqueIndex;not null"`
	UserId       int64
	UserCurrency string
	Consignee    Consignee   `gorm:"embedded"`
	OrderItem    []OrderItem `gorm:"foreignKey:OrderIdRefer;references:ID"`
	OrderState   OrderState
}

func (Order) TableName() string {
	return "order"
}

// 雪花算法生成id
func CreateId(flag int64) (id int64) {
	node, err := snowflake.NewNode(flag) // flga表示节点ID，每个人不一样
	if err != nil {
		log.Fatalf("failed to create snowflake node: %v", err)
	}
	id = node.Generate().Int64()
	return
}

func ListOrder(db *gorm.DB, ctx context.Context, userId int64) ([]*Order, error) {
	var orders []*Order
	err := db.WithContext(ctx).Where("user_id = ?", userId).Preload("OrderItems").Find(&orders).Error
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func GetOrder(db *gorm.DB, ctx context.Context, userId int64, orderId int64) (order Order, err error) {
	err = db.Where(&Order{UserId: userId, OrderId: orderId}).First(&order).Error
	return
}

func UpdateOrderState(db *gorm.DB, ctx context.Context, userId int64, orderId int64, state OrderState) error {
	return db.Model(&Order{}).Where(&Order{UserId: userId, OrderId: orderId}).Update("order_state", state).Error
}
