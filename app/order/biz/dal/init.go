package dal

import (
	"github.com/PiaoAdmin/gomall/app/order/biz/dal/mysql"
)

func Init() {
	// redis.Init()
	mysql.Init()
}
