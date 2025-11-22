package internal

import (
	"fmt"
	"strings"

	"github.com/Chateaubriand-g/bili/common/config"
	"github.com/spf13/viper"
)

func LoadConfig() (*config.Config, error) {
	viper.SetConfigName("configs")
	viper.AddConfigPath("./")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("viper raedinconfig failed: %w", err)
	}

	viper.SetEnvPrefix("COMMENT")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	var config config.Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("viper unmarshal failed: %w", err)
	}

	return &config, nil
}
