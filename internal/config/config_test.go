package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	// Create a temporary config file
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.yaml")

	configContent := `
accounts:
  - id: "test_account"
    name: "Test Account"
    type: "checking"
    currency: "EUR"
    initial_balance: 1500.0
    is_active: true

categories:
  - code: "FOOD"
    name: "Food"
    description: "Food expenses"
  - code: "TRANSPORT"
    name: "Transport"
    description: "Transport expenses"

tags:
  - code: "URGENT"
    name: "Urgent"
    description: "Urgent transaction"
`

	err := os.WriteFile(configPath, []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}

	// Test loading config
	config, err := LoadConfig(configPath)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify accounts
	if len(config.Accounts) != 1 {
		t.Errorf("Expected 1 account, got %d", len(config.Accounts))
	}

	account := config.Accounts[0]
	if account.ID != "test_account" {
		t.Errorf("Expected account ID 'test_account', got '%s'", account.ID)
	}

	if account.Name != "Test Account" {
		t.Errorf("Expected account name 'Test Account', got '%s'", account.Name)
	}

	if account.InitialBalance != 1500.0 {
		t.Errorf("Expected initial balance 1500.0, got %.2f", account.InitialBalance)
	}

	if account.Currency != "EUR" {
		t.Errorf("Expected currency 'EUR', got '%s'", account.Currency)
	}

	if !account.IsActive {
		t.Error("Expected account to be active")
	}

	// Verify categories
	if len(config.Categories) != 2 {
		t.Errorf("Expected 2 categories, got %d", len(config.Categories))
	}

	// Verify tags
	if len(config.Tags) != 1 {
		t.Errorf("Expected 1 tag, got %d", len(config.Tags))
	}
}

func TestLoadConfig_DefaultValues(t *testing.T) {
	// Create a minimal config file
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.yaml")

	configContent := `
accounts:
  - id: "minimal_account"
    name: "Minimal Account"

categories: []
tags: []
`

	err := os.WriteFile(configPath, []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}

	// Test loading config
	config, err := LoadConfig(configPath)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify default values
	account := config.Accounts[0]
	if account.Currency != "EUR" {
		t.Errorf("Expected default currency 'EUR', got '%s'", account.Currency)
	}

	if !account.IsActive {
		t.Error("Expected account to be active by default")
	}

	if account.CreatedAt.IsZero() {
		t.Error("Expected CreatedAt to be set")
	}
}

func TestLoadConfig_InvalidFile(t *testing.T) {
	// Test with non-existent file
	_, err := LoadConfig("/nonexistent/config.yaml")
	if err == nil {
		t.Error("Expected error for non-existent file, got nil")
	}
}

func TestLoadConfig_InvalidYAML(t *testing.T) {
	// Create a file with invalid YAML
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "invalid.yaml")

	invalidContent := `
accounts:
  - id: "test"
    name: "Test"
    invalid: [unclosed bracket
`

	err := os.WriteFile(configPath, []byte(invalidContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}

	// Test loading invalid config
	_, err = LoadConfig(configPath)
	if err == nil {
		t.Error("Expected error for invalid YAML, got nil")
	}
}

func TestLoadConfig_MissingAccountID(t *testing.T) {
	// Create a config with missing account ID
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.yaml")

	configContent := `
accounts:
  - name: "Account without ID"

categories: []
tags: []
`

	err := os.WriteFile(configPath, []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}

	// Test loading config with missing ID
	_, err = LoadConfig(configPath)
	if err == nil {
		t.Error("Expected error for missing account ID, got nil")
	}
}
