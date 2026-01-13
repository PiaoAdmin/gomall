package service

import (
	"context"

	auth "github.com/PiaoAdmin/pmall/app/api/biz/model/api/auth"
	"github.com/PiaoAdmin/pmall/app/api/biz/utils"
	"github.com/cloudwego/hertz/pkg/app"
)

type LogoutService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewLogoutService(Context context.Context, RequestContext *app.RequestContext) *LogoutService {
	return &LogoutService{RequestContext: RequestContext, Context: Context}
}

func (h *LogoutService) Run(req *auth.Empty) (resp *auth.Empty, err error) {
	err = utils.AddTokenToBlacklist(h.Context, h.RequestContext)
	if err != nil {
		return nil, err
	}

	return &auth.Empty{}, nil
}
