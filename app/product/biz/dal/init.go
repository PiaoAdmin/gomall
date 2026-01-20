package dal

import (
	"github.com/PiaoAdmin/pmall/app/product/biz/dal/mongo"
	"github.com/PiaoAdmin/pmall/app/product/biz/dal/mysql"
	"github.com/PiaoAdmin/pmall/app/product/biz/dal/redis"
	"github.com/PiaoAdmin/pmall/app/product/biz/model"
	"github.com/PiaoAdmin/pmall/app/product/conf"
)

func Init() {
	mysql.Init()
	redis.Init()
	mongo.Init()

	if conf.GetEnv() == "test" {
		model.InitProductDetailIndexes()
	}
}
