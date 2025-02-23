package service

import (
	"context"

	"github.com/PiaoAdmin/gomall/app/auth/biz/dal/model"
	"github.com/PiaoAdmin/gomall/app/auth/biz/dal/mysql"
	auth "github.com/PiaoAdmin/gomall/rpc_gen/kitex_gen/auth"
)

type AddUserRoleService struct {
	ctx context.Context
} // NewAddUserRoleService new AddUserRoleService
func NewAddUserRoleService(ctx context.Context) *AddUserRoleService {
	return &AddUserRoleService{ctx: ctx}
}

// Run create note info
func (s *AddUserRoleService) Run(req *auth.AddUserRoleRequest) (resp *auth.AddUserRoleResponse, err error) {
	// Finish your business logic.
	ur := model.UserRole{
		UserID:   req.UserId,
		RoleCode: int(req.RoleCode),
	}

	if err = mysql.DB.WithContext(s.ctx).Create(&ur).Error; err != nil {
		return &auth.AddUserRoleResponse{
			Success: false,
			Msg:     "添加用户角色失败: " + err.Error(),
		}, err
	}
	return &auth.AddUserRoleResponse{
		Success: true,
		Msg:     "添加用户角色成功",
	}, nil
}
