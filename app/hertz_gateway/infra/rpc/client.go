package rpc

import (
	"sync"

	"github.com/PiaoAdmin/gomall/app/hertz_gateway/conf"
	"github.com/PiaoAdmin/gomall/app/hertz_gateway/utils"
	"github.com/cloudwego/kitex/client"
	"github.com/PiaoAdmin/gomall/rpc_gen/kitex_gen/cart/cartservice"
	"github.com/PiaoAdmin/gomall/rpc_gen/kitex_gen/checkout/checkoutservice"
	"github.com/PiaoAdmin/gomall/rpc_gen/kitex_gen/product/productcatalogservice"
	"github.com/cloudwego/kitex/client"
	consul "github.com/kitex-contrib/registry-consul"
)

var (
	ProductClient  productcatalogservice.Client
	CheckoutClient checkoutservice.Client
  CartClient    cartservice.Client
	once           sync.Once
)

func Init() {
	once.Do(func() {
    initCartClient()
		initProductClient()
		initCheckoutClient()
	})
}

func initProductClient() {
	resolver, err := consul.NewConsulResolver("127.0.0.1:8500")
	utils.MustHandleError(err)
	ProductClient, err = productcatalogservice.NewClient("product", client.WithResolver(resolver))
	utils.MustHandleError(err)
}

func initCheckoutClient() {
	resolver, err := consul.NewConsulResolver("127.0.0.1:8500")
	utils.MustHandleError(err)
	CheckoutClient, err = checkoutservice.NewClient("checkout", client.WithResolver(resolver))
	utils.MustHandleError(err)
}
  
func initCartClient() {
	var opts []client.Option
	r, err := consul.NewConsulResolver(conf.GetConf().Hertz.RegistryAddr)
	utils.MustHandleError(err)
	opts = append(opts, client.WithResolver(r))
	CartClient, err = cartservice.NewClient("cart", opts...)
	utils.MustHandleError(err)
}
