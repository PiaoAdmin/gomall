package main

import (
	"database/sql"
	"gomall/biz/dal/queries/order"
	"gomall/biz/service/order"
	"gomall/handler/order"

	"net"

	"github.com/cloudwego/hertz/pkg/server"
	"github.com/cloudwego/kitex-contrib/consul-resolver"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
    db, err := sql.Open("mysql", "root:jinitaimei114514@tcp(39.103.237.155:10112)/gomall")
    if err != nil {
        panic(err)
    }
    defer db.Close()

    orderQuery := order.NewOrderQuery(db)
    orderService := order.NewOrderService(orderQuery)
    orderHandler := order.NewOrderHandler(orderService)

    // HTTP 服务
    r := server.New(server.WithHostPorts(":8000"))
    r.POST("/order/place", orderHandler.PlaceOrder)
    r.GET("/order/list", orderHandler.ListOrder)
    r.POST("/order/mark_paid", orderHandler.MarkOrderPaid)

    go func() {
        if err := r.Start(); err != nil {
            panic(err)
        }
    }()

    // Kitex 服务
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

    if err := svr.Run(); err != nil {
        panic(err)
    }
}
