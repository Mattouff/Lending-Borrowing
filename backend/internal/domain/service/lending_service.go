package service

import (
	"context"
	"math/big"

	"github.com/Mattouff/Lending-Borrowing/internal/domain/models"
	"github.com/ethereum/go-ethereum/common"
)

// LendingService defines the interface for lending business logic
type LendingService interface {
	// Deposit allows users to deposit tokens into the lending pool
	Deposit(ctx context.Context, userAddress common.Address, amount *big.Int) (string, error)

	// Withdraw allows users to withdraw tokens from the lending pool
	Withdraw(ctx context.Context, userAddress common.Address, amount *big.Int) (string, error)

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
}
