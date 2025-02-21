/**
 * @Author: ZhangHaoChen
 * @Date:   2/21/25 PM5:11
=======
/*
 * @Author: liaosijie
 * @Date: 2025-02-20 17:22:10
 * @Last Modified by: liaosijie
 * @Last Modified time: 2025-02-20 18:10:15
 */

package rpc

import (
	"github.com/PiaoAdmin/gomall/rpc_gen/kitex_gen/product/productcatalogservice"
	"github.com/cloudwego/kitex/client"
	consul "github.com/kitex-contrib/registry-consul"
	"log"
	"sync"
)

var (
	ProductClient productcatalogservice.Client
	once          sync.Once
)

func Init() {
	once.Do(func() {
		initUserClient()
	})
}
func initUserClient() {
	resolver, err := consul.NewConsulResolver("127.0.0.1:8500")
	if err != nil {
		log.Fatal(err)
	}
	ProductClient, err = productcatalogservice.NewClient("product", client.WithResolver(resolver))
	if err != nil {
		log.Fatal(err)
	}
}