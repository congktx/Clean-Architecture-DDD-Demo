package usecase

import (
	"fmt"
	"github.com/shopspring/decimal"
	"src/src/budget/domain"
	shared "src/src/shared/domain"
)

type BudgetUsecase interface {
	CreateBudget(req CreateBudgetRequest) error
	AllocateFunds(req AllocateFundsRequest) error
}

type budgetUsecaseImpl struct {
	repo       domain.BudgetRepository
	dispatcher shared.EventDispatcher
}

func NewBudgetUsecase(repo domain.BudgetRepository, dispatcher shared.EventDispatcher) BudgetUsecase {
	return &budgetUsecaseImpl{
		repo:       repo,
		dispatcher: dispatcher,
	}
}

type CreateBudgetRequest struct {
	OwnerID    string
	Month      int
	Year       int
	TotalLimit decimal.Decimal
	Currency   string
}

func (u *budgetUsecaseImpl) CreateBudget(req CreateBudgetRequest) error {
	period, err := shared.NewMonthYear(req.Month, req.Year)
	if err != nil {
		return fmt.Errorf("invalid period: %w", err)
	}

	totalLimit := shared.NewMoneyObject(req.TotalLimit, req.Currency)
	budget := domain.NewBudget(req.OwnerID, period, totalLimit)

	err = u.repo.Save(&budget)
	if err != nil {
		return fmt.Errorf("failed to save budget: %w", err)
	}

	events := budget.GetDomainEvents()
	if len(events) > 0 {
		if err := u.dispatcher.Dispatch(events); err != nil {
			return fmt.Errorf("failed to dispatch events: %w", err)
		}
		budget.ClearDomainEvents()
	}

	return nil
}

type AllocateFundsRequest struct {
	BudgetID   string
	CategoryID string
	Amount     decimal.Decimal
	Currency   string
}

func (u *budgetUsecaseImpl) AllocateFunds(req AllocateFundsRequest) error {
	budget, err := u.repo.FindByID(req.BudgetID)
	if err != nil {
		return fmt.Errorf("failed to find budget: %w", err)
	}
	if budget == nil {
		return fmt.Errorf("budget not found")
	}

	amount := shared.NewMoneyObject(req.Amount, req.Currency)
	err = budget.AllocateFunds(req.CategoryID, amount)
	if err != nil {
		return fmt.Errorf("failed to allocate funds: %w", err)
	}

	err = u.repo.Save(budget)
	if err != nil {
		return fmt.Errorf("failed to save budget: %w", err)
	}

	events := budget.GetDomainEvents()
	if len(events) > 0 {
		if err := u.dispatcher.Dispatch(events); err != nil {
			return fmt.Errorf("failed to dispatch events: %w", err)
		}
		budget.ClearDomainEvents()
	}

	return nil
}
