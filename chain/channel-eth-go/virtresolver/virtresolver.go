// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package virtresolver

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

// VirtContractResolverMetaData contains all meta data concerning the VirtContractResolver contract.
var VirtContractResolverMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"deploy\",\"inputs\":[{\"name\":\"_code\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"_nonce\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"resolve\",\"inputs\":[{\"name\":\"_virtAddr\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"Deploy\",\"inputs\":[{\"name\":\"virtAddr\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"}],\"anonymous\":false}]",
	Bin: "0x6080604052348015600e575f5ffd5b506103648061001c5f395ff3fe608060405234801561000f575f5ffd5b5060043610610034575f3560e01c80635c23bdf5146100385780639c4ae2d014610068575b5f5ffd5b61004b610046366004610293565b61008b565b6040516001600160a01b0390911681526020015b60405180910390f35b61007b6100763660046102aa565b61010e565b604051901515815260200161005f565b5f818152602081905260408120546001600160a01b03166100f35760405162461bcd60e51b815260206004820152601b60248201527f4e6f6e6578697374656e74207669727475616c2061646472657373000000000060448201526064015b60405180910390fd5b505f908152602081905260409020546001600160a01b031690565b5f5f8484846040516020016101259392919061031c565b6040516020818303038152906040528051906020012090505f85858080601f0160208091040260200160405190810160405280939291908181526020018383808284375f920182905250868152602081905260409020549394505050506001600160a01b0316156101d85760405162461bcd60e51b815260206004820152601d60248201527f43757272656e74207265616c2061646472657373206973206e6f74203000000060448201526064016100ea565b5f8151602083015ff090506001600160a01b0381166102395760405162461bcd60e51b815260206004820152601760248201527f43726561746520636f6e7472616374206661696c65642e00000000000000000060448201526064016100ea565b5f8381526020819052604080822080546001600160a01b0319166001600160a01b0385161790555184917f149208daa30a9306858cc9c171c3510e0e50ab5d59ed2027a37a728430dd02e491a25060019695505050505050565b5f602082840312156102a3575f5ffd5b5035919050565b5f5f5f604084860312156102bc575f5ffd5b833567ffffffffffffffff8111156102d2575f5ffd5b8401601f810186136102e2575f5ffd5b803567ffffffffffffffff8111156102f8575f5ffd5b866020828401011115610309575f5ffd5b6020918201979096509401359392505050565b8284823790910190815260200191905056fea26469706673582212201a345ef0281e1b23f511658aaf7664c8efc6b9a786cc0b36076021a932afe3e664736f6c634300081e0033",
}

// VirtContractResolverABI is the input ABI used to generate the binding from.
// Deprecated: Use VirtContractResolverMetaData.ABI instead.
var VirtContractResolverABI = VirtContractResolverMetaData.ABI

// VirtContractResolverBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use VirtContractResolverMetaData.Bin instead.
var VirtContractResolverBin = VirtContractResolverMetaData.Bin

// DeployVirtContractResolver deploys a new Ethereum contract, binding an instance of VirtContractResolver to it.
func DeployVirtContractResolver(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *VirtContractResolver, error) {
	parsed, err := VirtContractResolverMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(VirtContractResolverBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &VirtContractResolver{VirtContractResolverCaller: VirtContractResolverCaller{contract: contract}, VirtContractResolverTransactor: VirtContractResolverTransactor{contract: contract}, VirtContractResolverFilterer: VirtContractResolverFilterer{contract: contract}}, nil
}

// VirtContractResolver is an auto generated Go binding around an Ethereum contract.
type VirtContractResolver struct {
	VirtContractResolverCaller     // Read-only binding to the contract
	VirtContractResolverTransactor // Write-only binding to the contract
	VirtContractResolverFilterer   // Log filterer for contract events
}

// VirtContractResolverCaller is an auto generated read-only Go binding around an Ethereum contract.
type VirtContractResolverCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VirtContractResolverTransactor is an auto generated write-only Go binding around an Ethereum contract.
type VirtContractResolverTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VirtContractResolverFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type VirtContractResolverFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VirtContractResolverSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type VirtContractResolverSession struct {
	Contract     *VirtContractResolver // Generic contract binding to set the session for
	CallOpts     bind.CallOpts         // Call options to use throughout this session
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// VirtContractResolverCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type VirtContractResolverCallerSession struct {
	Contract *VirtContractResolverCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts               // Call options to use throughout this session
}

// VirtContractResolverTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type VirtContractResolverTransactorSession struct {
	Contract     *VirtContractResolverTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts               // Transaction auth options to use throughout this session
}

