package service

import (
	"context"
	"testing"

	"github.com/PiaoAdmin/gomall/app/product/biz/dal"
	product "github.com/PiaoAdmin/gomall/rpc_gen/kitex_gen/product"
	"github.com/joho/godotenv"
)

func TestAddCategory_Run(t *testing.T) {
	_ = godotenv.Load("../../.env")
	dal.Init()
	ctx := context.Background()
	s := NewAddCategoryService(ctx)
	// init req and assert value

	req := &product.AddCategoryReq{
		Category: &product.Category{
			CategoryName: "Category2", // 示例类别名称
			ParentId:     0,           // 示例父类别ID
			Status:       1,           // 示例状态（1 可能表示上架）
		},
	}
	resp, err := s.Run(req)
	t.Logf("err: %v", err)
	t.Logf("resp: %v", resp)

	// todo: edit your unit test

}
