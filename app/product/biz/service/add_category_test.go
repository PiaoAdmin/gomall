package service

import (
	"context"
	product "github.com/PiaoAdmin/gomall/rpc_gen/kitex_gen/product"
	"testing"
)

func TestAddCategory_Run(t *testing.T) {
	ctx := context.Background()
	s := NewAddCategoryService(ctx)
	// init req and assert value

	req := &product.AddCategoryReq{}
	resp, err := s.Run(req)
	t.Logf("err: %v", err)
	t.Logf("resp: %v", resp)

	// todo: edit your unit test

}
