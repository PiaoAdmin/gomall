package service

import (
	"context"
	"errors"

	"github.com/PiaoAdmin/gomall/app/user/biz/dal/model"
	"github.com/PiaoAdmin/gomall/app/user/biz/dal/mysql"
	user "github.com/PiaoAdmin/gomall/rpc_gen/kitex_gen/user"
)

type GetUserInfoService struct {
	ctx context.Context
} // NewGetUserInfoService new GetUserInfoService
func NewGetUserInfoService(ctx context.Context) *GetUserInfoService {
	return &GetUserInfoService{ctx: ctx}
}

// Run create note info
func (s *GetUserInfoService) Run(req *user.GetUserInfoRequest) (resp *user.GetUserInfoResponse, err error) {
	// Finish your business logic.
	if req == nil || req.UserId <= 0 {
		return nil, errors.New("参数错误")
	}
	userInfo, err := model.GetUserById(mysql.DB, s.ctx, req.UserId)
	if err != nil {
		return
	}
	return &user.GetUserInfoResponse{
		BaseUser: &user.BaseUser{
			UserId:    userInfo.ID,
			Username:  userInfo.Username,
			Email:     userInfo.Email,
			Phone:     userInfo.Phone,
			Nickname:  userInfo.Nickname,
			Avatar:    userInfo.Avatar,
			Gender:    int32(userInfo.Gender),
			BirthDate: userInfo.BirthDate.Format("2006-01-02"),
			Balance:   userInfo.Balance,
			Status:    int32(userInfo.Status),
		},
	}, nil
}
