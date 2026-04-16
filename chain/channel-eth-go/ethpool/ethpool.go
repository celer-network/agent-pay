// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package ethpool

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

// EthPoolMetaData contains all meta data concerning the EthPool contract.
var EthPoolMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"allowance\",\"inputs\":[{\"name\":\"_owner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_spender\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"approve\",\"inputs\":[{\"name\":\"_spender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"balanceOf\",\"inputs\":[{\"name\":\"_owner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"decimals\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"decreaseAllowance\",\"inputs\":[{\"name\":\"_spender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_subtractedValue\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"deposit\",\"inputs\":[{\"name\":\"_receiver\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"increaseAllowance\",\"inputs\":[{\"name\":\"_spender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_addedValue\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"name\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"symbol\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"transferFrom\",\"inputs\":[{\"name\":\"_from\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_to\",\"type\":\"address\",\"internalType\":\"addresspayable\"},{\"name\":\"_value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferToCelerWallet\",\"inputs\":[{\"name\":\"_from\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_walletAddr\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_walletId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"_value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"withdraw\",\"inputs\":[{\"name\":\"_value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"Approval\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"spender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Deposit\",\"inputs\":[{\"name\":\"receiver\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Transfer\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false}]",
	Bin: "0x6080604052348015600e575f5ffd5b50610a418061001c5f395ff3fe6080604052600436106100a5575f3560e01c806370a082311161006257806370a08231146101a75780637e1cd431146101e957806395d89b4114610208578063a457c2d714610238578063dd62ed3e14610257578063f340fa011461029b575f5ffd5b806306fdde03146100a9578063095ea7b3146100f357806323b872dd146101225780632e1a7d4d14610141578063313ce567146101625780633950935114610188575b5f5ffd5b3480156100b4575f5ffd5b506100dd60405180604001604052806009815260200168115d1a125b941bdbdb60ba1b81525081565b6040516100ea919061081f565b60405180910390f35b3480156100fe575f5ffd5b5061011261010d366004610868565b6102ae565b60405190151581526020016100ea565b34801561012d575f5ffd5b5061011261013c366004610892565b610330565b34801561014c575f5ffd5b5061016061015b3660046108d0565b6103b9565b005b34801561016d575f5ffd5b50610176601281565b60405160ff90911681526020016100ea565b348015610193575f5ffd5b506101126101a2366004610868565b6103c7565b3480156101b2575f5ffd5b506101db6101c13660046108e7565b6001600160a01b03165f9081526020819052604090205490565b6040519081526020016100ea565b3480156101f4575f5ffd5b50610112610203366004610909565b61045f565b348015610213575f5ffd5b506100dd60405180604001604052806005815260200164045746849560dc1b81525081565b348015610243575f5ffd5b50610112610252366004610868565b6105af565b348015610262575f5ffd5b506101db61027136600461094c565b6001600160a01b039182165f90815260016020908152604080832093909416825291909152205490565b6101606102a93660046108e7565b610604565b5f6001600160a01b0383166102de5760405162461bcd60e51b81526004016102d590610983565b60405180910390fd5b335f8181526001602090815260408083206001600160a01b03881680855290835292819020869055518581529192915f5160206109ec5f395f51905f5291015b60405180910390a35060015b92915050565b6001600160a01b0383165f90815260016020908152604080832033845290915281205461035e9083906109c5565b6001600160a01b0385165f81815260016020908152604080832033808552908352928190208590555193845290925f5160206109ec5f395f51905f52910160405180910390a36103af8484846106c6565b5060019392505050565b6103c43333836106c6565b50565b5f6001600160a01b0383166103ee5760405162461bcd60e51b81526004016102d590610983565b335f9081526001602090815260408083206001600160a01b038716845290915290205461041c9083906109d8565b335f8181526001602090815260408083206001600160a01b038916808552908352928190208590555193845290925f5160206109ec5f395f51905f52910161031e565b6001600160a01b0384165f90815260016020908152604080832033845290915281205461048d9083906109c5565b6001600160a01b0386165f81815260016020908152604080832033808552908352928190208590555193845290925f5160206109ec5f395f51905f52910160405180910390a36001600160a01b0385165f908152602081905260409020546104f69083906109c5565b6001600160a01b038681165f818152602081815260409182902094909455518581529187169290917fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef910160405180910390a3604051636b46cea760e11b81526004810184905284906001600160a01b0382169063d68d9d4e9085906024015f604051808303818588803b15801561058c575f5ffd5b505af115801561059e573d5f5f3e3d5ffd5b5060019a9950505050505050505050565b5f6001600160a01b0383166105d65760405162461bcd60e51b81526004016102d590610983565b335f9081526001602090815260408083206001600160a01b038716845290915290205461041c9083906109c5565b6001600160a01b0381166106525760405162461bcd60e51b8152602060048201526015602482015274052656365697665722061646472657373206973203605c1b60448201526064016102d5565b6001600160a01b0381165f908152602081905260409020546106759034906109d8565b6001600160a01b0382165f8181526020818152604091829020939093555134815290917fe1fffcc4923d04b559f4d29a8bfc6cda04eb5b0d3c460751c2402c5c5cc9109c910160405180910390a250565b6001600160a01b03821661070e5760405162461bcd60e51b815260206004820152600f60248201526e0546f2061646472657373206973203608c1b60448201526064016102d5565b6001600160a01b0383165f908152602081905260409020546107319082906109c5565b6001600160a01b038481165f818152602081815260409182902094909455518481529185169290917fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef910160405180910390a35f826001600160a01b0316826040515f6040518083038185875af1925050503d805f81146107cd576040519150601f19603f3d011682016040523d82523d5f602084013e6107d2565b606091505b50509050806108195760405162461bcd60e51b8152602060048201526013602482015272115512081d1c985b9cd9995c8819985a5b1959606a1b60448201526064016102d5565b50505050565b602081525f82518060208401528060208501604085015e5f604082850101526040601f19601f83011684010191505092915050565b6001600160a01b03811681146103c4575f5ffd5b5f5f60408385031215610879575f5ffd5b823561088481610854565b946020939093013593505050565b5f5f5f606084860312156108a4575f5ffd5b83356108af81610854565b925060208401356108bf81610854565b929592945050506040919091013590565b5f602082840312156108e0575f5ffd5b5035919050565b5f602082840312156108f7575f5ffd5b813561090281610854565b9392505050565b5f5f5f5f6080858703121561091c575f5ffd5b843561092781610854565b9350602085013561093781610854565b93969395505050506040820135916060013590565b5f5f6040838503121561095d575f5ffd5b823561096881610854565b9150602083013561097881610854565b809150509250929050565b60208082526014908201527305370656e646572206164647265737320697320360641b604082015260600190565b634e487b7160e01b5f52601160045260245ffd5b8181038181111561032a5761032a6109b1565b8082018082111561032a5761032a6109b156fe8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925a2646970667358221220302e2596484806f0db3abb1d5cfff35659525a7cb8185f2c3858bb82d8ad9d4664736f6c634300081e0033",
}

