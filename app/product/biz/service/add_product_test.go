package service

import (
	"context"
	"testing"

	"github.com/PiaoAdmin/gomall/app/product/biz/dal"
	product "github.com/PiaoAdmin/gomall/rpc_gen/kitex_gen/product"
	"github.com/joho/godotenv"
)

func TestAddProduct_Run(t *testing.T) {
	_ = godotenv.Load("../../.env")
	dal.Init()
	ctx := context.Background()
	s := NewAddProductService(ctx)
	// init req and assert value

	req := &product.AddProductReq{
		Product: &product.Product{
			ProdName:        "Pro2",                     // 示例产品名称
			Brief:           "Short description",        // 示例产品简短描述
			MainImage:       "main_image_url",           // 示例主图URL
			Price:           99,                         // 示例价格
			Status:          1,                          // 示例状态（1 可能表示上架）
			Categories:      nil,                        // 类别为空
			Content:         "Full product description", // 示例产品详细描述
			SecondaryImages: "image1_url",               // 示例副图URL列表
			SoldNum:         100,                        // 示例已售数量
			TotalStock:      500,                        // 示例库存总量
			ListingTime:     int64(1618486800),          // 示例上架时间
		},
	}
	resp, err := s.Run(req)
	t.Logf("err: %v", err)
	t.Logf("resp: %v", resp)

	// todo: edit your unit test

}
