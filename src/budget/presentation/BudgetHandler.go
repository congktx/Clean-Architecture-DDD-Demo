package presentation

import (
	"encoding/json"
	"net/http"
	"src/src/budget/usecase"

	"github.com/shopspring/decimal"
)

type BudgetHandler struct {
	usecase usecase.BudgetUsecase
}

func NewBudgetHandler(u usecase.BudgetUsecase) *BudgetHandler {
	return &BudgetHandler{usecase: u}
}

type CreateBudgetRequest struct {
	OwnerID    string `json:"owner_id"`
	Month      int    `json:"month"`
	Year       int    `json:"year"`
	TotalLimit string `json:"total_limit"`
	Currency   string `json:"currency"`
}

func (h *BudgetHandler) CreateBudget(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CreateBudgetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	limit, err := decimal.NewFromString(req.TotalLimit)
	if err != nil {
		http.Error(w, "Invalid total limit format", http.StatusBadRequest)
		return
	}

	err = h.usecase.CreateBudget(usecase.CreateBudgetRequest{
		OwnerID:    req.OwnerID,
		Month:      req.Month,
		Year:       req.Year,
		TotalLimit: limit,
		Currency:   req.Currency,
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"message": "Budget created successfully"}`))
}

type AllocateFundsRequest struct {
	BudgetID   string `json:"budget_id"`
	CategoryID string `json:"category_id"`
	Amount     string `json:"amount"`
	Currency   string `json:"currency"`
}

func (h *BudgetHandler) AllocateFunds(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req AllocateFundsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	amount, err := decimal.NewFromString(req.Amount)
	if err != nil {
		http.Error(w, "Invalid amount format", http.StatusBadRequest)
		return
	}

	err = h.usecase.AllocateFunds(usecase.AllocateFundsRequest{
		BudgetID:   req.BudgetID,
		CategoryID: req.CategoryID,
		Amount:     amount,
		Currency:   req.Currency,
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Funds allocated successfully"}`))
}

type RebalanceRequest struct {
	BudgetID       string `json:"budget_id"`
	FromCategoryID string `json:"from_category_id"`
	ToCategoryID   string `json:"to_category_id"`
	Amount         string `json:"amount"`
	Currency       string `json:"currency"`
}

func (h *BudgetHandler) Rebalance(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req RebalanceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	amount, err := decimal.NewFromString(req.Amount)
	if err != nil {
		http.Error(w, "Invalid amount format", http.StatusBadRequest)
		return
	}

	err = h.usecase.Rebalance(usecase.RebalanceRequest{
		BudgetID:       req.BudgetID,
		FromCategoryID: req.FromCategoryID,
		ToCategoryID:   req.ToCategoryID,
		Amount:         amount,
		Currency:       req.Currency,
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Budget rebalanced successfully"}`))
}

type RecordBudgetExpenseRequest struct {
	BudgetID   string `json:"budget_id"`
	CategoryID string `json:"category_id"`
	Amount     string `json:"amount"`
	Currency   string `json:"currency"`
}

func (h *BudgetHandler) RecordExpense(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req RecordBudgetExpenseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	amount, err := decimal.NewFromString(req.Amount)
	if err != nil {
		http.Error(w, "Invalid amount format", http.StatusBadRequest)
		return
	}

	err = h.usecase.RecordExpense(usecase.RecordBudgetExpenseRequest{
		BudgetID:   req.BudgetID,
		CategoryID: req.CategoryID,
		Amount:     amount,
		Currency:   req.Currency,
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Budget expense recorded successfully"}`))
}
