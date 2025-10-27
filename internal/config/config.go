package config

import (
	"comptes/internal/domain"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
)

// Config represents the configuration structure
type Config struct {
	Accounts   []domain.Account  `yaml:"accounts"`
	Categories []domain.Category `yaml:"categories"`
	Tags       []domain.Tag      `yaml:"tags"`
}

// LoadConfig loads configuration from YAML file
func LoadConfig(configPath string) (*Config, error) {
	// Check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file not found: %s", configPath)
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Set default values for accounts
	for i := range config.Accounts {
		if config.Accounts[i].ID == "" {
			return nil, fmt.Errorf("account ID is required")
		}
		if config.Accounts[i].Currency == "" {
			config.Accounts[i].Currency = "EUR"
		}
		if config.Accounts[i].CreatedAt.IsZero() {
			config.Accounts[i].CreatedAt = time.Now()
		}
		config.Accounts[i].IsActive = true
		// Note: initial_balance is preserved from YAML
	}

	return &config, nil
}

// GetConfigPath returns the path to the config file
func GetConfigPath() string {
	execPath, _ := os.Executable()
	return filepath.Join(filepath.Dir(execPath), "config", "config.yaml")
}
