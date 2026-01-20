package service

import (
	"context"

	apiProduct "github.com/PiaoAdmin/pmall/app/api/biz/model/api/product"
	"github.com/PiaoAdmin/pmall/app/api/rpc"
	"github.com/PiaoAdmin/pmall/rpc_gen/product"
	"github.com/cloudwego/hertz/pkg/app"
)

type GetHomeProductsService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewGetHomeProductsService(ctx context.Context, c *app.RequestContext) *GetHomeProductsService {
	return &GetHomeProductsService{
		RequestContext: c,
		Context:        ctx,
	}
}

func (s *GetHomeProductsService) Run(req *apiProduct.GetHomeProductsRequest) (resp *apiProduct.GetHomeProductsResponse, err error) {
	rpcResp, err := rpc.ProductClient.ListProducts(s.Context, &product.ListProductsRequest{
		Page:          req.Page,
		PageSize:      req.PageSize,
		CategoryId:    req.CategoryId,
		BrandId:       req.BrandId,
		PublishStatus: -1, // 不筛选状态
		VerifyStatus:  -1,
	})
	if err != nil {
		return nil, err
	}

	// 转换为 HomeSpuDTO
	list := make([]*apiProduct.HomeSpuDTO, 0, len(rpcResp.List))
	for _, spu := range rpcResp.List {
		list = append(list, &apiProduct.HomeSpuDTO{
			SpuId:     spu.Id,
			Name:      spu.Name,
			SubTitle:  spu.SubTitle,
			MainImage: spu.MainImage,
			LowPrice:  spu.LowPrice,
			HighPrice: spu.HighPrice,
			SaleCount: spu.SaleCount,
		})
	}

	return &apiProduct.GetHomeProductsResponse{
		List:  list,
		Total: rpcResp.Total,
	}, nil
}
