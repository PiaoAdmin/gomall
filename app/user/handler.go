package main

import (
	"context"

	"github.com/PiaoAdmin/gomall/app/user/biz/service"
	user "github.com/PiaoAdmin/gomall/rpc_gen/kitex_gen/user"
)

// UserServiceImpl implements the last service interface defined in the IDL.
type UserServiceImpl struct{}

// CreateUser implements the UserServiceImpl interface.
func (s *UserServiceImpl) CreateUser(ctx context.Context, req *user.CreateUserRequest) (resp *user.CreateUserResponse, err error) {
	resp, err = service.NewCreateUserService(ctx).Run(req)

	return resp, err
}

// DeleteUser implements the UserServiceImpl interface.
func (s *UserServiceImpl) DeleteUser(ctx context.Context, req *user.DeleteUserRequest) (resp *user.DeleteUserResponse, err error) {
	resp, err = service.NewDeleteUserService(ctx).Run(req)

	return resp, err
}

// GetUserInfo implements the UserServiceImpl interface.
func (s *UserServiceImpl) GetUserInfo(ctx context.Context, req *user.GetUserInfoRequest) (resp *user.GetUserInfoResponse, err error) {
	resp, err = service.NewGetUserInfoService(ctx).Run(req)

	return resp, err
}

// UpdateBaseUser implements the UserServiceImpl interface.
func (s *UserServiceImpl) UpdateBaseUser(ctx context.Context, req *user.UpdateBaseUserRequest) (resp *user.UpdateBaseUserResponse, err error) {
	resp, err = service.NewUpdateBaseUserService(ctx).Run(req)

	return resp, err
}

// UpdateUserPassword implements the UserServiceImpl interface.
func (s *UserServiceImpl) UpdateUserPassword(ctx context.Context, req *user.UpdateUserPasswordRequest) (resp *user.UpdateUserPasswordResponse, err error) {
	resp, err = service.NewUpdateUserPasswordService(ctx).Run(req)

	return resp, err
}

// UpdateUserBalance implements the UserServiceImpl interface.
func (s *UserServiceImpl) UpdateUserBalance(ctx context.Context, req *user.UpdateUserBalanceRequest) (resp *user.UpdateUserBalanceResponse, err error) {
	resp, err = service.NewUpdateUserBalanceService(ctx).Run(req)

	return resp, err
}

// UpdateUserStatus implements the UserServiceImpl interface.
func (s *UserServiceImpl) UpdateUserStatus(ctx context.Context, req *user.UpdateUserStatusRequest) (resp *user.UpdateUserStatusResponse, err error) {
	resp, err = service.NewUpdateUserStatusService(ctx).Run(req)

	return resp, err
}
