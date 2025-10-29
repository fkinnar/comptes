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

// GetTransactions reads transactions (movements) from JSON file
func (s *JSONStorage) GetTransactions() ([]domain.Transaction, error) {
	var transactions []domain.Transaction
	// Try movements.json first (new format), fallback to transactions.json for migration
	movementsPath := filepath.Join(s.dataDir, "movements.json")
	transactionsPath := filepath.Join(s.dataDir, "transactions.json")

	// Check if movements.json exists
	if _, err := os.Stat(movementsPath); err == nil {
		return transactions, s.readJSONFile("movements.json", &transactions)
	}

	// Check if old transactions.json exists (migration)
	if _, err := os.Stat(transactionsPath); err == nil {
		err := s.readJSONFile("transactions.json", &transactions)
		if err != nil {
			return transactions, err
		}
		// Migrate: rename file to movements.json
		if err := os.Rename(transactionsPath, movementsPath); err != nil {
			return transactions, fmt.Errorf("failed to migrate transactions.json to movements.json: %w", err)
		}
		return transactions, nil
	}

	// Neither file exists, return empty
	return transactions, nil
}

// SaveTransactions saves transactions (movements) to JSON file
func (s *JSONStorage) SaveTransactions(transactions []domain.Transaction) error {
	return s.writeJSONFile("movements.json", transactions)
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

// GetPendingBatches reads pending transaction batches from JSON file
func (s *JSONStorage) GetPendingBatches() ([]domain.TransactionBatch, error) {
	var batches []domain.TransactionBatch
	return batches, s.readJSONFile("pending_transactions.json", &batches)
}

// SavePendingBatches saves pending transaction batches to JSON file
func (s *JSONStorage) SavePendingBatches(batches []domain.TransactionBatch) error {
	return s.writeJSONFile("pending_transactions.json", batches)
}

// GetCommittedBatches reads committed transaction batches from JSON file
func (s *JSONStorage) GetCommittedBatches() ([]domain.TransactionBatch, error) {
	var batches []domain.TransactionBatch
	return batches, s.readJSONFile("committed_transactions.json", &batches)
}

// SaveCommittedBatches saves committed transaction batches to JSON file
func (s *JSONStorage) SaveCommittedBatches(batches []domain.TransactionBatch) error {
	return s.writeJSONFile("committed_transactions.json", batches)
}

// GetRolledBackBatches reads rolled back transaction batches from JSON file
func (s *JSONStorage) GetRolledBackBatches() ([]domain.TransactionBatch, error) {
	var batches []domain.TransactionBatch
	return batches, s.readJSONFile("rolled_back_transactions.json", &batches)
}

// SaveRolledBackBatches saves rolled back transaction batches to JSON file
func (s *JSONStorage) SaveRolledBackBatches(batches []domain.TransactionBatch) error {
	return s.writeJSONFile("rolled_back_transactions.json", batches)
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
