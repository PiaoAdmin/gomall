package model

type RolePermission struct {
	Base                  // 继承公共字段 ID自增
	RoleCode       int    `gorm:"not null"` // 角色ID
	PermissionCode string `gorm:"not null"` // 权限码
}

func (rp RolePermission) TableName() string {
	return "role_permission"
}
