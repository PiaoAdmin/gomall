package utils

import (
	"context"
	"fmt"

	"github.com/cloudwego/hertz/pkg/app"
)

func GetUserIdFromToken(ctx context.Context, c *app.RequestContext) (user_id *int64, err error) {
	token := string(c.GetHeader("token"))
	claim, err := NewARJWT().ParseAccessToken(token)
	if err != nil {
		fmt.Printf("parse token error: %v", err)
		return nil, err
	}
	return &claim.UserID, nil
}
