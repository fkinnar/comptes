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

// CreateDefaultConfig creates a default configuration
func CreateDefaultConfig() *Config {
	return &Config{
		Accounts: []domain.Account{
			{
				ID:             "BANQUE",
				Name:           "Compte Courant Principal",
				Type:           "checking",
				Currency:       "EUR",
				InitialBalance: 1500.00,
				IsActive:       true,
				CreatedAt:      time.Now(),
			},
		},
		Categories: []domain.Category{
			{Code: "ALM", Name: "Alimentation", Description: "Courses et repas"},
			{Code: "SLR", Name: "Salaire", Description: "Revenus professionnels"},
			{Code: "LGT", Name: "Logement", Description: "Loyer, charges, etc."},
		},
		Tags: []domain.Tag{
			{Code: "URG", Name: "Urgent", Description: "Transaction urgente"},
			{Code: "REC", Name: "Récurrent", Description: "Transaction récurrente"},
		},
	}
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

// SaveConfig saves configuration to YAML file
func SaveConfig(configPath string, cfg *Config) error {
	// Create directory if it doesn't exist
	dir := filepath.Dir(configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// GetConfigPath returns the path to the config file
func GetConfigPath() string {
	// Use environment variable for tests, or default to config/
	configDir := os.Getenv("COMPTES_CONFIG_DIR")
	if configDir == "" {
		// Default behavior: next to executable
		execPath, _ := os.Executable()
		configDir = filepath.Join(filepath.Dir(execPath), "config")
	}
	return filepath.Join(configDir, "config.yaml")
}
