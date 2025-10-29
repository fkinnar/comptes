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
			Account:     "account1",
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
			Account:     "account1",
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
	movementsFile := filepath.Join(tempDir, "movements.json")
	if _, err := os.Stat(movementsFile); os.IsNotExist(err) {
		t.Error("Expected movements.json file to be created")
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
	movementsFile := filepath.Join(tempDir, "movements.json")
	invalidJSON := `{"invalid": json, missing quotes}`
	err := os.WriteFile(movementsFile, []byte(invalidJSON), 0644)
	if err != nil {
		t.Fatalf("Failed to create invalid JSON file: %v", err)
	}

	// Test loading invalid JSON
	_, err = storage.GetTransactions()
	if err == nil {
		t.Error("Expected error loading invalid JSON, got nil")
	}
}

func TestJSONStorage_SaveAndLoadPendingBatches(t *testing.T) {
	// Create temporary directory
	tempDir := t.TempDir()
	storage := NewJSONStorage(tempDir)

	// Test data
	batches := []domain.TransactionBatch{
		{
			ID:          "batch1",
			Description: "Test batch 1",
			CreatedAt:   time.Now(),
			Transactions: []domain.Transaction{
				{
					ID:          "txn1",
					Account:     "account1",
					Amount:      -50.0,
					Description: "Transaction 1",
					IsActive:    true,
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				},
			},
		},
		{
			ID:           "batch2",
			Description:  "Test batch 2",
			CreatedAt:    time.Now(),
			Transactions: []domain.Transaction{},
		},
	}

	// Test saving batches
	err := storage.SavePendingBatches(batches)
	if err != nil {
		t.Errorf("Expected no error saving pending batches, got %v", err)
	}

	// Verify file was created
	pendingFile := filepath.Join(tempDir, "pending_transactions.json")
	if _, err := os.Stat(pendingFile); os.IsNotExist(err) {
		t.Error("Expected pending_transactions.json file to be created")
	}

	// Test loading batches
	loadedBatches, err := storage.GetPendingBatches()
	if err != nil {
		t.Errorf("Expected no error loading pending batches, got %v", err)
	}

	// Verify data integrity
	if len(loadedBatches) != 2 {
		t.Errorf("Expected 2 batches, got %d", len(loadedBatches))
	}

	if loadedBatches[0].ID != "batch1" {
		t.Errorf("Expected first batch ID 'batch1', got '%s'", loadedBatches[0].ID)
	}

	if len(loadedBatches[0].Transactions) != 1 {
		t.Errorf("Expected 1 transaction in first batch, got %d", len(loadedBatches[0].Transactions))
	}

	if loadedBatches[1].Description != "Test batch 2" {
		t.Errorf("Expected second batch description 'Test batch 2', got '%s'", loadedBatches[1].Description)
	}
}

