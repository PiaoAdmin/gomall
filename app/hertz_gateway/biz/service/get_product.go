package service

import (
	"context"
	"github.com/PiaoAdmin/gomall/app/hertz_gateway/infra/rpc"
	"github.com/cloudwego/hertz/pkg/common/utils"

	product "github.com/PiaoAdmin/gomall/app/hertz_gateway/hertz_gen/hertz_gateway/product"
	rpcProduct "github.com/PiaoAdmin/gomall/rpc_gen/kitex_gen/product"
	"github.com/cloudwego/hertz/pkg/app"
)

type GetProductService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewGetProductService(Context context.Context, RequestContext *app.RequestContext) *GetProductService {
	return &GetProductService{RequestContext: RequestContext, Context: Context}
}

func (h *GetProductService) Run(req *product.ProductReq) (resp map[string]any, err error) {
	r, err := rpc.ProductClient.GetProduct(h.Context, &rpcProduct.GetProductReq{Id: req.Id})
	if err != nil {
		return nil, err
	}
	return utils.H{
		"item": r.Product,
	}, nil
}
