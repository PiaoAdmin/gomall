package service

import (
	"context"
	"fmt"

	auth "github.com/PiaoAdmin/gomall/app/hertz_gateway/hertz_gen/hertz_gateway/auth"
	"github.com/PiaoAdmin/gomall/app/hertz_gateway/infra/rpc"
	rpcauth "github.com/PiaoAdmin/gomall/rpc_gen/kitex_gen/auth"
	"github.com/cloudwego/hertz/pkg/app"
)

type GetUserRoleService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewGetUserRoleService(Context context.Context, RequestContext *app.RequestContext) *GetUserRoleService {
	return &GetUserRoleService{RequestContext: RequestContext, Context: Context}
}

func (h *GetUserRoleService) Run(req *auth.GetUserRoleRequest) (resp *auth.GetUserRoleResponse, err error) {
	res, err := rpc.AuthClient.GetUserRole(h.Context, &rpcauth.GetUserRoleRequest{
		UserId: req.UserId,
	})
	fmt.Printf("req: %v\n", req)
	fmt.Printf("____________________")
	if err != nil {
		return nil, err
	}
	return &auth.GetUserRoleResponse{
		RoleCode: res.RoleCode,
		RoleName: res.RoleName,
	}, nil
}
