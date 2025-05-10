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

// DepositCollateral allows users to deposit tokens as collateral
func (s *collateralService) DepositCollateral(ctx context.Context, userAddress common.Address, amount *big.Int) (string, error) {
	// Check if amount is valid
	if amount.Cmp(big.NewInt(0)) <= 0 {
		return "", errors.New("deposit amount must be greater than 0")
	}

	// Get auth for transaction
	auth, err := getTransactOpts(ctx, userAddress)
	if err != nil {
		return "", err
	}

	// Execute the deposit collateral transaction
	tx, err := s.collateral.DepositCollateral(auth, amount)
	if err != nil {
		return "", err
	}

	// Store transaction in database
	user, err := s.userRepo.FindByAddress(ctx, userAddress.Hex())
	if err != nil {
		return "", err
	}

	transaction := &models.Transaction{
		UserID:       user.ID,
		Type:         models.TransactionDeposit,
		Status:       models.StatusPending,
		Hash:         tx.Hash().Hex(),
		Amount:       amount.String(),
		TokenAddress: s.collateral.ContractAddress().Hex(),
	}

	if err := s.transactionRepo.Create(ctx, transaction); err != nil {
		return "", err
	}

	// Update or create position
	positions, err := s.positionRepo.FindActiveByUserID(ctx, user.ID)
	if err != nil {
		return "", err
	}

	// Get current collateral balance
	newCollateralBalance, err := s.GetCollateralBalance(ctx, userAddress)
	if err != nil {
		return "", err
	}

	if len(positions) == 0 {
		// Create new position if user doesn't have one
		position := &models.Position{
			UserID:           user.ID,
			CollateralAmount: newCollateralBalance.String(),
			CollateralToken:  s.collateral.ContractAddress().Hex(),
			BorrowedAmount:   "0", // No borrowing yet
			BorrowedToken:    "",  // Will be set when borrowing
			InterestRate:     "0", // Will be set when borrowing
			Status:           models.StatusActive,
			HealthFactor:     "0", // Will be calculated later
		}

		if err := s.positionRepo.Create(ctx, position); err != nil {
			return "", err
		}
	} else {
		// Update existing position
		position := positions[0]
		position.CollateralAmount = newCollateralBalance.String()

		// Recalculate health factor
		healthFactor, err := s.GetCollateralRatio(ctx, userAddress)
		if err == nil && healthFactor != nil {
			position.HealthFactor = healthFactor.String()
		}

		if err := s.positionRepo.Update(ctx, position); err != nil {
			return "", err
		}
	}

	return tx.Hash().Hex(), nil
}

// WithdrawCollateral allows users to withdraw tokens from their collateral
func (s *collateralService) WithdrawCollateral(ctx context.Context, userAddress common.Address, amount *big.Int) (string, error) {
	// Check if amount is valid
	if amount.Cmp(big.NewInt(0)) <= 0 {
		return "", errors.New("withdrawal amount must be greater than 0")
	}

	// Check if user has enough collateral
	balance, err := s.GetCollateralBalance(ctx, userAddress)
	if err != nil {
		return "", err
	}

	if balance.Cmp(amount) < 0 {
		return "", errors.New("insufficient collateral balance for withdrawal")
	}

	// Check if withdrawal would put user's position at risk
	// First get current borrowed amount
	serviceFactory := services.GetInstance()
	borrowingService, err := serviceFactory.GetBorrowingService()
	if err != nil {
		return "", err
	}

	borrowedAmount, err := borrowingService.GetBorrowToken(ctx, userAddress)
	if err != nil {
		return "", err
	}

	// If user has borrowed, check if withdrawal would affect collateral ratio
	if borrowedAmount.Cmp(big.NewInt(0)) > 0 {
		// Calculate new collateral balance after withdrawal
		newBalance := new(big.Int).Sub(balance, amount)

		// Calculate minimum required collateral based on borrowed amount and min ratio
		minRatio, err := s.GetMinCollateralRatio(ctx)
		if err != nil {
			return "", err
		}

		// minCollateral = borrowedAmount * minRatio / 10^18 (assuming minRatio is in 18 decimals)
		divisor := big.NewInt(10).Exp(big.NewInt(10), big.NewInt(18), nil)
		minCollateral := new(big.Int).Mul(borrowedAmount, minRatio)
		minCollateral = minCollateral.Div(minCollateral, divisor)

		if newBalance.Cmp(minCollateral) < 0 {
			return "", errors.New("withdrawal would put position at risk of liquidation")
		}
	}

	// Get auth for transaction
	auth, err := getTransactOpts(ctx, userAddress)
	if err != nil {
		return "", err
	}

	// Execute the withdraw collateral transaction
	tx, err := s.collateral.WithdrawCollateral(auth, amount)
	if err != nil {
		return "", err
	}

	// Store transaction in database
	user, err := s.userRepo.FindByAddress(ctx, userAddress.Hex())
	if err != nil {
		return "", err
	}

	transaction := &models.Transaction{
		UserID:       user.ID,
		Type:         models.TransactionWithdraw,
		Status:       models.StatusPending,
		Hash:         tx.Hash().Hex(),
		Amount:       amount.String(),
		TokenAddress: s.collateral.ContractAddress().Hex(),
	}

	if err := s.transactionRepo.Create(ctx, transaction); err != nil {
		return "", err
	}

	// Update position
	positions, err := s.positionRepo.FindActiveByUserID(ctx, user.ID)
	if err != nil {
		return "", err
	}

	if len(positions) > 0 {
		position := positions[0]

		// Get new collateral balance
		newCollateralBalance, err := s.GetCollateralBalance(ctx, userAddress)
		if err != nil {
			return "", err
		}

		position.CollateralAmount = newCollateralBalance.String()

		// Recalculate health factor
		healthFactor, err := s.GetCollateralRatio(ctx, userAddress)
		if err == nil && healthFactor != nil {
			position.HealthFactor = healthFactor.String()
		}

		if err := s.positionRepo.Update(ctx, position); err != nil {
			return "", err
		}

		// If no more collateral and no borrowings, close the position
		if newCollateralBalance.Cmp(big.NewInt(0)) == 0 && borrowedAmount.Cmp(big.NewInt(0)) == 0 {
			position.Status = models.StatusClosed
			if err := s.positionRepo.Update(ctx, position); err != nil {
				return "", err
			}
		}
	}

	return tx.Hash().Hex(), nil
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
