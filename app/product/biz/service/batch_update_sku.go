package service

import (
	"context"

	"github.com/PiaoAdmin/pmall/app/product/biz/dal/mysql"
	"github.com/PiaoAdmin/pmall/app/product/biz/model"
	"github.com/PiaoAdmin/pmall/app/product/biz/utils"
	"github.com/PiaoAdmin/pmall/common/errs"
	product "github.com/PiaoAdmin/pmall/rpc_gen/product"
	"gorm.io/gorm"
)

type BatchUpdateSkuService struct {
	ctx context.Context
}

func NewBatchUpdateSkuService(ctx context.Context) *BatchUpdateSkuService {
	return &BatchUpdateSkuService{ctx: ctx}
}

func (s *BatchUpdateSkuService) Run(req *product.BatchUpdateSkuRequest) (*product.BatchUpdateSkuResponse, error) {
	if len(req.Skus) == 0 {
		return &product.BatchUpdateSkuResponse{Success: true}, nil
	}

	err := mysql.DB.Transaction(func(tx *gorm.DB) error {
		for _, sku := range req.Skus {
			if sku.Id == 0 {
				return errs.New(errs.ErrParam.Code, "sku id is required")
			}

			updates := make(map[string]interface{})

			if sku.Price != "" {
				price, err := utils.PriceConvert(sku.Price)
				if err != nil {
					return errs.New(errs.ErrParam.Code, "invalid price format")
				}
				updates["price"] = price
			}

			if sku.MarketPrice != "" {
				marketPrice, err := utils.PriceConvert(sku.MarketPrice)
				if err != nil {
					return errs.New(errs.ErrParam.Code, "invalid market_price format")
				}
				updates["market_price"] = marketPrice
			}
			// TODO： 这里库存更新最好不要直接更新，而是通过加减库存的方式更新
			if sku.Stock >= 0 {
				updates["stock"] = sku.Stock
			}

			if sku.Name != "" {
				updates["name"] = sku.Name
			}
			if sku.SubTitle != "" {
				updates["sub_title"] = sku.SubTitle
			}
			if sku.MainImage != "" {
				updates["main_image"] = sku.MainImage
			}
			if sku.SkuSpecData != "" {
				// 验证json格式
				err := utils.ValidateJsonFormat(sku.SkuSpecData)
				if err != nil {
					return errs.New(errs.ErrParam.Code, "invalid sku_spec_data format")
				}
				updates["sku_spec_data"] = sku.SkuSpecData
			}
			// TODO: 最好批量更新
			if len(updates) > 0 {
				err := model.UpdateSKU(s.ctx, tx, sku.Id, updates)
				if err != nil {
					return err
				}
			}
		}
		return nil
	})

	if err != nil {
		if e, ok := err.(*errs.Error); ok {
			return nil, e
		}
		return nil, errs.New(errs.ErrInternal.Code, "batch update sku failed: "+err.Error())
	}

	return &product.BatchUpdateSkuResponse{
		Success: true,
	}, nil
}
