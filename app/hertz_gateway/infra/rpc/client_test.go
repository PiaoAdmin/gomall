package rpc

import (
	"context"
	"testing"

	user "github.com/PiaoAdmin/gomall/rpc_gen/kitex_gen/user"
)

func Test_initUserClient(t *testing.T) {
	initUserClient()
	resp, err := UserClient.GetUserInfo(context.Background(), &user.GetUserInfoRequest{
		UserId: 1891709969154183168,
	})
	if err != nil {
		t.Errorf("err: %v", err)
		return
	}
	t.Logf("resp: %v", resp)
}
