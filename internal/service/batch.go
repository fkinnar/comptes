package service

import (
	"comptes/internal/domain"
	"comptes/internal/errors"
	"comptes/internal/storage"
	"strings"
	"time"

	"github.com/google/uuid"
)

// TransactionBatchService handles transaction batch operations
type TransactionBatchService struct {
	storage            storage.Storage
	transactionService *TransactionService
}

// NewTransactionBatchService creates a new transaction batch service
func NewTransactionBatchService(storage storage.Storage, transactionService *TransactionService) *TransactionBatchService {
	return &TransactionBatchService{
		storage:            storage,
		transactionService: transactionService,
	}
}

// BeginTransaction creates a new pending transaction batch
func (s *TransactionBatchService) BeginTransaction(description string) (*domain.TransactionBatch, error) {
	batch := domain.TransactionBatch{
		ID:           uuid.New().String(),
		Description:  description,
		CreatedAt:    time.Now(),
		Transactions: []domain.Transaction{},
	}

	// Get existing pending batches
	batches, err := s.storage.GetPendingBatches()
	if err != nil {
		return nil, errors.StorageReadFailed("pending_transactions", err)
	}

	// Add new batch
	batches = append(batches, batch)

	// Save pending batches
	if err := s.storage.SavePendingBatches(batches); err != nil {
		return nil, errors.StorageWriteFailed("pending_transactions", err)
	}

	return &batch, nil
}

// GetPendingBatches returns all pending transaction batches
func (s *TransactionBatchService) GetPendingBatches() ([]domain.TransactionBatch, error) {
	return s.storage.GetPendingBatches()
}

// GetPendingBatchByID finds a pending batch by ID (supports partial IDs)
func (s *TransactionBatchService) GetPendingBatchByID(id string) (*domain.TransactionBatch, error) {
	batches, err := s.storage.GetPendingBatches()
	if err != nil {
		return nil, errors.StorageReadFailed("pending_transactions", err)
	}

	var matches []domain.TransactionBatch
	for _, batch := range batches {
		if strings.HasPrefix(batch.ID, id) {
			matches = append(matches, batch)
		}
	}

	if len(matches) == 0 {
		return nil, errors.TransactionNotFound(id)
	}

	if len(matches) > 1 {
		return nil, errors.AmbiguousID(id)
	}

	return &matches[0], nil
}

// AddTransactionToBatch adds a transaction to a pending batch
func (s *TransactionBatchService) AddTransactionToBatch(batchID string, transaction domain.Transaction) error {
	// Get pending batches
	batches, err := s.storage.GetPendingBatches()
	if err != nil {
		return errors.StorageReadFailed("pending_transactions", err)
	}

	// Find the batch
	var batchFound bool
	for i, batch := range batches {
		if strings.HasPrefix(batch.ID, batchID) {
			// Generate ID if not provided
			if transaction.ID == "" {
				transaction.ID = uuid.New().String()
			}

			// Set timestamps if not provided
			if transaction.CreatedAt.IsZero() {
				transaction.CreatedAt = time.Now()
			}
			if transaction.UpdatedAt.IsZero() {
				transaction.UpdatedAt = time.Now()
			}
			if transaction.Date.IsZero() {
				transaction.Date = time.Now()
			}

			transaction.IsActive = true

			// Add transaction to batch
			batches[i].Transactions = append(batches[i].Transactions, transaction)
			batchFound = true
			break
		}
	}

	if !batchFound {
		return errors.TransactionNotFound(batchID)
	}

	// Save pending batches
	if err := s.storage.SavePendingBatches(batches); err != nil {
		return errors.StorageWriteFailed("pending_transactions", err)
	}

	return nil
}

