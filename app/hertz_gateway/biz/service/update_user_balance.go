package service

import (
	"context"

	user "github.com/PiaoAdmin/gomall/app/hertz_gateway/hertz_gen/hertz_gateway/user"
	"github.com/PiaoAdmin/gomall/app/hertz_gateway/infra/rpc"
	rpcuser "github.com/PiaoAdmin/gomall/rpc_gen/kitex_gen/user"
	"github.com/cloudwego/hertz/pkg/app"
)

type UpdateUserBalanceService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewUpdateUserBalanceService(Context context.Context, RequestContext *app.RequestContext) *UpdateUserBalanceService {
	return &UpdateUserBalanceService{RequestContext: RequestContext, Context: Context}
}

func (h *UpdateUserBalanceService) Run(req *user.UpdateUserBalanceRequest) (resp *user.UpdateUserBalanceResponse, err error) {
	// 调用后端 RPC 更新用户余额
	res, err := rpc.UserClient.UpdateUserBalance(h.Context, &rpcuser.UpdateUserBalanceRequest{
		UserId:  req.UserId,
		Balance: req.Balance,
		Type:    req.Type,
	})
	if err != nil {
		return nil, err
	}
	// 构造响应
	resp = &user.UpdateUserBalanceResponse{
		Success: res.Success,
		Msg:     res.Msg,
	}
	return resp, nil
}
