package main

import (
	"comptes/internal/cli"
	"comptes/internal/errors"
	"fmt"
	"os"
)

func main() {
	// Create CLI instance
	c, err := cli.NewCLI()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to initialize CLI: %v\n", err)
		os.Exit(1)
	}

	// Execute command
	if err := c.Execute(os.Args); err != nil {
		// Check if it's a known error type
		if e, ok := err.(*errors.ComptesError); ok {
			if e.Type == errors.ErrorTypeUserInput {
				// User errors (invalid command, missing args, etc.) - show error and exit 1
				cli.ShowError(err)
				os.Exit(1)
			} else {
				// System/validation/business errors - show error and exit 2
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(2)
			}
		} else {
			// Unknown error type - show error and exit 2
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(2)
		}
	}
}
