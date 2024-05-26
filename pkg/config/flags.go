package config

import (
	"github.com/spf13/cobra"
	v1 "github.com/sunyihoo/frp/pkg/config/v1"
)

type RegisterFlagOption func(*registerFlagOptions)

type registerFlagOptions struct {
	sshMode bool
}

// Todo RegisterServerConfigFlags
func RegisterServerConfigFlags(cmd *cobra.Command, c *v1.ServerConfig, opts ...RegisterFlagOption) {

}
