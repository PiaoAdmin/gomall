package user

import (
	"context"
	user "github.com/PiaoAdmin/gomall/rpc_gen/kitex_gen/user"
	"github.com/cloudwego/kitex/client/callopt"
	"github.com/cloudwego/kitex/pkg/klog"
)

func CreateUser(ctx context.Context, req *user.CreateUserRequest, callOptions ...callopt.Option) (resp *user.CreateUserResponse, err error) {
	resp, err = defaultClient.CreateUser(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "CreateUser call failed,err =%+v", err)
		return nil, err
	}
	return resp, nil
}

func DeleteUser(ctx context.Context, req *user.DeleteUserRequest, callOptions ...callopt.Option) (resp *user.DeleteUserResponse, err error) {
	resp, err = defaultClient.DeleteUser(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "DeleteUser call failed,err =%+v", err)
		return nil, err
	}
	return resp, nil
}

func UpdateBaseUser(ctx context.Context, req *user.UpdateBaseUserRequest, callOptions ...callopt.Option) (resp *user.UpdateBaseUserResponse, err error) {
	resp, err = defaultClient.UpdateBaseUser(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "UpdateBaseUser call failed,err =%+v", err)
		return nil, err
	}
	return resp, nil
}

func UpdateUserPassword(ctx context.Context, req *user.UpdateUserPasswordRequest, callOptions ...callopt.Option) (resp *user.UpdateUserPasswordResponse, err error) {
	resp, err = defaultClient.UpdateUserPassword(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "UpdateUserPassword call failed,err =%+v", err)
		return nil, err
	}
	return resp, nil
}

func UpdateUserBalance(ctx context.Context, req *user.UpdateUserBalanceRequest, callOptions ...callopt.Option) (resp *user.UpdateUserBalanceResponse, err error) {
	resp, err = defaultClient.UpdateUserBalance(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "UpdateUserBalance call failed,err =%+v", err)
		return nil, err
	}
	return resp, nil
}

func UpdateUserStatus(ctx context.Context, req *user.UpdateUserStatusRequest, callOptions ...callopt.Option) (resp *user.UpdateUserStatusResponse, err error) {
	resp, err = defaultClient.UpdateUserStatus(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "UpdateUserStatus call failed,err =%+v", err)
		return nil, err
	}
	return resp, nil
}

func GetUserInfo(ctx context.Context, req *user.GetUserInfoRequest, callOptions ...callopt.Option) (resp *user.GetUserInfoResponse, err error) {
	resp, err = defaultClient.GetUserInfo(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "GetUserInfo call failed,err =%+v", err)
		return nil, err
	}
	return resp, nil
}
