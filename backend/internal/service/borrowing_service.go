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
	filter := map[string]any{
		"user_id": user.ID,
		"type":    []models.TransactionType{models.TransactionBorrow, models.TransactionRepay},
	}

	return s.transactionRepo.List(ctx, filter, offset, limit)
}

// CountUserTransactions counts the number of transactions for a user with optional filtering
func (s *borrowingService) CountUserTransactions(ctx context.Context, address common.Address, filter map[string]any) (int64, error) {
	// Find the user by address
	user, err := s.userRepo.FindByAddress(ctx, address.Hex())
	if err != nil {
		return 0, err
	}

	if user == nil {
		return 0, errors.New("user not found")
	}

	// Add the user ID to the filter
	if filter == nil {
		filter = make(map[string]any)
	}
	filter["user_id"] = user.ID

	// Count the transactions
	return s.transactionRepo.Count(ctx, filter)
}

// PrepareBorrow prepares the contract interaction data for a borrow transaction
func (s *borrowingService) PrepareBorrow(ctx context.Context, userAddress common.Address, amount *big.Int) (*service.ContractInteractionData, error) {
	// Check if amount is valid
	if amount.Cmp(big.NewInt(0)) <= 0 {
		return nil, errors.New("borrow amount must be greater than 0")
	}

	// Check if user can borrow the requested amount
	maxBorrowable, err := s.collateralService.GetMaxBorrowableAmount(ctx, userAddress)
	if err != nil {
		return nil, err
	}

	if maxBorrowable.Cmp(amount) < 0 {
		return nil, errors.New("borrow amount exceeds maximum borrowable amount")
	}

	// Create contract interaction data
	contractData := &service.ContractInteractionData{
		ContractAddress: s.borrowing.ContractAddress().Hex(),
		MethodName:      "borrow",
		Parameters:      []interface{}{amount},
		Value:           "0",
	}

	return contractData, nil
}

// PrepareRepay prepares the contract interaction data for a repay transaction
func (s *borrowingService) PrepareRepay(ctx context.Context, userAddress common.Address, amount *big.Int) (*service.ContractInteractionData, error) {
	// Check if amount is valid
	if amount.Cmp(big.NewInt(0)) <= 0 {
		return nil, errors.New("repay amount must be greater than 0")
	}

	// Get current borrowed amount
	borrowed, err := s.GetBorrowedAmount(ctx, userAddress)
	if err != nil {
		return nil, err
	}

	// Check if repay amount is valid
	if borrowed.Cmp(amount) < 0 {
		return nil, errors.New("repay amount exceeds borrowed amount")
	}

	// Create contract interaction data
	contractData := &service.ContractInteractionData{
		ContractAddress: s.borrowing.ContractAddress().Hex(),
		MethodName:      "repay",
		Parameters:      []interface{}{amount},
		Value:           "0",
	}

	return contractData, nil
}

// ProcessWalletTransaction validates and processes a wallet-signed transaction
func (s *borrowingService) ProcessWalletTransaction(ctx context.Context, txHash string, userAddress common.Address, txType models.TransactionType, amount *big.Int) (*models.Transaction, error) {
	// Find the user
	user, err := s.userRepo.FindByAddress(ctx, userAddress.Hex())
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	// Check if transaction already exists
	existingTx, err := s.transactionRepo.FindByHash(ctx, txHash)
	if err != nil {
		return nil, err
	}

	if existingTx != nil {
		return existingTx, nil // Transaction already exists
	}

	// Create new transaction record
	transaction := &models.Transaction{
		UserID:       user.ID,
		Type:         txType,
		Status:       models.StatusPending,
		Hash:         txHash,
		Amount:       amount.String(),
		TokenAddress: s.borrowing.ContractAddress().Hex(),
	}

	if err := s.transactionRepo.Create(ctx, transaction); err != nil {
		return nil, err
	}

	// Update position based on transaction type
	switch txType {
	case models.TransactionBorrow:
		s.updatePositionAfterBorrow(ctx, user.ID, userAddress, amount)
	case models.TransactionRepay:
		s.updatePositionAfterRepay(ctx, user.ID, userAddress, amount)
	}

	return transaction, nil
}

// Helper method to update position after borrow
func (s *borrowingService) updatePositionAfterBorrow(ctx context.Context, userID uint, userAddress common.Address, amount *big.Int) error {
	// Get existing positions
	positions, err := s.positionRepo.FindActiveByUserID(ctx, userID)
	if err != nil {
		return err
	}

	// Get collateral balance
	collateralBalance, err := s.collateralService.GetCollateralBalance(ctx, userAddress)
	if err != nil {
		return err
	}

	// Get token address for underlying collateral
	serviceFactory := services.GetInstance()
	collateralService, err := serviceFactory.GetCollateralService()
	if err != nil {
		return err
	}

	if len(positions) == 0 {
		// Create new position
		position := &models.Position{
			UserID:           userID,
			CollateralAmount: collateralBalance.String(),
			CollateralToken:  collateralService.ContractAddress().Hex(),
			BorrowedAmount:   amount.String(),
			BorrowedToken:    s.borrowing.ContractAddress().Hex(),
			InterestRate:     "0", // Will be updated later
			Status:           models.StatusActive,
			HealthFactor:     "0", // Will be updated later
		}

		return s.positionRepo.Create(ctx, position)
	} else {
		// Update existing position
		position := positions[0]

		// Add new borrowed amount to existing
		currentBorrowed, success := new(big.Int).SetString(position.BorrowedAmount, 10)
		if !success {
			return errors.New("failed to parse borrowed amount")
		}

		newBorrowedAmount := new(big.Int).Add(currentBorrowed, amount)
		position.BorrowedAmount = newBorrowedAmount.String()
		position.CollateralAmount = collateralBalance.String()

		return s.positionRepo.Update(ctx, position)
	}
}

// Helper method to update position after repay
func (s *borrowingService) updatePositionAfterRepay(ctx context.Context, userID uint, userAddress common.Address, amount *big.Int) error {
	// Get existing positions
	positions, err := s.positionRepo.FindActiveByUserID(ctx, userID)
	if err != nil {
		return err
	}

	if len(positions) > 0 {
		position := positions[0]

		// Get current collateral balance
		collateralBalance, err := s.collateralService.GetCollateralBalance(ctx, userAddress)
		if err != nil {
			return err
		}

		// Subtract repaid amount from borrowed
		currentBorrowed, success := new(big.Int).SetString(position.BorrowedAmount, 10)
		if !success {
			return errors.New("failed to parse borrowed amount")
		}

		newBorrowedAmount := new(big.Int).Sub(currentBorrowed, amount)
		position.BorrowedAmount = newBorrowedAmount.String()
		position.CollateralAmount = collateralBalance.String()

		// If borrowed amount is 0, close the position
		if newBorrowedAmount.Cmp(big.NewInt(0)) == 0 {
			position.Status = models.StatusClosed
		}

		return s.positionRepo.Update(ctx, position)
	}

	return nil
}
