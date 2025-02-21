package dal

import (
	"github.com/PiaoAdmin/gomall/app/cart/biz/dal/mysql"
)

func Init() {
	//redis.Init()
	mysql.Init()
}
