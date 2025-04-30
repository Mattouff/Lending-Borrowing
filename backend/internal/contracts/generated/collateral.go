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

// CollateralMetaData contains all meta data concerning the Collateral contract.
var CollateralMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"_token\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_collateral\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"LIQUIDATION_BONUS\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"LIQUIDATION_THRESHOLD\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"MAX_BORROWING_PERCENTAGE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"MIN_COLLATERAL_RATIO\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"borrowing\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractBorrowing\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"canBorrow\",\"inputs\":[{\"name\":\"user\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"borrowAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"collateralBalance\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"depositCollateral\",\"inputs\":[{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getCollateralRatio\",\"inputs\":[{\"name\":\"user\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getMaxBorrowableAmount\",\"inputs\":[{\"name\":\"user\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"liquidate\",\"inputs\":[{\"name\":\"borrower\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"repayAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"token\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractToken\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"withdrawCollateral\",\"inputs\":[{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"}]",
}

// CollateralABI is the input ABI used to generate the binding from.
// Deprecated: Use CollateralMetaData.ABI instead.
var CollateralABI = CollateralMetaData.ABI

// Collateral is an auto generated Go binding around an Ethereum contract.
type Collateral struct {
	CollateralCaller     // Read-only binding to the contract
	CollateralTransactor // Write-only binding to the contract
	CollateralFilterer   // Log filterer for contract events
}

// CollateralCaller is an auto generated read-only Go binding around an Ethereum contract.
type CollateralCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CollateralTransactor is an auto generated write-only Go binding around an Ethereum contract.
type CollateralTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CollateralFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type CollateralFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CollateralSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type CollateralSession struct {
	Contract     *Collateral       // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// CollateralCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type CollateralCallerSession struct {
	Contract *CollateralCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts     // Call options to use throughout this session
}

// CollateralTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type CollateralTransactorSession struct {
	Contract     *CollateralTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// CollateralRaw is an auto generated low-level Go binding around an Ethereum contract.
type CollateralRaw struct {
	Contract *Collateral // Generic contract binding to access the raw methods on
}

// CollateralCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type CollateralCallerRaw struct {
	Contract *CollateralCaller // Generic read-only contract binding to access the raw methods on
}

// CollateralTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type CollateralTransactorRaw struct {
	Contract *CollateralTransactor // Generic write-only contract binding to access the raw methods on
}

