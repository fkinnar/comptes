package cli

import (
	"fmt"
)

func (c *CLI) handleMigrate() error {
	if err := c.migrateTransactionIDs(); err != nil {
		return fmt.Errorf("error migrating transaction IDs: %w", err)
	}
	fmt.Println("Transaction IDs migrated successfully!")
	return nil
}

func (c *CLI) migrateTransactionIDs() error {
	// Utiliser le service pour migrer
	if err := c.transactionService.MigrateTransactionIDs(c.generateShortID); err != nil {
		return err
	}

	return nil
}
