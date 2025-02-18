package rpc

import (
	"sync"

	"github.com/PiaoAdmin/gomall/app/hertz_gateway/conf"
	"github.com/PiaoAdmin/gomall/app/hertz_gateway/utils"
	"github.com/PiaoAdmin/gomall/rpc_gen/kitex_gen/user/userservice"
	"github.com/cloudwego/kitex/client"
	consul "github.com/kitex-contrib/registry-consul"
)

var (
	UserClient userservice.Client
	once       sync.Once
)

func InitClient() {
	once.Do(func() {
		initUserClient()
	})
}

func initUserClient() {
	r, err := consul.NewConsulResolver(conf.GetConf().Hertz.RegistryAddr)
	utils.MustHandleError(err)
	UserClient, err = userservice.NewClient("user", client.WithResolver(r))
	utils.MustHandleError(err)
}
