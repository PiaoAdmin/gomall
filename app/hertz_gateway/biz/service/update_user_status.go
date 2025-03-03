package service

import (
	"context"
	"fmt"

	user "github.com/PiaoAdmin/gomall/app/hertz_gateway/hertz_gen/hertz_gateway/user"
	"github.com/PiaoAdmin/gomall/app/hertz_gateway/infra/rpc"
	rpcuser "github.com/PiaoAdmin/gomall/rpc_gen/kitex_gen/user"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/kitex/pkg/kerrors"
)

type UpdateUserStatusService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewUpdateUserStatusService(Context context.Context, RequestContext *app.RequestContext) *UpdateUserStatusService {
	return &UpdateUserStatusService{RequestContext: RequestContext, Context: Context}
}

func (h *UpdateUserStatusService) Run(req *user.UpdateUserStatusRequest) (resp *user.UpdateUserStatusResponse, err error) {
	// 调用后端 RPC 更新用户状态
	res, err := rpc.UserClient.UpdateUserStatus(h.Context, &rpcuser.UpdateUserStatusRequest{
		UserId: req.UserId,
		Status: req.Status,
	})
	if err != nil {
		// TODO:如何返回业务异常
		bizErr, isBizErr := kerrors.FromBizStatusError(err)
		if isBizErr {
			fmt.Printf("bizErr: %v\n", bizErr)
			return nil, bizErr
		}
		return nil, err
	}
	// 构造响应
	resp = &user.UpdateUserStatusResponse{
		Success: res.Success,
		Msg:     res.Msg,
	}
	return resp, nil
}
