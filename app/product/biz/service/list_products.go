package service

import (
	"context"

	"github.com/PiaoAdmin/pmall/app/product/biz/dal/mysql"
	"github.com/PiaoAdmin/pmall/app/product/biz/model"
	"github.com/PiaoAdmin/pmall/common/errs"
	product "github.com/PiaoAdmin/pmall/rpc_gen/product"
)

type ListProductsService struct {
	ctx context.Context
}

func NewListProductsService(ctx context.Context) *ListProductsService {
	return &ListProductsService{ctx: ctx}
}

// TODO: add caching later
func (s *ListProductsService) Run(req *product.ListProductsRequest) (*product.ListProductsResponse, error) {
	page := int(req.Page)
	pageSize := int(req.PageSize)

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100 // 限制最大页面大小
	}

	spus, total, err := model.ListProducts(
		s.ctx,
		mysql.DB,
		page,
		pageSize,
		req.Keyword,
		req.CategoryId,
		req.BrandId,
		req.PublishStatus,
		req.VerifyStatus,
	)
	if err != nil {
		return nil, errs.New(errs.ErrInternal.Code, "list products failed: "+err.Error())
	}

	protoSPUs := make([]*product.ProductSPU, 0, len(spus))
	for _, spu := range spus {
		protoSPUs = append(protoSPUs, &product.ProductSPU{
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
		})
	}

	return &product.ListProductsResponse{
		List:  protoSPUs,
		Total: total,
	}, nil
}
