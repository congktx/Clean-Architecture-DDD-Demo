package infrastructure

import (
	"database/sql"
	"fmt"
	"github.com/shopspring/decimal"
	shared "src/src/shared/domain"
	"src/src/wallet/domain"
)

type PostgresWalletRepository struct {
	db *sql.DB
}

func NewPostgresWalletRepository(db *sql.DB) *PostgresWalletRepository {
	return &PostgresWalletRepository{db: db}
}

func (r *PostgresWalletRepository) Save(wallet *domain.WalletEntity) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `
		INSERT INTO wallets (id, owner_id, name, status, balance_amount, balance_currency)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (id) DO UPDATE SET
			name = EXCLUDED.name,
			status = EXCLUDED.status,
			balance_amount = EXCLUDED.balance_amount,
			balance_currency = EXCLUDED.balance_currency;
	`
	_, err = tx.Exec(query,
		wallet.ID(),
		wallet.OwnerID(),
		wallet.Name(),
		string(wallet.Status()),
		wallet.Balance().Amount().String(),
		wallet.Balance().Currency(),
	)
	if err != nil {
		return fmt.Errorf("failed to save wallet: %w", err)
	}

	for _, tr := range wallet.RecentTransactions() {
		trQuery := `
			INSERT INTO wallet_transactions (id, wallet_id, type, amount, currency, category_id, timestamp, description)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
			ON CONFLICT (id) DO NOTHING;
		`
		_, err = tx.Exec(trQuery,
			tr.ID(),
			wallet.ID(),
			string(tr.Type()),
			tr.Amount().Amount().String(),
			tr.Amount().Currency(),
			tr.CategoryID(),
			tr.Timestamp(),
			tr.Description(),
		)
		if err != nil {
			return fmt.Errorf("failed to save transaction: %w", err)
		}
	}

	return tx.Commit()
}

func (r *PostgresWalletRepository) FindByID(id string) (*domain.WalletEntity, error) {
	query := `SELECT owner_id, name, status, balance_amount, balance_currency FROM wallets WHERE id = $1`
	var ownerID, name, statusStr, balCurrency string
	var balAmountStr string

	err := r.db.QueryRow(query, id).Scan(&ownerID, &name, &statusStr, &balAmountStr, &balCurrency)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	balAmount, _ := decimal.NewFromString(balAmountStr)
	balance := shared.NewMoneyObject(balAmount, balCurrency)
	status := domain.WalletStatus(statusStr)

	wallet := domain.LoadWallet(id, ownerID, name, status, balance)
	return &wallet, nil
}

func (r *PostgresWalletRepository) FindByOwnerID(ownerId string) ([]*domain.WalletEntity, error) {
	query := `SELECT id, name, status, balance_amount, balance_currency FROM wallets WHERE owner_id = $1`
	rows, err := r.db.Query(query, ownerId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var wallets []*domain.WalletEntity
	for rows.Next() {
		var id, name, statusStr, balCurrency string
		var balAmountStr string

		if err := rows.Scan(&id, &name, &statusStr, &balAmountStr, &balCurrency); err != nil {
			return nil, err
		}

		balAmount, _ := decimal.NewFromString(balAmountStr)
		balance := shared.NewMoneyObject(balAmount, balCurrency)
		status := domain.WalletStatus(statusStr)

		wallet := domain.LoadWallet(id, ownerId, name, status, balance)
		wallets = append(wallets, &wallet)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return wallets, nil
}
