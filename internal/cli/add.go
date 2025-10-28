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
		fmt.Println("Usage: comptes add <json>")
		fmt.Println("Example: comptes add '{\"account\":\"BANQUE\",\"amount\":-25.50,\"description\":\"Achat\",\"categories\":[\"ALM\"]}'")
		fmt.Println("Example: comptes add '{\"account\":\"BANQUE\",\"amount\":-25.50,\"description\":\"Achat\",\"date\":\"today\"}'")
		fmt.Println("Date formats in JSON: 2024-01-15, 15/01/2024, yesterday, today, tomorrow")
		return fmt.Errorf("missing JSON data")
	}

	jsonData := args[2]
	if err := c.addTransaction(jsonData); err != nil {
		return fmt.Errorf("error adding transaction: %w", err)
	}
	fmt.Println("Transaction added successfully!")
	return nil
}

func (c *CLI) addTransaction(jsonData string) error {
	var input TransactionInput
	if err := json.Unmarshal([]byte(jsonData), &input); err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}

	// Convertir TransactionInput en domain.Transaction
	transaction := domain.Transaction{
		ID:          input.ID,
		Account:     input.Account,
		Amount:      input.Amount,
		Description: input.Description,
		Categories:  input.Categories,
		Tags:        input.Tags,
		IsActive:    input.IsActive,
	}

	// Gérer les dates
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

	return c.transactionService.AddTransaction(transaction)
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
