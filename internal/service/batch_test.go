package service

import (
	"comptes/internal/domain"
	"comptes/internal/errors"
	"testing"
	"time"
)

// MockStorageForBatch extends MockStorage with batch support
type MockStorageForBatch struct {
	*MockStorage
	pendingBatches    []domain.TransactionBatch
	committedBatches  []domain.TransactionBatch
	rolledBackBatches []domain.TransactionBatch
}

func NewMockStorageForBatch() *MockStorageForBatch {
	return &MockStorageForBatch{
		MockStorage: &MockStorage{
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
					Code:        "ALM",
					Name:        "Alimentation",
					Description: "Food expenses",
				},
			},
			tags: []domain.Tag{
				{
					Code:        "REC",
					Name:        "Récurrent",
					Description: "Recurring transaction",
				},
			},
		},
		pendingBatches:    []domain.TransactionBatch{},
		committedBatches:  []domain.TransactionBatch{},
		rolledBackBatches: []domain.TransactionBatch{},
	}
}

func (m *MockStorageForBatch) GetPendingBatches() ([]domain.TransactionBatch, error) {
	return m.pendingBatches, nil
}

func (m *MockStorageForBatch) SavePendingBatches(batches []domain.TransactionBatch) error {
	m.pendingBatches = batches
	return nil
}

func (m *MockStorageForBatch) GetCommittedBatches() ([]domain.TransactionBatch, error) {
	return m.committedBatches, nil
}

func (m *MockStorageForBatch) SaveCommittedBatches(batches []domain.TransactionBatch) error {
	m.committedBatches = batches
	return nil
}

func (m *MockStorageForBatch) GetRolledBackBatches() ([]domain.TransactionBatch, error) {
	return m.rolledBackBatches, nil
}

func (m *MockStorageForBatch) SaveRolledBackBatches(batches []domain.TransactionBatch) error {
	m.rolledBackBatches = batches
	return nil
}

func TestBatchService_BeginTransaction(t *testing.T) {
	mockStorage := NewMockStorageForBatch()
	transactionService := NewTransactionService(mockStorage.MockStorage)
	batchService := NewTransactionBatchService(mockStorage, transactionService)

	// Test creating a batch without description
	batch, err := batchService.BeginTransaction("")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if batch == nil {
		t.Fatal("Expected batch to be created, got nil")
	}

	if batch.ID == "" {
		t.Error("Expected batch to have an ID")
	}

	if batch.Description != "" {
		t.Errorf("Expected empty description, got '%s'", batch.Description)
	}

	if batch.CommittedAt != nil {
		t.Error("Expected CommittedAt to be nil")
	}

	if batch.RolledBackAt != nil {
		t.Error("Expected RolledBackAt to be nil")
	}

	if len(batch.Transactions) != 0 {
		t.Errorf("Expected empty transactions list, got %d", len(batch.Transactions))
	}

	// Verify batch was saved
	pendingBatches, _ := mockStorage.GetPendingBatches()
	if len(pendingBatches) != 1 {
		t.Errorf("Expected 1 pending batch, got %d", len(pendingBatches))
	}

	if pendingBatches[0].ID != batch.ID {
		t.Errorf("Expected batch ID '%s', got '%s'", batch.ID, pendingBatches[0].ID)
	}
}

func TestBatchService_BeginTransaction_WithDescription(t *testing.T) {
	mockStorage := NewMockStorageForBatch()
	transactionService := NewTransactionService(mockStorage.MockStorage)
	batchService := NewTransactionBatchService(mockStorage, transactionService)

	description := "Test batch description"
	batch, err := batchService.BeginTransaction(description)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if batch.Description != description {
		t.Errorf("Expected description '%s', got '%s'", description, batch.Description)
	}
}

