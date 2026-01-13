package response

import (
	perrors "github.com/PiaoAdmin/pmall/common/errs"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

type Response struct {
	Code    uint64      `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func Success(c *app.RequestContext, data interface{}) {
	c.JSON(consts.StatusOK, Response{
		Code:    20000,
		Message: "success",
		Data:    data,
	})
}

func Fail(c *app.RequestContext, httpStatusCode int, code uint64, message string, data interface{}) {
	c.JSON(httpStatusCode, Response{
		Code:    code,
		Message: message,
		Data:    data,
	})
}

func FailWithMessage(c *app.RequestContext, httpStatusCode int, code uint64, message string) {
	failWithMessage(c, httpStatusCode, code, message)
}

func FailWithErrorType(c *app.RequestContext, httpStatusCode int, code perrors.ErrorType, message string) {
	failWithMessage(c, httpStatusCode, uint64(code), message)
}

func failWithMessage(c *app.RequestContext, httpStatusCode int, code uint64, message string) {
	c.JSON(httpStatusCode, Response{
		Code:    code,
		Message: message,
		Data:    nil,
	})
}
