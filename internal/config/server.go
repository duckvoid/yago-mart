package config

import (
	"fmt"

	"github.com/spf13/pflag"
)

type ServerConfig struct {
	Address        string `mapstructure:"run_address" validate:"required,hostname_port"`
	LogLevel       string `mapstructure:"log_level"`
	Database       string `mapstructure:"database_uri" validate:"required"`
	Secret         string `mapstructure:"secret" validate:"required,uuid"`
	AccrualAddress string `mapstructure:"accrual_system_address" validate:"required"`
}

func LoadServerConfig() (*ServerConfig, error) {
	loader := NewLoader()

	loader.vp.SetDefault("run_address", "localhost:8080")
	loader.vp.SetDefault("log_level", "debug")
	loader.vp.SetDefault("database_uri", "")
	loader.vp.SetDefault("secret", "")
	loader.vp.SetDefault("accrual_system_address", "")

	fs := pflag.NewFlagSet("", pflag.ContinueOnError)
	fs.StringP("run_address", "a", "localhost:8080", "The address to public metrics.")
	fs.StringP("log_level", "l", "info", "Log level")
	fs.StringP("database_uri", "d", "", "Database")
	fs.StringP("accrual_system_address", "r", "", "Accrual System Address")

	loader.EnableEnvParse()
	_ = loader.SetFlags(fs)

	cfg := &ServerConfig{
		Address:        loader.vp.GetString("run_address"),
		LogLevel:       loader.vp.GetString("log_level"),
		Database:       loader.vp.GetString("database_uri"),
		Secret:         loader.vp.GetString("secret"),
		AccrualAddress: loader.vp.GetString("accrual_system_address"),
	}

	if err := loader.Validation(cfg); err != nil {
		return nil, fmt.Errorf("failed to validate configuration: %w", err)
	}

	return cfg, nil
}