func TestBatchService_BeginTransaction_MultipleBatches(t *testing.T) {
	mockStorage := NewMockStorageForBatch()
	transactionService := NewTransactionService(mockStorage.MockStorage)
	batchService := NewTransactionBatchService(mockStorage, transactionService)

	// Create multiple batches
	batch1, _ := batchService.BeginTransaction("Batch 1")
	batch2, _ := batchService.BeginTransaction("Batch 2")
	batch3, _ := batchService.BeginTransaction("Batch 3")

	// Verify all batches are saved
	pendingBatches, _ := mockStorage.GetPendingBatches()
	if len(pendingBatches) != 3 {
		t.Errorf("Expected 3 pending batches, got %d", len(pendingBatches))
	}

	// Verify batch IDs are unique
	batchIDs := make(map[string]bool)
	for _, batch := range []*domain.TransactionBatch{batch1, batch2, batch3} {
		if batchIDs[batch.ID] {
			t.Errorf("Duplicate batch ID: %s", batch.ID)
		}
		batchIDs[batch.ID] = true
	}
}

func TestBatchService_GetPendingBatches(t *testing.T) {
	mockStorage := NewMockStorageForBatch()
	transactionService := NewTransactionService(mockStorage.MockStorage)
	batchService := NewTransactionBatchService(mockStorage, transactionService)

	// Create some batches
	batchService.BeginTransaction("Batch 1")
	batchService.BeginTransaction("Batch 2")

	// Get pending batches
	batches, err := batchService.GetPendingBatches()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(batches) != 2 {
		t.Errorf("Expected 2 pending batches, got %d", len(batches))
	}
}

func TestBatchService_GetPendingBatches_Empty(t *testing.T) {
	mockStorage := NewMockStorageForBatch()
	transactionService := NewTransactionService(mockStorage.MockStorage)
	batchService := NewTransactionBatchService(mockStorage, transactionService)

	batches, err := batchService.GetPendingBatches()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(batches) != 0 {
		t.Errorf("Expected 0 pending batches, got %d", len(batches))
	}
}

func TestBatchService_GetPendingBatchByID_FullID(t *testing.T) {
	mockStorage := NewMockStorageForBatch()
	transactionService := NewTransactionService(mockStorage.MockStorage)
	batchService := NewTransactionBatchService(mockStorage, transactionService)

	// Create a batch
	createdBatch, _ := batchService.BeginTransaction("Test batch")

	// Find by full ID
	foundBatch, err := batchService.GetPendingBatchByID(createdBatch.ID)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if foundBatch.ID != createdBatch.ID {
		t.Errorf("Expected batch ID '%s', got '%s'", createdBatch.ID, foundBatch.ID)
	}
}

func TestBatchService_GetPendingBatchByID_PartialID(t *testing.T) {
	mockStorage := NewMockStorageForBatch()
	transactionService := NewTransactionService(mockStorage.MockStorage)
	batchService := NewTransactionBatchService(mockStorage, transactionService)

	// Create a batch
	createdBatch, _ := batchService.BeginTransaction("Test batch")

	// Find by partial ID (first 8 characters)
	partialID := createdBatch.ID[:8]
	foundBatch, err := batchService.GetPendingBatchByID(partialID)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if foundBatch.ID != createdBatch.ID {
		t.Errorf("Expected batch ID '%s', got '%s'", createdBatch.ID, foundBatch.ID)
	}
}

func TestBatchService_GetPendingBatchByID_Ambiguous(t *testing.T) {
	mockStorage := NewMockStorageForBatch()
	transactionService := NewTransactionService(mockStorage.MockStorage)
	batchService := NewTransactionBatchService(mockStorage, transactionService)

	// Create two batches with similar IDs
	batch1, _ := batchService.BeginTransaction("Batch 1")
	batchService.BeginTransaction("Batch 2")

	// Try to find with ambiguous partial ID (first character only)
	partialID := batch1.ID[:1]
	_, err := batchService.GetPendingBatchByID(partialID)
	if err == nil {
		t.Error("Expected error for ambiguous ID, got nil")
	}

	// Verify it's an AmbiguousID error
	var comptesErr *errors.ComptesError
	if err != nil {
		var ok bool
		comptesErr, ok = err.(*errors.ComptesError)
		if !ok {
			// Try unwrapping
			comptesErr, ok = err.(*errors.ComptesError)
		}
		if ok && comptesErr.Code != errors.CodeAmbiguousID {
			t.Errorf("Expected CodeAmbiguousID, got %s", comptesErr.Code)
		} else if !ok {
			// Just verify error message contains expected text
			if err.Error() == "" {
				t.Error("Expected error message, got empty")
			}
		}
	}
}

