package domain

import (
	shared "src/src/shared/domain"
	"time"
)

type ExpenseRecordedEvent struct {
	walletId   string
	amount     shared.MoneyObject
	categoryId string
	occurredOn int64
}

func NewExpenseRecordedEvent(walletId string, amount shared.MoneyObject, categoryId string) ExpenseRecordedEvent {
	return ExpenseRecordedEvent{
		walletId:   walletId,
		amount:     amount,
		categoryId: categoryId,
		occurredOn: time.Now().UnixNano(),
	}
}

func (e ExpenseRecordedEvent) EventName() string {
	return "ExpenseRecorded"
}

func (e ExpenseRecordedEvent) OccurredOn() int64 {
	return e.occurredOn
}

type IncomeRecordedEvent struct {
	walletId   string
	amount     shared.MoneyObject
	categoryId string
	occurredOn int64
}

func NewIncomeRecordedEvent(walletId string, amount shared.MoneyObject, categoryId string) IncomeRecordedEvent {
	return IncomeRecordedEvent{
		walletId:   walletId,
		amount:     amount,
		categoryId: categoryId,
		occurredOn: time.Now().UnixNano(),
	}
}

func (e IncomeRecordedEvent) EventName() string {
	return "IncomeRecorded"
}

func (e IncomeRecordedEvent) OccurredOn() int64 {
	return e.occurredOn
}
