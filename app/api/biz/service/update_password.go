package service

import (
	"context"

	auth "github.com/PiaoAdmin/pmall/app/api/biz/model/api/auth"
	"github.com/PiaoAdmin/pmall/app/api/biz/utils"
	"github.com/PiaoAdmin/pmall/app/api/md/jwt"
	"github.com/PiaoAdmin/pmall/app/api/rpc"
	perrors "github.com/PiaoAdmin/pmall/common/errs"
	"github.com/PiaoAdmin/pmall/rpc_gen/user"
	"github.com/cloudwego/hertz/pkg/app"
)

type UpdatePasswordService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewUpdatePasswordService(Context context.Context, RequestContext *app.RequestContext) *UpdatePasswordService {
	return &UpdatePasswordService{RequestContext: RequestContext, Context: Context}
}

func (h *UpdatePasswordService) Run(req *auth.UpdatePasswordReq) (resp *auth.UpdatePasswordResp, err error) {
	claims := jwt.ExtractClaims(h.Context, h.RequestContext)
	userId := uint64(claims[jwt.JwtMiddleware.IdentityKey].(float64))
	if req.NewPassword == "" || req.OldPassword == "" {
		return nil, perrors.New(perrors.ErrParam.Code, "Password cannot be empty")
	}
	_, err = rpc.UserClient.UpdatePassword(h.Context, &user.UpdatePasswordRequest{
		UserId:      userId,
		OldPassword: req.OldPassword,
		NewPassword: req.NewPassword,
	})
	if err != nil {
		return nil, err
	}
	// 将原来的token加入黑名单
	err = utils.AddTokenToBlacklist(h.Context, h.RequestContext)
	if err != nil {
		return nil, err
	}
	return &auth.UpdatePasswordResp{
		Success: true,
	}, nil
}
