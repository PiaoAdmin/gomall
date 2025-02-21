package dal

import (
	"github.com/PiaoAdmin/gomall/app/product/biz/dal/mysql"
	"github.com/PiaoAdmin/gomall/app/product/biz/dal/redis"
)

func Init() {
	redis.Init()
	mysql.Init()
}
