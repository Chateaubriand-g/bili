package config

import (
	"fmt"
	"time"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	DB struct {
		User     string `mapstructure:"user"`
		Password string `mapstructure:"password"`
		Host     string `mapstructure:"host"`
		Port     string `mapstructure:"port"`
	} `mapstructure:"db"`

	Consul struct {
		Addr    string        `mapstructure:"addr"`
		Port	string	      `mapstructure:"port"`
		Token   string        `mapstructure:"token"`
		Scheme  string        `mapstructure:"scheme"`
		Timeout time.Duration `mapstructure:"timeout"`
	} `mapstructure:"consul"`

	Server struct {
		ID   string `mapstructure:"id"`
		Name string `mapstructure:"name"`
		Addr string `mapstructure:"addr"`
		Port int    `mapstructure:"port"`
		/*
			HealthCheckPath string example /health
			HealthCheckInterval time.Duration
			HealthCheckTimeout time.Duration
		*/
	} `mapstructure:"server"`
}

func LoadConfig() (*Config, error) {
	viper.SetConfigName("configs")
	viper.AddConfigPath("./config")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("viper raedinconfig failed: %w", err)
	}

	viper.SetEnvPrefix("AUTH")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".","_"))
	viper.AutomaticEnv()

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("viper unmarshal failed: %w", err)
	}

	return &config, nil
}
