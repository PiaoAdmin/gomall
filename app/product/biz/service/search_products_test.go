package service

import (
	"context"
	"testing"

	"github.com/PiaoAdmin/gomall/app/product/biz/dal"
	product "github.com/PiaoAdmin/gomall/rpc_gen/kitex_gen/product"
	"github.com/joho/godotenv"
)

func TestSearchProducts_Run(t *testing.T) {
	_ = godotenv.Load("../../.env")
	dal.Init()
	ctx := context.Background()
	s := NewSearchProductsService(ctx)
	// init req and assert value

	req := &product.SearchProductsReq{
		Query: "Product",
	}
	resp, err := s.Run(req)
	t.Logf("err: %v", err)
	t.Logf("resp: %v", resp)

	// todo: edit your unit test

}
