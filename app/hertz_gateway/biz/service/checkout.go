/*
 * @Author: Jiaqin,Lu
 * @Date: 2025-02-21 12:41:48
 * @Last Modified by: Jiaqin,Lu
 * @Last Modified time: 2025-02-21 12:54:26
 */
package service

import (
	"context"

	checkout "github.com/PiaoAdmin/gomall/app/hertz_gateway/hertz_gen/hertz_gateway/checkout"
	"github.com/cloudwego/hertz/pkg/app"
)

type CheckoutService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewCheckoutService(Context context.Context, RequestContext *app.RequestContext) *CheckoutService {
	return &CheckoutService{RequestContext: RequestContext, Context: Context}
}

func (h *CheckoutService) Run(req *checkout.CheckoutReq) (resp map[string]any, err error) {
	//defer func() {
	// hlog.CtxInfof(h.Context, "req = %+v", req)
	// hlog.CtxInfof(h.Context, "resp = %+v", resp)
	//}()
	// todo edit your code

	return
}
