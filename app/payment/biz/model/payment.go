package model

import (
	"context"
	"time"

	"gorm.io/gorm"
)

// 数据库dao层，但数据定义也在这个部分
type PaymentLog struct {
	gorm.Model
	UserId        int64
	OrderId       int64
	TranscationId int64
	Amount        float32
	PayAt         time.Time
	State         int32
}

// 实现TableName方法用于指定表明，否则gorm本身会默认使用payment_logs作为表名
func (pl PaymentLog) TableName() string {
	return "payment"
}

// 获得数据库连接并执行相应的操作
// 插入一条数据
func CreatePaymentLog(db *gorm.DB, ctx context.Context, payment *PaymentLog) error {
	return db.WithContext(ctx).Model(&PaymentLog{}).Create(payment).Error
}
