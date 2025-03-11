package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	WalletAddress string `gorm:"unique;not null" json:"wallet_address"`
	Email         string `gorm:"unique" json:"email"`
	Nonce         string `gorm:"not null" json:"nonce"` // For wallet authentication
}
