package service

import (
	"context"
	"math/big"

	"github.com/Mattouff/Lending-Borrowing/internal/domain/models"
	"github.com/ethereum/go-ethereum/common"
)

// BorrowingService defines the interface for borrowing business logic
type BorrowingService interface {
	// Wallet-based transaction methods
	// PrepareBorrow prepares the contract interaction data for a borrow transaction
	PrepareBorrow(ctx context.Context, userAddress common.Address, amount *big.Int) (*ContractInteractionData, error)

	// PrepareRepay prepares the contract interaction data for a repay transaction
	PrepareRepay(ctx context.Context, userAddress common.Address, amount *big.Int) (*ContractInteractionData, error)

	// ProcessWalletTransaction validates and processes a wallet-signed transaction
	ProcessWalletTransaction(ctx context.Context, txHash string, userAddress common.Address, txType models.TransactionType, amount *big.Int) (*models.Transaction, error)

	// Read-only methods
	// GetBorrowedAmount returns the amount borrowed by a user
	GetBorrowedAmount(ctx context.Context, userAddress common.Address) (*big.Int, error)

	// GetTotalBorrowed returns the total amount borrowed from the protocol
	GetTotalBorrowed(ctx context.Context) (*big.Int, error)

	// GetCurrentInterestRate returns the current interest rate for borrowing
	GetCurrentInterestRate(ctx context.Context) (*big.Int, error)

	// GetUserInterestAccrued returns the interest accrued by a user on borrowed amount
	GetUserInterestAccrued(ctx context.Context, userAddress common.Address) (*big.Int, error)

	// GetUserTransactionHistory returns a user's borrowing transaction history
	GetUserTransactionHistory(ctx context.Context, userAddress common.Address, offset, limit int) ([]*models.Transaction, error)

	// CountUserTransactions counts the number of transactions for a user with optional filtering
	// CountUserTransactions counts the number of transactions for a user with optional filtering
	CountUserTransactions(ctx context.Context, address common.Address, filter map[string]any) (int64, error)
}
