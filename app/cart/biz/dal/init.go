package dal

import (
	"github.com/PiaoAdmin/pmall/app/cart/biz/dal/redis"
)

func Init() {
	redis.Init()
}
