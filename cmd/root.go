package cmd

import (
	"fmt"
	"os"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
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
	config.KfObject = koanf.New(".")
	if cfgFile == "" {
		home, err := homedir.Dir()
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		cfgFile = fmt.Sprintf("%s/%s", home, config.DefaultConfigName)
	}

	if err := config.KfObject.Load(file.Provider(cfgFile), yaml.Parser()); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
