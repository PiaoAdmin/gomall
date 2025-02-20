package service

import (
	"context"
	"github.com/PiaoAdmin/gomall/app/product/biz/dal/mysql"
	"github.com/PiaoAdmin/gomall/app/product/biz/model"
	"github.com/PiaoAdmin/gomall/rpc_gen/kitex_gen/product"
	"github.com/cloudwego/kitex/pkg/kerrors"
)

type GetProductService struct {
	ctx context.Context
} // NewGetProductService new GetProductService
func NewGetProductService(ctx context.Context) *GetProductService {
	return &GetProductService{ctx: ctx}
}

// Run create note info
func (s *GetProductService) Run(req *product.GetProductReq) (resp *product.GetProductResp, err error) {
	id := req.Id
	if id == 0 {
		return nil, kerrors.NewBizStatusError(2004001, "Product id is required.")
	}
	productQuery := model.NewProductQuery(s.ctx, mysql.DB)
	byId, err := productQuery.GetById(int64(id))
	if err != nil {
		return nil, err
	}
	resp = &product.GetProductResp{
		Product: &product.Product{
			Id:          uint32(byId.ID),
			Name:        byId.ProdName,
			Description: byId.Brief,
			Picture:     byId.MainImage,
			Price:       float32(byId.Price),
			Categories:  nil,
		},
	}
	//// Finish your business logic.
	//product1 := &model.Product{
	//	ProdName:        "testName",
	//	OriPrice:        114.5,
	//	Price:           114.5,
	//	Status:          1,
	//	ShopId:          1,
	//	TotalStocks:     114,
	//	SoldNum:         0,
	//	Brief:           "testBrief",
	//	Content:         "testContent",
	//	MainImage:       "testMainImage",
	//	SecondaryImages: "testSecondaryImages",
	//	ListingTime:     time.Now(),
	//}
	//err2 := model.CreateProduct(mysql.DB, product1)
	//if err2 != nil {
	//	return nil, err2
	//}
	return resp, nil
}
