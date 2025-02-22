package rpc

import (
	consul "github.com/kitex-contrib/registry-consul"
	"sync"

	"github.com/PiaoAdmin/gomall/app/cart/conf"
	cartutils "github.com/PiaoAdmin/gomall/app/cart/utils"
	//"github.com/PiaoAdmin/gomall/rpc_gen/kitex_gen/product/productcatalogservice" //这里需要商品服务的客户端，后续需要联调
	"github.com/cloudwego/kitex/client"
)

var (
	//ProductClient productcatalogservice.Client
	once sync.Once
)

func InitClient() {
	once.Do(func() {
		initProductClient()
	})
}

func initProductClient() {
	var opts []client.Option
	r, err := consul.NewConsulResolver(conf.GetConf().Registry.RegistryAddress[0])
	cartutils.MustHandleError(err)
	opts = append(opts, client.WithResolver(r))
	//ProductClient, err = productcatalogservice.NewClient("product", opts...)
	cartutils.MustHandleError(err)
}
