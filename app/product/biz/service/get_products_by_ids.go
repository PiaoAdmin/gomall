package service

import (
	"context"

	"github.com/PiaoAdmin/pmall/app/product/biz/dal/mysql"
	"github.com/PiaoAdmin/pmall/app/product/biz/model"
	"github.com/PiaoAdmin/pmall/common/errs"
	product "github.com/PiaoAdmin/pmall/rpc_gen/product"
)

type GetProductsByIdsService struct {
	ctx context.Context
}

func NewGetProductsByIdsService(ctx context.Context) *GetProductsByIdsService {
	return &GetProductsByIdsService{ctx: ctx}
}

func (s *GetProductsByIdsService) Run(req *product.GetProductsByIdsRequest) (*product.GetProductsByIdsResponse, error) {
	if len(req.Ids) == 0 {
		return &product.GetProductsByIdsResponse{
			Products: make(map[uint64]*product.ProductSPU),
		}, nil
	}

	spus, err := model.GetProductsByIds(s.ctx, mysql.DB, req.Ids)
	if err != nil {
		return nil, errs.New(errs.ErrInternal.Code, "get products by ids failed: "+err.Error())
	}

	productsMap := make(map[uint64]*product.ProductSPU)
	for _, spu := range spus {
		productsMap[spu.ID] = &product.ProductSPU{
			Id:            spu.ID,
			BrandId:       spu.BrandID,
			CategoryId:    spu.CategoryID,
			Name:          spu.Name,
			SubTitle:      spu.SubTitle,
			MainImage:     spu.MainImage,
			PublishStatus: int32(spu.PublishStatus),
			VerifyStatus:  int32(spu.VerifyStatus),
			SaleCount:     int32(spu.SaleCount),
			Sort:          int32(spu.Sort),
			ServiceBits:   spu.ServiceBits,
			Version:       int32(spu.Version),
		}
	}

	return &product.GetProductsByIdsResponse{
		Products: productsMap,
	}, nil
}
