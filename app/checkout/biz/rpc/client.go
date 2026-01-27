package rpc

import (
	"sync"

	"github.com/PiaoAdmin/pmall/app/checkout/conf"
	"github.com/PiaoAdmin/pmall/common/clientsuite"
	"github.com/PiaoAdmin/pmall/rpc_gen/order/orderservice"
	"github.com/PiaoAdmin/pmall/rpc_gen/payment/paymentservice"
	"github.com/PiaoAdmin/pmall/rpc_gen/product/productservice"
	"github.com/PiaoAdmin/pmall/rpc_gen/user/userservice"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/transmeta"
	"github.com/cloudwego/kitex/transport"
)

var (
	UserClient    userservice.Client
	ProductClient productservice.Client
	OrderClient   orderservice.Client
	PaymentClient paymentservice.Client
	once          sync.Once
	err           error
	registryAddr  string
	serviceName   string
)

func Init() {
	once.Do(func() {
		if conf.GetConf().Env == "test" {
			initUserClientDirect("127.0.0.1:8899")
			initProductClientDirect("127.0.0.1:9900")
			initOrderClientDirect("127.0.0.1:9902")
			initPaymentClientDirect("127.0.0.1:9904")
			return
		}
		registryAddr = conf.GetConf().Registry.RegistryAddress[0]
		serviceName = conf.GetConf().Kitex.Service
		initUserClient()
		initProductClient()
		initOrderClient()
		initPaymentClient()
	})
}

func initUserClientDirect(addr string) {
	UserClient, err = userservice.NewClient("user",
		client.WithHostPorts(addr),
		client.WithMetaHandler(transmeta.ClientHTTP2Handler),
		client.WithTransportProtocol(transport.GRPC))
	if err != nil {
		klog.Fatal(err)
	}
}

func initProductClientDirect(addr string) {
	ProductClient, err = productservice.NewClient("product",
		client.WithHostPorts(addr),
		client.WithMetaHandler(transmeta.ClientHTTP2Handler),
		client.WithTransportProtocol(transport.GRPC))
	if err != nil {
		klog.Fatal(err)
	}
}

func initOrderClientDirect(addr string) {
	OrderClient, err = orderservice.NewClient("order",
		client.WithHostPorts(addr),
		client.WithMetaHandler(transmeta.ClientHTTP2Handler),
		client.WithTransportProtocol(transport.GRPC))
	if err != nil {
		klog.Fatal(err)
	}
}

func initPaymentClientDirect(addr string) {
	PaymentClient, err = paymentservice.NewClient("payment",
		client.WithHostPorts(addr),
		client.WithMetaHandler(transmeta.ClientHTTP2Handler),
		client.WithTransportProtocol(transport.GRPC))
	if err != nil {
		klog.Fatal(err)
	}
}

func initUserClient() {
	opts := []client.Option{
		client.WithSuite(clientsuite.CommonGrpcClientSuite{
			RegistryAddr:       registryAddr,
			CurrentServiceName: serviceName,
		}),
	}
	UserClient, err = userservice.NewClient("user", opts...)
	if err != nil {
		klog.Fatalf(err.Error())
		panic(err)
	}
}

func initProductClient() {
	opts := []client.Option{
		client.WithSuite(clientsuite.CommonGrpcClientSuite{
			RegistryAddr:       registryAddr,
			CurrentServiceName: serviceName,
		}),
	}
	ProductClient, err = productservice.NewClient("product", opts...)
	if err != nil {
		klog.Fatalf(err.Error())
		panic(err)
	}
}

func initOrderClient() {
	opts := []client.Option{
		client.WithSuite(clientsuite.CommonGrpcClientSuite{
			RegistryAddr:       registryAddr,
			CurrentServiceName: serviceName,
		}),
	}
	OrderClient, err = orderservice.NewClient("order", opts...)
	if err != nil {
		klog.Fatalf(err.Error())
		panic(err)
	}
}

func initPaymentClient() {
	opts := []client.Option{
		client.WithSuite(clientsuite.CommonGrpcClientSuite{
			RegistryAddr:       registryAddr,
			CurrentServiceName: serviceName,
		}),
	}
	PaymentClient, err = paymentservice.NewClient("payment", opts...)
	if err != nil {
		klog.Fatalf(err.Error())
		panic(err)
	}
}
