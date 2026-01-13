package utils

import (
	"context"
	"time"

	"github.com/PiaoAdmin/pmall/app/api/biz/dal/redis"
	"github.com/PiaoAdmin/pmall/app/api/md/jwt"
	perrors "github.com/PiaoAdmin/pmall/common/errs"
	"github.com/cloudwego/hertz/pkg/app"
)

func AddTokenToBlacklist(ctx context.Context, c *app.RequestContext) error {
	claims, err := jwt.JwtMiddleware.CheckIfTokenExpire(ctx, c)
	if err != nil {
		return err
	}
	tokenStr := jwt.GetToken(ctx, c)
	if tokenStr == "" {
		return perrors.New(perrors.ErrAuthFailed.Code, "token is empty")
	}

	var deadline time.Time

	if origIatFloat, ok := claims["orig_iat"].(float64); ok {
		origIat := int64(origIatFloat)
		deadline = time.Unix(origIat, 0).Add(jwt.JwtMiddleware.MaxRefresh).Add(jwt.JwtMiddleware.Timeout)
	} else {
		if expFloat, ok := claims["exp"].(float64); ok {
			deadline = time.Unix(int64(expFloat), 0)
		}
	}
	ttl := time.Until(deadline)
	if ttl > 0 {
		redis.RedisClient.Set(ctx, "blacklist:token:"+tokenStr, "1", ttl)
	}
	return nil
}
