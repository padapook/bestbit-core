package service

import (
	"errors"
	"testing"

	"github.com/padapook/bestbit-core/internal/wallet/model"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockWalletRepository struct {
	mock.Mock
}

func (m *MockWalletRepository) GetWalletByUserID(userID string) ([]model.Wallet, error) {
	args := m.Called(userID)
	if args.Get(0) != nil {
		return args.Get(0).([]model.Wallet), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockWalletRepository) GetWalletByUserIDAndCurrency(userID, currency string) (*model.Wallet, error) {
	args := m.Called(userID, currency)
	if args.Get(0) != nil {
		return args.Get(0).(*model.Wallet), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockWalletRepository) Deposit(userID, currency string, amount decimal.Decimal, referenceID string) (*model.Wallet, error) {
	args := m.Called(userID, currency, amount, referenceID)
	if args.Get(0) != nil {
		return args.Get(0).(*model.Wallet), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockWalletRepository) Withdraw(userID, currency string, amount decimal.Decimal, referenceID string) (*model.Wallet, error) {
	args := m.Called(userID, currency, amount, referenceID)
	if args.Get(0) != nil {
		return args.Get(0).(*model.Wallet), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockWalletRepository) Transfer(fromUserID, toUserID, currency string, amount decimal.Decimal, referenceID string) error {
	args := m.Called(fromUserID, toUserID, currency, amount, referenceID)
	return args.Error(0)
}

func TestDepositMoney_Success(t *testing.T) {
	mockRepo := new(MockWalletRepository)
	service := NewWalletService(mockRepo)

	userID := "user-123"
	currency := "THB"
	amount := decimal.NewFromInt(1000)
	refID := "ref-123"

	expectedWallet := &model.Wallet{
		UserID:   userID,
		Currency: currency,
		Balance:  amount,
	}

	// เมื่อ Service เรียก Deposit ไปยัง Repo, ให้returnค่า expectedWallet
	mockRepo.On("Deposit", userID, currency, amount, refID).Return(expectedWallet, nil)

	wallet, err := service.DepositMoney(userID, currency, amount, refID)

	assert.NoError(t, err)
	assert.NotNil(t, wallet)
	assert.Equal(t, amount, wallet.Balance)

	mockRepo.AssertExpectations(t)
}

func TestDepositMoney_Fail_InvalidAmount(t *testing.T) {
	mockRepo := new(MockWalletRepository)
	service := NewWalletService(mockRepo)

	userID := "user-123"
	currency := "THB"
	// ยอดเงินติดลบ
	amount := decimal.NewFromInt(-500)
	refID := "ref-123"

	// กรณีนี้เราไม่คาดหวังให้ Repo.Deposit ถูกเรียกอเลย (เพราะควรโดนดักที่ Service ก่อน)
	wallet, err := service.DepositMoney(userID, currency, amount, refID)

	assert.Error(t, err)
	assert.Nil(t, wallet)
	assert.Equal(t, "deposit amount must be greater than zero", err.Error())
	mockRepo.AssertNotCalled(t, "Deposit")
}

func TestWithdrawMoney_Fail_InsufficientBalance(t *testing.T) {
	mockRepo := new(MockWalletRepository)
	service := NewWalletService(mockRepo)

	userID := "user-123"
	currency := "THB"
	amount := decimal.NewFromInt(2000)
	refID := "ref-123"

	expectedError := errors.New("insufficient balance")

	// จำลองว่ายอดเงินไม่พอ
	mockRepo.On("Withdraw", userID, currency, amount, refID).Return(nil, expectedError)

	wallet, err := service.WithdrawMoney(userID, currency, amount, refID)

	assert.Error(t, err)
	assert.Nil(t, wallet)
	assert.Equal(t, expectedError, err)
	mockRepo.AssertExpectations(t)
}

func TestTransferMoney_Fail_SelfTransfer(t *testing.T) {
	mockRepo := new(MockWalletRepository)
	service := NewWalletService(mockRepo)

	userID := "user-123"
	amount := decimal.NewFromInt(500)
	refID := "ref-123"

	// โอนให้ตัวเอง
	err := service.TransferMoney(userID, userID, "THB", amount, refID)

	assert.Error(t, err)
	assert.Equal(t, "cannot transfer to yourself", err.Error())
	mockRepo.AssertNotCalled(t, "Transfer")
}
