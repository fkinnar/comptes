package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// getCurrentBatchID reads the current batch ID from the file
func (c *CLI) getCurrentBatchID() (string, error) {
	currentBatchFile := filepath.Join(c.dataDir, ".current_batch")
	data, err := os.ReadFile(currentBatchFile)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil // No current batch
		}
		return "", fmt.Errorf("failed to read current batch file: %w", err)
	}
	batchID := strings.TrimSpace(string(data))
	return batchID, nil
}

// saveCurrentBatchID saves the current batch ID to the file
func (c *CLI) saveCurrentBatchID(batchID string) error {
	currentBatchFile := filepath.Join(c.dataDir, ".current_batch")
	if batchID == "" {
		// Remove file if batch ID is empty
		if err := os.Remove(currentBatchFile); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("failed to remove current batch file: %w", err)
		}
		return nil
	}
	if err := os.WriteFile(currentBatchFile, []byte(batchID), 0644); err != nil {
		return fmt.Errorf("failed to save current batch file: %w", err)
	}
	return nil
}

// resolveBatchID resolves batch ID (uses provided ID or current batch ID)
func (c *CLI) resolveBatchID(providedID string) (string, error) {
	if providedID != "" {
		return providedID, nil
	}

	currentID, err := c.getCurrentBatchID()
	if err != nil {
		return "", err
	}

	if currentID == "" {
		return "", fmt.Errorf("no batch ID provided and no current batch set. Use 'comptes begin' first or provide a batch ID")
	}

	return currentID, nil
}
