package auth

import (
	"context"
	auth "github.com/PiaoAdmin/gomall/rpc_gen/kitex_gen/auth"
	"github.com/cloudwego/kitex/client/callopt"
	"github.com/cloudwego/kitex/pkg/klog"
)

func IssueToken(ctx context.Context, req *auth.IssueTokenRequest, callOptions ...callopt.Option) (resp *auth.IssueTokenResponse, err error) {
	resp, err = defaultClient.IssueToken(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "IssueToken call failed,err =%+v", err)
		return nil, err
	}
	return resp, nil
}

func VerifyToken(ctx context.Context, req *auth.VerifyTokenRequest, callOptions ...callopt.Option) (resp *auth.VerifyTokenResponse, err error) {
	resp, err = defaultClient.VerifyToken(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "VerifyToken call failed,err =%+v", err)
		return nil, err
	}
	return resp, nil
}

func RefreshToken(ctx context.Context, req *auth.RefreshTokenRequest, callOptions ...callopt.Option) (resp *auth.RefreshTokenResponse, err error) {
	resp, err = defaultClient.RefreshToken(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "RefreshToken call failed,err =%+v", err)
		return nil, err
	}
	return resp, nil
}

func Logout(ctx context.Context, req *auth.LogoutRequest, callOptions ...callopt.Option) (resp *auth.LogoutResponse, err error) {
	resp, err = defaultClient.Logout(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "Logout call failed,err =%+v", err)
		return nil, err
	}
	return resp, nil
}

func CheckPermission(ctx context.Context, req *auth.PermissionCheckRequest, callOptions ...callopt.Option) (resp *auth.PermissionCheckResponse, err error) {
	resp, err = defaultClient.CheckPermission(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "CheckPermission call failed,err =%+v", err)
		return nil, err
	}
	return resp, nil
}

func AddUserRole(ctx context.Context, req *auth.AddUserRoleRequest, callOptions ...callopt.Option) (resp *auth.AddUserRoleResponse, err error) {
	resp, err = defaultClient.AddUserRole(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "AddUserRole call failed,err =%+v", err)
		return nil, err
	}
	return resp, nil
}

func RemoveUserRole(ctx context.Context, req *auth.RemoveUserRoleRequest, callOptions ...callopt.Option) (resp *auth.RemoveUserRoleResponse, err error) {
	resp, err = defaultClient.RemoveUserRole(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "RemoveUserRole call failed,err =%+v", err)
		return nil, err
	}
	return resp, nil
}

func GetUserRole(ctx context.Context, req *auth.GetUserRoleRequest, callOptions ...callopt.Option) (resp *auth.GetUserRoleResponse, err error) {
	resp, err = defaultClient.GetUserRole(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "GetUserRole call failed,err =%+v", err)
		return nil, err
	}
	return resp, nil
}
