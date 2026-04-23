package repository

import (
	"errors"

	"github.com/padapook/bestbit-core/internal/wallet/model"
	"gorm.io/gorm"
)

type WalletRepository interface {
	GetWalletByUserID(userID string) ([]model.Wallet, error)
	GetWalletByUserIDAndCurrency(userID, currency string) (*model.Wallet, error)
}

type walletRepository struct {
	db *gorm.DB
}

func NewWalletRepository(db *gorm.DB) WalletRepository {
	return &walletRepository{db: db}
}

func (r *walletRepository) GetWalletByUserID(userID string) ([]model.Wallet, error) {
	var wallets []model.Wallet
	err := r.db.Where("user_id = ?", userID).Find(&wallets).Error
	return wallets, err
}

func (r *walletRepository) GetWalletByUserIDAndCurrency(userID, currency string) (*model.Wallet, error) {
	var wallet model.Wallet
	err := r.db.Where("user_id = ? AND currency = ?", userID, currency).First(&wallet).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("wallet not found")
		}
		return nil, err
	}
	return &wallet, nil
}
