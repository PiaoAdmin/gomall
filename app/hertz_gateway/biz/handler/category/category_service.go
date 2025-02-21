package category

import (
	"context"

	"github.com/PiaoAdmin/gomall/app/hertz_gateway/biz/service"
	"github.com/PiaoAdmin/gomall/app/hertz_gateway/biz/utils"
	category "github.com/PiaoAdmin/gomall/app/hertz_gateway/hertz_gen/hertz_gateway/category"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

// Category .
// @router /category/:category [GET]
func Category(ctx context.Context, c *app.RequestContext) {
	var err error
	var req category.CategoryReq
	err = c.BindAndValidate(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	resp, err := service.NewCategoryService(ctx, c).Run(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	utils.SendSuccessResponse(ctx, c, consts.StatusOK, resp)
}
