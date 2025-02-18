package service

import (
	"context"

	user "github.com/PiaoAdmin/gomall/app/hertz_gateway/hertz_gen/hertz_gateway/user"
	"github.com/PiaoAdmin/gomall/app/hertz_gateway/infra/rpc"
	rpcuser "github.com/PiaoAdmin/gomall/rpc_gen/kitex_gen/user"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/jinzhu/copier"
)

type GetUserInfoService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewGetUserInfoService(Context context.Context, RequestContext *app.RequestContext) *GetUserInfoService {
	return &GetUserInfoService{RequestContext: RequestContext, Context: Context}
}

func (h *GetUserInfoService) Run(req *user.GetUserInfoRequest) (resp *user.GetUserInfoResponse, err error) {
	// 调用后端 RPC 获取用户信息
	res, err := rpc.UserClient.GetUserInfo(h.Context, &rpcuser.GetUserInfoRequest{
		UserId: req.UserId,
	})
	if err != nil {
		return nil, err
	}
	// 构造响应
	resp = &user.GetUserInfoResponse{
		BaseUser: &user.BaseUser{},
	}
	if err := copier.Copy(resp.BaseUser, res.BaseUser); err != nil {
		return nil, err
	}
	return resp, nil
}
