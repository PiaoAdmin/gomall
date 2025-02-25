package service

import (
	"context"
	"testing"

	"github.com/PiaoAdmin/gomall/app/product/biz/dal"
	product "github.com/PiaoAdmin/gomall/rpc_gen/kitex_gen/product"
	"github.com/joho/godotenv"
)

func TestListProducts_Run(t *testing.T) {
	_ = godotenv.Load("../../.env")
	dal.Init()
	ctx := context.Background()
	s := NewListProductsService(ctx)
	// init req and assert value

	req := &product.ListProductsReq{
		PageSize:     2,
		Page:         1,
		CategoryName: "Category2",
	}
	resp, err := s.Run(req)
	t.Logf("err: %v", err)
	t.Logf("resp: %v", resp)

	// todo: edit your unit test

}
