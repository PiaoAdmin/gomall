package service

import (
	"context"

	auth "github.com/PiaoAdmin/gomall/app/hertz_gateway/hertz_gen/hertz_gateway/auth"
	"github.com/PiaoAdmin/gomall/app/hertz_gateway/infra/rpc"
	rpcauth "github.com/PiaoAdmin/gomall/rpc_gen/kitex_gen/auth"
	"github.com/cloudwego/hertz/pkg/app"
)

type LoginService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewLoginService(Context context.Context, RequestContext *app.RequestContext) *LoginService {
	return &LoginService{RequestContext: RequestContext, Context: Context}
}

func (h *LoginService) Run(req *auth.LoginRequest) (resp *auth.LoginResponse, err error) {
	// 调用后端 RPC 更新基础用户信息
	res, err := rpc.AuthClient.IssueToken(h.Context, &rpcauth.IssueTokenRequest{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		return &auth.LoginResponse{
			Success: false,
			Msg:     err.Error(),
		}, err
	}

	// 构造响应
	resp = &auth.LoginResponse{
		Token:        res.Token,
		RefreshToken: res.RefreshToken,
		Success:      res.Success,
		Msg:          res.Msg,
	}
	return resp, nil
}
