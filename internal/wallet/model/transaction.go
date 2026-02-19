package model

import (
	"time"
)

type Transaction struct {
	ID 				uint64 		`gorm:"primaryKey" json:"id"`
	WalletID		uint64 		`gorm:"index" json:"wallet_id"`
	ReferenceID		string 		`gorm:"index" json:"reference_id"`
	TransactionType	string 		`gorm:"size:20" json:"transaction_type"`
	Amount          float64   	`gorm:"type:decimal(32,16)" json:"amount"`
    BalanceBefore   float64   	`gorm:"type:decimal(32,16)" json:"balance_before"`
    BalanceAfter    float64   	`gorm:"type:decimal(32,16)" json:"balance_after"`
    Status          string    	`gorm:"size:20" json:"status"`
    Description     string    	`gorm:"type:text" json:"description"`
    CreatedAt       time.Time 	`gorm:"index" json:"created_at"`
	CreatedBy		string 		`gorm:"size:100" json:"created_by"`
}