func TestBatchService_GetPendingBatchByID_NotFound(t *testing.T) {
	mockStorage := NewMockStorageForBatch()
	transactionService := NewTransactionService(mockStorage.MockStorage)
	batchService := NewTransactionBatchService(mockStorage, transactionService)

	// Try to find non-existent batch
	_, err := batchService.GetPendingBatchByID("nonexistent")
	if err == nil {
		t.Error("Expected error for non-existent batch, got nil")
	}

	// Verify it's a TransactionNotFound error
	var comptesErr *errors.ComptesError
	if err != nil {
		var ok bool
		comptesErr, ok = err.(*errors.ComptesError)
		if ok && comptesErr.Code != errors.CodeTransactionNotFound {
			t.Errorf("Expected CodeTransactionNotFound, got %s", comptesErr.Code)
		} else if !ok {
			// Just verify error message contains expected text
			if err.Error() == "" {
				t.Error("Expected error message, got empty")
			}
		}
	}
}

func TestBatchService_AddTransactionToBatch(t *testing.T) {
	mockStorage := NewMockStorageForBatch()
	transactionService := NewTransactionService(mockStorage.MockStorage)
	batchService := NewTransactionBatchService(mockStorage, transactionService)

	// Create a batch
	batch, _ := batchService.BeginTransaction("Test batch")

	// Create a transaction
	transaction := domain.Transaction{
		Account:     "account1",
		Amount:      -50.0,
		Description: "Test transaction",
		Categories:  []string{"ALM"},
		Tags:        []string{"REC"},
		IsActive:    true,
	}

	// Add transaction to batch
	err := batchService.AddTransactionToBatch(batch.ID, transaction)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify transaction was added
	foundBatch, _ := batchService.GetPendingBatchByID(batch.ID)
	if len(foundBatch.Transactions) != 1 {
		t.Errorf("Expected 1 transaction in batch, got %d", len(foundBatch.Transactions))
	}

	if foundBatch.Transactions[0].Description != transaction.Description {
		t.Errorf("Expected description '%s', got '%s'", transaction.Description, foundBatch.Transactions[0].Description)
	}

	// Verify transaction has ID and timestamps
	if foundBatch.Transactions[0].ID == "" {
		t.Error("Expected transaction to have an ID")
	}

	if foundBatch.Transactions[0].CreatedAt.IsZero() {
		t.Error("Expected transaction to have CreatedAt timestamp")
	}

	if foundBatch.Transactions[0].Date.IsZero() {
		t.Error("Expected transaction to have Date timestamp")
	}
}

func TestBatchService_AddTransactionToBatch_WithExistingID(t *testing.T) {
	mockStorage := NewMockStorageForBatch()
	transactionService := NewTransactionService(mockStorage.MockStorage)
	batchService := NewTransactionBatchService(mockStorage, transactionService)

	// Create a batch
	batch, _ := batchService.BeginTransaction("Test batch")

	// Create a transaction with existing ID
	existingID := "existing-transaction-id"
	transaction := domain.Transaction{
		ID:          existingID,
		Account:     "account1",
		Amount:      -50.0,
		Description: "Test transaction",
		IsActive:    true,
	}

	// Add transaction to batch
	err := batchService.AddTransactionToBatch(batch.ID, transaction)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify transaction ID is preserved
	foundBatch, _ := batchService.GetPendingBatchByID(batch.ID)
	if foundBatch.Transactions[0].ID != existingID {
		t.Errorf("Expected transaction ID '%s', got '%s'", existingID, foundBatch.Transactions[0].ID)
	}
}

