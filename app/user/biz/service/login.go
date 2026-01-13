package service

import (
	"context"

	"github.com/PiaoAdmin/pmall/app/user/biz/dal/mysql"
	"github.com/PiaoAdmin/pmall/app/user/biz/model"
	"github.com/PiaoAdmin/pmall/common/errs"
	user "github.com/PiaoAdmin/pmall/rpc_gen/user"
	"golang.org/x/crypto/bcrypt"
)

type LoginService struct {
	ctx context.Context
}

func NewLoginService(ctx context.Context) *LoginService {
	return &LoginService{ctx: ctx}
}

func (s *LoginService) Run(req *user.LoginRequest) (resp *user.LoginResponse, err error) {
	if req.Username == "" || req.Password == "" {
		return nil, errs.New(errs.ErrParam.Code, "username or password is empty")
	}
	userRow, err := model.GetByUsername(s.ctx, mysql.DB, req.Username)
	if err != nil {
		return nil, errs.New(errs.ErrRecordNotFound.Code, err.Error())
	}

	if err := bcrypt.CompareHashAndPassword([]byte(userRow.Password), []byte(req.Password)); err != nil {
		return nil, errs.New(errs.ErrParam.Code, "password incorrect")
	}

	return &user.LoginResponse{
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
