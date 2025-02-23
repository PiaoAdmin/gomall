/*
 * @Author: liaosijie
 * @Date: 2025-02-18 16:47:44
 * @Last Modified by: liaosijie
 * @Last Modified time: 2025-02-18 22:14:53
 */

package model

// import (
// 	"gorm.io/gorm"
// 	"time"
// )

type OrderItem struct {
	SnowflakeBase
	ProductID    int64
	OrderIdRefer int64 `gorm:"index"`
	Quantity     int32
	Cost         float32
}
