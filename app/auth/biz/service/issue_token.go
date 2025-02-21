package service

import (
	"context"

	"github.com/PiaoAdmin/gomall/app/auth/biz/dal/model"
	"github.com/PiaoAdmin/gomall/app/auth/biz/dal/mysql"
	"github.com/PiaoAdmin/gomall/app/auth/biz/utils"
	auth "github.com/PiaoAdmin/gomall/rpc_gen/kitex_gen/auth"
	"golang.org/x/crypto/bcrypt"
)

type IssueTokenService struct {
	ctx context.Context
} // NewIssueTokenService new IssueTokenService
func NewIssueTokenService(ctx context.Context) *IssueTokenService {
	return &IssueTokenService{ctx: ctx}
}

type UserClaims struct {
	Username string `json:"username"`
	UserID   int64  `json:"user_id"`
}

// Run create note info
func (s *IssueTokenService) Run(req *auth.IssueTokenRequest) (resp *auth.IssueTokenResponse, err error) {
	// Finish your business logic.
	// TODO: 验证用户名密码 统一配置
	u, err := model.GetUserByUsername(mysql.DB, s.ctx, req.Username)
	if err != nil {
		return
	}
	if u == nil {
		return &auth.IssueTokenResponse{
			Success: false,
			Msg:     "用户不存在",
		}, err
	}
	if err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(req.Password)); err != nil {
		return &auth.IssueTokenResponse{
			Success: false,
			Msg:     "旧密码错误",
		}, err
	}
	accessTokenString, refreshTokenString, err := utils.NewARJWT().GenerateToken(u.ID, u.Username)
	if err != nil {
		return
	}
	return &auth.IssueTokenResponse{
		Token:        accessTokenString,
		RefreshToken: refreshTokenString,
		Success:      true,
		Msg:          "Issue token success",
	}, nil
}
