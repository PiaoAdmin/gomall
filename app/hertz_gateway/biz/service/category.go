package service

import (
	"context"
	"github.com/PiaoAdmin/gomall/app/hertz_gateway/infra/rpc"
	rpcProduct "github.com/PiaoAdmin/gomall/rpc_gen/kitex_gen/product"
	"github.com/cloudwego/hertz/pkg/common/utils"

	category "github.com/PiaoAdmin/gomall/app/hertz_gateway/hertz_gen/hertz_gateway/category"
	"github.com/cloudwego/hertz/pkg/app"
)

type CategoryService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewCategoryService(Context context.Context, RequestContext *app.RequestContext) *CategoryService {
	return &CategoryService{RequestContext: RequestContext, Context: Context}
}

func (h *CategoryService) Run(req *category.CategoryReq) (resp map[string]any, err error) {

	r, err := rpc.ProductClient.ListProducts(h.Context, &rpcProduct.ListProductsReq{CategoryName: req.Category})
	if err != nil {
		return nil, err
	}

	return utils.H{
		"CategoryName": req.Category,
		"items":        r.Products,
	}, nil
}
