package config

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	v1 "github.com/sunyihoo/frp/pkg/config/v1"
	"strings"
)

func WordSepNormalizeFunc(f *pflag.FlagSet, name string) pflag.NormalizedName {
	if strings.Contains(name, "_") {
		return pflag.NormalizedName(strings.ReplaceAll(name, "_", "-"))
	}
	return pflag.NormalizedName(name)
}

type RegisterFlagOption func(*registerFlagOptions)

type registerFlagOptions struct {
	sshMode bool
}

func WithSSHMode() RegisterFlagOption {
	return func(o *registerFlagOptions) {
		o.sshMode = true
	}
}

// todo 未完成RegisterProxyFlags
func RegisterProxyFlags(cmd *cobra.Command, c *v1.ProxyConfigurer, opts ...RegisterFlagOption) {

}

// todo RegisterClientCommonConfigFlags
func RegisterClientCommonConfigFlags(cmd *cobra.Command, c *v1.ClientCommonConfig, opts ...RegisterFlagOption) {

}

// Todo 未完成RegisterServerConfigFlags
func RegisterServerConfigFlags(cmd *cobra.Command, c *v1.ServerConfig, opts ...RegisterFlagOption) {

}
