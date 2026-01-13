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

type UpdateUserService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewUpdateUserService(Context context.Context, RequestContext *app.RequestContext) *UpdateUserService {
	return &UpdateUserService{RequestContext: RequestContext, Context: Context}
}

func (h *UpdateUserService) Run(req *auth.UpdateUserReq) (resp *auth.UpdateUserResp, err error) {
	claims := jwt.ExtractClaims(h.Context, h.RequestContext)
	userId := uint64(claims[jwt.JwtMiddleware.IdentityKey].(float64))
	if userId == 0 {
		return nil, perrors.New(perrors.ErrParam.Code, "not logged in")
	}
	_, err = rpc.UserClient.UpdateUser(h.Context, &user.UpdateUserRequest{
		UserId:   userId,
		Email:    req.Email,
		Phone:    req.Phone,
		Avatar:   req.Avatar,
		Username: req.Username,
	})
	if err != nil {
		return nil, err
	}

	return &auth.UpdateUserResp{
		Success: true,
	}, nil
}
