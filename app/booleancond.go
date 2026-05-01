// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package app

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

// IBooleanCondMetaData contains all meta data concerning the IBooleanCond contract.
var IBooleanCondMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"getOutcome\",\"inputs\":[{\"name\":\"_query\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isFinalized\",\"inputs\":[{\"name\":\"_query\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"}]",
}

// IBooleanCondABI is the input ABI used to generate the binding from.
// Deprecated: Use IBooleanCondMetaData.ABI instead.
var IBooleanCondABI = IBooleanCondMetaData.ABI

// IBooleanCond is an auto generated Go binding around an Ethereum contract.
type IBooleanCond struct {
	IBooleanCondCaller     // Read-only binding to the contract
	IBooleanCondTransactor // Write-only binding to the contract
	IBooleanCondFilterer   // Log filterer for contract events
}

// IBooleanCondCaller is an auto generated read-only Go binding around an Ethereum contract.
type IBooleanCondCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IBooleanCondTransactor is an auto generated write-only Go binding around an Ethereum contract.
type IBooleanCondTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IBooleanCondFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type IBooleanCondFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IBooleanCondSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type IBooleanCondSession struct {
	Contract     *IBooleanCond     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// IBooleanCondCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type IBooleanCondCallerSession struct {
	Contract *IBooleanCondCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// IBooleanCondTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type IBooleanCondTransactorSession struct {
	Contract     *IBooleanCondTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// IBooleanCondRaw is an auto generated low-level Go binding around an Ethereum contract.
type IBooleanCondRaw struct {
	Contract *IBooleanCond // Generic contract binding to access the raw methods on
}

// IBooleanCondCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type IBooleanCondCallerRaw struct {
	Contract *IBooleanCondCaller // Generic read-only contract binding to access the raw methods on
}

// IBooleanCondTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type IBooleanCondTransactorRaw struct {
	Contract *IBooleanCondTransactor // Generic write-only contract binding to access the raw methods on
}

// NewIBooleanCond creates a new instance of IBooleanCond, bound to a specific deployed contract.
func NewIBooleanCond(address common.Address, backend bind.ContractBackend) (*IBooleanCond, error) {
	contract, err := bindIBooleanCond(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IBooleanCond{IBooleanCondCaller: IBooleanCondCaller{contract: contract}, IBooleanCondTransactor: IBooleanCondTransactor{contract: contract}, IBooleanCondFilterer: IBooleanCondFilterer{contract: contract}}, nil
}

// NewIBooleanCondCaller creates a new read-only instance of IBooleanCond, bound to a specific deployed contract.
func NewIBooleanCondCaller(address common.Address, caller bind.ContractCaller) (*IBooleanCondCaller, error) {
	contract, err := bindIBooleanCond(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IBooleanCondCaller{contract: contract}, nil
}

// NewIBooleanCondTransactor creates a new write-only instance of IBooleanCond, bound to a specific deployed contract.
func NewIBooleanCondTransactor(address common.Address, transactor bind.ContractTransactor) (*IBooleanCondTransactor, error) {
	contract, err := bindIBooleanCond(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IBooleanCondTransactor{contract: contract}, nil
}

// NewIBooleanCondFilterer creates a new log filterer instance of IBooleanCond, bound to a specific deployed contract.
func NewIBooleanCondFilterer(address common.Address, filterer bind.ContractFilterer) (*IBooleanCondFilterer, error) {
	contract, err := bindIBooleanCond(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IBooleanCondFilterer{contract: contract}, nil
}

// bindIBooleanCond binds a generic wrapper to an already deployed contract.
func bindIBooleanCond(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := IBooleanCondMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IBooleanCond *IBooleanCondRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IBooleanCond.Contract.IBooleanCondCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IBooleanCond *IBooleanCondRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IBooleanCond.Contract.IBooleanCondTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IBooleanCond *IBooleanCondRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IBooleanCond.Contract.IBooleanCondTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IBooleanCond *IBooleanCondCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IBooleanCond.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IBooleanCond *IBooleanCondTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IBooleanCond.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IBooleanCond *IBooleanCondTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IBooleanCond.Contract.contract.Transact(opts, method, params...)
}

// GetOutcome is a free data retrieval call binding the contract method 0xea4ba8eb.
//
// Solidity: function getOutcome(bytes _query) view returns(bool)
func (_IBooleanCond *IBooleanCondCaller) GetOutcome(opts *bind.CallOpts, _query []byte) (bool, error) {
	var out []interface{}
	err := _IBooleanCond.contract.Call(opts, &out, "getOutcome", _query)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// GetOutcome is a free data retrieval call binding the contract method 0xea4ba8eb.
//
// Solidity: function getOutcome(bytes _query) view returns(bool)
func (_IBooleanCond *IBooleanCondSession) GetOutcome(_query []byte) (bool, error) {
	return _IBooleanCond.Contract.GetOutcome(&_IBooleanCond.CallOpts, _query)
}

// GetOutcome is a free data retrieval call binding the contract method 0xea4ba8eb.
//
// Solidity: function getOutcome(bytes _query) view returns(bool)
func (_IBooleanCond *IBooleanCondCallerSession) GetOutcome(_query []byte) (bool, error) {
	return _IBooleanCond.Contract.GetOutcome(&_IBooleanCond.CallOpts, _query)
}

// IsFinalized is a free data retrieval call binding the contract method 0xbcdbda94.
//
// Solidity: function isFinalized(bytes _query) view returns(bool)
func (_IBooleanCond *IBooleanCondCaller) IsFinalized(opts *bind.CallOpts, _query []byte) (bool, error) {
	var out []interface{}
	err := _IBooleanCond.contract.Call(opts, &out, "isFinalized", _query)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsFinalized is a free data retrieval call binding the contract method 0xbcdbda94.
//
// Solidity: function isFinalized(bytes _query) view returns(bool)
func (_IBooleanCond *IBooleanCondSession) IsFinalized(_query []byte) (bool, error) {
	return _IBooleanCond.Contract.IsFinalized(&_IBooleanCond.CallOpts, _query)
}

// IsFinalized is a free data retrieval call binding the contract method 0xbcdbda94.
//
// Solidity: function isFinalized(bytes _query) view returns(bool)
func (_IBooleanCond *IBooleanCondCallerSession) IsFinalized(_query []byte) (bool, error) {
	return _IBooleanCond.Contract.IsFinalized(&_IBooleanCond.CallOpts, _query)
}
