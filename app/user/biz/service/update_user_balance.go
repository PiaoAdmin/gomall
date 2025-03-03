package service

import (
	"context"

	"github.com/PiaoAdmin/gomall/app/user/biz/dal/mysql"
	"github.com/PiaoAdmin/gomall/app/user/biz/model"
	"github.com/PiaoAdmin/gomall/common/constant"
	user "github.com/PiaoAdmin/gomall/rpc_gen/kitex_gen/user"
)

type UpdateUserBalanceService struct {
	ctx context.Context
} // NewUpdateUserBalanceService new UpdateUserBalanceService
func NewUpdateUserBalanceService(ctx context.Context) *UpdateUserBalanceService {
	return &UpdateUserBalanceService{ctx: ctx}
}

// Run create note info
func (s *UpdateUserBalanceService) Run(req *user.UpdateUserBalanceRequest) (resp *user.UpdateUserBalanceResponse, err error) {
	// Finish your business logic.
	if req.UserId <= 0 {
		return nil, constant.ParametersError("用户id错误")
	}
	// 获取用户现有余额
	u, err := model.GetUserById(mysql.DB, s.ctx, req.UserId)
	if err != nil {
		return
	}
	// 更新余额
	balance := u.Balance + req.Balance
	if balance < 0 {
		return &user.UpdateUserBalanceResponse{
			Success: false,
			Msg:     "余额不足",
		}, nil
	}
	if err = model.UpdateUser(mysql.DB, s.ctx, req.UserId, &model.User{Balance: balance}); err != nil {
		return
	}
	return &user.UpdateUserBalanceResponse{
		Success: true,
		Msg:     "更新成功",
	}, nil
}
