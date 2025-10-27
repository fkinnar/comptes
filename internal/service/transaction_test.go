package service

import (
	"comptes/internal/domain"
	"testing"
	"time"
)

// MockStorage implements the Storage interface for testing
type MockStorage struct {
	transactions []domain.Transaction
	accounts     []domain.Account
	categories   []domain.Category
	tags         []domain.Tag
}

func (m *MockStorage) GetTransactions() ([]domain.Transaction, error) {
	return m.transactions, nil
}

func (m *MockStorage) SaveTransactions(transactions []domain.Transaction) error {
	m.transactions = transactions
	return nil
}

func (m *MockStorage) GetAccounts() ([]domain.Account, error) {
	return m.accounts, nil
}

func (m *MockStorage) SaveAccounts(accounts []domain.Account) error {
	m.accounts = accounts
	return nil
}

func (m *MockStorage) GetCategories() ([]domain.Category, error) {
	return m.categories, nil
}

func (m *MockStorage) SaveCategories(categories []domain.Category) error {
	m.categories = categories
	return nil
}

func (m *MockStorage) GetTags() ([]domain.Tag, error) {
	return m.tags, nil
}

func (m *MockStorage) SaveTags(tags []domain.Tag) error {
	m.tags = tags
	return nil
}

func TestTransactionService_AddTransaction(t *testing.T) {
	// Setup
	mockStorage := &MockStorage{
		transactions: []domain.Transaction{},
		accounts: []domain.Account{
			{
				ID:             "account1",
				Name:           "Test Account",
				Type:           "checking",
				Currency:       "EUR",
				InitialBalance: 1000.0,
				IsActive:       true,
				CreatedAt:      time.Now(),
			},
		},
		categories: []domain.Category{
			{
				Code:        "food",
				Name:        "Food",
				Description: "Food expenses",
			},
		},
		tags: []domain.Tag{
			{
				Code:        "test",
				Name:        "Test",
				Description: "Test tag",
			},
		},
	}
	
	service := NewTransactionService(mockStorage)
	
	// Test adding a valid transaction
	transaction := domain.Transaction{
		ID:          "txn1",
		AccountID:   "account1",
		Amount:      -50.0,
		Description: "Test purchase",
		Categories:  []string{"food"},
		Tags:        []string{"test"},
		IsActive:    true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	
	err := service.AddTransaction(transaction)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	
	// Verify transaction was saved
	transactions, _ := mockStorage.GetTransactions()
	if len(transactions) != 1 {
		t.Errorf("Expected 1 transaction, got %d", len(transactions))
	}
	
	if transactions[0].ID != "txn1" {
		t.Errorf("Expected transaction ID 'txn1', got '%s'", transactions[0].ID)
	}
}

func TestTransactionService_GetAccountBalance(t *testing.T) {
	// Setup
	mockStorage := &MockStorage{
		transactions: []domain.Transaction{
			{
				ID:          "txn1",
				AccountID:   "account1",
				Amount:      -50.0,
				Description: "Purchase",
				IsActive:    true,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				ID:          "txn2",
				AccountID:   "account1",
				Amount:      100.0,
				Description: "Income",
				IsActive:    true,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				ID:          "txn3",
				AccountID:   "account1",
				Amount:      -25.0,
				Description: "Inactive transaction",
				IsActive:    false, // This should not be counted
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
		},
		accounts: []domain.Account{
			{
				ID:             "account1",
				Name:           "Test Account",
				Type:           "checking",
				Currency:       "EUR",
				InitialBalance: 1000.0,
				IsActive:       true,
				CreatedAt:      time.Now(),
			},
		},
		categories: []domain.Category{},
		tags:       []domain.Tag{},
	}
	
	service := NewTransactionService(mockStorage)
	
	// Test balance calculation
	balance, err := service.GetAccountBalance("account1")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	
	// Expected balance: 1000 (initial) + 100 (income) - 50 (expense) = 1050
	// Inactive transaction (-25) should not be counted
	expectedBalance := 1050.0
	if balance != expectedBalance {
		t.Errorf("Expected balance %.2f, got %.2f", expectedBalance, balance)
	}
}

func TestTransactionService_GetAccountBalance_NonExistentAccount(t *testing.T) {
	// Setup
	mockStorage := &MockStorage{
		transactions: []domain.Transaction{},
		accounts:     []domain.Account{},
		categories:   []domain.Category{},
		tags:         []domain.Tag{},
	}
	
	service := NewTransactionService(mockStorage)
	
	// Test non-existent account
	_, err := service.GetAccountBalance("nonexistent")
	if err == nil {
		t.Error("Expected error for non-existent account, got nil")
	}
}

func TestTransactionService_AddTransaction_NonExistentAccount(t *testing.T) {
	// Setup
	mockStorage := &MockStorage{
		transactions: []domain.Transaction{},
		accounts:     []domain.Account{},
		categories:   []domain.Category{},
		tags:         []domain.Tag{},
	}
	
	service := NewTransactionService(mockStorage)
	
	// Test adding transaction to non-existent account
	transaction := domain.Transaction{
		ID:          "txn1",
		AccountID:   "nonexistent",
		Amount:      -50.0,
		Description: "Test purchase",
		IsActive:    true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	
	err := service.AddTransaction(transaction)
	if err == nil {
		t.Error("Expected error for non-existent account, got nil")
	}
}
