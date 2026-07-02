package domains

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type BuggetItemEntity struct {
	id              string
	categoryId      string
	allocatedAmount MoneyObject
	spentAmount     MoneyObject
}

func NewBuggetItem(categoryId string, allocatedAmount MoneyObject) BuggetItemEntity {
	return BuggetItemEntity{
		id:              uuid.New().String(),
		categoryId:      categoryId,
		allocatedAmount: allocatedAmount,
		spentAmount:     NewMoneyObject(decimal.NewFromInt(0), allocatedAmount.currency),
	}
}

func (b *BuggetItemEntity) addAllocation(amount MoneyObject) error {
	if b.allocatedAmount.currency != amount.currency {
		return fmt.Errorf("currency mismatch")
	}

	b.allocatedAmount.amount = b.allocatedAmount.amount.Add(amount.amount)
	return nil
}

func (b *BuggetItemEntity) reduceAllocation(amount MoneyObject) error {
	if b.allocatedAmount.currency != amount.currency {
		return fmt.Errorf("currency mismatch")
	}

	if b.allocatedAmount.amount.LessThan(amount.amount) {
		return fmt.Errorf("insufficient allocated amount")
	}

	b.allocatedAmount.amount = b.allocatedAmount.amount.Sub(amount.amount)
	return nil
}

func (b *BuggetItemEntity) recordExpense(amount MoneyObject) error {
	if b.allocatedAmount.currency != amount.currency {
		return fmt.Errorf("currency mismatch")
	}

	b.spentAmount.amount = b.spentAmount.amount.Add(amount.amount)
	return nil
}

func (b *BuggetItemEntity) isOverspent() bool {
	return b.spentAmount.amount.GreaterThan(b.allocatedAmount.amount)
}
