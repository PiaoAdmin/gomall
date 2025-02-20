package service

import (
	"context"
	"github.com/PiaoAdmin/gomall/app/product/biz/dal/mysql"
	"github.com/PiaoAdmin/gomall/app/product/biz/model"
	product "github.com/PiaoAdmin/gomall/rpc_gen/kitex_gen/product"
)

type ListProductsService struct {
	ctx context.Context
} // NewListProductsService new ListProductsService
func NewListProductsService(ctx context.Context) *ListProductsService {
	return &ListProductsService{ctx: ctx}
}

// Run create note info
func (s *ListProductsService) Run(req *product.ListProductsReq) (resp *product.ListProductsResp, err error) {
	// Finish your business logic.
	categoryQuery := model.NewCategoryQuery(s.ctx, mysql.DB)
	products, err := categoryQuery.GetProductsByCategoryName(req.CategoryName)
	resp = &product.ListProductsResp{}
	resp.Products = make([]*product.Product, len(products))
	for i, p := range products {
		resp.Products[i] = &product.Product{
			Id:          uint32(p.ID),
			Name:        p.ProdName,
			Description: p.Brief,
			Picture:     p.MainImage,
			Price:       float32(p.Price),
			Categories:  nil,
		}
	}
	return resp, nil
}
