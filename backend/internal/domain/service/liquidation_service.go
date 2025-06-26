package service

import (
	"context"
	"math/big"

	"github.com/Mattouff/Lending-Borrowing/internal/domain/models"
)

// LiquidationService defines the interface for liquidation business logic
type LiquidationService interface {
	// GetLiquidatablePositions returns all positions that can be liquidated
	GetLiquidatablePositions(ctx context.Context) ([]*models.Position, error)

	// GetLiquidationBonus returns the bonus a liquidator receives for liquidating a position
	GetLiquidationBonus(ctx context.Context) (*big.Int, error)

	// GetLiquidationHistory returns the history of liquidation events
	GetLiquidationHistory(ctx context.Context, offset, limit int) ([]*models.Transaction, error)

	// CountLiquidations counts the total number of liquidation transactions
	CountLiquidations(ctx context.Context) (int64, error)
}
