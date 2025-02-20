package service

import (
	"context"

	auth "github.com/PiaoAdmin/gomall/app/hertz_gateway/hertz_gen/hertz_gateway/auth"
	"github.com/PiaoAdmin/gomall/app/hertz_gateway/infra/rpc"
	rpcauth "github.com/PiaoAdmin/gomall/rpc_gen/kitex_gen/auth"
	"github.com/cloudwego/hertz/pkg/app"
)

type CheckPermissionService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewCheckPermissionService(Context context.Context, RequestContext *app.RequestContext) *CheckPermissionService {
	return &CheckPermissionService{RequestContext: RequestContext, Context: Context}
}

func (h *CheckPermissionService) Run(req *auth.PermissionCheckRequest) (resp *auth.PermissionCheckResponse, err error) {
	rpcResp, err := rpc.AuthClient.CheckPermission(h.Context, &rpcauth.PermissionCheckRequest{
		Token:        req.Token,
		ServiceName:  req.ServiceName,
		ResourceName: req.ResourceName,
	})
	if err != nil {
		return &auth.PermissionCheckResponse{
			HasPermission: false,
			Msg:           "权限校验失败: " + err.Error(),
		}, err
	}

	return &auth.PermissionCheckResponse{
		HasPermission: rpcResp.HasPermission,
		Msg:           rpcResp.Msg,
	}, nil
}
