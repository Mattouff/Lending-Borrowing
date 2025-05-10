// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package generated

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// LendingPoolMetaData contains all meta data concerning the LendingPool contract.
var LendingPoolMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"_underlying\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_annualInterestRate\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"allowance\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"spender\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"annualInterestRate\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"approve\",\"inputs\":[{\"name\":\"spender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"balanceOf\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"decimals\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"deposit\",\"inputs\":[{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getAllLendingToken\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getLendingToken\",\"inputs\":[{\"name\":\"user\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"lastUpdate\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"lendingBalance\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"name\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"symbol\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"totalLending\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"totalSupply\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"transfer\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferFrom\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"underlying\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractToken\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"updateUserInterest\",\"inputs\":[{\"name\":\"user\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"withdraw\",\"inputs\":[{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"Approval\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"spender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Transfer\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"ERC20InsufficientAllowance\",\"inputs\":[{\"name\":\"spender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"allowance\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"needed\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ERC20InsufficientBalance\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"balance\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"needed\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ERC20InvalidApprover\",\"inputs\":[{\"name\":\"approver\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC20InvalidReceiver\",\"inputs\":[{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC20InvalidSender\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC20InvalidSpender\",\"inputs\":[{\"name\":\"spender\",\"type\":\"address\",\"internalType\":\"address\"}]}]",
}

// LendingPoolABI is the input ABI used to generate the binding from.
// Deprecated: Use LendingPoolMetaData.ABI instead.
var LendingPoolABI = LendingPoolMetaData.ABI

// LendingPool is an auto generated Go binding around an Ethereum contract.
type LendingPool struct {
	LendingPoolCaller     // Read-only binding to the contract
	LendingPoolTransactor // Write-only binding to the contract
	LendingPoolFilterer   // Log filterer for contract events
}

// LendingPoolCaller is an auto generated read-only Go binding around an Ethereum contract.
type LendingPoolCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LendingPoolTransactor is an auto generated write-only Go binding around an Ethereum contract.
type LendingPoolTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LendingPoolFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type LendingPoolFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LendingPoolSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type LendingPoolSession struct {
	Contract     *LendingPool      // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// LendingPoolCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type LendingPoolCallerSession struct {
	Contract *LendingPoolCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts      // Call options to use throughout this session
}

// LendingPoolTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type LendingPoolTransactorSession struct {
	Contract     *LendingPoolTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// LendingPoolRaw is an auto generated low-level Go binding around an Ethereum contract.
type LendingPoolRaw struct {
	Contract *LendingPool // Generic contract binding to access the raw methods on
}

// LendingPoolCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type LendingPoolCallerRaw struct {
	Contract *LendingPoolCaller // Generic read-only contract binding to access the raw methods on
}

// LendingPoolTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type LendingPoolTransactorRaw struct {
	Contract *LendingPoolTransactor // Generic write-only contract binding to access the raw methods on
}

