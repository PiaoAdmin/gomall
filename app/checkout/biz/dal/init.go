package dal

import (
	"github.com/PiaoAdmin/gomall/app/checkout/biz/dal/mysql"
	"github.com/PiaoAdmin/gomall/app/checkout/biz/dal/redis"
)

func Init() {
	redis.Init()
	mysql.Init()
}
