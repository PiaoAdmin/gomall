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
		// BaseUser: &rpcuser.BaseUser{
		// 	Username:  req.BaseUser.Username,
		// 	Email:     req.BaseUser.Email,
		// 	Phone:     req.BaseUser.Phone,
		// 	Nickname:  req.BaseUser.Nickname,
		// 	Avatar:    req.BaseUser.Avatar,
		// 	Gender:    req.BaseUser.Gender,
		// 	BirthDate: req.BaseUser.BirthDate,
		// },
	})
	if err != nil {
		return
	}
	resp = &user.CreateUserResponse{
		BaseUser: &user.BaseUser{
			// UserId: res.SafeUser.UserId,
			// BaseUser: &user.BaseUser{
			// 	Username:  res.SafeUser.BaseUser.Username,
			// 	Email:     res.SafeUser.BaseUser.Email,
			// 	Phone:     res.SafeUser.BaseUser.Phone,
			// 	Nickname:  res.SafeUser.BaseUser.Nickname,
			// 	Avatar:    res.SafeUser.BaseUser.Avatar,
			// 	Gender:    res.SafeUser.BaseUser.Gender,
			// 	BirthDate: res.SafeUser.BaseUser.BirthDate,
			// },
		},
	}
	if err := copier.Copy(resp.BaseUser, res.BaseUser); err != nil {
		return nil, err
	}
	return resp, nil
}
