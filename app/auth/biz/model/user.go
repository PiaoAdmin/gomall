package model

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type User struct {
	SnowflakeBase           // 继承公共字段 使用雪花算法生成ID
	Username      string    `gorm:"unique;size:50;not null"`
	Nickname      string    `gorm:"size:50"`
	Avatar        string    `gorm:"type:varchar(512)"`
	Phone         string    `gorm:"unique;size:20;not null;index"`
	Email         string    `gorm:"unique;size:80;index"`
	Password      string    `gorm:"size:256;not null"`
	Gender        int8      `gorm:"type:tinyint;default:2"`
	BirthDate     time.Time `gorm:"type:date"`
	Status        int32     `gorm:"default:1"`                       // 1 表示正常，0 表示禁用，2 表示未激活等
	Balance       float64   `gorm:"type:decimal(10,2);default:0.00"` // 余额字段，精确到小数点后两位
}

func (u User) TableName() string {
	return "user"
}

func GetUserByUsername(db *gorm.DB, ctx context.Context, username string) (user *User, err error) {
	err = db.WithContext(ctx).Model(&User{}).Where("username = ?", username).First(&user).Error
	return
}
