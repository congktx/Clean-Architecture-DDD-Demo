package domain

import (
	"fmt"

	shared "src/src/shared/domain"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type WalletStatus string

const (
	WalletStatusActive   WalletStatus = "ACTIVE"
	WalletStatusInactive WalletStatus = "INACTIVE"
)

type WalletEntity struct {
	shared.AggregateRoot
	id                 string
	ownerId            string
	name               string
	status             WalletStatus
	balance            shared.MoneyObject
	recentTransactions []WalletTransactionEntity
}

func NewWallet(ownerId string, name string, balance shared.MoneyObject) WalletEntity {
	return WalletEntity{
		id:                 uuid.New().String(),
		ownerId:            ownerId,
		name:               name,
		status:             WalletStatusActive,
		balance:            balance,
		recentTransactions: []WalletTransactionEntity{},
	}
}

func (w *WalletEntity) RecordExpense(amount shared.MoneyObject, categoryId string, timestamp int64, description string) error {
	if w.balance.Currency() != amount.Currency() {
		return fmt.Errorf("currency mismatch")
	}

	if w.status != WalletStatusActive {
		return fmt.Errorf("wallet is not active")
	}

	if amount.Amount().LessThanOrEqual(decimal.NewFromInt(0)) {
		return fmt.Errorf("amount must be greater than zero")
	}

	if w.balance.Amount().LessThan(amount.Amount()) {
		return fmt.Errorf("insufficient balance")
	}

	transaction := NewWalletTransaction(
		TransactionTypeExpense,
		amount,
		categoryId,
		timestamp,
		description,
	)

	w.recentTransactions = append(w.recentTransactions, transaction)

	var err error
	w.balance, err = w.balance.Sub(amount)
	if err != nil {
		return err
	}

	w.AddDomainEvent(NewExpenseRecordedEvent(w.id, amount, categoryId))

	return nil
}

func (w *WalletEntity) RecordIncome(amount shared.MoneyObject, categoryId string, timestamp int64, description string) error {
	if w.balance.Currency() != amount.Currency() {
		return fmt.Errorf("currency mismatch")
	}

	if w.status != WalletStatusActive {
		return fmt.Errorf("wallet is not ACTIVE")
	}

	if amount.Amount().LessThanOrEqual(decimal.NewFromInt(0)) {
		return fmt.Errorf("amount must be greater than zero")
	}

	transaction := NewWalletTransaction(
		TransactionTypeIncome,
		amount,
		categoryId,
		timestamp,
		description,
	)

	w.recentTransactions = append(w.recentTransactions, transaction)

	var err error
	w.balance, err = w.balance.Add(amount)
	if err != nil {
		return err
	}

	w.AddDomainEvent(NewIncomeRecordedEvent(w.id, amount, categoryId))

	return nil
}

func (w *WalletEntity) ClearRecentTransactions() {
	w.recentTransactions = []WalletTransactionEntity{}
}

func (w *WalletEntity) ID() string {
	return w.id
}

func (w *WalletEntity) OwnerID() string {
	return w.ownerId
}

func (w *WalletEntity) Name() string {
	return w.name
}

func (w *WalletEntity) Status() WalletStatus {
	return w.status
}

func (w *WalletEntity) Balance() shared.MoneyObject {
	return w.balance
}

func (w *WalletEntity) RecentTransactions() []WalletTransactionEntity {
	return w.recentTransactions
}

func LoadWallet(id string, ownerId string, name string, status WalletStatus, balance shared.MoneyObject) WalletEntity {
	return WalletEntity{
		id:                 id,
		ownerId:            ownerId,
		name:               name,
		status:             status,
		balance:            balance,
		recentTransactions: []WalletTransactionEntity{},
	}
}
