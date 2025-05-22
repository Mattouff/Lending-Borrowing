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

// CollateralService provides methods to interact with the Collateral contract
type CollateralService struct {
	client    *ethclient.Client
	contract  *generated.Collateral
	address   common.Address
	ethClient *blockchain.EthClient
}

// NewCollateralService creates a new instance of CollateralService
func NewCollateralService() (*CollateralService, error) {
	ethClient := blockchain.GetInstance()
	client, err := ethClient.GetClient()
	if err != nil {
		return nil, err
	}

	address, err := ethClient.GetContractAddress("Collateral")
	if err != nil {
		return nil, err
	}

	collateralContract, err := generated.NewCollateral(address, client)
	if err != nil {
		return nil, err
	}

	return &CollateralService{
		client:    client,
		contract:  collateralContract,
		address:   address,
		ethClient: ethClient,
	}, nil
}

// DepositCollateral allows users to deposit tokens as collateral
func (s *CollateralService) DepositCollateral(auth *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return s.contract.DepositCollateral(auth, amount)
}

// WithdrawCollateral allows users to withdraw tokens from their collateral
func (s *CollateralService) WithdrawCollateral(auth *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return s.contract.WithdrawCollateral(auth, amount)
}

// Liquidate allows liquidators to liquidate an under-collateralized position
func (s *CollateralService) Liquidate(auth *bind.TransactOpts, borrower common.Address, repayAmount *big.Int) (*types.Transaction, error) {
	return s.contract.Liquidate(auth, borrower, repayAmount)
}

// GetCollateralRatio returns the collateral ratio for a specific user
func (s *CollateralService) GetCollateralRatio(ctx context.Context, user common.Address) (*big.Int, error) {
	opts := &bind.CallOpts{Context: ctx}
	return s.contract.GetCollateralRatio(opts, user)
}

// GetMaxBorrowableAmount returns the maximum amount a user can borrow
func (s *CollateralService) GetMaxBorrowableAmount(ctx context.Context, user common.Address) (*big.Int, error) {
	opts := &bind.CallOpts{Context: ctx}
	return s.contract.GetMaxBorrowableAmount(opts, user)
}

// CanBorrow checks if a user can borrow a specific amount
func (s *CollateralService) CanBorrow(ctx context.Context, user common.Address, borrowAmount *big.Int) (bool, error) {
	opts := &bind.CallOpts{Context: ctx}
	return s.contract.CanBorrow(opts, user, borrowAmount)
}

// GetCollateralBalance returns the collateral balance of a user
func (s *CollateralService) GetCollateralBalance(ctx context.Context, user common.Address) (*big.Int, error) {
	opts := &bind.CallOpts{Context: ctx}
	return s.contract.CollateralBalance(opts, user)
}

// GetMinCollateralRatio returns the minimum collateral ratio required
func (s *CollateralService) GetMinCollateralRatio(ctx context.Context) (*big.Int, error) {
	opts := &bind.CallOpts{Context: ctx}
	return s.contract.MINCOLLATERALRATIO(opts)
}

// GetLiquidationThreshold returns the liquidation threshold
func (s *CollateralService) GetLiquidationThreshold(ctx context.Context) (*big.Int, error) {
	opts := &bind.CallOpts{Context: ctx}
	return s.contract.LIQUIDATIONTHRESHOLD(opts)
}

// GetLiquidationBonus returns the liquidation bonus
func (s *CollateralService) GetLiquidationBonus(ctx context.Context) (*big.Int, error) {
	opts := &bind.CallOpts{Context: ctx}
	return s.contract.LIQUIDATIONBONUS(opts)
}

// ContractAddress returns the address of the collateral contract
func (s *CollateralService) ContractAddress() common.Address {
	return s.address
}
