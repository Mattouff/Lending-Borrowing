package service

import (
	"context"
	"math/big"

	"github.com/Mattouff/Lending-Borrowing/internal/domain/models"
	"github.com/ethereum/go-ethereum/common"
)

// LiquidationService defines the interface for liquidation business logic
type LiquidationService interface {
	// Liquidate allows liquidators to liquidate an under-collateralized position
	Liquidate(ctx context.Context, liquidatorAddress, borrowerAddress common.Address, repayAmount *big.Int) (string, error)

	// GetLiquidatablePositions returns all positions that can be liquidated
	GetLiquidatablePositions(ctx context.Context) ([]*models.Position, error)

	// GetLiquidationBonus returns the bonus a liquidator receives for liquidating a position
	GetLiquidationBonus(ctx context.Context) (*big.Int, error)

	// GetLiquidationHistory returns the history of liquidation events
	GetLiquidationHistory(ctx context.Context, offset, limit int) ([]*models.Transaction, error)

	// CountLiquidations counts the total number of liquidation transactions
	CountLiquidations(ctx context.Context) (int64, error)
}
