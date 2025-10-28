package cli

import (
	"comptes/internal/config"
	"comptes/internal/domain"
	"fmt"
	"os"
)

func (c *CLI) handleInit() error {
	fmt.Println("Initializing comptes project...")
	if err := c.initProject(); err != nil {
		return fmt.Errorf("error initializing project: %w", err)
	}
	fmt.Println("Project initialized successfully!")
	return nil
}

func (c *CLI) initProject() error {
	// Load configuration from YAML file, or create default if not exists
	configPath := config.GetConfigPath()
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		// If config file doesn't exist, create default configuration
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			fmt.Println("Creating default configuration...")
			cfg = config.CreateDefaultConfig()

			// Save default config to file
			if err := config.SaveConfig(configPath, cfg); err != nil {
				return fmt.Errorf("failed to save default configuration: %w", err)
			}
		} else {
			return fmt.Errorf("failed to load configuration: %w", err)
		}
	}

	// Save accounts from config
	if err := c.storage.SaveAccounts(cfg.Accounts); err != nil {
		return err
	}

	// Save categories from config
	if err := c.storage.SaveCategories(cfg.Categories); err != nil {
		return err
	}

	// Save tags from config
	if err := c.storage.SaveTags(cfg.Tags); err != nil {
		return err
	}

	// Create empty transactions file
	if err := c.storage.SaveTransactions([]domain.Transaction{}); err != nil {
		return err
	}

	return nil
}
