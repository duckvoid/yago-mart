package config

import (
	"fmt"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type Loader struct {
	vp *viper.Viper
	vl *validator.Validate
}

func NewLoader() *Loader {
	return &Loader{
		vp: viper.New(),
		vl: validator.New(),
	}
}

func (c *Loader) SetFlags(fs *pflag.FlagSet) error {
	if err := c.vp.BindPFlags(fs); err != nil {
		return fmt.Errorf("failed to bind flags: %w", err)
	}

	if err := fs.Parse(os.Args[1:]); err != nil {
		if errors.Is(err, pflag.ErrHelp) {
			os.Exit(0)
		}
		return fmt.Errorf("failed to parse flags: %w", err)
	}

	return nil
}

func (c *Loader) SetFileConfig(filename string, paths []string) error {
	c.vp.SetConfigName(filename)

	for _, path := range paths {
		c.vp.AddConfigPath(path)
	}

	if err := c.vp.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read config: %w", err)
	}

	return nil
}

func (c *Loader) EnableEnvParse() {
	c.vp.AutomaticEnv()
}

func (c *Loader) Validation(cfg interface{}) error {
	if err := c.vp.Unmarshal(cfg); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	if err := c.vl.Struct(cfg); err != nil {
		return fmt.Errorf("failed to validate config: %w", err)
	}

	return nil
}
