package errs

import (
	"fmt"
	"strings"
)

type ErrorType uint64

type Error struct {
	Code    ErrorType
	Message string
	Extra   map[string]string
}

func (e *Error) Error() string {
	baseMsg := fmt.Sprintf("%s: %s", convertErrorTypeToString(e.Code), e.Message)

	if len(e.Extra) > 0 {
		var extraParts []string
		for k, v := range e.Extra {
			extraParts = append(extraParts, fmt.Sprintf("%s=%s", k, v))
		}
		return fmt.Sprintf("%s [%s]", baseMsg, strings.Join(extraParts, ", "))
	}

	return baseMsg
}

// 实现 Kitex BizStatusErrorIface 接口
func (e *Error) BizStatusCode() int32 {
	return int32(e.Code)
}

func (e *Error) BizMessage() string {
	return e.Message
}

func (e *Error) BizExtra() map[string]string {
	if e.Extra == nil {
		return make(map[string]string)
	}
	return e.Extra
}

func New(code ErrorType, msg string) *Error {
	return &Error{
		Code:    code,
		Message: msg,
		Extra:   nil,
	}
}

func NewWithExtra(code ErrorType, msg string, extra map[string]string) *Error {
	return &Error{
		Code:    code,
		Message: msg,
		Extra:   extra,
	}
}

var (
	Success = New(20000, "success")

	ErrInternal = New(50000, "internal server error")
	ErrRending  = New(50001, "rending error")

	ErrBing = New(40001, "binding error")

	ErrParam           = New(40002, "invalid params")
	ErrRecordNotFound  = New(40003, "record not found")
	ErrRecordAlreadyEx = New(40004, "record already exists")
	ErrAuthFailed      = New(40005, "authorization failed")
	ErrNotLogin        = New(40006, "user not login")
)

func convertErrorTypeToString(code ErrorType) string {
	template := "(" + fmt.Sprintf("%d", code) + ")"
	var str string
	switch code {
	case 20000:
		str = "SUCCESS"
	case 40002:
		str = "ERR_INVALID_PARAMS"
	case 40003:
		str = "ERR_RECORD_NOT_FOUND"
	case 40004:
		str = "ERR_RECORD_ALREADY_EXISTS"
	case 40005:
		str = "ERR_AUTH_FAILED"
	case 40006:
		str = "ERR_NOT_LOGIN"
	case 50000:
		str = "ERR_INTERNAL"
	}
	return str + template
}

func ConvertErr(err error) *Error {
	if err == nil {
		return Success
	}
	if e, ok := err.(*Error); ok {
		return e
	}
	return New(ErrInternal.Code, err.Error())
}
