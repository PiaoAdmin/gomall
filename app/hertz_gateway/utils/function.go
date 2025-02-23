package utils

import "context"

func GetUserIdFromCtx(ctx context.Context) int64 {
	if ctx.Value(UserIdKey) == nil {
		return 0
	}
	return ctx.Value(UserIdKey).(int64)
}
