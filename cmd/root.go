package cmd

import (
	"fmt"
	"github.com/spf13/viper"
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/siyibai/remote_debug/cmd/debug"
	"github.com/siyibai/remote_debug/config"
	"github.com/spf13/cobra"
)

var cfgFile string

var RootCmd = &cobra.Command{
	Use:   "remote_debug",
	Short: "goland dlv调试工具",
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	RootCmd.AddCommand(debug.DebugCmd)
	RootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", fmt.Sprintf("config file (default is $HOME/%s)",
		config.DefaultConfigName))
}

func initConfig() {
	if cfgFile == "" {
		home, err := homedir.Dir()
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		cfgFile = fmt.Sprintf("%s/%s", home, config.DefaultConfigName)
	}
	viper.SetConfigFile(cfgFile)
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()
}
