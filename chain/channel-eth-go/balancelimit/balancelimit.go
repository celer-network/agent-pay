// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package balancelimit

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

// LedgerBalanceLimitMetaData contains all meta data concerning the LedgerBalanceLimit contract.
var LedgerBalanceLimitMetaData = &bind.MetaData{
	ABI: "[]",
	Bin: "0x610381610034600b8282823980515f1a607314602857634e487b7160e01b5f525f60045260245ffd5b305f52607381538281f3fe7300000000000000000000000000000000000000003014608060405260043610610060575f3560e01c80635930e0e1146100645780636ad1dc2d1461008d5780636ae97472146100b7578063bdca79a7146100e4578063c88c62651461011e575b5f5ffd5b81801561006f575f5ffd5b5061008b61007e3660046101fa565b600501805460ff19169055565b005b818015610098575f5ffd5b5061008b6100a73660046101fa565b600501805460ff19166001179055565b6100cf6100c53660046101fa565b6005015460ff1690565b60405190151581526020015b60405180910390f35b6101106100f236600461022c565b6001600160a01b03165f908152600491909101602052604090205490565b6040519081526020016100db565b818015610129575f5ffd5b5061008b61013836600461029e565b8281146101825760405162461bcd60e51b8152602060048201526014602482015273098cadccee8d0e640c8de40dcdee840dac2e8c6d60631b604482015260640160405180910390fd5b5f5b838110156101f25782828281811061019e5761019e610317565b90506020020135866004015f8787858181106101bc576101bc610317565b90506020020160208101906101d1919061032b565b6001600160a01b0316815260208101919091526040015f2055600101610184565b505050505050565b5f6020828403121561020a575f5ffd5b5035919050565b80356001600160a01b0381168114610227575f5ffd5b919050565b5f5f6040838503121561023d575f5ffd5b8235915061024d60208401610211565b90509250929050565b5f5f83601f840112610266575f5ffd5b50813567ffffffffffffffff81111561027d575f5ffd5b6020830191508360208260051b8501011115610297575f5ffd5b9250929050565b5f5f5f5f5f606086880312156102b2575f5ffd5b85359450602086013567ffffffffffffffff8111156102cf575f5ffd5b6102db88828901610256565b909550935050604086013567ffffffffffffffff8111156102fa575f5ffd5b61030688828901610256565b969995985093965092949392505050565b634e487b7160e01b5f52603260045260245ffd5b5f6020828403121561033b575f5ffd5b61034482610211565b939250505056fea26469706673582212200a2f75d9c1fdc82c6c73ad9859c159a691858fd85e661d5ba5d08f6134e037a164736f6c634300081e0033",
}

// LedgerBalanceLimitABI is the input ABI used to generate the binding from.
// Deprecated: Use LedgerBalanceLimitMetaData.ABI instead.
var LedgerBalanceLimitABI = LedgerBalanceLimitMetaData.ABI

// LedgerBalanceLimitBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use LedgerBalanceLimitMetaData.Bin instead.
var LedgerBalanceLimitBin = LedgerBalanceLimitMetaData.Bin

// DeployLedgerBalanceLimit deploys a new Ethereum contract, binding an instance of LedgerBalanceLimit to it.
func DeployLedgerBalanceLimit(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *LedgerBalanceLimit, error) {
	parsed, err := LedgerBalanceLimitMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(LedgerBalanceLimitBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &LedgerBalanceLimit{LedgerBalanceLimitCaller: LedgerBalanceLimitCaller{contract: contract}, LedgerBalanceLimitTransactor: LedgerBalanceLimitTransactor{contract: contract}, LedgerBalanceLimitFilterer: LedgerBalanceLimitFilterer{contract: contract}}, nil
}

// LedgerBalanceLimit is an auto generated Go binding around an Ethereum contract.
type LedgerBalanceLimit struct {
	LedgerBalanceLimitCaller     // Read-only binding to the contract
	LedgerBalanceLimitTransactor // Write-only binding to the contract
	LedgerBalanceLimitFilterer   // Log filterer for contract events
}

// LedgerBalanceLimitCaller is an auto generated read-only Go binding around an Ethereum contract.
type LedgerBalanceLimitCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LedgerBalanceLimitTransactor is an auto generated write-only Go binding around an Ethereum contract.
type LedgerBalanceLimitTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LedgerBalanceLimitFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type LedgerBalanceLimitFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LedgerBalanceLimitSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type LedgerBalanceLimitSession struct {
	Contract     *LedgerBalanceLimit // Generic contract binding to set the session for
	CallOpts     bind.CallOpts       // Call options to use throughout this session
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// LedgerBalanceLimitCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type LedgerBalanceLimitCallerSession struct {
	Contract *LedgerBalanceLimitCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts             // Call options to use throughout this session
}

// LedgerBalanceLimitTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type LedgerBalanceLimitTransactorSession struct {
	Contract     *LedgerBalanceLimitTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts             // Transaction auth options to use throughout this session
}

