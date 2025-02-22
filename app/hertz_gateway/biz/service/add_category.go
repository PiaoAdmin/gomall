package service

import (
	"context"

	category "github.com/PiaoAdmin/gomall/app/hertz_gateway/hertz_gen/hertz_gateway/category"
	common "github.com/PiaoAdmin/gomall/app/hertz_gateway/hertz_gen/hertz_gateway/common"
	"github.com/cloudwego/hertz/pkg/app"
)

type AddCategoryService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewAddCategoryService(Context context.Context, RequestContext *app.RequestContext) *AddCategoryService {
	return &AddCategoryService{RequestContext: RequestContext, Context: Context}
}

func (h *AddCategoryService) Run(req *category.AddCategoryReq) (resp *common.Empty, err error) {
	//defer func() {
	// hlog.CtxInfof(h.Context, "req = %+v", req)
	// hlog.CtxInfof(h.Context, "resp = %+v", resp)
	//}()
	// todo edit your code
	return
}
