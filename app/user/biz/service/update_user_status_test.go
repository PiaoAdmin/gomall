package service

import (
	"context"
	"testing"

	"github.com/PiaoAdmin/gomall/app/user/biz/dal"
	user "github.com/PiaoAdmin/gomall/rpc_gen/kitex_gen/user"
	"github.com/joho/godotenv"
)

func TestUpdateUserStatus_Run(t *testing.T) {
	_ = godotenv.Load("../../.env")
	dal.Init()
	ctx := context.Background()
	s := NewUpdateUserStatusService(ctx)
	// init req and assert value

	req := &user.UpdateUserStatusRequest{
		UserId: 1891341893795581952,
		Status: 2,
	}
	resp, err := s.Run(req)
	t.Logf("err: %v", err)
	t.Logf("resp: %v", resp)

	// todo: edit your unit test

}
