package dal

import (
	"github.com/PiaoAdmin/gomall/app/hertz_gateway/biz/dal/mysql"
)

func Init() {
	// redis.Init()
	mysql.Init()
}
