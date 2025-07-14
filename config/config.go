package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

const (
	DefaultConfigName = "config.yaml"
	DefaultTimeout    = 10 * time.Second
)

type RemoteAddress struct {
	RemoteIP   string        `mapstructure:"remote_ip"`
	RemotePort string        `mapstructure:"remote_port"`
	DestPath   string        `mapstructure:"dest_path"`
	Timeout    time.Duration `mapstructure:"timeout"`
}

type Config struct {
	ProjectPath     string        `mapstructure:"project_path"`
	BuildCmdFmt     string        `mapstructure:"build_cmd_fmt"`
	CrossCompileCmd string        `mapstructure:"cross_compile_cmd"`
	RunCmdFmt       string        `mapstructure:"run_cmd_fmt"`
	RemoteAddress   RemoteAddress `mapstructure:"remote_address"`
}

// GetConfig 使用全局viper解析yaml配置
func GetConfig() (*Config, error) {
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	c := &Config{}
	if err := viper.Unmarshal(c); err != nil {
		return nil, fmt.Errorf("解析配置失败: %w", err)
	}

	if err := checkConfig(c); err != nil {
		return nil, fmt.Errorf("配置不完整, 请检查, err: %s", err.Error())
	}
	return c, nil
}

func checkConfig(c *Config) error {
	if c.ProjectPath == "" {
		return fmt.Errorf("配置缺失:ProjectPath")
	}
	if c.RemoteAddress.RemoteIP == "" {
		return fmt.Errorf("配置缺失:RemoteAddress.RemoteIP")
	}
	if c.RemoteAddress.RemotePort == "" {
		return fmt.Errorf("配置缺失:RemoteAddress.RemotePort")
	}
	if c.RemoteAddress.DestPath == "" {
		return fmt.Errorf("配置缺失:RemoteAddress.DestPath")
	}
	if c.RemoteAddress.Timeout == 0 {
		c.RemoteAddress.Timeout = DefaultTimeout
	}
	return nil
}