func TestBatchService_AddTransactionToBatch_MultipleTransactions(t *testing.T) {
	mockStorage := NewMockStorageForBatch()
	transactionService := NewTransactionService(mockStorage.MockStorage)
	batchService := NewTransactionBatchService(mockStorage, transactionService)

	// Create a batch
	batch, _ := batchService.BeginTransaction("Test batch")

	// Add multiple transactions
	transaction1 := domain.Transaction{
		Account:     "account1",
		Amount:      -50.0,
		Description: "Transaction 1",
		IsActive:    true,
	}
	transaction2 := domain.Transaction{
		Account:     "account1",
		Amount:      -25.0,
		Description: "Transaction 2",
		IsActive:    true,
	}

	batchService.AddTransactionToBatch(batch.ID, transaction1)
	batchService.AddTransactionToBatch(batch.ID, transaction2)

	// Verify both transactions were added
	foundBatch, _ := batchService.GetPendingBatchByID(batch.ID)
	if len(foundBatch.Transactions) != 2 {
		t.Errorf("Expected 2 transactions in batch, got %d", len(foundBatch.Transactions))
	}
}

func TestBatchService_AddTransactionToBatch_NotFound(t *testing.T) {
	mockStorage := NewMockStorageForBatch()
	transactionService := NewTransactionService(mockStorage.MockStorage)
	batchService := NewTransactionBatchService(mockStorage, transactionService)

	transaction := domain.Transaction{
		Account:     "account1",
		Amount:      -50.0,
		Description: "Test transaction",
		IsActive:    true,
	}

	// Try to add to non-existent batch
	err := batchService.AddTransactionToBatch("nonexistent", transaction)
	if err == nil {
		t.Error("Expected error for non-existent batch, got nil")
	}
}

func TestBatchService_AddTransactionToBatch_PartialID(t *testing.T) {
	mockStorage := NewMockStorageForBatch()
	transactionService := NewTransactionService(mockStorage.MockStorage)
	batchService := NewTransactionBatchService(mockStorage, transactionService)

	// Create a batch
	batch, _ := batchService.BeginTransaction("Test batch")

	transaction := domain.Transaction{
		Account:     "account1",
		Amount:      -50.0,
		Description: "Test transaction",
		IsActive:    true,
	}

	// Add using partial ID
	partialID := batch.ID[:8]
	err := batchService.AddTransactionToBatch(partialID, transaction)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify transaction was added
	foundBatch, _ := batchService.GetPendingBatchByID(batch.ID)
	if len(foundBatch.Transactions) != 1 {
		t.Errorf("Expected 1 transaction in batch, got %d", len(foundBatch.Transactions))
	}
}

func TestBatchService_CommitBatch(t *testing.T) {
	mockStorage := NewMockStorageForBatch()
	transactionService := NewTransactionService(mockStorage.MockStorage)
	batchService := NewTransactionBatchService(mockStorage, transactionService)

	// Create a batch
	batch, _ := batchService.BeginTransaction("Test batch")

	// Commit empty batch
	err := batchService.CommitBatch(batch.ID)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify batch was moved to committed
	pendingBatches, _ := mockStorage.GetPendingBatches()
	if len(pendingBatches) != 0 {
		t.Errorf("Expected 0 pending batches, got %d", len(pendingBatches))
	}

	committedBatches, _ := mockStorage.GetCommittedBatches()
	if len(committedBatches) != 1 {
		t.Errorf("Expected 1 committed batch, got %d", len(committedBatches))
	}

	if committedBatches[0].ID != batch.ID {
		t.Errorf("Expected batch ID '%s', got '%s'", batch.ID, committedBatches[0].ID)
	}

	// Verify batch has CommittedAt timestamp
	if committedBatches[0].CommittedAt == nil {
		t.Error("Expected batch to have CommittedAt timestamp")
	}
}

