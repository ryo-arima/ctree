package config

import (
	"fmt"
	"os"
	"path/filepath"
)

// Config represents the configuration for ctree
type Config struct {
	SourcePath   string   `yaml:"source_path"`
	OutputPath   string   `yaml:"output_path"`
	Languages    []string `yaml:"languages,omitempty"`
	Recursive    bool     `yaml:"recursive"`
	ExcludeFiles []string `yaml:"exclude_files,omitempty"`
	IncludeFiles []string `yaml:"include_files,omitempty"`
	MaxDepth     int      `yaml:"max_depth,omitempty"`
	ConfigPath   string   `yaml:"-"` // not serialized
}

// NewConfig creates a new default configuration
func NewConfig() *Config {
	return &Config{
		Recursive: true,
		MaxDepth:  10,
	}
}

// LoadConfigFromFile loads configuration from a TOML file
func LoadConfigFromFile(configPath string) (*Config, error) {
	config := NewConfig()
	config.ConfigPath = configPath

	// For now, return default config since we don't have TOML parsing yet
	// TODO: Add TOML parsing when needed
	return config, nil
}

// GetConfigPath returns the path to the configuration file
func GetConfigPath() string {
	// Check for config file in etc directory
	configPath := filepath.Join("etc", "app.toml")
	if _, err := os.Stat(configPath); err == nil {
		return configPath
	}

	// Check for config file in home directory
	homeDir, err := os.UserHomeDir()
	if err == nil {
		configPath = filepath.Join(homeDir, ".ctree", "config.toml")
		if _, err := os.Stat(configPath); err == nil {
			return configPath
		}
	}

	return ""
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.SourcePath == "" {
		return fmt.Errorf("source_path is required")
	}

	if c.OutputPath == "" {
		return fmt.Errorf("output_path is required")
	}

	return nil
}
