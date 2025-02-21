package service

import (
	"context"
	"github.com/PiaoAdmin/gomall/app/product/biz/dal/mysql"
	product "github.com/PiaoAdmin/gomall/rpc_gen/kitex_gen/product"
	"github.com/joho/godotenv"
	"testing"
)

func TestGetProduct_Run(t *testing.T) {
	err2 := godotenv.Load("../../.env")
	if err2 != nil {
		t.Log("load .env success")
	}
	mysql.Init()
	ctx := context.Background()
	s := NewGetProductService(ctx)
	//l := NewListProductsService(ctx)
	// init req and assert value

	req := &product.GetProductReq{
		Id: 1892154393032900608,
	}
	resp, err := s.Run(req)

	t.Logf("err: %v", err)
	t.Logf("resp: %v", resp)

	// todo: edit your unit test

}
