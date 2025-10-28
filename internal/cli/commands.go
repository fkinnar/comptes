package cli

import (
	"comptes/internal/errors"
	"comptes/internal/service"
	"comptes/internal/storage"
	"os"
	"path/filepath"
)

// CLI represents the command line interface
type CLI struct {
	transactionService *service.TransactionService
	storage            storage.Storage
}

// NewCLI creates a new CLI instance
func NewCLI() (*CLI, error) {
	// Get data directory (use environment variable for tests, or default to data/)
	dataDir := os.Getenv("COMPTES_DATA_DIR")
	if dataDir == "" {
		// Default behavior: next to executable for MVP
		execPath, err := os.Executable()
		if err != nil {
			return nil, errors.Wrap(errors.ErrorTypeSystem, "exec_path_failed", "Failed to get executable path", err)
		}
		dataDir = filepath.Join(filepath.Dir(execPath), "data")
	}

	// Create data directory if it doesn't exist
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, errors.Wrap(errors.ErrorTypeSystem, "data_dir_failed", "Failed to create data directory", err)
	}

	// Initialize storage and services
	storage := storage.NewJSONStorage(dataDir)
	transactionService := service.NewTransactionService(storage)

	return &CLI{
		transactionService: transactionService,
		storage:            storage,
	}, nil
}

// Execute runs the CLI command
func (c *CLI) Execute(args []string) error {
	if len(args) < 2 {
		ShowHelp("")
		return errors.MissingArguments("comptes")
	}

	command := args[1]

	switch command {
	case "init":
		return c.handleInit()
	case "add":
		return c.handleAdd(args)
	case "list":
		return c.handleList(args)
	case "edit":
		return c.handleEdit(args)
	case "delete":
		return c.handleDelete(args)
	case "undo":
		return c.handleUndo(args)
	case "balance":
		return c.handleBalance()
	case "migrate":
		return c.handleMigrate()
	default:
		return errors.InvalidCommand(command)
	}
}
