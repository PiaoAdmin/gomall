package service

import (
	"context"
	"regexp"

	"github.com/PiaoAdmin/pmall/app/user/biz/dal/mysql"
	"github.com/PiaoAdmin/pmall/app/user/biz/model"
	"github.com/PiaoAdmin/pmall/common/errs"
	user "github.com/PiaoAdmin/pmall/rpc_gen/user"
)

type UpdateUserService struct {
	ctx context.Context
}

func NewUpdateUserService(ctx context.Context) *UpdateUserService {
	return &UpdateUserService{ctx: ctx}
}

func (s *UpdateUserService) Run(req *user.UpdateUserRequest) (resp *user.UpdateUserResponse, err error) {
	userRow, err := model.GetById(s.ctx, mysql.DB, req.UserId)
	if err != nil {
		return nil, errs.New(errs.ErrRecordNotFound.Code, "user not found")
	}

	if req.Username != "" {
		u, _ := model.GetByEmail(s.ctx, mysql.DB, req.Email)
		if u != nil {
			return nil, errs.New(errs.ErrRecordAlreadyEx.Code, "username already exists")
		}
		userRow.Username = req.Username
	}
	if req.Email != "" {
		var emailRegex = regexp.MustCompile(
			`^[A-Za-z0-9._%+\-]+@[A-Za-z0-9.\-]+\.[A-Za-z]{2,40}$`,
		)
		if !emailRegex.MatchString(req.Email) {
			return nil, errs.New(errs.ErrParam.Code, "invalid email format")
		}
		userRow.Email = req.Email
	}
	if req.Phone != "" {
		userRow.Phone = req.Phone
	}
	if req.Avatar != "" {
		userRow.Avatar = req.Avatar
	}

	if err := model.Update(s.ctx, mysql.DB, userRow); err != nil {
		return nil, errs.ConvertErr(err)
	}

	return &user.UpdateUserResponse{
		Success: true,
		User: &user.User{
			Id:       uint64(userRow.ID),
			Username: userRow.Username,
			Email:    userRow.Email,
			Phone:    userRow.Phone,
			Avatar:   userRow.Avatar,
			Status:   userRow.Status,
		},
	}, nil
}
