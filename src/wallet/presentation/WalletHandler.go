package presentation

import (
	"encoding/json"
	"net/http"
	"github.com/shopspring/decimal"
	"src/src/wallet/usecase"
	"time"
)

type WalletHandler struct {
	walletUsecase usecase.WalletUsecase
}

func NewWalletHandler(u usecase.WalletUsecase) *WalletHandler {
	return &WalletHandler{walletUsecase: u}
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

	err := h.walletUsecase.CreateWallet(usecase.CreateWalletRequest{
		OwnerID:  req.OwnerID,
		Name:     req.Name,
		Currency: req.Currency,
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"message": "Wallet created successfully"}`))
}

type RecordExpenseRequest struct {
	WalletID    string `json:"wallet_id"`
	CategoryID  string `json:"category_id"`
	Amount      string `json:"amount"` // String to maintain decimal precision
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

	err = h.walletUsecase.RecordExpense(usecase.RecordExpenseRequest{
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

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Expense recorded successfully"}`))
}
