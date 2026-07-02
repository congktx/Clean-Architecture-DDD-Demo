package domains

import (
	"github.com/google/uuid"
)

type WalletTransactionEntity struct {
	id          string
	typ         string
	amount      MoneyObject
	categoryId  string
	timestamp   int64
	description string
}

func NewWalletTransaction(typ string, amount MoneyObject, categoryId string, timestamp int64, description string) WalletTransactionEntity {
	return WalletTransactionEntity{
		id:          uuid.New().String(),
		typ:         typ,
		amount:      amount,
		categoryId:  categoryId,
		timestamp:   timestamp,
		description: description,
	}
}
