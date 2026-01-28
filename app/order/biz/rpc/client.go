package rpc

import (
	"sync"

	"github.com/PiaoAdmin/pmall/app/order/conf"
	"github.com/PiaoAdmin/pmall/common/clientsuite"
	"github.com/PiaoAdmin/pmall/rpc_gen/product/productservice"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/transmeta"
	"github.com/cloudwego/kitex/transport"
)

var (
	ProductClient productservice.Client
	once          sync.Once
	err           error
	registryAddr  string
	serviceName   string
)

func Init() {
	once.Do(func() {
		if conf.GetConf().Env == "test" {
			initProductClientDirect("127.0.0.1:9900")
			return
		}
		registryAddr = conf.GetConf().Registry.RegistryAddress[0]
		serviceName = conf.GetConf().Kitex.Service
		initProductClient()
	})
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