// CommitBatch commits a pending batch by adding all transactions to the main transactions file
func (s *TransactionBatchService) CommitBatch(batchID string) error {
	// Get pending batches
	batches, err := s.storage.GetPendingBatches()
	if err != nil {
		return errors.StorageReadFailed("pending_transactions", err)
	}

	// Find the batch
	var batch *domain.TransactionBatch
	var batchIndex int
	var batchFound bool
	for i, b := range batches {
		if strings.HasPrefix(b.ID, batchID) {
			batch = &b
			batchIndex = i
			batchFound = true
			break
		}
	}

	if !batchFound {
		return errors.TransactionNotFound(batchID)
	}

	// Validate all transactions in the batch
	for _, transaction := range batch.Transactions {
		if err := s.transactionService.ValidateTransaction(transaction); err != nil {
			return errors.Wrap(errors.ErrorTypeValidation, "validation_failed", "Transaction validation failed in batch", err)
		}
	}

	// Get existing transactions
	existingTransactions, err := s.storage.GetTransactions()
	if err != nil {
		return errors.StorageReadFailed("transactions", err)
	}

	// Add all transactions from batch to main transactions
	for _, transaction := range batch.Transactions {
		existingTransactions = append(existingTransactions, transaction)
	}

	// Save transactions
	if err := s.storage.SaveTransactions(existingTransactions); err != nil {
		return errors.StorageWriteFailed("transactions", err)
	}

	// Mark batch as committed
	now := time.Now()
	batch.CommittedAt = &now

	// Remove from pending batches
	var newPendingBatches []domain.TransactionBatch
	for i, b := range batches {
		if i != batchIndex {
			newPendingBatches = append(newPendingBatches, b)
		}
	}

	if err := s.storage.SavePendingBatches(newPendingBatches); err != nil {
		return errors.StorageWriteFailed("pending_transactions", err)
	}

	// Add to committed batches
	committedBatches, err := s.storage.GetCommittedBatches()
	if err != nil {
		return errors.StorageReadFailed("committed_transactions", err)
	}

	committedBatches = append(committedBatches, *batch)

	if err := s.storage.SaveCommittedBatches(committedBatches); err != nil {
		return errors.StorageWriteFailed("committed_transactions", err)
	}

	return nil
}

// RollbackBatch rolls back a pending batch
func (s *TransactionBatchService) RollbackBatch(batchID string) error {
	// Get pending batches
	batches, err := s.storage.GetPendingBatches()
	if err != nil {
		return errors.StorageReadFailed("pending_transactions", err)
	}

	// Find the batch
	var batch *domain.TransactionBatch
	var batchIndex int
	var batchFound bool
	for i, b := range batches {
		if strings.HasPrefix(b.ID, batchID) {
			batch = &b
			batchIndex = i
			batchFound = true
			break
		}
	}

	if !batchFound {
		return errors.TransactionNotFound(batchID)
	}

	// Mark batch as rolled back
	now := time.Now()
	batch.RolledBackAt = &now

	// Remove from pending batches
	var newPendingBatches []domain.TransactionBatch
	for i, b := range batches {
		if i != batchIndex {
			newPendingBatches = append(newPendingBatches, b)
		}
	}

	if err := s.storage.SavePendingBatches(newPendingBatches); err != nil {
		return errors.StorageWriteFailed("pending_transactions", err)
	}

	// Add to rolled back batches
	rolledBackBatches, err := s.storage.GetRolledBackBatches()
	if err != nil {
		return errors.StorageReadFailed("rolled_back_transactions", err)
	}

	rolledBackBatches = append(rolledBackBatches, *batch)

	if err := s.storage.SaveRolledBackBatches(rolledBackBatches); err != nil {
		return errors.StorageWriteFailed("rolled_back_transactions", err)
	}

	return nil
}

// GetCommittedBatches returns all committed transaction batches
func (s *TransactionBatchService) GetCommittedBatches() ([]domain.TransactionBatch, error) {
	return s.storage.GetCommittedBatches()
}

// GetRolledBackBatches returns all rolled back transaction batches
func (s *TransactionBatchService) GetRolledBackBatches() ([]domain.TransactionBatch, error) {
	return s.storage.GetRolledBackBatches()
}
