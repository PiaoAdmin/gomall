package service

import (
	"context"

	"github.com/PiaoAdmin/gomall/app/hertz_gateway/biz/utils"
	auth "github.com/PiaoAdmin/gomall/app/hertz_gateway/hertz_gen/hertz_gateway/auth"
	common "github.com/PiaoAdmin/gomall/app/hertz_gateway/hertz_gen/hertz_gateway/common"
	"github.com/PiaoAdmin/gomall/app/hertz_gateway/infra/rpc"
	rpcauth "github.com/PiaoAdmin/gomall/rpc_gen/kitex_gen/auth"
	"github.com/cloudwego/hertz/pkg/app"
)

type RefreshTokenService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewRefreshTokenService(Context context.Context, RequestContext *app.RequestContext) *RefreshTokenService {
	return &RefreshTokenService{RequestContext: RequestContext, Context: Context}
}

func (h *RefreshTokenService) Run(req *common.Empty) (resp *auth.RefreshTokenResponse, err error) {
	token, err := utils.GetToken(h.Context, h.RequestContext)
	if err != nil {
		return
	}
	refreshToken, err := utils.GetRefreshToken(h.Context, h.RequestContext)
	if err != nil {
		return
	}

	res, err := rpc.AuthClient.RefreshToken(h.Context, &rpcauth.RefreshTokenRequest{
		Token:        token,
		RefreshToken: refreshToken,
	})
	if err != nil {
		return nil, err
	}

	return &auth.RefreshTokenResponse{
		IsValid:      res.IsValid,
		Msg:          res.Msg,
		Token:        res.Token,
		RefreshToken: res.RefreshToken,
	}, nil
}
