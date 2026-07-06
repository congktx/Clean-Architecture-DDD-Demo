package config

import (
	"database/sql"
	"fmt"
	"log"
)

func InitDB(db *sql.DB) error {
	log.Println("Initializing database tables...")

	query := `
		-- Table for Wallets
		CREATE TABLE IF NOT EXISTS wallets (
			id VARCHAR(36) PRIMARY KEY,
			owner_id VARCHAR(36) NOT NULL,
			name VARCHAR(255) NOT NULL,
			status VARCHAR(50) NOT NULL,
			balance_amount VARCHAR(50) NOT NULL,
			balance_currency VARCHAR(10) NOT NULL
		);

		-- Table for Wallet Transactions
		CREATE TABLE IF NOT EXISTS wallet_transactions (
			id VARCHAR(36) PRIMARY KEY,
			wallet_id VARCHAR(36) NOT NULL,
			type VARCHAR(50) NOT NULL,
			amount VARCHAR(50) NOT NULL,
			currency VARCHAR(10) NOT NULL,
			category_id VARCHAR(36) NOT NULL,
			timestamp BIGINT NOT NULL,
			description TEXT,
			FOREIGN KEY (wallet_id) REFERENCES wallets(id) ON DELETE CASCADE
		);

		-- Table for Budgets
		CREATE TABLE IF NOT EXISTS budgets (
			id VARCHAR(36) PRIMARY KEY,
			owner_id VARCHAR(36) NOT NULL,
			month INT NOT NULL,
			year INT NOT NULL,
			total_limit_amount VARCHAR(50) NOT NULL,
			total_limit_currency VARCHAR(10) NOT NULL
		);

		-- Table for Budget Items
		CREATE TABLE IF NOT EXISTS budget_items (
			id VARCHAR(36) PRIMARY KEY,
			budget_id VARCHAR(36) NOT NULL,
			category_id VARCHAR(36) NOT NULL,
			allocated_amount VARCHAR(50) NOT NULL,
			spent_amount VARCHAR(50) NOT NULL,
			currency VARCHAR(10) NOT NULL,
			FOREIGN KEY (budget_id) REFERENCES budgets(id) ON DELETE CASCADE
		);
	`

	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create tables: %w", err)
	}

	log.Println("Database tables initialized successfully.")
	return nil
}
