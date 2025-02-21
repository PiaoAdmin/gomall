package service

import (
	"context"

	"github.com/PiaoAdmin/gomall/app/auth/biz/dal/model"
	"github.com/PiaoAdmin/gomall/app/auth/biz/dal/mysql"
	auth "github.com/PiaoAdmin/gomall/rpc_gen/kitex_gen/auth"
)

type GetUserRoleService struct {
	ctx context.Context
} // NewGetUserRoleService new GetUserRoleService
func NewGetUserRoleService(ctx context.Context) *GetUserRoleService {
	return &GetUserRoleService{ctx: ctx}
}

// Run create note info
func (s *GetUserRoleService) Run(req *auth.GetUserRoleRequest) (resp *auth.GetUserRoleResponse, err error) {
	// Finish your business logic.
	var roles []model.Role
	err = mysql.DB.WithContext(s.ctx).
		Table("user_role").
		Joins("JOIN role ON user_role.role_code = role.role_code").
		Where("user_role.user_id = ?", req.UserId).
		Select("role.*").
		Find(&roles).Error
	if err != nil {
		return
	}

	// 将查询结果转换为 RPC 层定义的角色结构，此处以简单映射为例
	resp = &auth.GetUserRoleResponse{}
	for _, r := range roles {
		resp.RoleCode = append(resp.RoleCode, int32(r.RoleCode))
		resp.RoleName = append(resp.RoleName, r.RoleName)
	}
	return resp, nil
}
