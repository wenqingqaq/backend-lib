package tecutil

import (
	"gitee.com/yanwenqing/backend-lib/tecons"
	"github.com/go-kratos/kratos/v2/errors"
)

func InternalErr(e error) error {
	return errors.InternalServer(tecons.InternalServer, e.Error()).WithCause(e)
}

func UnauthorizedErr(e error) error {
	return errors.Unauthorized(tecons.Unauthorized, e.Error())
}

func Forbidden() error {
	return errors.Forbidden(tecons.Forbidden, "")
}