func TestBatchService_CommitBatch_WithTransactions(t *testing.T) {
	mockStorage := NewMockStorageForBatch()
	transactionService := NewTransactionService(mockStorage.MockStorage)
	batchService := NewTransactionBatchService(mockStorage, transactionService)

	// Create a batch
	batch, _ := batchService.BeginTransaction("Test batch")

	// Add transactions
	transaction1 := domain.Transaction{
		Account:     "account1",
		Amount:      -50.0,
		Description: "Transaction 1",
		Categories:  []string{"ALM"},
		IsActive:    true,
	}
	transaction2 := domain.Transaction{
		Account:     "account1",
		Amount:      -25.0,
		Description: "Transaction 2",
		Tags:        []string{"REC"},
		IsActive:    true,
	}

	batchService.AddTransactionToBatch(batch.ID, transaction1)
	batchService.AddTransactionToBatch(batch.ID, transaction2)

	// Commit batch
	err := batchService.CommitBatch(batch.ID)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify transactions were added to main transactions
	mainTransactions, _ := mockStorage.GetTransactions()
	if len(mainTransactions) != 2 {
		t.Errorf("Expected 2 transactions in main storage, got %d", len(mainTransactions))
	}

	// Verify batch was moved to committed
	committedBatches, _ := mockStorage.GetCommittedBatches()
	if len(committedBatches) != 1 {
		t.Errorf("Expected 1 committed batch, got %d", len(committedBatches))
	}

	if len(committedBatches[0].Transactions) != 2 {
		t.Errorf("Expected 2 transactions in committed batch, got %d", len(committedBatches[0].Transactions))
	}
}

func TestBatchService_CommitBatch_ValidationFailure(t *testing.T) {
	mockStorage := NewMockStorageForBatch()
	transactionService := NewTransactionService(mockStorage.MockStorage)
	batchService := NewTransactionBatchService(mockStorage, transactionService)

	// Create a batch
	batch, _ := batchService.BeginTransaction("Test batch")

	// Add invalid transaction (non-existent account)
	invalidTransaction := domain.Transaction{
		Account:     "nonexistent",
		Amount:      -50.0,
		Description: "Invalid transaction",
		IsActive:    true,
	}

	batchService.AddTransactionToBatch(batch.ID, invalidTransaction)

	// Try to commit batch - should fail validation
	err := batchService.CommitBatch(batch.ID)
	if err == nil {
		t.Error("Expected error for invalid transaction, got nil")
	}

	// Verify batch was NOT moved to committed
	pendingBatches, _ := mockStorage.GetPendingBatches()
	if len(pendingBatches) != 1 {
		t.Errorf("Expected 1 pending batch, got %d", len(pendingBatches))
	}

	committedBatches, _ := mockStorage.GetCommittedBatches()
	if len(committedBatches) != 0 {
		t.Errorf("Expected 0 committed batches, got %d", len(committedBatches))
	}

	// Verify transactions were NOT added to main storage
	mainTransactions, _ := mockStorage.GetTransactions()
	if len(mainTransactions) != 0 {
		t.Errorf("Expected 0 transactions in main storage, got %d", len(mainTransactions))
	}
}

func TestBatchService_CommitBatch_NotFound(t *testing.T) {
	mockStorage := NewMockStorageForBatch()
	transactionService := NewTransactionService(mockStorage.MockStorage)
	batchService := NewTransactionBatchService(mockStorage, transactionService)

	// Try to commit non-existent batch
	err := batchService.CommitBatch("nonexistent")
	if err == nil {
		t.Error("Expected error for non-existent batch, got nil")
	}
}

func TestBatchService_CommitBatch_PartialID(t *testing.T) {
	mockStorage := NewMockStorageForBatch()
	transactionService := NewTransactionService(mockStorage.MockStorage)
	batchService := NewTransactionBatchService(mockStorage, transactionService)

	// Create a batch
	batch, _ := batchService.BeginTransaction("Test batch")

	// Commit using partial ID
	partialID := batch.ID[:8]
	err := batchService.CommitBatch(partialID)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify batch was committed
	committedBatches, _ := mockStorage.GetCommittedBatches()
	if len(committedBatches) != 1 {
		t.Errorf("Expected 1 committed batch, got %d", len(committedBatches))
	}
}