// NewCollateral creates a new instance of Collateral, bound to a specific deployed contract.
func NewCollateral(address common.Address, backend bind.ContractBackend) (*Collateral, error) {
	contract, err := bindCollateral(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Collateral{CollateralCaller: CollateralCaller{contract: contract}, CollateralTransactor: CollateralTransactor{contract: contract}, CollateralFilterer: CollateralFilterer{contract: contract}}, nil
}

// NewCollateralCaller creates a new read-only instance of Collateral, bound to a specific deployed contract.
func NewCollateralCaller(address common.Address, caller bind.ContractCaller) (*CollateralCaller, error) {
	contract, err := bindCollateral(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &CollateralCaller{contract: contract}, nil
}

// NewCollateralTransactor creates a new write-only instance of Collateral, bound to a specific deployed contract.
func NewCollateralTransactor(address common.Address, transactor bind.ContractTransactor) (*CollateralTransactor, error) {
	contract, err := bindCollateral(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &CollateralTransactor{contract: contract}, nil
}

// NewCollateralFilterer creates a new log filterer instance of Collateral, bound to a specific deployed contract.
func NewCollateralFilterer(address common.Address, filterer bind.ContractFilterer) (*CollateralFilterer, error) {
	contract, err := bindCollateral(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &CollateralFilterer{contract: contract}, nil
}

// bindCollateral binds a generic wrapper to an already deployed contract.
func bindCollateral(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := CollateralMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Collateral *CollateralRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Collateral.Contract.CollateralCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Collateral *CollateralRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Collateral.Contract.CollateralTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Collateral *CollateralRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Collateral.Contract.CollateralTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Collateral *CollateralCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Collateral.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Collateral *CollateralTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Collateral.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Collateral *CollateralTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Collateral.Contract.contract.Transact(opts, method, params...)
}

// LIQUIDATIONBONUS is a free data retrieval call binding the contract method 0x3574d4c4.
//
// Solidity: function LIQUIDATION_BONUS() view returns(uint256)
func (_Collateral *CollateralCaller) LIQUIDATIONBONUS(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Collateral.contract.Call(opts, &out, "LIQUIDATION_BONUS")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LIQUIDATIONBONUS is a free data retrieval call binding the contract method 0x3574d4c4.
//
// Solidity: function LIQUIDATION_BONUS() view returns(uint256)
func (_Collateral *CollateralSession) LIQUIDATIONBONUS() (*big.Int, error) {
	return _Collateral.Contract.LIQUIDATIONBONUS(&_Collateral.CallOpts)
}

// LIQUIDATIONBONUS is a free data retrieval call binding the contract method 0x3574d4c4.
//
// Solidity: function LIQUIDATION_BONUS() view returns(uint256)
func (_Collateral *CollateralCallerSession) LIQUIDATIONBONUS() (*big.Int, error) {
	return _Collateral.Contract.LIQUIDATIONBONUS(&_Collateral.CallOpts)
}

// LIQUIDATIONTHRESHOLD is a free data retrieval call binding the contract method 0x90a8ae9b.
//
// Solidity: function LIQUIDATION_THRESHOLD() view returns(uint256)
func (_Collateral *CollateralCaller) LIQUIDATIONTHRESHOLD(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Collateral.contract.Call(opts, &out, "LIQUIDATION_THRESHOLD")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LIQUIDATIONTHRESHOLD is a free data retrieval call binding the contract method 0x90a8ae9b.
//
// Solidity: function LIQUIDATION_THRESHOLD() view returns(uint256)
func (_Collateral *CollateralSession) LIQUIDATIONTHRESHOLD() (*big.Int, error) {
	return _Collateral.Contract.LIQUIDATIONTHRESHOLD(&_Collateral.CallOpts)
}

// LIQUIDATIONTHRESHOLD is a free data retrieval call binding the contract method 0x90a8ae9b.
//
// Solidity: function LIQUIDATION_THRESHOLD() view returns(uint256)
func (_Collateral *CollateralCallerSession) LIQUIDATIONTHRESHOLD() (*big.Int, error) {
	return _Collateral.Contract.LIQUIDATIONTHRESHOLD(&_Collateral.CallOpts)
}

// MAXBORROWINGPERCENTAGE is a free data retrieval call binding the contract method 0x07b037ba.
//
// Solidity: function MAX_BORROWING_PERCENTAGE() view returns(uint256)
func (_Collateral *CollateralCaller) MAXBORROWINGPERCENTAGE(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Collateral.contract.Call(opts, &out, "MAX_BORROWING_PERCENTAGE")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MAXBORROWINGPERCENTAGE is a free data retrieval call binding the contract method 0x07b037ba.
//
// Solidity: function MAX_BORROWING_PERCENTAGE() view returns(uint256)
func (_Collateral *CollateralSession) MAXBORROWINGPERCENTAGE() (*big.Int, error) {
	return _Collateral.Contract.MAXBORROWINGPERCENTAGE(&_Collateral.CallOpts)
}

// MAXBORROWINGPERCENTAGE is a free data retrieval call binding the contract method 0x07b037ba.
//
// Solidity: function MAX_BORROWING_PERCENTAGE() view returns(uint256)
func (_Collateral *CollateralCallerSession) MAXBORROWINGPERCENTAGE() (*big.Int, error) {
	return _Collateral.Contract.MAXBORROWINGPERCENTAGE(&_Collateral.CallOpts)
}

// MINCOLLATERALRATIO is a free data retrieval call binding the contract method 0x7a9fffb7.
//
// Solidity: function MIN_COLLATERAL_RATIO() view returns(uint256)
func (_Collateral *CollateralCaller) MINCOLLATERALRATIO(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Collateral.contract.Call(opts, &out, "MIN_COLLATERAL_RATIO")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MINCOLLATERALRATIO is a free data retrieval call binding the contract method 0x7a9fffb7.
//
// Solidity: function MIN_COLLATERAL_RATIO() view returns(uint256)
func (_Collateral *CollateralSession) MINCOLLATERALRATIO() (*big.Int, error) {
	return _Collateral.Contract.MINCOLLATERALRATIO(&_Collateral.CallOpts)
}

// MINCOLLATERALRATIO is a free data retrieval call binding the contract method 0x7a9fffb7.
//
// Solidity: function MIN_COLLATERAL_RATIO() view returns(uint256)
func (_Collateral *CollateralCallerSession) MINCOLLATERALRATIO() (*big.Int, error) {
	return _Collateral.Contract.MINCOLLATERALRATIO(&_Collateral.CallOpts)
}

// Borrowing is a free data retrieval call binding the contract method 0x60b415ed.
//
// Solidity: function borrowing() view returns(address)
func (_Collateral *CollateralCaller) Borrowing(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Collateral.contract.Call(opts, &out, "borrowing")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Borrowing is a free data retrieval call binding the contract method 0x60b415ed.
//
// Solidity: function borrowing() view returns(address)
func (_Collateral *CollateralSession) Borrowing() (common.Address, error) {
	return _Collateral.Contract.Borrowing(&_Collateral.CallOpts)
}

// Borrowing is a free data retrieval call binding the contract method 0x60b415ed.
//
// Solidity: function borrowing() view returns(address)
func (_Collateral *CollateralCallerSession) Borrowing() (common.Address, error) {
	return _Collateral.Contract.Borrowing(&_Collateral.CallOpts)
}

// CanBorrow is a free data retrieval call binding the contract method 0xcf8c8c44.
//
// Solidity: function canBorrow(address user, uint256 borrowAmount) view returns(bool)
func (_Collateral *CollateralCaller) CanBorrow(opts *bind.CallOpts, user common.Address, borrowAmount *big.Int) (bool, error) {
	var out []interface{}
	err := _Collateral.contract.Call(opts, &out, "canBorrow", user, borrowAmount)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// CanBorrow is a free data retrieval call binding the contract method 0xcf8c8c44.
//
// Solidity: function canBorrow(address user, uint256 borrowAmount) view returns(bool)
func (_Collateral *CollateralSession) CanBorrow(user common.Address, borrowAmount *big.Int) (bool, error) {
	return _Collateral.Contract.CanBorrow(&_Collateral.CallOpts, user, borrowAmount)
}

// CanBorrow is a free data retrieval call binding the contract method 0xcf8c8c44.
//
// Solidity: function canBorrow(address user, uint256 borrowAmount) view returns(bool)
func (_Collateral *CollateralCallerSession) CanBorrow(user common.Address, borrowAmount *big.Int) (bool, error) {
	return _Collateral.Contract.CanBorrow(&_Collateral.CallOpts, user, borrowAmount)
}

// CollateralBalance is a free data retrieval call binding the contract method 0xa1bf2840.
//
// Solidity: function collateralBalance(address ) view returns(uint256)
func (_Collateral *CollateralCaller) CollateralBalance(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Collateral.contract.Call(opts, &out, "collateralBalance", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// CollateralBalance is a free data retrieval call binding the contract method 0xa1bf2840.
//
// Solidity: function collateralBalance(address ) view returns(uint256)
func (_Collateral *CollateralSession) CollateralBalance(arg0 common.Address) (*big.Int, error) {
	return _Collateral.Contract.CollateralBalance(&_Collateral.CallOpts, arg0)
}

// CollateralBalance is a free data retrieval call binding the contract method 0xa1bf2840.
//
// Solidity: function collateralBalance(address ) view returns(uint256)
func (_Collateral *CollateralCallerSession) CollateralBalance(arg0 common.Address) (*big.Int, error) {
	return _Collateral.Contract.CollateralBalance(&_Collateral.CallOpts, arg0)
}

// GetCollateralRatio is a free data retrieval call binding the contract method 0x15a3ba43.
//
// Solidity: function getCollateralRatio(address user) view returns(uint256)
func (_Collateral *CollateralCaller) GetCollateralRatio(opts *bind.CallOpts, user common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Collateral.contract.Call(opts, &out, "getCollateralRatio", user)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetCollateralRatio is a free data retrieval call binding the contract method 0x15a3ba43.
//
// Solidity: function getCollateralRatio(address user) view returns(uint256)
func (_Collateral *CollateralSession) GetCollateralRatio(user common.Address) (*big.Int, error) {
	return _Collateral.Contract.GetCollateralRatio(&_Collateral.CallOpts, user)
}

// GetCollateralRatio is a free data retrieval call binding the contract method 0x15a3ba43.
//
// Solidity: function getCollateralRatio(address user) view returns(uint256)
func (_Collateral *CollateralCallerSession) GetCollateralRatio(user common.Address) (*big.Int, error) {
	return _Collateral.Contract.GetCollateralRatio(&_Collateral.CallOpts, user)
}

// GetMaxBorrowableAmount is a free data retrieval call binding the contract method 0x2bda1f39.
//
// Solidity: function getMaxBorrowableAmount(address user) view returns(uint256)
func (_Collateral *CollateralCaller) GetMaxBorrowableAmount(opts *bind.CallOpts, user common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Collateral.contract.Call(opts, &out, "getMaxBorrowableAmount", user)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetMaxBorrowableAmount is a free data retrieval call binding the contract method 0x2bda1f39.
//
// Solidity: function getMaxBorrowableAmount(address user) view returns(uint256)
func (_Collateral *CollateralSession) GetMaxBorrowableAmount(user common.Address) (*big.Int, error) {
	return _Collateral.Contract.GetMaxBorrowableAmount(&_Collateral.CallOpts, user)
}

// GetMaxBorrowableAmount is a free data retrieval call binding the contract method 0x2bda1f39.
//
// Solidity: function getMaxBorrowableAmount(address user) view returns(uint256)
func (_Collateral *CollateralCallerSession) GetMaxBorrowableAmount(user common.Address) (*big.Int, error) {
	return _Collateral.Contract.GetMaxBorrowableAmount(&_Collateral.CallOpts, user)
}

// Token is a free data retrieval call binding the contract method 0xfc0c546a.
//
// Solidity: function token() view returns(address)
func (_Collateral *CollateralCaller) Token(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Collateral.contract.Call(opts, &out, "token")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Token is a free data retrieval call binding the contract method 0xfc0c546a.
//
// Solidity: function token() view returns(address)
func (_Collateral *CollateralSession) Token() (common.Address, error) {
	return _Collateral.Contract.Token(&_Collateral.CallOpts)
}

// Token is a free data retrieval call binding the contract method 0xfc0c546a.
//
// Solidity: function token() view returns(address)
func (_Collateral *CollateralCallerSession) Token() (common.Address, error) {
	return _Collateral.Contract.Token(&_Collateral.CallOpts)
}

// DepositCollateral is a paid mutator transaction binding the contract method 0xbad4a01f.
//
// Solidity: function depositCollateral(uint256 amount) returns()
func (_Collateral *CollateralTransactor) DepositCollateral(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _Collateral.contract.Transact(opts, "depositCollateral", amount)
}

// DepositCollateral is a paid mutator transaction binding the contract method 0xbad4a01f.
//
// Solidity: function depositCollateral(uint256 amount) returns()
func (_Collateral *CollateralSession) DepositCollateral(amount *big.Int) (*types.Transaction, error) {
	return _Collateral.Contract.DepositCollateral(&_Collateral.TransactOpts, amount)
}

// DepositCollateral is a paid mutator transaction binding the contract method 0xbad4a01f.
//
// Solidity: function depositCollateral(uint256 amount) returns()
func (_Collateral *CollateralTransactorSession) DepositCollateral(amount *big.Int) (*types.Transaction, error) {
	return _Collateral.Contract.DepositCollateral(&_Collateral.TransactOpts, amount)
}

// Liquidate is a paid mutator transaction binding the contract method 0xbcbaf487.
//
// Solidity: function liquidate(address borrower, uint256 repayAmount) returns()
func (_Collateral *CollateralTransactor) Liquidate(opts *bind.TransactOpts, borrower common.Address, repayAmount *big.Int) (*types.Transaction, error) {
	return _Collateral.contract.Transact(opts, "liquidate", borrower, repayAmount)
}

// Liquidate is a paid mutator transaction binding the contract method 0xbcbaf487.
//
// Solidity: function liquidate(address borrower, uint256 repayAmount) returns()
func (_Collateral *CollateralSession) Liquidate(borrower common.Address, repayAmount *big.Int) (*types.Transaction, error) {
	return _Collateral.Contract.Liquidate(&_Collateral.TransactOpts, borrower, repayAmount)
}

// Liquidate is a paid mutator transaction binding the contract method 0xbcbaf487.
//
// Solidity: function liquidate(address borrower, uint256 repayAmount) returns()
func (_Collateral *CollateralTransactorSession) Liquidate(borrower common.Address, repayAmount *big.Int) (*types.Transaction, error) {
	return _Collateral.Contract.Liquidate(&_Collateral.TransactOpts, borrower, repayAmount)
}

// WithdrawCollateral is a paid mutator transaction binding the contract method 0x6112fe2e.
//
// Solidity: function withdrawCollateral(uint256 amount) returns()
func (_Collateral *CollateralTransactor) WithdrawCollateral(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _Collateral.contract.Transact(opts, "withdrawCollateral", amount)
}

// WithdrawCollateral is a paid mutator transaction binding the contract method 0x6112fe2e.
//
// Solidity: function withdrawCollateral(uint256 amount) returns()
func (_Collateral *CollateralSession) WithdrawCollateral(amount *big.Int) (*types.Transaction, error) {
	return _Collateral.Contract.WithdrawCollateral(&_Collateral.TransactOpts, amount)
}

// WithdrawCollateral is a paid mutator transaction binding the contract method 0x6112fe2e.
//
// Solidity: function withdrawCollateral(uint256 amount) returns()
func (_Collateral *CollateralTransactorSession) WithdrawCollateral(amount *big.Int) (*types.Transaction, error) {
	return _Collateral.Contract.WithdrawCollateral(&_Collateral.TransactOpts, amount)
}
