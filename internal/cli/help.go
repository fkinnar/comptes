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

	HelpList = `Usage: comptes list [options]

Options:
  --transactions     List transactions (default)
  --categories, -c   Show available categories
  --tags, -t         Show available tags
  --history, -h      Show all transactions (including deleted/edited)
  --format <fmt>     Output format: text (default), csv, json
  --codes            Show category/tag codes instead of names

Examples:
  comptes list                           # Show active transactions with full names
  comptes list --transactions --format csv  # Export transactions as CSV
  comptes list --categories              # Show available categories
  comptes list --categories --format csv # Export categories as CSV
  comptes list --categories --format json # Export categories as JSON
  comptes list --tags --format json      # Export tags as JSON
  comptes list --history                 # Show all transactions
  comptes list --codes                   # Show codes instead of names`

	HelpEdit = `Usage: comptes edit <id> <json> -m <message>

Example: comptes edit fd6647d8 '{"amount": -30.00}' -m "Correction montant"

Note: You can use partial IDs like 'fd66' if unique
Note: Message is mandatory for edit operations`

	HelpDelete = `Usage: comptes delete <id> -m <message> [options]

Options:
  --hard            Permanently delete the transaction (cannot be undone)
  -f, --force       Skip confirmation prompt for destructive operations

Examples:
  comptes delete fd6647d8 -m "Transaction erronée"           # Soft delete
  comptes delete fd6647d8 --hard -m "Dupliquée"              # Permanent delete
  comptes delete fd6647d8 --hard --force -m "Dupliquée"      # No confirmation

Note: You can use partial IDs like 'fd66' if unique
Note: Message is mandatory for delete operations`

	HelpUndo = `Usage: comptes undo <id> [options]

Options:
  --hard            Permanently remove the transaction (cannot be undone)
  -f, --force       Skip confirmation prompt for destructive operations

Examples:
  comptes undo fd6647d8                    # Soft undo (restore transaction)
  comptes undo fd6647d8 --hard             # Permanent removal
  comptes undo fd6647d8 --hard --force     # No confirmation

Note: You can use partial IDs like 'fd66' if unique
Note: Undoes the last operation (add/edit/delete) on the transaction`
)

// ShowHelp displays help for a specific command
func ShowHelp(command string) {
	switch command {
	case "add":
		fmt.Println(HelpAdd)
	case "list":
		fmt.Println(HelpList)
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
