package service

import (
	"context"
	"testing"
	auth "github.com/PiaoAdmin/gomall/rpc_gen/kitex_gen/auth"
)

func TestRemoveUserRole_Run(t *testing.T) {
	ctx := context.Background()
	s := NewRemoveUserRoleService(ctx)
	// init req and assert value

	req := &auth.RemoveUserRoleRequest{}
	resp, err := s.Run(req)
	t.Logf("err: %v", err)
	t.Logf("resp: %v", resp)

	// todo: edit your unit test

}
