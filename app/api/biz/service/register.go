package service

import (
	"context"

	auth "github.com/PiaoAdmin/pmall/app/api/biz/model/api/auth"
	"github.com/PiaoAdmin/pmall/app/api/md/jwt"
	"github.com/PiaoAdmin/pmall/app/api/rpc"
	perrors "github.com/PiaoAdmin/pmall/common/errs"
	"github.com/PiaoAdmin/pmall/rpc_gen/user"
	"github.com/cloudwego/hertz/pkg/app"
)

type RegisterService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewRegisterService(Context context.Context, RequestContext *app.RequestContext) *RegisterService {
	return &RegisterService{RequestContext: RequestContext, Context: Context}
}

func (h *RegisterService) Run(req *auth.RegisterReq) (resp *auth.RegisterResp, err error) {
	if req.Password != req.PasswordConfirm {
		return nil, perrors.New(perrors.ErrParam.Code, "Password and password confirmation do not match")
	}

	rpcResp, err := rpc.UserClient.Register(h.Context, &user.RegisterRequest{
		Username:        req.Username,
		Password:        req.Password,
		Email:           req.Email,
		PasswordConfirm: req.PasswordConfirm,
	})
	if err != nil {
		return nil, err
	}

	token, expire, err := jwt.JwtMiddleware.TokenGenerator(rpcResp.User.Id)
	if err != nil {
		return nil, err
	}

	return &auth.RegisterResp{
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
