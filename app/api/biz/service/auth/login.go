package auth

import (
	"context"

	"github.com/PiaoAdmin/pmall/app/api/biz/model/api/auth"
	"github.com/PiaoAdmin/pmall/app/api/md/jwt"
	"github.com/PiaoAdmin/pmall/app/api/rpc"
	"github.com/PiaoAdmin/pmall/rpc_gen/user"
	"github.com/cloudwego/hertz/pkg/app"
)

type LoginService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewLoginService(ctx context.Context, c *app.RequestContext) *LoginService {
	return &LoginService{
		RequestContext: c,
		Context:        ctx,
	}
}

func (s *LoginService) Run(req *auth.LoginReq) (resp *auth.LoginResp, err error) {
	rpcResp, err := rpc.UserClient.Login(s.Context, &user.LoginRequest{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		return nil, err
	}

	token, expire, err := jwt.JwtMiddleware.TokenGenerator(rpcResp.User.Id)
	if err != nil {
		return nil, err
	}

	return &auth.LoginResp{
		User: &auth.UserInfo{
			Id:       rpcResp.User.Id,
			Username: rpcResp.User.Username,
			Email:    rpcResp.User.Email,
			Phone:    rpcResp.User.Phone,
			Avatar:   rpcResp.User.Avatar,
			Status:   rpcResp.User.Status,
		},
		Token:    token,
		ExpireIn: expire.Unix(),
	}, nil
}
