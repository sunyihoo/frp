package validation

import (
	"errors"
	v1 "github.com/sunyihoo/frp/pkg/config/v1"
	splugin "github.com/sunyihoo/frp/pkg/plugin/server"
)

var (
	// SupportedAuthMethods 支持的身份验证方法
	SupportedAuthMethods = []v1.AuthMethod{
		"token",
		"oidc",
	}

	// SupportedAuthAdditionalScopes 支持的身份验证其他范围
	SupportedAuthAdditionalScopes = []v1.AuthScope{
		"HeartBeats",
		"NewWorkConns",
	}

	// SupportedLogLevels 支持的日志等级
	SupportedLogLevels = []string{
		"trace",
		"debug",
		"info",
		"warn",
		"error",
	}

	SupportedHTTPPlugins = []string{
		splugin.OpLogin,
		splugin.OpNewProxy,
		splugin.OpCloseProxy,
		splugin.OpPing,
		splugin.OpNewWorkConn,
		splugin.OpNewUserConn,
	}
)

type Warning error

func AppendError(err error, errs ...error) error {
	if len(errs) == 0 {
		return err
	}
	return errors.Join(append([]error{err}, errs...)...)
}
