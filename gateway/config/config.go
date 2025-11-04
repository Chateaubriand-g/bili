package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Gateway struct {
		Addr string `mapstructure:"addr"`
	} `mapstructure:"gateway"`

	Consul struct {
		Addr    string        `mapstructure:"addr"`
		Token   string        `mapstructure:"token"`
		Scheme  string        `mapstructure:"scheme"`
		Timeout time.Duration `mapstructure:"timeout"`
	} `mapstructure:"consul"`

	Jwt struct {
		Secret string `mapstructure:"secret"`
	} `mapstructure:"jwt"`
}

func LoadConfig() (*Config, error) {
	viper.SetConfigName("configs")
	viper.AddConfigPath("./config")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("viper readinconfig failed: %w", err)
	}
	viper.AutomaticEnv()
	viper.SetEnvPrefix("GATE")

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("viper unmarshal failed: %w", err)
	}

	return &config, nil
}
