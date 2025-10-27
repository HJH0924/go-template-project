// Package config provides configuration management for the application.
package config

import (
	"fmt"
	"log/slog"
	"sync"

	"github.com/spf13/viper"
)

var (
	instance *Config
	once     sync.Once
)

// Config 应用配置.
type Config struct {
	Server ServerConfig `mapstructure:"server"`
}

// ServerConfig 服务器配置.
type ServerConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

// Load 从指定路径加载配置文件到全局单例.
func Load(configPath string) error {
	var err error

	once.Do(func() {
		v := viper.New()

		// 设置配置文件路径
		v.SetConfigFile(configPath)
		v.SetConfigType("yaml")

		// 读取配置文件
		if err = v.ReadInConfig(); err != nil {
			err = fmt.Errorf("failed to read config file: %w", err)
			return
		}

		slog.Info("config file loaded", slog.String("path", v.ConfigFileUsed()))

		// 解析配置到结构体
		var cfg Config
		if err = v.Unmarshal(&cfg); err != nil {
			err = fmt.Errorf("failed to unmarshal config: %w", err)
			return
		}

		instance = &cfg
	})

	return err
}

// Get 获取全局配置实例.
func Get() *Config {
	if instance == nil {
		panic("config not loaded, please call Load() first")
	}

	return instance
}

// Address 返回服务器地址.
func (s *ServerConfig) Address() string {
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}
