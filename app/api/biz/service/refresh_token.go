package service

import (
	"context"

	perrors "github.com/PiaoAdmin/pmall/common/errs"

	"github.com/PiaoAdmin/pmall/app/api/biz/dal/redis"
	"github.com/PiaoAdmin/pmall/app/api/biz/model/api/auth"
	"github.com/PiaoAdmin/pmall/app/api/md/jwt" // 引入本地封装的 jwt 包
	"github.com/cloudwego/hertz/pkg/app"
	jwtGo "github.com/golang-jwt/jwt/v4" // 引入官方包用于判断错误类型
)

type RefreshTokenService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewRefreshTokenService(Context context.Context, RequestContext *app.RequestContext) *RefreshTokenService {
	return &RefreshTokenService{RequestContext: RequestContext, Context: Context}
}

func (h *RefreshTokenService) Run(req *auth.Empty) (resp *auth.RefreshResp, err error) {
	// 1. 尝试解析 Token
	// ParseToken 会自动验证签名和过期时间
	tokenObj, err := jwt.JwtMiddleware.ParseToken(h.Context, h.RequestContext)

	// 2. 关键判断：Token 是否有效
	// 如果 err 为 nil，说明 Token 验证通过，完全没有过期 --> 拒绝刷新
	if err == nil {
		return nil, perrors.New(perrors.ErrAuthFailed.Code, "token is still valid, no need to refresh")
	}

	// 3. 判断是否是 "过期错误"
	// 我们只允许 "过期" 这一种错误通过，如果是签名错误等其他问题，直接报错
	validationErr, ok := err.(*jwtGo.ValidationError)
	if !ok || validationErr.Errors != jwtGo.ValidationErrorExpired {
		return nil, err
	}

	// 4. 黑名单检查
	// 既然是刷新，我们需要确保旧 Token 没有被拉黑
	tokenStr := ""
	if tokenObj != nil {
		tokenStr = tokenObj.Raw
	}
	if tokenStr != "" {
		exist, _ := redis.RedisClient.Exists(h.Context, "blacklist:token:"+tokenStr).Result()
		if exist > 0 {
			return nil, perrors.New(perrors.ErrNotLogin.Code, "token has been blacklisted")
		}
	}

	// 5. 复用核心逻辑：调用 JwtMiddleware.RefreshToken
	// 这个函数内部会再次调用 CheckIfTokenExpire：
	// - 检查是否在 MaxRefresh 可刷新窗口内
	// - 生成新 Token (继承旧 Claims)
	tokenString, expire, err := jwt.JwtMiddleware.RefreshToken(h.Context, h.RequestContext)
	if err != nil {
		// 如果这里报错，通常意味着超过了 MaxRefresh 时间，或者生成失败
		return nil, err
	}

	return &auth.RefreshResp{
		Token:    tokenString,
		ExpireIn: expire.Unix(),
	}, nil
}
