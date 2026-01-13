package service

import (
	"context"

	"github.com/PiaoAdmin/pmall/app/api/biz/model/api/auth"
	"github.com/PiaoAdmin/pmall/app/api/md/jwt"
	"github.com/PiaoAdmin/pmall/app/api/rpc"
	"github.com/PiaoAdmin/pmall/rpc_gen/user"
	"github.com/cloudwego/hertz/pkg/app"
)

type GetUserInfoService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewGetUserInfoService(Context context.Context, RequestContext *app.RequestContext) *GetUserInfoService {
	return &GetUserInfoService{
		RequestContext: RequestContext,
		Context:        Context,
	}
}

func (h *GetUserInfoService) Run(req *auth.GetUserInfoReq) (resp *auth.GetUserInfoResp, err error) {
	claims := jwt.ExtractClaims(h.Context, h.RequestContext)
	userId := uint64(claims[jwt.JwtMiddleware.IdentityKey].(float64))

	rpcResp, err := rpc.UserClient.GetUserInfo(h.Context, &user.GetUserInfoRequest{
		UserId: userId,
	})
	if err != nil {
		return nil, err
	}

	return &auth.GetUserInfoResp{
		User: &auth.UserInfo{
			Id:       rpcResp.User.Id,
			Username: rpcResp.User.Username,
			Email:    rpcResp.User.Email,
			Phone:    rpcResp.User.Phone,
			Avatar:   rpcResp.User.Avatar,
			Status:   rpcResp.User.Status,
		},
	}, nil
}
