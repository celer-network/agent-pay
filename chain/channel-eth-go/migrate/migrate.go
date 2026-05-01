// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package migrate

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

// LedgerMigrateMetaData contains all meta data concerning the LedgerMigrate contract.
var LedgerMigrateMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"event\",\"name\":\"MigrateChannelFrom\",\"inputs\":[{\"name\":\"channelId\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"oldLedgerAddr\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"MigrateChannelTo\",\"inputs\":[{\"name\":\"channelId\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"newLedgerAddr\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"ECDSAInvalidSignature\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ECDSAInvalidSignatureLength\",\"inputs\":[{\"name\":\"length\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ECDSAInvalidSignatureS\",\"inputs\":[{\"name\":\"s\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}]",
	Bin: "0x611534610034600b8282823980515f1a607314602857634e487b7160e01b5f525f60045260245ffd5b305f52607381538281f3fe730000000000000000000000000000000000000000301460806040526004361061003f575f3560e01c80633c50ec721461004357806382b4338a14610074575b5f5ffd5b81801561004e575f5ffd5b5061006261005d36600461119d565b610095565b60405190815260200160405180910390f35b81801561007f575f5ffd5b5061009361008e3660046111fc565b61038f565b005b5f5f6100d584848080601f0160208091040260200160405190810160405280939291908181526020018383808284375f920191909152506105a692505050565b90505f6100e4825f0151610741565b80515f818152600689016020526040908190209083015192935090916001600383015460ff16600481111561011b5761011b611254565b148061013f57506002600383015460ff16600481111561013d5761013d611254565b145b610147575f5ffd5b8451805160209182012090860151610162908490839061083a565b6101aa5760405162461bcd60e51b815260206004820152601460248201527310da1958dac818dbcb5cda59dcc819985a5b195960621b60448201526064015b60405180910390fd5b60208501516001600160a01b031630146102065760405162461bcd60e51b815260206004820152601f60248201527f46726f6d206c65646765722061646472657373206973206e6f7420746869730060448201526064016101a1565b6001600160a01b038216331461026a5760405162461bcd60e51b815260206004820152602360248201527f546f206c65646765722061646472657373206973206e6f74206d73672e73656e6044820152623232b960e91b60648201526084016101a1565b84606001514211156102be5760405162461bcd60e51b815260206004820152601960248201527f506173736564206d6967726174696f6e20646561646c696e650000000000000060448201526064016101a1565b6102ca8a8460046108ff565b600383018054610100600160a81b0319166101006001600160a01b0385169081029190911790915560405185907fdefb8a94bbfc44ef5297b035407a7dd1314f369e39c3301f5b90f8810fb9fe4f905f90a360038a015460405163283226a360e21b8152600481018690526001600160a01b0384811660248301529091169063a0c89a8c906044015f604051808303815f87803b158015610369575f5ffd5b505af115801561037b573d5f5f3e3d5ffd5b5095985050505050505050505b9392505050565b60405163e0a515b760e01b815283905f906001600160a01b0383169063e0a515b7906103c19087908790600401611268565b6020604051808303815f875af11580156103dd573d5f5f3e3d5ffd5b505050506040513d601f19601f820116820180604052508101906104019190611296565b5f8181526006880160205260408120919250600382015460ff16600481111561042c5761042c611254565b146104835760405162461bcd60e51b815260206004820152602160248201527f496d6d69677261746564206368616e6e656c20616c72656164792065786973746044820152607360f81b60648201526084016101a1565b6003870154604051632a5a97e560e21b81526004810184905230916001600160a01b03169063a96a5f9490602401602060405180830381865afa1580156104cc573d5f5f3e3d5ffd5b505050506040513d601f19601f820116820180604052508101906104f091906112ad565b6001600160a01b0316146105465760405162461bcd60e51b815260206004820152601c60248201527f4f70657261746f7273686970206e6f74207472616e736665727265640000000060448201526064016101a1565b610552878260016108ff565b61055d818484610a36565b610568818484610b18565b6040516001600160a01b0387169083907f141a72a1d915a7c4205104b6e564cc991aa827c5f2c672a5d6a1da8bef99d6eb905f90a350505050505050565b60408051808201909152606080825260208201525f6105d783604080518082019091525f8152602081019190915290565b90505f6105e5826002610c97565b9050806002815181106105fa576105fa6112c8565b602002602001015167ffffffffffffffff81111561061a5761061a6112dc565b60405190808252806020026020018201604052801561064d57816020015b60608152602001906001900390816106385790505b5083602001819052505f8160028151811061066a5761066a6112c8565b6020026020010181815250505f5f5b602084015151845110156107385761069084610d50565b9092509050816001036106ad576106a684610d88565b8552610679565b81600203610729576106be84610d88565b8560200151846002815181106106d6576106d66112c8565b6020026020010151815181106106ee576106ee6112c8565b60200260200101819052508260028151811061070c5761070c6112c8565b60200260200101805180919061072190611304565b905250610679565b6107338482610e40565b610679565b50505050919050565b604080516080810182525f80825260208083018290528284018290526060830182905283518085019094528184528301849052909190805b602083015151835110156108325761079083610d50565b9092509050816001036107b5576107ae6107a984610d88565b610eb5565b8452610779565b816002036107e1576107ce6107c984610d88565b610ecb565b6001600160a01b03166020850152610779565b81600303610808576107f56107c984610d88565b6001600160a01b03166040850152610779565b816004036108235761081983610edb565b6060850152610779565b61082d8382610e40565b610779565b505050919050565b5f815160021461084b57505f610388565b7f19457468657265756d205369676e6564204d6573736167653a0a3332000000005f908152601c849052603c812090805b60028110156108f2576108b185828151811061089a5761089a6112c8565b602002602001015184610f4a90919063ffffffff16565b91508660040181600281106108c8576108c86112c8565b60080201546001600160a01b038381169116146108ea575f9350505050610388565b60010161087c565b5060019695505050505050565b80600481111561091157610911611254565b600383015460ff16600481111561092a5761092a611254565b0361093457505050565b5f600383015460ff16600481111561094e5761094e611254565b146109b957600382015460019084905f9060ff16600481111561097357610973611254565b81526020019081526020015f205461098b919061131c565b600383015484905f9060ff1660048111156109a8576109a8611254565b815260208101919091526040015f20555b825f8260048111156109cd576109cd611254565b81526020019081526020015f205460016109e7919061132f565b835f8360048111156109fb576109fb611254565b815260208101919091526040015f205560038201805482919060ff19166001836004811115610a2c57610a2c611254565b0217905550505050565b604051630bc2b0c160e21b8152600481018290525f906001600160a01b03841690632f0ac30490602401608060405180830381865afa158015610a7b573d5f5f3e3d5ffd5b505050506040513d601f19601f82011682018060405250810190610a9f9190611342565b6014880155600280880180546001600160a01b0390931661010002610100600160a81b031990931692909217909155600187019290925591508190811115610ae957610ae9611254565b60028086018054909160ff19909116906001908490811115610b0d57610b0d611254565b021790555050505050565b5f5f5f5f5f5f876001600160a01b03166388f41465886040518263ffffffff1660e01b8152600401610b4c91815260200190565b61018060405180830381865afa158015610b68573d5f5f3e3d5ffd5b505050506040513d601f19601f82011682018060405250810190610b8c919061140b565b9550955095509550955095505f5f90505b6002811015610c8b575f8a6004018260028110610bbc57610bbc6112c8565b600802019050878260028110610bd457610bd46112c8565b602002015181546001600160a01b0319166001600160a01b03909116178155868260028110610c0557610c056112c8565b60200201516001820155858260028110610c2157610c216112c8565b60200201518160020181905550848260028110610c4057610c406112c8565b60200201516003820155838260028110610c5c57610c5c6112c8565b60200201516004820155828260028110610c7857610c786112c8565b6020020151600790910155600101610b9d565b50505050505050505050565b8151606090610ca783600161132f565b67ffffffffffffffff811115610cbf57610cbf6112dc565b604051908082528060200260200182016040528015610ce8578160200160208202803683370190505b5091505f5f5b60208601515186511015610d4757610d0586610d50565b80925081935050506001848381518110610d2157610d216112c8565b60200260200101818151610d35919061132f565b905250610d428682610e40565b610cee565b50509092525090565b5f5f5f610d5c84610edb565b9050610d696008826114c8565b9250806007166005811115610d8057610d80611254565b915050915091565b60605f610d9483610edb565b90505f81845f0151610da6919061132f565b9050836020015151811115610db9575f5ffd5b8167ffffffffffffffff811115610dd257610dd26112dc565b6040519080825280601f01601f191660200182016040528015610dfc576020820181803683370190505b5060208086015186519295509181860191908301015f5b85811015610e35578181015183820152610e2e60208261132f565b9050610e13565b505050935250919050565b5f816005811115610e5357610e53611254565b03610e6657610e6182610edb565b505050565b6002816005811115610e7a57610e7a611254565b0361003f575f610e8983610edb565b905080835f01818151610e9c919061132f565b90525060208301515183511115610e61575f5ffd5b5050565b5f8151602014610ec3575f5ffd5b506020015190565b5f610ed582610f72565b92915050565b60208082015182518101909101515f9182805b600a81101561003f5783811a9150610f078160076114e7565b82607f16901b85179450816080165f03610f4257610f2681600161132f565b86518790610f3590839061132f565b9052509395945050505050565b600101610eee565b5f5f5f5f610f588686610f8f565b925092509250610f688282610fd8565b5090949350505050565b5f8151601414610f80575f5ffd5b5060200151600160601b900490565b5f5f5f8351604103610fc6576020840151604085015160608601515f1a610fb888828585611090565b955095509550505050610fd1565b505081515f91506002905b9250925092565b5f826003811115610feb57610feb611254565b03610ff4575050565b600182600381111561100857611008611254565b036110265760405163f645eedf60e01b815260040160405180910390fd5b600282600381111561103a5761103a611254565b0361105b5760405163fce698f760e01b8152600481018290526024016101a1565b600382600381111561106f5761106f611254565b03610eb1576040516335e2f38360e21b8152600481018290526024016101a1565b5f80807f7fffffffffffffffffffffffffffffff5d576e7357a4501ddfe92f46681b20a08411156110c957505f9150600390508261114e565b604080515f808252602082018084528a905260ff891692820192909252606081018790526080810186905260019060a0016020604051602081039080840390855afa15801561111a573d5f5f3e3d5ffd5b5050604051601f1901519150506001600160a01b03811661114557505f92506001915082905061114e565b92505f91508190505b9450945094915050565b5f5f83601f840112611168575f5ffd5b50813567ffffffffffffffff81111561117f575f5ffd5b602083019150836020828501011115611196575f5ffd5b9250929050565b5f5f5f604084860312156111af575f5ffd5b83359250602084013567ffffffffffffffff8111156111cc575f5ffd5b6111d886828701611158565b9497909650939450505050565b6001600160a01b03811681146111f9575f5ffd5b50565b5f5f5f5f6060858703121561120f575f5ffd5b843593506020850135611221816111e5565b9250604085013567ffffffffffffffff81111561123c575f5ffd5b61124887828801611158565b95989497509550505050565b634e487b7160e01b5f52602160045260245ffd5b60208152816020820152818360408301375f818301604090810191909152601f909201601f19160101919050565b5f602082840312156112a6575f5ffd5b5051919050565b5f602082840312156112bd575f5ffd5b8151610388816111e5565b634e487b7160e01b5f52603260045260245ffd5b634e487b7160e01b5f52604160045260245ffd5b634e487b7160e01b5f52601160045260245ffd5b5f60018201611315576113156112f0565b5060010190565b81810381811115610ed557610ed56112f0565b80820180821115610ed557610ed56112f0565b5f5f5f5f60808587031215611355575f5ffd5b845160208601516040870151919550935061136f816111e5565b6060959095015193969295505050565b604051601f8201601f1916810167ffffffffffffffff811182821017156113b457634e487b7160e01b5f52604160045260245ffd5b604052919050565b5f82601f8301126113cb575f5ffd5b6113d5604061137f565b8060408401858111156113e6575f5ffd5b845b818110156114005780518452602093840193016113e8565b509095945050505050565b5f5f5f5f5f5f6101808789031215611421575f5ffd5b87601f88011261142f575f5ffd5b611439604061137f565b80604089018a81111561144a575f5ffd5b895b8181101561146d57805161145f816111e5565b84526020938401930161144c565b5081985061147b8b826113bc565b975050505061148d88608089016113bc565b935061149c8860c089016113bc565b92506114ac8861010089016113bc565b91506114bc8861014089016113bc565b90509295509295509295565b5f826114e257634e487b7160e01b5f52601260045260245ffd5b500490565b8082028115828204841417610ed557610ed56112f056fea2646970667358221220c1c39f09a79ffbaba558bdaebfb0e1c26c1ef42552397fa776a6dd23857a69fb64736f6c634300081e0033",
}

