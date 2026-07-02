package domains

import (
	"fmt"

	"github.com/google/uuid"
)

type BuggetEntity struct {
	id          string
	ownerId     string
	period      MonthYearObject
	totalLimit  MoneyObject
	buggetItems []BuggetItemEntity
}

func NewBugget(ownerId string, period MonthYearObject, totalLimit MoneyObject) BuggetEntity {
	return BuggetEntity{
		id:          uuid.New().String(),
		ownerId:     ownerId,
		period:      period,
		totalLimit:  totalLimit,
		buggetItems: []BuggetItemEntity{},
	}
}

func (b *BuggetEntity) allocateFunds(categoryId string, amount MoneyObject) error {
	if b.totalLimit.currency != amount.currency {
		return fmt.Errorf("currency mismatch")
	}

	if b.totalLimit.amount.LessThan(amount.amount) {
		return fmt.Errorf("insufficient total limit")
	}

	for i := range b.buggetItems {
		if b.buggetItems[i].categoryId == categoryId {
			return b.buggetItems[i].addAllocation(amount)
		}
	}

	newBuggetItem := NewBuggetItem(categoryId, amount)
	b.buggetItems = append(b.buggetItems, newBuggetItem)

	return nil
}

func (b *BuggetEntity) rebalance(fromCategoryId string, toCategoryId string, amount MoneyObject) error {
	if b.totalLimit.currency != amount.currency {
		return fmt.Errorf("currency mismatch")
	}

	var fromItem *BuggetItemEntity
	var toItem *BuggetItemEntity

	for i := range b.buggetItems {
		if b.buggetItems[i].categoryId == fromCategoryId {
			fromItem = &b.buggetItems[i]
		}
		if b.buggetItems[i].categoryId == toCategoryId {
			toItem = &b.buggetItems[i]
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

func (b *BuggetEntity) recordExpense(categoryId string, amount MoneyObject) error {
	if b.totalLimit.currency != amount.currency {
		return fmt.Errorf("currency mismatch")
	}

	var buggetItem *BuggetItemEntity

	for i := range b.buggetItems {
		if b.buggetItems[i].categoryId == categoryId {
			buggetItem = &b.buggetItems[i]
			break
		}
	}

	if buggetItem != nil {
		if err := buggetItem.recordExpense(amount); err != nil {
			return err
		}
		if buggetItem.isOverspent() {
			return fmt.Errorf("category %s is overspent", categoryId)
		}
		return nil
	}

	return fmt.Errorf("category not found")
}
