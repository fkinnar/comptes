package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// TransactionContext holds shared context for transactions
type TransactionContext struct {
	Account    string   `json:"account,omitempty"`
	Categories []string `json:"categories,omitempty"`
	Tags       []string `json:"tags,omitempty"`
}

// getCurrentContext reads the current context from the file
func (c *CLI) getCurrentContext() (*TransactionContext, error) {
	contextFile := filepath.Join(c.dataDir, ".current_context")
	data, err := os.ReadFile(contextFile)
	if err != nil {
		if os.IsNotExist(err) {
			return &TransactionContext{}, nil // Empty context
		}
		return nil, fmt.Errorf("failed to read context file: %w", err)
	}

	var context TransactionContext
	if err := json.Unmarshal(data, &context); err != nil {
		return nil, fmt.Errorf("failed to parse context file: %w", err)
	}

	return &context, nil
}

// saveCurrentContext saves the current context to the file
func (c *CLI) saveCurrentContext(context *TransactionContext) error {
	contextFile := filepath.Join(c.dataDir, ".current_context")

	// Check if context is empty
	isEmpty := context.Account == "" && len(context.Categories) == 0 && len(context.Tags) == 0

	if isEmpty {
		// Remove file if context is empty
		if err := os.Remove(contextFile); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("failed to remove context file: %w", err)
		}
		return nil
	}

	data, err := json.MarshalIndent(context, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal context: %w", err)
	}

	if err := os.WriteFile(contextFile, data, 0644); err != nil {
		return fmt.Errorf("failed to save context file: %w", err)
	}

	return nil
}

// handleAccount sets the default account in the context
func (c *CLI) handleAccount(args []string) error {
	if len(args) < 3 {
		// Show current context
		context, err := c.getCurrentContext()
		if err != nil {
			return fmt.Errorf("error reading context: %w", err)
		}

		if context.Account == "" {
			fmt.Println("No account set in context.")
			fmt.Println("Usage: comptes account <account-id>")
			fmt.Println("Example: comptes account BANQUE")
			fmt.Println("Note: Context is only valid within a transaction batch. Use 'comptes begin' first.")
		} else {
			fmt.Printf("Current account: %s\n", context.Account)
		}
		return nil
	}

	accountID := args[2]

	// Check if a batch is active
	currentBatchID, err := c.getCurrentBatchID()
	if err != nil || currentBatchID == "" {
		return fmt.Errorf("no active transaction batch. Use 'comptes begin' to start a batch first, then set the context.")
	}

	// Validate account exists
	accounts, err := c.storage.GetAccounts()
	if err != nil {
		return fmt.Errorf("error loading accounts: %w", err)
	}

	accountExists := false
	for _, acc := range accounts {
		if acc.ID == accountID {
			accountExists = true
			break
		}
	}

	if !accountExists {
		return fmt.Errorf("account not found: %s", accountID)
	}

	// Update context
	context, err := c.getCurrentContext()
	if err != nil {
		return fmt.Errorf("error reading context: %w", err)
	}

	context.Account = accountID
	if err := c.saveCurrentContext(context); err != nil {
		return fmt.Errorf("error saving context: %w", err)
	}

	fmt.Printf("Account set to: %s\n", accountID)
	return nil
}

