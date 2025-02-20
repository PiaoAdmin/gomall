package service

import (
	"context"

	"github.com/PiaoAdmin/gomall/app/auth/biz/dal/model"
	"github.com/PiaoAdmin/gomall/app/auth/biz/dal/mysql"
	auth "github.com/PiaoAdmin/gomall/rpc_gen/kitex_gen/auth"
)

type RemoveUserRoleService struct {
	ctx context.Context
} // NewRemoveUserRoleService new RemoveUserRoleService
func NewRemoveUserRoleService(ctx context.Context) *RemoveUserRoleService {
	return &RemoveUserRoleService{ctx: ctx}
}

// Run create note info
func (s *RemoveUserRoleService) Run(req *auth.RemoveUserRoleRequest) (resp *auth.RemoveUserRoleResponse, err error) {
	// Finish your business logic.
	if err = mysql.DB.WithContext(s.ctx).
		Where("user_id = ? AND role_code = ?", req.UserId, req.RoleCode).
		Delete(&model.UserRole{}).Error; err != nil {
		return &auth.RemoveUserRoleResponse{
			Success: false,
			Msg:     "删除用户角色失败: " + err.Error(),
		}, err
	}
	return &auth.RemoveUserRoleResponse{
		Success: true,
		Msg:     "删除用户角色成功",
	}, nil
}
