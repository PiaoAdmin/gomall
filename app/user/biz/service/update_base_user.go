package service

import (
	"context"
	"errors"
	"time"

	"github.com/PiaoAdmin/gomall/app/user/biz/dal/model"
	"github.com/PiaoAdmin/gomall/app/user/biz/dal/mysql"
	user "github.com/PiaoAdmin/gomall/rpc_gen/kitex_gen/user"
)

type UpdateBaseUserService struct {
	ctx context.Context
} // NewUpdateBaseUserService new UpdateBaseUserService
func NewUpdateBaseUserService(ctx context.Context) *UpdateBaseUserService {
	return &UpdateBaseUserService{ctx: ctx}
}

// Run create note info
func (s *UpdateBaseUserService) Run(req *user.UpdateBaseUserRequest) (resp *user.UpdateBaseUserResponse, err error) {
	// Finish your business logic.
	// 基本校验
	if req.BaseUser == nil {
		return nil, errors.New("user信息不能为空")
	}
	if req.BaseUser.UserId <= 0 {
		return nil, errors.New("无效的用户ID")
	}

	// 更新字段
	u := &model.User{}
	if req.BaseUser != nil {
		if req.BaseUser.BirthDate != "" {
			birthDate, err := time.Parse("2006-01-02", req.BaseUser.BirthDate)
			if err != nil {
				return nil, err
			}
			u.BirthDate = birthDate
		}
		if req.BaseUser.Nickname != "" {
			u.Nickname = req.BaseUser.Nickname
		}
		if req.BaseUser.Avatar != "" {
			u.Avatar = req.BaseUser.Avatar
		}
		if req.BaseUser.Gender != 0 {
			u.Gender = int8(req.BaseUser.Gender)
		}
		if req.BaseUser.Username != "" {
			u.Username = req.BaseUser.Username
		}
		if req.BaseUser.Email != "" {
			u.Email = req.BaseUser.Email
		}
		if req.BaseUser.Phone != "" {
			u.Phone = req.BaseUser.Phone
		}
	}

	if err = model.UpdateUser(mysql.DB, s.ctx, req.BaseUser.UserId, u); err != nil {
		return
	}

	return &user.UpdateBaseUserResponse{
		Success: true,
		Msg:     "更新成功",
	}, nil

}
