package service

import (
	"context"

	apiProduct "github.com/PiaoAdmin/pmall/app/api/biz/model/api/product"
	"github.com/PiaoAdmin/pmall/app/api/rpc"
	"github.com/PiaoAdmin/pmall/rpc_gen/product"
	"github.com/cloudwego/hertz/pkg/app"
)

type SearchProductsService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewSearchProductsService(ctx context.Context, c *app.RequestContext) *SearchProductsService {
	return &SearchProductsService{
		RequestContext: c,
		Context:        ctx,
	}
}

func (s *SearchProductsService) Run(req *apiProduct.SearchProductsRequest) (resp *apiProduct.SearchProductsResponse, err error) {
	rpcResp, err := rpc.ProductClient.SearchProducts(s.Context, &product.SearchProductsRequest{
		Page:       req.Page,
		PageSize:   req.PageSize,
		Keyword:    req.Keyword,
		CategoryId: req.CategoryId,
		BrandId:    req.BrandId,
		MinPrice:   req.MinPrice,
		MaxPrice:   req.MaxPrice,
		SortType:   req.SortType,
	})
	if err != nil {
		return nil, err
	}

	list := make([]*apiProduct.SearchSkuDTO, 0, len(rpcResp.List))
	for _, item := range rpcResp.List {
		list = append(list, &apiProduct.SearchSkuDTO{
			SkuId:        item.Sku.Id,
			SkuName:      item.Sku.Name,
			SubTitle:     item.Sku.SubTitle,
			MainImage:    item.Sku.MainImage,
			Price:        item.Sku.Price,
			MarketPrice:  item.Sku.MarketPrice,
			Stock:        item.Sku.Stock,
			SkuSpecData:  item.Sku.SkuSpecData,
			SpuId:        item.SpuId,
			SpuName:      item.SpuName,
			CategoryId:   item.CategoryId,
			CategoryName: item.CategoryName,
			BrandId:      item.BrandId,
			BrandName:    item.BrandName,
			SaleCount:    item.SpuSaleCount,
		})
	}

	return &apiProduct.SearchProductsResponse{
		List:  list,
		Total: rpcResp.Total,
	}, nil
}
