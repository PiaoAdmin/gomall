package service

import (
	"context"

	"github.com/PiaoAdmin/pmall/app/product/biz/dal/mysql"
	"github.com/PiaoAdmin/pmall/app/product/biz/model"
	"github.com/PiaoAdmin/pmall/common/errs"
	product "github.com/PiaoAdmin/pmall/rpc_gen/product"
)

type UpdateProductService struct {
	ctx context.Context
}

func NewUpdateProductService(ctx context.Context) *UpdateProductService {
	return &UpdateProductService{ctx: ctx}
}

func (s *UpdateProductService) Run(req *product.UpdateProductRequest) (*product.UpdateProductResponse, error) {
	if req.Spu == nil || req.Spu.Id == 0 {
		return nil, errs.New(errs.ErrParam.Code, "spu id is required")
	}

	spuUpdates := make(map[string]interface{})
	if req.Spu.Name != "" {
		spuUpdates["name"] = req.Spu.Name
	}
	if req.Spu.SubTitle != "" {
		spuUpdates["sub_title"] = req.Spu.SubTitle
	}
	if req.Spu.MainImage != "" {
		spuUpdates["main_image"] = req.Spu.MainImage
	}
	if req.Spu.BrandId != 0 {
		spuUpdates["brand_id"] = req.Spu.BrandId
	}
	if req.Spu.CategoryId != 0 {
		spuUpdates["category_id"] = req.Spu.CategoryId
	}
	if req.Spu.ServiceBits != 0 {
		spuUpdates["service_bits"] = req.Spu.ServiceBits
	}

	if len(spuUpdates) > 0 {
		_, err := model.UpdateSPU(s.ctx, mysql.DB, req.Spu.Id, spuUpdates)
		if err != nil {
			return nil, errs.New(errs.ErrInternal.Code, "update spu failed: "+err.Error())
		}
	}

	if req.Detail != nil {
		detailUpdates := make(map[string]interface{})
		if req.Detail.Description != "" {
			detailUpdates["description"] = req.Detail.Description
		}
		if len(req.Detail.Images) > 0 {
			detailUpdates["images"] = req.Detail.Images
		}
		if len(req.Detail.Videos) > 0 {
			detailUpdates["videos"] = req.Detail.Videos
		}
		if req.Detail.MarketTagJson != "" {
			detailUpdates["market_tag_json"] = req.Detail.MarketTagJson
		}
		if req.Detail.TechTagJson != "" {
			detailUpdates["tech_tag_json"] = req.Detail.TechTagJson
		}
		if req.Detail.FaqJson != "" {
			detailUpdates["faq_json"] = req.Detail.FaqJson
		}

		if len(detailUpdates) > 0 {
			_, err := model.UpdateProductDetail(s.ctx, req.Spu.Id, detailUpdates)
			if err != nil {
				return nil, errs.New(errs.ErrInternal.Code, "update product detail failed: "+err.Error())
			}
		}
	}

	return &product.UpdateProductResponse{
		Success: true,
	}, nil
}
