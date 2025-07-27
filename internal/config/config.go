package config

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/spf13/pflag"
	"gopkg.in/yaml.v3"
)

type Config struct {
	RefreshInterval     time.Duration `yaml:"refresh_interval" env:"REFRESH_INTERVAL"`
	ColumnBordersColor  string        `yaml:"column_borders_color" env:"COLUMN_BORDERS_COLOR"`
	ShowColumnSeparator bool          `yaml:"show_column_separator" env:"SHOW_COLUMN_SEPARATOR"`
}

func Load(cfgPath string) (Config, error) {
	file, err := os.Open(cfgPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return Config{}, fmt.Errorf("file %s: not exist: %w", cfgPath, err)
		}
		return Config{}, fmt.Errorf("open file<%s>: %w", cfgPath, err)
	}

	defer file.Close()

	data, err := io.ReadAll(bufio.NewReader(file))

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return Config{}, fmt.Errorf("parse config file<%s>: %w", cfgPath, err)
	}

	return cfg, nil
}

const (
	CONFIG_PATH    = "CONFIG_PATH"
	defaultCfgPath = "configs/config.yaml"
)

func FetchConfigPath() string {
	var cfgPath string
	pflag.StringVarP(&cfgPath, "config-path", "c", defaultCfgPath, "path to file with config")
	pflag.Parse()

	return cfgPath
}
