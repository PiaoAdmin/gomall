package service

import (
	"context"
	"time"

	"github.com/PiaoAdmin/gomall/app/product/biz/dal/mysql"
	"github.com/PiaoAdmin/gomall/app/product/biz/model"
	"github.com/PiaoAdmin/gomall/common/constant"
	"github.com/PiaoAdmin/gomall/rpc_gen/kitex_gen/product"
)

type AddProductService struct {
	ctx context.Context
} // NewAddProductService new AddProductService
func NewAddProductService(ctx context.Context) *AddProductService {
	return &AddProductService{ctx: ctx}
}

// Run create note info
func (s *AddProductService) Run(req *product.AddProductReq) (resp *product.AddProductResp, err error) {
	if req == nil || req.Product == nil {
		return nil, constant.ParametersError("请求为空")
	}
	newProduct := &model.Product{
		ProdName:        req.Product.ProdName,
		ShopId:          req.Product.ShopId,
		Brief:           req.Product.Brief,
		MainImage:       req.Product.MainImage,
		Price:           float64(req.Product.Price),
		Status:          int(req.Product.Status),
		Categories:      nil,
		Content:         req.Product.Content,
		SecondaryImages: req.Product.SecondaryImages,
		SoldNum:         int(req.Product.SoldNum),
		TotalStock:      int(req.Product.TotalStock),
		ListingTime:     time.Unix(req.Product.ListingTime, 0),
	}
	err = model.CreateProduct(mysql.DB, newProduct)
	if err != nil {
		return nil, err
	}
	return &product.AddProductResp{
		Id: newProduct.ID,
	}, nil
}
