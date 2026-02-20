package repository

import (
	"errors"
	"time"

	"github.com/padapook/bestbit-core/internal/wallet/model"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type WalletRepository interface {
	GetWalletByUserID(userID string) ([]model.Wallet, error)
	GetWalletByUserIDAndCurrency(userID, currency string) (*model.Wallet, error)
	Deposit(userID, currency string, amount decimal.Decimal, referenceID string) (*model.Wallet, error)
	Withdraw(userID, currency string, amount decimal.Decimal, referenceID string) (*model.Wallet, error)
	Transfer(fromUserID, toUserID, currency string, amount decimal.Decimal, referenceID string) error
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

func (r *walletRepository) Deposit(userID, currency string, amount decimal.Decimal, referenceID string) (*model.Wallet, error) {
	var wallet model.Wallet

	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("user_id = ? AND currency = ?", userID, currency).
			First(&wallet).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("wallet not found")
			}
			return err
		}

		balanceBefore := wallet.Balance
		wallet.Balance = wallet.Balance.Add(amount)
		balanceAfter := wallet.Balance

		if err := tx.Save(&wallet).Error; err != nil {
			return err
		}

		trx := model.Transaction{
			WalletID:        wallet.ID,
			ReferenceID:     referenceID,
			TransactionType: "DEPOSIT",
			Amount:          amount,
			BalanceBefore:   balanceBefore,
			BalanceAfter:    balanceAfter,
			Status:          "COMPLETED",
			Description:     "Deposit via API",
			CreatedAt:       time.Now(),
			CreatedBy:       userID,
		}

		if err := tx.Create(&trx).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &wallet, nil
}

func (r *walletRepository) Withdraw(userID, currency string, amount decimal.Decimal, referenceID string) (*model.Wallet, error) {
	var wallet model.Wallet

	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("user_id = ? AND currency = ?", userID, currency).
			First(&wallet).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("wallet not found")
			}
			return err
		}

		if wallet.Balance.LessThan(amount) {
			return errors.New("insufficient balance")
		}

		balanceBefore := wallet.Balance
		wallet.Balance = wallet.Balance.Sub(amount)
		balanceAfter := wallet.Balance

		if err := tx.Save(&wallet).Error; err != nil {
			return err
		}

		trx := model.Transaction{
			WalletID:        wallet.ID,
			ReferenceID:     referenceID,
			TransactionType: "WITHDRAW",
			Amount:          amount,
			BalanceBefore:   balanceBefore,
			BalanceAfter:    balanceAfter,
			Status:          "COMPLETED",
			Description:     "Withdraw via API",
			CreatedAt:       time.Now(),
			CreatedBy:       userID,
		}

		if err := tx.Create(&trx).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &wallet, nil
}

func (r *walletRepository) Transfer(fromUserID, toUserID, currency string, amount decimal.Decimal, referenceID string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		firstID, secondID := fromUserID, toUserID
		if firstID > secondID {
			firstID, secondID = secondID, firstID
		}

		var firstWallet model.Wallet
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("user_id = ? AND currency = ?", firstID, currency).First(&firstWallet).Error; err != nil {
			return errors.New("wallet not found for user: " + firstID)
		}

		var secondWallet model.Wallet
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("user_id = ? AND currency = ?", secondID, currency).First(&secondWallet).Error; err != nil {
			return errors.New("wallet not found for user: " + secondID)
		}

		var senderWallet, receiverWallet *model.Wallet
		if fromUserID == firstID {
			senderWallet = &firstWallet
			receiverWallet = &secondWallet
		} else {
			senderWallet = &secondWallet
			receiverWallet = &firstWallet
		}

		if senderWallet.Balance.LessThan(amount) {
			return errors.New("insufficient balance")
		}

		senderBalBefore := senderWallet.Balance
		senderWallet.Balance = senderWallet.Balance.Sub(amount)
		senderBalAfter := senderWallet.Balance

		recvBalBefore := receiverWallet.Balance
		receiverWallet.Balance = receiverWallet.Balance.Add(amount)
		recvBalAfter := receiverWallet.Balance

		if err := tx.Save(senderWallet).Error; err != nil {
			return err
		}
		if err := tx.Save(receiverWallet).Error; err != nil {
			return err
		}

		// tx ฝั่ง sender
		txSender := model.Transaction{
			WalletID:        senderWallet.ID,
			ReferenceID:     referenceID,
			TransactionType: "TRANSFER_OUT",
			Amount:          amount,
			BalanceBefore:   senderBalBefore,
			BalanceAfter:    senderBalAfter,
			Status:          "COMPLETED",
			Description:     "Transfer to " + toUserID,
			CreatedAt:       time.Now(),
			CreatedBy:       fromUserID,
		}
		if err := tx.Create(&txSender).Error; err != nil {
			return err
		}

		// tx ฝั่ง receive
		txReceiver := model.Transaction{
			WalletID:        receiverWallet.ID,
			ReferenceID:     referenceID,
			TransactionType: "TRANSFER_IN",
			Amount:          amount,
			BalanceBefore:   recvBalBefore,
			BalanceAfter:    recvBalAfter,
			Status:          "COMPLETED",
			Description:     "Transfer from " + fromUserID,
			CreatedAt:       time.Now(),
			CreatedBy:       fromUserID,
		}
		if err := tx.Create(&txReceiver).Error; err != nil {
			return err
		}

		return nil
	})
}
