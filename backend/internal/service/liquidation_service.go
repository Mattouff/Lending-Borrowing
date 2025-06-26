package service

import (
	"context"
	"math/big"

	"github.com/Mattouff/Lending-Borrowing/internal/domain/models"
	"github.com/Mattouff/Lending-Borrowing/internal/domain/repository"
	"github.com/Mattouff/Lending-Borrowing/internal/domain/service"
	"github.com/Mattouff/Lending-Borrowing/internal/infrastructure/blockchain/services"
)

type liquidationService struct {
	transactionRepo   repository.TransactionRepository
	userRepo          repository.UserRepository
	positionRepo      repository.PositionRepository
	collateral        *services.CollateralService
	borrowing         *services.BorrowingService
	collateralService service.CollateralService
}

// NewLiquidationService creates a new liquidation service
func NewLiquidationService(
	transactionRepo repository.TransactionRepository,
	userRepo repository.UserRepository,
	positionRepo repository.PositionRepository,
	collateralService service.CollateralService,
) (service.LiquidationService, error) {
	serviceFactory := services.GetInstance()

	collateral, err := serviceFactory.GetCollateralService()
	if err != nil {
		return nil, err
	}

	borrowing, err := serviceFactory.GetBorrowingService()
	if err != nil {
		return nil, err
	}

	return &liquidationService{
		transactionRepo:   transactionRepo,
		userRepo:          userRepo,
		positionRepo:      positionRepo,
		collateral:        collateral,
		borrowing:         borrowing,
		collateralService: collateralService,
	}, nil
}

// Liquidate allows liquidators to liquidate an under-collateralized position

// GetLiquidatablePositions returns all positions that can be liquidated
func (s *liquidationService) GetLiquidatablePositions(ctx context.Context) ([]*models.Position, error) {
	// Get liquidation threshold from contract
	threshold, err := s.collateral.GetLiquidationThreshold(ctx)
	if err != nil {
		return nil, err
	}

	// Find positions with health factor below liquidation threshold
	return s.positionRepo.FindAtRisk(ctx, threshold.String())
}

// GetLiquidationBonus returns the bonus a liquidator receives for liquidating a position
func (s *liquidationService) GetLiquidationBonus(ctx context.Context) (*big.Int, error) {
	return s.collateral.GetLiquidationBonus(ctx)
}

// GetLiquidationHistory returns the history of liquidation events
func (s *liquidationService) GetLiquidationHistory(ctx context.Context, offset, limit int) ([]*models.Transaction, error) {
	filter := map[string]any{
		"type": models.TransactionLiquidate,
	}

	return s.transactionRepo.List(ctx, filter, offset, limit)
}

// CountLiquidations counts the total number of liquidation transactions
func (s *liquidationService) CountLiquidations(ctx context.Context) (int64, error) {
	filter := map[string]any{
		"type": models.TransactionLiquidate,
	}

	// Count the liquidation transactions
	return s.transactionRepo.Count(ctx, filter)
}
