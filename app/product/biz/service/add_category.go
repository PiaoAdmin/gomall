package service

import (
	"context"
	"github.com/PiaoAdmin/gomall/app/product/biz/dal/mysql"
	"github.com/PiaoAdmin/gomall/app/product/biz/model"
	product "github.com/PiaoAdmin/gomall/rpc_gen/kitex_gen/product"
)

type AddCategoryService struct {
	ctx context.Context
} // NewAddCategoryService new AddCategoryService
func NewAddCategoryService(ctx context.Context) *AddCategoryService {
	return &AddCategoryService{ctx: ctx}
}

// Run create note info
func (s *AddCategoryService) Run(req *product.AddCategoryReq) (resp *product.AddCategoryResp, err error) {
	category := &model.Category{
		Name:     req.Category.CategoryName,
		ParentId: req.Category.ParentId,
		Status:   int(req.Category.Status),
	}
	err = model.CreateCategory(mysql.DB, category)
	return &product.AddCategoryResp{
		Id: category.ID,
	}, err
}
