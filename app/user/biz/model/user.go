package model

import (
	"context"

	"github.com/PiaoAdmin/pmall/common/uniqueid"
	"gorm.io/gorm"
)

type Status int32

const (
	StatusOK      = iota // 0
	StatusBAN            // 1
	StatusDELETED        // 2
)

type User struct {
	Model
	Username string `gorm:"uniqueIndex;type:varchar(100);not null"`
	Password string `gorm:"type:varchar(255);not null"`
	Email    string `gorm:"uniqueIndex;type:varchar(100);not null"`
	Phone    string `gorm:"type:varchar(20)"`
	Avatar   string `gorm:"type:varchar(255)"`
	Status   int32  `gorm:"default:1"` // 1: active, 0: inactive
}

func (User) TableName() string {
	return "user"
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	// 雪花算法生成id
	u.ID = uint64(uniqueid.GenId())
	return
}

// Create 创建用户
func Create(ctx context.Context, db *gorm.DB, user *User) error {
	return db.WithContext(ctx).Create(user).Error
}

// GetByUsername 根据用户名获取用户
func GetByUsername(ctx context.Context, db *gorm.DB, username string) (*User, error) {
	var user User
	err := db.WithContext(ctx).Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByEmail 根据邮箱获取用户
func GetByEmail(ctx context.Context, db *gorm.DB, email string) (*User, error) {
	var user User
	err := db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetById 根据ID获取用户
func GetById(ctx context.Context, db *gorm.DB, id uint64) (*User, error) {
	var user User
	err := db.WithContext(ctx).Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Update 更新用户信息
func Update(ctx context.Context, db *gorm.DB, user *User) error {
	return db.WithContext(ctx).Save(user).Error
}
