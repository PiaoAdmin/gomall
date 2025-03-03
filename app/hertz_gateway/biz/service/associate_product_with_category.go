package service

import (
	"context"

	product "github.com/PiaoAdmin/gomall/app/hertz_gateway/hertz_gen/hertz_gateway/product"
	"github.com/PiaoAdmin/gomall/app/hertz_gateway/infra/rpc"
	rpcproduct "github.com/PiaoAdmin/gomall/rpc_gen/kitex_gen/product"
	"github.com/cloudwego/hertz/pkg/app"
)

type AssociateProductWithCategoryService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewAssociateProductWithCategoryService(Context context.Context, RequestContext *app.RequestContext) *AssociateProductWithCategoryService {
	return &AssociateProductWithCategoryService{RequestContext: RequestContext, Context: Context}
}

func (h *AssociateProductWithCategoryService) Run(req *product.AssociateProductWithCategoryReq) (resp *bool, err error) {
	//defer func() {
	// hlog.CtxInfof(h.Context, "req = %+v", req)
	// hlog.CtxInfof(h.Context, "resp = %+v", resp)
	//}()
	// todo edit your code
	res, err := rpc.ProductClient.AssociateProductWithCategory(h.Context, &rpcproduct.AssociateProductWithCategoryReq{
		ProductId:  req.ProductId,
		CategoryId: req.CategoryId,
	})
	if err != nil {
		return nil, err
	}
	return &res.Success, nil
}
