package repository

import (
	"context"
	"errors"

	"github.com/padapook/bestbit-core/internal/wallet/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/shopspring/decimal"
)

type WalletRepository interface {
	GetWalletByUserID(ctx context.Context, userID string) ([]model.Wallet, error)
	GetWalletByUserIDAndCurrency(ctx context.Context, userID, currency string) (*model.Wallet, error)
	LockFunds(ctx context.Context, tx *gorm.DB, userId, currency string, amount decimal.Decimal) (*model.Wallet, error)
	UnlockFunds(ctx context.Context, tx *gorm.DB, userId, currency string, amount decimal.Decimal) (*model.Wallet, error)
}

type walletRepository struct {
	db *gorm.DB
}

func NewWalletRepository(db *gorm.DB) WalletRepository {
	return &walletRepository{db: db}
}

func (r *walletRepository) GetWalletByUserID(ctx context.Context, userID string) ([]model.Wallet, error) {
	var wallets []model.Wallet
	err := r.db.Where("user_id = ?", userID).Find(&wallets).Error
	return wallets, err
}

func (r *walletRepository) GetWalletByUserIDAndCurrency(ctx context.Context, userID, currency string) (*model.Wallet, error) {
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

// lock เงิน user กรณีที่มี order
func (r *walletRepository) LockFunds(ctx context.Context, tx *gorm.DB, userID, currency string, amount decimal.Decimal) (*model.Wallet, error) {
	var wallet model.Wallet
	//ส่งค่า userid, currency และ amount มา
	// เช็คว่า userId มี currentcy นี้ไหม
	err := tx.WithContext(ctx).Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("user_id = ? AND currency = ?", userID, currency).
		First(&wallet).Error
	if err != nil {
		return nil, model.ErrWalletNotFound
	}

	//ถ้ามี จะเช็คว่า balance ของ currency นี้ มีมากกว่า amount ไหม
	if wallet.Balance.LessThan(amount) {
		return nil, model.ErrInSufficientFunds
	}
	// ทำการหัก balance ไปเก็บไว้ที่ amount_locked
	wallet.Balance = wallet.Balance.Sub(amount)
	wallet.AmountLocked = wallet.AmountLocked.Add(amount)

	err = tx.Save(&wallet).Error

	return &wallet, nil
}

// ปลดล็อคเงินเมื่อ user ยกเลิกออเดอร์
func (r *walletRepository) UnlockFunds(ctx context.Context, tx *gorm.DB, userId, currency string, amount decimal.Decimal) (*model.Wallet, error) {
	//
	var wallet model.Wallet

	err := tx.WithContext(ctx).Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("user_id = ? AND currency = ?", userId, currency).
		First(&wallet).Error
	if err != nil {
		return nil, model.ErrWalletNotFound
	}

	if wallet.AmountLocked.LessThan(amount) {
		return  nil, model.ErrInSufficientFunds
	}

	// คืนยอดตามจำนวน amount
	wallet.AmountLocked = wallet.AmountLocked.Sub(amount)
	wallet.Balance = wallet.Balance.Add(amount)

	err = tx.Save(&wallet).Error

	return &wallet, nil
}
