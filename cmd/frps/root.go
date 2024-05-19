package main

import (
	"fmt"
	"github.com/spf13/cobra"
	v1 "github.com/sunyihoo/frp/pkg/config/v1"
	"github.com/sunyihoo/frp/pkg/util/version"
)

var (
	cfgFile string
	showVersion bool

	serverCfg v1.ServerConfig
)

var rootCmd = &cobra.Command{
	Use:   "frps",
	Short: "frps is the server of frp (https://github.com/fatedier/frp)",
	RunE: func(cmd *cobra.Command, args []string) error {
		if showVersion {
			fmt.Println(version.Full())
			return nil
		}

		var (
			svrCfg         *v1.ServerConfig
			isLegacyFormat bool
			err            error
		)
		if cfgFile != "" {
			svrCfg.isLegacyFormat,err = config.
		}
	},
}

func Execute() {

}
