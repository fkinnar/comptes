package domain

import (
	"time"
)

// Account represents a bank account
type Account struct {
	ID             string    `json:"id" yaml:"id"`
	Name           string    `json:"name" yaml:"name"`
	Type           string    `json:"type" yaml:"type"`
	Currency       string    `json:"currency" yaml:"currency"`
	InitialBalance float64   `json:"initial_balance" yaml:"initial_balance"`
	IsActive       bool      `json:"is_active" yaml:"is_active"`
	CreatedAt      time.Time `json:"created_at" yaml:"created_at"`
}

// Transaction represents a financial transaction
type Transaction struct {
	ID          string    `json:"id"`
	Account     string    `json:"account"`
	Date        time.Time `json:"date"`
	Amount      float64   `json:"amount"`
	Description string    `json:"description"`
	Categories  []string  `json:"categories"`
	Tags        []string  `json:"tags"`
	IsActive    bool      `json:"is_active"`
	EditComment string    `json:"edit_comment,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Category represents a transaction category
type Category struct {
	Code        string   `json:"code"`
	Name        string   `json:"name"`
	Parent      *string  `json:"parent,omitempty"`
	Children    []string `json:"children"`
	Description string   `json:"description"`
}

// Tag represents a transaction tag
type Tag struct {
	Code        string   `json:"code"`
	Name        string   `json:"name"`
	Parent      *string  `json:"parent,omitempty"`
	Children    []string `json:"children"`
	Description string   `json:"description"`
}
