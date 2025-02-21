package rpc

import (
	"github.com/PiaoAdmin/gomall/rpc_gen/kitex_gen/cart/cartservice"
	consul "github.com/kitex-contrib/registry-consul"
	"sync"

	"github.com/PiaoAdmin/gomall/app/hertz_gateway/conf"
	gateutils "github.com/PiaoAdmin/gomall/app/hertz_gateway/utils"
	"github.com/cloudwego/kitex/client"
)

var (
	CartClient    cartservice.Client
	UserClient    userservice.Client
	ProductClient productcatalogservice.Client
	once          sync.Once
)

func InitClient() {
	once.Do(func() {
		initUserClient()
		initProductClient()
		initCartClient()
	})
}

func initUserClient() {
	r, err := consul.NewConsulResolver(conf.GetConf().Hertz.RegistryAddr)
	gateutils.MustHandleError(err)
	UserClient, err = userservice.NewClient("user", client.WithResolver(r))
	gateutils.MustHandleError(err)
}

func initProductClient() {
	var opts []client.Option
	r, err := consul.NewConsulResolver(conf.GetConf().Hertz.RegistryAddr)
	gateutils.MustHandleError(err)
	opts = append(opts, client.WithResolver(r))
	ProductClient, err = productcatalogservice.NewClient("cart", opts...)
	gateutils.MustHandleError(err)
}

func initCartClient() {
	var opts []client.Option
	r, err := consul.NewConsulResolver(conf.GetConf().Hertz.RegistryAddr)
	gateutils.MustHandleError(err)
	opts = append(opts, client.WithResolver(r))
	CartClient, err = cartservice.NewClient("cart", opts...)
	gateutils.MustHandleError(err)
}
