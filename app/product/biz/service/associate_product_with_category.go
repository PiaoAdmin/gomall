package service

import (
	"context"
	"github.com/PiaoAdmin/gomall/app/product/biz/dal/mysql"
	"github.com/PiaoAdmin/gomall/app/product/biz/model"
	product "github.com/PiaoAdmin/gomall/rpc_gen/kitex_gen/product"
)

type AssociateProductWithCategoryService struct {
	ctx context.Context
} // NewAssociateProductWithCategoryService new AssociateProductWithCategoryService
func NewAssociateProductWithCategoryService(ctx context.Context) *AssociateProductWithCategoryService {
	return &AssociateProductWithCategoryService{ctx: ctx}
}

// Run create note info
func (s *AssociateProductWithCategoryService) Run(req *product.AssociateProductWithCategoryReq) (resp *product.AssociateProductWithCategoryResp, err error) {
	err = model.AssociateProductWithCategory(mysql.DB, req.ProductId, req.CategoryId)
	return
}
