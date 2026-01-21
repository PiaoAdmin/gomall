package rpc

import (
	"sync"

	"github.com/PiaoAdmin/pmall/app/api/conf"
	"github.com/PiaoAdmin/pmall/common/clientsuite"
	"github.com/PiaoAdmin/pmall/rpc_gen/cart/cartservice"
	"github.com/PiaoAdmin/pmall/rpc_gen/product/productservice"
	"github.com/PiaoAdmin/pmall/rpc_gen/user/userservice"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/transmeta"
	"github.com/cloudwego/kitex/transport"
)

var (
	UserClient    userservice.Client
	ProductClient productservice.Client
	CartClient    cartservice.Client
	once          sync.Once
	err           error
	commonSuite   client.Option
)

func Init() {
	once.Do(func() {
		if conf.GetConf().Env == "test" {
			initUserClientDirect("127.0.0.1:8899")
			initProductClientDirect("127.0.0.1:9900")
			initCartClientDirect("127.0.0.1:9901")
			return
		}
		registryAddr := conf.GetConf().Hertz.RegistryAddr
		commonSuite = client.WithSuite(clientsuite.CommonGrpcClientSuite{
			CurrentServiceName: "api",
			RegistryAddr:       registryAddr,
		})
		initUserClient()
		initProductClient()
		initCartClient()
	})
}

// 测试时使用
func initUserClientDirect(addr string) {
	UserClient, err = userservice.NewClient("user",
		client.WithHostPorts(addr),
		client.WithMetaHandler(transmeta.ClientHTTP2Handler),
		client.WithTransportProtocol(transport.GRPC))
	if err != nil {
		hlog.Fatal(err)
	}
}

func initProductClientDirect(addr string) {
	ProductClient, err = productservice.NewClient("product",
		client.WithHostPorts(addr),
		client.WithMetaHandler(transmeta.ClientHTTP2Handler),
		client.WithTransportProtocol(transport.GRPC))
	if err != nil {
		hlog.Fatal(err)
	}
}

func initCartClientDirect(addr string) {
	CartClient, err = cartservice.NewClient("cart",
		client.WithHostPorts(addr),
		client.WithMetaHandler(transmeta.ClientHTTP2Handler),
		client.WithTransportProtocol(transport.GRPC))
	if err != nil {
		hlog.Fatal(err)
	}
}

func initUserClient() {
	UserClient, err = userservice.NewClient("user", commonSuite)
	if err != nil {
		hlog.Fatal(err)
	}
}

func initProductClient() {
	ProductClient, err = productservice.NewClient("product", commonSuite)
	if err != nil {
		hlog.Fatal(err)
	}
}

func initCartClient() {
	CartClient, err = cartservice.NewClient("cart", commonSuite)
	if err != nil {
		hlog.Fatal(err)
	}
}
