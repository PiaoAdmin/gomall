/*
 * @Author: liaosijie
 * @Date: 2025-02-16 10:16:07
 * @Last Modified by:   liaosijie
 * @Last Modified time: 2025-02-16 13:16:07
 */

package main

import (
	"database/sql"
	"gomall/biz/dal/queries/order"
	"gomall/biz/service/order"
	"net"

	"github.com/cloudwego/kitex-examples/kitex_gen/order"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	"github.com/kitex-contrib/consul-resolver"
)

func main() {
    db, err := sql.Open("mysql", "root:jinitaimei114514@tcp(39.103.237.155:10112)/gomall")
    if err != nil {
        panic(err)
    }
    defer db.Close()

    orderQuery := order.NewOrderQuery(db)
    orderService := order.NewOrderService(orderQuery)

    resolver, err := consul.NewConsulResolver("39.103.237.155:8500")
    if err != nil {
        panic(err)
    }

    svr := ordersvr.NewServer(orderService, server.WithServiceAddr(&net.TCPAddr{
        IP:   net.ParseIP("0.0.0.0"),
        Port: 8000,
    }), server.WithRegistry(resolver), server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{
        ServiceName: "order_service",
    }))

    err = svr.Run()
    if err != nil {
        panic(err)
    }
}
