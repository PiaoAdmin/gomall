package service

import (
	"context"

	auth "github.com/PiaoAdmin/gomall/app/hertz_gateway/hertz_gen/hertz_gateway/auth"
	"github.com/PiaoAdmin/gomall/app/hertz_gateway/infra/rpc"
	rpcauth "github.com/PiaoAdmin/gomall/rpc_gen/kitex_gen/auth"
	"github.com/cloudwego/hertz/pkg/app"
)

type LogoutService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewLogoutService(Context context.Context, RequestContext *app.RequestContext) *LogoutService {
	return &LogoutService{RequestContext: RequestContext, Context: Context}
}

func (h *LogoutService) Run(req *auth.LogoutRequest) (resp *auth.LogoutResponse, err error) {
	res, err := rpc.AuthClient.Logout(h.Context, &rpcauth.LogoutRequest{
		Token:        req.Token,
		RefreshToken: req.RefreshToken,
	})
	if err != nil {
		return &auth.LogoutResponse{
			Success: false,
			Msg:     err.Error(),
		}, err
	}

	return &auth.LogoutResponse{
		Success: res.Success,
		Msg:     res.Msg,
	}, nil
}
