package redis

import (
	"context"

<<<<<<< HEAD
	"github.com/redis/go-redis/v9"
	// "douyin-gomall/gomall/app/order/conf"
	"github.com/PiaoAdmin/gomall/app/order/conf"
=======
	"github.com/PiaoAdmin/gomall/app/order/conf"
	"github.com/redis/go-redis/v9"
>>>>>>> b6e73c27fce12b01552c5334097a847176b8f26a
)

var (
	RedisClient *redis.Client
)

func Init() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     conf.GetConf().Redis.Address,
		Username: conf.GetConf().Redis.Username,
		Password: conf.GetConf().Redis.Password,
		DB:       conf.GetConf().Redis.DB,
	})
	if err := RedisClient.Ping(context.Background()).Err(); err != nil {
		panic(err)
	}
}
