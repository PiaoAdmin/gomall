package constant

import "github.com/cloudwego/kitex/pkg/kerrors"

type EerrorCode int32

const (
	ReqIsNil EerrorCode = iota + 4001
	ParameterError
	NotLogin
	NoPermission
)

func ReqIsNilError(msg string, err ...error) kerrors.BizStatusErrorIface {
	if len(err) > 0 && err[0] != nil {
		return kerrors.NewBizStatusErrorWithExtra(int32(ReqIsNil), msg, map[string]string{"err": err[0].Error()})
	}
	return kerrors.NewBizStatusError(int32(ReqIsNil), msg)
}

func ParametersError(msg string, err ...error) kerrors.BizStatusErrorIface {
	if len(err) > 0 && err[0] != nil {
		return kerrors.NewBizStatusErrorWithExtra(int32(ReqIsNil), msg, map[string]string{"err": err[0].Error()})
	}
	return kerrors.NewBizStatusError(int32(ParameterError), msg)
}

func NotLoginError(msg string, err ...error) kerrors.BizStatusErrorIface {
	if len(err) > 0 && err[0] != nil {
		return kerrors.NewBizStatusErrorWithExtra(int32(ReqIsNil), msg, map[string]string{"err": err[0].Error()})
	}
	return kerrors.NewBizStatusError(int32(NotLogin), msg)
}

func NoPermissionErr(msg string, err ...error) kerrors.BizStatusErrorIface {
	if len(err) > 0 && err[0] != nil {
		return kerrors.NewBizStatusErrorWithExtra(int32(ReqIsNil), msg, map[string]string{"err": err[0].Error()})
	}
	return kerrors.NewBizStatusError(int32(NoPermission), msg)
}
