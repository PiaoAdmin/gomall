package service

import (
	"context"
	"time"

	"github.com/PiaoAdmin/gomall/app/user/biz/dal/model"
	"github.com/PiaoAdmin/gomall/app/user/biz/dal/mysql"

	user "github.com/PiaoAdmin/gomall/rpc_gen/kitex_gen/user"
	"golang.org/x/crypto/bcrypt"
)

type CreateUserService struct {
	ctx context.Context
} // NewCreateUserService new CreateUserService
func NewCreateUserService(ctx context.Context) *CreateUserService {
	return &CreateUserService{ctx: ctx}
}

// Run create note info
func (s *CreateUserService) Run(req *user.CreateUserRequest) (resp *user.CreateUserResponse, err error) {
	// Finish your business logic.
	// TODO: 判空和密码校验

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.User.Password), bcrypt.DefaultCost)
	if err != nil {
		return
	}
	birthDate, err := time.Parse("2006-01-02", req.User.BaseUser.BirthDate)
	if err != nil {
		return
	}
	newUser := &model.User{
		Username:  req.User.BaseUser.Username,
		Password:  string(hashedPassword),
		Email:     req.User.BaseUser.Email,
		Phone:     req.User.BaseUser.Phone,
		Nickname:  req.User.BaseUser.Nickname,
		Avatar:    req.User.BaseUser.Avatar,
		Gender:    int8(req.User.BaseUser.Gender),
		BirthDate: birthDate,
	}
	if err = model.CreateUser(mysql.DB, s.ctx, newUser); err != nil {
		return
	}
	return &user.CreateUserResponse{
		BaseUser: &user.BaseUser{
			UserId:    newUser.ID,
			Username:  newUser.Username,
			Nickname:  newUser.Nickname,
			Avatar:    newUser.Avatar,
			Phone:     newUser.Phone,
			Email:     newUser.Email,
			Gender:    int32(newUser.Gender),
			BirthDate: newUser.BirthDate.Format("2006-01-02"),
		},
	}, nil
}
