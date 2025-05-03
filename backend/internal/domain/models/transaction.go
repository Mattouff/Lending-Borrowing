package models

import (
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"gorm.io/gorm"
)

// TransactionType defines the type of blockchain transaction
type TransactionType string

const (
	// TransactionDeposit represents a deposit of collateral
	TransactionDeposit TransactionType = "deposit"
	// TransactionWithdraw represents a withdrawal of collateral
	TransactionWithdraw TransactionType = "withdraw"
	// TransactionBorrow represents borrowing tokens
	TransactionBorrow TransactionType = "borrow"
	// TransactionRepay represents repaying borrowed tokens
	TransactionRepay TransactionType = "repay"
	// TransactionLiquidate represents liquidation of a position
	TransactionLiquidate TransactionType = "liquidate"
)

// TransactionStatus defines the status of a blockchain transaction
type TransactionStatus string

const (
	// StatusPending means transaction is pending
	StatusPending TransactionStatus = "pending"
	// StatusCompleted means transaction is completed
	StatusCompleted TransactionStatus = "completed"
	// StatusFailed means transaction failed
	StatusFailed TransactionStatus = "failed"
)

// Transaction represents a blockchain transaction in the platform
type Transaction struct {
	ID           uint              `json:"id" gorm:"primaryKey"`
	UserID       uint              `json:"userId" gorm:"index;not null"`
	User         *User             `json:"user" gorm:"foreignKey:UserID"`
	Type         TransactionType   `json:"type" gorm:"type:varchar(20);not null"`
	Status       TransactionStatus `json:"status" gorm:"type:varchar(20);not null;default:'pending'"`
	Hash         string            `json:"hash" gorm:"type:varchar(66);unique"`
	BlockNumber  *uint64           `json:"blockNumber"`
	Amount       string            `json:"amount" gorm:"type:varchar(78);not null"` // Big numbers stored as strings
	TokenAddress string            `json:"tokenAddress" gorm:"type:varchar(42);not null"`
	GasUsed      *uint64           `json:"gasUsed"`
	GasPrice     string            `json:"gasPrice" gorm:"type:varchar(78)"` // Big numbers stored as strings
	ErrorMessage string            `json:"errorMessage" gorm:"type:text"`
	CreatedAt    time.Time         `json:"createdAt"`
	UpdatedAt    time.Time         `json:"updatedAt"`
	DeletedAt    gorm.DeletedAt    `json:"deletedAt" gorm:"index"`
}

// BigIntAmount converts the amount string to a big.Int
func (t *Transaction) BigIntAmount() (*big.Int, bool) {
	amount := new(big.Int)
	amount, success := amount.SetString(t.Amount, 10)
	return amount, success
}

// GetTokenAddress returns the token address as a common.Address
func (t *Transaction) GetTokenAddress() common.Address {
	return common.HexToAddress(t.TokenAddress)
}
