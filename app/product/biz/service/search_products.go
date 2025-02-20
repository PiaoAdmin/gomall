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
		results = append(results, &product.Product{
			Id:          uint32(p.ID),
			Name:        p.ProdName,
			Description: p.Brief,
			Picture:     p.MainImage,
			Price:       float32(p.Price),
			Categories:  nil,
		})
	}
	return &product.SearchProductsResp{
		Results: results,
	}, nil
}
