package cli

import (
	"comptes/internal/domain"
	"encoding/json"
	"fmt"
	"strings"
)

func (c *CLI) handleList(args []string) error {
	format := "text" // Format par défaut
	showHistory := false

	// Check for flags
	for i, arg := range args {
		if arg == "--format" && i+1 < len(args) {
			format = args[i+1]
		}
		if arg == "--history" {
			showHistory = true
		}
	}

	if err := c.listTransactions(format, showHistory); err != nil {
		return fmt.Errorf("error listing transactions: %w", err)
	}
	return nil
}

func (c *CLI) listTransactions(format string, showHistory bool) error {
	transactions, err := c.transactionService.GetTransactions()
	if err != nil {
		return err
	}

	if len(transactions) == 0 {
		fmt.Println("No transactions found.")
		return nil
	}

	// Filtrer selon le flag --history
	var filteredTransactions []domain.Transaction
	if showHistory {
		filteredTransactions = transactions // Toutes les transactions
	} else {
		// Seulement les transactions actives
		for _, txn := range transactions {
			if txn.IsActive {
				filteredTransactions = append(filteredTransactions, txn)
			}
		}
	}

	switch format {
	case "csv":
		return c.listTransactionsCSV(filteredTransactions, showHistory)
	case "json":
		return c.listTransactionsJSON(filteredTransactions, showHistory)
	case "text":
		fallthrough
	default:
		return c.listTransactionsText(filteredTransactions, showHistory)
	}
}

func (c *CLI) listTransactionsText(transactions []domain.Transaction, showHistory bool) error {
	if showHistory {
		fmt.Println("Transactions (all):")
	} else {
		fmt.Println("Transactions:")
	}

	for _, txn := range transactions {
		// Format de base avec ID
		var line string
		if showHistory {
			// Indicateur d'état seulement pour --history
			status := "✅"
			if !txn.IsActive {
				status = "❌"
			}
			line = fmt.Sprintf("- [%s] %s %s: %.2f EUR - %s (Categories: %v)",
				txn.ID, status, txn.Date.Format("2006-01-02"), txn.Amount, txn.Description, txn.Categories)
		} else {
			// Pas d'indicateur de statut pour list normal (toutes sont actives)
			line = fmt.Sprintf("- [%s] %s: %.2f EUR - %s (Categories: %v)",
				txn.ID, txn.Date.Format("2006-01-02"), txn.Amount, txn.Description, txn.Categories)
		}

		// Ajouter les tags seulement s'il y en a
		if len(txn.Tags) > 0 {
			line += fmt.Sprintf(", Tags: %v", txn.Tags)
		}

		// Ajouter le commentaire d'edit si présent
		if txn.EditComment != "" {
			line += fmt.Sprintf(" | Edit: %s", txn.EditComment)
		}

		fmt.Println(line)
	}
	return nil
}

func (c *CLI) listTransactionsCSV(transactions []domain.Transaction, showHistory bool) error {
	// En-tête CSV
	if showHistory {
		fmt.Println("id,date,amount,description,categories,tags,is_active,edit_comment")
	} else {
		fmt.Println("id,date,amount,description,categories,tags")
	}

	for _, txn := range transactions {
		// Convertir les catégories en string
		categoriesStr := ""
		if len(txn.Categories) > 0 {
			categoriesStr = strings.Join(txn.Categories, ";")
		}

		// Convertir les tags en string
		tagsStr := ""
		if len(txn.Tags) > 0 {
			tagsStr = strings.Join(txn.Tags, ";")
		}

		if showHistory {
			fmt.Printf("%s,%s,%.2f,%s,%s,%s,%t,%s\n",
				txn.ID, txn.Date.Format("2006-01-02"), txn.Amount, txn.Description, categoriesStr, tagsStr, txn.IsActive, txn.EditComment)
		} else {
			fmt.Printf("%s,%s,%.2f,%s,%s,%s\n",
				txn.ID, txn.Date.Format("2006-01-02"), txn.Amount, txn.Description, categoriesStr, tagsStr)
		}
	}
	return nil
}

func (c *CLI) listTransactionsJSON(transactions []domain.Transaction, showHistory bool) error {
	// Créer une structure simplifiée pour le JSON
	type TransactionOutput struct {
		ID          string   `json:"id"`
		Date        string   `json:"date"`
		Amount      float64  `json:"amount"`
		Description string   `json:"description"`
		Categories  []string `json:"categories"`
		Tags        []string `json:"tags"`
		IsActive    *bool    `json:"is_active,omitempty"`
		EditComment string   `json:"edit_comment,omitempty"`
	}

	var output []TransactionOutput
	for _, txn := range transactions {
		// S'assurer que les slices ne sont pas nil
		categories := txn.Categories
		if categories == nil {
			categories = []string{}
		}
		tags := txn.Tags
		if tags == nil {
			tags = []string{}
		}

		transactionOutput := TransactionOutput{
			ID:          txn.ID,
			Date:        txn.Date.Format("2006-01-02"),
			Amount:      txn.Amount,
			Description: txn.Description,
			Categories:  categories,
			Tags:        tags,
			EditComment: txn.EditComment,
		}

		// Ajouter is_active seulement si --history
		if showHistory {
			transactionOutput.IsActive = &txn.IsActive
		}

		output = append(output, transactionOutput)
	}

	jsonData, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	fmt.Println(string(jsonData))
	return nil
}
