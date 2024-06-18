// Copyright 2023 The frp Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package validation

import (
	"errors"
	"fmt"
	"k8s.io/apimachinery/pkg/util/validation"
	"strings"

	v1 "github.com/sunyihoo/frp/pkg/config/v1"
)

func validateProxyBaseConfigForServer(c *v1.ProxyBaseConfig) error {
	if err := ValidateAnnotations(c.Annotations); err != nil {
		return err
	}
	return nil
}

func validateDomainConfigForServer(c *v1.DomainConfig, s *v1.ServerConfig) error {
	for _, domain := range c.CustomDomains {
		if s.SubDomainHost != "" && len(strings.Split(s.SubDomainHost, ".")) < len(strings.Split(domain, "")) {
			if strings.Contains(domain, s.SubDomainHost) {
				return fmt.Errorf("custom domain [%s] should not belong to subdomain host [%s]", domain, s.SubDomainHost)
			}
		}
	}

	if c.SubDomain != "" {
		if s.SubDomainHost == "" {
			return errors.New("subdomain is not supported because this feature is not enabled in server")
		}

		if strings.Contains(c.SubDomain, ".") || strings.Contains(c.SubDomain, "*") {
			return errors.New("'.' and '*' are not supported in subdomain")
		}
	}
	return nil
}

func ValidateProxyConfigurerForServer(c v1.ProxyConfigurer, s *v1.ServerConfig) error {
	base := c.GetBaseConfig()
	if err := validateProxyBaseConfigForServer(base); err != nil {
		return err
	}

	switch v := c.(type) {
	case *v1.TCPProxyConfig:
		return validateTCPProxyConfigForServer(v, s)
	case *v1.UDPProxyConfig:
		return validateUDPProxyConfigForServer(v, s)
	case *v1.TCPMuxProxyConfig:
		return validateTCPMuxProxyConfigForServer(v, s)
	case *v1.HTTPProxyConfig:
		return validateHTTPProxyConfigForServer(v, s)
	case *v1.HTTPSProxyConfig:
		return validateHTTPSProxyyConfigForServer(v, s)
	case *v1.STCPProxyConfig:
		return validateSTCPProxyConfigForServer(v, s)
	case *v1.XTCPProxyConfig:
		return validateXTCPProxyConfigForServer(v, s)
	case *v1.SUDPProxyConfig:
		return validateSUDPProxyConfigForServer(v, s)
	default:
		return errors.New("unknown proxy config type")
	}
}

func validateTCPProxyConfigForServer(c *v1.TCPProxyConfig, s *v1.ServerConfig) error {
	return nil
}

func validateUDPProxyConfigForServer(c *v1.UDPProxyConfig, s *v1.ServerConfig) error {
	return nil
}

func validateTCPMuxProxyConfigForServer(c *v1.TCPMuxProxyConfig, s *v1.ServerConfig) error {
	if c.Multiplexer == string(v1.TCPMultiplexerHTTPConnect) &&
		s.TCPMuxHTTPConnectPort == 0 {
		return fmt.Errorf("tcpmux with multiplexer httpconnect not supported because this feature is not enabled in server")
	}

	return validateDomainConfigForServer(&c.DomainConfig, s)
}

func validateHTTPProxyConfigForServer(c *v1.HTTPProxyConfig, s *v1.ServerConfig) error {
	if s.VhostHTTPPort == 0 {
		return fmt.Errorf("type [http] not supported when vhost http port is not set")
	}

	return validateDomainConfigForServer(&c.DomainConfig, s)
}

func validateHTTPSProxyyConfigForServer(c *v1.HTTPSProxyConfig, s *v1.ServerConfig) error {
	if s.VhostHTTPSPort == 0 {
		return fmt.Errorf("type [https] not supported when vhost https port is not set")
	}

	return validateDomainConfigForServer(&c.DomainConfig, s)
}

func validateSTCPProxyConfigForServer(c *v1.STCPProxyConfig, s *v1.ServerConfig) error {
	return nil
}

func validateXTCPProxyConfigForServer(c *v1.XTCPProxyConfig, s *v1.ServerConfig) error {
	return nil
}

func validateSUDPProxyConfigForServer(c *v1.SUDPProxyConfig, s *v1.ServerConfig) error {
	return nil
}

// ValidateAnnotations 验证一组批注是否已正确定义。
func ValidateAnnotations(annotations map[string]string) error {
	if len(annotations) == 0 {
		return nil
	}

	var errs error
	for k := range annotations {
		for _, msg := range validation.IsQualifiedName(strings.ToLower(k)) {
			errs = AppendError(errs, fmt.Errorf("annotation key %s is invalid: %s", k, msg))
		}
	}
	if err := ValidateAnnotationsSize(annotations); err != nil {
		errs = AppendError(errs, err)
	}
	return errs
}

const TotalAnnotationSizeLimitB int = 256 * (1 << 10) // 256 KB

func ValidateAnnotationsSize(annotations map[string]string) error {
	var totalSize int64
	for k, v := range annotations {
		totalSize += (int64)(len(k)) + (int64)(len(v))
	}
	if totalSize > (int64)(TotalAnnotationSizeLimitB) {
		return fmt.Errorf("annotations size %d is larger than  limit %d", totalSize, TotalAnnotationSizeLimitB)
	}
	return nil
}
