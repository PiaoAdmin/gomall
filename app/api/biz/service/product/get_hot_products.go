package service

import (
	"context"
	"fmt"

	apiProduct "github.com/PiaoAdmin/pmall/app/api/biz/model/api/product"
	"github.com/PiaoAdmin/pmall/app/api/rpc"
	"github.com/PiaoAdmin/pmall/rpc_gen/product"
	"github.com/cloudwego/hertz/pkg/app"
)

type GetHotProductsService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewGetHotProductsService(ctx context.Context, c *app.RequestContext) *GetHotProductsService {
	return &GetHotProductsService{
		RequestContext: c,
		Context:        ctx,
	}
}

func (s *GetHotProductsService) Run(req *apiProduct.GetHotProductsRequest) (resp *apiProduct.GetHotProductsResponse, err error) {
	limit := req.Limit
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	rpcResp, err := rpc.ProductClient.GetHotProducts(s.Context, &product.GetHotProductsRequest{
		Limit: limit,
	})
	if err != nil {
		return nil, err
	}

	// 转换为 HotProductDTO
	list := make([]*apiProduct.HotProductDTO, 0, len(rpcResp.Products))
	for _, p := range rpcResp.Products {
		list = append(list, &apiProduct.HotProductDTO{
			SpuId:     p.Id,
			Name:      p.Name,
			SubTitle:  p.SubTitle,
			MainImage: p.MainImage,
			LowPrice:  fmt.Sprintf("%.2f", p.LowPrice),
			SaleCount: p.SaleCount,
		})
	}

	return &apiProduct.GetHotProductsResponse{
		Products: list,
	}, nil
}
