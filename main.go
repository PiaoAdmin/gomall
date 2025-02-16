/*
 * @Author: liaosijie
 * @Date: 2025-02-16 11:53:18
 * @Last Modified by: liaosijie
 * @Last Modified time: 2025-02-16 17:55:55
 */

package main

import (
	"database/sql"
	"gomall/biz/dal/queries/order"
	"gomall/biz/service/order"
	"gomall/handler/order"
	"log"
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
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	orderQuery := order.NewOrderQuery(db)
	orderService := order.NewOrderService(orderQuery)
	orderHandler := order.NewOrderHandler(orderService)

	r := server.New(server.WithHostPorts(":8000"))

	// 注册订单相关的路由
	r.POST("/order/place", orderHandler.PlaceOrder)
	r.GET("/order/list", orderHandler.ListOrder)
	r.POST("/order/mark_paid", orderHandler.MarkOrderPaid)

	resolver, err := consul.NewConsulResolver("39.103.237.155:8500")
	if err != nil {
		log.Fatalf("failed to create consul resolver: %v", err)
	}

	svr := ordersvr.NewServer(orderService, server.WithServiceAddr(&net.TCPAddr{
		IP:   net.ParseIP("0.0.0.0"),
		Port: 8000,
	}), server.WithRegistry(resolver), server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{
		ServiceName: "order_service",
	}))

	if err := svr.Run(); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}