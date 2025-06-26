package service

import (
	"context"
	"math/big"

	"github.com/Mattouff/Lending-Borrowing/internal/domain/models"
	"github.com/ethereum/go-ethereum/common"
)

// LendingService defines the interface for lending business logic
type LendingService interface {
	// Wallet-based transaction methods
	// PrepareDeposit prepares the contract interaction data for a deposit transaction
	PrepareDeposit(ctx context.Context, userAddress common.Address, amount *big.Int) (*ContractInteractionData, error)

	// PrepareWithdraw prepares the contract interaction data for a withdraw transaction
	PrepareWithdraw(ctx context.Context, userAddress common.Address, amount *big.Int) (*ContractInteractionData, error)

	// ProcessWalletTransaction validates and processes a wallet-signed transaction
	ProcessWalletTransaction(ctx context.Context, txHash string, userAddress common.Address, txType models.TransactionType, amount *big.Int) (*models.Transaction, error)

	// Read-only methods
	// GetUserBalance returns the user's balance in the lending pool
	GetUserBalance(ctx context.Context, userAddress common.Address) (*big.Int, error)

	// GetTotalDeposited returns the total amount deposited in the lending pool
	GetTotalDeposited(ctx context.Context) (*big.Int, error)

	// GetCurrentInterestRate returns the current interest rate for lending
	GetCurrentInterestRate(ctx context.Context) (*big.Int, error)

	// GetUserInterestEarned returns the interest earned by a user
	GetUserInterestEarned(ctx context.Context, userAddress common.Address) (*big.Int, error)

	// GetUserTransactionHistory returns a user's lending transaction history
	GetUserTransactionHistory(ctx context.Context, userAddress common.Address, offset, limit int) ([]*models.Transaction, error)

	// CountUserTransactions counts the number of transactions for a user with optional filtering
	CountUserTransactions(ctx context.Context, address common.Address, filter map[string]any) (int64, error)
}
