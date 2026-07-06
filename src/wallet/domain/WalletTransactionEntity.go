package domain

import (
	shared "src/src/shared/domain"

	"github.com/google/uuid"
)

type TransactionType string

const (
	TransactionTypeExpense TransactionType = "EXPENSE"
	TransactionTypeIncome  TransactionType = "INCOME"
)

type WalletTransactionEntity struct {
	id          string
	typ         TransactionType
	amount      shared.MoneyObject
	categoryId  string
	timestamp   int64
	description string
}

func NewWalletTransaction(typ TransactionType, amount shared.MoneyObject, categoryId string, timestamp int64, description string) WalletTransactionEntity {
	return WalletTransactionEntity{
		id:          uuid.New().String(),
		typ:         typ,
		amount:      amount,
		categoryId:  categoryId,
		timestamp:   timestamp,
		description: description,
	}
}

func (t *WalletTransactionEntity) ID() string {
	return t.id
}

func (t *WalletTransactionEntity) Type() TransactionType {
	return t.typ
}

func (t *WalletTransactionEntity) Amount() shared.MoneyObject {
	return t.amount
}

func (t *WalletTransactionEntity) CategoryID() string {
	return t.categoryId
}

func (t *WalletTransactionEntity) Timestamp() int64 {
	return t.timestamp
}

func (t *WalletTransactionEntity) Description() string {
	return t.description
}
