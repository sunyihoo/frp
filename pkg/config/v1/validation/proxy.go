package validation

import v1 "github.com/sunyihoo/frp/pkg/config/v1"

func ValidateProxyConfigurerForServer(c v1.ProxyConfigurer, s *v1.ServerConfig) error {
	base := c.GetBaseConfig()
	if err := validateProxyBaseConfigForServer(base); err != nil {
		return err
	}

}
