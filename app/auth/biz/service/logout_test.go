package service

import (
	"context"
	"testing"
	auth "github.com/PiaoAdmin/gomall/rpc_gen/kitex_gen/auth"
)

func TestLogout_Run(t *testing.T) {
	ctx := context.Background()
	s := NewLogoutService(ctx)
	// init req and assert value

	req := &auth.LogoutRequest{}
	resp, err := s.Run(req)
	t.Logf("err: %v", err)
	t.Logf("resp: %v", resp)

	// todo: edit your unit test

}
