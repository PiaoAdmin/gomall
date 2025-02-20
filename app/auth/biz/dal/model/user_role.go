package model

type UserRole struct {
	Base           // 继承公共字段 ID自增
	UserID   int64 `gorm:"not null"` // 用户ID
	RoleCode int   `gorm:"not null"` // 角色ID
}

func (ur UserRole) TableName() string {
	return "user_role"
}
