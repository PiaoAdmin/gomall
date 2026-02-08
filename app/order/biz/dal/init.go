package dal

import (
	"github.com/PiaoAdmin/pmall/app/order/biz/dal/mysql"
	"github.com/PiaoAdmin/pmall/app/order/biz/dal/rabbitmq"
)

func Init() {
	mysql.Init()
	rabbitmq.Init()
}

func Close() {
	rabbitmq.Close()
}
