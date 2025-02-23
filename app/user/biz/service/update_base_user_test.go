package service

import (
	"context"
	"testing"

	"github.com/PiaoAdmin/gomall/app/user/biz/dal"
	user "github.com/PiaoAdmin/gomall/rpc_gen/kitex_gen/user"
	"github.com/joho/godotenv"
)

func TestUpdateBaseUser_Run(t *testing.T) {
	_ = godotenv.Load("../../.env")
	dal.Init()
	ctx := context.Background()
	s := NewUpdateBaseUserService(ctx)
	// init req and assert value

	req := &user.UpdateBaseUserRequest{
		BaseUser: &user.BaseUser{
			UserId:    1890655699399086080,
			Username:  "test2",
			Email:     "1998888@test.com",
			Phone:     "12345666601",
			Avatar:    "test",
			Nickname:  "test",
			Gender:    1,
			BirthDate: "1991-01-01",
		},
	}
	resp, err := s.Run(req)
	t.Logf("err: %v", err)
	t.Logf("resp: %v", resp)

	// todo: edit your unit test

}
