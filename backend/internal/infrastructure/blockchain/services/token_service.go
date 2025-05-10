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

// TokenService provides methods to interact with ERC20 tokens
type TokenService struct {
	client    *ethclient.Client
	contract  *generated.Token
	address   common.Address
	ethClient *blockchain.EthClient
}

// NewTokenService creates a new instance of TokenService
func NewTokenService() (*TokenService, error) {
	ethClient := blockchain.GetInstance()
	client, err := ethClient.GetClient()
	if err != nil {
		return nil, err
	}

	address, err := ethClient.GetContractAddress("Token")
	if err != nil {
		return nil, err
	}

	tokenContract, err := generated.NewToken(address, client)
	if err != nil {
		return nil, err
	}

	return &TokenService{
		client:    client,
		contract:  tokenContract,
		address:   address,
		ethClient: ethClient,
	}, nil
}

// BalanceOf retrieves the token balance of a specific address
func (s *TokenService) BalanceOf(ctx context.Context, address common.Address) (*big.Int, error) {
	opts := &bind.CallOpts{Context: ctx}
	return s.contract.BalanceOf(opts, address)
}

// TotalSupply returns the total token supply
func (s *TokenService) TotalSupply(ctx context.Context) (*big.Int, error) {
	opts := &bind.CallOpts{Context: ctx}
	return s.contract.TotalSupply(opts)
}

// Name returns the name of the token
func (s *TokenService) Name(ctx context.Context) (string, error) {
	opts := &bind.CallOpts{Context: ctx}
	return s.contract.Name(opts)
}

// Symbol returns the symbol of the token
func (s *TokenService) Symbol(ctx context.Context) (string, error) {
	opts := &bind.CallOpts{Context: ctx}
	return s.contract.Symbol(opts)
}

// Decimals returns the number of decimals of the token
func (s *TokenService) Decimals(ctx context.Context) (uint8, error) {
	opts := &bind.CallOpts{Context: ctx}
	return s.contract.Decimals(opts)
}

// Approve approves the spender to transfer tokens on behalf of the sender
func (s *TokenService) Approve(auth *bind.TransactOpts, spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return s.contract.Approve(auth, spender, amount)
}

// Transfer transfers tokens to the specified address
func (s *TokenService) Transfer(auth *bind.TransactOpts, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return s.contract.Transfer(auth, to, amount)
}

// TransferFrom transfers tokens from one address to another using the allowance mechanism
func (s *TokenService) TransferFrom(auth *bind.TransactOpts, from common.Address, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return s.contract.TransferFrom(auth, from, to, amount)
}

// Allowance returns the remaining number of tokens that spender will be allowed to spend on behalf of owner
func (s *TokenService) Allowance(ctx context.Context, owner common.Address, spender common.Address) (*big.Int, error) {
	opts := &bind.CallOpts{Context: ctx}
	return s.contract.Allowance(opts, owner, spender)
}

// ContractAddress returns the address of the token contract
func (s *TokenService) ContractAddress() common.Address {
	return s.address
}
