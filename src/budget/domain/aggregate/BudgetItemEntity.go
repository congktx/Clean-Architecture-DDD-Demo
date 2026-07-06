package aggregate

import (
	"fmt"

	shared "src/src/shared/domain"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type BudgetItemEntity struct {
	id              string
	categoryId      string
	allocatedAmount shared.MoneyObject
	spentAmount     shared.MoneyObject
}

func NewBudgetItem(categoryId string, allocatedAmount shared.MoneyObject) BudgetItemEntity {
	return BudgetItemEntity{
		id:              uuid.New().String(),
		categoryId:      categoryId,
		allocatedAmount: allocatedAmount,
		spentAmount:     shared.NewMoneyObject(decimal.NewFromInt(0), allocatedAmount.Currency()),
	}
}

func (b *BudgetItemEntity) addAllocation(amount shared.MoneyObject) error {
	if b.allocatedAmount.Currency() != amount.Currency() {
		return fmt.Errorf("currency mismatch")
	}

	var err error
	b.allocatedAmount, err = b.allocatedAmount.Add(amount)
	if err != nil {
		return err
	}

	return nil
}

func (b *BudgetItemEntity) reduceAllocation(amount shared.MoneyObject) error {
	if b.allocatedAmount.Currency() != amount.Currency() {
		return fmt.Errorf("currency mismatch")
	}

	if b.allocatedAmount.Amount().LessThan(amount.Amount()) {
		return fmt.Errorf("insufficient allocated amount")
	}

	var err error
	b.allocatedAmount, err = b.allocatedAmount.Sub(amount)
	if err != nil {
		return err
	}

	return nil
}

func (b *BudgetItemEntity) recordExpense(amount shared.MoneyObject) error {
	if b.allocatedAmount.Currency() != amount.Currency() {
		return fmt.Errorf("currency mismatch")
	}

	var err error
	b.spentAmount, err = b.spentAmount.Add(amount)
	if err != nil {
		return err
	}

	return nil
}

func (b *BudgetItemEntity) isOverspent() bool {
	isOver, _ := b.spentAmount.GreaterThan(b.allocatedAmount)
	return isOver
}
