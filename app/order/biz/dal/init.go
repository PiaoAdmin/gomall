package dal

import (
	"douyin-gomall/gomall/app/order/biz/dal/mysql"
	"douyin-gomall/gomall/app/order/biz/dal/redis"
)

func Init() {
	redis.Init()
	mysql.Init()
}
