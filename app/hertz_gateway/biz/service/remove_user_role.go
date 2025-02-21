package service

import (
	"context"

	auth "github.com/PiaoAdmin/gomall/app/hertz_gateway/hertz_gen/hertz_gateway/auth"
	"github.com/PiaoAdmin/gomall/app/hertz_gateway/infra/rpc"
	rpcauth "github.com/PiaoAdmin/gomall/rpc_gen/kitex_gen/auth"
	"github.com/cloudwego/hertz/pkg/app"
)

type RemoveUserRoleService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewRemoveUserRoleService(Context context.Context, RequestContext *app.RequestContext) *RemoveUserRoleService {
	return &RemoveUserRoleService{RequestContext: RequestContext, Context: Context}
}

func (h *RemoveUserRoleService) Run(req *auth.RemoveUserRoleRequest) (resp *auth.RemoveUserRoleResponse, err error) {
	rpcResp, err := rpc.AuthClient.RemoveUserRole(h.Context, &rpcauth.RemoveUserRoleRequest{
		UserId:   req.UserId,
		RoleCode: req.RoleCode,
		RoleName: req.RoleName,
	})
	if err != nil {
		return &auth.RemoveUserRoleResponse{
			Success: false,
			Msg:     "删除用户角色失败: " + err.Error(),
		}, err
	}

	return &auth.RemoveUserRoleResponse{
		Success: rpcResp.Success,
		Msg:     rpcResp.Msg,
	}, nil
}
