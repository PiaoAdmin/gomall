package service

import (
	"context"
	"testing"

	"github.com/PiaoAdmin/gomall/app/auth/biz/dal"
	auth "github.com/PiaoAdmin/gomall/rpc_gen/kitex_gen/auth"
	"github.com/joho/godotenv"
)

func TestIssueToken_Run(t *testing.T) {
	_ = godotenv.Load("../../.env")
	dal.Init()
	ctx := context.Background()
	s := NewIssueTokenService(ctx)
	// init req and assert value

	req := &auth.IssueTokenRequest{
		Username: "test2",
		Password: "123456",
	}
	resp, err := s.Run(req)
	t.Logf("err: %v", err)
	t.Logf("resp: %v", resp)

	// todo: edit your unit test

}
