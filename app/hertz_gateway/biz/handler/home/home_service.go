package home

import (
	"context"

	"github.com/PiaoAdmin/gomall/app/hertz_gateway/biz/service"
	"github.com/PiaoAdmin/gomall/app/hertz_gateway/biz/utils"
	home "github.com/PiaoAdmin/gomall/app/hertz_gateway/hertz_gen/hertz_gateway/home"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

// Home .
// @router / [GET]
func Home(ctx context.Context, c *app.RequestContext) {
	var err error
	var req home.Empty
	err = c.BindAndValidate(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	resp := &home.Empty{}
	resp, err = service.NewHomeService(ctx, c).Run(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	// utils.SendSuccessResponse(ctx, c, consts.StatusOK, resp)
	c.HTML(consts.StatusOK, "home", resp)
}
