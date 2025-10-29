package service

import (
	"comptes/internal/domain"
	"fmt"
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

func (m *MockStorage) GetAccountBalance(accountID string) (float64, error) {
	// Find the account
	var initialBalance float64
	var accountFound bool
	for _, acc := range m.accounts {
		if acc.ID == accountID {
			initialBalance = acc.InitialBalance
			accountFound = true
			break
		}
	}

	if !accountFound {
		return 0, fmt.Errorf("account not found: %s", accountID)
	}

	// Calculate balance
	balance := initialBalance
	for _, txn := range m.transactions {
		if txn.Account == accountID && txn.IsActive {
			balance += txn.Amount
		}
	}

	return balance, nil
}

func (m *MockStorage) GetPendingBatches() ([]domain.TransactionBatch, error) {
	return []domain.TransactionBatch{}, nil
}

func (m *MockStorage) SavePendingBatches(batches []domain.TransactionBatch) error {
	return nil
}

func (m *MockStorage) GetCommittedBatches() ([]domain.TransactionBatch, error) {
	return []domain.TransactionBatch{}, nil
}

func (m *MockStorage) SaveCommittedBatches(batches []domain.TransactionBatch) error {
	return nil
}

func (m *MockStorage) GetRolledBackBatches() ([]domain.TransactionBatch, error) {
	return []domain.TransactionBatch{}, nil
}

func (m *MockStorage) SaveRolledBackBatches(batches []domain.TransactionBatch) error {
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
		Account:     "account1",
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
				Account:     "account1",
				Amount:      -50.0,
				Description: "Purchase",
				IsActive:    true,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				ID:          "txn2",
				Account:     "account1",
				Amount:      100.0,
				Description: "Income",
				IsActive:    true,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				ID:          "txn3",
				Account:     "account1",
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
		Account:     "nonexistent",
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

func TestTransactionService_EditTransaction(t *testing.T) {
	// Setup
	now := time.Now()
	mockStorage := &MockStorage{
		transactions: []domain.Transaction{
			{
				ID:          "txn1",
				Account:     "account1",
				Amount:      -50.0,
				Description: "Original purchase",
				Categories:  []string{"food"},
				Tags:        []string{"test"},
				IsActive:    true,
				CreatedAt:   now,
				UpdatedAt:   now,
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
				CreatedAt:      now,
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

	// Test editing a transaction
	modifications := domain.Transaction{
		ID:          "txn2",             // New ID
		Amount:      -75.0,              // Changed amount
		Description: "Updated purchase", // Changed description
	}

	newTransaction, err := service.EditTransaction("txn1", modifications, "Price correction")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify new transaction was created
	if newTransaction.ID != "txn2" {
		t.Errorf("Expected new transaction ID 'txn2', got '%s'", newTransaction.ID)
	}

	if newTransaction.Amount != -75.0 {
		t.Errorf("Expected amount -75.0, got %.2f", newTransaction.Amount)
	}

	if newTransaction.Description != "Updated purchase" {
		t.Errorf("Expected description 'Updated purchase', got '%s'", newTransaction.Description)
	}

	if newTransaction.ParentID != "txn1" {
		t.Errorf("Expected parent ID 'txn1', got '%s'", newTransaction.ParentID)
	}

	// Verify old transaction was soft-deleted
	transactions, _ := mockStorage.GetTransactions()
	if len(transactions) != 2 {
		t.Errorf("Expected 2 transactions, got %d", len(transactions))
	}

	// Find the old transaction
	var oldTransaction *domain.Transaction
	for _, txn := range transactions {
		if txn.ID == "txn1" {
			oldTransaction = &txn
			break
		}
	}

	if oldTransaction == nil {
		t.Error("Old transaction not found")
	} else {
		if oldTransaction.IsActive {
			t.Error("Expected old transaction to be inactive")
		}
		if oldTransaction.EditComment != "Price correction" {
			t.Errorf("Expected edit comment 'Price correction', got '%s'", oldTransaction.EditComment)
		}
	}
}

func TestTransactionService_DeleteTransaction(t *testing.T) {
	// Setup
	now := time.Now()
	mockStorage := &MockStorage{
		transactions: []domain.Transaction{
			{
				ID:          "txn1",
				Account:     "account1",
				Amount:      -50.0,
				Description: "Test purchase",
				IsActive:    true,
				CreatedAt:   now,
				UpdatedAt:   now,
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
				CreatedAt:      now,
			},
		},
		categories: []domain.Category{},
		tags:       []domain.Tag{},
	}

	service := NewTransactionService(mockStorage)

	// Test deleting a transaction
	err := service.DeleteTransaction("txn1", "Mistake")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify transaction was soft-deleted
	transactions, _ := mockStorage.GetTransactions()
	if len(transactions) != 1 {
		t.Errorf("Expected 1 transaction, got %d", len(transactions))
	}

	txn := transactions[0]
	if txn.IsActive {
		t.Error("Expected transaction to be inactive")
	}
	if txn.EditComment != "Mistake" {
		t.Errorf("Expected edit comment 'Mistake', got '%s'", txn.EditComment)
	}
}

func TestTransactionService_DeleteTransaction_AlreadyDeleted(t *testing.T) {
	// Setup
	now := time.Now()
	mockStorage := &MockStorage{
		transactions: []domain.Transaction{
			{
				ID:          "txn1",
				Account:     "account1",
				Amount:      -50.0,
				Description: "Test purchase",
				IsActive:    false, // Already deleted
				CreatedAt:   now,
				UpdatedAt:   now,
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
				CreatedAt:      now,
			},
		},
		categories: []domain.Category{},
		tags:       []domain.Tag{},
	}

	service := NewTransactionService(mockStorage)

	// Test deleting an already deleted transaction
	err := service.DeleteTransaction("txn1", "Mistake")
	if err == nil {
		t.Error("Expected error for already deleted transaction, got nil")
	}
}

func TestTransactionService_UndoTransaction_UndoAdd(t *testing.T) {
	// Setup
	now := time.Now()
	mockStorage := &MockStorage{
		transactions: []domain.Transaction{
			{
				ID:          "txn1",
				Account:     "account1",
				Amount:      -50.0,
				Description: "Test purchase",
				IsActive:    true,
				EditComment: "", // No edit comment = added transaction
				CreatedAt:   now,
				UpdatedAt:   now,
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
				CreatedAt:      now,
			},
		},
		categories: []domain.Category{},
		tags:       []domain.Tag{},
	}

	service := NewTransactionService(mockStorage)

	// Test undoing an add operation
	err := service.UndoTransaction("txn1")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify transaction was soft-deleted
	transactions, _ := mockStorage.GetTransactions()
	if len(transactions) != 1 {
		t.Errorf("Expected 1 transaction, got %d", len(transactions))
	}

	txn := transactions[0]
	if txn.IsActive {
		t.Error("Expected transaction to be inactive")
	}
	if txn.EditComment != "Undo add operation" {
		t.Errorf("Expected edit comment 'Undo add operation', got '%s'", txn.EditComment)
	}
}

func TestTransactionService_UndoTransaction_UndoDelete(t *testing.T) {
	// Setup
	now := time.Now()
	mockStorage := &MockStorage{
		transactions: []domain.Transaction{
			{
				ID:          "txn1",
				Account:     "account1",
				Amount:      -50.0,
				Description: "Test purchase",
				IsActive:    false, // Deleted transaction
				EditComment: "Mistake",
				CreatedAt:   now,
				UpdatedAt:   now,
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
				CreatedAt:      now,
			},
		},
		categories: []domain.Category{},
		tags:       []domain.Tag{},
	}

	service := NewTransactionService(mockStorage)

	// Test undoing a delete operation
	err := service.UndoTransaction("txn1")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify transaction was restored
	transactions, _ := mockStorage.GetTransactions()
	if len(transactions) != 1 {
		t.Errorf("Expected 1 transaction, got %d", len(transactions))
	}

	txn := transactions[0]
	if !txn.IsActive {
		t.Error("Expected transaction to be active")
	}
	if txn.EditComment != "" {
		t.Errorf("Expected empty edit comment, got '%s'", txn.EditComment)
	}
}

func TestTransactionService_UndoTransaction_UndoEdit(t *testing.T) {
	// Setup
	now := time.Now()
	mockStorage := &MockStorage{
		transactions: []domain.Transaction{
			{
				ID:          "txn1",
				Account:     "account1",
				Amount:      -50.0,
				Description: "Original purchase",
				IsActive:    false, // Parent transaction (soft-deleted)
				EditComment: "Price correction",
				CreatedAt:   now,
				UpdatedAt:   now,
			},
			{
				ID:          "txn2",
				Account:     "account1",
				Amount:      -75.0,
				Description: "Updated purchase",
				IsActive:    true,
				ParentID:    "txn1", // Child transaction
				CreatedAt:   now,
				UpdatedAt:   now,
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
				CreatedAt:      now,
			},
		},
		categories: []domain.Category{},
		tags:       []domain.Tag{},
	}

	service := NewTransactionService(mockStorage)

	// Test undoing an edit operation
	err := service.UndoTransaction("txn2")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify parent was restored and child was removed
	transactions, _ := mockStorage.GetTransactions()
	if len(transactions) != 1 {
		t.Errorf("Expected 1 transaction, got %d", len(transactions))
	}

	txn := transactions[0]
	if txn.ID != "txn1" {
		t.Errorf("Expected transaction ID 'txn1', got '%s'", txn.ID)
	}
	if !txn.IsActive {
		t.Error("Expected transaction to be active")
	}
	if txn.EditComment != "" {
		t.Errorf("Expected empty edit comment, got '%s'", txn.EditComment)
	}
}

func TestTransactionService_findTransactionByID(t *testing.T) {
	// Setup
	now := time.Now()
	mockStorage := &MockStorage{
		transactions: []domain.Transaction{
			{
				ID:          "abc12345",
				Account:     "account1",
				Amount:      -50.0,
				Description: "Transaction 1",
				IsActive:    true,
				CreatedAt:   now,
				UpdatedAt:   now,
			},
			{
				ID:          "abdef67890",
				Account:     "account1",
				Amount:      -25.0,
				Description: "Transaction 2",
				IsActive:    true,
				CreatedAt:   now,
				UpdatedAt:   now,
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
				CreatedAt:      now,
			},
		},
		categories: []domain.Category{},
		tags:       []domain.Tag{},
	}

	service := NewTransactionService(mockStorage)

	// Test finding transaction by full ID
	transactions, _ := mockStorage.GetTransactions()
	txn, err := service.findTransactionByID(transactions, "abc12345")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if txn.ID != "abc12345" {
		t.Errorf("Expected transaction ID 'abc12345', got '%s'", txn.ID)
	}

	// Test finding transaction by partial ID
	txn, err = service.findTransactionByID(transactions, "abc")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if txn.ID != "abc12345" {
		t.Errorf("Expected transaction ID 'abc12345', got '%s'", txn.ID)
	}

	// Test finding non-existent transaction
	_, err = service.findTransactionByID(transactions, "xyz")
	if err == nil {
		t.Error("Expected error for non-existent transaction, got nil")
	}

	// Test ambiguous partial ID
	_, err = service.findTransactionByID(transactions, "ab")
	if err == nil {
		t.Error("Expected error for ambiguous partial ID, got nil")
	}
}

func TestTransactionService_DeleteTransactionHard(t *testing.T) {
	now := time.Now()
	mockStorage := &MockStorage{
		transactions: []domain.Transaction{
			{
				ID:          "test123",
				Account:     "account1",
				Amount:      -25.0,
				Description: "Test transaction",
				IsActive:    true,
				CreatedAt:   now,
				UpdatedAt:   now,
			},
			{
				ID:          "test456",
				Account:     "account2",
				Amount:      100.0,
				Description: "Another transaction",
				IsActive:    true,
				CreatedAt:   now,
				UpdatedAt:   now,
			},
		},
	}

	service := NewTransactionService(mockStorage)

	// Test hard delete
	err := service.DeleteTransactionHard("test123", "Permanent deletion")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify transaction was permanently removed
	transactions := mockStorage.transactions
	if len(transactions) != 1 {
		t.Errorf("Expected 1 transaction remaining, got %d", len(transactions))
	}
	if transactions[0].ID != "test456" {
		t.Errorf("Expected remaining transaction ID 'test456', got '%s'", transactions[0].ID)
	}

	// Test hard delete non-existent transaction
	err = service.DeleteTransactionHard("nonexistent", "Test")
	if err == nil {
		t.Error("Expected error for non-existent transaction, got nil")
	}
}

func TestTransactionService_UndoTransactionHard(t *testing.T) {
	now := time.Now()
	mockStorage := &MockStorage{
		transactions: []domain.Transaction{
			{
				ID:          "test123",
				Account:     "account1",
				Amount:      -25.0,
				Description: "Test transaction",
				IsActive:    true,
				CreatedAt:   now,
				UpdatedAt:   now,
			},
			{
				ID:          "test456",
				Account:     "account2",
				Amount:      100.0,
				Description: "Another transaction",
				IsActive:    true,
				CreatedAt:   now,
				UpdatedAt:   now,
			},
		},
	}

	service := NewTransactionService(mockStorage)

	// Test hard undo
	err := service.UndoTransactionHard("test123")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify transaction was permanently removed
	transactions := mockStorage.transactions
	if len(transactions) != 1 {
		t.Errorf("Expected 1 transaction remaining, got %d", len(transactions))
	}
	if transactions[0].ID != "test456" {
		t.Errorf("Expected remaining transaction ID 'test456', got '%s'", transactions[0].ID)
	}

	// Test hard undo non-existent transaction
	err = service.UndoTransactionHard("nonexistent")
	if err == nil {
		t.Error("Expected error for non-existent transaction, got nil")
	}
}
