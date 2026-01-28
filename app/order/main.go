package main

import (
	"io"
	"net"
	"os"

	"github.com/PiaoAdmin/pmall/app/order/biz/dal"
	"github.com/PiaoAdmin/pmall/app/order/biz/rpc"
	"github.com/PiaoAdmin/pmall/app/order/conf"
	order "github.com/PiaoAdmin/pmall/rpc_gen/order/orderservice"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/pkg/transmeta"
	"github.com/cloudwego/kitex/server"
	kitexconsul "github.com/kitex-contrib/registry-consul"
)

func main() {
	dal.Init()
	rpc.Init()
	opts := kitexInit()

	logFile, err := os.OpenFile(conf.GetConf().Kitex.LogFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer logFile.Close()
	fileWriter := io.MultiWriter(logFile, os.Stdout)
	klog.SetOutput(fileWriter)
	klog.SetLevel(conf.LogLevel())

	svr := order.NewServer(new(OrderServiceImpl), opts...)

	err = svr.Run()

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

	// registry
	r, err := kitexconsul.NewConsulRegister(conf.GetConf().Registry.RegistryAddress[0])
	if err != nil {
		panic(err)
	}
	opts = append(opts, server.WithRegistry(r))
	opts = append(opts, server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{
		ServiceName: conf.GetConf().Kitex.Service,
	}))

	opts = append(opts, server.WithMetaHandler(transmeta.ServerHTTP2Handler))
	return
}
