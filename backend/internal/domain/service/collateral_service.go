package service

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// CollateralService defines the interface for collateral business logic
type CollateralService interface {
	// DepositCollateral allows users to deposit tokens as collateral
	DepositCollateral(ctx context.Context, userAddress common.Address, amount *big.Int) (string, error)

	// WithdrawCollateral allows users to withdraw tokens from their collateral
	WithdrawCollateral(ctx context.Context, userAddress common.Address, amount *big.Int) (string, error)

	// GetCollateralBalance returns the collateral balance of a user
	GetCollateralBalance(ctx context.Context, userAddress common.Address) (*big.Int, error)

	// GetTotalCollateral returns the total collateral in the protocol
	GetTotalCollateral(ctx context.Context) (*big.Int, error)

	// GetCollateralRatio returns the collateral ratio for a user
	GetCollateralRatio(ctx context.Context, userAddress common.Address) (*big.Int, error)

	// GetMinCollateralRatio returns the minimum collateral ratio required
	GetMinCollateralRatio(ctx context.Context) (*big.Int, error)

	// GetLiquidationThreshold returns the threshold at which a position can be liquidated
	GetLiquidationThreshold(ctx context.Context) (*big.Int, error)

	// GetMaxBorrowableAmount returns the maximum amount a user can borrow based on their collateral
	GetMaxBorrowableAmount(ctx context.Context, userAddress common.Address) (*big.Int, error)

	// IsAtRisk checks if a user's position is at risk of liquidation
	IsAtRisk(ctx context.Context, userAddress common.Address) (bool, error)
}