// LedgerMigrateABI is the input ABI used to generate the binding from.
// Deprecated: Use LedgerMigrateMetaData.ABI instead.
var LedgerMigrateABI = LedgerMigrateMetaData.ABI

// LedgerMigrateBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use LedgerMigrateMetaData.Bin instead.
var LedgerMigrateBin = LedgerMigrateMetaData.Bin

// DeployLedgerMigrate deploys a new Ethereum contract, binding an instance of LedgerMigrate to it.
func DeployLedgerMigrate(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *LedgerMigrate, error) {
	parsed, err := LedgerMigrateMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(LedgerMigrateBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &LedgerMigrate{LedgerMigrateCaller: LedgerMigrateCaller{contract: contract}, LedgerMigrateTransactor: LedgerMigrateTransactor{contract: contract}, LedgerMigrateFilterer: LedgerMigrateFilterer{contract: contract}}, nil
}

// LedgerMigrate is an auto generated Go binding around an Ethereum contract.
type LedgerMigrate struct {
	LedgerMigrateCaller     // Read-only binding to the contract
	LedgerMigrateTransactor // Write-only binding to the contract
	LedgerMigrateFilterer   // Log filterer for contract events
}

// LedgerMigrateCaller is an auto generated read-only Go binding around an Ethereum contract.
type LedgerMigrateCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LedgerMigrateTransactor is an auto generated write-only Go binding around an Ethereum contract.
type LedgerMigrateTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LedgerMigrateFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type LedgerMigrateFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LedgerMigrateSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type LedgerMigrateSession struct {
	Contract     *LedgerMigrate    // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// LedgerMigrateCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type LedgerMigrateCallerSession struct {
	Contract *LedgerMigrateCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts        // Call options to use throughout this session
}

// LedgerMigrateTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type LedgerMigrateTransactorSession struct {
	Contract     *LedgerMigrateTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// LedgerMigrateRaw is an auto generated low-level Go binding around an Ethereum contract.
type LedgerMigrateRaw struct {
	Contract *LedgerMigrate // Generic contract binding to access the raw methods on
}

// LedgerMigrateCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type LedgerMigrateCallerRaw struct {
	Contract *LedgerMigrateCaller // Generic read-only contract binding to access the raw methods on
}

// LedgerMigrateTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type LedgerMigrateTransactorRaw struct {
	Contract *LedgerMigrateTransactor // Generic write-only contract binding to access the raw methods on
}

// NewLedgerMigrate creates a new instance of LedgerMigrate, bound to a specific deployed contract.
func NewLedgerMigrate(address common.Address, backend bind.ContractBackend) (*LedgerMigrate, error) {
	contract, err := bindLedgerMigrate(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &LedgerMigrate{LedgerMigrateCaller: LedgerMigrateCaller{contract: contract}, LedgerMigrateTransactor: LedgerMigrateTransactor{contract: contract}, LedgerMigrateFilterer: LedgerMigrateFilterer{contract: contract}}, nil
}

// NewLedgerMigrateCaller creates a new read-only instance of LedgerMigrate, bound to a specific deployed contract.
func NewLedgerMigrateCaller(address common.Address, caller bind.ContractCaller) (*LedgerMigrateCaller, error) {
	contract, err := bindLedgerMigrate(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &LedgerMigrateCaller{contract: contract}, nil
}

// NewLedgerMigrateTransactor creates a new write-only instance of LedgerMigrate, bound to a specific deployed contract.
func NewLedgerMigrateTransactor(address common.Address, transactor bind.ContractTransactor) (*LedgerMigrateTransactor, error) {
	contract, err := bindLedgerMigrate(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &LedgerMigrateTransactor{contract: contract}, nil
}

// NewLedgerMigrateFilterer creates a new log filterer instance of LedgerMigrate, bound to a specific deployed contract.
func NewLedgerMigrateFilterer(address common.Address, filterer bind.ContractFilterer) (*LedgerMigrateFilterer, error) {
	contract, err := bindLedgerMigrate(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &LedgerMigrateFilterer{contract: contract}, nil
}

// bindLedgerMigrate binds a generic wrapper to an already deployed contract.
func bindLedgerMigrate(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := LedgerMigrateMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_LedgerMigrate *LedgerMigrateRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _LedgerMigrate.Contract.LedgerMigrateCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_LedgerMigrate *LedgerMigrateRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LedgerMigrate.Contract.LedgerMigrateTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_LedgerMigrate *LedgerMigrateRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _LedgerMigrate.Contract.LedgerMigrateTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_LedgerMigrate *LedgerMigrateCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _LedgerMigrate.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_LedgerMigrate *LedgerMigrateTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LedgerMigrate.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_LedgerMigrate *LedgerMigrateTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _LedgerMigrate.Contract.contract.Transact(opts, method, params...)
}

// LedgerMigrateMigrateChannelFromIterator is returned from FilterMigrateChannelFrom and is used to iterate over the raw logs and unpacked data for MigrateChannelFrom events raised by the LedgerMigrate contract.
type LedgerMigrateMigrateChannelFromIterator struct {
	Event *LedgerMigrateMigrateChannelFrom // Event containing the contract specifics and raw log

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
func (it *LedgerMigrateMigrateChannelFromIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LedgerMigrateMigrateChannelFrom)
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
		it.Event = new(LedgerMigrateMigrateChannelFrom)
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
func (it *LedgerMigrateMigrateChannelFromIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LedgerMigrateMigrateChannelFromIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LedgerMigrateMigrateChannelFrom represents a MigrateChannelFrom event raised by the LedgerMigrate contract.
type LedgerMigrateMigrateChannelFrom struct {
	ChannelId     [32]byte
	OldLedgerAddr common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterMigrateChannelFrom is a free log retrieval operation binding the contract event 0x141a72a1d915a7c4205104b6e564cc991aa827c5f2c672a5d6a1da8bef99d6eb.
//
// Solidity: event MigrateChannelFrom(bytes32 indexed channelId, address indexed oldLedgerAddr)
func (_LedgerMigrate *LedgerMigrateFilterer) FilterMigrateChannelFrom(opts *bind.FilterOpts, channelId [][32]byte, oldLedgerAddr []common.Address) (*LedgerMigrateMigrateChannelFromIterator, error) {

	var channelIdRule []interface{}
	for _, channelIdItem := range channelId {
		channelIdRule = append(channelIdRule, channelIdItem)
	}
	var oldLedgerAddrRule []interface{}
	for _, oldLedgerAddrItem := range oldLedgerAddr {
		oldLedgerAddrRule = append(oldLedgerAddrRule, oldLedgerAddrItem)
	}

	logs, sub, err := _LedgerMigrate.contract.FilterLogs(opts, "MigrateChannelFrom", channelIdRule, oldLedgerAddrRule)
	if err != nil {
		return nil, err
	}
	return &LedgerMigrateMigrateChannelFromIterator{contract: _LedgerMigrate.contract, event: "MigrateChannelFrom", logs: logs, sub: sub}, nil
}

// WatchMigrateChannelFrom is a free log subscription operation binding the contract event 0x141a72a1d915a7c4205104b6e564cc991aa827c5f2c672a5d6a1da8bef99d6eb.
//
// Solidity: event MigrateChannelFrom(bytes32 indexed channelId, address indexed oldLedgerAddr)
func (_LedgerMigrate *LedgerMigrateFilterer) WatchMigrateChannelFrom(opts *bind.WatchOpts, sink chan<- *LedgerMigrateMigrateChannelFrom, channelId [][32]byte, oldLedgerAddr []common.Address) (event.Subscription, error) {

	var channelIdRule []interface{}
	for _, channelIdItem := range channelId {
		channelIdRule = append(channelIdRule, channelIdItem)
	}
	var oldLedgerAddrRule []interface{}
	for _, oldLedgerAddrItem := range oldLedgerAddr {
		oldLedgerAddrRule = append(oldLedgerAddrRule, oldLedgerAddrItem)
	}

	logs, sub, err := _LedgerMigrate.contract.WatchLogs(opts, "MigrateChannelFrom", channelIdRule, oldLedgerAddrRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LedgerMigrateMigrateChannelFrom)
				if err := _LedgerMigrate.contract.UnpackLog(event, "MigrateChannelFrom", log); err != nil {
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

// ParseMigrateChannelFrom is a log parse operation binding the contract event 0x141a72a1d915a7c4205104b6e564cc991aa827c5f2c672a5d6a1da8bef99d6eb.
//
// Solidity: event MigrateChannelFrom(bytes32 indexed channelId, address indexed oldLedgerAddr)
func (_LedgerMigrate *LedgerMigrateFilterer) ParseMigrateChannelFrom(log types.Log) (*LedgerMigrateMigrateChannelFrom, error) {
	event := new(LedgerMigrateMigrateChannelFrom)
	if err := _LedgerMigrate.contract.UnpackLog(event, "MigrateChannelFrom", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// LedgerMigrateMigrateChannelToIterator is returned from FilterMigrateChannelTo and is used to iterate over the raw logs and unpacked data for MigrateChannelTo events raised by the LedgerMigrate contract.
type LedgerMigrateMigrateChannelToIterator struct {
	Event *LedgerMigrateMigrateChannelTo // Event containing the contract specifics and raw log

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
func (it *LedgerMigrateMigrateChannelToIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LedgerMigrateMigrateChannelTo)
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
		it.Event = new(LedgerMigrateMigrateChannelTo)
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
func (it *LedgerMigrateMigrateChannelToIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LedgerMigrateMigrateChannelToIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LedgerMigrateMigrateChannelTo represents a MigrateChannelTo event raised by the LedgerMigrate contract.
type LedgerMigrateMigrateChannelTo struct {
	ChannelId     [32]byte
	NewLedgerAddr common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterMigrateChannelTo is a free log retrieval operation binding the contract event 0xdefb8a94bbfc44ef5297b035407a7dd1314f369e39c3301f5b90f8810fb9fe4f.
//
// Solidity: event MigrateChannelTo(bytes32 indexed channelId, address indexed newLedgerAddr)
func (_LedgerMigrate *LedgerMigrateFilterer) FilterMigrateChannelTo(opts *bind.FilterOpts, channelId [][32]byte, newLedgerAddr []common.Address) (*LedgerMigrateMigrateChannelToIterator, error) {

	var channelIdRule []interface{}
	for _, channelIdItem := range channelId {
		channelIdRule = append(channelIdRule, channelIdItem)
	}
	var newLedgerAddrRule []interface{}
	for _, newLedgerAddrItem := range newLedgerAddr {
		newLedgerAddrRule = append(newLedgerAddrRule, newLedgerAddrItem)
	}

	logs, sub, err := _LedgerMigrate.contract.FilterLogs(opts, "MigrateChannelTo", channelIdRule, newLedgerAddrRule)
	if err != nil {
		return nil, err
	}
	return &LedgerMigrateMigrateChannelToIterator{contract: _LedgerMigrate.contract, event: "MigrateChannelTo", logs: logs, sub: sub}, nil
}

// WatchMigrateChannelTo is a free log subscription operation binding the contract event 0xdefb8a94bbfc44ef5297b035407a7dd1314f369e39c3301f5b90f8810fb9fe4f.
//
// Solidity: event MigrateChannelTo(bytes32 indexed channelId, address indexed newLedgerAddr)
func (_LedgerMigrate *LedgerMigrateFilterer) WatchMigrateChannelTo(opts *bind.WatchOpts, sink chan<- *LedgerMigrateMigrateChannelTo, channelId [][32]byte, newLedgerAddr []common.Address) (event.Subscription, error) {

	var channelIdRule []interface{}
	for _, channelIdItem := range channelId {
		channelIdRule = append(channelIdRule, channelIdItem)
	}
	var newLedgerAddrRule []interface{}
	for _, newLedgerAddrItem := range newLedgerAddr {
		newLedgerAddrRule = append(newLedgerAddrRule, newLedgerAddrItem)
	}

	logs, sub, err := _LedgerMigrate.contract.WatchLogs(opts, "MigrateChannelTo", channelIdRule, newLedgerAddrRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LedgerMigrateMigrateChannelTo)
				if err := _LedgerMigrate.contract.UnpackLog(event, "MigrateChannelTo", log); err != nil {
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

// ParseMigrateChannelTo is a log parse operation binding the contract event 0xdefb8a94bbfc44ef5297b035407a7dd1314f369e39c3301f5b90f8810fb9fe4f.
//
// Solidity: event MigrateChannelTo(bytes32 indexed channelId, address indexed newLedgerAddr)
func (_LedgerMigrate *LedgerMigrateFilterer) ParseMigrateChannelTo(log types.Log) (*LedgerMigrateMigrateChannelTo, error) {
	event := new(LedgerMigrateMigrateChannelTo)
	if err := _LedgerMigrate.contract.UnpackLog(event, "MigrateChannelTo", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
