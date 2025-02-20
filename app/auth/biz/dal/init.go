package dal

import (
	"github.com/PiaoAdmin/gomall/app/auth/biz/dal/mysql"
)

func Init() {
	// redis.Init()
	mysql.Init()
}
