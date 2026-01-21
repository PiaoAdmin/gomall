package dal

import (
	"github.com/PiaoAdmin/pmall/app/order/biz/dal/mysql"
)

func Init() {
	mysql.Init()
}
