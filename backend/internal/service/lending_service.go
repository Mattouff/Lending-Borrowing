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

type lendingService struct {
	transactionRepo repository.TransactionRepository
	userRepo        repository.UserRepository
	lendingPool     *services.LendingPoolService
}

// NewLendingService creates a new lending service
func NewLendingService(
	transactionRepo repository.TransactionRepository,
	userRepo repository.UserRepository,
) (service.LendingService, error) {
	serviceFactory := services.GetInstance()
	lendingPool, err := serviceFactory.GetLendingPoolService()
	if err != nil {
		return nil, err
	}

	return &lendingService{
		transactionRepo: transactionRepo,
		userRepo:        userRepo,
		lendingPool:     lendingPool,
	}, nil
}

// GetUserBalance returns the user's balance in the lending pool
func (s *lendingService) GetUserBalance(ctx context.Context, userAddress common.Address) (*big.Int, error) {
	return s.lendingPool.GetLendingToken(ctx, userAddress)
}

// GetTotalDeposited returns the total amount deposited in the lending pool
func (s *lendingService) GetTotalDeposited(ctx context.Context) (*big.Int, error) {
	return s.lendingPool.GetTotalLending(ctx)
}

// GetCurrentInterestRate returns the current interest rate for lending
func (s *lendingService) GetCurrentInterestRate(ctx context.Context) (*big.Int, error) {
	return s.lendingPool.GetAnnualInterestRate(ctx)
}

// GetUserInterestEarned returns the interest earned by a user
// Note: This is a simplified implementation. In reality, you might need to calculate this based on time and rate
func (s *lendingService) GetUserInterestEarned(ctx context.Context, userAddress common.Address) (*big.Int, error) {
	// Implementation depends on your contract design
	// For this example, we'll just return 0 as a placeholder
	return big.NewInt(0), nil
}

// GetUserTransactionHistory returns a user's lending transaction history
func (s *lendingService) GetUserTransactionHistory(ctx context.Context, userAddress common.Address, offset, limit int) ([]*models.Transaction, error) {
	user, err := s.userRepo.FindByAddress(ctx, userAddress.Hex())
	if err != nil {
		return nil, err
	}

	// Get transactions with filter for deposit and withdraw types
	filter := map[string]any{
		"user_id": user.ID,
		"type":    []models.TransactionType{models.TransactionDeposit, models.TransactionWithdraw},
	}

	return s.transactionRepo.List(ctx, filter, offset, limit)
}

// PrepareDeposit prepares the contract interaction data for a deposit transaction
func (s *lendingService) PrepareDeposit(ctx context.Context, userAddress common.Address, amount *big.Int) (*service.ContractInteractionData, error) {
	// Check if amount is valid
	if amount.Cmp(big.NewInt(0)) <= 0 {
		return nil, errors.New("deposit amount must be greater than 0")
	}

	// Create contract interaction data
	contractData := &service.ContractInteractionData{
		ContractAddress: s.lendingPool.ContractAddress().Hex(),
		MethodName:      "deposit",
		Parameters:      []interface{}{amount},
		Value:           "0", // No ETH value needed for ERC20 deposits
	}

	return contractData, nil
}

// PrepareWithdraw prepares the contract interaction data for a withdraw transaction
func (s *lendingService) PrepareWithdraw(ctx context.Context, userAddress common.Address, amount *big.Int) (*service.ContractInteractionData, error) {
	// Check if amount is valid
	if amount.Cmp(big.NewInt(0)) <= 0 {
		return nil, errors.New("withdrawal amount must be greater than 0")
	}

	// Check if user has enough balance
	balance, err := s.GetUserBalance(ctx, userAddress)
	if err != nil {
		return nil, err
	}

	if balance.Cmp(amount) < 0 {
		return nil, errors.New("insufficient balance for withdrawal")
	}

	// Create contract interaction data
	contractData := &service.ContractInteractionData{
		ContractAddress: s.lendingPool.ContractAddress().Hex(),
		MethodName:      "withdraw",
		Parameters:      []interface{}{amount},
		Value:           "0",
	}

	return contractData, nil
}

// ProcessWalletTransaction validates and processes a wallet-signed transaction
func (s *lendingService) ProcessWalletTransaction(ctx context.Context, txHash string, userAddress common.Address, txType models.TransactionType, amount *big.Int) (*models.Transaction, error) {
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
		TokenAddress: s.lendingPool.ContractAddress().Hex(),
	}

	if err := s.transactionRepo.Create(ctx, transaction); err != nil {
		return nil, err
	}

	return transaction, nil
}

// CountUserTransactions counts the number of transactions for a user with optional filtering
func (s *lendingService) CountUserTransactions(ctx context.Context, address common.Address, filter map[string]any) (int64, error) {
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
