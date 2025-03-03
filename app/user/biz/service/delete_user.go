package service

import (
	"context"

	"github.com/PiaoAdmin/gomall/app/user/biz/dal/mysql"
	"github.com/PiaoAdmin/gomall/app/user/biz/model"
	"github.com/PiaoAdmin/gomall/common/constant"
	user "github.com/PiaoAdmin/gomall/rpc_gen/kitex_gen/user"
)

type DeleteUserService struct {
	ctx context.Context
} // NewDeleteUserService new DeleteUserService
func NewDeleteUserService(ctx context.Context) *DeleteUserService {
	return &DeleteUserService{ctx: ctx}
}

// Run create note info
func (s *DeleteUserService) Run(req *user.DeleteUserRequest) (resp *user.DeleteUserResponse, err error) {
	// Finish your business logic.
	if req == nil {
		return nil, constant.ReqIsNilError("请求为空")
	}
	if req.UserId <= 0 {
		return nil, constant.ParametersError("用户id错误")
	}
	is_user, err := model.GetUserById(mysql.DB, s.ctx, req.UserId)
	if err != nil {
		return
	}
	if err = model.DeleteUserById(mysql.DB, s.ctx, is_user.ID); err != nil {
		return
	}
	return &user.DeleteUserResponse{
		Success: true,
		Msg:     "删除成功",
	}, nil
}