// LedgerBalanceLimitRaw is an auto generated low-level Go binding around an Ethereum contract.
type LedgerBalanceLimitRaw struct {
	Contract *LedgerBalanceLimit // Generic contract binding to access the raw methods on
}

// LedgerBalanceLimitCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type LedgerBalanceLimitCallerRaw struct {
	Contract *LedgerBalanceLimitCaller // Generic read-only contract binding to access the raw methods on
}

// LedgerBalanceLimitTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type LedgerBalanceLimitTransactorRaw struct {
	Contract *LedgerBalanceLimitTransactor // Generic write-only contract binding to access the raw methods on
}

// NewLedgerBalanceLimit creates a new instance of LedgerBalanceLimit, bound to a specific deployed contract.
func NewLedgerBalanceLimit(address common.Address, backend bind.ContractBackend) (*LedgerBalanceLimit, error) {
	contract, err := bindLedgerBalanceLimit(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &LedgerBalanceLimit{LedgerBalanceLimitCaller: LedgerBalanceLimitCaller{contract: contract}, LedgerBalanceLimitTransactor: LedgerBalanceLimitTransactor{contract: contract}, LedgerBalanceLimitFilterer: LedgerBalanceLimitFilterer{contract: contract}}, nil
}

// NewLedgerBalanceLimitCaller creates a new read-only instance of LedgerBalanceLimit, bound to a specific deployed contract.
func NewLedgerBalanceLimitCaller(address common.Address, caller bind.ContractCaller) (*LedgerBalanceLimitCaller, error) {
	contract, err := bindLedgerBalanceLimit(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &LedgerBalanceLimitCaller{contract: contract}, nil
}

// NewLedgerBalanceLimitTransactor creates a new write-only instance of LedgerBalanceLimit, bound to a specific deployed contract.
func NewLedgerBalanceLimitTransactor(address common.Address, transactor bind.ContractTransactor) (*LedgerBalanceLimitTransactor, error) {
	contract, err := bindLedgerBalanceLimit(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &LedgerBalanceLimitTransactor{contract: contract}, nil
}

// NewLedgerBalanceLimitFilterer creates a new log filterer instance of LedgerBalanceLimit, bound to a specific deployed contract.
func NewLedgerBalanceLimitFilterer(address common.Address, filterer bind.ContractFilterer) (*LedgerBalanceLimitFilterer, error) {
	contract, err := bindLedgerBalanceLimit(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &LedgerBalanceLimitFilterer{contract: contract}, nil
}

// bindLedgerBalanceLimit binds a generic wrapper to an already deployed contract.
func bindLedgerBalanceLimit(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := LedgerBalanceLimitMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_LedgerBalanceLimit *LedgerBalanceLimitRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _LedgerBalanceLimit.Contract.LedgerBalanceLimitCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_LedgerBalanceLimit *LedgerBalanceLimitRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LedgerBalanceLimit.Contract.LedgerBalanceLimitTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_LedgerBalanceLimit *LedgerBalanceLimitRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _LedgerBalanceLimit.Contract.LedgerBalanceLimitTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_LedgerBalanceLimit *LedgerBalanceLimitCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _LedgerBalanceLimit.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_LedgerBalanceLimit *LedgerBalanceLimitTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LedgerBalanceLimit.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_LedgerBalanceLimit *LedgerBalanceLimitTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _LedgerBalanceLimit.Contract.contract.Transact(opts, method, params...)
}
