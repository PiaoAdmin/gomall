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
	categoryQuery := model.NewCategoryQuery(s.ctx, mysql.DB)
	products, err := categoryQuery.GetProductsByCategoryName(req.CategoryName)
	resp = &product.ListProductsResp{}
	resp.Products = make([]*product.Product, len(products))
	for i, p := range products {
		categories := make([]string, len(p.Categories))
		for i, category := range p.Categories {
			categories[i] = category.Name
		}
		resp.Products[i] = &product.Product{
			Id:              p.ID,
			ProdName:        p.ProdName,
			Brief:           p.Brief,
			MainImage:       p.MainImage,
			Price:           float32(p.Price),
			Status:          int32(p.Status),
			Categories:      categories,
			Content:         p.Content,
			SecondaryImages: p.SecondaryImages,
			SoldNum:         int32(p.SoldNum),
			TotalStock:      int32(p.TotalStock),
			ListingTime:     p.ListingTime.Unix(),
		}
	}
	return resp, err
}
