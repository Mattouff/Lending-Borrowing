package service

import (
	"context"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	"github.com/Mattouff/Lending-Borrowing/internal/domain/repository"
	"github.com/Mattouff/Lending-Borrowing/internal/domain/service"
	"github.com/Mattouff/Lending-Borrowing/internal/infrastructure/blockchain/services"
)

type collateralService struct {
	transactionRepo repository.TransactionRepository
	userRepo        repository.UserRepository
	positionRepo    repository.PositionRepository
	collateral      *services.CollateralService
}

// NewCollateralService creates a new collateral service
func NewCollateralService(
	transactionRepo repository.TransactionRepository,
	userRepo repository.UserRepository,
	positionRepo repository.PositionRepository,
) (service.CollateralService, error) {
	serviceFactory := services.GetInstance()
	collateral, err := serviceFactory.GetCollateralService()
	if err != nil {
		return nil, err
	}

	return &collateralService{
		transactionRepo: transactionRepo,
		userRepo:        userRepo,
		positionRepo:    positionRepo,
		collateral:      collateral,
	}, nil
}

// GetCollateralBalance returns the collateral balance of a user
func (s *collateralService) GetCollateralBalance(ctx context.Context, userAddress common.Address) (*big.Int, error) {
	return s.collateral.GetCollateralBalance(ctx, userAddress)
}

// GetTotalCollateral returns the total collateral in the protocol
func (s *collateralService) GetTotalCollateral(ctx context.Context) (*big.Int, error) {
	// This is just a placeholder. The actual implementation depends on your contract
	return big.NewInt(0), errors.New("not implemented in contract")
}

// GetCollateralRatio returns the collateral ratio for a user
func (s *collateralService) GetCollateralRatio(ctx context.Context, userAddress common.Address) (*big.Int, error) {
	return s.collateral.GetCollateralRatio(ctx, userAddress)
}

// GetMinCollateralRatio returns the minimum collateral ratio required
func (s *collateralService) GetMinCollateralRatio(ctx context.Context) (*big.Int, error) {
	return s.collateral.GetMinCollateralRatio(ctx)
}

// GetLiquidationThreshold returns the threshold at which a position can be liquidated
func (s *collateralService) GetLiquidationThreshold(ctx context.Context) (*big.Int, error) {
	return s.collateral.GetLiquidationThreshold(ctx)
}

// GetMaxBorrowableAmount returns the maximum amount a user can borrow based on their collateral
func (s *collateralService) GetMaxBorrowableAmount(ctx context.Context, userAddress common.Address) (*big.Int, error) {
	return s.collateral.GetMaxBorrowableAmount(ctx, userAddress)
}

// IsAtRisk checks if a user's position is at risk of liquidation
func (s *collateralService) IsAtRisk(ctx context.Context, userAddress common.Address) (bool, error) {
	// Check if user has any borrowings
	serviceFactory := services.GetInstance()
	borrowingService, err := serviceFactory.GetBorrowingService()
	if err != nil {
		return false, err
	}

	borrowedAmount, err := borrowingService.GetBorrowToken(ctx, userAddress)
	if err != nil {
		return false, err
	}

	// If no borrowings, not at risk
	if borrowedAmount.Cmp(big.NewInt(0)) == 0 {
		return false, nil
	}

	// Get current collateral ratio
	ratio, err := s.GetCollateralRatio(ctx, userAddress)
	if err != nil {
		return false, err
	}

	// Get liquidation threshold
	threshold, err := s.GetLiquidationThreshold(ctx)
	if err != nil {
		return false, err
	}

	// If ratio is below threshold, position is at risk
	return ratio.Cmp(threshold) <= 0, nil
}
