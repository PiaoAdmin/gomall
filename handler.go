/*
 * @Author: liaosijie
 * @Date: 2025-02-16 10:15:55
 * @Last Modified by:   liaosijie
 * @Last Modified time: 2025-02-16 13:15:55
 */

package main

import (
	"database/sql"
	"gomall/biz/dal/queries/order"
	"gomall/biz/service/order"
	"gomall/handler/order"

	"github.com/cloudwego/hertz/pkg/server"
	"github.com/cloudwego/kitex-examples/kitex_gen/order"
	"github.com/cloudwego/kitex/server"
	_ "github.com/go-sql-driver/mysql"
)

func handler() {
    db, err := sql.Open("mysql", "root:jinitaimei114514@tcp(39.103.237.155:10112)/gomall")
    if err != nil {
        panic(err)
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

    err = r.Start()
    if err != nil {
        panic(err)
    }
}