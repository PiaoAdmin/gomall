package rpc

import (
	"sync"

	"github.com/PiaoAdmin/pmall/app/api/conf"
	"github.com/PiaoAdmin/pmall/common/clientsuite"
	"github.com/PiaoAdmin/pmall/rpc_gen/product/productservice"
	"github.com/PiaoAdmin/pmall/rpc_gen/user/userservice"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/kitex/client"
)

var (
	UserClient    userservice.Client
	ProductClient productservice.Client
	once          sync.Once
	err           error
	commonSuite   client.Option
)

func Init() {
	once.Do(func() {
		registryAddr := conf.GetConf().Hertz.RegistryAddr
		commonSuite = client.WithSuite(clientsuite.CommonGrpcClientSuite{
			CurrentServiceName: "api",
			RegistryAddr:       registryAddr,
		})
		initUserClient()
		initProductClient()
	})
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
