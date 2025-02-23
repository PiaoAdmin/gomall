package service

import (
	"context"
	"errors"

	"github.com/PiaoAdmin/gomall/app/user/biz/dal/model"
	"github.com/PiaoAdmin/gomall/app/user/biz/dal/mysql"
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
	if req == nil || req.UserId <= 0 {
		return nil, errors.New("参数错误")
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
