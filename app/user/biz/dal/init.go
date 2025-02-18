package dal

import (
	"github.com/PiaoAdmin/gomall/app/user/biz/dal/mysql"
)

func Init() {
	// redis.Init()
	mysql.Init()
}
