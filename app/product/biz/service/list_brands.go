package service

import (
	"context"

	"github.com/PiaoAdmin/pmall/app/product/biz/dal/mysql"
	"github.com/PiaoAdmin/pmall/app/product/biz/model"
	"github.com/PiaoAdmin/pmall/common/errs"
	product "github.com/PiaoAdmin/pmall/rpc_gen/product"
)

type ListBrandsService struct {
	ctx context.Context
}

func NewListBrandsService(ctx context.Context) *ListBrandsService {
	return &ListBrandsService{ctx: ctx}
}

func (s *ListBrandsService) Run(req *product.ListBrandsRequest) (*product.ListBrandsResponse, error) {
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

	brands, total, err := model.ListBrands(s.ctx, mysql.DB, page, pageSize)
	if err != nil {
		return nil, errs.New(errs.ErrInternal.Code, "list brands failed: "+err.Error())
	}

	respBrands := make([]*product.Brand, 0, len(brands))
	for _, b := range brands {
		respBrands = append(respBrands, &product.Brand{
			Id:          b.ID,
			Name:        b.Name,
			FirstLetter: b.FirstLetter,
			Logo:        b.Logo,
			Sort:        int32(b.Sort),
		})
	}

	return &product.ListBrandsResponse{
		Brands: respBrands,
		Total:  total,
	}, nil
}
