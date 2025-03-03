package utils

import (
	"context"

	"github.com/PiaoAdmin/gomall/common/constant"
	"github.com/cloudwego/hertz/pkg/app"
)

func GetToken(ctx context.Context, c *app.RequestContext) (token string, err error) {
	token = string(c.GetHeader("Token"))
	if token == "" {
		return "", constant.NotLoginError("Token is empty")
	}
	return
}

func GetRefreshToken(ctx context.Context, c *app.RequestContext) (refreshToken string, err error) {
	refreshToken = string(c.GetHeader("RefreshToken"))
	if refreshToken == "" {
		return "", constant.NotLoginError("RefreshToken is empty")
	}
	return
}
