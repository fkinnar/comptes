package cli

import "fmt"

// Help messages for CLI commands
const (
	HelpUsage = `Usage: comptes <command>

Commands:
  init     - Initialize the project
  add      - Add transactions
  list     - List transactions
  edit     - Edit a transaction (soft delete + new)
  delete   - Delete a transaction (soft delete)
  undo     - Undo the last operation on a transaction
  balance  - Show account balances
  migrate  - Migrate old transaction IDs to UUID format`

	HelpAdd = `Usage: comptes add <json>

Example: comptes add '{"account":"BANQUE","amount":-25.50,"description":"Achat","categories":["ALM"]}'
Example: comptes add '{"account":"BANQUE","amount":-25.50,"description":"Achat","date":"today"}'

Date formats in JSON: 2024-01-15, 15/01/2024, yesterday, today, tomorrow`

	HelpEdit = `Usage: comptes edit <id> <json> -m <message>

Example: comptes edit fd6647d8 '{"amount": -30.00}' -m "Correction montant"

Note: You can use partial IDs like 'fd66' if unique
Note: Message is mandatory for edit operations`

	HelpDelete = `Usage: comptes delete <id> -m <message>

Example: comptes delete fd6647d8 -m "Transaction erronée"

Note: You can use partial IDs like 'fd66' if unique
Note: Message is mandatory for delete operations`

	HelpUndo = `Usage: comptes undo <id>

Example: comptes undo fd6647d8

Note: You can use partial IDs like 'fd66' if unique
Note: Undoes the last operation (add/edit/delete) on the transaction`
)

// ShowHelp displays help for a specific command
func ShowHelp(command string) {
	switch command {
	case "add":
		fmt.Println(HelpAdd)
	case "edit":
		fmt.Println(HelpEdit)
	case "delete":
		fmt.Println(HelpDelete)
	case "undo":
		fmt.Println(HelpUndo)
	default:
		fmt.Println(HelpUsage)
	}
}

// ShowError displays an error message with context
func ShowError(err error) {
	fmt.Printf("Error: %v\n", err)
}
