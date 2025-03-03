package service

import (
	"context"

	category "github.com/PiaoAdmin/gomall/app/hertz_gateway/hertz_gen/hertz_gateway/category"
	"github.com/PiaoAdmin/gomall/app/hertz_gateway/infra/rpc"
	"github.com/PiaoAdmin/gomall/rpc_gen/kitex_gen/product"
	"github.com/cloudwego/hertz/pkg/app"
)

type AddCategoryService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewAddCategoryService(Context context.Context, RequestContext *app.RequestContext) *AddCategoryService {
	return &AddCategoryService{RequestContext: RequestContext, Context: Context}
}

func (h *AddCategoryService) Run(req *category.AddCategoryReq) (resp *int64, err error) {
	//defer func() {
	// hlog.CtxInfof(h.Context, "req = %+v", req)
	// hlog.CtxInfof(h.Context, "resp = %+v", resp)
	//}()
	// todo edit your code
	res, err := rpc.ProductClient.AddCategory(h.Context, &product.AddCategoryReq{
		Category: &product.Category{
			CategoryName: req.Category.CategoryName,
			ParentId:     req.Category.ParentId,
			Status:       req.Category.Status,
		},
	})
	if err != nil {
		return nil, err
	}
	return &res.Id, nil
}
