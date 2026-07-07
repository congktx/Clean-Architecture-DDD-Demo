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

func writeJSONResponse(w http.ResponseWriter, status int, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": message,
		"data":    data,
	})
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

	res, err := h.usecase.CreateBudget(usecase.CreateBudgetRequest{
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

	writeJSONResponse(w, http.StatusCreated, "Budget created successfully", res)
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

	res, err := h.usecase.AllocateFunds(usecase.AllocateFundsRequest{
		BudgetID:   req.BudgetID,
		CategoryID: req.CategoryID,
		Amount:     amount,
		Currency:   req.Currency,
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSONResponse(w, http.StatusOK, "Funds allocated successfully", res)
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

	res, err := h.usecase.Rebalance(usecase.RebalanceRequest{
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

	writeJSONResponse(w, http.StatusOK, "Budget rebalanced successfully", res)
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

	res, err := h.usecase.RecordExpense(usecase.RecordBudgetExpenseRequest{
		BudgetID:   req.BudgetID,
		CategoryID: req.CategoryID,
		Amount:     amount,
		Currency:   req.Currency,
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSONResponse(w, http.StatusOK, "Budget expense recorded successfully", res)
}

func (h *BudgetHandler) GetBudget(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Budget ID is required", http.StatusBadRequest)
		return
	}

	res, err := h.usecase.GetBudgetByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSONResponse(w, http.StatusOK, "Budget retrieved successfully", res)
}

func (h *BudgetHandler) GetBudgets(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	ownerID := r.URL.Query().Get("owner_id")
	if ownerID == "" {
		http.Error(w, "owner_id query parameter is required", http.StatusBadRequest)
		return
	}

	res, err := h.usecase.GetBudgetsByOwnerID(ownerID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSONResponse(w, http.StatusOK, "Budgets retrieved successfully", res)
}
