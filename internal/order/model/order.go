package model

import (
	"time"
)

type Order struct {
	ID           uint64    `gorm:"primaryKey" json:"id"`
	UserID       uint64    `gorm:"index" json:"user_id"`
	Symbol       string    `gorm:"size:20;index" json:"symbol"`
	Side         string    `gorm:"size:10;index" json:"side" comment:"BUY, SELL"`
	OrderType    string    `gorm:"size:10" json:"order_type" comment:"LIMIT, MARKET"`
	Status       string    `gorm:"size:20;index" json:"status" comment:"PENDING, PARTIAL_FILLED, FILLED, CANCELED"`
	Price        float64   `gorm:"type:decimal(32,16)" json:"price"`
	Amount       float64   `gorm:"type:decimal(32,16)" json:"amount"`
	FilledAmount float64   `gorm:"type:decimal(32,16);default:0" json:"filled_amount" comment:""`
	CreatedAt    time.Time `gorm:"index" json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
