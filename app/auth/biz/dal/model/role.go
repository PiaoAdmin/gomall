package model

type Role struct {
	Base               // 继承公共字段 ID自增
	RoleCode    int    `gorm:"unique;not null"`          // 角色编码，100\200\300等
	RoleName    string `gorm:"size:100;not null;unique"` // 角色名称，如"admin", "user"等
	Description string `gorm:"size:255"`                 // 角色描述
	User        []User `gorm:"many2many:user_role;foreignKey:RoleCode;joinForeignKey:RoleCode;References:ID;joinReferences:UserID"`
}

func (r Role) TableName() string {
	return "role"
}

type UserRole struct {
	Base           // 继承公共字段 ID自增
	UserID   int64 `gorm:"not null"` // 用户ID
	RoleCode int   `gorm:"not null"` // 角色ID
}

func (ur UserRole) TableName() string {
	return "user_role"
}
