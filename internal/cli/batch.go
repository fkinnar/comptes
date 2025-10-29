package cli

import (
	"fmt"
	"strings"
)

func (c *CLI) handleBegin(args []string) error {
	// Optional description
	description := ""
	if len(args) > 2 {
		description = strings.Join(args[2:], " ")
	}

	batch, err := c.batchService.BeginTransaction(description)
	if err != nil {
		return fmt.Errorf("error beginning transaction: %w", err)
	}

	// Save as current batch
	if err := c.saveCurrentBatchID(batch.ID); err != nil {
		return fmt.Errorf("error saving current batch: %w", err)
	}

	fmt.Printf("Transaction batch started: %s\n", batch.ID)
	if description != "" {
		fmt.Printf("Description: %s\n", description)
	}
	fmt.Printf("You can now add transactions to this batch.\n")
	fmt.Printf("You can set context with 'comptes account <id>', 'comptes category <code>', 'comptes tags <code>'.\n")
	fmt.Printf("Use 'comptes commit' (or 'comptes commit %s') to commit or 'comptes rollback' to rollback.\n", batch.ID[:8])

	return nil
}

func (c *CLI) handleCommit(args []string) error {
	var providedBatchID string
	if len(args) >= 3 {
		providedBatchID = args[2]
	}

	batchID, err := c.resolveBatchID(providedBatchID)
	if err != nil {
		return fmt.Errorf("error resolving batch ID: %w", err)
	}

	// Get batch before committing to show transaction count
	batch, err := c.batchService.GetPendingBatchByID(batchID)
	if err != nil {
		return fmt.Errorf("error getting batch: %w", err)
	}

	transactionCount := len(batch.Transactions)

	if err := c.batchService.CommitBatch(batchID); err != nil {
		return fmt.Errorf("error committing transaction batch: %w", err)
	}

	// Clear current batch if it was the one committed
	currentID, _ := c.getCurrentBatchID()
	if currentID == batchID {
		c.saveCurrentBatchID("")
		// Clear context when batch is committed
		c.saveCurrentContext(&TransactionContext{})
	}

	fmt.Printf("Transaction batch %s committed successfully!\n", batchID)
	fmt.Printf("Committed %d transaction(s).\n", transactionCount)
	return nil
}

func (c *CLI) handleRollback(args []string) error {
	var providedBatchID string
	if len(args) >= 3 {
		providedBatchID = args[2]
	}

	batchID, err := c.resolveBatchID(providedBatchID)
	if err != nil {
		return fmt.Errorf("error resolving batch ID: %w", err)
	}

	if err := c.batchService.RollbackBatch(batchID); err != nil {
		return fmt.Errorf("error rolling back transaction batch: %w", err)
	}

	// Clear current batch if it was the one rolled back
	currentID, _ := c.getCurrentBatchID()
	if currentID == batchID {
		c.saveCurrentBatchID("")
		// Clear context when batch is rolled back
		c.saveCurrentContext(&TransactionContext{})
	}

	fmt.Printf("Transaction batch %s rolled back successfully!\n", batchID)

	return nil
}
