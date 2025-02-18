package service

import (
	"context"

	user "github.com/PiaoAdmin/gomall/app/hertz_gateway/hertz_gen/hertz_gateway/user"
	"github.com/PiaoAdmin/gomall/app/hertz_gateway/infra/rpc"
	rpcuser "github.com/PiaoAdmin/gomall/rpc_gen/kitex_gen/user"
	"github.com/cloudwego/hertz/pkg/app"
)

type UpdateUserPasswordService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewUpdateUserPasswordService(Context context.Context, RequestContext *app.RequestContext) *UpdateUserPasswordService {
	return &UpdateUserPasswordService{RequestContext: RequestContext, Context: Context}
}

func (h *UpdateUserPasswordService) Run(req *user.UpdateUserPasswordRequest) (resp *user.UpdateUserPasswordResponse, err error) {
	// 调用后端 RPC 更新用户密码
	res, err := rpc.UserClient.UpdateUserPassword(h.Context, &rpcuser.UpdateUserPasswordRequest{
		UserId:      req.UserId,
		OldPassword: req.OldPassword,
		NewPassword: req.NewPassword,
	})
	if err != nil {
		return nil, err
	}
	// 构造响应
	resp = &user.UpdateUserPasswordResponse{
		Success: res.Success,
		Msg:     res.Msg,
	}
	return resp, nil
}
