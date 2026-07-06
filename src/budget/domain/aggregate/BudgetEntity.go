package aggregate

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

func (b *BudgetEntity) allocateFunds(categoryId string, amount shared.MoneyObject) error {
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
			err = b.budgetItems[i].addAllocation(amount)
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

func (b *BudgetEntity) rebalance(fromCategoryId string, toCategoryId string, amount shared.MoneyObject) error {
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

	if err := fromItem.reduceAllocation(amount); err != nil {
		return err
	}

	return toItem.addAllocation(amount)
}

func (b *BudgetEntity) recordExpense(categoryId string, amount shared.MoneyObject) error {
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
		if err := buggetItem.recordExpense(amount); err != nil {
			return err
		}
		if buggetItem.isOverspent() {
			b.AddDomainEvent(NewBudgetOverspentEvent(b.id, categoryId))
			return fmt.Errorf("category %s is overspent", categoryId)
		}
		return nil
	}

	return fmt.Errorf("category not found")
}
