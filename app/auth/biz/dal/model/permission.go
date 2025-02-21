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
