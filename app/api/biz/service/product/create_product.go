package service

import (
	"context"
	"fmt"

	apiProduct "github.com/PiaoAdmin/pmall/app/api/biz/model/api/product"
	"github.com/PiaoAdmin/pmall/app/api/rpc"
	"github.com/PiaoAdmin/pmall/rpc_gen/product"
	"github.com/cloudwego/hertz/pkg/app"
)

type CreateProductService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewCreateProductService(ctx context.Context, c *app.RequestContext) *CreateProductService {
	return &CreateProductService{
		RequestContext: c,
		Context:        ctx,
	}
}

func (s *CreateProductService) Run(req *apiProduct.CreateProductRequest) (resp *apiProduct.CreateProductResponse, err error) {
	spu := &product.ProductSPU{
		BrandId:     req.Spu.BrandId,
		CategoryId:  req.Spu.CategoryId,
		Name:        req.Spu.Name,
		SubTitle:    req.Spu.SubTitle,
		MainImage:   req.Spu.MainImage,
		Sort:        req.Spu.Sort,
		ServiceBits: req.Spu.ServiceBits,
	}

	// 转换 SKUs
	skus := make([]*product.ProductSKU, 0, len(req.Skus))
	for _, apiSku := range req.Skus {
		skus = append(skus, &product.ProductSKU{
			SkuCode:     apiSku.SkuCode,
			Name:        apiSku.Name,
			SubTitle:    apiSku.SubTitle,
			MainImage:   apiSku.MainImage,
			Price:       apiSku.Price,
			MarketPrice: apiSku.MarketPrice,
			Stock:       apiSku.Stock,
			SkuSpecData: apiSku.SkuSpecData,
		})
	}

	// 转换 Detail
	var detail *product.ProductDetail
	if req.Detail != nil {
		detail = &product.ProductDetail{
			Description:   req.Detail.Description,
			Images:        req.Detail.Images,
			Videos:        req.Detail.Videos,
			MarketTagJson: req.Detail.MarketTagJson,
			TechTagJson:   req.Detail.TechTagJson,
		}
	}

	// 调用 RPC CreateProduct
	rpcResp, err := rpc.ProductClient.CreateProduct(s.Context, &product.CreateProductRequest{
		Spu:    spu,
		Skus:   skus,
		Detail: detail,
	})
	if err != nil {
		fmt.Print(err.Error())
		return nil, err
	}

	return &apiProduct.CreateProductResponse{
		SpuId:   rpcResp.SpuId,
		Message: "商品创建成功",
	}, nil
}
