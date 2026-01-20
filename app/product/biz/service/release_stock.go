package service

import (
	"context"

	"github.com/PiaoAdmin/pmall/app/product/biz/dal/mysql"
	"github.com/PiaoAdmin/pmall/app/product/biz/model"
	"github.com/PiaoAdmin/pmall/common/errs"
	product "github.com/PiaoAdmin/pmall/rpc_gen/product"
	"gorm.io/gorm"
)

type ReleaseStockService struct {
	ctx context.Context
}

func NewReleaseStockService(ctx context.Context) *ReleaseStockService {
	return &ReleaseStockService{ctx: ctx}
}

// only order service call this interface
func (s *ReleaseStockService) Run(req *product.ReleaseStockRequest) (*product.ReleaseStockResponse, error) {
	if req.OrderSn == "" {
		return nil, errs.New(errs.ErrParam.Code, "order_sn is required")
	}
	if len(req.Items) == 0 {
		return nil, errs.New(errs.ErrParam.Code, "items is empty")
	}

	err := mysql.DB.Transaction(func(tx *gorm.DB) error {
		for _, item := range req.Items {
			if item.SkuId == 0 || item.Count <= 0 {
				return errs.New(errs.ErrParam.Code, "invalid sku_id or count")
			}

			err := model.ReleaseStock(s.ctx, tx, item.SkuId, int(item.Count))
			if err != nil {
				if err == gorm.ErrRecordNotFound {
					return errs.New(errs.ErrInternal.Code, "lock stock not enough or sku not found")
				}
				return err
			}
		}
		return nil
	})

	if err != nil {
		if e, ok := err.(*errs.Error); ok {
			return nil, e
		}
		return nil, errs.New(errs.ErrInternal.Code, "release stock failed: "+err.Error())
	}

	return &product.ReleaseStockResponse{
		Success: true,
	}, nil
}
