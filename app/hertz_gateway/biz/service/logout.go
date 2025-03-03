package service

import (
	"context"

	"github.com/PiaoAdmin/gomall/app/hertz_gateway/biz/utils"
	auth "github.com/PiaoAdmin/gomall/app/hertz_gateway/hertz_gen/hertz_gateway/auth"
	"github.com/PiaoAdmin/gomall/app/hertz_gateway/hertz_gen/hertz_gateway/common"
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

func (h *LogoutService) Run(req *common.Empty) (resp *auth.LogoutResponse, err error) {
	token, err := utils.GetToken(h.Context, h.RequestContext)
	if err != nil {
		return
	}
	refreshToken, err := utils.GetRefreshToken(h.Context, h.RequestContext)
	if err != nil {
		return
	}
	res, err := rpc.AuthClient.Logout(h.Context, &rpcauth.LogoutRequest{
		Token:        token,
		RefreshToken: refreshToken,
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
