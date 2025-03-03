package service

import (
	"context"

	product "github.com/PiaoAdmin/gomall/app/hertz_gateway/hertz_gen/hertz_gateway/product"
	"github.com/PiaoAdmin/gomall/app/hertz_gateway/infra/rpc"
	"github.com/PiaoAdmin/gomall/app/hertz_gateway/utils"
	rpcProduct "github.com/PiaoAdmin/gomall/rpc_gen/kitex_gen/product"
	"github.com/jinzhu/copier"

	"github.com/cloudwego/hertz/pkg/app"
)

type CreateProductService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewCreateProductService(Context context.Context, RequestContext *app.RequestContext) *CreateProductService {
	return &CreateProductService{RequestContext: RequestContext, Context: Context}
}

func (h *CreateProductService) Run(req *product.CreateProductReq) (resp *int64, err error) {
	//defer func() {
	// hlog.CtxInfof(h.Context, "req = %+v", req)
	// hlog.CtxInfof(h.Context, "resp = %+v", resp)
	//}()
	// todo edit your code
	if req == nil || req.Product == nil {
		return nil, err
	}
	product := rpcProduct.Product{}
	err = copier.Copy(&product, &req.Product)
	if err != nil {
		return nil, err
	}
	id, err := utils.GetUserIdFromToken(h.Context, h.RequestContext)
	if err != nil {
		return nil, err
	}
	product.ShopId = *id
	res, err := rpc.ProductClient.AddProduct(h.Context, &rpcProduct.AddProductReq{
		Product: &product})
	return &res.Id, err
}
