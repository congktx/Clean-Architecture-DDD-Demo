package usecase

import (
	"fmt"
	"github.com/shopspring/decimal"
	shared "src/src/shared/domain"
	"src/src/wallet/domain"
)

type WalletUsecase interface {
	CreateWallet(req CreateWalletRequest) (*WalletResponse, error)
	RecordExpense(req RecordExpenseRequest) (*WalletResponse, error)
	RecordIncome(req RecordIncomeRequest) (*WalletResponse, error)
	GetWalletByID(id string) (*WalletResponse, error)
	GetWalletsByOwnerID(ownerId string) ([]*WalletResponse, error)
}

type WalletResponse struct {
	ID       string `json:"id"`
	OwnerID  string `json:"owner_id"`
	Name     string `json:"name"`
	Balance  string `json:"balance"`
	Currency string `json:"currency"`
	Status   string `json:"status"`
}

func mapWalletToResponse(w *domain.WalletEntity) *WalletResponse {
	return &WalletResponse{
		ID:       w.ID(),
		OwnerID:  w.OwnerID(),
		Name:     w.Name(),
		Balance:  w.Balance().Amount().String(),
		Currency: w.Balance().Currency(),
		Status:   string(w.Status()),
	}
}

type walletUsecaseImpl struct {
	walletRepo domain.WalletRepository
	dispatcher shared.EventDispatcher
}

func NewWalletUsecase(repo domain.WalletRepository, dispatcher shared.EventDispatcher) WalletUsecase {
	return &walletUsecaseImpl{
		walletRepo: repo,
		dispatcher: dispatcher,
	}
}

type CreateWalletRequest struct {
	OwnerID  string
	Name     string
	Currency string
}

func (u *walletUsecaseImpl) CreateWallet(req CreateWalletRequest) (*WalletResponse, error) {
	balance := shared.NewMoneyObject(decimal.NewFromInt(0), req.Currency)
	wallet := domain.NewWallet(req.OwnerID, req.Name, balance)

	err := u.walletRepo.Save(&wallet)
	if err != nil {
		return nil, err
	}

	events := wallet.GetDomainEvents()
	if len(events) > 0 {
		err = u.dispatcher.Dispatch(events)
		if err != nil {
			return nil, fmt.Errorf("failed to dispatch events: %w", err)
		}
		wallet.ClearDomainEvents()
	}

	return mapWalletToResponse(&wallet), nil
}

type RecordExpenseRequest struct {
	WalletID    string
	CategoryID  string
	Amount      decimal.Decimal
	Currency    string
	Description string
	Timestamp   int64
}

func (u *walletUsecaseImpl) RecordExpense(req RecordExpenseRequest) (*WalletResponse, error) {
	wallet, err := u.walletRepo.FindByID(req.WalletID)
	if err != nil {
		return nil, fmt.Errorf("wallet not found: %w", err)
	}
	if wallet == nil {
		return nil, fmt.Errorf("wallet not found")
	}

	money := shared.NewMoneyObject(req.Amount, req.Currency)
	err = wallet.RecordExpense(money, req.CategoryID, req.Timestamp, req.Description)
	if err != nil {
		return nil, fmt.Errorf("failed to record expense: %w", err)
	}

	err = u.walletRepo.Save(wallet)
	if err != nil {
		return nil, fmt.Errorf("failed to save wallet: %w", err)
	}

	events := wallet.GetDomainEvents()
	if len(events) > 0 {
		err = u.dispatcher.Dispatch(events)
		if err != nil {
			return nil, fmt.Errorf("failed to dispatch events: %w", err)
		}
		wallet.ClearDomainEvents()
	}

	return mapWalletToResponse(wallet), nil
}

type RecordIncomeRequest struct {
	WalletID    string
	CategoryID  string
	Amount      decimal.Decimal
	Currency    string
	Description string
	Timestamp   int64
}

func (u *walletUsecaseImpl) RecordIncome(req RecordIncomeRequest) (*WalletResponse, error) {
	wallet, err := u.walletRepo.FindByID(req.WalletID)
	if err != nil {
		return nil, fmt.Errorf("wallet not found: %w", err)
	}
	if wallet == nil {
		return nil, fmt.Errorf("wallet not found")
	}

	money := shared.NewMoneyObject(req.Amount, req.Currency)
	err = wallet.RecordIncome(money, req.CategoryID, req.Timestamp, req.Description)
	if err != nil {
		return nil, fmt.Errorf("failed to record income: %w", err)
	}

	err = u.walletRepo.Save(wallet)
	if err != nil {
		return nil, fmt.Errorf("failed to save wallet: %w", err)
	}

	events := wallet.GetDomainEvents()
	if len(events) > 0 {
		err = u.dispatcher.Dispatch(events)
		if err != nil {
			return nil, fmt.Errorf("failed to dispatch events: %w", err)
		}
		wallet.ClearDomainEvents()
	}

	return mapWalletToResponse(wallet), nil
}

func (u *walletUsecaseImpl) GetWalletByID(id string) (*WalletResponse, error) {
	wallet, err := u.walletRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get wallet: %w", err)
	}
	if wallet == nil {
		return nil, fmt.Errorf("wallet not found")
	}
	return mapWalletToResponse(wallet), nil
}

func (u *walletUsecaseImpl) GetWalletsByOwnerID(ownerId string) ([]*WalletResponse, error) {
	wallets, err := u.walletRepo.FindByOwnerID(ownerId)
	if err != nil {
		return nil, fmt.Errorf("failed to get wallets: %w", err)
	}
	var res []*WalletResponse
	for _, w := range wallets {
		res = append(res, mapWalletToResponse(w))
	}
	return res, nil
}
