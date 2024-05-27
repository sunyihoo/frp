package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/sunyihoo/frp/pkg/config"
	v1 "github.com/sunyihoo/frp/pkg/config/v1"
	"github.com/sunyihoo/frp/pkg/util/version"
	"os"
)

var (
	cfgFile          string
	showVersion      bool
	strictConfigMode bool

	serverCfg v1.ServerConfig
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file of frps")
	rootCmd.PersistentFlags().BoolVarP(&showVersion, "version", "v", false, "version of frps")
	rootCmd.PersistentFlags().BoolVarP(&strictConfigMode, "strict_config", "", true, "strict config parsing mode, unknown fields will cause error")

	config.RegisterServerConfigFlags(rootCmd, &serverCfg)
}

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
			svrCfg, isLegacyFormat, err = config.LoadServerConfig(cfgFile, strictConfigMode)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			if isLegacyFormat {
				fmt.Printf("WARING: ini format is deprecated and the support will be removed in the future, " +
					"please use yaml/json/toml format instead!\n")
			}
		} else {
			serverCfg.Complete()
			svrCfg = &serverCfg
		}
		// todo
		warning,err :=

		return nil
	},
}

func Execute() {

}
