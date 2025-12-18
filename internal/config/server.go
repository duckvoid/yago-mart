package config

import (
	"fmt"

	"github.com/spf13/pflag"
)

type ServerConfig struct {
	Address  string `mapstructure:"address" validate:"required,hostname_port"`
	LogLevel string `mapstructure:"log_level"`
	Database string `mapstructure:"database" validate:"required"`
	Secret   string `mapstructure:"secret" validate:"required,uuid"`
}

func LoadServerConfig() (*ServerConfig, error) {
	loader := NewLoader()

	loader.vp.SetDefault("address", "localhost:8080")
	loader.vp.SetDefault("log_level", "debug")
	loader.vp.SetDefault("database", "")
	loader.vp.SetDefault("secret", "")

	fs := pflag.NewFlagSet("", pflag.ContinueOnError)
	fs.StringP("address", "a", "", "The address to public metrics.")
	fs.StringP("log_level", "l", "info", "Log level")
	fs.StringP("database", "d", "", "Database")

	loader.EnableEnvParse()
	_ = loader.SetFlags(fs)

	cfg := &ServerConfig{
		Address:  loader.vp.GetString("address"),
		LogLevel: loader.vp.GetString("log_level"),
		Database: loader.vp.GetString("database"),
		Secret:   loader.vp.GetString("secret"),
	}

	if err := loader.Validation(cfg); err != nil {
		return nil, fmt.Errorf("failed to validate configuration: %w", err)
	}

	return cfg, nil
}
