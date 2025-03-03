package service

import (
	"context"

	"github.com/PiaoAdmin/gomall/app/auth/biz/utils"
	auth "github.com/PiaoAdmin/gomall/rpc_gen/kitex_gen/auth"
)

type RefreshTokenService struct {
	ctx context.Context
} // NewRefreshTokenService new RefreshTokenService
func NewRefreshTokenService(ctx context.Context) *RefreshTokenService {
	return &RefreshTokenService{ctx: ctx}
}

// Run create note info
func (s *RefreshTokenService) Run(req *auth.RefreshTokenRequest) (resp *auth.RefreshTokenResponse, err error) {
	newAToken, newRToken, err := utils.NewARJWT().RefreshToken(req.Token, req.RefreshToken)
	if err != nil {
		if err.Error() == "access token still valid" {
			return &auth.RefreshTokenResponse{
				Token:        req.Token,
				RefreshToken: req.RefreshToken,
				IsValid:      true,
				Msg:          "access token still valid",
			}, nil
		}
		return nil, err
	}
	return &auth.RefreshTokenResponse{
		IsValid:      true,
		Msg:          "token refreshed successfully",
		Token:        newAToken,
		RefreshToken: newRToken,
	}, nil
}
