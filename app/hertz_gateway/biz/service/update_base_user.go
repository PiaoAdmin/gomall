package service

import (
	"context"

	user "github.com/PiaoAdmin/gomall/app/hertz_gateway/hertz_gen/hertz_gateway/user"
	"github.com/PiaoAdmin/gomall/app/hertz_gateway/infra/rpc"
	rpcuser "github.com/PiaoAdmin/gomall/rpc_gen/kitex_gen/user"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/jinzhu/copier"
)

type UpdateBaseUserService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewUpdateBaseUserService(Context context.Context, RequestContext *app.RequestContext) *UpdateBaseUserService {
	return &UpdateBaseUserService{RequestContext: RequestContext, Context: Context}
}

func (h *UpdateBaseUserService) Run(req *user.UpdateBaseUserRequest) (resp *user.UpdateBaseUserResponse, err error) {
	// 调用后端 RPC 更新基础用户信息
	tempUser := rpcuser.BaseUser{}
	if err := copier.Copy(&tempUser, req.BaseUser); err != nil {
		return nil, err
	}
	res, err := rpc.UserClient.UpdateBaseUser(h.Context, &rpcuser.UpdateBaseUserRequest{
		BaseUser: &tempUser,
	})
	if err != nil {
		return nil, err
	}

	// 构造响应
	resp = &user.UpdateBaseUserResponse{
		Success: res.Success,
		Msg:     res.Msg,
	}
	return
}
