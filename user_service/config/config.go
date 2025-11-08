package config

import (
	"fmt"
	"time"

	"github.com/Chateaubriand-g/bili/user_service/config"
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
		Port    string        `mapstructure:"port"`
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

	Zipkin struct {
		URL         string  `mapstructure:"url"`
		ServiceName string  `mapstructure:"servicename"`
		SampleRate  float64 `mapstructure:"samplerate"`
	} `mapstructure:"zipkin"`
}

func LoadConfig() (*config.Config, error) {
	viper.SetConfigName("configs")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("viper readinconfig failed: %w", err)
	}

	viper.SetEnvPrefix("USER")
	viper.AutomaticEnv()

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("viper unmarshal failed: %w", err)
	}

	return &config, nil
}
