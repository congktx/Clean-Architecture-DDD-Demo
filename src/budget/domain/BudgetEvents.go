package domain

import (
	shared "src/src/shared/domain"
	"time"
)

type FundsAllocatedEvent struct {
	budgetId   string
	categoryId string
	amount     shared.MoneyObject
	occurredOn int64
}

func NewFundsAllocatedEvent(budgetId string, categoryId string, amount shared.MoneyObject) FundsAllocatedEvent {
	return FundsAllocatedEvent{
		budgetId:   budgetId,
		categoryId: categoryId,
		amount:     amount,
		occurredOn: time.Now().UnixNano(),
	}
}

func (e FundsAllocatedEvent) EventName() string {
	return "FundsAllocated"
}

func (e FundsAllocatedEvent) OccurredOn() int64 {
	return e.occurredOn
}

type BudgetOverspentEvent struct {
	budgetId   string
	categoryId string
	occurredOn int64
}

func NewBudgetOverspentEvent(budgetId string, categoryId string) BudgetOverspentEvent {
	return BudgetOverspentEvent{
		budgetId:   budgetId,
		categoryId: categoryId,
		occurredOn: time.Now().UnixNano(),
	}
}

func (e BudgetOverspentEvent) EventName() string {
	return "BudgetOverspent"
}

func (e BudgetOverspentEvent) OccurredOn() int64 {
	return e.occurredOn
}
