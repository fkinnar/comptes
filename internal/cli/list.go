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
	showCategories := false
	showTags := false
	showTransactions := true // Par défaut, on liste les transactions
	showCodes := false

	// Check for flags
	for i, arg := range args {
		if arg == "--format" && i+1 < len(args) {
			format = args[i+1]
		}
		if arg == "--history" {
			showHistory = true
		}
		if arg == "--categories" || arg == "-c" {
			showCategories = true
			showTransactions = false
		}
		if arg == "--tags" || arg == "-t" {
			showTags = true
			showTransactions = false
		}
		if arg == "--transactions" {
			showTransactions = true
		}
		if arg == "--codes" {
			showCodes = true
		}
	}

	// Handle different list types
	if showCategories {
		return c.showCategories(format)
	}
	if showTags {
		return c.showTags(format)
	}
	if showTransactions {
		return c.listTransactions(format, showHistory, showCodes)
	}

	// Fallback: liste les transactions par défaut
	return c.listTransactions(format, showHistory, showCodes)
}

func (c *CLI) listTransactions(format string, showHistory bool, showCodes bool) error {
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
		return c.listTransactionsCSV(filteredTransactions, showHistory, showCodes)
	case "json":
		return c.listTransactionsJSON(filteredTransactions, showHistory, showCodes)
	case "text":
		fallthrough
	default:
		return c.listTransactionsText(filteredTransactions, showHistory, showCodes)
	}
}

func (c *CLI) listTransactionsText(transactions []domain.Transaction, showHistory bool, showCodes bool) error {
	if showHistory {
		fmt.Println("Transactions (all):")
	} else {
		fmt.Println("Transactions:")
	}

	// Charger les catégories et tags pour la conversion des codes en noms
	categories, _ := c.storage.GetCategories()
	tags, _ := c.storage.GetTags()

	// Créer des maps pour la conversion rapide
	categoryMap := make(map[string]string)
	for _, cat := range categories {
		categoryMap[cat.Code] = cat.Name
	}

	tagMap := make(map[string]string)
	for _, tag := range tags {
		tagMap[tag.Code] = tag.Name
	}

	for _, txn := range transactions {
		// Convertir les codes en noms si nécessaire
		var categoriesDisplay []string
		var tagsDisplay []string

		if showCodes {
			categoriesDisplay = txn.Categories
			tagsDisplay = txn.Tags
		} else {
			// Convertir les codes en noms
			for _, code := range txn.Categories {
				if name, exists := categoryMap[code]; exists {
					categoriesDisplay = append(categoriesDisplay, name)
				} else {
					categoriesDisplay = append(categoriesDisplay, code) // Fallback si pas trouvé
				}
			}

			for _, code := range txn.Tags {
				if name, exists := tagMap[code]; exists {
					tagsDisplay = append(tagsDisplay, name)
				} else {
					tagsDisplay = append(tagsDisplay, code) // Fallback si pas trouvé
				}
			}
		}

		// Format de base avec ID
		var line string
		if showHistory {
			// Indicateur d'état seulement pour --history
			status := "✅"
			if !txn.IsActive {
				status = "❌"
			}
			line = fmt.Sprintf("- [%s] %s %s: %.2f EUR - %s (Categories: %v)",
				txn.ID, status, txn.Date.Format("2006-01-02"), txn.Amount, txn.Description, categoriesDisplay)
		} else {
			// Pas d'indicateur de statut pour list normal (toutes sont actives)
			line = fmt.Sprintf("- [%s] %s: %.2f EUR - %s (Categories: %v)",
				txn.ID, txn.Date.Format("2006-01-02"), txn.Amount, txn.Description, categoriesDisplay)
		}

		// Ajouter les tags seulement s'il y en a
		if len(tagsDisplay) > 0 {
			line += fmt.Sprintf(", Tags: %v", tagsDisplay)
		}

		// Ajouter le commentaire d'edit si présent
		if txn.EditComment != "" {
			line += fmt.Sprintf(" | Edit: %s", txn.EditComment)
		}

		fmt.Println(line)
	}
	return nil
}

