package service

import (
	"context"
	"github.com/PiaoAdmin/gomall/app/product/biz/dal/mysql"
	"github.com/PiaoAdmin/gomall/app/product/biz/model"
	product "github.com/PiaoAdmin/gomall/rpc_gen/kitex_gen/product"
)

type SearchProductsService struct {
	ctx context.Context
} // NewSearchProductsService new SearchProductsService
func NewSearchProductsService(ctx context.Context) *SearchProductsService {
	return &SearchProductsService{ctx: ctx}
}

// Run create note info
func (s *SearchProductsService) Run(req *product.SearchProductsReq) (resp *product.SearchProductsResp, err error) {
	// Finish your business logic.
	productQuery := model.NewProductQuery(s.ctx, mysql.DB)
	products, err := productQuery.SearchProducts(req.Query)
	var results []*product.Product
	for _, p := range products {
		// 把category结构体切片转成string切片
		categories := make([]string, len(p.Categories))
		for i, category := range p.Categories {
			categories[i] = category.Name
		}
		results = append(results, &product.Product{
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
		})
	}
	return &product.SearchProductsResp{
		Results: results,
	}, nil
}
