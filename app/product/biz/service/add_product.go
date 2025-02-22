package service

import (
	"context"
	"github.com/PiaoAdmin/gomall/app/product/biz/dal/mysql"
	"github.com/PiaoAdmin/gomall/app/product/biz/model"
	"github.com/PiaoAdmin/gomall/rpc_gen/kitex_gen/product"
	"time"
)

type AddProductService struct {
	ctx context.Context
} // NewAddProductService new AddProductService
func NewAddProductService(ctx context.Context) *AddProductService {
	return &AddProductService{ctx: ctx}
}

// Run create note info
func (s *AddProductService) Run(p *product.AddProductReq) (resp *product.AddProductResp, err error) {
	newProduct := &model.Product{
		ProdName:        p.Product.ProdName,
		Brief:           p.Product.Brief,
		MainImage:       p.Product.MainImage,
		Price:           float64(p.Product.Price),
		Status:          int(p.Product.Status),
		Categories:      nil,
		Content:         p.Product.Content,
		SecondaryImages: p.Product.SecondaryImages,
		SoldNum:         int(p.Product.SoldNum),
		TotalStock:      int(p.Product.TotalStock),
		ListingTime:     time.Unix(p.Product.ListingTime, 0),
	}
	err = model.CreateProduct(mysql.DB, newProduct)

	return &product.AddProductResp{
		Id: newProduct.ID,
	}, nil
}
