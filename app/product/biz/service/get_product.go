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
	p, err := productQuery.GetById(int64(id))
	if err != nil {
		return nil, err
	}
	categories := make([]string, len(p.Categories))
	for i, category := range p.Categories {
		categories[i] = category.Name
	}
	resp = &product.GetProductResp{
		Product: &product.Product{
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
		},
	}

	return resp, nil
}
