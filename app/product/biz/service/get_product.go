package service

import (
	"context"
	"github.com/PiaoAdmin/gomall/app/product/biz/dal/mysql"
	"github.com/PiaoAdmin/gomall/app/product/biz/model"
	"github.com/PiaoAdmin/gomall/rpc_gen/kitex_gen/product"
	"time"
)

type GetProductService struct {
	ctx context.Context
} // NewGetProductService new GetProductService
func NewGetProductService(ctx context.Context) *GetProductService {
	return &GetProductService{ctx: ctx}
}

// Run create note info
func (s *GetProductService) Run(req *product.GetProductReq) (resp *product.GetProductResp, err error) {
	// Finish your business logic.
	product1 := &model.Product{
		ProdName:        "testName",
		OriPrice:        114.5,
		Price:           114.5,
		Status:          1,
		ShopId:          1,
		TotalStocks:     114,
		SoldNum:         0,
		Brief:           "testBrief",
		Content:         "testContent",
		MainImage:       "testMainImage",
		SecondaryImages: "testSecondaryImages",
		ListingTime:     time.Now(),
	}
	err2 := model.CreateProduct(mysql.DB, product1)
	if err2 != nil {
		return nil, err2
	}
	return
}
