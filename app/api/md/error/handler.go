package error

import (
	"context"

	"github.com/PiaoAdmin/pmall/app/api/pkg/response"
	"github.com/PiaoAdmin/pmall/app/user/conf"
	perrors "github.com/PiaoAdmin/pmall/common/errs"
	"github.com/cloudwego/hertz/pkg/app"
	herrors "github.com/cloudwego/hertz/pkg/common/errors"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/cloudwego/kitex/pkg/kerrors"
)

func GlobalErrorHandler() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		c.Next(ctx)

		if len(c.Errors) == 0 {
			return
		}

		hlog.CtxErrorf(ctx, "errors: %s", c.Errors.String())

		lastErr := c.Errors.Last()

		statusCode, resp := resolveError(lastErr)
		c.JSON(statusCode, resp)
		// c.Abort()
	}
}

func resolveError(e *herrors.Error) (int, *response.Response) {
	if e.IsType(herrors.ErrorTypeBind) {
		return consts.StatusBadRequest, &response.Response{
			Code:    uint64(perrors.ErrBing.Code),
			Message: e.Error(),
			Data:    e.Meta,
		}
	}

	if e.IsType(herrors.ErrorTypeRender) {
		return consts.StatusInternalServerError, &response.Response{
			Code:    uint64(perrors.ErrRending.Code),
			Message: e.Error(),
			Data:    e.Meta,
		}
	}

	// 业务错误处理
	err := e.Unwrap()

	if bizErr, ok := kerrors.FromBizStatusError(err); ok {
		code := uint64(bizErr.BizStatusCode())
		var myErr *perrors.Error
		if len(bizErr.BizExtra()) > 0 {
			myErr = perrors.NewWithExtra(perrors.ErrorType(code), bizErr.BizMessage(), bizErr.BizExtra())
		} else {
			myErr = perrors.New(perrors.ErrorType(code), bizErr.BizMessage())
		}
		if code >= 50000 {
			return consts.StatusInternalServerError, &response.Response{
				Code:    code,
				Message: myErr.Error(),
				Data:    nil,
			}
		}
		return consts.StatusOK, &response.Response{
			Code:    code,
			Message: myErr.Error(),
			Data:    nil,
		}
	}

	if e, ok := err.(*perrors.Error); ok {
		if uint64(e.Code) >= 50000 {
			return consts.StatusInternalServerError, &response.Response{
				Code:    uint64(e.Code),
				Message: e.Error(),
				Data:    nil,
			}
		}
		return consts.StatusOK, &response.Response{
			Code:    uint64(e.Code),
			Message: e.Error(),
			Data:    nil,
		}
	}

	message := "Internal Server Error"
	if conf.GetConf().Env == "test" {
		message = err.Error()
	}
	// 默认错误处理
	return consts.StatusInternalServerError, &response.Response{
		Code:    uint64(perrors.ErrInternal.Code),
		Message: message,
		Data:    nil,
	}
}
