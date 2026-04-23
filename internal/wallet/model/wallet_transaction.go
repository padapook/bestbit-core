package model

import (
	"github.com/shopspring/decimal"
	"time"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type WalletTransaction struct {
	ID              uuid.UUID       `gorm:"primaryKey" json:"id"`
	WalletID        uuid.UUID       `gorm:"type:uuid;index;not null" json:"wallet_id"`
	TargetWalletID	*string 		`gorm:"index" json:"target_wallet_id,omitempty"`
	ReferenceID     string          `gorm:"uniqueIndex;not null" json:"reference_id"`
	TransactionType	string 			`gorm:"type:varchar(20);not null" json:"transaction_type"`
	Status          string          `gorm:"type:varchar(20);not null;default:'PENDING'" json:"status"`
	Amount          decimal.Decimal `gorm:"type:decimal(32,16); default:0" json:"amount"`
	Currency 		string 			`gorm:"type:varchar(20);not null;default:'THB'" json:"currency"`
	BalanceBefore   decimal.Decimal `gorm:"type:decimal(32,16);not null" json:"balance_before"`
    BalanceAfter    decimal.Decimal `gorm:"type:decimal(32,16);not null" json:"balance_after"`
	Description     string          `gorm:"type:text" json:"description"`
	Remark			string			`gorm:"type:text" json:"remark"`
    CreatedAt       time.Time       `gorm:"index" json:"created_at"`
    CreatedBy       string          `gorm:"size:100" json:"created_by"`
}

func (m *WalletTransaction) BeforeCreate(tx *gorm.DB) error {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	return nil
}