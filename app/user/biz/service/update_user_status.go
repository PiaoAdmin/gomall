package service

import (
	"context"

	"github.com/PiaoAdmin/gomall/app/user/biz/dal/mysql"
	"github.com/PiaoAdmin/gomall/app/user/biz/model"
	"github.com/PiaoAdmin/gomall/common/constant"
	user "github.com/PiaoAdmin/gomall/rpc_gen/kitex_gen/user"
)

type UpdateUserStatusService struct {
	ctx context.Context
} // NewUpdateUserStatusService new UpdateUserStatusService
func NewUpdateUserStatusService(ctx context.Context) *UpdateUserStatusService {
	return &UpdateUserStatusService{ctx: ctx}
}

// Run create note info
func (s *UpdateUserStatusService) Run(req *user.UpdateUserStatusRequest) (resp *user.UpdateUserStatusResponse, err error) {
	// Finish your business logic.
	if req.Status == 0 {
		return nil, constant.ParametersError("无效的状态")
	}
	if req.UserId <= 0 {
		return nil, constant.ParametersError("无效的用户ID")
	}

	// 更新状态
	if err = model.UpdateUser(mysql.DB, s.ctx, req.UserId, &model.User{Status: req.Status}); err != nil {
		return
	}
	return &user.UpdateUserStatusResponse{
		Success: true,
		Msg:     "更新成功",
	}, nil
}
