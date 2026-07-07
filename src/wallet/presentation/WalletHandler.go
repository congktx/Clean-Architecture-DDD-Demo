package presentation

import (
	"encoding/json"
	"net/http"
	"src/src/wallet/usecase"
	"time"

	"github.com/shopspring/decimal"
)

type WalletHandler struct {
	walletUsecase usecase.WalletUsecase
}

func NewWalletHandler(u usecase.WalletUsecase) *WalletHandler {
	return &WalletHandler{walletUsecase: u}
}

func writeJSONResponse(w http.ResponseWriter, status int, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": message,
		"data":    data,
	})
}

type CreateWalletRequest struct {
	OwnerID  string `json:"owner_id"`
	Name     string `json:"name"`
	Currency string `json:"currency"`
}

func (h *WalletHandler) CreateWallet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CreateWalletRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	res, err := h.walletUsecase.CreateWallet(usecase.CreateWalletRequest{
		OwnerID:  req.OwnerID,
		Name:     req.Name,
		Currency: req.Currency,
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSONResponse(w, http.StatusCreated, "Wallet created successfully", res)
}

type RecordExpenseRequest struct {
	WalletID    string `json:"wallet_id"`
	CategoryID  string `json:"category_id"`
	Amount      string `json:"amount"`
	Currency    string `json:"currency"`
	Description string `json:"description"`
}

func (h *WalletHandler) RecordExpense(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req RecordExpenseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	amount, err := decimal.NewFromString(req.Amount)
	if err != nil {
		http.Error(w, "Invalid amount format", http.StatusBadRequest)
		return
	}

	res, err := h.walletUsecase.RecordExpense(usecase.RecordExpenseRequest{
		WalletID:    req.WalletID,
		CategoryID:  req.CategoryID,
		Amount:      amount,
		Currency:    req.Currency,
		Description: req.Description,
		Timestamp:   time.Now().UnixNano(),
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSONResponse(w, http.StatusOK, "Expense recorded successfully", res)
}

type RecordIncomeRequest struct {
	WalletID    string `json:"wallet_id"`
	CategoryID  string `json:"category_id"`
	Amount      string `json:"amount"`
	Currency    string `json:"currency"`
	Description string `json:"description"`
}

func (h *WalletHandler) RecordIncome(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req RecordIncomeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	amount, err := decimal.NewFromString(req.Amount)
	if err != nil {
		http.Error(w, "Invalid amount format", http.StatusBadRequest)
		return
	}

	res, err := h.walletUsecase.RecordIncome(usecase.RecordIncomeRequest{
		WalletID:    req.WalletID,
		CategoryID:  req.CategoryID,
		Amount:      amount,
		Currency:    req.Currency,
		Description: req.Description,
		Timestamp:   time.Now().UnixNano(),
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSONResponse(w, http.StatusOK, "Income recorded successfully", res)
}

func (h *WalletHandler) GetWallet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Wallet ID is required", http.StatusBadRequest)
		return
	}

	res, err := h.walletUsecase.GetWalletByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSONResponse(w, http.StatusOK, "Wallet retrieved successfully", res)
}

func (h *WalletHandler) GetWallets(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	ownerID := r.URL.Query().Get("owner_id")
	if ownerID == "" {
		http.Error(w, "owner_id query parameter is required", http.StatusBadRequest)
		return
	}

	res, err := h.walletUsecase.GetWalletsByOwnerID(ownerID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSONResponse(w, http.StatusOK, "Wallets retrieved successfully", res)
}
