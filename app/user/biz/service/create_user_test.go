package service

import (
	"context"
	"testing"

	"github.com/PiaoAdmin/gomall/app/user/biz/dal"
	user "github.com/PiaoAdmin/gomall/rpc_gen/kitex_gen/user"
	"github.com/joho/godotenv"
)

func TestCreateUser_Run(t *testing.T) {
	_ = godotenv.Load("../../.env")
	dal.Init()
	ctx := context.Background()
	s := NewCreateUserService(ctx)
	// init req and assert value

	req := &user.CreateUserRequest{
		User: &user.User{
			BaseUser: &user.BaseUser{
				Username:  "test2",
				Email:     "19912@test.com",
				Phone:     "123425678901",
				Avatar:    "test",
				Nickname:  "test",
				Gender:    1,
				BirthDate: "1991-01-01",
			},
			Password: "123456",
		},
	}
	resp, err := s.Run(req)
	t.Logf("err: %v", err)
	t.Logf("resp: %v", resp)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if resp == nil {
		t.Errorf("unexpected nil response")
	}
	// todo: edit your unit test
	t.Logf("resp: %v", resp)
	t.Log("unit test passed")
}
