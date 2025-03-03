package middleware

import (
	"context"
	"fmt"

	"github.com/PiaoAdmin/gomall/app/auth/biz/utils"
	resputils "github.com/PiaoAdmin/gomall/app/hertz_gateway/biz/utils"
	"github.com/PiaoAdmin/gomall/app/hertz_gateway/infra/rpc"
	"github.com/PiaoAdmin/gomall/common/constant"
	rpcauth "github.com/PiaoAdmin/gomall/rpc_gen/kitex_gen/auth"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

func hasCommonValue(list1, list2 []int32) bool {
	// 创建一个map存储list1的元素
	numSet := make(map[int32]struct{})
	for _, num := range list1 {
		numSet[num] = struct{}{}
	}

	// 遍历list2，检查是否有值在map中
	for _, num := range list2 {
		if _, exists := numSet[num]; exists {
			return true
		}
	}
	return false
}

func PermissionCheck(RoleCode []int32) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		// your code...
		token := c.GetHeader("Token")
		claims, err := utils.NewARJWT().ParseAccessToken(string(token))
		if err != nil {
			resputils.SendErrResponse(ctx, c, consts.StatusOK, err)
			c.Abort()
			return
		}
		res, err := rpc.AuthClient.GetUserRole(ctx, &rpcauth.GetUserRoleRequest{
			UserId: claims.UserID,
		})
		if err != nil {
			resputils.SendErrResponse(ctx, c, consts.StatusOK, err)
			c.Abort()
			return
		}
		if !hasCommonValue(RoleCode, res.RoleCode) {
			resputils.SendErrResponse(ctx, c, consts.StatusOK, constant.NoPermissionErr("用户权限不足"))
			c.Abort()
		}
		c.Next(ctx)
	}
}

func LoginCheck(t string) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		fmt.Printf("%s\n", t)
		token := c.GetHeader("Token")
		if len(token) == 0 {
			resputils.SendErrResponse(ctx, c, consts.StatusOK, constant.NotLoginError("未登录或登录过期"))
			c.Abort()
			return
		}
		res, err := rpc.AuthClient.VerifyToken(ctx, &rpcauth.VerifyTokenRequest{
			Token: string(token),
		})
		if err != nil || res == nil || !res.IsValid {
			resputils.SendErrResponse(ctx, c, consts.StatusOK, constant.NotLoginError("未登录或登录过期"))
			c.Abort()
			return
		}
		c.Next(ctx)
	}
}
