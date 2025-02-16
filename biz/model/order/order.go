/*
 * @Author: liaosijie
 * @Date: 2025-02-16 10:14:49
 * @Last Modified by: liaosijie
 * @Last Modified time: 2025-02-16 17:50:08
 */

package order

import (
	"time"
)
 
 type SnowflakeBase struct {
	 ID        int64     `gorm:"primarykey"`
	 CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	 UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	 DeletedAt time.Time `gorm:"index"`
 }
 
 type Order struct {
	 SnowflakeBase
	 UserID       int64    `gorm:"not null"`
	 UserCurrency string   `gorm:"size:10;not null"`
	 AddressID    int64    `gorm:"not null"`
	 Email        string   `gorm:"size:255;not null"`
	 Status       string   `gorm:"size:20;not null"`
	 OrderItems   []*OrderItem
 }
 
 type OrderItem struct {
	 SnowflakeBase
	 OrderID  int64 `gorm:"not null"`
	 ProductID int64 `gorm:"not null"`
	 Quantity int   `gorm:"not null"`
	 Price    float64 `gorm:"not null"`
 }
 
 type Address struct {
	 SnowflakeBase
	 UserID       int64  `gorm:"not null"`
	 StreetAddress string `gorm:"size:255;not null"`
	 City         string `gorm:"size:100;not null"`
	 State        string `gorm:"size:100;not null"`
	 Country      string `gorm:"size:100;not null"`
	 ZipCode      int    `gorm:"not null"`
 }