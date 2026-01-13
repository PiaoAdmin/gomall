package service

import (
	"context"
	"regexp"

	"github.com/PiaoAdmin/pmall/app/user/biz/dal/mysql"
	"github.com/PiaoAdmin/pmall/app/user/biz/model"
	"github.com/PiaoAdmin/pmall/common/errs"
	user "github.com/PiaoAdmin/pmall/rpc_gen/user"
	"golang.org/x/crypto/bcrypt"
)

type RegisterService struct {
	ctx context.Context
}

func NewRegisterService(ctx context.Context) *RegisterService {
	return &RegisterService{ctx: ctx}
}

func (s *RegisterService) Run(req *user.RegisterRequest) (resp *user.RegisterResponse, err error) {
	if req.Username == "" || req.Password == "" || req.PasswordConfirm == "" || req.Email == "" {
		return nil, errs.New(errs.ErrParam.Code, "missing required fields")
	}
	if req.Password != req.PasswordConfirm {
		return nil, errs.New(errs.ErrParam.Code, "password not match")
	}
	var emailRegex = regexp.MustCompile(
		`^[A-Za-z0-9._%+\-]+@[A-Za-z0-9.\-]+\.[A-Za-z]{2,40}$`,
	)
	if !emailRegex.MatchString(req.Email) {
		return nil, errs.New(errs.ErrParam.Code, "invalid email format")
	}
	// 用户名只能包含数字和字母以及下划线
	var usernameRegex = regexp.MustCompile(
		`^[a-zA-Z0-9_]{3,20}$`,
	)
	if !usernameRegex.MatchString(req.Username) {
		return nil, errs.New(errs.ErrParam.Code, "invalid username format")
	}
	// 检查用户是否存在
	u, _ := model.GetByEmail(s.ctx, mysql.DB, req.Email)
	if u != nil {
		return nil, errs.New(errs.ErrRecordAlreadyEx.Code, "user already exists")
	}
	u, _ = model.GetByUsername(s.ctx, mysql.DB, req.Username)
	if u != nil {
		return nil, errs.New(errs.ErrRecordAlreadyEx.Code, "user already exists")
	}

	// 加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errs.ConvertErr(err)
	}
	newUser := &model.User{
		Username: req.Username,
		Password: string(hashedPassword),
		Email:    req.Email,
		Phone:    "",
		Avatar:   "",
		Status:   model.StatusOK,
	}
	if err := model.Create(s.ctx, mysql.DB, newUser); err != nil {
		return nil, errs.New(errs.ErrInternal.Code, err.Error())
	}

	return &user.RegisterResponse{
		User: &user.User{
			Id:       uint64(newUser.ID),
			Username: newUser.Username,
			Email:    newUser.Email,
			Phone:    newUser.Phone,
			Avatar:   newUser.Avatar,
			Status:   newUser.Status,
		},
	}, nil
}
