package cli

import (
	"fmt"
)

func (c *CLI) handleBalance() error {
	if err := c.showBalances(); err != nil {
		return fmt.Errorf("error showing balances: %w", err)
	}
	return nil
}

func (c *CLI) showBalances() error {
	accounts, err := c.storage.GetAccounts()
	if err != nil {
		return err
	}

	fmt.Println("Account Balances:")
	for _, account := range accounts {
		if account.IsActive {
			balance, err := c.transactionService.GetAccountBalance(account.ID)
			if err != nil {
				return err
			}
			fmt.Printf("- %s: %.2f %s\n", account.Name, balance, account.Currency)
		}
	}

	return nil
}
