package dal

import (
	"github.com/PiaoAdmin/gomall/app/order/biz/dal/mysql"
	"github.com/PiaoAdmin/gomall/app/order/biz/dal/redis"
)

func Init() {
	redis.Init()
	mysql.Init()
}
