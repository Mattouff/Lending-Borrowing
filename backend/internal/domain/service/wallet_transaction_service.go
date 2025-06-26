package service

import (
	"context"
	"math/big"

	"github.com/Mattouff/Lending-Borrowing/internal/domain/models"
	"github.com/ethereum/go-ethereum/common"
)

// WalletTransactionService defines the interface for wallet-based transaction handling
type WalletTransactionService interface {
	// PrepareDeposit prepares the contract interaction data for a deposit transaction
	PrepareDeposit(ctx context.Context, userAddress common.Address, amount *big.Int) (*ContractInteractionData, error)

	// PrepareWithdraw prepares the contract interaction data for a withdraw transaction
	PrepareWithdraw(ctx context.Context, userAddress common.Address, amount *big.Int) (*ContractInteractionData, error)

	// PrepareBorrow prepares the contract interaction data for a borrow transaction
	PrepareBorrow(ctx context.Context, userAddress common.Address, amount *big.Int) (*ContractInteractionData, error)

	// PrepareRepay prepares the contract interaction data for a repay transaction
	PrepareRepay(ctx context.Context, userAddress common.Address, amount *big.Int) (*ContractInteractionData, error)

	// PrepareCollateralDeposit prepares the contract interaction data for a collateral deposit transaction
	PrepareCollateralDeposit(ctx context.Context, userAddress common.Address, amount *big.Int) (*ContractInteractionData, error)

	// PrepareCollateralWithdraw prepares the contract interaction data for a collateral withdraw transaction
	PrepareCollateralWithdraw(ctx context.Context, userAddress common.Address, amount *big.Int) (*ContractInteractionData, error)

	// ValidateAndProcessTransaction validates a signed transaction and updates the database
	ValidateAndProcessTransaction(ctx context.Context, txHash string, userAddress common.Address, txType models.TransactionType, amount *big.Int, tokenAddress common.Address) (*models.Transaction, error)

	// GetTransactionStatus gets the current status of a transaction from the blockchain
	GetTransactionStatus(ctx context.Context, txHash string) (models.TransactionStatus, error)
}
