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
  begin    - Begin a new transaction batch
  commit   - Commit a pending transaction batch
  rollback - Rollback a pending transaction batch
  account  - Set default account in context
  category - Set default categories in context
  tags     - Set default tags in context
  context  - Show or clear transaction context`

	HelpAdd = `Usage: comptes add <json> [batch-id] [--immediate]
       comptes add [flags] [batch-id] [--immediate]

Add a transaction using either JSON format or flags.

JSON format:
  comptes add '{"account":"BANQUE","amount":-25.50,"description":"Achat","categories":["ALM"]}'
  comptes add '{"account":"BANQUE","amount":-25.50,"description":"Achat","date":"today"}' [batch-id]
  comptes add '{"account":"BANQUE","amount":-25.50,"description":"Achat"}' --immediate

Flags format:
  comptes add -a BANQUE -m -25.50 -d "Achat" -c ALM -t REC
  comptes add --account BANQUE --amount -25.50 --description "Achat" --categories "ALM,REC" --tags "URG"
  comptes add -a BANQUE -m -25.50 -d "Achat" -o today [batch-id]
  comptes add -a BANQUE -m -25.50 -d "Achat" --immediate

Flags:
  -a, --account <id>        Account ID (required)
  -m, --amount <value>      Amount (required, negative for expense, positive for income)
  -d, --desc, --description <text>  Description (required)
  -c, --categories <codes>  Categories (comma-separated, e.g., "ALM,SLR")
  -t, --tags <codes>        Tags (comma-separated, e.g., "REC,URG")
  -o, --on, --date <date>   Date (optional, formats: today, yesterday, 2024-01-15)
  -i, --immediate           Force immediate addition even if a batch is in progress

If batch-id is provided (or current batch is set), the transaction is added to the pending batch instead of directly.
Use --immediate (or -i) flag to force immediate addition even if a batch is in progress.

Date formats: 2024-01-15, 15/01/2024, yesterday, today, tomorrow`

	HelpList = `Usage: comptes list [options]

Options:
  --transactions, -T  List transactions (default)
  --categories, -c   Show available categories
  --tags, -t         Show available tags
  --accounts, -a     Show available accounts with balances
  --history, -h      Show all transactions (including deleted/edited)
  --format <fmt>, -F Output format: text (default), csv, json
  --codes, -k        Show category/tag codes instead of names
  --help, -?         Show this help message

