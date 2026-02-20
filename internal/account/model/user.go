package model

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID           uint64         `gorm:"primaryKey" json:"uid"`
	AccountId    string         `gorm:"size:100;not null;uniqueIndex" json:"account_id"`
	Username     string         `gorm:"size:50;not null;uniqueIndex" json:"username"`
	Password     string         `gorm:"size:255;not null" json:"-"`
	TitleName    string         `gorm:"size:20" json:"title_name"`
	FirstName    string         `gorm:"size:100" json:"first_name"`
	MiddleName   string         `gorm:"size:100" json:"middle_name"`
	LastName     string         `gorm:"size:100" json:"last_name"`
	Email        string         `gorm:"size:100;index;uniqueIndex" json:"email"`
	MobileNumber string         `gorm:"size:20;index" json:"mobile_number"`
	IsActive     bool           `gorm:"default:true" json:"is_active"`
	CreatedAt    time.Time      `gorm:"autoCreateTime" json:"created_at"`
	CreatedBy    string         `gorm:"size:50;default:'SYSTEM'" json:"created_by"`
	UpdatedAt    time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	UpdatedBy    string         `gorm:"size:50" json:"updated_by"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
	DeletedBy    string         `gorm:"size:50" json:"deleted_by"`
}
