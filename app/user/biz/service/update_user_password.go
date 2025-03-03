package service

import (
	"context"

	"github.com/PiaoAdmin/gomall/app/user/biz/dal/mysql"
	"github.com/PiaoAdmin/gomall/app/user/biz/model"
	"github.com/PiaoAdmin/gomall/common/constant"
	user "github.com/PiaoAdmin/gomall/rpc_gen/kitex_gen/user"
	"golang.org/x/crypto/bcrypt"
)

type UpdateUserPasswordService struct {
	ctx context.Context
} // NewUpdateUserPasswordService new UpdateUserPasswordService
func NewUpdateUserPasswordService(ctx context.Context) *UpdateUserPasswordService {
	return &UpdateUserPasswordService{ctx: ctx}
}

// Run create note info
func (s *UpdateUserPasswordService) Run(req *user.UpdateUserPasswordRequest) (resp *user.UpdateUserPasswordResponse, err error) {
	// Finish your business logic.
	// 基本校验
	if req.NewPassword == "" {
		return nil, constant.ParametersError("请求参数错误")
	}
	if req.UserId <= 0 {
		return nil, constant.ParametersError("无效的用户ID")
	}
	// 验证旧密码
	u, err := model.GetUserById(mysql.DB, s.ctx, req.UserId)
	if err != nil {
		return
	}
	if err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(req.OldPassword)); err != nil {
		return &user.UpdateUserPasswordResponse{
			Success: false,
			Msg:     "旧密码错误",
		}, nil
	}
	// 更新密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, constant.ParametersError("密码加密错误", err)
	}
	if err = model.UpdateUser(mysql.DB, s.ctx, req.UserId, &model.User{Password: string(hashedPassword)}); err != nil {
		return
	}
	return &user.UpdateUserPasswordResponse{
		Success: true,
		Msg:     "更新成功",
	}, nil
}
