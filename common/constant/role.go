package constant

type Role struct {
	RoleCode int32
	RoleName string
}

var (
	Admin    = Role{RoleCode: 300, RoleName: "admin"}
	User     = Role{RoleCode: 100, RoleName: "user"}
	Merchant = Role{RoleCode: 200, RoleName: "merchant"}
)
