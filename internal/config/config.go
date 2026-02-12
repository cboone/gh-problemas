package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Config holds the application configuration.
type Config struct {
	Version  int      `mapstructure:"version"`
	Defaults Defaults `mapstructure:"defaults"`
	Theme    string   `mapstructure:"theme"`
}

// Defaults holds default configuration values.
type Defaults struct {
	Repo            string `mapstructure:"repo"`
	RefreshInterval int    `mapstructure:"refresh_interval"`
	PageSize        int    `mapstructure:"page_size"`
	DateFormat      string `mapstructure:"date_format"`
}

// Load reads configuration from the config file with sensible defaults.
func Load() (*Config, error) {
	v := viper.New()

	// Set defaults
	v.SetDefault("version", 1)
	v.SetDefault("defaults.repo", "")
	v.SetDefault("defaults.refresh_interval", 300)
	v.SetDefault("defaults.page_size", 50)
	v.SetDefault("defaults.date_format", "relative")
	v.SetDefault("theme", "dark")

	// Config path
	configDir := configDirectory()
	v.SetConfigName("config")
	v.AddConfigPath(configDir)

	// Read config file; ignore not-found
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("reading config: %w", err)
		}
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("parsing config: %w", err)
	}

	return &cfg, nil
}

func configDirectory() string {
	if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
		return filepath.Join(xdg, "gh-problemas")
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return filepath.Join(".", ".config", "gh-problemas")
	}
	return filepath.Join(home, ".config", "gh-problemas")
}
