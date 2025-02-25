package dal

import (
	"github.com/PiaoAdmin/gomall/app/checkout/biz/dal/mysql"
)

func Init() {
	// redis.Init()
	mysql.Init()
}
