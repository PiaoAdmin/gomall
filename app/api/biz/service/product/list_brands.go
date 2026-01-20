package service

import (
	"context"

	apiProduct "github.com/PiaoAdmin/pmall/app/api/biz/model/api/product"
	"github.com/PiaoAdmin/pmall/app/api/rpc"
	"github.com/PiaoAdmin/pmall/rpc_gen/product"
	"github.com/cloudwego/hertz/pkg/app"
)

type ListBrandsService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewListBrandsService(ctx context.Context, c *app.RequestContext) *ListBrandsService {
	return &ListBrandsService{
		RequestContext: c,
		Context:        ctx,
	}
}

func (s *ListBrandsService) Run(req *apiProduct.ListBrandsRequest) (resp *apiProduct.ListBrandsResponse, err error) {
	rpcResp, err := rpc.ProductClient.ListBrands(s.Context, &product.ListBrandsRequest{
		Page:     req.Page,
		PageSize: req.PageSize,
	})
	if err != nil {
		return nil, err
	}

	// 转换为 BrandDTO
	brands := make([]*apiProduct.BrandDTO, 0, len(rpcResp.Brands))
	for _, brand := range rpcResp.Brands {
		brands = append(brands, &apiProduct.BrandDTO{
			Id:          brand.Id,
			Name:        brand.Name,
			FirstLetter: brand.FirstLetter,
			Logo:        brand.Logo,
		})
	}

	return &apiProduct.ListBrandsResponse{
		Brands: brands,
		Total:  rpcResp.Total,
	}, nil
}