Examples:
  comptes list                           # Show active transactions with full names
  comptes list --transactions --format csv  # Export transactions as CSV
  comptes list --categories              # Show available categories
  comptes list --categories --format csv # Export categories as CSV
  comptes list --categories --format json # Export categories as JSON
  comptes list --tags --format json      # Export tags as JSON
  comptes list --accounts                # Show available accounts with balances
  comptes list --accounts --format csv   # Export accounts as CSV
  comptes list --accounts --format json # Export accounts as JSON
  comptes list --history                 # Show all transactions
  comptes list --codes                   # Show codes instead of names`

	HelpEdit = `Usage: comptes edit <id> <json> -m <message>

Example: comptes edit fd6647d8 '{"amount": -30.00}' -m "Correction montant"

Note: You can use partial IDs like 'fd66' if unique
Note: Message is mandatory for edit operations`

	HelpDelete = `Usage: comptes delete <id> -m <message> [options]

Options:
  -m, --message     Message explaining the deletion (required)
  --hard, -H        Permanently delete the transaction (cannot be undone)
  -f, --force       Skip confirmation prompt for destructive operations
  --help, -?        Show this help message

Examples:
  comptes delete fd6647d8 -m "Transaction erronée"           # Soft delete
  comptes delete fd6647d8 --hard -m "Dupliquée"              # Permanent delete
  comptes delete fd6647d8 --hard --force -m "Dupliquée"      # No confirmation
  comptes delete fd6647d8 -H -f -m "Dupliquée"                # Short form

Note: You can use partial IDs like 'fd66' if unique
Note: Message is mandatory for delete operations`

	HelpUndo = `Usage: comptes undo <id> [options]

Options:
  --hard, -H        Permanently remove the transaction (cannot be undone)
  -f, --force       Skip confirmation prompt for destructive operations
  --help, -?        Show this help message

Examples:
  comptes undo fd6647d8                    # Soft undo (restore transaction)
  comptes undo fd6647d8 --hard             # Permanent removal
  comptes undo fd6647d8 --hard --force     # No confirmation
  comptes undo fd6647d8 -H -f              # Short form

Note: You can use partial IDs like 'fd66' if unique
Note: Undoes the last operation (add/edit/delete) on the transaction`

	HelpBegin = `Usage: comptes begin [description]

Examples:
  comptes begin
  comptes begin "Monthly expenses"

Creates a new pending transaction batch and sets it as the current batch. You can then add transactions to it using:
  comptes add '{"account":"BANQUE","amount":-25.50}'

Use 'comptes commit' (or 'comptes commit <batch-id>') to commit or 'comptes rollback' to rollback.`

	HelpCommit = `Usage: comptes commit [batch-id]

Examples:
  comptes commit                    # Commits the current batch
  comptes commit abc12345           # Commits batch with partial ID
  comptes commit abc12345-b623-...  # Commits batch with full UUID

Commits a pending transaction batch. All transactions in the batch are added to the main transactions file.
If no batch-id is provided, uses the current batch (set by 'comptes begin').
You can use partial batch IDs if unique.`

	HelpRollback = `Usage: comptes rollback [batch-id]

Examples:
  comptes rollback                    # Rollbacks the current batch
  comptes rollback abc12345           # Rollbacks batch with partial ID
  comptes rollback abc12345-b623-...  # Rollbacks batch with full UUID

Rolls back a pending transaction batch. All transactions in the batch are discarded.
If no batch-id is provided, uses the current batch (set by 'comptes begin').
You can use partial batch IDs if unique.`

	HelpAccount = `Usage: comptes account [account-id]

Sets the default account in the transaction context. This account will be used automatically
when adding transactions if no account is specified.

Examples:
  comptes account BANQUE              # Set default account
  comptes account                     # Show current account

When account is set, you can use simplified add commands:
  comptes add -m -25.50 --desc "Courses"  # Uses context account automatically`

	HelpCategory = `Usage: comptes category [category-code ...]

Sets the default categories in the transaction context. These categories will be used automatically
when adding transactions if no categories are specified.

Examples:
  comptes category ALM                # Set single category
  comptes category ALM SLR            # Set multiple categories
  comptes category                    # Show current categories

When categories are set, you can use simplified add commands:
  comptes add -a BANQUE -m -25.50 --desc "Courses"  # Uses context categories automatically`

	HelpTags = `Usage: comptes tags [tag-code ...]

Sets the default tags in the transaction context. These tags will be used automatically
when adding transactions if no tags are specified.

Examples:
  comptes tags REC                    # Set single tag
  comptes tags REC URG                # Set multiple tags
  comptes tags                        # Show current tags

When tags are set, you can use simplified add commands:
  comptes add -a BANQUE -m -25.50 --desc "Courses"  # Uses context tags automatically`

	HelpContext = `Usage: comptes context [clear]

Shows or clears the current transaction context.

Examples:
  comptes context                     # Show current context
  comptes context clear               # Clear all context settings

The context includes:
  - Default account (set with 'comptes account')
  - Default categories (set with 'comptes category')
  - Default tags (set with 'comptes tags')

Context is used automatically when adding transactions if fields are not specified.`
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
	case "begin":
		fmt.Println(HelpBegin)
	case "commit":
		fmt.Println(HelpCommit)
	case "rollback":
		fmt.Println(HelpRollback)
	case "account":
		fmt.Println(HelpAccount)
	case "category":
		fmt.Println(HelpCategory)
	case "tags":
		fmt.Println(HelpTags)
	case "context":
		fmt.Println(HelpContext)
	default:
		fmt.Println(HelpUsage)
	}
}

// ShowError displays an error message with context
func ShowError(err error) {
	fmt.Printf("Error: %v\n", err)
}
