/*
 * @Author: liaosijie
 * @Date: 2025-02-20 11:31:46
 * @Last Modified by: liaosijie
 * @Last Modified time: 2025-02-20 20:47:46
 */

package order

import (
	"context"

	// "douyin-gomall/gomall/app/hertz_gateway/biz/service"
	// "douyin-gomall/gomall/app/hertz_gateway/biz/utils"
	// common "douyin-gomall/gomall/app/hertz_gateway/hertz_gen/hertz_gateway/common"
	"github.com/PiaoAdmin/gomall/app/hertz_gateway/biz/service"
	"github.com/PiaoAdmin/gomall/app/hertz_gateway/biz/utils"
	common "github.com/PiaoAdmin/gomall/app/hertz_gateway/hertz_gen/hertz_gateway/common"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

// OrderList .
// @router / [GET]
func OrderList(ctx context.Context, c *app.RequestContext) {
	var err error
	var req common.Empty
	err = c.BindAndValidate(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	resp, err := service.NewOrderListService(ctx, c).Run(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	utils.SendSuccessResponse(ctx, c, consts.StatusOK, resp)
}