func TestBatchService_CommitBatch_MultipleBatches(t *testing.T) {
	mockStorage := NewMockStorageForBatch()
	transactionService := NewTransactionService(mockStorage.MockStorage)
	batchService := NewTransactionBatchService(mockStorage, transactionService)

	// Create multiple batches
	batch1, _ := batchService.BeginTransaction("Batch 1")
	batch2, _ := batchService.BeginTransaction("Batch 2")
	batch3, _ := batchService.BeginTransaction("Batch 3")

	// Commit only batch2
	err := batchService.CommitBatch(batch2.ID)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify only batch2 was committed
	pendingBatches, _ := mockStorage.GetPendingBatches()
	if len(pendingBatches) != 2 {
		t.Errorf("Expected 2 pending batches, got %d", len(pendingBatches))
	}

	committedBatches, _ := mockStorage.GetCommittedBatches()
	if len(committedBatches) != 1 {
		t.Errorf("Expected 1 committed batch, got %d", len(committedBatches))
	}

	if committedBatches[0].ID != batch2.ID {
		t.Errorf("Expected committed batch ID '%s', got '%s'", batch2.ID, committedBatches[0].ID)
	}

	// Verify batch1 and batch3 are still pending
	pendingIDs := make(map[string]bool)
	for _, b := range pendingBatches {
		pendingIDs[b.ID] = true
	}

	if !pendingIDs[batch1.ID] {
		t.Error("Expected batch1 to still be pending")
	}

	if !pendingIDs[batch3.ID] {
		t.Error("Expected batch3 to still be pending")
	}
}

func TestBatchService_RollbackBatch(t *testing.T) {
	mockStorage := NewMockStorageForBatch()
	transactionService := NewTransactionService(mockStorage.MockStorage)
	batchService := NewTransactionBatchService(mockStorage, transactionService)

	// Create a batch
	batch, _ := batchService.BeginTransaction("Test batch")

	// Add transactions
	transaction := domain.Transaction{
		Account:     "account1",
		Amount:      -50.0,
		Description: "Test transaction",
		IsActive:    true,
	}

	batchService.AddTransactionToBatch(batch.ID, transaction)

	// Rollback batch
	err := batchService.RollbackBatch(batch.ID)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify batch was moved to rolled back
	pendingBatches, _ := mockStorage.GetPendingBatches()
	if len(pendingBatches) != 0 {
		t.Errorf("Expected 0 pending batches, got %d", len(pendingBatches))
	}

	rolledBackBatches, _ := mockStorage.GetRolledBackBatches()
	if len(rolledBackBatches) != 1 {
		t.Errorf("Expected 1 rolled back batch, got %d", len(rolledBackBatches))
	}

	if rolledBackBatches[0].ID != batch.ID {
		t.Errorf("Expected batch ID '%s', got '%s'", batch.ID, rolledBackBatches[0].ID)
	}

	// Verify batch has RolledBackAt timestamp
	if rolledBackBatches[0].RolledBackAt == nil {
		t.Error("Expected batch to have RolledBackAt timestamp")
	}

	// Verify transactions were NOT added to main storage
	mainTransactions, _ := mockStorage.GetTransactions()
	if len(mainTransactions) != 0 {
		t.Errorf("Expected 0 transactions in main storage, got %d", len(mainTransactions))
	}
}

func TestBatchService_RollbackBatch_NotFound(t *testing.T) {
	mockStorage := NewMockStorageForBatch()
	transactionService := NewTransactionService(mockStorage.MockStorage)
	batchService := NewTransactionBatchService(mockStorage, transactionService)

	// Try to rollback non-existent batch
	err := batchService.RollbackBatch("nonexistent")
	if err == nil {
		t.Error("Expected error for non-existent batch, got nil")
	}
}

