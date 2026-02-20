package service

import (
	"errors"

	"github.com/padapook/bestbit-core/internal/wallet/model"
	"github.com/padapook/bestbit-core/internal/wallet/repository"
	"github.com/shopspring/decimal"
)

type WalletService interface {
	GetUserWallets(userID string) ([]model.Wallet, error)
	GetWalletBalance(userID, currency string) (*model.Wallet, error)
	DepositMoney(userID, currency string, amount decimal.Decimal, referenceID string) (*model.Wallet, error)
	WithdrawMoney(userID, currency string, amount decimal.Decimal, referenceID string) (*model.Wallet, error)
	TransferMoney(fromUserID, toUserID, currency string, amount decimal.Decimal, referenceID string) error
}

type walletService struct {
	repo repository.WalletRepository
}

func NewWalletService(repo repository.WalletRepository) WalletService {
	return &walletService{repo: repo}
}

func (s *walletService) GetUserWallets(userID string) ([]model.Wallet, error) {
	return s.repo.GetWalletByUserID(userID)
}

func (s *walletService) GetWalletBalance(userID, currency string) (*model.Wallet, error) {
	return s.repo.GetWalletByUserIDAndCurrency(userID, currency)
}

func (s *walletService) DepositMoney(userID, currency string, amount decimal.Decimal, referenceID string) (*model.Wallet, error) {
	if amount.LessThanOrEqual(decimal.Zero) {
		return nil, errors.New("deposit amount must be greater than zero")
	}

	return s.repo.Deposit(userID, currency, amount, referenceID)
}

func (s *walletService) WithdrawMoney(userID, currency string, amount decimal.Decimal, referenceID string) (*model.Wallet, error) {
	if amount.LessThanOrEqual(decimal.Zero) {
		return nil, errors.New("withdraw amount must be greater than zero")
	}

	return s.repo.Withdraw(userID, currency, amount, referenceID)
}

func (s *walletService) TransferMoney(fromUserID, toUserID, currency string, amount decimal.Decimal, referenceID string) error {
	if fromUserID == toUserID {
		return errors.New("cannot transfer to yourself")
	}

	if amount.LessThanOrEqual(decimal.Zero) {
		return errors.New("transfer amount must be greater than zero")
	}

	return s.repo.Transfer(fromUserID, toUserID, currency, amount, referenceID)
}
