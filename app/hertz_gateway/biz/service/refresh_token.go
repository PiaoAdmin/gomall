package service

import (
	"context"

	auth "github.com/PiaoAdmin/gomall/app/hertz_gateway/hertz_gen/hertz_gateway/auth"
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

func (h *RefreshTokenService) Run(req *auth.RefreshTokenRequest) (resp *auth.RefreshTokenResponse, err error) {
	res, err := rpc.AuthClient.RefreshToken(h.Context, &rpcauth.RefreshTokenRequest{
		Token:        req.Token,
		RefreshToken: req.RefreshToken,
	})
	if err != nil {
		return &auth.RefreshTokenResponse{
			IsValid:      false,
			Msg:          err.Error(),
			Token:        res.Token,        // 可能为空
			RefreshToken: res.RefreshToken, // 可能为空
		}, err
	}
	return &auth.RefreshTokenResponse{
		IsValid:      res.IsValid,
		Msg:          res.Msg,
		Token:        res.Token,
		RefreshToken: res.RefreshToken,
	}, nil
}
