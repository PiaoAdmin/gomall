package rpc

import (
	"sync"

	"github.com/PiaoAdmin/gomall/app/user/conf"
	"github.com/PiaoAdmin/gomall/rpc_gen/kitex_gen/auth/authservice"
	"github.com/cloudwego/kitex/client"
	consul "github.com/kitex-contrib/registry-consul"
)

var (
	AuthClient authservice.Client
	once       sync.Once
)

func InitClient() {
	once.Do(func() {
		initAuthClient()
	})
}

func initAuthClient() {
	resolver, err := consul.NewConsulResolver(conf.GetConf().Registry.RegistryAddress[0])
	if err != nil {
		return
	}
	AuthClient, err = authservice.NewClient("auth", client.WithResolver(resolver))
	if err != nil {
		return
	}
}
