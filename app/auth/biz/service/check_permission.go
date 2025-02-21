package service

import (
	"context"

	"github.com/PiaoAdmin/gomall/app/auth/biz/dal/model"
	"github.com/PiaoAdmin/gomall/app/auth/biz/dal/mysql"
	"github.com/PiaoAdmin/gomall/app/auth/biz/utils"
	auth "github.com/PiaoAdmin/gomall/rpc_gen/kitex_gen/auth"
)

type CheckPermissionService struct {
	ctx context.Context
} // NewCheckPermissionService new CheckPermissionService
func NewCheckPermissionService(ctx context.Context) *CheckPermissionService {
	return &CheckPermissionService{ctx: ctx}
}

// Run create note info
func (s *CheckPermissionService) Run(req *auth.PermissionCheckRequest) (resp *auth.PermissionCheckResponse, err error) {
	// Finish your business logic.
	claims, err := utils.NewARJWT().ParseAccessToken(req.Token)
	if err != nil {
		return &auth.PermissionCheckResponse{
			HasPermission: false,
			Msg:           err.Error(),
		}, nil
	}

	// check permission
	hasPermission, err := model.HasUserPermission(mysql.DB, s.ctx, claims.UserID, req.ServiceName, req.ResourceName)
	if err != nil {
		return &auth.PermissionCheckResponse{
			HasPermission: false,
			Msg:           "查询权限失败: " + err.Error(),
		}, err
	}
	return &auth.PermissionCheckResponse{
		HasPermission: hasPermission,
		Msg:           "查询权限成功",
	}, nil
}
