package service

import (
	"context"

	"github.com/PiaoAdmin/pmall/app/user/biz/dal/mysql"
	"github.com/PiaoAdmin/pmall/app/user/biz/model"
	"github.com/PiaoAdmin/pmall/common/errs"
	user "github.com/PiaoAdmin/pmall/rpc_gen/user"
	"golang.org/x/crypto/bcrypt"
)

type UpdatePasswordService struct {
	ctx context.Context
}

func NewUpdatePasswordService(ctx context.Context) *UpdatePasswordService {
	return &UpdatePasswordService{ctx: ctx}
}

func (s *UpdatePasswordService) Run(req *user.UpdatePasswordRequest) (resp *user.UpdatePasswordResponse, err error) {
	if req.NewPassword == "" || req.OldPassword == "" {
		return nil, errs.New(errs.ErrParam.Code, "password cannot be empty")
	}
	if req.UserId == 0 {
		return nil, errs.New(errs.ErrParam.Code, "not logged in")
	}
	userRow, err := model.GetById(s.ctx, mysql.DB, req.UserId)
	if err != nil {
		return nil, errs.New(errs.ErrRecordNotFound.Code, "user not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(userRow.Password), []byte(req.OldPassword)); err != nil {
		return nil, errs.New(errs.ErrParam.Code, "old password incorrect")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, errs.ConvertErr(err)
	}

	userRow.Password = string(hashedPassword)
	if err := model.Update(s.ctx, mysql.DB, userRow); err != nil {
		return nil, errs.ConvertErr(err)
	}

	return &user.UpdatePasswordResponse{Success: true}, nil
}
