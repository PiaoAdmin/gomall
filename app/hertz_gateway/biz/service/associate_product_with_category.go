package service

import (
	"context"

	common "github.com/PiaoAdmin/gomall/app/hertz_gateway/hertz_gen/hertz_gateway/common"
	product "github.com/PiaoAdmin/gomall/app/hertz_gateway/hertz_gen/hertz_gateway/product"
	"github.com/cloudwego/hertz/pkg/app"
)

type AssociateProductWithCategoryService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewAssociateProductWithCategoryService(Context context.Context, RequestContext *app.RequestContext) *AssociateProductWithCategoryService {
	return &AssociateProductWithCategoryService{RequestContext: RequestContext, Context: Context}
}

func (h *AssociateProductWithCategoryService) Run(req *product.AssociateProductWithCategoryReq) (resp *common.Empty, err error) {
	//defer func() {
	// hlog.CtxInfof(h.Context, "req = %+v", req)
	// hlog.CtxInfof(h.Context, "resp = %+v", resp)
	//}()
	// todo edit your code
	return
}
