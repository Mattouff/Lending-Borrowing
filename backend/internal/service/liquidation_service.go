package service

import (
	"context"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/common"

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
func (s *liquidationService) Liquidate(ctx context.Context, liquidatorAddress, borrowerAddress common.Address, repayAmount *big.Int) (string, error) {
	// Check if amount is valid
	if repayAmount.Cmp(big.NewInt(0)) <= 0 {
		return "", errors.New("liquidation amount must be greater than 0")
	}

	// Check if the position is eligible for liquidation
	isAtRisk, err := s.collateralService.IsAtRisk(ctx, borrowerAddress)
	if err != nil {
		return "", err
	}

	if !isAtRisk {
		return "", errors.New("position is not eligible for liquidation")
	}

	// Get auth for transaction
	auth, err := getTransactOpts(ctx, liquidatorAddress)
	if err != nil {
		return "", err
	}

	// Execute the liquidation transaction
	tx, err := s.collateral.Liquidate(auth, borrowerAddress, repayAmount)
	if err != nil {
		return "", err
	}

	// Store transaction in database for both liquidator and borrower
	liquidator, err := s.userRepo.FindByAddress(ctx, liquidatorAddress.Hex())
	if err != nil {
		return "", err
	}

	borrower, err := s.userRepo.FindByAddress(ctx, borrowerAddress.Hex())
	if err != nil {
		return "", err
	}

	// Record liquidation transaction for liquidator
	liquidatorTx := &models.Transaction{
		UserID:       liquidator.ID,
		Type:         models.TransactionLiquidate,
		Status:       models.StatusPending,
		Hash:         tx.Hash().Hex(),
		Amount:       repayAmount.String(),
		TokenAddress: s.borrowing.ContractAddress().Hex(),
	}

	if err := s.transactionRepo.Create(ctx, liquidatorTx); err != nil {
		return "", err
	}

	// Record liquidation transaction for borrower
	borrowerTx := &models.Transaction{
		UserID:       borrower.ID,
		Type:         models.TransactionLiquidate,
		Status:       models.StatusPending,
		Hash:         tx.Hash().Hex(),
		Amount:       repayAmount.String(),
		TokenAddress: s.collateral.ContractAddress().Hex(),
	}

	if err := s.transactionRepo.Create(ctx, borrowerTx); err != nil {
		return "", err
	}

	// Update borrower's position
	positions, err := s.positionRepo.FindActiveByUserID(ctx, borrower.ID)
	if err != nil {
		return "", err
	}

	if len(positions) > 0 {
		position := positions[0]

		// Get updated collateral balance
		collateralBalance, err := s.collateral.GetCollateralBalance(ctx, borrowerAddress)
		if err != nil {
			return "", err
		}

		// Get updated borrowed amount
		borrowedAmount, err := s.borrowing.GetBorrowToken(ctx, borrowerAddress)
		if err != nil {
			return "", err
		}

		// Update position data
		position.CollateralAmount = collateralBalance.String()
		position.BorrowedAmount = borrowedAmount.String()

		// Recalculate health factor
		healthFactor, err := s.collateral.GetCollateralRatio(ctx, borrowerAddress)
		if err == nil && healthFactor != nil {
			position.HealthFactor = healthFactor.String()
		}

		// If collateral or borrowed amount is 0, mark position as liquidated
		if collateralBalance.Cmp(big.NewInt(0)) == 0 || borrowedAmount.Cmp(big.NewInt(0)) == 0 {
			position.Status = models.StatusLiquidated
		}

		if err := s.positionRepo.Update(ctx, position); err != nil {
			return "", err
		}
	}

	return tx.Hash().Hex(), nil
}

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
	filter := map[string]interface{}{
		"type": models.TransactionLiquidate,
	}

	return s.transactionRepo.List(ctx, filter, offset, limit)
}
