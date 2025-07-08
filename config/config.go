package config

import (
	"fmt"
	"time"

	"github.com/knadh/koanf/v2"
)

var KfObject *koanf.Koanf

const (
	DefaultConfigName = "./config.yaml"
	DefaultTimeout    = 10 * time.Second
)

type Config struct {
	ProjectPath   string `yaml:"project_path"`
	RemoteAddress struct {
		RemoteIP   string        `yaml:"remote_ip"`
		RemotePort string        `yaml:"remote_port"`
		DestPath   string        `yaml:"dest_path"`
		Timeout    time.Duration `yaml:"timeout"`
	}
}

// GetConfig 获取配置
func GetConfig() (*Config, error) {
	c := &Config{}
	err := KfObject.Unmarshal("", c)
	if err != nil {
		return nil, err
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