func TestBatchService_RollbackBatch_PartialID(t *testing.T) {
	mockStorage := NewMockStorageForBatch()
	transactionService := NewTransactionService(mockStorage.MockStorage)
	batchService := NewTransactionBatchService(mockStorage, transactionService)

	// Create a batch
	batch, _ := batchService.BeginTransaction("Test batch")

	// Rollback using partial ID
	partialID := batch.ID[:8]
	err := batchService.RollbackBatch(partialID)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify batch was rolled back
	rolledBackBatches, _ := mockStorage.GetRolledBackBatches()
	if len(rolledBackBatches) != 1 {
		t.Errorf("Expected 1 rolled back batch, got %d", len(rolledBackBatches))
	}
}

func TestBatchService_GetCommittedBatches(t *testing.T) {
	mockStorage := NewMockStorageForBatch()
	transactionService := NewTransactionService(mockStorage.MockStorage)
	batchService := NewTransactionBatchService(mockStorage, transactionService)

	// Create and commit batches
	batch1, _ := batchService.BeginTransaction("Batch 1")
	batch2, _ := batchService.BeginTransaction("Batch 2")

	batchService.CommitBatch(batch1.ID)
	batchService.CommitBatch(batch2.ID)

	// Get committed batches
	committedBatches, err := batchService.GetCommittedBatches()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(committedBatches) != 2 {
		t.Errorf("Expected 2 committed batches, got %d", len(committedBatches))
	}
}

func TestBatchService_GetRolledBackBatches(t *testing.T) {
	mockStorage := NewMockStorageForBatch()
	transactionService := NewTransactionService(mockStorage.MockStorage)
	batchService := NewTransactionBatchService(mockStorage, transactionService)

	// Create and rollback batches
	batch1, _ := batchService.BeginTransaction("Batch 1")
	batch2, _ := batchService.BeginTransaction("Batch 2")

	batchService.RollbackBatch(batch1.ID)
	batchService.RollbackBatch(batch2.ID)

	// Get rolled back batches
	rolledBackBatches, err := batchService.GetRolledBackBatches()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(rolledBackBatches) != 2 {
		t.Errorf("Expected 2 rolled back batches, got %d", len(rolledBackBatches))
	}
}

func TestBatchService_CommitBatch_Atomicity(t *testing.T) {
	mockStorage := NewMockStorageForBatch()
	transactionService := NewTransactionService(mockStorage.MockStorage)
	batchService := NewTransactionBatchService(mockStorage, transactionService)

	// Create a batch
	batch, _ := batchService.BeginTransaction("Test batch")

	// Add valid and invalid transactions
	validTransaction := domain.Transaction{
		Account:     "account1",
		Amount:      -50.0,
		Description: "Valid transaction",
		Categories:  []string{"ALM"},
		IsActive:    true,
	}

	invalidTransaction := domain.Transaction{
		Account:     "nonexistent",
		Amount:      -25.0,
		Description: "Invalid transaction",
		IsActive:    true,
	}

	batchService.AddTransactionToBatch(batch.ID, validTransaction)
	batchService.AddTransactionToBatch(batch.ID, invalidTransaction)

	// Try to commit - should fail entirely
	err := batchService.CommitBatch(batch.ID)
	if err == nil {
		t.Error("Expected error for invalid transaction, got nil")
	}

	// Verify NO transactions were added to main storage (atomicity)
	mainTransactions, _ := mockStorage.GetTransactions()
	if len(mainTransactions) != 0 {
		t.Errorf("Expected 0 transactions in main storage (atomicity), got %d", len(mainTransactions))
	}

	// Verify batch is still pending
	pendingBatches, _ := mockStorage.GetPendingBatches()
	if len(pendingBatches) != 1 {
		t.Errorf("Expected 1 pending batch, got %d", len(pendingBatches))
	}
}
