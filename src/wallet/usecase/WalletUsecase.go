package usecase

import (
	"fmt"
	"github.com/shopspring/decimal"
	shared "src/src/shared/domain"
	"src/src/wallet/domain"
)

type WalletUsecase interface {
	CreateWallet(req CreateWalletRequest) error
	RecordExpense(req RecordExpenseRequest) error
	RecordIncome(req RecordIncomeRequest) error
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

func (u *walletUsecaseImpl) CreateWallet(req CreateWalletRequest) error {
	balance := shared.NewMoneyObject(decimal.NewFromInt(0), req.Currency)
	wallet := domain.NewWallet(req.OwnerID, req.Name, balance)

	err := u.walletRepo.Save(&wallet)
	if err != nil {
		return err
	}

	events := wallet.GetDomainEvents()
	if len(events) > 0 {
		err = u.dispatcher.Dispatch(events)
		if err != nil {
			return fmt.Errorf("failed to dispatch events: %w", err)
		}
		wallet.ClearDomainEvents()
	}

	return nil
}

type RecordExpenseRequest struct {
	WalletID    string
	CategoryID  string
	Amount      decimal.Decimal
	Currency    string
	Description string
	Timestamp   int64
}

func (u *walletUsecaseImpl) RecordExpense(req RecordExpenseRequest) error {
	wallet, err := u.walletRepo.FindByID(req.WalletID)
	if err != nil {
		return fmt.Errorf("wallet not found: %w", err)
	}
	if wallet == nil {
		return fmt.Errorf("wallet not found")
	}

	money := shared.NewMoneyObject(req.Amount, req.Currency)
	err = wallet.RecordExpense(money, req.CategoryID, req.Timestamp, req.Description)
	if err != nil {
		return fmt.Errorf("failed to record expense: %w", err)
	}

	err = u.walletRepo.Save(wallet)
	if err != nil {
		return fmt.Errorf("failed to save wallet: %w", err)
	}

	events := wallet.GetDomainEvents()
	if len(events) > 0 {
		err = u.dispatcher.Dispatch(events)
		if err != nil {
			return fmt.Errorf("failed to dispatch events: %w", err)
		}
		wallet.ClearDomainEvents()
	}

	return nil
}

type RecordIncomeRequest struct {
	WalletID    string
	CategoryID  string
	Amount      decimal.Decimal
	Currency    string
	Description string
	Timestamp   int64
}

func (u *walletUsecaseImpl) RecordIncome(req RecordIncomeRequest) error {
	wallet, err := u.walletRepo.FindByID(req.WalletID)
	if err != nil {
		return fmt.Errorf("wallet not found: %w", err)
	}
	if wallet == nil {
		return fmt.Errorf("wallet not found")
	}

	money := shared.NewMoneyObject(req.Amount, req.Currency)
	err = wallet.RecordIncome(money, req.CategoryID, req.Timestamp, req.Description)
	if err != nil {
		return fmt.Errorf("failed to record income: %w", err)
	}

	err = u.walletRepo.Save(wallet)
	if err != nil {
		return fmt.Errorf("failed to save wallet: %w", err)
	}

	events := wallet.GetDomainEvents()
	if len(events) > 0 {
		err = u.dispatcher.Dispatch(events)
		if err != nil {
			return fmt.Errorf("failed to dispatch events: %w", err)
		}
		wallet.ClearDomainEvents()
	}

	return nil
}
