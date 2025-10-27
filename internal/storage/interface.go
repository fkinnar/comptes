package storage

import (
	"comptes/internal/domain"
)

// Storage defines the interface for data persistence
type Storage interface {
	// Accounts
	GetAccounts() ([]domain.Account, error)
	SaveAccounts(accounts []domain.Account) error

	// Transactions
	GetTransactions() ([]domain.Transaction, error)
	SaveTransactions(transactions []domain.Transaction) error

	// Categories
	GetCategories() ([]domain.Category, error)
	SaveCategories(categories []domain.Category) error

	// Tags
	GetTags() ([]domain.Tag, error)
	SaveTags(tags []domain.Tag) error
}

