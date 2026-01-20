package service

import (
	"context"

	apiProduct "github.com/PiaoAdmin/pmall/app/api/biz/model/api/product"
	"github.com/PiaoAdmin/pmall/app/api/rpc"
	"github.com/PiaoAdmin/pmall/rpc_gen/product"
	"github.com/cloudwego/hertz/pkg/app"
)

type ListCategoriesService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewListCategoriesService(ctx context.Context, c *app.RequestContext) *ListCategoriesService {
	return &ListCategoriesService{
		RequestContext: c,
		Context:        ctx,
	}
}

func (s *ListCategoriesService) Run(req *apiProduct.ListCategoriesRequest) (resp *apiProduct.ListCategoriesResponse, err error) {
	rpcResp, err := rpc.ProductClient.ListCategories(s.Context, &product.ListCategoriesRequest{
		ParentId: req.ParentId,
	})
	if err != nil {
		return nil, err
	}

	// 转换为 CategoryDTO
	categories := make([]*apiProduct.CategoryDTO, 0, len(rpcResp.Categories))
	for _, cat := range rpcResp.Categories {
		categories = append(categories, convertCategory(cat))
	}

	return &apiProduct.ListCategoriesResponse{
		Categories: categories,
	}, nil
}
