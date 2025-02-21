package auth

import (
	"context"

	"github.com/PiaoAdmin/gomall/app/hertz_gateway/biz/service"
	"github.com/PiaoAdmin/gomall/app/hertz_gateway/biz/utils"
	auth "github.com/PiaoAdmin/gomall/app/hertz_gateway/hertz_gen/hertz_gateway/auth"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

// Login .
// @router /auth/login [POST]
func Login(ctx context.Context, c *app.RequestContext) {
	var err error
	var req auth.LoginRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	resp := &auth.LoginResponse{}
	resp, err = service.NewLoginService(ctx, c).Run(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	utils.SendSuccessResponse(ctx, c, consts.StatusOK, resp)
}

// VerifyToken .
// @router /auth/verify [POST]
func VerifyToken(ctx context.Context, c *app.RequestContext) {
	var err error
	var req auth.VerifyTokenRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	resp := &auth.VerifyTokenResponse{}
	resp, err = service.NewVerifyTokenService(ctx, c).Run(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	utils.SendSuccessResponse(ctx, c, consts.StatusOK, resp)
}

// RefreshToken .
// @router /auth/refresh [POST]
func RefreshToken(ctx context.Context, c *app.RequestContext) {
	var err error
	var req auth.RefreshTokenRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	resp := &auth.RefreshTokenResponse{}
	resp, err = service.NewRefreshTokenService(ctx, c).Run(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	utils.SendSuccessResponse(ctx, c, consts.StatusOK, resp)
}

// Logout .
// @router /auth/logout [POST]
func Logout(ctx context.Context, c *app.RequestContext) {
	var err error
	var req auth.LogoutRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	resp := &auth.LogoutResponse{}
	resp, err = service.NewLogoutService(ctx, c).Run(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	utils.SendSuccessResponse(ctx, c, consts.StatusOK, resp)
}

// CheckPermission .
// @router /auth/permission/check [POST]
func CheckPermission(ctx context.Context, c *app.RequestContext) {
	var err error
	var req auth.PermissionCheckRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	resp := &auth.PermissionCheckResponse{}
	resp, err = service.NewCheckPermissionService(ctx, c).Run(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	utils.SendSuccessResponse(ctx, c, consts.StatusOK, resp)
}

// AddUserRole .
// @router /auth/role/add [POST]
func AddUserRole(ctx context.Context, c *app.RequestContext) {
	var err error
	var req auth.AddUserRoleRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	resp := &auth.AddUserRoleResponse{}
	resp, err = service.NewAddUserRoleService(ctx, c).Run(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	utils.SendSuccessResponse(ctx, c, consts.StatusOK, resp)
}

// RemoveUserRole .
// @router /auth/role/remove [POST]
func RemoveUserRole(ctx context.Context, c *app.RequestContext) {
	var err error
	var req auth.RemoveUserRoleRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	resp := &auth.RemoveUserRoleResponse{}
	resp, err = service.NewRemoveUserRoleService(ctx, c).Run(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	utils.SendSuccessResponse(ctx, c, consts.StatusOK, resp)
}

// GetUserRole .
// @router /auth/role/get [GET]
func GetUserRole(ctx context.Context, c *app.RequestContext) {
	var err error
	var req auth.GetUserRoleRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	resp := &auth.GetUserRoleResponse{}
	resp, err = service.NewGetUserRoleService(ctx, c).Run(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	utils.SendSuccessResponse(ctx, c, consts.StatusOK, resp)
}
