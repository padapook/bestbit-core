package model

import (
	"github.com/shopspring/decimal"
	"time"
)

type Order struct {
	ID           uint64          `gorm:"primaryKey" json:"id"`
	UserID       uint64          `gorm:"index" json:"userId"`
	Symbol       string          `gorm:"size:20;index" json:"symbol"`
	Side         OrderSide       `gorm:"size:10;index" json:"orderSide"`
	OrderType    OrderType       `gorm:"size:10" json:"orderType"`
	Status       OrderStatus     `gorm:"size:20;index" json:"orderStatus"`
	Price        decimal.Decimal `gorm:"type:decimal(32,16)" json:"price"`
	Amount       decimal.Decimal `gorm:"type:decimal(32,16)" json:"amount"`
	FilledAmount decimal.Decimal `gorm:"type:decimal(32,16);default:0" json:"filledAmount" comment:""`
	CreatedAt    time.Time       `gorm:"index" json:"createdAt"`
	UpdatedAt    time.Time       `json:"updatedAt"`
}
