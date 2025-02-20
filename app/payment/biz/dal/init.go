package dal

import (
	"github.com/PiaoAdmin/gomall/app/payment/biz/dal/mysql"
)

func Init() {
	//redis.Init()
	mysql.Init()
}
