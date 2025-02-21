/**
 * @Author: ZhangHaoChen
 * @Date:   2/21/25 PM5:11
 */

package rpc

import (
	"log"
	"sync"

	"github.com/PiaoAdmin/gomall/rpc_gen/kitex_gen/checkout/checkoutservice"
	"github.com/PiaoAdmin/gomall/rpc_gen/kitex_gen/product/productcatalogservice"
	"github.com/cloudwego/kitex/client"
	consul "github.com/kitex-contrib/registry-consul"
)

var (
	ProductClient  productcatalogservice.Client
	CheckoutClient checkoutservice.Client
	once           sync.Once
)

func Init() {
	once.Do(func() {
		initUserClient()
		initCheckoutClient()
	})
}
func initUserClient() {
	resolver, err := consul.NewConsulResolver("127.0.0.1:8500")
	if err != nil {
		log.Fatal(err)
	}
	ProductClient, err = productcatalogservice.NewClient("product", client.WithResolver(resolver))
	if err != nil {
		log.Fatal(err)
	}
}

func initCheckoutClient() {
	resolver, err := consul.NewConsulResolver("127.0.0.1:8500")
	if err != nil {
		log.Fatal(err)
	}
	CheckoutClient, err = checkoutservice.NewClient("checkout", client.WithResolver(resolver))
	if err != nil {
		log.Fatal(err)
	}
}
