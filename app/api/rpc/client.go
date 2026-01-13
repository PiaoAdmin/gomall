package rpc

import (
	"sync"

	"github.com/PiaoAdmin/pmall/app/api/conf"
	"github.com/PiaoAdmin/pmall/rpc_gen/user/userservice"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/transmeta"
	"github.com/cloudwego/kitex/transport"
	consul "github.com/kitex-contrib/registry-consul"
)

var (
	UserClient  userservice.Client
	once        sync.Once
	err         error
	commonSuite client.Option
)

// TODO: 记得修改
type optionSuite struct {
	opts []client.Option
}

func (s *optionSuite) Options() []client.Option {
	return s.opts
}

func Init() {
	once.Do(func() {
		registryAddr := conf.GetConf().Hertz.RegistryAddr
		r, err := consul.NewConsulResolver(registryAddr)
		if err != nil {
			panic(err)
		}
		suite := &optionSuite{
			opts: []client.Option{
				client.WithMetaHandler(transmeta.ClientHTTP2Handler),
				client.WithTransportProtocol(transport.GRPC),
			},
		}
		if conf.GetConf().Env == "test" {
			suite.opts = append(suite.opts, client.WithHostPorts("localhost:8899"))
		} else {
			suite.opts = append(suite.opts, client.WithResolver(r))
		}
		commonSuite = client.WithSuite(suite)
		initUserClient()
	})
}

func initUserClient() {
	UserClient, err = userservice.NewClient("user", commonSuite)
	if err != nil {
		panic(err)
	}
}
