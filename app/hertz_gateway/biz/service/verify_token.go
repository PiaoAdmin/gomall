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

type VerifyTokenService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewVerifyTokenService(Context context.Context, RequestContext *app.RequestContext) *VerifyTokenService {
	return &VerifyTokenService{RequestContext: RequestContext, Context: Context}
}

func (h *VerifyTokenService) Run(req *common.Empty) (resp *auth.VerifyTokenResponse, err error) {
	token, err := utils.GetToken(h.Context, h.RequestContext)
	if err != nil {
		return
	}
	res, err := rpc.AuthClient.VerifyToken(h.Context, &rpcauth.VerifyTokenRequest{
		Token: string(token),
	})
	if err != nil {
		return &auth.VerifyTokenResponse{
			IsValid: false,
			Msg:     err.Error(),
		}, err
	}

	return &auth.VerifyTokenResponse{
		IsValid: res.IsValid,
		Msg:     res.Msg,
	}, nil
}
