package services

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/Mattouff/Lending-Borrowing/internal/contracts/generated"
	"github.com/Mattouff/Lending-Borrowing/internal/infrastructure/blockchain"
)

// BorrowingService provides methods to interact with the Borrowing contract
type BorrowingService struct {
	client    *ethclient.Client
	contract  *generated.Borrowing
	address   common.Address
	ethClient *blockchain.EthClient
}

// NewBorrowingService creates a new instance of BorrowingService
func NewBorrowingService() (*BorrowingService, error) {
	ethClient := blockchain.GetInstance()
	client, err := ethClient.GetClient()
	if err != nil {
		return nil, err
	}

	address, err := ethClient.GetContractAddress("Borrowing")
	if err != nil {
		return nil, err
	}

	borrowingContract, err := generated.NewBorrowing(address, client)
	if err != nil {
		return nil, err
	}

	return &BorrowingService{
		client:    client,
		contract:  borrowingContract,
		address:   address,
		ethClient: ethClient,
	}, nil
}

// Borrow allows users to borrow tokens
func (s *BorrowingService) Borrow(auth *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return s.contract.Borrow(auth, amount)
}

// Repay allows users to repay borrowed tokens
func (s *BorrowingService) Repay(auth *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return s.contract.Repay(auth, amount)
}

// ReduceDebt allows liquidators to reduce a borrower's debt
func (s *BorrowingService) ReduceDebt(auth *bind.TransactOpts, borrower common.Address, amount *big.Int) (*types.Transaction, error) {
	return s.contract.ReduceDebt(auth, borrower, amount)
}

// GetBorrowToken returns the amount borrowed by a user
func (s *BorrowingService) GetBorrowToken(ctx context.Context, user common.Address) (*big.Int, error) {
	opts := &bind.CallOpts{Context: ctx}
	return s.contract.GetBorrowToken(opts, user)
}

// GetAllBorrowToken returns the total amount borrowed from the protocol
func (s *BorrowingService) GetAllBorrowToken(ctx context.Context) (*big.Int, error) {
	opts := &bind.CallOpts{Context: ctx}
	return s.contract.GetAllBorrowToken(opts)
}

// GetCurrentRate returns the current interest rate for borrowing
func (s *BorrowingService) GetCurrentRate(ctx context.Context) (*big.Int, error) {
	opts := &bind.CallOpts{Context: ctx}
	return s.contract.GetCurrentRate(opts)
}

// GetBorrowedPrincipal returns the principal amount borrowed by a user
func (s *BorrowingService) GetBorrowedPrincipal(ctx context.Context, user common.Address) (*big.Int, error) {
	opts := &bind.CallOpts{Context: ctx}
	return s.contract.BorrowedPrincipal(opts, user)
}

// GetTotalBorrowed returns the total amount borrowed from the protocol
func (s *BorrowingService) GetTotalBorrowed(ctx context.Context) (*big.Int, error) {
	opts := &bind.CallOpts{Context: ctx}
	return s.contract.TotalBorrowed(opts)
}

// GetMinInterestRate returns the minimum interest rate
func (s *BorrowingService) GetMinInterestRate(ctx context.Context) (*big.Int, error) {
	opts := &bind.CallOpts{Context: ctx}
	return s.contract.RMin(opts)
}

// GetMaxInterestRate returns the maximum interest rate
func (s *BorrowingService) GetMaxInterestRate(ctx context.Context) (*big.Int, error) {
	opts := &bind.CallOpts{Context: ctx}
	return s.contract.RMax(opts)
}

// ContractAddress returns the address of the borrowing contract
func (s *BorrowingService) ContractAddress() common.Address {
	return s.address
}
