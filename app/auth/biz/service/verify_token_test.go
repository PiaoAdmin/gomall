package service

import (
	"context"
	"testing"

	"github.com/PiaoAdmin/gomall/app/auth/biz/dal"
	auth "github.com/PiaoAdmin/gomall/rpc_gen/kitex_gen/auth"
	"github.com/joho/godotenv"
)

func TestVerifyToken_Run(t *testing.T) {
	_ = godotenv.Load("../../.env")
	dal.Init()
	ctx := context.Background()
	s := NewVerifyTokenService(ctx)
	// init req and assert value

	req := &auth.VerifyTokenRequest{
		Token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxODkzNTk3NDAxMDMwMDcwMjcyLCJ1c2VybmFtZSI6InRlc3QyIiwiaXNzIjoiZ29tYWxsIiwiZXhwIjoxNzQwNjU3NzQyLCJpYXQiOjE3NDA2NTcxNDJ9.Qx5HRkSL9BjJptqojpRiFC_u9aKwNx2u-t-tKxyoFTA",
	}
	resp, err := s.Run(req)
	t.Logf("err: %v", err)
	t.Logf("resp: %v", resp)

	// todo: edit your unit test

}
