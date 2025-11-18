package main

import (
	"fmt"

	"github.com/spf13/viper"
)

type TargetConfig struct {
	Name string `mapstructure:"name"`
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type PushoverConfig struct {
	AppToken  string `mapstructure:"appToken"`
	UserToken string `mapstructure:"userToken"`
}

type Config struct {
	Interval uint           `mapstructure:"interval"`
	Targets  []TargetConfig `mapstructure:"targets"`

	Pushover PushoverConfig `mapstructure:"pushover"`
}

func LoadConfig(configPath string) (*Config, error) {
	viper.SetConfigFile(configPath)
	viper.SetDefault("interval", 10)
	viper.SetDefault("targets", []TargetConfig{})

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &config, nil
}
