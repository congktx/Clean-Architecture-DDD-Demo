package domains

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type WalletEntity struct {
	id                 string
	ownerId            string
	name               string
	status             string
	balance            MoneyObject
	recentTransactions []WalletTransactionEntity
}

func NewWallet(ownerId string, name string, balance MoneyObject) WalletEntity {
	return WalletEntity{
		id:                 uuid.New().String(),
		ownerId:            ownerId,
		name:               name,
		status:             "ACTIVE",
		balance:            balance,
		recentTransactions: []WalletTransactionEntity{},
	}
}

func (w *WalletEntity) recordExpense(amount MoneyObject, categoryId string, timestamp int64, description string) error {
	if w.balance.currency != amount.currency {
		return fmt.Errorf("currency mismatch")
	}

	if w.status != "ACTIVE" {
		return fmt.Errorf("wallet is not active")
	}

	if amount.amount.LessThanOrEqual(decimal.NewFromInt(0)) {
		return fmt.Errorf("amount must be greater than zero")
	}

	if w.balance.amount.LessThan(amount.amount) {
		return fmt.Errorf("insufficient balance")
	}

	transaction := NewWalletTransaction(
		"EXPENSE",
		amount,
		categoryId,
		timestamp,
		description,
	)

	w.recentTransactions = append(w.recentTransactions, transaction)
	w.balance.amount = w.balance.amount.Sub(amount.amount)

	return nil
}

func (w *WalletEntity) recordIncome(amount MoneyObject, categoryId string, timestamp int64, description string) error {
	if w.balance.currency != amount.currency {
		return fmt.Errorf("currency mismatch")
	}

	if w.status != "ACTIVE" {
		return fmt.Errorf("wallet is not ACTIVE")
	}

	if amount.amount.LessThanOrEqual(decimal.NewFromInt(0)) {
		return fmt.Errorf("amount must be greater than zero")
	}

	transaction := NewWalletTransaction(
		"INCOME",
		amount,
		categoryId,
		timestamp,
		description,
	)

	w.recentTransactions = append(w.recentTransactions, transaction)
	w.balance.amount = w.balance.amount.Add(amount.amount)

	return nil
}

func (w *WalletEntity) clearRecentTransactions() {
	w.recentTransactions = []WalletTransactionEntity{}
}
