package dal

import (
	"github.com/PiaoAdmin/pmall/app/user/biz/dal/mysql"
)

func Init() {
	mysql.Init()
}