// EthPoolABI is the input ABI used to generate the binding from.
// Deprecated: Use EthPoolMetaData.ABI instead.
var EthPoolABI = EthPoolMetaData.ABI

// EthPoolBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use EthPoolMetaData.Bin instead.
var EthPoolBin = EthPoolMetaData.Bin

// DeployEthPool deploys a new Ethereum contract, binding an instance of EthPool to it.
func DeployEthPool(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *EthPool, error) {
	parsed, err := EthPoolMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(EthPoolBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &EthPool{EthPoolCaller: EthPoolCaller{contract: contract}, EthPoolTransactor: EthPoolTransactor{contract: contract}, EthPoolFilterer: EthPoolFilterer{contract: contract}}, nil
}

// EthPool is an auto generated Go binding around an Ethereum contract.
type EthPool struct {
	EthPoolCaller     // Read-only binding to the contract
	EthPoolTransactor // Write-only binding to the contract
	EthPoolFilterer   // Log filterer for contract events
}

// EthPoolCaller is an auto generated read-only Go binding around an Ethereum contract.
type EthPoolCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EthPoolTransactor is an auto generated write-only Go binding around an Ethereum contract.
type EthPoolTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EthPoolFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type EthPoolFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EthPoolSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type EthPoolSession struct {
	Contract     *EthPool          // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// EthPoolCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type EthPoolCallerSession struct {
	Contract *EthPoolCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts  // Call options to use throughout this session
}

// EthPoolTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type EthPoolTransactorSession struct {
	Contract     *EthPoolTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// EthPoolRaw is an auto generated low-level Go binding around an Ethereum contract.
type EthPoolRaw struct {
	Contract *EthPool // Generic contract binding to access the raw methods on
}

// EthPoolCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type EthPoolCallerRaw struct {
	Contract *EthPoolCaller // Generic read-only contract binding to access the raw methods on
}

// EthPoolTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type EthPoolTransactorRaw struct {
	Contract *EthPoolTransactor // Generic write-only contract binding to access the raw methods on
}

// NewEthPool creates a new instance of EthPool, bound to a specific deployed contract.
func NewEthPool(address common.Address, backend bind.ContractBackend) (*EthPool, error) {
	contract, err := bindEthPool(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &EthPool{EthPoolCaller: EthPoolCaller{contract: contract}, EthPoolTransactor: EthPoolTransactor{contract: contract}, EthPoolFilterer: EthPoolFilterer{contract: contract}}, nil
}

// NewEthPoolCaller creates a new read-only instance of EthPool, bound to a specific deployed contract.
func NewEthPoolCaller(address common.Address, caller bind.ContractCaller) (*EthPoolCaller, error) {
	contract, err := bindEthPool(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &EthPoolCaller{contract: contract}, nil
}

// NewEthPoolTransactor creates a new write-only instance of EthPool, bound to a specific deployed contract.
func NewEthPoolTransactor(address common.Address, transactor bind.ContractTransactor) (*EthPoolTransactor, error) {
	contract, err := bindEthPool(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &EthPoolTransactor{contract: contract}, nil
}

// NewEthPoolFilterer creates a new log filterer instance of EthPool, bound to a specific deployed contract.
func NewEthPoolFilterer(address common.Address, filterer bind.ContractFilterer) (*EthPoolFilterer, error) {
	contract, err := bindEthPool(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &EthPoolFilterer{contract: contract}, nil
}

// bindEthPool binds a generic wrapper to an already deployed contract.
func bindEthPool(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := EthPoolMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_EthPool *EthPoolRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _EthPool.Contract.EthPoolCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_EthPool *EthPoolRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _EthPool.Contract.EthPoolTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_EthPool *EthPoolRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _EthPool.Contract.EthPoolTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_EthPool *EthPoolCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _EthPool.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_EthPool *EthPoolTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _EthPool.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_EthPool *EthPoolTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _EthPool.Contract.contract.Transact(opts, method, params...)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address _owner, address _spender) view returns(uint256)
func (_EthPool *EthPoolCaller) Allowance(opts *bind.CallOpts, _owner common.Address, _spender common.Address) (*big.Int, error) {
	var out []interface{}
	err := _EthPool.contract.Call(opts, &out, "allowance", _owner, _spender)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address _owner, address _spender) view returns(uint256)
func (_EthPool *EthPoolSession) Allowance(_owner common.Address, _spender common.Address) (*big.Int, error) {
	return _EthPool.Contract.Allowance(&_EthPool.CallOpts, _owner, _spender)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address _owner, address _spender) view returns(uint256)
func (_EthPool *EthPoolCallerSession) Allowance(_owner common.Address, _spender common.Address) (*big.Int, error) {
	return _EthPool.Contract.Allowance(&_EthPool.CallOpts, _owner, _spender)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address _owner) view returns(uint256)
func (_EthPool *EthPoolCaller) BalanceOf(opts *bind.CallOpts, _owner common.Address) (*big.Int, error) {
	var out []interface{}
	err := _EthPool.contract.Call(opts, &out, "balanceOf", _owner)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address _owner) view returns(uint256)
func (_EthPool *EthPoolSession) BalanceOf(_owner common.Address) (*big.Int, error) {
	return _EthPool.Contract.BalanceOf(&_EthPool.CallOpts, _owner)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address _owner) view returns(uint256)
func (_EthPool *EthPoolCallerSession) BalanceOf(_owner common.Address) (*big.Int, error) {
	return _EthPool.Contract.BalanceOf(&_EthPool.CallOpts, _owner)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_EthPool *EthPoolCaller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _EthPool.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_EthPool *EthPoolSession) Decimals() (uint8, error) {
	return _EthPool.Contract.Decimals(&_EthPool.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_EthPool *EthPoolCallerSession) Decimals() (uint8, error) {
	return _EthPool.Contract.Decimals(&_EthPool.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_EthPool *EthPoolCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _EthPool.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_EthPool *EthPoolSession) Name() (string, error) {
	return _EthPool.Contract.Name(&_EthPool.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_EthPool *EthPoolCallerSession) Name() (string, error) {
	return _EthPool.Contract.Name(&_EthPool.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_EthPool *EthPoolCaller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _EthPool.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_EthPool *EthPoolSession) Symbol() (string, error) {
	return _EthPool.Contract.Symbol(&_EthPool.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_EthPool *EthPoolCallerSession) Symbol() (string, error) {
	return _EthPool.Contract.Symbol(&_EthPool.CallOpts)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address _spender, uint256 _value) returns(bool)
func (_EthPool *EthPoolTransactor) Approve(opts *bind.TransactOpts, _spender common.Address, _value *big.Int) (*types.Transaction, error) {
	return _EthPool.contract.Transact(opts, "approve", _spender, _value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address _spender, uint256 _value) returns(bool)
func (_EthPool *EthPoolSession) Approve(_spender common.Address, _value *big.Int) (*types.Transaction, error) {
	return _EthPool.Contract.Approve(&_EthPool.TransactOpts, _spender, _value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address _spender, uint256 _value) returns(bool)
func (_EthPool *EthPoolTransactorSession) Approve(_spender common.Address, _value *big.Int) (*types.Transaction, error) {
	return _EthPool.Contract.Approve(&_EthPool.TransactOpts, _spender, _value)
}

// DecreaseAllowance is a paid mutator transaction binding the contract method 0xa457c2d7.
//
// Solidity: function decreaseAllowance(address _spender, uint256 _subtractedValue) returns(bool)
func (_EthPool *EthPoolTransactor) DecreaseAllowance(opts *bind.TransactOpts, _spender common.Address, _subtractedValue *big.Int) (*types.Transaction, error) {
	return _EthPool.contract.Transact(opts, "decreaseAllowance", _spender, _subtractedValue)
}

// DecreaseAllowance is a paid mutator transaction binding the contract method 0xa457c2d7.
//
// Solidity: function decreaseAllowance(address _spender, uint256 _subtractedValue) returns(bool)
func (_EthPool *EthPoolSession) DecreaseAllowance(_spender common.Address, _subtractedValue *big.Int) (*types.Transaction, error) {
	return _EthPool.Contract.DecreaseAllowance(&_EthPool.TransactOpts, _spender, _subtractedValue)
}

// DecreaseAllowance is a paid mutator transaction binding the contract method 0xa457c2d7.
//
// Solidity: function decreaseAllowance(address _spender, uint256 _subtractedValue) returns(bool)
func (_EthPool *EthPoolTransactorSession) DecreaseAllowance(_spender common.Address, _subtractedValue *big.Int) (*types.Transaction, error) {
	return _EthPool.Contract.DecreaseAllowance(&_EthPool.TransactOpts, _spender, _subtractedValue)
}

// Deposit is a paid mutator transaction binding the contract method 0xf340fa01.
//
// Solidity: function deposit(address _receiver) payable returns()
func (_EthPool *EthPoolTransactor) Deposit(opts *bind.TransactOpts, _receiver common.Address) (*types.Transaction, error) {
	return _EthPool.contract.Transact(opts, "deposit", _receiver)
}

// Deposit is a paid mutator transaction binding the contract method 0xf340fa01.
//
// Solidity: function deposit(address _receiver) payable returns()
func (_EthPool *EthPoolSession) Deposit(_receiver common.Address) (*types.Transaction, error) {
	return _EthPool.Contract.Deposit(&_EthPool.TransactOpts, _receiver)
}

// Deposit is a paid mutator transaction binding the contract method 0xf340fa01.
//
// Solidity: function deposit(address _receiver) payable returns()
func (_EthPool *EthPoolTransactorSession) Deposit(_receiver common.Address) (*types.Transaction, error) {
	return _EthPool.Contract.Deposit(&_EthPool.TransactOpts, _receiver)
}

// IncreaseAllowance is a paid mutator transaction binding the contract method 0x39509351.
//
// Solidity: function increaseAllowance(address _spender, uint256 _addedValue) returns(bool)
func (_EthPool *EthPoolTransactor) IncreaseAllowance(opts *bind.TransactOpts, _spender common.Address, _addedValue *big.Int) (*types.Transaction, error) {
	return _EthPool.contract.Transact(opts, "increaseAllowance", _spender, _addedValue)
}

// IncreaseAllowance is a paid mutator transaction binding the contract method 0x39509351.
//
// Solidity: function increaseAllowance(address _spender, uint256 _addedValue) returns(bool)
func (_EthPool *EthPoolSession) IncreaseAllowance(_spender common.Address, _addedValue *big.Int) (*types.Transaction, error) {
	return _EthPool.Contract.IncreaseAllowance(&_EthPool.TransactOpts, _spender, _addedValue)
}

// IncreaseAllowance is a paid mutator transaction binding the contract method 0x39509351.
//
// Solidity: function increaseAllowance(address _spender, uint256 _addedValue) returns(bool)
func (_EthPool *EthPoolTransactorSession) IncreaseAllowance(_spender common.Address, _addedValue *big.Int) (*types.Transaction, error) {
	return _EthPool.Contract.IncreaseAllowance(&_EthPool.TransactOpts, _spender, _addedValue)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address _from, address _to, uint256 _value) returns(bool)
func (_EthPool *EthPoolTransactor) TransferFrom(opts *bind.TransactOpts, _from common.Address, _to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _EthPool.contract.Transact(opts, "transferFrom", _from, _to, _value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address _from, address _to, uint256 _value) returns(bool)
func (_EthPool *EthPoolSession) TransferFrom(_from common.Address, _to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _EthPool.Contract.TransferFrom(&_EthPool.TransactOpts, _from, _to, _value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address _from, address _to, uint256 _value) returns(bool)
func (_EthPool *EthPoolTransactorSession) TransferFrom(_from common.Address, _to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _EthPool.Contract.TransferFrom(&_EthPool.TransactOpts, _from, _to, _value)
}

// TransferToCelerWallet is a paid mutator transaction binding the contract method 0x7e1cd431.
//
// Solidity: function transferToCelerWallet(address _from, address _walletAddr, bytes32 _walletId, uint256 _value) returns(bool)
func (_EthPool *EthPoolTransactor) TransferToCelerWallet(opts *bind.TransactOpts, _from common.Address, _walletAddr common.Address, _walletId [32]byte, _value *big.Int) (*types.Transaction, error) {
	return _EthPool.contract.Transact(opts, "transferToCelerWallet", _from, _walletAddr, _walletId, _value)
}

// TransferToCelerWallet is a paid mutator transaction binding the contract method 0x7e1cd431.
//
// Solidity: function transferToCelerWallet(address _from, address _walletAddr, bytes32 _walletId, uint256 _value) returns(bool)
func (_EthPool *EthPoolSession) TransferToCelerWallet(_from common.Address, _walletAddr common.Address, _walletId [32]byte, _value *big.Int) (*types.Transaction, error) {
	return _EthPool.Contract.TransferToCelerWallet(&_EthPool.TransactOpts, _from, _walletAddr, _walletId, _value)
}

// TransferToCelerWallet is a paid mutator transaction binding the contract method 0x7e1cd431.
//
// Solidity: function transferToCelerWallet(address _from, address _walletAddr, bytes32 _walletId, uint256 _value) returns(bool)
func (_EthPool *EthPoolTransactorSession) TransferToCelerWallet(_from common.Address, _walletAddr common.Address, _walletId [32]byte, _value *big.Int) (*types.Transaction, error) {
	return _EthPool.Contract.TransferToCelerWallet(&_EthPool.TransactOpts, _from, _walletAddr, _walletId, _value)
}

// Withdraw is a paid mutator transaction binding the contract method 0x2e1a7d4d.
//
// Solidity: function withdraw(uint256 _value) returns()
func (_EthPool *EthPoolTransactor) Withdraw(opts *bind.TransactOpts, _value *big.Int) (*types.Transaction, error) {
	return _EthPool.contract.Transact(opts, "withdraw", _value)
}

// Withdraw is a paid mutator transaction binding the contract method 0x2e1a7d4d.
//
// Solidity: function withdraw(uint256 _value) returns()
func (_EthPool *EthPoolSession) Withdraw(_value *big.Int) (*types.Transaction, error) {
	return _EthPool.Contract.Withdraw(&_EthPool.TransactOpts, _value)
}

// Withdraw is a paid mutator transaction binding the contract method 0x2e1a7d4d.
//
// Solidity: function withdraw(uint256 _value) returns()
func (_EthPool *EthPoolTransactorSession) Withdraw(_value *big.Int) (*types.Transaction, error) {
	return _EthPool.Contract.Withdraw(&_EthPool.TransactOpts, _value)
}

// EthPoolApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the EthPool contract.
type EthPoolApprovalIterator struct {
	Event *EthPoolApproval // Event containing the contract specifics and raw log

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
func (it *EthPoolApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EthPoolApproval)
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
		it.Event = new(EthPoolApproval)
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
func (it *EthPoolApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EthPoolApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EthPoolApproval represents a Approval event raised by the EthPool contract.
type EthPoolApproval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_EthPool *EthPoolFilterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*EthPoolApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _EthPool.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &EthPoolApprovalIterator{contract: _EthPool.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_EthPool *EthPoolFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *EthPoolApproval, owner []common.Address, spender []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _EthPool.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EthPoolApproval)
				if err := _EthPool.contract.UnpackLog(event, "Approval", log); err != nil {
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
func (_EthPool *EthPoolFilterer) ParseApproval(log types.Log) (*EthPoolApproval, error) {
	event := new(EthPoolApproval)
	if err := _EthPool.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// EthPoolDepositIterator is returned from FilterDeposit and is used to iterate over the raw logs and unpacked data for Deposit events raised by the EthPool contract.
type EthPoolDepositIterator struct {
	Event *EthPoolDeposit // Event containing the contract specifics and raw log

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
func (it *EthPoolDepositIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EthPoolDeposit)
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
		it.Event = new(EthPoolDeposit)
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
func (it *EthPoolDepositIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EthPoolDepositIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EthPoolDeposit represents a Deposit event raised by the EthPool contract.
type EthPoolDeposit struct {
	Receiver common.Address
	Value    *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterDeposit is a free log retrieval operation binding the contract event 0xe1fffcc4923d04b559f4d29a8bfc6cda04eb5b0d3c460751c2402c5c5cc9109c.
//
// Solidity: event Deposit(address indexed receiver, uint256 value)
func (_EthPool *EthPoolFilterer) FilterDeposit(opts *bind.FilterOpts, receiver []common.Address) (*EthPoolDepositIterator, error) {

	var receiverRule []interface{}
	for _, receiverItem := range receiver {
		receiverRule = append(receiverRule, receiverItem)
	}

	logs, sub, err := _EthPool.contract.FilterLogs(opts, "Deposit", receiverRule)
	if err != nil {
		return nil, err
	}
	return &EthPoolDepositIterator{contract: _EthPool.contract, event: "Deposit", logs: logs, sub: sub}, nil
}

// WatchDeposit is a free log subscription operation binding the contract event 0xe1fffcc4923d04b559f4d29a8bfc6cda04eb5b0d3c460751c2402c5c5cc9109c.
//
// Solidity: event Deposit(address indexed receiver, uint256 value)
func (_EthPool *EthPoolFilterer) WatchDeposit(opts *bind.WatchOpts, sink chan<- *EthPoolDeposit, receiver []common.Address) (event.Subscription, error) {

	var receiverRule []interface{}
	for _, receiverItem := range receiver {
		receiverRule = append(receiverRule, receiverItem)
	}

	logs, sub, err := _EthPool.contract.WatchLogs(opts, "Deposit", receiverRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EthPoolDeposit)
				if err := _EthPool.contract.UnpackLog(event, "Deposit", log); err != nil {
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

// ParseDeposit is a log parse operation binding the contract event 0xe1fffcc4923d04b559f4d29a8bfc6cda04eb5b0d3c460751c2402c5c5cc9109c.
//
// Solidity: event Deposit(address indexed receiver, uint256 value)
func (_EthPool *EthPoolFilterer) ParseDeposit(log types.Log) (*EthPoolDeposit, error) {
	event := new(EthPoolDeposit)
	if err := _EthPool.contract.UnpackLog(event, "Deposit", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// EthPoolTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the EthPool contract.
type EthPoolTransferIterator struct {
	Event *EthPoolTransfer // Event containing the contract specifics and raw log

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
func (it *EthPoolTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EthPoolTransfer)
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
		it.Event = new(EthPoolTransfer)
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
func (it *EthPoolTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EthPoolTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EthPoolTransfer represents a Transfer event raised by the EthPool contract.
type EthPoolTransfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_EthPool *EthPoolFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*EthPoolTransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _EthPool.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &EthPoolTransferIterator{contract: _EthPool.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_EthPool *EthPoolFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *EthPoolTransfer, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _EthPool.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EthPoolTransfer)
				if err := _EthPool.contract.UnpackLog(event, "Transfer", log); err != nil {
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
func (_EthPool *EthPoolFilterer) ParseTransfer(log types.Log) (*EthPoolTransfer, error) {
	event := new(EthPoolTransfer)
	if err := _EthPool.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
