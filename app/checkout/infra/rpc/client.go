/*
 * @Author: Jiaqin,Lu
 * @Date: 2025-02-20 15:38:27
 * @Last Modified by: Jiaqin,Lu
 * @Last Modified time: 2025-02-20 16:57:22
 */
package rpc

import (
	"sync"

	"github.com/cloudwego/kitex/client"

	"github.com/PiaoAdmin/gomall/rpc_gen/kitex_gen/cart/cartservice"
	"github.com/PiaoAdmin/gomall/rpc_gen/kitex_gen/order/orderservice"
	"github.com/PiaoAdmin/gomall/rpc_gen/kitex_gen/payment/paymentservice"
	productservice "github.com/PiaoAdmin/gomall/rpc_gen/kitex_gen/product/productcatalogservice"

	consul "github.com/kitex-contrib/registry-consul"
)

//checkout需要用到一连串的rpc调用，需要用到其他微服务进行校验
//此部分用于获得对应服务的client并对其进行调用
//需要使用CartClient，ProductClient，OrderClient，PaymentClient

var (
	PaymentClient paymentservice.Client
	ProductClient productservice.Client
	OrderClient   orderservice.Client
	CartClient    cartservice.Client
	once          sync.Once
)

func InitClient() {
	once.Do(func() {
		initPaymentClient()
		initCartClient()
		initProductClient()
		initOrderClient()
	})
}

func initProductClient() {
	resolver, err := consul.NewConsulResolver("127.0.0.1:8500")
	if err != nil {
		return
	}
	ProductClient, err = productservice.NewClient("product", client.WithResolver(resolver))
	if err != nil {
		return
	}
}

func initPaymentClient() {
	resolver, err := consul.NewConsulResolver("127.0.0.1:8500")
	if err != nil {
		return
	}
	PaymentClient, err = paymentservice.NewClient("payment", client.WithResolver(resolver))
	if err != nil {
		return
	}
}
func initOrderClient() {
	resolver, err := consul.NewConsulResolver("127.0.0.1:8500")
	if err != nil {
		return
	}
	OrderClient, err = orderservice.NewClient("order", client.WithResolver(resolver))
	if err != nil {
		return
	}
}
func initCartClient() {
	resolver, err := consul.NewConsulResolver("127.0.0.1:8500")
	if err != nil {
		return
	}
	CartClient, err = cartservice.NewClient("cart", client.WithResolver(resolver))
	if err != nil {
		return
	}
}
