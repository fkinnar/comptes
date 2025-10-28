package storage

import (
	"comptes/internal/domain"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// JSONStorage implements Storage interface using JSON files
type JSONStorage struct {
	dataDir string
}

// NewJSONStorage creates a new JSON storage instance
func NewJSONStorage(dataDir string) *JSONStorage {
	return &JSONStorage{
		dataDir: dataDir,
	}
}

// GetAccounts reads accounts from JSON file
func (s *JSONStorage) GetAccounts() ([]domain.Account, error) {
	var accounts []domain.Account
	return accounts, s.readJSONFile("accounts.json", &accounts)
}

// SaveAccounts saves accounts to JSON file
func (s *JSONStorage) SaveAccounts(accounts []domain.Account) error {
	return s.writeJSONFile("accounts.json", accounts)
}

// GetAccountBalance calculates the current balance for an account
func (s *JSONStorage) GetAccountBalance(accountID string) (float64, error) {
	// Get accounts to find initial balance
	accounts, err := s.GetAccounts()
	if err != nil {
		return 0, err
	}

	// Find the account
	var initialBalance float64
	var accountFound bool
	for _, acc := range accounts {
		if acc.ID == accountID {
			initialBalance = acc.InitialBalance
			accountFound = true
			break
		}
	}

	if !accountFound {
		return 0, fmt.Errorf("account not found: %s", accountID)
	}

	// Get transactions
	transactions, err := s.GetTransactions()
	if err != nil {
		return 0, err
	}

	// Calculate balance
	balance := initialBalance
	for _, txn := range transactions {
		if txn.Account == accountID && txn.IsActive {
			balance += txn.Amount
		}
	}

	return balance, nil
}

// GetTransactions reads transactions from JSON file
func (s *JSONStorage) GetTransactions() ([]domain.Transaction, error) {
	var transactions []domain.Transaction
	return transactions, s.readJSONFile("transactions.json", &transactions)
}

// SaveTransactions saves transactions to JSON file
func (s *JSONStorage) SaveTransactions(transactions []domain.Transaction) error {
	return s.writeJSONFile("transactions.json", transactions)
}

// GetCategories reads categories from JSON file
func (s *JSONStorage) GetCategories() ([]domain.Category, error) {
	var categories []domain.Category
	return categories, s.readJSONFile("categories.json", &categories)
}

// SaveCategories saves categories to JSON file
func (s *JSONStorage) SaveCategories(categories []domain.Category) error {
	return s.writeJSONFile("categories.json", categories)
}

// GetTags reads tags from JSON file
func (s *JSONStorage) GetTags() ([]domain.Tag, error) {
	var tags []domain.Tag
	return tags, s.readJSONFile("tags.json", &tags)
}

// SaveTags saves tags to JSON file
func (s *JSONStorage) SaveTags(tags []domain.Tag) error {
	return s.writeJSONFile("tags.json", tags)
}

// Helper methods

func (s *JSONStorage) readJSONFile(filename string, v interface{}) error {
	filepath := filepath.Join(s.dataDir, filename)

	// If file doesn't exist, return empty slice
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return nil
	}

	data, err := os.ReadFile(filepath)
	if err != nil {
		return fmt.Errorf("failed to read %s: %w", filename, err)
	}

	if len(data) == 0 {
		return nil
	}

	if err := json.Unmarshal(data, v); err != nil {
		return fmt.Errorf("failed to parse %s: %w", filename, err)
	}

	return nil
}

func (s *JSONStorage) writeJSONFile(filename string, v interface{}) error {
	filepath := filepath.Join(s.dataDir, filename)

	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal %s: %w", filename, err)
	}

	if err := os.WriteFile(filepath, data, 0644); err != nil {
		return fmt.Errorf("failed to write %s: %w", filename, err)
	}

	return nil
}
