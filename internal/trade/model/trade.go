package model

import (
	"time"
)

type Trade struct {
	ID           uint64    `gorm:"primaryKey" json:"id"`
	Symbol       string    `gorm:"size:20;index" json:"symbol" comment:""`
	MakerOrderID uint64    `gorm:"index" json:"maker_order_id" comment:""`
	TakerOrderID uint64    `gorm:"index" json:"taker_order_id" comment:""`
	Price        float64   `gorm:"type:decimal(32,16)" json:"price"`
	Amount       float64   `gorm:"type:decimal(32,16)" json:"amount"`
	ExecutedAt   time.Time `gorm:"index" json:"executed_at"`
}
