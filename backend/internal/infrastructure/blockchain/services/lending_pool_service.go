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

// LendingPoolService provides methods to interact with the LendingPool contract
type LendingPoolService struct {
	client    *ethclient.Client
	contract  *generated.LendingPool
	address   common.Address
	ethClient *blockchain.EthClient
}

// NewLendingPoolService creates a new instance of LendingPoolService
func NewLendingPoolService() (*LendingPoolService, error) {
	ethClient := blockchain.GetInstance()
	client, err := ethClient.GetClient()
	if err != nil {
		return nil, err
	}

	address, err := ethClient.GetContractAddress("LendingPool")
	if err != nil {
		return nil, err
	}

	lendingPoolContract, err := generated.NewLendingPool(address, client)
	if err != nil {
		return nil, err
	}

	return &LendingPoolService{
		client:    client,
		contract:  lendingPoolContract,
		address:   address,
		ethClient: ethClient,
	}, nil
}

// Deposit allows users to deposit tokens into the lending pool
func (s *LendingPoolService) Deposit(auth *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return s.contract.Deposit(auth, amount)
}

// Withdraw allows users to withdraw tokens from the lending pool
func (s *LendingPoolService) Withdraw(auth *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return s.contract.Withdraw(auth, amount)
}

// GetLendingToken returns the user's lent token amount
func (s *LendingPoolService) GetLendingToken(ctx context.Context, user common.Address) (*big.Int, error) {
	opts := &bind.CallOpts{Context: ctx}
	return s.contract.GetLendingToken(opts, user)
}

// GetAllLendingToken returns the total amount of tokens lent to the pool
func (s *LendingPoolService) GetAllLendingToken(ctx context.Context) (*big.Int, error) {
	opts := &bind.CallOpts{Context: ctx}
	return s.contract.GetAllLendingToken(opts)
}

// GetAnnualInterestRate returns the current annual interest rate
func (s *LendingPoolService) GetAnnualInterestRate(ctx context.Context) (*big.Int, error) {
	opts := &bind.CallOpts{Context: ctx}
	return s.contract.AnnualInterestRate(opts)
}

// UpdateUserInterest updates the user's interest earnings
func (s *LendingPoolService) UpdateUserInterest(auth *bind.TransactOpts, user common.Address) (*types.Transaction, error) {
	return s.contract.UpdateUserInterest(auth, user)
}

// GetTotalLending returns the total amount of tokens lent
func (s *LendingPoolService) GetTotalLending(ctx context.Context) (*big.Int, error) {
	opts := &bind.CallOpts{Context: ctx}
	return s.contract.TotalLending(opts)
}

// GetUnderlying returns the address of the underlying token
func (s *LendingPoolService) GetUnderlying(ctx context.Context) (common.Address, error) {
	opts := &bind.CallOpts{Context: ctx}
	return s.contract.Underlying(opts)
}

// ContractAddress returns the address of the lending pool contract
func (s *LendingPoolService) ContractAddress() common.Address {
	return s.address
}
