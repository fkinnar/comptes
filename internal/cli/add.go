package cli

import (
	"comptes/internal/domain"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

// FlexibleDate permet de parser les dates relatives dans le JSON
type FlexibleDate struct {
	time.Time
}

// UnmarshalJSON implémente l'interface json.Unmarshaler
func (fd *FlexibleDate) UnmarshalJSON(data []byte) error {
	// Enlever les guillemets
	dateStr := strings.Trim(string(data), `"`)

	// Parser avec notre fonction parseDate
	parsedTime, err := parseDate(dateStr)
	if err != nil {
		return err
	}

	fd.Time = parsedTime
	return nil
}

// MarshalJSON implémente l'interface json.Marshaler
func (fd FlexibleDate) MarshalJSON() ([]byte, error) {
	return json.Marshal(fd.Time.Format("2006-01-02T15:04:05Z07:00"))
}

// TransactionInput structure pour le parsing JSON avec dates flexibles
type TransactionInput struct {
	ID          string       `json:"id"`
	Account     string       `json:"account"`
	Date        FlexibleDate `json:"date"`
	Amount      float64      `json:"amount"`
	Description string       `json:"description"`
	Categories  []string     `json:"categories"`
	Tags        []string     `json:"tags"`
	IsActive    bool         `json:"is_active"`
	CreatedAt   FlexibleDate `json:"created_at"`
	UpdatedAt   FlexibleDate `json:"updated_at"`
}

func (c *CLI) handleAdd(args []string) error {
	if len(args) < 3 {
		ShowHelp("add")
		return fmt.Errorf("missing arguments")
	}

	// Check if first argument is JSON (starts with { or [) or flags
	firstArg := args[2]
	isJSON := strings.HasPrefix(firstArg, "{") || strings.HasPrefix(firstArg, "[")

	var transaction domain.Transaction
	var providedBatchID string
	forceDirect := false

	if isJSON {
		// JSON mode (existing behavior)
		jsonData := firstArg

		// Parse remaining arguments for batch-id and --immediate
		for i := 3; i < len(args); i++ {
			if args[i] == "--immediate" || args[i] == "-i" {
				forceDirect = true
			} else if providedBatchID == "" {
				// First non-flag argument is treated as batch-id
				providedBatchID = args[i]
			}
		}

		// Parse JSON and create transaction
		var input TransactionInput
		if err := json.Unmarshal([]byte(jsonData), &input); err != nil {
			return fmt.Errorf("failed to parse JSON: %w", err)
		}

		// Load context for JSON mode (only if fields are missing in JSON AND batch is active)
		// Context is only valid within a transaction batch
		currentBatchID, err := c.getCurrentBatchID()
		if err == nil && currentBatchID != "" {
			context, err := c.getCurrentContext()
			if err == nil {
				// Use context only if fields are not provided in JSON
				if input.Account == "" && context.Account != "" {
					input.Account = context.Account
				}
				if len(input.Categories) == 0 && len(context.Categories) > 0 {
					input.Categories = context.Categories
				}
				if len(input.Tags) == 0 && len(context.Tags) > 0 {
					input.Tags = context.Tags
				}
			}
		}

		transaction = c.buildTransactionFromInput(input)
	} else {
		// Flags mode (new behavior)
		var account, description, date string
		var amount float64
		var categories, tags []string
		var hasAmount bool

		// Parse flags
		for i := 2; i < len(args); i++ {
			arg := args[i]

			switch {
			case arg == "-a" || arg == "--account":
				if i+1 < len(args) {
					account = args[i+1]
					i++
				} else {
					return fmt.Errorf("--account requires a value")
				}
			case arg == "-m" || arg == "--amount":
				if i+1 < len(args) {
					var err error
					amount, err = parseFloat(args[i+1])
					if err != nil {
						return fmt.Errorf("invalid amount: %w", err)
					}
					hasAmount = true
					i++
				} else {
					return fmt.Errorf("--amount requires a value")
				}
			case arg == "--desc" || arg == "--description" || arg == "-d":
				if i+1 < len(args) {
					description = args[i+1]
					i++
				} else {
					return fmt.Errorf("--description requires a value")
				}
			case arg == "-c" || arg == "--categories":
				if i+1 < len(args) {
					categories = parseList(args[i+1])
					i++
				} else {
					return fmt.Errorf("--categories requires a value")
				}
			case arg == "-t" || arg == "--tags":
				if i+1 < len(args) {
					tags = parseList(args[i+1])
					i++
				} else {
					return fmt.Errorf("--tags requires a value")
				}
			case arg == "--date" || arg == "--on" || arg == "-o":
				if i+1 < len(args) {
					date = args[i+1]
					i++
				} else {
					return fmt.Errorf("--date requires a value")
				}
			case arg == "--immediate" || arg == "-i":
				forceDirect = true
			case strings.HasPrefix(arg, "-"):
				// Unknown flag
				return fmt.Errorf("unknown flag: %s", arg)
			default:
				// Treat as batch-id if not a flag
				if providedBatchID == "" {
					providedBatchID = arg
				}
			}
		}

		// Load context if fields are missing AND if we're in a batch
		// Context is only valid within a transaction batch
		currentBatchID, err := c.getCurrentBatchID()
		if err == nil && currentBatchID != "" {
			context, err := c.getCurrentContext()
			if err != nil {
				return fmt.Errorf("error loading context: %w", err)
			}

			// Use context for missing fields only if batch is active
			if account == "" {
				account = context.Account
			}
			if len(categories) == 0 && len(context.Categories) > 0 {
				categories = context.Categories
			}
			if len(tags) == 0 && len(context.Tags) > 0 {
				tags = context.Tags
			}
		}

		// Validate required fields
		if account == "" {
			// Check if we're in a batch context
			currentBatchID, _ := c.getCurrentBatchID()
			if currentBatchID != "" {
				return fmt.Errorf("account is required (use -a/--account or set context with 'comptes account <id>')")
			}
			return fmt.Errorf("account is required (use -a/--account)")
		}
		if !hasAmount {
			return fmt.Errorf("amount is required (use -m/--amount)")
		}
		if description == "" {
			return fmt.Errorf("description is required (use --desc/--description)")
		}

		// Build transaction from flags
		transaction.ID = c.generateShortID()
		transaction.Account = account
		transaction.Amount = amount
		transaction.Description = description
		transaction.Categories = categories
		transaction.Tags = tags
		transaction.IsActive = true
		transaction.CreatedAt = time.Now()
		transaction.UpdatedAt = time.Now()

		// Parse date if provided
		if date != "" {
			parsedDate, err := parseDate(date)
			if err != nil {
				return fmt.Errorf("invalid date format: %w", err)
			}
			transaction.Date = parsedDate
		} else {
			transaction.Date = time.Now()
		}
	}

	// If --immediate flag is set, add directly regardless of batch
	if forceDirect {
		if err := c.transactionService.AddTransaction(transaction); err != nil {
			return fmt.Errorf("error adding transaction: %w", err)
		}
		fmt.Println("Transaction added successfully!")
		return nil
	}

	// Normal behavior: check for batch
	batchID, err := c.resolveBatchID(providedBatchID)
	if err == nil && batchID != "" {
		// Add to batch
		if err := c.batchService.AddTransactionToBatch(batchID, transaction); err != nil {
			return fmt.Errorf("error adding transaction to batch: %w", err)
		}
		fmt.Printf("Transaction added to batch %s successfully!\n", batchID)
	} else {
		// Add directly (no batch ID provided and no current batch)
		if err := c.transactionService.AddTransaction(transaction); err != nil {
			return fmt.Errorf("error adding transaction: %w", err)
		}
		fmt.Println("Transaction added successfully!")
	}
	return nil
}

// buildTransactionFromInput converts TransactionInput to domain.Transaction
func (c *CLI) buildTransactionFromInput(input TransactionInput) domain.Transaction {
	transaction := domain.Transaction{
		ID:          input.ID,
		Account:     input.Account,
		Amount:      input.Amount,
		Description: input.Description,
		Categories:  input.Categories,
		Tags:        input.Tags,
		IsActive:    input.IsActive,
	}

	// Generate ID if not provided
	if transaction.ID == "" {
		transaction.ID = c.generateShortID()
	}

	// Date : JSON ou par défaut
	if !input.Date.IsZero() {
		transaction.Date = input.Date.Time
	} else {
		transaction.Date = time.Now()
	}

	// CreatedAt et UpdatedAt
	if !input.CreatedAt.IsZero() {
		transaction.CreatedAt = input.CreatedAt.Time
	} else {
		transaction.CreatedAt = time.Now()
	}

	if !input.UpdatedAt.IsZero() {
		transaction.UpdatedAt = input.UpdatedAt.Time
	} else {
		transaction.UpdatedAt = time.Now()
	}

	transaction.IsActive = true
	return transaction
}

// parseFloat parses a float64 string
func parseFloat(s string) (float64, error) {
	var result float64
	_, err := fmt.Sscanf(s, "%f", &result)
	return result, err
}

// parseList parses a comma-separated list into a []string
func parseList(s string) []string {
	if s == "" {
		return []string{}
	}
	parts := strings.Split(s, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			result = append(result, part)
		}
	}
	return result
}

// generateShortID génère un UUID complet (36 caractères)
// UUID v4 garantit l'unicité absolue sans vérification de collision
func (c *CLI) generateShortID() string {
	return uuid.New().String()
}

// parseDate parses various date formats
func parseDate(dateStr string) (time.Time, error) {
	// Handle relative dates
	switch strings.ToLower(dateStr) {
	case "today":
		return time.Now(), nil
	case "yesterday":
		return time.Now().AddDate(0, 0, -1), nil
	case "tomorrow":
		return time.Now().AddDate(0, 0, 1), nil
	}

	// Try different date formats
	formats := []string{
		"2006-01-02",          // 2024-01-15
		"02/01/2006",          // 15/01/2024
		"02-01-2006",          // 15-01-2024
		"2006-01-02 15:04:05", // 2024-01-15 14:30:00
		"02/01/2006 15:04:05", // 15/01/2024 14:30:00
	}

	for _, format := range formats {
		if t, err := time.Parse(format, dateStr); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("unable to parse date: %s", dateStr)
}
