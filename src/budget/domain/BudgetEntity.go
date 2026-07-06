package domain

import (
	"fmt"

	shared "src/src/shared/domain"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type BudgetEntity struct {
	shared.AggregateRoot
	id          string
	ownerId     string
	period      shared.MonthYearObject
	totalLimit  shared.MoneyObject
	budgetItems []BudgetItemEntity
}

func NewBudget(ownerId string, period shared.MonthYearObject, totalLimit shared.MoneyObject) BudgetEntity {
	return BudgetEntity{
		id:          uuid.New().String(),
		ownerId:     ownerId,
		period:      period,
		totalLimit:  totalLimit,
		budgetItems: []BudgetItemEntity{},
	}
}

func (b *BudgetEntity) AllocateFunds(categoryId string, amount shared.MoneyObject) error {
	if b.totalLimit.Currency() != amount.Currency() {
		return fmt.Errorf("currency mismatch")
	}

	totalAllocated := shared.NewMoneyObject(decimal.NewFromInt(0), amount.Currency())
	for i := range b.budgetItems {
		var err error
		totalAllocated, err = totalAllocated.Add(b.budgetItems[i].allocatedAmount)
		if err != nil {
			return err
		}
	}

	newTotal, err := totalAllocated.Add(amount)
	if err != nil {
		return err
	}

	isOver, err := newTotal.GreaterThan(b.totalLimit)
	if err != nil {
		return err
	}

	if isOver {
		return fmt.Errorf("insufficient total limit")
	}

	for i := range b.budgetItems {
		if b.budgetItems[i].categoryId == categoryId {
			err = b.budgetItems[i].AddAllocation(amount)
			if err == nil {
				b.AddDomainEvent(NewFundsAllocatedEvent(b.id, categoryId, amount))
			}
			return err
		}
	}

	newBudgetItem := NewBudgetItem(categoryId, amount)
	b.budgetItems = append(b.budgetItems, newBudgetItem)
	b.AddDomainEvent(NewFundsAllocatedEvent(b.id, categoryId, amount))

	return nil
}

func (b *BudgetEntity) Rebalance(fromCategoryId string, toCategoryId string, amount shared.MoneyObject) error {
	if b.totalLimit.Currency() != amount.Currency() {
		return fmt.Errorf("currency mismatch")
	}

	var fromItem *BudgetItemEntity
	var toItem *BudgetItemEntity

	for i := range b.budgetItems {
		if b.budgetItems[i].categoryId == fromCategoryId {
			fromItem = &b.budgetItems[i]
		}
		if b.budgetItems[i].categoryId == toCategoryId {
			toItem = &b.budgetItems[i]
		}
	}

	if fromItem == nil || toItem == nil {
		return fmt.Errorf("one or both categories not found")
	}

	if err := fromItem.ReduceAllocation(amount); err != nil {
		return err
	}

	return toItem.AddAllocation(amount)
}

func (b *BudgetEntity) RecordExpense(categoryId string, amount shared.MoneyObject) error {
	if b.totalLimit.Currency() != amount.Currency() {
		return fmt.Errorf("currency mismatch")
	}

	var buggetItem *BudgetItemEntity

	for i := range b.budgetItems {
		if b.budgetItems[i].categoryId == categoryId {
			buggetItem = &b.budgetItems[i]
			break
		}
	}

	if buggetItem != nil {
		if err := buggetItem.RecordExpense(amount); err != nil {
			return err
		}
		if buggetItem.IsOverspent() {
			b.AddDomainEvent(NewBudgetOverspentEvent(b.id, categoryId))
			return fmt.Errorf("category %s is overspent", categoryId)
		}
		return nil
	}

	return fmt.Errorf("category not found")
}

func LoadBudget(id, ownerId string, period shared.MonthYearObject, totalLimit shared.MoneyObject, items []BudgetItemEntity) BudgetEntity {
	return BudgetEntity{
		id:          id,
		ownerId:     ownerId,
		period:      period,
		totalLimit:  totalLimit,
		budgetItems: items,
	}
}

func (b *BudgetEntity) ID() string {
	return b.id
}

func (b *BudgetEntity) OwnerID() string {
	return b.ownerId
}

func (b *BudgetEntity) Period() shared.MonthYearObject {
	return b.period
}

func (b *BudgetEntity) TotalLimit() shared.MoneyObject {
	return b.totalLimit
}

func (b *BudgetEntity) BudgetItems() []BudgetItemEntity {
	return b.budgetItems
}
