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

// BorrowingMetaData contains all meta data concerning the Borrowing contract.
var BorrowingMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"_token\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_collateral\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_rMin\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_rMax\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_beta\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"beta\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"borrow\",\"inputs\":[{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"borrowedPrincipal\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"collateral\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractCollateral\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getAllBorrowToken\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getBorrowToken\",\"inputs\":[{\"name\":\"user\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getCurrentRate\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"lastUpdateTime\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"rMax\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"rMin\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"reduceDebt\",\"inputs\":[{\"name\":\"borrower\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"repay\",\"inputs\":[{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"token\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractToken\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"totalBorrowed\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"}]",
}

// BorrowingABI is the input ABI used to generate the binding from.
// Deprecated: Use BorrowingMetaData.ABI instead.
var BorrowingABI = BorrowingMetaData.ABI

// Borrowing is an auto generated Go binding around an Ethereum contract.
type Borrowing struct {
	BorrowingCaller     // Read-only binding to the contract
	BorrowingTransactor // Write-only binding to the contract
	BorrowingFilterer   // Log filterer for contract events
}

// BorrowingCaller is an auto generated read-only Go binding around an Ethereum contract.
type BorrowingCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BorrowingTransactor is an auto generated write-only Go binding around an Ethereum contract.
type BorrowingTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BorrowingFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type BorrowingFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BorrowingSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type BorrowingSession struct {
	Contract     *Borrowing        // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// BorrowingCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type BorrowingCallerSession struct {
	Contract *BorrowingCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts    // Call options to use throughout this session
}

// BorrowingTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type BorrowingTransactorSession struct {
	Contract     *BorrowingTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// BorrowingRaw is an auto generated low-level Go binding around an Ethereum contract.
type BorrowingRaw struct {
	Contract *Borrowing // Generic contract binding to access the raw methods on
}

// BorrowingCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type BorrowingCallerRaw struct {
	Contract *BorrowingCaller // Generic read-only contract binding to access the raw methods on
}

// BorrowingTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type BorrowingTransactorRaw struct {
	Contract *BorrowingTransactor // Generic write-only contract binding to access the raw methods on
}