// NewLendingPool creates a new instance of LendingPool, bound to a specific deployed contract.
func NewLendingPool(address common.Address, backend bind.ContractBackend) (*LendingPool, error) {
	contract, err := bindLendingPool(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &LendingPool{LendingPoolCaller: LendingPoolCaller{contract: contract}, LendingPoolTransactor: LendingPoolTransactor{contract: contract}, LendingPoolFilterer: LendingPoolFilterer{contract: contract}}, nil
}

// NewLendingPoolCaller creates a new read-only instance of LendingPool, bound to a specific deployed contract.
func NewLendingPoolCaller(address common.Address, caller bind.ContractCaller) (*LendingPoolCaller, error) {
	contract, err := bindLendingPool(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &LendingPoolCaller{contract: contract}, nil
}

// NewLendingPoolTransactor creates a new write-only instance of LendingPool, bound to a specific deployed contract.
func NewLendingPoolTransactor(address common.Address, transactor bind.ContractTransactor) (*LendingPoolTransactor, error) {
	contract, err := bindLendingPool(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &LendingPoolTransactor{contract: contract}, nil
}

// NewLendingPoolFilterer creates a new log filterer instance of LendingPool, bound to a specific deployed contract.
func NewLendingPoolFilterer(address common.Address, filterer bind.ContractFilterer) (*LendingPoolFilterer, error) {
	contract, err := bindLendingPool(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &LendingPoolFilterer{contract: contract}, nil
}

// bindLendingPool binds a generic wrapper to an already deployed contract.
func bindLendingPool(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := LendingPoolMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_LendingPool *LendingPoolRaw) Call(opts *bind.CallOpts, result *[]any, method string, params ...any) error {
	return _LendingPool.Contract.LendingPoolCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_LendingPool *LendingPoolRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LendingPool.Contract.LendingPoolTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_LendingPool *LendingPoolRaw) Transact(opts *bind.TransactOpts, method string, params ...any) (*types.Transaction, error) {
	return _LendingPool.Contract.LendingPoolTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_LendingPool *LendingPoolCallerRaw) Call(opts *bind.CallOpts, result *[]any, method string, params ...any) error {
	return _LendingPool.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_LendingPool *LendingPoolTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LendingPool.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_LendingPool *LendingPoolTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...any) (*types.Transaction, error) {
	return _LendingPool.Contract.contract.Transact(opts, method, params...)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_LendingPool *LendingPoolCaller) Allowance(opts *bind.CallOpts, owner common.Address, spender common.Address) (*big.Int, error) {
	var out []any
	err := _LendingPool.contract.Call(opts, &out, "allowance", owner, spender)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_LendingPool *LendingPoolSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _LendingPool.Contract.Allowance(&_LendingPool.CallOpts, owner, spender)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_LendingPool *LendingPoolCallerSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _LendingPool.Contract.Allowance(&_LendingPool.CallOpts, owner, spender)
}

// AnnualInterestRate is a free data retrieval call binding the contract method 0x1a9703da.
//
// Solidity: function annualInterestRate() view returns(uint256)
func (_LendingPool *LendingPoolCaller) AnnualInterestRate(opts *bind.CallOpts) (*big.Int, error) {
	var out []any
	err := _LendingPool.contract.Call(opts, &out, "annualInterestRate")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// AnnualInterestRate is a free data retrieval call binding the contract method 0x1a9703da.
//
// Solidity: function annualInterestRate() view returns(uint256)
func (_LendingPool *LendingPoolSession) AnnualInterestRate() (*big.Int, error) {
	return _LendingPool.Contract.AnnualInterestRate(&_LendingPool.CallOpts)
}

// AnnualInterestRate is a free data retrieval call binding the contract method 0x1a9703da.
//
// Solidity: function annualInterestRate() view returns(uint256)
func (_LendingPool *LendingPoolCallerSession) AnnualInterestRate() (*big.Int, error) {
	return _LendingPool.Contract.AnnualInterestRate(&_LendingPool.CallOpts)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_LendingPool *LendingPoolCaller) BalanceOf(opts *bind.CallOpts, account common.Address) (*big.Int, error) {
	var out []any
	err := _LendingPool.contract.Call(opts, &out, "balanceOf", account)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_LendingPool *LendingPoolSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _LendingPool.Contract.BalanceOf(&_LendingPool.CallOpts, account)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_LendingPool *LendingPoolCallerSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _LendingPool.Contract.BalanceOf(&_LendingPool.CallOpts, account)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_LendingPool *LendingPoolCaller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []any
	err := _LendingPool.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_LendingPool *LendingPoolSession) Decimals() (uint8, error) {
	return _LendingPool.Contract.Decimals(&_LendingPool.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_LendingPool *LendingPoolCallerSession) Decimals() (uint8, error) {
	return _LendingPool.Contract.Decimals(&_LendingPool.CallOpts)
}

// GetAllLendingToken is a free data retrieval call binding the contract method 0xb64e0cb0.
//
// Solidity: function getAllLendingToken() view returns(uint256)
func (_LendingPool *LendingPoolCaller) GetAllLendingToken(opts *bind.CallOpts) (*big.Int, error) {
	var out []any
	err := _LendingPool.contract.Call(opts, &out, "getAllLendingToken")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetAllLendingToken is a free data retrieval call binding the contract method 0xb64e0cb0.
//
// Solidity: function getAllLendingToken() view returns(uint256)
func (_LendingPool *LendingPoolSession) GetAllLendingToken() (*big.Int, error) {
	return _LendingPool.Contract.GetAllLendingToken(&_LendingPool.CallOpts)
}

// GetAllLendingToken is a free data retrieval call binding the contract method 0xb64e0cb0.
//
// Solidity: function getAllLendingToken() view returns(uint256)
func (_LendingPool *LendingPoolCallerSession) GetAllLendingToken() (*big.Int, error) {
	return _LendingPool.Contract.GetAllLendingToken(&_LendingPool.CallOpts)
}

// GetLendingToken is a free data retrieval call binding the contract method 0x10c76077.
//
// Solidity: function getLendingToken(address user) view returns(uint256)
func (_LendingPool *LendingPoolCaller) GetLendingToken(opts *bind.CallOpts, user common.Address) (*big.Int, error) {
	var out []any
	err := _LendingPool.contract.Call(opts, &out, "getLendingToken", user)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetLendingToken is a free data retrieval call binding the contract method 0x10c76077.
//
// Solidity: function getLendingToken(address user) view returns(uint256)
func (_LendingPool *LendingPoolSession) GetLendingToken(user common.Address) (*big.Int, error) {
	return _LendingPool.Contract.GetLendingToken(&_LendingPool.CallOpts, user)
}

// GetLendingToken is a free data retrieval call binding the contract method 0x10c76077.
//
// Solidity: function getLendingToken(address user) view returns(uint256)
func (_LendingPool *LendingPoolCallerSession) GetLendingToken(user common.Address) (*big.Int, error) {
	return _LendingPool.Contract.GetLendingToken(&_LendingPool.CallOpts, user)
}

// LastUpdate is a free data retrieval call binding the contract method 0xcb03fb1e.
//
// Solidity: function lastUpdate(address ) view returns(uint256)
func (_LendingPool *LendingPoolCaller) LastUpdate(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []any
	err := _LendingPool.contract.Call(opts, &out, "lastUpdate", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LastUpdate is a free data retrieval call binding the contract method 0xcb03fb1e.
//
// Solidity: function lastUpdate(address ) view returns(uint256)
func (_LendingPool *LendingPoolSession) LastUpdate(arg0 common.Address) (*big.Int, error) {
	return _LendingPool.Contract.LastUpdate(&_LendingPool.CallOpts, arg0)
}

// LastUpdate is a free data retrieval call binding the contract method 0xcb03fb1e.
//
// Solidity: function lastUpdate(address ) view returns(uint256)
func (_LendingPool *LendingPoolCallerSession) LastUpdate(arg0 common.Address) (*big.Int, error) {
	return _LendingPool.Contract.LastUpdate(&_LendingPool.CallOpts, arg0)
}

// LendingBalance is a free data retrieval call binding the contract method 0x7f10883b.
//
// Solidity: function lendingBalance(address ) view returns(uint256)
func (_LendingPool *LendingPoolCaller) LendingBalance(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []any
	err := _LendingPool.contract.Call(opts, &out, "lendingBalance", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LendingBalance is a free data retrieval call binding the contract method 0x7f10883b.
//
// Solidity: function lendingBalance(address ) view returns(uint256)
func (_LendingPool *LendingPoolSession) LendingBalance(arg0 common.Address) (*big.Int, error) {
	return _LendingPool.Contract.LendingBalance(&_LendingPool.CallOpts, arg0)
}

// LendingBalance is a free data retrieval call binding the contract method 0x7f10883b.
//
// Solidity: function lendingBalance(address ) view returns(uint256)
func (_LendingPool *LendingPoolCallerSession) LendingBalance(arg0 common.Address) (*big.Int, error) {
	return _LendingPool.Contract.LendingBalance(&_LendingPool.CallOpts, arg0)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_LendingPool *LendingPoolCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []any
	err := _LendingPool.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_LendingPool *LendingPoolSession) Name() (string, error) {
	return _LendingPool.Contract.Name(&_LendingPool.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_LendingPool *LendingPoolCallerSession) Name() (string, error) {
	return _LendingPool.Contract.Name(&_LendingPool.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_LendingPool *LendingPoolCaller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []any
	err := _LendingPool.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_LendingPool *LendingPoolSession) Symbol() (string, error) {
	return _LendingPool.Contract.Symbol(&_LendingPool.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_LendingPool *LendingPoolCallerSession) Symbol() (string, error) {
	return _LendingPool.Contract.Symbol(&_LendingPool.CallOpts)
}

// TotalLending is a free data retrieval call binding the contract method 0xd79e738e.
//
// Solidity: function totalLending() view returns(uint256)
func (_LendingPool *LendingPoolCaller) TotalLending(opts *bind.CallOpts) (*big.Int, error) {
	var out []any
	err := _LendingPool.contract.Call(opts, &out, "totalLending")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalLending is a free data retrieval call binding the contract method 0xd79e738e.
//
// Solidity: function totalLending() view returns(uint256)
func (_LendingPool *LendingPoolSession) TotalLending() (*big.Int, error) {
	return _LendingPool.Contract.TotalLending(&_LendingPool.CallOpts)
}

// TotalLending is a free data retrieval call binding the contract method 0xd79e738e.
//
// Solidity: function totalLending() view returns(uint256)
func (_LendingPool *LendingPoolCallerSession) TotalLending() (*big.Int, error) {
	return _LendingPool.Contract.TotalLending(&_LendingPool.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_LendingPool *LendingPoolCaller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []any
	err := _LendingPool.contract.Call(opts, &out, "totalSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_LendingPool *LendingPoolSession) TotalSupply() (*big.Int, error) {
	return _LendingPool.Contract.TotalSupply(&_LendingPool.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_LendingPool *LendingPoolCallerSession) TotalSupply() (*big.Int, error) {
	return _LendingPool.Contract.TotalSupply(&_LendingPool.CallOpts)
}

// Underlying is a free data retrieval call binding the contract method 0x6f307dc3.
//
// Solidity: function underlying() view returns(address)
func (_LendingPool *LendingPoolCaller) Underlying(opts *bind.CallOpts) (common.Address, error) {
	var out []any
	err := _LendingPool.contract.Call(opts, &out, "underlying")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Underlying is a free data retrieval call binding the contract method 0x6f307dc3.
//
// Solidity: function underlying() view returns(address)
func (_LendingPool *LendingPoolSession) Underlying() (common.Address, error) {
	return _LendingPool.Contract.Underlying(&_LendingPool.CallOpts)
}

// Underlying is a free data retrieval call binding the contract method 0x6f307dc3.
//
// Solidity: function underlying() view returns(address)
func (_LendingPool *LendingPoolCallerSession) Underlying() (common.Address, error) {
	return _LendingPool.Contract.Underlying(&_LendingPool.CallOpts)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_LendingPool *LendingPoolTransactor) Approve(opts *bind.TransactOpts, spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _LendingPool.contract.Transact(opts, "approve", spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_LendingPool *LendingPoolSession) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _LendingPool.Contract.Approve(&_LendingPool.TransactOpts, spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_LendingPool *LendingPoolTransactorSession) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _LendingPool.Contract.Approve(&_LendingPool.TransactOpts, spender, value)
}

// Deposit is a paid mutator transaction binding the contract method 0xb6b55f25.
//
// Solidity: function deposit(uint256 amount) returns()
func (_LendingPool *LendingPoolTransactor) Deposit(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _LendingPool.contract.Transact(opts, "deposit", amount)
}

// Deposit is a paid mutator transaction binding the contract method 0xb6b55f25.
//
// Solidity: function deposit(uint256 amount) returns()
func (_LendingPool *LendingPoolSession) Deposit(amount *big.Int) (*types.Transaction, error) {
	return _LendingPool.Contract.Deposit(&_LendingPool.TransactOpts, amount)
}

// Deposit is a paid mutator transaction binding the contract method 0xb6b55f25.
//
// Solidity: function deposit(uint256 amount) returns()
func (_LendingPool *LendingPoolTransactorSession) Deposit(amount *big.Int) (*types.Transaction, error) {
	return _LendingPool.Contract.Deposit(&_LendingPool.TransactOpts, amount)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_LendingPool *LendingPoolTransactor) Transfer(opts *bind.TransactOpts, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _LendingPool.contract.Transact(opts, "transfer", to, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_LendingPool *LendingPoolSession) Transfer(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _LendingPool.Contract.Transfer(&_LendingPool.TransactOpts, to, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_LendingPool *LendingPoolTransactorSession) Transfer(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _LendingPool.Contract.Transfer(&_LendingPool.TransactOpts, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_LendingPool *LendingPoolTransactor) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _LendingPool.contract.Transact(opts, "transferFrom", from, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_LendingPool *LendingPoolSession) TransferFrom(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _LendingPool.Contract.TransferFrom(&_LendingPool.TransactOpts, from, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_LendingPool *LendingPoolTransactorSession) TransferFrom(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _LendingPool.Contract.TransferFrom(&_LendingPool.TransactOpts, from, to, value)
}

// UpdateUserInterest is a paid mutator transaction binding the contract method 0xafb442e9.
//
// Solidity: function updateUserInterest(address user) returns()
func (_LendingPool *LendingPoolTransactor) UpdateUserInterest(opts *bind.TransactOpts, user common.Address) (*types.Transaction, error) {
	return _LendingPool.contract.Transact(opts, "updateUserInterest", user)
}

// UpdateUserInterest is a paid mutator transaction binding the contract method 0xafb442e9.
//
// Solidity: function updateUserInterest(address user) returns()
func (_LendingPool *LendingPoolSession) UpdateUserInterest(user common.Address) (*types.Transaction, error) {
	return _LendingPool.Contract.UpdateUserInterest(&_LendingPool.TransactOpts, user)
}

// UpdateUserInterest is a paid mutator transaction binding the contract method 0xafb442e9.
//
// Solidity: function updateUserInterest(address user) returns()
func (_LendingPool *LendingPoolTransactorSession) UpdateUserInterest(user common.Address) (*types.Transaction, error) {
	return _LendingPool.Contract.UpdateUserInterest(&_LendingPool.TransactOpts, user)
}

// Withdraw is a paid mutator transaction binding the contract method 0x2e1a7d4d.
//
// Solidity: function withdraw(uint256 amount) returns()
func (_LendingPool *LendingPoolTransactor) Withdraw(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _LendingPool.contract.Transact(opts, "withdraw", amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0x2e1a7d4d.
//
// Solidity: function withdraw(uint256 amount) returns()
func (_LendingPool *LendingPoolSession) Withdraw(amount *big.Int) (*types.Transaction, error) {
	return _LendingPool.Contract.Withdraw(&_LendingPool.TransactOpts, amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0x2e1a7d4d.
//
// Solidity: function withdraw(uint256 amount) returns()
func (_LendingPool *LendingPoolTransactorSession) Withdraw(amount *big.Int) (*types.Transaction, error) {
	return _LendingPool.Contract.Withdraw(&_LendingPool.TransactOpts, amount)
}

// LendingPoolApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the LendingPool contract.
type LendingPoolApprovalIterator struct {
	Event *LendingPoolApproval // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *LendingPoolApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LendingPoolApproval)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(LendingPoolApproval)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *LendingPoolApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LendingPoolApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LendingPoolApproval represents a Approval event raised by the LendingPool contract.
type LendingPoolApproval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_LendingPool *LendingPoolFilterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*LendingPoolApprovalIterator, error) {

	var ownerRule []any
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []any
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _LendingPool.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &LendingPoolApprovalIterator{contract: _LendingPool.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_LendingPool *LendingPoolFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *LendingPoolApproval, owner []common.Address, spender []common.Address) (event.Subscription, error) {

	var ownerRule []any
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []any
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _LendingPool.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LendingPoolApproval)
				if err := _LendingPool.contract.UnpackLog(event, "Approval", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseApproval is a log parse operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_LendingPool *LendingPoolFilterer) ParseApproval(log types.Log) (*LendingPoolApproval, error) {
	event := new(LendingPoolApproval)
	if err := _LendingPool.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// LendingPoolTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the LendingPool contract.
type LendingPoolTransferIterator struct {
	Event *LendingPoolTransfer // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *LendingPoolTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LendingPoolTransfer)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(LendingPoolTransfer)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *LendingPoolTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LendingPoolTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LendingPoolTransfer represents a Transfer event raised by the LendingPool contract.
type LendingPoolTransfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_LendingPool *LendingPoolFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*LendingPoolTransferIterator, error) {

	var fromRule []any
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []any
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _LendingPool.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &LendingPoolTransferIterator{contract: _LendingPool.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_LendingPool *LendingPoolFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *LendingPoolTransfer, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []any
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []any
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _LendingPool.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LendingPoolTransfer)
				if err := _LendingPool.contract.UnpackLog(event, "Transfer", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseTransfer is a log parse operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_LendingPool *LendingPoolFilterer) ParseTransfer(log types.Log) (*LendingPoolTransfer, error) {
	event := new(LendingPoolTransfer)
	if err := _LendingPool.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