// VirtContractResolverRaw is an auto generated low-level Go binding around an Ethereum contract.
type VirtContractResolverRaw struct {
	Contract *VirtContractResolver // Generic contract binding to access the raw methods on
}

// VirtContractResolverCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type VirtContractResolverCallerRaw struct {
	Contract *VirtContractResolverCaller // Generic read-only contract binding to access the raw methods on
}

// VirtContractResolverTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type VirtContractResolverTransactorRaw struct {
	Contract *VirtContractResolverTransactor // Generic write-only contract binding to access the raw methods on
}

// NewVirtContractResolver creates a new instance of VirtContractResolver, bound to a specific deployed contract.
func NewVirtContractResolver(address common.Address, backend bind.ContractBackend) (*VirtContractResolver, error) {
	contract, err := bindVirtContractResolver(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &VirtContractResolver{VirtContractResolverCaller: VirtContractResolverCaller{contract: contract}, VirtContractResolverTransactor: VirtContractResolverTransactor{contract: contract}, VirtContractResolverFilterer: VirtContractResolverFilterer{contract: contract}}, nil
}

// NewVirtContractResolverCaller creates a new read-only instance of VirtContractResolver, bound to a specific deployed contract.
func NewVirtContractResolverCaller(address common.Address, caller bind.ContractCaller) (*VirtContractResolverCaller, error) {
	contract, err := bindVirtContractResolver(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &VirtContractResolverCaller{contract: contract}, nil
}

// NewVirtContractResolverTransactor creates a new write-only instance of VirtContractResolver, bound to a specific deployed contract.
func NewVirtContractResolverTransactor(address common.Address, transactor bind.ContractTransactor) (*VirtContractResolverTransactor, error) {
	contract, err := bindVirtContractResolver(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &VirtContractResolverTransactor{contract: contract}, nil
}

// NewVirtContractResolverFilterer creates a new log filterer instance of VirtContractResolver, bound to a specific deployed contract.
func NewVirtContractResolverFilterer(address common.Address, filterer bind.ContractFilterer) (*VirtContractResolverFilterer, error) {
	contract, err := bindVirtContractResolver(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &VirtContractResolverFilterer{contract: contract}, nil
}

// bindVirtContractResolver binds a generic wrapper to an already deployed contract.
func bindVirtContractResolver(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := VirtContractResolverMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_VirtContractResolver *VirtContractResolverRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VirtContractResolver.Contract.VirtContractResolverCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_VirtContractResolver *VirtContractResolverRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VirtContractResolver.Contract.VirtContractResolverTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_VirtContractResolver *VirtContractResolverRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VirtContractResolver.Contract.VirtContractResolverTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_VirtContractResolver *VirtContractResolverCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VirtContractResolver.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_VirtContractResolver *VirtContractResolverTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VirtContractResolver.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_VirtContractResolver *VirtContractResolverTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VirtContractResolver.Contract.contract.Transact(opts, method, params...)
}

// Resolve is a free data retrieval call binding the contract method 0x5c23bdf5.
//
// Solidity: function resolve(bytes32 _virtAddr) view returns(address)
func (_VirtContractResolver *VirtContractResolverCaller) Resolve(opts *bind.CallOpts, _virtAddr [32]byte) (common.Address, error) {
	var out []interface{}
	err := _VirtContractResolver.contract.Call(opts, &out, "resolve", _virtAddr)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Resolve is a free data retrieval call binding the contract method 0x5c23bdf5.
//
// Solidity: function resolve(bytes32 _virtAddr) view returns(address)
func (_VirtContractResolver *VirtContractResolverSession) Resolve(_virtAddr [32]byte) (common.Address, error) {
	return _VirtContractResolver.Contract.Resolve(&_VirtContractResolver.CallOpts, _virtAddr)
}

// Resolve is a free data retrieval call binding the contract method 0x5c23bdf5.
//
// Solidity: function resolve(bytes32 _virtAddr) view returns(address)
func (_VirtContractResolver *VirtContractResolverCallerSession) Resolve(_virtAddr [32]byte) (common.Address, error) {
	return _VirtContractResolver.Contract.Resolve(&_VirtContractResolver.CallOpts, _virtAddr)
}

// Deploy is a paid mutator transaction binding the contract method 0x9c4ae2d0.
//
// Solidity: function deploy(bytes _code, uint256 _nonce) returns(bool)
func (_VirtContractResolver *VirtContractResolverTransactor) Deploy(opts *bind.TransactOpts, _code []byte, _nonce *big.Int) (*types.Transaction, error) {
	return _VirtContractResolver.contract.Transact(opts, "deploy", _code, _nonce)
}

// Deploy is a paid mutator transaction binding the contract method 0x9c4ae2d0.
//
// Solidity: function deploy(bytes _code, uint256 _nonce) returns(bool)
func (_VirtContractResolver *VirtContractResolverSession) Deploy(_code []byte, _nonce *big.Int) (*types.Transaction, error) {
	return _VirtContractResolver.Contract.Deploy(&_VirtContractResolver.TransactOpts, _code, _nonce)
}

// Deploy is a paid mutator transaction binding the contract method 0x9c4ae2d0.
//
// Solidity: function deploy(bytes _code, uint256 _nonce) returns(bool)
func (_VirtContractResolver *VirtContractResolverTransactorSession) Deploy(_code []byte, _nonce *big.Int) (*types.Transaction, error) {
	return _VirtContractResolver.Contract.Deploy(&_VirtContractResolver.TransactOpts, _code, _nonce)
}

// VirtContractResolverDeployIterator is returned from FilterDeploy and is used to iterate over the raw logs and unpacked data for Deploy events raised by the VirtContractResolver contract.
type VirtContractResolverDeployIterator struct {
	Event *VirtContractResolverDeploy // Event containing the contract specifics and raw log

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
func (it *VirtContractResolverDeployIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VirtContractResolverDeploy)
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
		it.Event = new(VirtContractResolverDeploy)
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
func (it *VirtContractResolverDeployIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *VirtContractResolverDeployIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// VirtContractResolverDeploy represents a Deploy event raised by the VirtContractResolver contract.
type VirtContractResolverDeploy struct {
	VirtAddr [32]byte
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterDeploy is a free log retrieval operation binding the contract event 0x149208daa30a9306858cc9c171c3510e0e50ab5d59ed2027a37a728430dd02e4.
//
// Solidity: event Deploy(bytes32 indexed virtAddr)
func (_VirtContractResolver *VirtContractResolverFilterer) FilterDeploy(opts *bind.FilterOpts, virtAddr [][32]byte) (*VirtContractResolverDeployIterator, error) {

	var virtAddrRule []interface{}
	for _, virtAddrItem := range virtAddr {
		virtAddrRule = append(virtAddrRule, virtAddrItem)
	}

	logs, sub, err := _VirtContractResolver.contract.FilterLogs(opts, "Deploy", virtAddrRule)
	if err != nil {
		return nil, err
	}
	return &VirtContractResolverDeployIterator{contract: _VirtContractResolver.contract, event: "Deploy", logs: logs, sub: sub}, nil
}

// WatchDeploy is a free log subscription operation binding the contract event 0x149208daa30a9306858cc9c171c3510e0e50ab5d59ed2027a37a728430dd02e4.
//
// Solidity: event Deploy(bytes32 indexed virtAddr)
func (_VirtContractResolver *VirtContractResolverFilterer) WatchDeploy(opts *bind.WatchOpts, sink chan<- *VirtContractResolverDeploy, virtAddr [][32]byte) (event.Subscription, error) {

	var virtAddrRule []interface{}
	for _, virtAddrItem := range virtAddr {
		virtAddrRule = append(virtAddrRule, virtAddrItem)
	}

	logs, sub, err := _VirtContractResolver.contract.WatchLogs(opts, "Deploy", virtAddrRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(VirtContractResolverDeploy)
				if err := _VirtContractResolver.contract.UnpackLog(event, "Deploy", log); err != nil {
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

// ParseDeploy is a log parse operation binding the contract event 0x149208daa30a9306858cc9c171c3510e0e50ab5d59ed2027a37a728430dd02e4.
//
// Solidity: event Deploy(bytes32 indexed virtAddr)
func (_VirtContractResolver *VirtContractResolverFilterer) ParseDeploy(log types.Log) (*VirtContractResolverDeploy, error) {
	event := new(VirtContractResolverDeploy)
	if err := _VirtContractResolver.contract.UnpackLog(event, "Deploy", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