func TestJSONStorage_SaveAndLoadCommittedBatches(t *testing.T) {
	// Create temporary directory
	tempDir := t.TempDir()
	storage := NewJSONStorage(tempDir)

	now := time.Now()
	batches := []domain.TransactionBatch{
		{
			ID:          "batch1",
			Description: "Committed batch",
			CreatedAt:   now.Add(-time.Hour),
			CommittedAt: &now,
			Transactions: []domain.Transaction{
				{
					ID:          "txn1",
					Account:     "account1",
					Amount:      -50.0,
					Description: "Transaction 1",
					IsActive:    true,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
			},
		},
	}

	// Test saving committed batches
	err := storage.SaveCommittedBatches(batches)
	if err != nil {
		t.Errorf("Expected no error saving committed batches, got %v", err)
	}

	// Verify file was created
	committedFile := filepath.Join(tempDir, "committed_transactions.json")
	if _, err := os.Stat(committedFile); os.IsNotExist(err) {
		t.Error("Expected committed_transactions.json file to be created")
	}

	// Test loading committed batches
	loadedBatches, err := storage.GetCommittedBatches()
	if err != nil {
		t.Errorf("Expected no error loading committed batches, got %v", err)
	}

	if len(loadedBatches) != 1 {
		t.Errorf("Expected 1 committed batch, got %d", len(loadedBatches))
	}

	if loadedBatches[0].ID != "batch1" {
		t.Errorf("Expected batch ID 'batch1', got '%s'", loadedBatches[0].ID)
	}

	if loadedBatches[0].CommittedAt == nil {
		t.Error("Expected CommittedAt timestamp to be set")
	}
}

func TestJSONStorage_SaveAndLoadRolledBackBatches(t *testing.T) {
	// Create temporary directory
	tempDir := t.TempDir()
	storage := NewJSONStorage(tempDir)

	now := time.Now()
	batches := []domain.TransactionBatch{
		{
			ID:           "batch1",
			Description:  "Rolled back batch",
			CreatedAt:    now.Add(-time.Hour),
			RolledBackAt: &now,
			Transactions: []domain.Transaction{
				{
					ID:          "txn1",
					Account:     "account1",
					Amount:      -50.0,
					Description: "Transaction 1",
					IsActive:    true,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
			},
		},
	}

	// Test saving rolled back batches
	err := storage.SaveRolledBackBatches(batches)
	if err != nil {
		t.Errorf("Expected no error saving rolled back batches, got %v", err)
	}

	// Verify file was created
	rolledBackFile := filepath.Join(tempDir, "rolled_back_transactions.json")
	if _, err := os.Stat(rolledBackFile); os.IsNotExist(err) {
		t.Error("Expected rolled_back_transactions.json file to be created")
	}

	// Test loading rolled back batches
	loadedBatches, err := storage.GetRolledBackBatches()
	if err != nil {
		t.Errorf("Expected no error loading rolled back batches, got %v", err)
	}

	if len(loadedBatches) != 1 {
		t.Errorf("Expected 1 rolled back batch, got %d", len(loadedBatches))
	}

	if loadedBatches[0].ID != "batch1" {
		t.Errorf("Expected batch ID 'batch1', got '%s'", loadedBatches[0].ID)
	}

	if loadedBatches[0].RolledBackAt == nil {
		t.Error("Expected RolledBackAt timestamp to be set")
	}
}

func TestJSONStorage_BatchFiles_Empty(t *testing.T) {
	// Create temporary directory
	tempDir := t.TempDir()
	storage := NewJSONStorage(tempDir)

	// Test loading empty batch files
	pendingBatches, err := storage.GetPendingBatches()
	if err != nil {
		t.Errorf("Expected no error loading empty pending batches, got %v", err)
	}

	if len(pendingBatches) != 0 {
		t.Errorf("Expected empty pending batches list, got %d batches", len(pendingBatches))
	}

	committedBatches, err := storage.GetCommittedBatches()
	if err != nil {
		t.Errorf("Expected no error loading empty committed batches, got %v", err)
	}

	if len(committedBatches) != 0 {
		t.Errorf("Expected empty committed batches list, got %d batches", len(committedBatches))
	}

	rolledBackBatches, err := storage.GetRolledBackBatches()
	if err != nil {
		t.Errorf("Expected no error loading empty rolled back batches, got %v", err)
	}

	if len(rolledBackBatches) != 0 {
		t.Errorf("Expected empty rolled back batches list, got %d batches", len(rolledBackBatches))
	}
}

func TestJSONStorage_BatchFiles_InvalidJSON(t *testing.T) {
	// Create temporary directory
	tempDir := t.TempDir()
	storage := NewJSONStorage(tempDir)

	// Create invalid JSON file for pending batches
	pendingFile := filepath.Join(tempDir, "pending_transactions.json")
	invalidJSON := `{"invalid": json, missing quotes}`
	err := os.WriteFile(pendingFile, []byte(invalidJSON), 0644)
	if err != nil {
		t.Fatalf("Failed to create invalid JSON file: %v", err)
	}

	// Test loading invalid JSON
	_, err = storage.GetPendingBatches()
	if err == nil {
		t.Error("Expected error loading invalid JSON, got nil")
	}

	// Test invalid JSON for committed batches
	committedFile := filepath.Join(tempDir, "committed_transactions.json")
	err = os.WriteFile(committedFile, []byte(invalidJSON), 0644)
	if err != nil {
		t.Fatalf("Failed to create invalid JSON file: %v", err)
	}

	_, err = storage.GetCommittedBatches()
	if err == nil {
		t.Error("Expected error loading invalid JSON, got nil")
	}

	// Test invalid JSON for rolled back batches
	rolledBackFile := filepath.Join(tempDir, "rolled_back_transactions.json")
	err = os.WriteFile(rolledBackFile, []byte(invalidJSON), 0644)
	if err != nil {
		t.Fatalf("Failed to create invalid JSON file: %v", err)
	}

	_, err = storage.GetRolledBackBatches()
	if err == nil {
		t.Error("Expected error loading invalid JSON, got nil")
	}
}
