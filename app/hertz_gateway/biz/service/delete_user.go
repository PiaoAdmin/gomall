package service

import (
	"context"

	user "github.com/PiaoAdmin/gomall/app/hertz_gateway/hertz_gen/hertz_gateway/user"
	"github.com/PiaoAdmin/gomall/app/hertz_gateway/infra/rpc"
	rpcuser "github.com/PiaoAdmin/gomall/rpc_gen/kitex_gen/user"
	"github.com/cloudwego/hertz/pkg/app"
)

type DeleteUserService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewDeleteUserService(Context context.Context, RequestContext *app.RequestContext) *DeleteUserService {
	return &DeleteUserService{RequestContext: RequestContext, Context: Context}
}

func (h *DeleteUserService) Run(req *user.DeleteUserRequest) (resp *user.DeleteUserResponse, err error) {
	// 调用后端 RPC 删除用户
	res, err := rpc.UserClient.DeleteUser(h.Context, &rpcuser.DeleteUserRequest{
		UserId: req.UserId,
	})
	if err != nil {
		return nil, err
	}

	// 构造响应
	resp = &user.DeleteUserResponse{
		Success: res.Success,
		Msg:     res.Msg,
	}
	return
}
