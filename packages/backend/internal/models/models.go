package models

import (
	"time"
)

type User struct {
	ID           uint      `gorm:"primarykey" json:"id"`
	WalletAddr   string    `gorm:"uniqueIndex;size:42" json:"wallet_addr"`
	Email        string    `gorm:"index" json:"email"`
	TokenID      uint      `json:"token_id"`
	IsActive     bool      `gorm:"default:true" json:"is_active"`
	RegisteredAt time.Time `json:"registered_at"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Operator struct {
	ID              uint      `gorm:"primarykey" json:"id"`
	Name            string    `json:"name"`
	Region          string    `json:"region"`
	RequiredDeposit string    `json:"required_deposit"`
	IsActive        bool      `gorm:"default:true" json:"is_active"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type Deposit struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	UserID    uint      `gorm:"index" json:"user_id"`
	Amount    string    `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Bill struct {
	ID          uint       `gorm:"primarykey" json:"id"`
	UserID      uint       `gorm:"index" json:"user_id"`
	OperatorID  uint       `gorm:"index" json:"operator_id"`
	Amount      string     `json:"amount"`
	PlatformFee string     `json:"platform_fee"`
	IsPaid      bool       `gorm:"default:false" json:"is_paid"`
	CreatedAt   time.Time  `json:"created_at"`
	PaidAt      *time.Time `json:"paid_at,omitempty"`
}

type UserService struct {
	ID            uint      `gorm:"primarykey" json:"id"`
	UserID        uint      `gorm:"uniqueIndex" json:"user_id"`
	OperatorID    uint      `gorm:"index" json:"operator_id"`
	VirtualNumber string    `json:"virtual_number"`
	Password      string    `json:"password"`
	IsActive      bool      `gorm:"default:true" json:"is_active"`
	ActivatedAt   time.Time `json:"activated_at"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type UsageData struct {
	ID         uint      `gorm:"primarykey" json:"id"`
	UserID     uint      `gorm:"index" json:"user_id"`
	OperatorID uint      `gorm:"index" json:"operator_id"`
	DataUsage  uint64    `json:"data_usage"`
	CallUsage  uint64    `json:"call_usage"`
	Timestamp  time.Time `json:"timestamp"`
	Signature  string    `gorm:"type:text" json:"signature"`
	CreatedAt  time.Time `json:"created_at"`
}