// handleCategory sets the default categories in the context
func (c *CLI) handleCategory(args []string) error {
	if len(args) < 3 {
		// Show current context
		context, err := c.getCurrentContext()
		if err != nil {
			return fmt.Errorf("error reading context: %w", err)
		}

		if len(context.Categories) == 0 {
			fmt.Println("No categories set in context.")
			fmt.Println("Usage: comptes category <category-code> [category-code2 ...]")
			fmt.Println("Example: comptes category ALM")
			fmt.Println("Example: comptes category ALM SLR")
			fmt.Println("Note: Context is only valid within a transaction batch. Use 'comptes begin' first.")
		} else {
			fmt.Printf("Current categories: %v\n", context.Categories)
		}
		return nil
	}

	categories := args[2:]

	// Check if a batch is active
	currentBatchID, err := c.getCurrentBatchID()
	if err != nil || currentBatchID == "" {
		return fmt.Errorf("no active transaction batch. Use 'comptes begin' to start a batch first, then set the context.")
	}

	// Validate categories exist
	allCategories, err := c.storage.GetCategories()
	if err != nil {
		return fmt.Errorf("error loading categories: %w", err)
	}

	categoryMap := make(map[string]bool)
	for _, cat := range allCategories {
		categoryMap[cat.Code] = true
	}

	for _, catCode := range categories {
		if !categoryMap[catCode] {
			return fmt.Errorf("category not found: %s", catCode)
		}
	}

	// Update context
	context, err := c.getCurrentContext()
	if err != nil {
		return fmt.Errorf("error reading context: %w", err)
	}

	context.Categories = categories
	if err := c.saveCurrentContext(context); err != nil {
		return fmt.Errorf("error saving context: %w", err)
	}

	fmt.Printf("Categories set to: %v\n", categories)
	return nil
}

// handleTags sets the default tags in the context
func (c *CLI) handleTags(args []string) error {
	if len(args) < 3 {
		// Show current context
		context, err := c.getCurrentContext()
		if err != nil {
			return fmt.Errorf("error reading context: %w", err)
		}

		if len(context.Tags) == 0 {
			fmt.Println("No tags set in context.")
			fmt.Println("Usage: comptes tags <tag-code> [tag-code2 ...]")
			fmt.Println("Example: comptes tags REC")
			fmt.Println("Example: comptes tags REC URG")
			fmt.Println("Note: Context is only valid within a transaction batch. Use 'comptes begin' first.")
		} else {
			fmt.Printf("Current tags: %v\n", context.Tags)
		}
		return nil
	}

	tags := args[2:]

	// Check if a batch is active
	currentBatchID, err := c.getCurrentBatchID()
	if err != nil || currentBatchID == "" {
		return fmt.Errorf("no active transaction batch. Use 'comptes begin' to start a batch first, then set the context.")
	}

	// Validate tags exist
	allTags, err := c.storage.GetTags()
	if err != nil {
		return fmt.Errorf("error loading tags: %w", err)
	}

	tagMap := make(map[string]bool)
	for _, tag := range allTags {
		tagMap[tag.Code] = true
	}

	for _, tagCode := range tags {
		if !tagMap[tagCode] {
			return fmt.Errorf("tag not found: %s", tagCode)
		}
	}

	// Update context
	context, err := c.getCurrentContext()
	if err != nil {
		return fmt.Errorf("error reading context: %w", err)
	}

	context.Tags = tags
	if err := c.saveCurrentContext(context); err != nil {
		return fmt.Errorf("error saving context: %w", err)
	}

	fmt.Printf("Tags set to: %v\n", tags)
	return nil
}

// handleContextClear clears the current context
func (c *CLI) handleContextClear() error {
	context := &TransactionContext{}
	if err := c.saveCurrentContext(context); err != nil {
		return fmt.Errorf("error clearing context: %w", err)
	}

	fmt.Println("Context cleared.")
	return nil
}

// handleContextShow shows the current context
func (c *CLI) handleContextShow() error {
	context, err := c.getCurrentContext()
	if err != nil {
		return fmt.Errorf("error reading context: %w", err)
	}

	fmt.Println("Current context:")
	if context.Account != "" {
		fmt.Printf("  Account: %s\n", context.Account)
	} else {
		fmt.Println("  Account: (not set)")
	}

	if len(context.Categories) > 0 {
		fmt.Printf("  Categories: %v\n", context.Categories)
	} else {
		fmt.Println("  Categories: (not set)")
	}

	if len(context.Tags) > 0 {
		fmt.Printf("  Tags: %v\n", context.Tags)
	} else {
		fmt.Println("  Tags: (not set)")
	}

	return nil
}
