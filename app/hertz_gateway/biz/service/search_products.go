package service

import (
	"context"
	"github.com/PiaoAdmin/gomall/app/hertz_gateway/infra/rpc"
	rpcProduct "github.com/PiaoAdmin/gomall/rpc_gen/kitex_gen/product"
	"github.com/cloudwego/hertz/pkg/common/utils"

	product "github.com/PiaoAdmin/gomall/app/hertz_gateway/hertz_gen/hertz_gateway/product"
	"github.com/cloudwego/hertz/pkg/app"
)

type SearchProductsService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewSearchProductsService(Context context.Context, RequestContext *app.RequestContext) *SearchProductsService {
	return &SearchProductsService{RequestContext: RequestContext, Context: Context}
}

func (h *SearchProductsService) Run(req *product.SearchProductsReq) (resp map[string]any, err error) {
	r, err := rpc.ProductClient.SearchProducts(h.Context, &rpcProduct.SearchProductsReq{Query: req.Query})
	if err != nil {
		return
	}
	return utils.H{
		"items": r.Results,
		"q":     req.Query,
	}, nil
}
