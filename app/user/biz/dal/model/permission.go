package model

type Permission struct {
	Base                  // 继承公共字段 ID自增
	PermissionCode int    `gorm:"unique;not null"`  // 一个服务对应一个唯一的权限码
	ServiceName    string `gorm:"size:50;not null"` // 微服务名称 如user_service
	ResourceName   string `gorm:"size:50;not null"` // 资源名称 如create_user
	IsPublic       bool   `gorm:"default:false"`    // 是否公共权限，公共权限不需要校验
	Desc           string `gorm:"size:255"`         // 描述
}

func (p Permission) TableName() string {
	return "permission"
}

type Role struct {
	Base               // 继承公共字段 ID自增
	RoleCode    int    `gorm:"unique;not null"`          // 角色编码，100\200\300等
	RoleName    string `gorm:"size:100;not null;unique"` // 角色名称，如"admin", "user"等
	Description string `gorm:"size:255"`                 // 角色描述
}

func (r Role) TableName() string {
	return "role"
}

type UserRole struct {
	Base         // 继承公共字段 ID自增
	UserID   int `gorm:"not null"` // 用户ID
	RoleCode int `gorm:"not null"` // 角色ID
}

func (ur UserRole) TableName() string {
	return "user_role"
}

type RolePermission struct {
	Base                  // 继承公共字段 ID自增
	RoleCode       int    `gorm:"not null"` // 角色ID
	PermissionCode string `gorm:"not null"` // 权限码
}

func (rp RolePermission) TableName() string {
	return "role_permission"
}
