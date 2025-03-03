package service

import (
	"context"
	"fmt"
	"time"

	constant "github.com/PiaoAdmin/gomall/common/constant"

	"github.com/PiaoAdmin/gomall/app/user/biz/dal/mysql"
	"github.com/PiaoAdmin/gomall/app/user/biz/model"
	"github.com/PiaoAdmin/gomall/app/user/infra/rpc"
	"github.com/PiaoAdmin/gomall/rpc_gen/kitex_gen/auth"

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
	if req.User == nil || req.User.BaseUser == nil {
		return nil, constant.ReqIsNilError("请求为空")
	}
	if len(req.User.Password) < 8 || len(req.User.Password) > 16 {
		return nil, constant.ParametersError("密码长度不符合要求")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.User.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, constant.ParametersError("密码加密错误", err)
	}
	birthDate, err := time.Parse("2006-01-02", req.User.BaseUser.BirthDate)
	if err != nil {
		return nil, constant.ParametersError("生日格式错误", err)
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
	// TODO: 事务，未成功添加角色回滚
	if err = model.CreateUser(mysql.DB, s.ctx, newUser); err != nil {
		return
	}
	_, err = rpc.AuthClient.AddUserRole(s.ctx, &auth.AddUserRoleRequest{
		UserId:   newUser.ID,
		RoleCode: constant.User.RoleCode,
		RoleName: constant.User.RoleName,
	})
	if err != nil {
		fmt.Print(err.Error())
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
