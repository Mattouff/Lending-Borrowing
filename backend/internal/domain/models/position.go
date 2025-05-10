package models

import (
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"gorm.io/gorm"
)

// PositionStatus represents the status of a lending/borrowing position
type PositionStatus string

const (
	// StatusActive means the position is active
	StatusActive PositionStatus = "active"
	// StatusLiquidated means the position was liquidated
	StatusLiquidated PositionStatus = "liquidated"
	// StatusClosed means the position was closed by the user
	StatusClosed PositionStatus = "closed"
)

// Position represents a user's lending or borrowing position
type Position struct {
	ID                 uint           `json:"id" gorm:"primaryKey"`
	UserID             uint           `json:"userId" gorm:"index;not null"`
	User               *User          `json:"user" gorm:"foreignKey:UserID"`
	CollateralAmount   string         `json:"collateralAmount" gorm:"type:varchar(78);not null"` // Big numbers stored as strings
	CollateralToken    string         `json:"collateralToken" gorm:"type:varchar(42);not null"`
	BorrowedAmount     string         `json:"borrowedAmount" gorm:"type:varchar(78);not null"` // Big numbers stored as strings
	BorrowedToken      string         `json:"borrowedToken" gorm:"type:varchar(42);not null"`
	InterestRate       string         `json:"interestRate" gorm:"type:varchar(78);not null"` // Interest rate as a big number (e.g. 5% = 5 * 10^16)
	LastInterestUpdate time.Time      `json:"lastInterestUpdate"`
	Status             PositionStatus `json:"status" gorm:"type:varchar(20);default:'active'"`
	HealthFactor       string         `json:"healthFactor" gorm:"type:varchar(78)"`     // Current health factor of the position
	LiquidationPrice   string         `json:"liquidationPrice" gorm:"type:varchar(78)"` // Price at which position is liquidated
	CreatedAt          time.Time      `json:"createdAt"`
	UpdatedAt          time.Time      `json:"updatedAt"`
	DeletedAt          gorm.DeletedAt `json:"deletedAt" gorm:"index"`
}

// CollateralBigInt converts the collateral amount to a big.Int
func (p *Position) CollateralBigInt() (*big.Int, bool) {
	amount := new(big.Int)
	amount, success := amount.SetString(p.CollateralAmount, 10)
	return amount, success
}

// BorrowedBigInt converts the borrowed amount to a big.Int
func (p *Position) BorrowedBigInt() (*big.Int, bool) {
	amount := new(big.Int)
	amount, success := amount.SetString(p.BorrowedAmount, 10)
	return amount, success
}

// GetCollateralTokenAddress returns the collateral token address as a common.Address
func (p *Position) GetCollateralTokenAddress() common.Address {
	return common.HexToAddress(p.CollateralToken)
}

// GetBorrowedTokenAddress returns the borrowed token address as a common.Address
func (p *Position) GetBorrowedTokenAddress() common.Address {
	return common.HexToAddress(p.BorrowedToken)
}
