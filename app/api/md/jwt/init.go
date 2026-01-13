package jwt

import (
	"context"
	"time"

	"github.com/PiaoAdmin/pmall/app/api/biz/dal/redis"
	"github.com/PiaoAdmin/pmall/app/api/conf"
	"github.com/PiaoAdmin/pmall/app/api/pkg/response"
	perrors "github.com/PiaoAdmin/pmall/common/errs"
	"github.com/cloudwego/hertz/pkg/app"
)

var (
	JwtMiddleware *HertzJWTMiddleware
)

func Init() {
	var err error
	// 检查 Redis 是否已初始化
	if redis.RedisClient == nil {
		panic("JWT Init Error: redis.RedisClient is nil")
	}
	jwtConf := conf.GetConf().JWT
	JwtMiddleware, err = New(&HertzJWTMiddleware{
		Realm:             jwtConf.Realm,
		Key:               []byte(jwtConf.Key),
		Timeout:           time.Minute * time.Duration(jwtConf.Timeout),
		MaxRefresh:        time.Minute * time.Duration(jwtConf.MaxRefresh),
		IdentityKey:       jwtConf.IdentityKey,
		TokenLookup:       jwtConf.TokenLookup,
		DisabledAbort:     false,
		SendAuthorization: true,
		PayloadFunc: func(data interface{}) MapClaims { // 负载
			if v, ok := data.(uint64); ok {
				return MapClaims{
					jwtConf.IdentityKey: v,
				}
			}
			return MapClaims{}
		},
		Authorizator: func(data interface{}, ctx context.Context, c *app.RequestContext) bool { // Token 认证
			token := GetToken(ctx, c)
			exist, _ := redis.RedisClient.Exists(ctx, "blacklist:token:"+token).Result()
			if exist > 0 {
				return false
			}
			return true
		},
		Unauthorized: func(ctx context.Context, c *app.RequestContext, code int, message string) {
			response.FailWithErrorType(c, code, perrors.ErrAuthFailed.Code, message)
		},
	})
	if err != nil {
		panic("JWT Error:" + err.Error())
	}
	if JwtMiddleware == nil {
		panic("JWT Init Error: JwtMiddleware is nil after New()")
	}
}
