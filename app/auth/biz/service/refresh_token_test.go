package service

import (
	"context"
	"testing"

	"github.com/PiaoAdmin/gomall/app/auth/biz/dal"
	auth "github.com/PiaoAdmin/gomall/rpc_gen/kitex_gen/auth"
	"github.com/joho/godotenv"
)

func TestRefreshToken_Run(t *testing.T) {
	_ = godotenv.Load("../../.env")
	dal.Init()
	ctx := context.Background()
	s := NewRefreshTokenService(ctx)
	// init req and assert value

	req := &auth.RefreshTokenRequest{
		Token:        "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxODkxODI3MTM4NDgzNDU4MDQ4LCJ1c2VybmFtZSI6InRlc3QyIiwiaXNzIjoiZ29tYWxsIiwiZXhwIjoxNzM5OTUzODEwLCJpYXQiOjE3Mzk5NTM3NTB9.0PErixrOQ4P9hmpnXE3SLw5J1eBSDSFwb1FMPOSEFaw",
		RefreshToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJnb21hbGwiLCJleHAiOjE3Mzk5NTM5MzB9.mWb446CHzKRQ-zaS0VaQ6jl6jtSZ3_CvdQr1lofvXdM",
	}
	resp, err := s.Run(req)
	t.Logf("err: %v", err)
	t.Logf("resp: %v", resp)

	// todo: edit your unit test

}
