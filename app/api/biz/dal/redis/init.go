package redis

import (
	"context"
	"fmt"

	"github.com/PiaoAdmin/pmall/app/api/conf"
	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func Init() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     conf.GetConf().Redis.Address,
		Username: conf.GetConf().Redis.Username,
		Password: conf.GetConf().Redis.Password,
		DB:       conf.GetConf().Redis.DB,
	})
	if err := RedisClient.Ping(context.Background()).Err(); err != nil {
		fmt.Print("\n+++++++++++++++++++++++++++++++++++++++++++++\n")
		panic(err)
	}
}
