package service

import (
	"context"

	user "github.com/PiaoAdmin/gomall/app/hertz_gateway/hertz_gen/hertz_gateway/user"
	"github.com/PiaoAdmin/gomall/app/hertz_gateway/infra/rpc"
	rpcuser "github.com/PiaoAdmin/gomall/rpc_gen/kitex_gen/user"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/jinzhu/copier"
)

type CreateUserService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewCreateUserService(Context context.Context, RequestContext *app.RequestContext) *CreateUserService {
	return &CreateUserService{RequestContext: RequestContext, Context: Context}
}

func (h *CreateUserService) Run(req *user.CreateUserRequest) (resp *user.CreateUserResponse, err error) {
	newUser := rpcuser.User{}
	if err := copier.Copy(&newUser, req.User); err != nil {
		return nil, err
	}
	res, err := rpc.UserClient.CreateUser(h.Context, &rpcuser.CreateUserRequest{
		User: &newUser,
	})
	if err != nil {
		return
	}
	resp = &user.CreateUserResponse{
		BaseUser: &user.BaseUser{},
	}
	if err := copier.Copy(resp.BaseUser, res.BaseUser); err != nil {
		return nil, err
	}
	return resp, nil
}
