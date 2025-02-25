package dal

import (
	"github.com/PiaoAdmin/gomall/app/product/biz/dal/mysql"
)

func Init() {
	// redis.Init()
	mysql.Init()
}
