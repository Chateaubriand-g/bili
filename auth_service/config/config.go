package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	DSN struct {
		User     string `mapstructure:"user"`
		Password string `mapstructure:"password"`
		Host     string `mapstructure:"host"`
		Port     string `mapstructure:"port"`
	} `mapstructure:"dsn"`
}

func LoadConfig() (*Config, error) {
	viper.SetConfigName("configs")
	viper.AddConfigPath("./config")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("viper raedinconfig failed: %w", err)
	}

	viper.AutomaticEnv()
	viper.SetEnvPrefix("AUTH")

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("viper unmarshal failed: %w", err)
	}

	return &config, nil
}
