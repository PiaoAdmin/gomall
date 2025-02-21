package dal

import (
	// "douyin-gomall/gomall/app/order/biz/dal/mysql"
	// "douyin-gomall/gomall/app/order/biz/dal/redis"
	"github.com/PiaoAdmin/gomall/app/order/biz/dal/mysql"
	"github.com/PiaoAdmin/gomall/app/order/biz/dal/redis"
)

func Init() {
	redis.Init()
	mysql.Init()
}
