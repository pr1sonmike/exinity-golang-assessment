package config

import (
	"github.com/spf13/viper"
	"time"
)

type GatewayConfig struct {
	Name        string        `mapstructure:"name"`
	APIEndpoint string        `mapstructure:"api_endpoint"`
	APITimeout  time.Duration `mapstructure:"api_timeout"`
	Auth        struct {
		Username string `mapstructure:"username"`
		Password string `mapstructure:"password"`
	} `mapstructure:"auth"`
}

type Config struct {
	Database struct {
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
		User     string `mapstructure:"user"`
		Password string `mapstructure:"password"`
		DBName   string `mapstructure:"dbname"`
		SSLMode  string `mapstructure:"sslmode"`
	} `mapstructure:"database"`
	Gateways []GatewayConfig `mapstructure:"gateways"`
	Server   struct {
		Port string `mapstructure:"port"`
	} `mapstructure:"server"`
}

func LoadConfig(path string) (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(path)
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
