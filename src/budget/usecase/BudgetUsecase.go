package usecase

import (
	"fmt"
	"github.com/shopspring/decimal"
	"src/src/budget/domain"
	shared "src/src/shared/domain"
)

type BudgetUsecase interface {
	CreateBudget(req CreateBudgetRequest) (*BudgetResponse, error)
	AllocateFunds(req AllocateFundsRequest) (*BudgetResponse, error)
	Rebalance(req RebalanceRequest) (*BudgetResponse, error)
	RecordExpense(req RecordBudgetExpenseRequest) (*BudgetResponse, error)
	GetBudgetByID(id string) (*BudgetResponse, error)
	GetBudgetsByOwnerID(ownerId string) ([]*BudgetResponse, error)
}

type BudgetItemResponse struct {
	ID              string `json:"id"`
	CategoryID      string `json:"category_id"`
	AllocatedAmount string `json:"allocated_amount"`
	SpentAmount     string `json:"spent_amount"`
	Currency        string `json:"currency"`
}

type BudgetResponse struct {
	ID         string               `json:"id"`
	OwnerID    string               `json:"owner_id"`
	Month      int                  `json:"month"`
	Year       int                  `json:"year"`
	TotalLimit string               `json:"total_limit"`
	Currency   string               `json:"currency"`
	Items      []BudgetItemResponse `json:"items"`
}

func mapBudgetToResponse(b *domain.BudgetEntity) *BudgetResponse {
	items := make([]BudgetItemResponse, 0, len(b.BudgetItems()))
	for _, item := range b.BudgetItems() {
		items = append(items, BudgetItemResponse{
			ID:              item.ID(),
			CategoryID:      item.CategoryID(),
			AllocatedAmount: item.AllocatedAmount().Amount().String(),
			SpentAmount:     item.SpentAmount().Amount().String(),
			Currency:        item.AllocatedAmount().Currency(),
		})
	}
	return &BudgetResponse{
		ID:         b.ID(),
		OwnerID:    b.OwnerID(),
		Month:      b.Period().Month(),
		Year:       b.Period().Year(),
		TotalLimit: b.TotalLimit().Amount().String(),
		Currency:   b.TotalLimit().Currency(),
		Items:      items,
	}
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

func (u *budgetUsecaseImpl) CreateBudget(req CreateBudgetRequest) (*BudgetResponse, error) {
	period, err := shared.NewMonthYear(req.Month, req.Year)
	if err != nil {
		return nil, fmt.Errorf("invalid period: %w", err)
	}

	totalLimit := shared.NewMoneyObject(req.TotalLimit, req.Currency)
	budget := domain.NewBudget(req.OwnerID, period, totalLimit)

	err = u.repo.Save(&budget)
	if err != nil {
		return nil, fmt.Errorf("failed to save budget: %w", err)
	}

	events := budget.GetDomainEvents()
	if len(events) > 0 {
		if err := u.dispatcher.Dispatch(events); err != nil {
			return nil, fmt.Errorf("failed to dispatch events: %w", err)
		}
		budget.ClearDomainEvents()
	}

	return mapBudgetToResponse(&budget), nil
}

type AllocateFundsRequest struct {
	BudgetID   string
	CategoryID string
	Amount     decimal.Decimal
	Currency   string
}

func (u *budgetUsecaseImpl) AllocateFunds(req AllocateFundsRequest) (*BudgetResponse, error) {
	budget, err := u.repo.FindByID(req.BudgetID)
	if err != nil {
		return nil, fmt.Errorf("failed to find budget: %w", err)
	}
	if budget == nil {
		return nil, fmt.Errorf("budget not found")
	}

	amount := shared.NewMoneyObject(req.Amount, req.Currency)
	err = budget.AllocateFunds(req.CategoryID, amount)
	if err != nil {
		return nil, fmt.Errorf("failed to allocate funds: %w", err)
	}

	err = u.repo.Save(budget)
	if err != nil {
		return nil, fmt.Errorf("failed to save budget: %w", err)
	}

	events := budget.GetDomainEvents()
	if len(events) > 0 {
		if err := u.dispatcher.Dispatch(events); err != nil {
			return nil, fmt.Errorf("failed to dispatch events: %w", err)
		}
		budget.ClearDomainEvents()
	}

	return mapBudgetToResponse(budget), nil
}

type RebalanceRequest struct {
	BudgetID       string
	FromCategoryID string
	ToCategoryID   string
	Amount         decimal.Decimal
	Currency       string
}

func (u *budgetUsecaseImpl) Rebalance(req RebalanceRequest) (*BudgetResponse, error) {
	budget, err := u.repo.FindByID(req.BudgetID)
	if err != nil {
		return nil, fmt.Errorf("failed to find budget: %w", err)
	}
	if budget == nil {
		return nil, fmt.Errorf("budget not found")
	}

	amount := shared.NewMoneyObject(req.Amount, req.Currency)
	err = budget.Rebalance(req.FromCategoryID, req.ToCategoryID, amount)
	if err != nil {
		return nil, fmt.Errorf("failed to rebalance: %w", err)
	}

	err = u.repo.Save(budget)
	if err != nil {
		return nil, fmt.Errorf("failed to save budget: %w", err)
	}

	events := budget.GetDomainEvents()
	if len(events) > 0 {
		if err := u.dispatcher.Dispatch(events); err != nil {
			return nil, fmt.Errorf("failed to dispatch events: %w", err)
		}
		budget.ClearDomainEvents()
	}

	return mapBudgetToResponse(budget), nil
}

type RecordBudgetExpenseRequest struct {
	BudgetID   string
	CategoryID string
	Amount     decimal.Decimal
	Currency   string
}

func (u *budgetUsecaseImpl) RecordExpense(req RecordBudgetExpenseRequest) (*BudgetResponse, error) {
	budget, err := u.repo.FindByID(req.BudgetID)
	if err != nil {
		return nil, fmt.Errorf("failed to find budget: %w", err)
	}
	if budget == nil {
		return nil, fmt.Errorf("budget not found")
	}

	amount := shared.NewMoneyObject(req.Amount, req.Currency)
	err = budget.RecordExpense(req.CategoryID, amount)
	if err != nil {
		return nil, fmt.Errorf("failed to record expense: %w", err)
	}

	err = u.repo.Save(budget)
	if err != nil {
		return nil, fmt.Errorf("failed to save budget: %w", err)
	}

	events := budget.GetDomainEvents()
	if len(events) > 0 {
		if err := u.dispatcher.Dispatch(events); err != nil {
			return nil, fmt.Errorf("failed to dispatch events: %w", err)
		}
		budget.ClearDomainEvents()
	}

	return mapBudgetToResponse(budget), nil
}

func (u *budgetUsecaseImpl) GetBudgetByID(id string) (*BudgetResponse, error) {
	budget, err := u.repo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get budget: %w", err)
	}
	if budget == nil {
		return nil, fmt.Errorf("budget not found")
	}
	return mapBudgetToResponse(budget), nil
}

func (u *budgetUsecaseImpl) GetBudgetsByOwnerID(ownerId string) ([]*BudgetResponse, error) {
	budgets, err := u.repo.FindByOwnerID(ownerId)
	if err != nil {
		return nil, fmt.Errorf("failed to get budgets: %w", err)
	}
	var res []*BudgetResponse
	for _, b := range budgets {
		res = append(res, mapBudgetToResponse(b))
	}
	return res, nil
}
