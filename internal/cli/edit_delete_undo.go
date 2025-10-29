package cli

import (
	"comptes/internal/domain"
	"comptes/internal/errors"
	"encoding/json"
	"fmt"
)

func (c *CLI) handleEdit(args []string) error {
	if len(args) < 4 {
		ShowHelp("edit")
		return errors.MissingArguments("edit")
	}

	transactionID := args[2]
	jsonData := args[3]

	// Parse message flag (-m or --message)
	var message string
	for i, arg := range args {
		if (arg == "-m" || arg == "--message") && i+1 < len(args) {
			message = args[i+1]
			break
		}
	}

	if message == "" {
		ShowHelp("edit")
		return errors.MissingMessage("edit")
	}

	if err := c.editTransaction(transactionID, jsonData, message); err != nil {
		return errors.Wrap(errors.ErrorTypeUserInput, "edit_failed", "Failed to edit transaction", err)
	}
	fmt.Println("Transaction edited successfully!")
	return nil
}

func (c *CLI) handleDelete(args []string) error {
	if len(args) < 3 {
		ShowHelp("delete")
		return errors.MissingArguments("delete")
	}

	transactionID := args[2]

	// Check for help flag
	if transactionID == "--help" || transactionID == "-?" {
		ShowHelp("delete")
		return nil
	}

	// Parse flags
	var message string
	var hardDelete bool
	var force bool

	for i, arg := range args {
		if (arg == "-m" || arg == "--message") && i+1 < len(args) {
			message = args[i+1]
		}
		if arg == "--hard" || arg == "-H" {
			hardDelete = true
		}
		if arg == "-f" || arg == "--force" {
			force = true
		}
	}

	if message == "" {
		ShowHelp("delete")
		return errors.MissingMessage("delete")
	}

	// Confirmation for hard delete unless --force
	if hardDelete && !force {
		fmt.Printf("⚠️  WARNING: This will permanently delete transaction %s. This action cannot be undone!\n", transactionID)
		fmt.Print("Are you sure? (y/N): ")
		var response string
		fmt.Scanln(&response)
		if response != "y" && response != "Y" {
			fmt.Println("Operation cancelled.")
			return nil
		}
	}

	if hardDelete {
		if err := c.deleteTransactionHard(transactionID, message); err != nil {
			return errors.Wrap(errors.ErrorTypeUserInput, "delete_hard_failed", "Failed to permanently delete transaction", err)
		}
		fmt.Println("Transaction permanently deleted!")
	} else {
		if err := c.deleteTransaction(transactionID, message); err != nil {
			return errors.Wrap(errors.ErrorTypeUserInput, "delete_failed", "Failed to delete transaction", err)
		}
		fmt.Println("Transaction deleted successfully!")
	}
	return nil
}

func (c *CLI) handleUndo(args []string) error {
	if len(args) < 3 {
		ShowHelp("undo")
		return errors.MissingArguments("undo")
	}

	transactionID := args[2]

	// Check for help flag
	if transactionID == "--help" || transactionID == "-?" {
		ShowHelp("undo")
		return nil
	}

	// Parse flags
	var hardUndo bool
	var force bool

	for _, arg := range args {
		if arg == "--hard" || arg == "-H" {
			hardUndo = true
		}
		if arg == "-f" || arg == "--force" {
			force = true
		}
	}

	// Confirmation for hard undo unless --force
	if hardUndo && !force {
		fmt.Printf("⚠️  WARNING: This will permanently remove transaction %s. This action cannot be undone!\n", transactionID)
		fmt.Print("Are you sure? (y/N): ")
		var response string
		fmt.Scanln(&response)
		if response != "y" && response != "Y" {
			fmt.Println("Operation cancelled.")
			return nil
		}
	}

	if hardUndo {
		if err := c.undoTransactionHard(transactionID); err != nil {
			return errors.Wrap(errors.ErrorTypeUserInput, "undo_hard_failed", "Failed to permanently undo transaction", err)
		}
		fmt.Println("Transaction permanently removed!")
	} else {
		if err := c.undoTransaction(transactionID); err != nil {
			return errors.Wrap(errors.ErrorTypeUserInput, "undo_failed", "Failed to undo transaction", err)
		}
		fmt.Println("Transaction undone successfully!")
	}
	return nil
}

func (c *CLI) editTransaction(transactionID string, jsonData string, message string) error {
	// Parser le JSON des modifications
	var modifications TransactionInput
	if err := json.Unmarshal([]byte(jsonData), &modifications); err != nil {
		return errors.InvalidJSON(err)
	}

	// Convertir TransactionInput en domain.Transaction
	modificationsDomain := domain.Transaction{
		ID:          c.generateShortID(), // Nouvel ID généré
		Account:     modifications.Account,
		Amount:      modifications.Amount,
		Description: modifications.Description,
		Categories:  modifications.Categories,
		Tags:        modifications.Tags,
	}

	// Gérer les dates
	if !modifications.Date.IsZero() {
		modificationsDomain.Date = modifications.Date.Time
	}

	// Utiliser le service pour éditer
	newTransaction, err := c.transactionService.EditTransaction(transactionID, modificationsDomain, message)
	if err != nil {
		return err
	}

	fmt.Printf("Edited transaction %s -> %s\n", transactionID, newTransaction.ID)
	return nil
}

func (c *CLI) deleteTransaction(transactionID string, message string) error {
	// Utiliser le service pour supprimer
	if err := c.transactionService.DeleteTransaction(transactionID, message); err != nil {
		return err
	}

	fmt.Printf("Deleted transaction %s\n", transactionID)
	return nil
}

func (c *CLI) undoTransaction(transactionID string) error {
	// Utiliser le service pour annuler
	if err := c.transactionService.UndoTransaction(transactionID); err != nil {
		return err
	}

	return nil
}

func (c *CLI) deleteTransactionHard(transactionID string, message string) error {
	// Utiliser le service pour suppression définitive
	if err := c.transactionService.DeleteTransactionHard(transactionID, message); err != nil {
		return err
	}

	fmt.Printf("Permanently deleted transaction %s\n", transactionID)
	return nil
}

func (c *CLI) undoTransactionHard(transactionID string) error {
	// Utiliser le service pour undo définitif
	if err := c.transactionService.UndoTransactionHard(transactionID); err != nil {
		return err
	}

	fmt.Printf("Permanently removed transaction %s\n", transactionID)
	return nil
}
