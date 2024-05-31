package validation

import (
	"fmt"
	"github.com/samber/lo"
	v1 "github.com/sunyihoo/frp/pkg/config/v1"
	"slices"
)

func ValidateServerConfig(c *v1.ServerConfig) (Warning, error) {
	var (
		warnings Warning
		errs     error
	)
	if !slices.Contains(SupportedAuthMethods, c.Auth.Method) {
		errs = AppendError(errs, fmt.Errorf("invalid auth method, optional values are %v", SupportedAuthMethods))
	}
	if !lo.Every(SupportedAuthAdditionalScopes, c.Auth.AdditionalScopes) {
		errs = AppendError(errs, fmt.Errorf("invalid auth addtional scopes, optional values are %v", SupportedAuthAdditionalScopes))
	}

	if err := validateLogConfig(&c.Log); err != nil {
		errs = AppendError(errs, err)
	}

	if err := validateWebServerConfig(&c.WebServer); err != nil {
		errs = AppendError(errs, err)
	}

	errs = AppendError(errs, ValidatePort(c.BindPort, "bindPort"))
	errs = AppendError(errs, ValidatePort(c.KCPBindPort, "kcpBindPort"))
	errs = AppendError(errs, ValidatePort(c.QUICBindPort, "quicBindPort"))
	errs = AppendError(errs, ValidatePort(c.VhostHTTPPort, "vhostHTTPPort"))
	errs = AppendError(errs, ValidatePort(c.VhostHTTPSPort, "vhostHTTPSPort"))
	errs = AppendError(errs, ValidatePort(c.TCPMuxHTTPConnectPort, "tcpMuxHTTPConnectPort"))

	for _, p := range c.HTTPPlugins {
		if !lo.Every(SupportedHTTPPlugins, p.Ops) {
			errs = AppendError(errs, fmt.Errorf("invalid http plugin ops, optional values are %v", SupportedHTTPPlugins))
		}
	}
	return warnings, errs
}
