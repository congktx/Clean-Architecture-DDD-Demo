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
