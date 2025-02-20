package main

import (
	"net"
	"time"

	"github.com/PiaoAdmin/gomall/app/auth/biz/dal"
	"github.com/PiaoAdmin/gomall/app/auth/conf"
	"github.com/PiaoAdmin/gomall/rpc_gen/kitex_gen/auth/authservice"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	consulapi "github.com/hashicorp/consul/api"
	"github.com/joho/godotenv"
	kitexlogrus "github.com/kitex-contrib/obs-opentelemetry/logging/logrus"
	consul "github.com/kitex-contrib/registry-consul"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func main() {
	_ = godotenv.Load()
	dal.Init()
	opts := kitexInit()

	svr := authservice.NewServer(new(AuthServiceImpl), opts...)

	err := svr.Run()
	if err != nil {
		klog.Error(err.Error())
	}
}

func kitexInit() (opts []server.Option) {
	// address
	addr, err := net.ResolveTCPAddr("tcp", conf.GetConf().Kitex.Address)
	if err != nil {
		panic(err)
	}
	opts = append(opts, server.WithServiceAddr(addr))

	// service info
	opts = append(opts, server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{
		ServiceName: conf.GetConf().Kitex.Service,
	}))

	// consul registry
	r, err := consul.NewConsulRegister(conf.GetConf().Registry.RegistryAddress[0],
		consul.WithCheck(&consulapi.AgentServiceCheck{
			Interval:                       "7s",
			Timeout:                        "5s",
			DeregisterCriticalServiceAfter: "1m",
		}))
	if err != nil {
		klog.Fatal(err)
	}
	opts = append(opts, server.WithRegistry(r))

	// klog
	logger := kitexlogrus.NewLogger()
	klog.SetLogger(logger)
	klog.SetLevel(conf.LogLevel())
	asyncWriter := &zapcore.BufferedWriteSyncer{
		WS: zapcore.AddSync(&lumberjack.Logger{
			Filename:   conf.GetConf().Kitex.LogFileName,
			MaxSize:    conf.GetConf().Kitex.LogMaxSize,
			MaxBackups: conf.GetConf().Kitex.LogMaxBackups,
			MaxAge:     conf.GetConf().Kitex.LogMaxAge,
		}),
		FlushInterval: time.Minute,
	}
	klog.SetOutput(asyncWriter)
	server.RegisterShutdownHook(func() {
		asyncWriter.Sync()
	})
	return
}
