package utils

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/kitex/pkg/kerrors"
)

type BaseResp struct {
	Code int
	Data interface{}
	Msg  string
}

// SendErrResponse  pack error response
func SendErrResponse(ctx context.Context, c *app.RequestContext, code int, err error) {
	// todo edit custom code
	//处理业务异常
	bizErr, isBizErr := kerrors.FromBizStatusError(err)
	if isBizErr {
		resp := new(BaseResp)
		resp.Code = int(bizErr.BizStatusCode())
		resp.Data = nil
		resp.Msg = bizErr.BizMessage()
		c.JSON(code, resp)
		return
	}
	//非业务异常
	c.JSON(code, err.Error())
}

// SendSuccessResponse  pack success response
func SendSuccessResponse(ctx context.Context, c *app.RequestContext, code int, data interface{}) {
	// todo edit custom code
	resp := new(BaseResp)
	resp.Code = code
	resp.Data = data
	resp.Msg = "success"
	c.JSON(code, resp)
}

// func WarpResponse(ctx context.Context, c *app.RequestContext, content map[string]any) map[string]any {
// 	userId, err := gateutils.GetUserIdFromToken(ctx, c)
// 	if err != nil {
// 		return nil
// 	}
// 	content["userId"] = *userId

// 	if *userId > 0 {
// 		cartResp, err := rpc.CartClient.GetCart(ctx, &cart.GetCartReq{
// 			UserId: *userId,
// 		})
// 		if err != nil && cartResp != nil {
// 			content["cart_num"] = len(cartResp.Cart.Items)
// 		}
// 	}

// 	// func WarpResponse(ctx context.Context, c *app.RequestContext, content map[string]any) map[string]any {
// 	// 	var cartNum int
// 	// 	userId := frontendutils.GetUserIdFromCtx(ctx)
// 	// 	cartResp, _ := rpc.CartClient.GetCart(ctx, &cart.GetCartReq{UserId: userId})
// 	// 	if cartResp != nil && cartResp.Cart != nil {
// 	// 		cartNum = len(cartResp.Cart.Items)
// 	// 	}
// 	// 	content["user_id"] = ctx.Value(frontendutils.UserIdKey)
// 	// 	content["cart_num"] = cartNum

// 	return content
// }
