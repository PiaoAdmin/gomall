package service

import (
	"context"
	"testing"

	"github.com/PiaoAdmin/gomall/app/auth/biz/dal"
	auth "github.com/PiaoAdmin/gomall/rpc_gen/kitex_gen/auth"
	"github.com/joho/godotenv"
)

func TestCheckPermission_Run(t *testing.T) {
	_ = godotenv.Load("../../.env")
	dal.Init()
	ctx := context.Background()
	s := NewCheckPermissionService(ctx)
	// init req and assert value

	req := &auth.PermissionCheckRequest{
		Token:        "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxODkxODI3MTM4NDgzNDU4MDQ4LCJ1c2VybmFtZSI6InRlc3QyIiwiaXNzIjoiZ29tYWxsIiwiZXhwIjoxNzQwMjI2NDAyLCJpYXQiOjE3NDAyMjU4MDJ9.t-wBnzot5oLlIGd8VTQE6CRiF6rNjv4pxTkBcwa1GQY",
		ServiceName:  "user_service",
		ResourceName: "create_user",
	}
	resp, err := s.Run(req)
	t.Logf("err: %v", err)
	t.Logf("resp: %v", resp)

	// todo: edit your unit test

}
