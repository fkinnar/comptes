package service

import (
	"comptes/internal/domain"
	"comptes/internal/storage"
	"fmt"
)

// TransactionService handles transaction operations
type TransactionService struct {
	storage storage.Storage
}

// NewTransactionService creates a new transaction service
func NewTransactionService(storage storage.Storage) *TransactionService {
	return &TransactionService{
		storage: storage,
	}
}

// AddTransaction adds a new transaction
func (s *TransactionService) AddTransaction(transaction domain.Transaction) error {
	// Validate transaction
	if err := s.validateTransaction(transaction); err != nil {
		return err
	}

	// Get existing transactions
	transactions, err := s.storage.GetTransactions()
	if err != nil {
		return fmt.Errorf("failed to get transactions: %w", err)
	}

	// Add new transaction
	transactions = append(transactions, transaction)

	// Save back to storage
	if err := s.storage.SaveTransactions(transactions); err != nil {
		return fmt.Errorf("failed to save transactions: %w", err)
	}

	return nil
}

// GetTransactions returns all transactions
func (s *TransactionService) GetTransactions() ([]domain.Transaction, error) {
	return s.storage.GetTransactions()
}

// GetAccountBalance calculates the current balance for an account
func (s *TransactionService) GetAccountBalance(accountID string) (float64, error) {
	// Verify account exists
	accounts, err := s.storage.GetAccounts()
	if err != nil {
		return 0, fmt.Errorf("failed to get accounts: %w", err)
	}

	// Find the account
	var accountFound bool
	for _, acc := range accounts {
		if acc.ID == accountID {
			accountFound = true
			break
		}
	}

	if !accountFound {
		return 0, fmt.Errorf("account not found: %s", accountID)
	}

	// Delegate balance calculation to storage
	return s.storage.GetAccountBalance(accountID)
}

// validateTransaction validates a transaction
func (s *TransactionService) validateTransaction(transaction domain.Transaction) error {
	// Check if account exists
	accounts, err := s.storage.GetAccounts()
	if err != nil {
		return fmt.Errorf("failed to get accounts: %w", err)
	}

	accountExists := false
	for _, account := range accounts {
		if account.ID == transaction.AccountID {
			accountExists = true
			break
		}
	}

	if !accountExists {
		return fmt.Errorf("account not found: %s", transaction.AccountID)
	}

	// Check if categories exist
	categories, err := s.storage.GetCategories()
	if err != nil {
		return fmt.Errorf("failed to get categories: %w", err)
	}

	categoryMap := make(map[string]bool)
	for _, category := range categories {
		categoryMap[category.Code] = true
	}

	for _, categoryCode := range transaction.Categories {
		if !categoryMap[categoryCode] {
			return fmt.Errorf("category not found: %s", categoryCode)
		}
	}

	// Check if tags exist
	tags, err := s.storage.GetTags()
	if err != nil {
		return fmt.Errorf("failed to get tags: %w", err)
	}

	tagMap := make(map[string]bool)
	for _, tag := range tags {
		tagMap[tag.Code] = true
	}

	for _, tagCode := range transaction.Tags {
		if !tagMap[tagCode] {
			return fmt.Errorf("tag not found: %s", tagCode)
		}
	}

	return nil
}
