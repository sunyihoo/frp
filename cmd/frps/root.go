package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/sunyihoo/frp/pkg/util/version"
)

var (
	showVersion bool
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
			svrCfg *
		)
	},
}

func Execute() {

}
