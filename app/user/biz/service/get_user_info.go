package service

import (
	"context"

	"github.com/PiaoAdmin/pmall/app/user/biz/dal/mysql"
	"github.com/PiaoAdmin/pmall/app/user/biz/model"
	"github.com/PiaoAdmin/pmall/common/errs"
	user "github.com/PiaoAdmin/pmall/rpc_gen/user"
)

type GetUserInfoService struct {
	ctx context.Context
}

func NewGetUserInfoService(ctx context.Context) *GetUserInfoService {
	return &GetUserInfoService{ctx: ctx}
}

func (s *GetUserInfoService) Run(req *user.GetUserInfoRequest) (resp *user.GetUserInfoResponse, err error) {
	if req.UserId == 0 {
		return nil, errs.New(errs.ErrParam.Code, "user_id is empty")
	}
	userRow, err := model.GetById(s.ctx, mysql.DB, req.UserId)
	if err != nil {
		return nil, errs.New(errs.ErrRecordNotFound.Code, err.Error())
	}

	return &user.GetUserInfoResponse{
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
