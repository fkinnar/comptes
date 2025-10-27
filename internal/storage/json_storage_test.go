package storage

import (
	"comptes/internal/domain"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestJSONStorage_SaveAndLoadTransactions(t *testing.T) {
	// Create temporary directory
	tempDir := t.TempDir()
	storage := NewJSONStorage(tempDir)
	
	// Test data
	transactions := []domain.Transaction{
		{
			ID:          "txn1",
			AccountID:   "account1",
			Amount:      -50.0,
			Description: "Test purchase",
			Categories:  []string{"food"},
			Tags:        []string{"test"},
			IsActive:    true,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          "txn2",
			AccountID:   "account1",
			Amount:      100.0,
			Description: "Test income",
			Categories:  []string{"salary"},
			Tags:        []string{"recurring"},
			IsActive:    true,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}
	
	// Test saving transactions
	err := storage.SaveTransactions(transactions)
	if err != nil {
		t.Errorf("Expected no error saving transactions, got %v", err)
	}
	
	// Verify file was created
	transactionsFile := filepath.Join(tempDir, "transactions.json")
	if _, err := os.Stat(transactionsFile); os.IsNotExist(err) {
		t.Error("Expected transactions.json file to be created")
	}
	
	// Test loading transactions
	loadedTransactions, err := storage.GetTransactions()
	if err != nil {
		t.Errorf("Expected no error loading transactions, got %v", err)
	}
	
	// Verify data integrity
	if len(loadedTransactions) != 2 {
		t.Errorf("Expected 2 transactions, got %d", len(loadedTransactions))
	}
	
	if loadedTransactions[0].ID != "txn1" {
		t.Errorf("Expected first transaction ID 'txn1', got '%s'", loadedTransactions[0].ID)
	}
	
	if loadedTransactions[1].Amount != 100.0 {
		t.Errorf("Expected second transaction amount 100.0, got %.2f", loadedTransactions[1].Amount)
	}
}

func TestJSONStorage_SaveAndLoadAccounts(t *testing.T) {
	// Create temporary directory
	tempDir := t.TempDir()
	storage := NewJSONStorage(tempDir)
	
	// Test data
	accounts := []domain.Account{
		{
			ID:             "account1",
			Name:           "Test Account",
			Type:           "checking",
			Currency:       "EUR",
			InitialBalance: 1500.0,
			IsActive:       true,
			CreatedAt:      time.Now(),
		},
		{
			ID:             "account2",
			Name:           "Savings Account",
			Type:           "savings",
			Currency:       "EUR",
			InitialBalance: 5000.0,
			IsActive:       true,
			CreatedAt:      time.Now(),
		},
	}
	
	// Test saving accounts
	err := storage.SaveAccounts(accounts)
	if err != nil {
		t.Errorf("Expected no error saving accounts, got %v", err)
	}
	
	// Verify file was created
	accountsFile := filepath.Join(tempDir, "accounts.json")
	if _, err := os.Stat(accountsFile); os.IsNotExist(err) {
		t.Error("Expected accounts.json file to be created")
	}
	
	// Test loading accounts
	loadedAccounts, err := storage.GetAccounts()
	if err != nil {
		t.Errorf("Expected no error loading accounts, got %v", err)
	}
	
	// Verify data integrity
	if len(loadedAccounts) != 2 {
		t.Errorf("Expected 2 accounts, got %d", len(loadedAccounts))
	}
	
	if loadedAccounts[0].Name != "Test Account" {
		t.Errorf("Expected first account name 'Test Account', got '%s'", loadedAccounts[0].Name)
	}
	
	if loadedAccounts[1].InitialBalance != 5000.0 {
		t.Errorf("Expected second account initial balance 5000.0, got %.2f", loadedAccounts[1].InitialBalance)
	}
}

func TestJSONStorage_LoadEmptyFiles(t *testing.T) {
	// Create temporary directory
	tempDir := t.TempDir()
	storage := NewJSONStorage(tempDir)
	
	// Test loading transactions from non-existent file
	transactions, err := storage.GetTransactions()
	if err != nil {
		t.Errorf("Expected no error loading non-existent transactions, got %v", err)
	}
	
	if len(transactions) != 0 {
		t.Errorf("Expected empty transactions list, got %d transactions", len(transactions))
	}
	
	// Test loading accounts from non-existent file
	accounts, err := storage.GetAccounts()
	if err != nil {
		t.Errorf("Expected no error loading non-existent accounts, got %v", err)
	}
	
	if len(accounts) != 0 {
		t.Errorf("Expected empty accounts list, got %d accounts", len(accounts))
	}
}

func TestJSONStorage_InvalidJSON(t *testing.T) {
	// Create temporary directory
	tempDir := t.TempDir()
	storage := NewJSONStorage(tempDir)
	
	// Create invalid JSON file
	transactionsFile := filepath.Join(tempDir, "transactions.json")
	invalidJSON := `{"invalid": json, missing quotes}`
	err := os.WriteFile(transactionsFile, []byte(invalidJSON), 0644)
	if err != nil {
		t.Fatalf("Failed to create invalid JSON file: %v", err)
	}
	
	// Test loading invalid JSON
	_, err = storage.GetTransactions()
	if err == nil {
		t.Error("Expected error loading invalid JSON, got nil")
	}
}
