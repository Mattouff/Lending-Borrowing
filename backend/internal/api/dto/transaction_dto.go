package dto

import "time"

// TransactionType for DTO
type TransactionType string

const (
	TransactionTypeDeposit   TransactionType = "deposit"
	TransactionTypeWithdraw  TransactionType = "withdraw"
	TransactionTypeBorrow    TransactionType = "borrow"
	TransactionTypeRepay     TransactionType = "repay"
	TransactionTypeLiquidate TransactionType = "liquidate"
)

// TransactionStatus for DTO
type TransactionStatus string

const (
	TransactionStatusPending   TransactionStatus = "pending"
	TransactionStatusConfirmed TransactionStatus = "confirmed"
	TransactionStatusFailed    TransactionStatus = "failed"
)

// TransactionRequest represents data for a new transaction
type TransactionRequest struct {
	Amount string `json:"amount" validate:"required"`
}

// TransactionLiquidationRequest represents data for a liquidation transaction
type TransactionLiquidationRequest struct {
	BorrowerAddress string `json:"borrowerAddress" validate:"required,eth_addr"`
	Amount          string `json:"amount" validate:"required"`
}

// TransactionResponse represents a transaction in API responses
type TransactionResponse struct {
	ID           uint              `json:"id"`
	UserID       uint              `json:"userId"`
	Type         TransactionType   `json:"type"`
	Status       TransactionStatus `json:"status"`
	Hash         string            `json:"hash"`
	Amount       string            `json:"amount"`
	TokenAddress string            `json:"tokenAddress"`
	BlockNumber  uint64            `json:"blockNumber,omitempty"`
	GasUsed      uint64            `json:"gasUsed,omitempty"`
	GasPrice     string            `json:"gasPrice,omitempty"`
	CreatedAt    time.Time         `json:"createdAt"`
	UpdatedAt    time.Time         `json:"updatedAt"`
}

// TransactionListResponse represents a list of transactions for API responses
type TransactionListResponse struct {
	Transactions []TransactionResponse `json:"transactions"`
	Total        int64                 `json:"total"`
	Page         int                   `json:"page"`
	PageSize     int                   `json:"pageSize"`
	TotalPage    int                   `json:"totalPage"`
}
