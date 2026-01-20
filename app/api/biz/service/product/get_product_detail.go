package service

import (
	"context"

	apiProduct "github.com/PiaoAdmin/pmall/app/api/biz/model/api/product"
	"github.com/PiaoAdmin/pmall/app/api/rpc"
	"github.com/PiaoAdmin/pmall/rpc_gen/product"
	"github.com/cloudwego/hertz/pkg/app"
)

type GetProductDetailService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewGetProductDetailService(ctx context.Context, c *app.RequestContext) *GetProductDetailService {
	return &GetProductDetailService{
		RequestContext: c,
		Context:        ctx,
	}
}

func (s *GetProductDetailService) Run(req *apiProduct.GetProductDetailRequest) (resp *apiProduct.GetProductDetailResponse, err error) {
	rpcResp, err := rpc.ProductClient.GetProductDetail(s.Context, &product.GetProductDetailRequest{
		Id: req.Id,
	})
	if err != nil {
		return nil, err
	}

	// 转换 Category
	var categoryDTO *apiProduct.CategoryDTO
	if rpcResp.Category != nil {
		categoryDTO = convertCategory(rpcResp.Category)
	}

	// 转换 Brand
	var brandDTO *apiProduct.BrandDTO
	if rpcResp.Brand != nil {
		brandDTO = &apiProduct.BrandDTO{
			Id:          rpcResp.Brand.Id,
			Name:        rpcResp.Brand.Name,
			FirstLetter: rpcResp.Brand.FirstLetter,
			Logo:        rpcResp.Brand.Logo,
		}
	}

	// 转换 SKUs
	skus := make([]*apiProduct.SkuDTO, 0, len(rpcResp.Skus))
	for _, sku := range rpcResp.Skus {
		skus = append(skus, &apiProduct.SkuDTO{
			Id:          sku.Id,
			Name:        sku.Name,
			SubTitle:    sku.SubTitle,
			MainImage:   sku.MainImage,
			Price:       sku.Price,
			MarketPrice: sku.MarketPrice,
			Stock:       sku.Stock,
			SkuSpecData: sku.SkuSpecData,
		})
	}

	// 转换 Detail
	var detailDTO *apiProduct.DetailInfoDTO
	if rpcResp.Detail != nil {
		detailDTO = &apiProduct.DetailInfoDTO{
			Description:   rpcResp.Detail.Description,
			Images:        rpcResp.Detail.Images,
			Videos:        rpcResp.Detail.Videos,
			MarketTagJson: rpcResp.Detail.MarketTagJson,
			TechTagJson:   rpcResp.Detail.TechTagJson,
		}
	}

	return &apiProduct.GetProductDetailResponse{
		Product: &apiProduct.ProductDetailDTO{
			Id:          rpcResp.Spu.Id,
			Name:        rpcResp.Spu.Name,
			SubTitle:    rpcResp.Spu.SubTitle,
			MainImage:   rpcResp.Spu.MainImage,
			LowPrice:    rpcResp.Spu.LowPrice,
			HighPrice:   rpcResp.Spu.HighPrice,
			SaleCount:   rpcResp.Spu.SaleCount,
			ServiceBits: rpcResp.Spu.ServiceBits,
			Category:    categoryDTO,
			Brand:       brandDTO,
			Skus:        skus,
			Detail:      detailDTO,
		},
	}, nil
}

func convertCategory(cat *product.Category) *apiProduct.CategoryDTO {
	if cat == nil {
		return nil
	}

	result := &apiProduct.CategoryDTO{
		Id:       cat.Id,
		ParentId: cat.ParentId,
		Name:     cat.Name,
		Level:    cat.Level,
		Icon:     cat.Icon,
	}

	// 递归转换子分类
	if len(cat.Children) > 0 {
		children := make([]*apiProduct.CategoryDTO, 0, len(cat.Children))
		for _, child := range cat.Children {
			children = append(children, convertCategory(child))
		}
		result.Children = children
	}

	return result
}
