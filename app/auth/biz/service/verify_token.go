package service

import (
	"context"

	"github.com/PiaoAdmin/gomall/app/auth/biz/utils"
	auth "github.com/PiaoAdmin/gomall/rpc_gen/kitex_gen/auth"
)

type VerifyTokenService struct {
	ctx context.Context
} // NewVerifyTokenService new VerifyTokenService
func NewVerifyTokenService(ctx context.Context) *VerifyTokenService {
	return &VerifyTokenService{ctx: ctx}
}

// Run create note info
func (s *VerifyTokenService) Run(req *auth.VerifyTokenRequest) (resp *auth.VerifyTokenResponse, err error) {
	// Finish your business logic.
	_, err = utils.NewARJWT().ParseAccessToken(req.Token)
	if err != nil {
		return &auth.VerifyTokenResponse{
			IsValid: false,
			Msg:     err.Error(),
		}, nil
	}

	return &auth.VerifyTokenResponse{
		IsValid: true,
		Msg:     "token is valid",
	}, nil
}