// NewBorrowing creates a new instance of Borrowing, bound to a specific deployed contract.
func NewBorrowing(address common.Address, backend bind.ContractBackend) (*Borrowing, error) {
	contract, err := bindBorrowing(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Borrowing{BorrowingCaller: BorrowingCaller{contract: contract}, BorrowingTransactor: BorrowingTransactor{contract: contract}, BorrowingFilterer: BorrowingFilterer{contract: contract}}, nil
}

// NewBorrowingCaller creates a new read-only instance of Borrowing, bound to a specific deployed contract.
func NewBorrowingCaller(address common.Address, caller bind.ContractCaller) (*BorrowingCaller, error) {
	contract, err := bindBorrowing(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &BorrowingCaller{contract: contract}, nil
}

// NewBorrowingTransactor creates a new write-only instance of Borrowing, bound to a specific deployed contract.
func NewBorrowingTransactor(address common.Address, transactor bind.ContractTransactor) (*BorrowingTransactor, error) {
	contract, err := bindBorrowing(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &BorrowingTransactor{contract: contract}, nil
}

// NewBorrowingFilterer creates a new log filterer instance of Borrowing, bound to a specific deployed contract.
func NewBorrowingFilterer(address common.Address, filterer bind.ContractFilterer) (*BorrowingFilterer, error) {
	contract, err := bindBorrowing(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &BorrowingFilterer{contract: contract}, nil
}

// bindBorrowing binds a generic wrapper to an already deployed contract.
func bindBorrowing(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := BorrowingMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Borrowing *BorrowingRaw) Call(opts *bind.CallOpts, result *[]any, method string, params ...any) error {
	return _Borrowing.Contract.BorrowingCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Borrowing *BorrowingRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Borrowing.Contract.BorrowingTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Borrowing *BorrowingRaw) Transact(opts *bind.TransactOpts, method string, params ...any) (*types.Transaction, error) {
	return _Borrowing.Contract.BorrowingTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Borrowing *BorrowingCallerRaw) Call(opts *bind.CallOpts, result *[]any, method string, params ...any) error {
	return _Borrowing.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Borrowing *BorrowingTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Borrowing.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Borrowing *BorrowingTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...any) (*types.Transaction, error) {
	return _Borrowing.Contract.contract.Transact(opts, method, params...)
}

// Beta is a free data retrieval call binding the contract method 0x9faa3c91.
//
// Solidity: function beta() view returns(uint256)
func (_Borrowing *BorrowingCaller) Beta(opts *bind.CallOpts) (*big.Int, error) {
	var out []any
	err := _Borrowing.contract.Call(opts, &out, "beta")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Beta is a free data retrieval call binding the contract method 0x9faa3c91.
//
// Solidity: function beta() view returns(uint256)
func (_Borrowing *BorrowingSession) Beta() (*big.Int, error) {
	return _Borrowing.Contract.Beta(&_Borrowing.CallOpts)
}

// Beta is a free data retrieval call binding the contract method 0x9faa3c91.
//
// Solidity: function beta() view returns(uint256)
func (_Borrowing *BorrowingCallerSession) Beta() (*big.Int, error) {
	return _Borrowing.Contract.Beta(&_Borrowing.CallOpts)
}

// BorrowedPrincipal is a free data retrieval call binding the contract method 0x7c523690.
//
// Solidity: function borrowedPrincipal(address ) view returns(uint256)
func (_Borrowing *BorrowingCaller) BorrowedPrincipal(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []any
	err := _Borrowing.contract.Call(opts, &out, "borrowedPrincipal", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BorrowedPrincipal is a free data retrieval call binding the contract method 0x7c523690.
//
// Solidity: function borrowedPrincipal(address ) view returns(uint256)
func (_Borrowing *BorrowingSession) BorrowedPrincipal(arg0 common.Address) (*big.Int, error) {
	return _Borrowing.Contract.BorrowedPrincipal(&_Borrowing.CallOpts, arg0)
}

// BorrowedPrincipal is a free data retrieval call binding the contract method 0x7c523690.
//
// Solidity: function borrowedPrincipal(address ) view returns(uint256)
func (_Borrowing *BorrowingCallerSession) BorrowedPrincipal(arg0 common.Address) (*big.Int, error) {
	return _Borrowing.Contract.BorrowedPrincipal(&_Borrowing.CallOpts, arg0)
}

// Collateral is a free data retrieval call binding the contract method 0xd8dfeb45.
//
// Solidity: function collateral() view returns(address)
func (_Borrowing *BorrowingCaller) Collateral(opts *bind.CallOpts) (common.Address, error) {
	var out []any
	err := _Borrowing.contract.Call(opts, &out, "collateral")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Collateral is a free data retrieval call binding the contract method 0xd8dfeb45.
//
// Solidity: function collateral() view returns(address)
func (_Borrowing *BorrowingSession) Collateral() (common.Address, error) {
	return _Borrowing.Contract.Collateral(&_Borrowing.CallOpts)
}

// Collateral is a free data retrieval call binding the contract method 0xd8dfeb45.
//
// Solidity: function collateral() view returns(address)
func (_Borrowing *BorrowingCallerSession) Collateral() (common.Address, error) {
	return _Borrowing.Contract.Collateral(&_Borrowing.CallOpts)
}

// GetAllBorrowToken is a free data retrieval call binding the contract method 0x2e021184.
//
// Solidity: function getAllBorrowToken() view returns(uint256)
func (_Borrowing *BorrowingCaller) GetAllBorrowToken(opts *bind.CallOpts) (*big.Int, error) {
	var out []any
	err := _Borrowing.contract.Call(opts, &out, "getAllBorrowToken")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetAllBorrowToken is a free data retrieval call binding the contract method 0x2e021184.
//
// Solidity: function getAllBorrowToken() view returns(uint256)
func (_Borrowing *BorrowingSession) GetAllBorrowToken() (*big.Int, error) {
	return _Borrowing.Contract.GetAllBorrowToken(&_Borrowing.CallOpts)
}

// GetAllBorrowToken is a free data retrieval call binding the contract method 0x2e021184.
//
// Solidity: function getAllBorrowToken() view returns(uint256)
func (_Borrowing *BorrowingCallerSession) GetAllBorrowToken() (*big.Int, error) {
	return _Borrowing.Contract.GetAllBorrowToken(&_Borrowing.CallOpts)
}

// GetBorrowToken is a free data retrieval call binding the contract method 0x3765891f.
//
// Solidity: function getBorrowToken(address user) view returns(uint256)
func (_Borrowing *BorrowingCaller) GetBorrowToken(opts *bind.CallOpts, user common.Address) (*big.Int, error) {
	var out []any
	err := _Borrowing.contract.Call(opts, &out, "getBorrowToken", user)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetBorrowToken is a free data retrieval call binding the contract method 0x3765891f.
//
// Solidity: function getBorrowToken(address user) view returns(uint256)
func (_Borrowing *BorrowingSession) GetBorrowToken(user common.Address) (*big.Int, error) {
	return _Borrowing.Contract.GetBorrowToken(&_Borrowing.CallOpts, user)
}

// GetBorrowToken is a free data retrieval call binding the contract method 0x3765891f.
//
// Solidity: function getBorrowToken(address user) view returns(uint256)
func (_Borrowing *BorrowingCallerSession) GetBorrowToken(user common.Address) (*big.Int, error) {
	return _Borrowing.Contract.GetBorrowToken(&_Borrowing.CallOpts, user)
}

// GetCurrentRate is a free data retrieval call binding the contract method 0xf7fb07b0.
//
// Solidity: function getCurrentRate() view returns(uint256)
func (_Borrowing *BorrowingCaller) GetCurrentRate(opts *bind.CallOpts) (*big.Int, error) {
	var out []any
	err := _Borrowing.contract.Call(opts, &out, "getCurrentRate")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetCurrentRate is a free data retrieval call binding the contract method 0xf7fb07b0.
//
// Solidity: function getCurrentRate() view returns(uint256)
func (_Borrowing *BorrowingSession) GetCurrentRate() (*big.Int, error) {
	return _Borrowing.Contract.GetCurrentRate(&_Borrowing.CallOpts)
}

// GetCurrentRate is a free data retrieval call binding the contract method 0xf7fb07b0.
//
// Solidity: function getCurrentRate() view returns(uint256)
func (_Borrowing *BorrowingCallerSession) GetCurrentRate() (*big.Int, error) {
	return _Borrowing.Contract.GetCurrentRate(&_Borrowing.CallOpts)
}

// LastUpdateTime is a free data retrieval call binding the contract method 0x2ce9aead.
//
// Solidity: function lastUpdateTime(address ) view returns(uint256)
func (_Borrowing *BorrowingCaller) LastUpdateTime(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []any
	err := _Borrowing.contract.Call(opts, &out, "lastUpdateTime", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LastUpdateTime is a free data retrieval call binding the contract method 0x2ce9aead.
//
// Solidity: function lastUpdateTime(address ) view returns(uint256)
func (_Borrowing *BorrowingSession) LastUpdateTime(arg0 common.Address) (*big.Int, error) {
	return _Borrowing.Contract.LastUpdateTime(&_Borrowing.CallOpts, arg0)
}

// LastUpdateTime is a free data retrieval call binding the contract method 0x2ce9aead.
//
// Solidity: function lastUpdateTime(address ) view returns(uint256)
func (_Borrowing *BorrowingCallerSession) LastUpdateTime(arg0 common.Address) (*big.Int, error) {
	return _Borrowing.Contract.LastUpdateTime(&_Borrowing.CallOpts, arg0)
}

// RMax is a free data retrieval call binding the contract method 0x8301f1d6.
//
// Solidity: function rMax() view returns(uint256)
func (_Borrowing *BorrowingCaller) RMax(opts *bind.CallOpts) (*big.Int, error) {
	var out []any
	err := _Borrowing.contract.Call(opts, &out, "rMax")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// RMax is a free data retrieval call binding the contract method 0x8301f1d6.
//
// Solidity: function rMax() view returns(uint256)
func (_Borrowing *BorrowingSession) RMax() (*big.Int, error) {
	return _Borrowing.Contract.RMax(&_Borrowing.CallOpts)
}

// RMax is a free data retrieval call binding the contract method 0x8301f1d6.
//
// Solidity: function rMax() view returns(uint256)
func (_Borrowing *BorrowingCallerSession) RMax() (*big.Int, error) {
	return _Borrowing.Contract.RMax(&_Borrowing.CallOpts)
}

// RMin is a free data retrieval call binding the contract method 0xd28a1259.
//
// Solidity: function rMin() view returns(uint256)
func (_Borrowing *BorrowingCaller) RMin(opts *bind.CallOpts) (*big.Int, error) {
	var out []any
	err := _Borrowing.contract.Call(opts, &out, "rMin")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// RMin is a free data retrieval call binding the contract method 0xd28a1259.
//
// Solidity: function rMin() view returns(uint256)
func (_Borrowing *BorrowingSession) RMin() (*big.Int, error) {
	return _Borrowing.Contract.RMin(&_Borrowing.CallOpts)
}

// RMin is a free data retrieval call binding the contract method 0xd28a1259.
//
// Solidity: function rMin() view returns(uint256)
func (_Borrowing *BorrowingCallerSession) RMin() (*big.Int, error) {
	return _Borrowing.Contract.RMin(&_Borrowing.CallOpts)
}

// Token is a free data retrieval call binding the contract method 0xfc0c546a.
//
// Solidity: function token() view returns(address)
func (_Borrowing *BorrowingCaller) Token(opts *bind.CallOpts) (common.Address, error) {
	var out []any
	err := _Borrowing.contract.Call(opts, &out, "token")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Token is a free data retrieval call binding the contract method 0xfc0c546a.
//
// Solidity: function token() view returns(address)
func (_Borrowing *BorrowingSession) Token() (common.Address, error) {
	return _Borrowing.Contract.Token(&_Borrowing.CallOpts)
}

// Token is a free data retrieval call binding the contract method 0xfc0c546a.
//
// Solidity: function token() view returns(address)
func (_Borrowing *BorrowingCallerSession) Token() (common.Address, error) {
	return _Borrowing.Contract.Token(&_Borrowing.CallOpts)
}

// TotalBorrowed is a free data retrieval call binding the contract method 0x4c19386c.
//
// Solidity: function totalBorrowed() view returns(uint256)
func (_Borrowing *BorrowingCaller) TotalBorrowed(opts *bind.CallOpts) (*big.Int, error) {
	var out []any
	err := _Borrowing.contract.Call(opts, &out, "totalBorrowed")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalBorrowed is a free data retrieval call binding the contract method 0x4c19386c.
//
// Solidity: function totalBorrowed() view returns(uint256)
func (_Borrowing *BorrowingSession) TotalBorrowed() (*big.Int, error) {
	return _Borrowing.Contract.TotalBorrowed(&_Borrowing.CallOpts)
}

// TotalBorrowed is a free data retrieval call binding the contract method 0x4c19386c.
//
// Solidity: function totalBorrowed() view returns(uint256)
func (_Borrowing *BorrowingCallerSession) TotalBorrowed() (*big.Int, error) {
	return _Borrowing.Contract.TotalBorrowed(&_Borrowing.CallOpts)
}

// Borrow is a paid mutator transaction binding the contract method 0xc5ebeaec.
//
// Solidity: function borrow(uint256 amount) returns()
func (_Borrowing *BorrowingTransactor) Borrow(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _Borrowing.contract.Transact(opts, "borrow", amount)
}

// Borrow is a paid mutator transaction binding the contract method 0xc5ebeaec.
//
// Solidity: function borrow(uint256 amount) returns()
func (_Borrowing *BorrowingSession) Borrow(amount *big.Int) (*types.Transaction, error) {
	return _Borrowing.Contract.Borrow(&_Borrowing.TransactOpts, amount)
}

// Borrow is a paid mutator transaction binding the contract method 0xc5ebeaec.
//
// Solidity: function borrow(uint256 amount) returns()
func (_Borrowing *BorrowingTransactorSession) Borrow(amount *big.Int) (*types.Transaction, error) {
	return _Borrowing.Contract.Borrow(&_Borrowing.TransactOpts, amount)
}

// ReduceDebt is a paid mutator transaction binding the contract method 0x807d3f1d.
//
// Solidity: function reduceDebt(address borrower, uint256 amount) returns()
func (_Borrowing *BorrowingTransactor) ReduceDebt(opts *bind.TransactOpts, borrower common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Borrowing.contract.Transact(opts, "reduceDebt", borrower, amount)
}

// ReduceDebt is a paid mutator transaction binding the contract method 0x807d3f1d.
//
// Solidity: function reduceDebt(address borrower, uint256 amount) returns()
func (_Borrowing *BorrowingSession) ReduceDebt(borrower common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Borrowing.Contract.ReduceDebt(&_Borrowing.TransactOpts, borrower, amount)
}

// ReduceDebt is a paid mutator transaction binding the contract method 0x807d3f1d.
//
// Solidity: function reduceDebt(address borrower, uint256 amount) returns()
func (_Borrowing *BorrowingTransactorSession) ReduceDebt(borrower common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Borrowing.Contract.ReduceDebt(&_Borrowing.TransactOpts, borrower, amount)
}

// Repay is a paid mutator transaction binding the contract method 0x371fd8e6.
//
// Solidity: function repay(uint256 amount) returns()
func (_Borrowing *BorrowingTransactor) Repay(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _Borrowing.contract.Transact(opts, "repay", amount)
}

// Repay is a paid mutator transaction binding the contract method 0x371fd8e6.
//
// Solidity: function repay(uint256 amount) returns()
func (_Borrowing *BorrowingSession) Repay(amount *big.Int) (*types.Transaction, error) {
	return _Borrowing.Contract.Repay(&_Borrowing.TransactOpts, amount)
}

// Repay is a paid mutator transaction binding the contract method 0x371fd8e6.
//
// Solidity: function repay(uint256 amount) returns()
func (_Borrowing *BorrowingTransactorSession) Repay(amount *big.Int) (*types.Transaction, error) {
	return _Borrowing.Contract.Repay(&_Borrowing.TransactOpts, amount)
}
