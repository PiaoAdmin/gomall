package service

import (
	"context"

	auth "github.com/PiaoAdmin/gomall/app/hertz_gateway/hertz_gen/hertz_gateway/auth"
	"github.com/PiaoAdmin/gomall/app/hertz_gateway/infra/rpc"
	rpcauth "github.com/PiaoAdmin/gomall/rpc_gen/kitex_gen/auth"
	"github.com/cloudwego/hertz/pkg/app"
)

type AddUserRoleService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewAddUserRoleService(Context context.Context, RequestContext *app.RequestContext) *AddUserRoleService {
	return &AddUserRoleService{RequestContext: RequestContext, Context: Context}
}

func (h *AddUserRoleService) Run(req *auth.AddUserRoleRequest) (resp *auth.AddUserRoleResponse, err error) {
	res, err := rpc.AuthClient.AddUserRole(h.Context, &rpcauth.AddUserRoleRequest{
		UserId:   req.UserId,
		RoleCode: req.RoleCode,
		RoleName: req.RoleName,
	})
	if err != nil {
		return &auth.AddUserRoleResponse{
			Success: false,
			Msg:     "添加用户角色失败: " + err.Error(),
		}, err
	}

	return &auth.AddUserRoleResponse{
		Success: res.Success,
		Msg:     res.Msg,
	}, nil
}
