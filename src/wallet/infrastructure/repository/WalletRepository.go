package repository

import (
	domainWallet "src/src/wallet/domain"
	"src/src/wallet/infrastructure/database"
)

type walletRepository struct {
	db *database.DB
}

func NewWalletRepository(db *database.DB) domainWallet.WalletRepository {
	return &walletRepository{
		db: db,
	}
}

func (w *walletRepository) FindByID(id string) (*domainWallet.WalletEntity, error) {

	row := w.db.Conn.QueryRow(
		"SELECT id FROM users WHERE id=$1",
		id,
	)

	wallet := &domainWallet.WalletEntity{}

	err := row.Scan(
		&wallet.id,
	)

	if err != nil {
		return nil, err
	}

	return wallet, nil
}