func (c *CLI) listTransactionsCSV(transactions []domain.Transaction, showHistory bool, showCodes bool) error {
	// En-tête CSV
	if showHistory {
		fmt.Println("id,date,amount,description,categories,tags,is_active,edit_comment")
	} else {
		fmt.Println("id,date,amount,description,categories,tags")
	}

	// Charger les catégories et tags pour la conversion des codes en noms
	categories, _ := c.storage.GetCategories()
	tags, _ := c.storage.GetTags()

	// Créer des maps pour la conversion rapide
	categoryMap := make(map[string]string)
	for _, cat := range categories {
		categoryMap[cat.Code] = cat.Name
	}

	tagMap := make(map[string]string)
	for _, tag := range tags {
		tagMap[tag.Code] = tag.Name
	}

	for _, txn := range transactions {
		// Convertir les codes en noms si nécessaire
		var categoriesDisplay []string
		var tagsDisplay []string

		if showCodes {
			categoriesDisplay = txn.Categories
			tagsDisplay = txn.Tags
		} else {
			// Convertir les codes en noms
			for _, code := range txn.Categories {
				if name, exists := categoryMap[code]; exists {
					categoriesDisplay = append(categoriesDisplay, name)
				} else {
					categoriesDisplay = append(categoriesDisplay, code) // Fallback si pas trouvé
				}
			}

			for _, code := range txn.Tags {
				if name, exists := tagMap[code]; exists {
					tagsDisplay = append(tagsDisplay, name)
				} else {
					tagsDisplay = append(tagsDisplay, code) // Fallback si pas trouvé
				}
			}
		}

		// Convertir les catégories en string
		categoriesStr := ""
		if len(categoriesDisplay) > 0 {
			categoriesStr = strings.Join(categoriesDisplay, ";")
		}

		// Convertir les tags en string
		tagsStr := ""
		if len(tagsDisplay) > 0 {
			tagsStr = strings.Join(tagsDisplay, ";")
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

func (c *CLI) listTransactionsJSON(transactions []domain.Transaction, showHistory bool, showCodes bool) error {
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

// showCategories displays all available categories
func (c *CLI) showCategories(format string) error {
	categories, err := c.storage.GetCategories()
	if err != nil {
		return fmt.Errorf("error loading categories: %w", err)
	}

	if len(categories) == 0 {
		fmt.Println("No categories found.")
		return nil
	}

	switch format {
	case "csv":
		return c.showCategoriesCSV(categories)
	case "json":
		return c.showCategoriesJSON(categories)
	default:
		return c.showCategoriesText(categories)
	}
}

// showCategoriesText displays categories in text format
func (c *CLI) showCategoriesText(categories []domain.Category) error {
	fmt.Println("Available categories:")
	fmt.Println("===================")
	for _, cat := range categories {
		fmt.Printf("• %s (%s) - %s\n", cat.Name, cat.Code, cat.Description)
	}
	return nil
}

// showCategoriesCSV displays categories in CSV format
func (c *CLI) showCategoriesCSV(categories []domain.Category) error {
	fmt.Println("code,name,description")
	for _, cat := range categories {
		// Échapper les virgules dans la description en utilisant des guillemets
		description := strings.ReplaceAll(cat.Description, "\"", "\"\"") // Échapper les guillemets existants
		fmt.Printf("%s,%s,\"%s\"\n", cat.Code, cat.Name, description)
	}
	return nil
}

// showCategoriesJSON displays categories in JSON format
func (c *CLI) showCategoriesJSON(categories []domain.Category) error {
	jsonData, err := json.MarshalIndent(categories, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}
	fmt.Println(string(jsonData))
	return nil
}

// showTags displays all available tags
func (c *CLI) showTags(format string) error {
	tags, err := c.storage.GetTags()
	if err != nil {
		return fmt.Errorf("error loading tags: %w", err)
	}

	if len(tags) == 0 {
		fmt.Println("No tags found.")
		return nil
	}

	switch format {
	case "csv":
		return c.showTagsCSV(tags)
	case "json":
		return c.showTagsJSON(tags)
	default:
		return c.showTagsText(tags)
	}
}

// showTagsText displays tags in text format
func (c *CLI) showTagsText(tags []domain.Tag) error {
	fmt.Println("Available tags:")
	fmt.Println("==============")
	for _, tag := range tags {
		fmt.Printf("• %s (%s) - %s\n", tag.Name, tag.Code, tag.Description)
	}
	return nil
}

// showTagsCSV displays tags in CSV format
func (c *CLI) showTagsCSV(tags []domain.Tag) error {
	fmt.Println("code,name,description")
	for _, tag := range tags {
		// Échapper les virgules dans la description en utilisant des guillemets
		description := strings.ReplaceAll(tag.Description, "\"", "\"\"") // Échapper les guillemets existants
		fmt.Printf("%s,%s,\"%s\"\n", tag.Code, tag.Name, description)
	}
	return nil
}

// showTagsJSON displays tags in JSON format
func (c *CLI) showTagsJSON(tags []domain.Tag) error {
	jsonData, err := json.MarshalIndent(tags, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}
	fmt.Println(string(jsonData))
	return nil
}
