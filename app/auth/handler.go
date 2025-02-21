package main

import (
	"context"

	"github.com/PiaoAdmin/gomall/app/auth/biz/service"
	auth "github.com/PiaoAdmin/gomall/rpc_gen/kitex_gen/auth"
)

// AuthServiceImpl implements the last service interface defined in the IDL.
type AuthServiceImpl struct{}

// IssueToken implements the AuthServiceImpl interface.
func (s *AuthServiceImpl) IssueToken(ctx context.Context, req *auth.IssueTokenRequest) (resp *auth.IssueTokenResponse, err error) {
	resp, err = service.NewIssueTokenService(ctx).Run(req)

	return resp, err
}

// VerifyToken implements the AuthServiceImpl interface.
func (s *AuthServiceImpl) VerifyToken(ctx context.Context, req *auth.VerifyTokenRequest) (resp *auth.VerifyTokenResponse, err error) {
	resp, err = service.NewVerifyTokenService(ctx).Run(req)

	return resp, err
}

// RefreshToken implements the AuthServiceImpl interface.
func (s *AuthServiceImpl) RefreshToken(ctx context.Context, req *auth.RefreshTokenRequest) (resp *auth.RefreshTokenResponse, err error) {
	resp, err = service.NewRefreshTokenService(ctx).Run(req)

	return resp, err
}

// Logout implements the AuthServiceImpl interface.
func (s *AuthServiceImpl) Logout(ctx context.Context, req *auth.LogoutRequest) (resp *auth.LogoutResponse, err error) {
	resp, err = service.NewLogoutService(ctx).Run(req)

	return resp, err
}

// CheckPermission implements the AuthServiceImpl interface.
func (s *AuthServiceImpl) CheckPermission(ctx context.Context, req *auth.PermissionCheckRequest) (resp *auth.PermissionCheckResponse, err error) {
	resp, err = service.NewCheckPermissionService(ctx).Run(req)

	return resp, err
}

// AddUserRole implements the AuthServiceImpl interface.
func (s *AuthServiceImpl) AddUserRole(ctx context.Context, req *auth.AddUserRoleRequest) (resp *auth.AddUserRoleResponse, err error) {
	resp, err = service.NewAddUserRoleService(ctx).Run(req)

	return resp, err
}

// RemoveUserRole implements the AuthServiceImpl interface.
func (s *AuthServiceImpl) RemoveUserRole(ctx context.Context, req *auth.RemoveUserRoleRequest) (resp *auth.RemoveUserRoleResponse, err error) {
	resp, err = service.NewRemoveUserRoleService(ctx).Run(req)

	return resp, err
}

// GetUserRole implements the AuthServiceImpl interface.
func (s *AuthServiceImpl) GetUserRole(ctx context.Context, req *auth.GetUserRoleRequest) (resp *auth.GetUserRoleResponse, err error) {
	resp, err = service.NewGetUserRoleService(ctx).Run(req)

	return resp, err
}
