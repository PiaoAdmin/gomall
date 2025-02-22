/*
 * @Author: Jiaqin,Lu
 * @Date: 2025-02-20 15:38:27
 * @Last Modified by: Jiaqin,Lu
 * @Last Modified time: 2025-02-20 16:57:22
 */
package rpc

import (
	"sync"

	"github.com/PiaoAdmin/gomall/app/checkout/conf"
	"github.com/cloudwego/biz-demo/gomall/rpc_gen/kitex_gen/cart/cartservice"
	"github.com/cloudwego/biz-demo/gomall/rpc_gen/kitex_gen/order/orderservice"
	productservice "github.com/cloudwego/biz-demo/gomall/rpc_gen/kitex_gen/product/productcatalogservice"

	//TODO: 修改为PiaoAdmin
	// "github.com/PiaoAdmin/gomall/rpc_gen/kitex_gen/cart/cartservice"
	// "github.com/PiaoAdmin/gomall/rpc_gen/kitex_gen/order/orderservice"
	//productservice "github.com/PiaoAdmin/gomall/rpc_gen/kitex_gen/product/productcatalogservice"
	"github.com/PiaoAdmin/gomall/rpc_gen/kitex_gen/payment/paymentservice"

	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/pkg/transmeta"
	"github.com/cloudwego/kitex/transport"
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
	err           error
)

func InitClient() {
	once.Do(func() {
		initPaymentClient()
		initCartClient()
		initProductClient()
		initOrderClient()
	})
}

func initPaymentClient() {
	var opts []client.Option
	r, err := consul.NewConsulResolver(conf.GetConf().Registry.RegistryAddress[0])
	if err != nil {
		panic(err)
	}
	//获得对应的的客户端
	opts = append(opts, client.WithResolver(r))
	opts = append(opts,
		client.WithClientBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: conf.GetConf().Kitex.Service}),
		client.WithTransportProtocol(transport.GRPC),
		client.WithMetaHandler(transmeta.ClientHTTP2Handler),
	)
	PaymentClient, err = paymentservice.NewClient(conf.GetConf().Kitex.Service, opts...)
	if err != nil {
		panic(err)
	}
}

func initCartClient() {
	var opts []client.Option
	r, err := consul.NewConsulResolver(conf.GetConf().Registry.RegistryAddress[0])
	if err != nil {
		panic(err)
	}
	//获得对应的的客户端
	opts = append(opts, client.WithResolver(r))
	opts = append(opts,
		client.WithClientBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: conf.GetConf().Kitex.Service}),
		client.WithTransportProtocol(transport.GRPC),
		client.WithMetaHandler(transmeta.ClientHTTP2Handler),
	)
	CartClient, err = cartservice.NewClient(conf.GetConf().Kitex.Service, opts...)
	if err != nil {
		panic(err)
	}
}

func initProductClient() {
	var opts []client.Option
	r, err := consul.NewConsulResolver(conf.GetConf().Registry.RegistryAddress[0])
	if err != nil {
		panic(err)
	}
	//获得对应的的客户端
	opts = append(opts, client.WithResolver(r))
	opts = append(opts,
		client.WithClientBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: conf.GetConf().Kitex.Service}),
		client.WithTransportProtocol(transport.GRPC),
		client.WithMetaHandler(transmeta.ClientHTTP2Handler),
	)
	ProductClient, err = productservice.NewClient(conf.GetConf().Kitex.Service, opts...)
	if err != nil {
		panic(err)
	}
}

func initOrderClient() {
	var opts []client.Option
	r, err := consul.NewConsulResolver(conf.GetConf().Registry.RegistryAddress[0])
	if err != nil {
		panic(err)
	}
	//获得对应的的客户端
	opts = append(opts, client.WithResolver(r))
	opts = append(opts,
		client.WithClientBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: conf.GetConf().Kitex.Service}),
		client.WithTransportProtocol(transport.GRPC),
		client.WithMetaHandler(transmeta.ClientHTTP2Handler),
	)
	OrderClient, err = orderservice.NewClient(conf.GetConf().Kitex.Service, opts...)
	if err != nil {
		panic(err)
	}
}
