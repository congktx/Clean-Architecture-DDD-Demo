package domain

import (
	"fmt"

	"github.com/shopspring/decimal"
)

type MoneyObject struct {
	amount   decimal.Decimal
	currency string
}

func NewMoneyObject(amount decimal.Decimal, currency string) MoneyObject {
	return MoneyObject{
		amount:   amount,
		currency: currency,
	}
}

func (m MoneyObject) Currency() string {
	return m.currency
}

func (m MoneyObject) Amount() decimal.Decimal {
	return m.amount
}

func (m MoneyObject) Add(other MoneyObject) (MoneyObject, error) {
	if m.currency != other.currency {
		return MoneyObject{}, fmt.Errorf("currency mismatch")
	}

	return MoneyObject{
		amount:   m.amount.Add(other.amount),
		currency: m.currency,
	}, nil
}

func (m MoneyObject) Sub(other MoneyObject) (MoneyObject, error) {
	if m.currency != other.currency {
		return MoneyObject{}, fmt.Errorf("currency mismatch")
	}

	return MoneyObject{
		amount:   m.amount.Sub(other.amount),
		currency: m.currency,
	}, nil
}
func (m MoneyObject) LessThan(other MoneyObject) (bool, error) {
	if m.currency != other.currency {
		return false, fmt.Errorf("currency mismatch")
	}

	return m.amount.LessThan(other.amount), nil
}

func (m MoneyObject) GreaterThan(other MoneyObject) (bool, error) {
	if m.currency != other.currency {
		return false, fmt.Errorf("currency mismatch")
	}

	return m.amount.GreaterThan(other.amount), nil
}

func (m MoneyObject) Equal(other MoneyObject) (bool, error) {
	if m.currency != other.currency {
		return false, fmt.Errorf("currency mismatch")
	}

	return m.amount.Equal(other.amount), nil
}
func (m MoneyObject) LessThanOrEqual(other MoneyObject) (bool, error) {
	if m.currency != other.currency {
		return false, fmt.Errorf("currency mismatch")
	}

	return m.amount.LessThanOrEqual(other.amount), nil
}

func (m MoneyObject) GreaterThanOrEqual(other MoneyObject) (bool, error) {
	if m.currency != other.currency {
		return false, fmt.Errorf("currency mismatch")
	}

	return m.amount.GreaterThanOrEqual(other.amount), nil
}
