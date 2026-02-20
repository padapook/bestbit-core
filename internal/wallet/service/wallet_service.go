package service

import (
	"github.com/padapook/bestbit-core/internal/wallet/model"
	"github.com/padapook/bestbit-core/internal/wallet/repository"
)

type WalletService interface {
	GetUserWallets(userID string) ([]model.Wallet, error)
	GetWalletBalance(userID, currency string) (*model.Wallet, error)
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
