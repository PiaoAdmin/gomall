/*
 * @Author: liaosijie
 * @Date: 2025-02-18 16:47:44
 * @Last Modified by: liaosijie
 * @Last Modified time: 2025-02-18 17:05:00
 */

package model

// import (
// 	"gorm.io/gorm"
// 	"time"
// )

// type SnowflakeBase struct {
//     ID        int64 `gorm:"primarykey;autoIncrement:false"`
//     CreatedAt time.Time
//     UpdatedAt time.Time
//     DeletedAt gorm.DeletedAt `gorm:"index"`
// }

type OrderItem struct {
    SnowflakeBase
    ProductID    uint32 `gorm:"type:int(11)"`
    OrderIdRefer int64  `gorm:"type:bigint;index"`
    Quantity     uint32 `gorm:"type:int(11)"`
    Cost         float64 `gorm:"type:decimal(10,2)"`
}