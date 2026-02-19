package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type Wallet struct {
	ID  			uint64 		`gorm:"primaryKey" json:"id"`
	UserID 			uint64 		`gorm:"index;uniqueIndex:idx_user_currency" json:"user_id"`
	Currency 		string  	`gorm:"size:10;uniqueIndex:idx_user_currency" json:"currency"`
	Balance 		decimal.Decimal 	`gorm:"type:decimal(32,16);default:0" json:"balance"`
	AmountLocked 	decimal.Decimal 	`gorm:"type:decimal(32,16);default:0" json:"amount_locked"`
	UpdatedAt 		time.Time   `json:"updated_at"`
	IsActive 		bool 		`json:"is_active"`
}