package config

import (
	"fmt"
	"strings"

	"github.com/Chateaubriand-g/bili/common/config"
	"github.com/spf13/viper"
)

func LoadConfig() (*config.Config, error) {
	viper.AddConfigPath("./config")
	viper.SetConfigName("configs")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("viper readinconfig error: %w", err)
	}

	viper.SetEnvPrefix("INTER")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	var config config.Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("viper unmarshal config error: %w", err)
	}
	return &config, nil
}
