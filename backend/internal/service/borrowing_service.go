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

type borrowingService struct {
	transactionRepo   repository.TransactionRepository
	userRepo          repository.UserRepository
	positionRepo      repository.PositionRepository
	borrowing         *services.BorrowingService
	collateralService service.CollateralService
}

// NewBorrowingService creates a new borrowing service
func NewBorrowingService(
	transactionRepo repository.TransactionRepository,
	userRepo repository.UserRepository,
	positionRepo repository.PositionRepository,
	collateralService service.CollateralService,
) (service.BorrowingService, error) {
	serviceFactory := services.GetInstance()
	borrowing, err := serviceFactory.GetBorrowingService()
	if err != nil {
		return nil, err
	}

	return &borrowingService{
		transactionRepo:   transactionRepo,
		userRepo:          userRepo,
		positionRepo:      positionRepo,
		borrowing:         borrowing,
		collateralService: collateralService,
	}, nil
}

// Borrow allows users to borrow tokens based on their collateral
func (s *borrowingService) Borrow(ctx context.Context, userAddress common.Address, amount *big.Int) (string, error) {
	// Check if amount is valid
	if amount.Cmp(big.NewInt(0)) <= 0 {
		return "", errors.New("borrow amount must be greater than 0")
	}

	// Check if user can borrow the requested amount
	maxBorrowable, err := s.collateralService.GetMaxBorrowableAmount(ctx, userAddress)
	if err != nil {
		return "", err
	}

	if maxBorrowable.Cmp(amount) < 0 {
		return "", errors.New("borrow amount exceeds maximum borrowable amount")
	}

	// Get auth for transaction
	auth, err := getTransactOpts(ctx, userAddress)
	if err != nil {
		return "", err
	}

	// Execute the borrow transaction
	tx, err := s.borrowing.Borrow(auth, amount)
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
		Type:         models.TransactionBorrow,
		Status:       models.StatusPending,
		Hash:         tx.Hash().Hex(),
		Amount:       amount.String(),
		TokenAddress: s.borrowing.ContractAddress().Hex(),
	}

	if err := s.transactionRepo.Create(ctx, transaction); err != nil {
		return "", err
	}

	// Update or create position
	positions, err := s.positionRepo.FindActiveByUserID(ctx, user.ID)
	if err != nil {
		return "", err
	}

	// Get collateral balance
	collateralBalance, err := s.collateralService.GetCollateralBalance(ctx, userAddress)
	if err != nil {
		return "", err
	}

	// Get token address for underlying collateral
	serviceFactory := services.GetInstance()
	collateralService, err := serviceFactory.GetCollateralService()
	if err != nil {
		return "", err
	}

	if len(positions) == 0 {
		// Create new position
		position := &models.Position{
			UserID:           user.ID,
			CollateralAmount: collateralBalance.String(),
			CollateralToken:  collateralService.ContractAddress().Hex(),
			BorrowedAmount:   amount.String(),
			BorrowedToken:    s.borrowing.ContractAddress().Hex(),
			InterestRate:     "0", // Will be updated later
			Status:           models.StatusActive,
			HealthFactor:     "0", // Will be updated later
		}

		if err := s.positionRepo.Create(ctx, position); err != nil {
			return "", err
		}
	} else {
		// Update existing position
		position := positions[0]

		// Add new borrowed amount to existing
		currentBorrowed, success := new(big.Int).SetString(position.BorrowedAmount, 10)
		if !success {
			return "", errors.New("failed to parse borrowed amount")
		}

		newBorrowedAmount := new(big.Int).Add(currentBorrowed, amount)
		position.BorrowedAmount = newBorrowedAmount.String()
		position.CollateralAmount = collateralBalance.String()

		if err := s.positionRepo.Update(ctx, position); err != nil {
			return "", err
		}
	}

	return tx.Hash().Hex(), nil
}

// Repay allows users to repay borrowed tokens
func (s *borrowingService) Repay(ctx context.Context, userAddress common.Address, amount *big.Int) (string, error) {
	// Check if amount is valid
	if amount.Cmp(big.NewInt(0)) <= 0 {
		return "", errors.New("repay amount must be greater than 0")
	}

	// Get current borrowed amount
	borrowed, err := s.GetBorrowedAmount(ctx, userAddress)
	if err != nil {
		return "", err
	}

	// Check if repay amount is valid
	if borrowed.Cmp(amount) < 0 {
		return "", errors.New("repay amount exceeds borrowed amount")
	}

	// Get auth for transaction
	auth, err := getTransactOpts(ctx, userAddress)
	if err != nil {
		return "", err
	}

	// Execute the repay transaction
	tx, err := s.borrowing.Repay(auth, amount)
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
		Type:         models.TransactionRepay,
		Status:       models.StatusPending,
		Hash:         tx.Hash().Hex(),
		Amount:       amount.String(),
		TokenAddress: s.borrowing.ContractAddress().Hex(),
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

		// Subtract repaid amount from borrowed
		currentBorrowed, success := new(big.Int).SetString(position.BorrowedAmount, 10)
		if !success {
			return "", errors.New("failed to parse borrowed amount")
		}

		newBorrowedAmount := new(big.Int).Sub(currentBorrowed, amount)
		position.BorrowedAmount = newBorrowedAmount.String()

		// If borrowed amount is 0, close the position
		if newBorrowedAmount.Cmp(big.NewInt(0)) == 0 {
			position.Status = models.StatusClosed
		}

		if err := s.positionRepo.Update(ctx, position); err != nil {
			return "", err
		}
	}

	return tx.Hash().Hex(), nil
}

// GetBorrowedAmount returns the amount borrowed by a user
func (s *borrowingService) GetBorrowedAmount(ctx context.Context, userAddress common.Address) (*big.Int, error) {
	return s.borrowing.GetBorrowToken(ctx, userAddress)
}

// GetTotalBorrowed returns the total amount borrowed from the protocol
func (s *borrowingService) GetTotalBorrowed(ctx context.Context) (*big.Int, error) {
	return s.borrowing.GetTotalBorrowed(ctx)
}

// GetCurrentInterestRate returns the current interest rate for borrowing
func (s *borrowingService) GetCurrentInterestRate(ctx context.Context) (*big.Int, error) {
	return s.borrowing.GetCurrentRate(ctx)
}

// GetUserInterestAccrued returns the interest accrued by a user on borrowed amount
// Note: This is a simplified implementation. In reality, you might need to calculate this based on time and rate
func (s *borrowingService) GetUserInterestAccrued(ctx context.Context, userAddress common.Address) (*big.Int, error) {
	// Implementation depends on your contract design
	// For this example, we'll just return 0 as a placeholder
	return big.NewInt(0), nil
}

// GetUserTransactionHistory returns a user's borrowing transaction history
func (s *borrowingService) GetUserTransactionHistory(ctx context.Context, userAddress common.Address, offset, limit int) ([]*models.Transaction, error) {
	user, err := s.userRepo.FindByAddress(ctx, userAddress.Hex())
	if err != nil {
		return nil, err
	}

	// Get transactions with filter for borrow and repay types
	filter := map[string]interface{}{
		"user_id": user.ID,
		"type":    []models.TransactionType{models.TransactionBorrow, models.TransactionRepay},
	}

	return s.transactionRepo.List(ctx, filter, offset, limit)
}
