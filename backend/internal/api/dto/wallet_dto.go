package dto

// WalletTransactionRequest represents a transaction signed by the user's wallet
type WalletTransactionRequest struct {
	// Transaction details
	Amount       string          `json:"amount" validate:"required"`
	TokenAddress string          `json:"tokenAddress" validate:"required,eth_addr"`
	Type         TransactionType `json:"type" validate:"required"` // Added transaction type

	// Signed transaction data
	TransactionHash string `json:"transactionHash" validate:"required"`
	RawTransaction  string `json:"rawTransaction" validate:"required"` // Signed transaction hex

	// Additional verification data
	UserAddress string `json:"userAddress" validate:"required,eth_addr"`
	Nonce       string `json:"nonce,omitempty"`
}

// WalletTransactionResponse represents the response after submitting a wallet transaction
type WalletTransactionResponse struct {
	TransactionID   uint   `json:"transactionId"`
	TransactionHash string `json:"transactionHash"`
	Status          string `json:"status"`
	Message         string `json:"message"`
}
