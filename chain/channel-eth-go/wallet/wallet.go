// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package wallet

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

// CelerWalletMetaData contains all meta data concerning the CelerWallet contract.
var CelerWalletMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"create\",\"inputs\":[{\"name\":\"_owners\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"_operator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_nonce\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"depositERC20\",\"inputs\":[{\"name\":\"_walletId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"_tokenAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"depositETH\",\"inputs\":[{\"name\":\"_walletId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"drainToken\",\"inputs\":[{\"name\":\"_tokenAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_receiver\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getBalance\",\"inputs\":[{\"name\":\"_walletId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"_tokenAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getOperator\",\"inputs\":[{\"name\":\"_walletId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getProposalVote\",\"inputs\":[{\"name\":\"_walletId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"_owner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getProposedNewOperator\",\"inputs\":[{\"name\":\"_walletId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getWalletOwners\",\"inputs\":[{\"name\":\"_walletId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pause\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"paused\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"proposeNewOperator\",\"inputs\":[{\"name\":\"_walletId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"_newOperator\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferOperatorship\",\"inputs\":[{\"name\":\"_walletId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"_newOperator\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferToWallet\",\"inputs\":[{\"name\":\"_fromWalletId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"_toWalletId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"_tokenAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_receiver\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"unpause\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"walletNum\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"withdraw\",\"inputs\":[{\"name\":\"_walletId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"_tokenAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_receiver\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"ChangeOperator\",\"inputs\":[{\"name\":\"walletId\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"oldOperator\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOperator\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"CreateWallet\",\"inputs\":[{\"name\":\"walletId\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"owners\",\"type\":\"address[]\",\"indexed\":true,\"internalType\":\"address[]\"},{\"name\":\"operator\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"DepositToWallet\",\"inputs\":[{\"name\":\"walletId\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"tokenAddress\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"DrainToken\",\"inputs\":[{\"name\":\"tokenAddress\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"receiver\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Paused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ProposeNewOperator\",\"inputs\":[{\"name\":\"walletId\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"newOperator\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"proposer\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"TransferToWallet\",\"inputs\":[{\"name\":\"fromWalletId\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"toWalletId\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"tokenAddress\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"receiver\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Unpaused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"WithdrawFromWallet\",\"inputs\":[{\"name\":\"walletId\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"tokenAddress\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"receiver\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"EnforcedPause\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ExpectedPause\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OwnableInvalidOwner\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OwnableUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"SafeERC20FailedOperation\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"}]}]",
	Bin: "0x6080604052348015600e575f5ffd5b503380603357604051631e4fbdf760e01b81525f600482015260240160405180910390fd5b603a81603f565b506097565b5f80546001600160a01b03838116610100818102610100600160a81b0319851617855560405193049190911692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a35050565b611546806100a45f395ff3fe60806040526004361061011b575f3560e01c80638456cb591161009d578063bfa2c1d211610062578063bfa2c1d214610334578063c108bb4014610353578063cafd460014610372578063d68d9d4e14610391578063f2fde38b146103a4575f5ffd5b80638456cb591461028a5780638da5cb5b1461029e5780638e0cc176146102bf578063a0c89a8c146102de578063a96a5f94146102fd575f5ffd5b80633f4ba83a116100e35780633f4ba83a14610202578063530e931c146102165780635c975abb14610235578063715018a61461025757806380ba952e1461026b575f5ffd5b80630d63a1fd1461011f57806314da2906146101515780631687cc60146101a0578063323c4480146101cc57806336cc9e8d146101ed575b5f5ffd5b34801561012a575f5ffd5b5061013e61013936600461112d565b6103c3565b6040519081526020015b60405180910390f35b34801561015c575f5ffd5b5061018861016b366004611213565b5f908152600260205260409020600301546001600160a01b031690565b6040516001600160a01b039091168152602001610148565b3480156101ab575f5ffd5b506101bf6101ba366004611213565b61053a565b604051610148919061122a565b3480156101d7575f5ffd5b506101eb6101e6366004611275565b6105a3565b005b3480156101f8575f5ffd5b5061013e60015481565b34801561020d575f5ffd5b506101eb6106b2565b348015610221575f5ffd5b5061013e610230366004611275565b6106c4565b348015610240575f5ffd5b505f5460ff165b6040519015158152602001610148565b348015610262575f5ffd5b506101eb6106ef565b348015610276575f5ffd5b506101eb61028536600461129f565b610700565b348015610295575f5ffd5b506101eb610801565b3480156102a9575f5ffd5b505f5461010090046001600160a01b0316610188565b3480156102ca575f5ffd5b506101eb6102d93660046112e9565b610811565b3480156102e9575f5ffd5b506101eb6102f8366004611275565b6108ea565b348015610308575f5ffd5b50610188610317366004611213565b5f908152600260205260409020600101546001600160a01b031690565b34801561033f575f5ffd5b506101eb61034e36600461132a565b61093b565b34801561035e575f5ffd5b506101eb61036d366004611353565b6109a3565b34801561037d575f5ffd5b5061024761038c366004611275565b610a10565b6101eb61039f366004611213565b610a6a565b3480156103af575f5ffd5b506101eb6103be366004611375565b610aba565b5f6103cc610af7565b6001600160a01b0383166103fb5760405162461bcd60e51b81526004016103f290611395565b60405180910390fd5b6040516bffffffffffffffffffffffff1930606090811b8216602084015233901b166034820152604881018390525f9060680160408051601f1981840301815291815281516020928301205f818152600290935291206001810154919250906001600160a01b0316156104a55760405162461bcd60e51b815260206004820152601260248201527113d8d8dd5c1a5959081dd85b1b195d081a5960721b60448201526064016103f2565b85516104b79082906020890190611087565b50600181810180546001600160a01b0319166001600160a01b0388161790558054905f6104e3836113e0565b9190505550846001600160a01b03168660405161050091906113f8565b6040519081900381209084907fe778e91533ef049a5fc99752bc4efb2b50ca4c967dfc0d4bb4782fb128070c34905f90a450949350505050565b5f8181526002602090815260409182902080548351818402810184019094528084526060939283018282801561059757602002820191905f5260205f20905b81546001600160a01b03168152600190910190602001808311610579575b50505050509050919050565b81336105af8282610b1a565b6105cb5760405162461bcd60e51b81526004016103f290611436565b6001600160a01b0383166105f15760405162461bcd60e51b81526004016103f290611395565b5f84815260026020526040902060038101546001600160a01b0385811691161461063c5761061e81610b81565b6003810180546001600160a01b0319166001600160a01b0386161790555b335f818152600483016020526040808220805460ff19166001179055516001600160a01b0387169188917f71f9e7796b33cb192d1670169ee7f4af7c5364f8f01bab4b95466787593745c39190a461069381610be5565b156106ab576106a28585610c52565b6106ab81610b81565b5050505050565b6106ba610cde565b6106c2610d10565b565b5f8281526002602081815260408084206001600160a01b038616855290920190529020545b92915050565b6106f7610cde565b6106c25f610d61565b610708610af7565b5f8581526002602052604090206001015485906001600160a01b031633146107425760405162461bcd60e51b81526004016103f290611477565b858361074e8282610b1a565b61076a5760405162461bcd60e51b81526004016103f290611436565b86856107768282610b1a565b6107925760405162461bcd60e51b81526004016103f290611436565b61079f8a89886001610db9565b6107ab8989885f610db9565b604080516001600160a01b038981168252602082018990528a16918b918d917f1b56f805e5edb1e61b0d3f46feffdcbab5e591aa0e70e978ada9fc22093601c8910160405180910390a450505050505050505050565b610809610cde565b6106c2610e68565b610819610af7565b5f8481526002602052604090206001015484906001600160a01b031633146108535760405162461bcd60e51b81526004016103f290611477565b848361085f8282610b1a565b61087b5760405162461bcd60e51b81526004016103f290611436565b6108888787866001610db9565b846001600160a01b0316866001600160a01b0316887fd897e862036b62a0f770979fbd2227f3210565bba2eb4d9acd1dc8ccc00c928b876040516108ce91815260200190565b60405180910390a46108e1868686610ea4565b50505050505050565b6108f2610af7565b5f8281526002602052604090206001015482906001600160a01b0316331461092c5760405162461bcd60e51b81526004016103f290611477565b6109368383610c52565b505050565b610943610f61565b61094b610cde565b816001600160a01b0316836001600160a01b03167f896ecb17b26927fb33933fc5f413873193bced3c59fe736c42968a9778bf6b588360405161099091815260200190565b60405180910390a3610936838383610ea4565b6109ab610af7565b6109b78383835f610db9565b816001600160a01b0316837fbc8e388b96ba8b9f627cb6d72d3513182f763c33c6107ecd31191de1f71abc1a836040516109f391815260200190565b60405180910390a36109366001600160a01b038316333084610f83565b5f8282610a1d8282610b1a565b610a395760405162461bcd60e51b81526004016103f290611436565b5050505f9182526002602090815260408084206001600160a01b039390931684526004909201905290205460ff1690565b610a72610af7565b34610a7f825f8381610db9565b6040518181525f9083907fbc8e388b96ba8b9f627cb6d72d3513182f763c33c6107ecd31191de1f71abc1a9060200160405180910390a35050565b610ac2610cde565b6001600160a01b038116610aeb57604051631e4fbdf760e01b81525f60048201526024016103f2565b610af481610d61565b50565b5f5460ff16156106c25760405163d93c066560e01b815260040160405180910390fd5b5f828152600260205260408120815b8154811015610b7757815f018181548110610b4657610b466114ae565b5f918252602090912001546001600160a01b0390811690851603610b6f576001925050506106e9565b600101610b29565b505f949350505050565b5f5b8154811015610be1575f826004015f845f018481548110610ba657610ba66114ae565b5f918252602080832091909101546001600160a01b031683528201929092526040019020805460ff1916911515919091179055600101610b83565b5050565b5f805b8254811015610c4957826004015f845f018381548110610c0a57610c0a6114ae565b5f9182526020808320909101546001600160a01b0316835282019290925260400181205460ff1615159003610c4157505f92915050565b600101610be8565b50600192915050565b6001600160a01b038116610c785760405162461bcd60e51b81526004016103f290611395565b5f828152600260205260408082206001810180546001600160a01b038681166001600160a01b031983168117909355935192949316929091839187917f118c3f8030bc3c8254e737a0bd0584403c33646afbcbee8321c3bd5b26543cda9190a450505050565b5f546001600160a01b036101009091041633146106c25760405163118cdaa760e01b81523360048201526024016103f2565b610d18610f61565b5f805460ff191690557f5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa335b6040516001600160a01b03909116815260200160405180910390a1565b5f80546001600160a01b03838116610100818102610100600160a81b0319851617855560405193049190911692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a35050565b5f84815260026020526040812090826001811115610dd957610dd96114c2565b03610e22576001600160a01b0384165f908152600282016020526040902054610e039084906114d6565b6001600160a01b0385165f9081526002830160205260409020556106ab565b6001826001811115610e3657610e366114c2565b03610e60576001600160a01b0384165f908152600282016020526040902054610e039084906114e9565b6106ab6114fc565b610e70610af7565b5f805460ff191660011790557f62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258610d443390565b6001600160a01b038316610f4d575f826001600160a01b0316826040515f6040518083038185875af1925050503d805f8114610efb576040519150601f19603f3d011682016040523d82523d5f602084013e610f00565b606091505b5050905080610f475760405162461bcd60e51b8152602060048201526013602482015272115512081d1c985b9cd9995c8819985a5b1959606a1b60448201526064016103f2565b50505050565b6109366001600160a01b0384168383610fea565b5f5460ff166106c257604051638dfc202b60e01b815260040160405180910390fd5b6040516001600160a01b038481166024830152838116604483015260648201839052610f479186918216906323b872dd906084015b604051602081830303815290604052915060e01b6020820180516001600160e01b03838183161783525050505061101b565b6040516001600160a01b0383811660248301526044820183905261093691859182169063a9059cbb90606401610fb8565b5f5f60205f8451602086015f885af18061103a576040513d5f823e3d81fd5b50505f513d9150811561105157806001141561105e565b6001600160a01b0384163b155b15610f4757604051635274afe760e01b81526001600160a01b03851660048201526024016103f2565b828054828255905f5260205f209081019282156110da579160200282015b828111156110da57825182546001600160a01b0319166001600160a01b039091161782556020909201916001909101906110a5565b506110e69291506110ea565b5090565b5b808211156110e6575f81556001016110eb565b634e487b7160e01b5f52604160045260245ffd5b80356001600160a01b0381168114611128575f5ffd5b919050565b5f5f5f6060848603121561113f575f5ffd5b833567ffffffffffffffff811115611155575f5ffd5b8401601f81018613611165575f5ffd5b803567ffffffffffffffff81111561117f5761117f6110fe565b8060051b604051601f19603f830116810181811067ffffffffffffffff821117156111ac576111ac6110fe565b6040529182526020818401810192908101898411156111c9575f5ffd5b6020850194505b838510156111ef576111e185611112565b8152602094850194016111d0565b5095506112029250505060208501611112565b929592945050506040919091013590565b5f60208284031215611223575f5ffd5b5035919050565b602080825282518282018190525f918401906040840190835b8181101561126a5783516001600160a01b0316835260209384019390920191600101611243565b509095945050505050565b5f5f60408385031215611286575f5ffd5b8235915061129660208401611112565b90509250929050565b5f5f5f5f5f60a086880312156112b3575f5ffd5b85359450602086013593506112ca60408701611112565b92506112d860608701611112565b949793965091946080013592915050565b5f5f5f5f608085870312156112fc575f5ffd5b8435935061130c60208601611112565b925061131a60408601611112565b9396929550929360600135925050565b5f5f5f6060848603121561133c575f5ffd5b61134584611112565b925061120260208501611112565b5f5f5f60608486031215611365575f5ffd5b8335925061120260208501611112565b5f60208284031215611385575f5ffd5b61138e82611112565b9392505050565b6020808252601a908201527f4e6577206f70657261746f722069732061646472657373283029000000000000604082015260600190565b634e487b7160e01b5f52601160045260245ffd5b5f600182016113f1576113f16113cc565b5060010190565b81515f90829060208501835b8281101561142b5781516001600160a01b0316845260209384019390910190600101611404565b509195945050505050565b60208082526021908201527f476976656e2061646472657373206973206e6f742077616c6c6574206f776e656040820152603960f91b606082015260800190565b6020808252601a908201527f6d73672e73656e646572206973206e6f74206f70657261746f72000000000000604082015260600190565b634e487b7160e01b5f52603260045260245ffd5b634e487b7160e01b5f52602160045260245ffd5b808201808211156106e9576106e96113cc565b818103818111156106e9576106e96113cc565b634e487b7160e01b5f52600160045260245ffdfea26469706673582212201b9b4b6c183e79cec68d9f0ccf2b588af7aa700a4362844ad962a1fcf60f1f7d64736f6c634300081e0033",
}

// CelerWalletABI is the input ABI used to generate the binding from.
// Deprecated: Use CelerWalletMetaData.ABI instead.
var CelerWalletABI = CelerWalletMetaData.ABI

// CelerWalletBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use CelerWalletMetaData.Bin instead.
var CelerWalletBin = CelerWalletMetaData.Bin

// DeployCelerWallet deploys a new Ethereum contract, binding an instance of CelerWallet to it.
func DeployCelerWallet(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *CelerWallet, error) {
	parsed, err := CelerWalletMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(CelerWalletBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &CelerWallet{CelerWalletCaller: CelerWalletCaller{contract: contract}, CelerWalletTransactor: CelerWalletTransactor{contract: contract}, CelerWalletFilterer: CelerWalletFilterer{contract: contract}}, nil
}

// CelerWallet is an auto generated Go binding around an Ethereum contract.
type CelerWallet struct {
	CelerWalletCaller     // Read-only binding to the contract
	CelerWalletTransactor // Write-only binding to the contract
	CelerWalletFilterer   // Log filterer for contract events
}

// CelerWalletCaller is an auto generated read-only Go binding around an Ethereum contract.
type CelerWalletCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CelerWalletTransactor is an auto generated write-only Go binding around an Ethereum contract.
type CelerWalletTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CelerWalletFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type CelerWalletFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CelerWalletSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type CelerWalletSession struct {
	Contract     *CelerWallet      // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// CelerWalletCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type CelerWalletCallerSession struct {
	Contract *CelerWalletCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts      // Call options to use throughout this session
}

// CelerWalletTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type CelerWalletTransactorSession struct {
	Contract     *CelerWalletTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// CelerWalletRaw is an auto generated low-level Go binding around an Ethereum contract.
type CelerWalletRaw struct {
	Contract *CelerWallet // Generic contract binding to access the raw methods on
}

// CelerWalletCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type CelerWalletCallerRaw struct {
	Contract *CelerWalletCaller // Generic read-only contract binding to access the raw methods on
}

// CelerWalletTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type CelerWalletTransactorRaw struct {
	Contract *CelerWalletTransactor // Generic write-only contract binding to access the raw methods on
}

// NewCelerWallet creates a new instance of CelerWallet, bound to a specific deployed contract.
func NewCelerWallet(address common.Address, backend bind.ContractBackend) (*CelerWallet, error) {
	contract, err := bindCelerWallet(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &CelerWallet{CelerWalletCaller: CelerWalletCaller{contract: contract}, CelerWalletTransactor: CelerWalletTransactor{contract: contract}, CelerWalletFilterer: CelerWalletFilterer{contract: contract}}, nil
}

// NewCelerWalletCaller creates a new read-only instance of CelerWallet, bound to a specific deployed contract.
func NewCelerWalletCaller(address common.Address, caller bind.ContractCaller) (*CelerWalletCaller, error) {
	contract, err := bindCelerWallet(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &CelerWalletCaller{contract: contract}, nil
}

// NewCelerWalletTransactor creates a new write-only instance of CelerWallet, bound to a specific deployed contract.
func NewCelerWalletTransactor(address common.Address, transactor bind.ContractTransactor) (*CelerWalletTransactor, error) {
	contract, err := bindCelerWallet(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &CelerWalletTransactor{contract: contract}, nil
}

// NewCelerWalletFilterer creates a new log filterer instance of CelerWallet, bound to a specific deployed contract.
func NewCelerWalletFilterer(address common.Address, filterer bind.ContractFilterer) (*CelerWalletFilterer, error) {
	contract, err := bindCelerWallet(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &CelerWalletFilterer{contract: contract}, nil
}

// bindCelerWallet binds a generic wrapper to an already deployed contract.
func bindCelerWallet(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := CelerWalletMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_CelerWallet *CelerWalletRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _CelerWallet.Contract.CelerWalletCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_CelerWallet *CelerWalletRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CelerWallet.Contract.CelerWalletTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_CelerWallet *CelerWalletRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _CelerWallet.Contract.CelerWalletTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_CelerWallet *CelerWalletCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _CelerWallet.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_CelerWallet *CelerWalletTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CelerWallet.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_CelerWallet *CelerWalletTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _CelerWallet.Contract.contract.Transact(opts, method, params...)
}

// GetBalance is a free data retrieval call binding the contract method 0x530e931c.
//
// Solidity: function getBalance(bytes32 _walletId, address _tokenAddress) view returns(uint256)
func (_CelerWallet *CelerWalletCaller) GetBalance(opts *bind.CallOpts, _walletId [32]byte, _tokenAddress common.Address) (*big.Int, error) {
	var out []interface{}
	err := _CelerWallet.contract.Call(opts, &out, "getBalance", _walletId, _tokenAddress)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetBalance is a free data retrieval call binding the contract method 0x530e931c.
//
// Solidity: function getBalance(bytes32 _walletId, address _tokenAddress) view returns(uint256)
func (_CelerWallet *CelerWalletSession) GetBalance(_walletId [32]byte, _tokenAddress common.Address) (*big.Int, error) {
	return _CelerWallet.Contract.GetBalance(&_CelerWallet.CallOpts, _walletId, _tokenAddress)
}

// GetBalance is a free data retrieval call binding the contract method 0x530e931c.
//
// Solidity: function getBalance(bytes32 _walletId, address _tokenAddress) view returns(uint256)
func (_CelerWallet *CelerWalletCallerSession) GetBalance(_walletId [32]byte, _tokenAddress common.Address) (*big.Int, error) {
	return _CelerWallet.Contract.GetBalance(&_CelerWallet.CallOpts, _walletId, _tokenAddress)
}

// GetOperator is a free data retrieval call binding the contract method 0xa96a5f94.
//
// Solidity: function getOperator(bytes32 _walletId) view returns(address)
func (_CelerWallet *CelerWalletCaller) GetOperator(opts *bind.CallOpts, _walletId [32]byte) (common.Address, error) {
	var out []interface{}
	err := _CelerWallet.contract.Call(opts, &out, "getOperator", _walletId)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetOperator is a free data retrieval call binding the contract method 0xa96a5f94.
//
// Solidity: function getOperator(bytes32 _walletId) view returns(address)
func (_CelerWallet *CelerWalletSession) GetOperator(_walletId [32]byte) (common.Address, error) {
	return _CelerWallet.Contract.GetOperator(&_CelerWallet.CallOpts, _walletId)
}

// GetOperator is a free data retrieval call binding the contract method 0xa96a5f94.
//
// Solidity: function getOperator(bytes32 _walletId) view returns(address)
func (_CelerWallet *CelerWalletCallerSession) GetOperator(_walletId [32]byte) (common.Address, error) {
	return _CelerWallet.Contract.GetOperator(&_CelerWallet.CallOpts, _walletId)
}

// GetProposalVote is a free data retrieval call binding the contract method 0xcafd4600.
//
// Solidity: function getProposalVote(bytes32 _walletId, address _owner) view returns(bool)
func (_CelerWallet *CelerWalletCaller) GetProposalVote(opts *bind.CallOpts, _walletId [32]byte, _owner common.Address) (bool, error) {
	var out []interface{}
	err := _CelerWallet.contract.Call(opts, &out, "getProposalVote", _walletId, _owner)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// GetProposalVote is a free data retrieval call binding the contract method 0xcafd4600.
//
// Solidity: function getProposalVote(bytes32 _walletId, address _owner) view returns(bool)
func (_CelerWallet *CelerWalletSession) GetProposalVote(_walletId [32]byte, _owner common.Address) (bool, error) {
	return _CelerWallet.Contract.GetProposalVote(&_CelerWallet.CallOpts, _walletId, _owner)
}

// GetProposalVote is a free data retrieval call binding the contract method 0xcafd4600.
//
// Solidity: function getProposalVote(bytes32 _walletId, address _owner) view returns(bool)
func (_CelerWallet *CelerWalletCallerSession) GetProposalVote(_walletId [32]byte, _owner common.Address) (bool, error) {
	return _CelerWallet.Contract.GetProposalVote(&_CelerWallet.CallOpts, _walletId, _owner)
}

// GetProposedNewOperator is a free data retrieval call binding the contract method 0x14da2906.
//
// Solidity: function getProposedNewOperator(bytes32 _walletId) view returns(address)
func (_CelerWallet *CelerWalletCaller) GetProposedNewOperator(opts *bind.CallOpts, _walletId [32]byte) (common.Address, error) {
	var out []interface{}
	err := _CelerWallet.contract.Call(opts, &out, "getProposedNewOperator", _walletId)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetProposedNewOperator is a free data retrieval call binding the contract method 0x14da2906.
//
// Solidity: function getProposedNewOperator(bytes32 _walletId) view returns(address)
func (_CelerWallet *CelerWalletSession) GetProposedNewOperator(_walletId [32]byte) (common.Address, error) {
	return _CelerWallet.Contract.GetProposedNewOperator(&_CelerWallet.CallOpts, _walletId)
}

// GetProposedNewOperator is a free data retrieval call binding the contract method 0x14da2906.
//
// Solidity: function getProposedNewOperator(bytes32 _walletId) view returns(address)
func (_CelerWallet *CelerWalletCallerSession) GetProposedNewOperator(_walletId [32]byte) (common.Address, error) {
	return _CelerWallet.Contract.GetProposedNewOperator(&_CelerWallet.CallOpts, _walletId)
}

// GetWalletOwners is a free data retrieval call binding the contract method 0x1687cc60.
//
// Solidity: function getWalletOwners(bytes32 _walletId) view returns(address[])
func (_CelerWallet *CelerWalletCaller) GetWalletOwners(opts *bind.CallOpts, _walletId [32]byte) ([]common.Address, error) {
	var out []interface{}
	err := _CelerWallet.contract.Call(opts, &out, "getWalletOwners", _walletId)

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

// GetWalletOwners is a free data retrieval call binding the contract method 0x1687cc60.
//
// Solidity: function getWalletOwners(bytes32 _walletId) view returns(address[])
func (_CelerWallet *CelerWalletSession) GetWalletOwners(_walletId [32]byte) ([]common.Address, error) {
	return _CelerWallet.Contract.GetWalletOwners(&_CelerWallet.CallOpts, _walletId)
}

// GetWalletOwners is a free data retrieval call binding the contract method 0x1687cc60.
//
// Solidity: function getWalletOwners(bytes32 _walletId) view returns(address[])
func (_CelerWallet *CelerWalletCallerSession) GetWalletOwners(_walletId [32]byte) ([]common.Address, error) {
	return _CelerWallet.Contract.GetWalletOwners(&_CelerWallet.CallOpts, _walletId)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_CelerWallet *CelerWalletCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _CelerWallet.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_CelerWallet *CelerWalletSession) Owner() (common.Address, error) {
	return _CelerWallet.Contract.Owner(&_CelerWallet.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_CelerWallet *CelerWalletCallerSession) Owner() (common.Address, error) {
	return _CelerWallet.Contract.Owner(&_CelerWallet.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_CelerWallet *CelerWalletCaller) Paused(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _CelerWallet.contract.Call(opts, &out, "paused")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_CelerWallet *CelerWalletSession) Paused() (bool, error) {
	return _CelerWallet.Contract.Paused(&_CelerWallet.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_CelerWallet *CelerWalletCallerSession) Paused() (bool, error) {
	return _CelerWallet.Contract.Paused(&_CelerWallet.CallOpts)
}

// WalletNum is a free data retrieval call binding the contract method 0x36cc9e8d.
//
// Solidity: function walletNum() view returns(uint256)
func (_CelerWallet *CelerWalletCaller) WalletNum(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _CelerWallet.contract.Call(opts, &out, "walletNum")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// WalletNum is a free data retrieval call binding the contract method 0x36cc9e8d.
//
// Solidity: function walletNum() view returns(uint256)
func (_CelerWallet *CelerWalletSession) WalletNum() (*big.Int, error) {
	return _CelerWallet.Contract.WalletNum(&_CelerWallet.CallOpts)
}

// WalletNum is a free data retrieval call binding the contract method 0x36cc9e8d.
//
// Solidity: function walletNum() view returns(uint256)
func (_CelerWallet *CelerWalletCallerSession) WalletNum() (*big.Int, error) {
	return _CelerWallet.Contract.WalletNum(&_CelerWallet.CallOpts)
}

// Create is a paid mutator transaction binding the contract method 0x0d63a1fd.
//
// Solidity: function create(address[] _owners, address _operator, bytes32 _nonce) returns(bytes32)
func (_CelerWallet *CelerWalletTransactor) Create(opts *bind.TransactOpts, _owners []common.Address, _operator common.Address, _nonce [32]byte) (*types.Transaction, error) {
	return _CelerWallet.contract.Transact(opts, "create", _owners, _operator, _nonce)
}

// Create is a paid mutator transaction binding the contract method 0x0d63a1fd.
//
// Solidity: function create(address[] _owners, address _operator, bytes32 _nonce) returns(bytes32)
func (_CelerWallet *CelerWalletSession) Create(_owners []common.Address, _operator common.Address, _nonce [32]byte) (*types.Transaction, error) {
	return _CelerWallet.Contract.Create(&_CelerWallet.TransactOpts, _owners, _operator, _nonce)
}

// Create is a paid mutator transaction binding the contract method 0x0d63a1fd.
//
// Solidity: function create(address[] _owners, address _operator, bytes32 _nonce) returns(bytes32)
func (_CelerWallet *CelerWalletTransactorSession) Create(_owners []common.Address, _operator common.Address, _nonce [32]byte) (*types.Transaction, error) {
	return _CelerWallet.Contract.Create(&_CelerWallet.TransactOpts, _owners, _operator, _nonce)
}

// DepositERC20 is a paid mutator transaction binding the contract method 0xc108bb40.
//
// Solidity: function depositERC20(bytes32 _walletId, address _tokenAddress, uint256 _amount) returns()
func (_CelerWallet *CelerWalletTransactor) DepositERC20(opts *bind.TransactOpts, _walletId [32]byte, _tokenAddress common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _CelerWallet.contract.Transact(opts, "depositERC20", _walletId, _tokenAddress, _amount)
}

// DepositERC20 is a paid mutator transaction binding the contract method 0xc108bb40.
//
// Solidity: function depositERC20(bytes32 _walletId, address _tokenAddress, uint256 _amount) returns()
func (_CelerWallet *CelerWalletSession) DepositERC20(_walletId [32]byte, _tokenAddress common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _CelerWallet.Contract.DepositERC20(&_CelerWallet.TransactOpts, _walletId, _tokenAddress, _amount)
}

// DepositERC20 is a paid mutator transaction binding the contract method 0xc108bb40.
//
// Solidity: function depositERC20(bytes32 _walletId, address _tokenAddress, uint256 _amount) returns()
func (_CelerWallet *CelerWalletTransactorSession) DepositERC20(_walletId [32]byte, _tokenAddress common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _CelerWallet.Contract.DepositERC20(&_CelerWallet.TransactOpts, _walletId, _tokenAddress, _amount)
}

// DepositETH is a paid mutator transaction binding the contract method 0xd68d9d4e.
//
// Solidity: function depositETH(bytes32 _walletId) payable returns()
func (_CelerWallet *CelerWalletTransactor) DepositETH(opts *bind.TransactOpts, _walletId [32]byte) (*types.Transaction, error) {
	return _CelerWallet.contract.Transact(opts, "depositETH", _walletId)
}

// DepositETH is a paid mutator transaction binding the contract method 0xd68d9d4e.
//
// Solidity: function depositETH(bytes32 _walletId) payable returns()
func (_CelerWallet *CelerWalletSession) DepositETH(_walletId [32]byte) (*types.Transaction, error) {
	return _CelerWallet.Contract.DepositETH(&_CelerWallet.TransactOpts, _walletId)
}

// DepositETH is a paid mutator transaction binding the contract method 0xd68d9d4e.
//
// Solidity: function depositETH(bytes32 _walletId) payable returns()
func (_CelerWallet *CelerWalletTransactorSession) DepositETH(_walletId [32]byte) (*types.Transaction, error) {
	return _CelerWallet.Contract.DepositETH(&_CelerWallet.TransactOpts, _walletId)
}

// DrainToken is a paid mutator transaction binding the contract method 0xbfa2c1d2.
//
// Solidity: function drainToken(address _tokenAddress, address _receiver, uint256 _amount) returns()
func (_CelerWallet *CelerWalletTransactor) DrainToken(opts *bind.TransactOpts, _tokenAddress common.Address, _receiver common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _CelerWallet.contract.Transact(opts, "drainToken", _tokenAddress, _receiver, _amount)
}

// DrainToken is a paid mutator transaction binding the contract method 0xbfa2c1d2.
//
// Solidity: function drainToken(address _tokenAddress, address _receiver, uint256 _amount) returns()
func (_CelerWallet *CelerWalletSession) DrainToken(_tokenAddress common.Address, _receiver common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _CelerWallet.Contract.DrainToken(&_CelerWallet.TransactOpts, _tokenAddress, _receiver, _amount)
}

// DrainToken is a paid mutator transaction binding the contract method 0xbfa2c1d2.
//
// Solidity: function drainToken(address _tokenAddress, address _receiver, uint256 _amount) returns()
func (_CelerWallet *CelerWalletTransactorSession) DrainToken(_tokenAddress common.Address, _receiver common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _CelerWallet.Contract.DrainToken(&_CelerWallet.TransactOpts, _tokenAddress, _receiver, _amount)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_CelerWallet *CelerWalletTransactor) Pause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CelerWallet.contract.Transact(opts, "pause")
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_CelerWallet *CelerWalletSession) Pause() (*types.Transaction, error) {
	return _CelerWallet.Contract.Pause(&_CelerWallet.TransactOpts)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_CelerWallet *CelerWalletTransactorSession) Pause() (*types.Transaction, error) {
	return _CelerWallet.Contract.Pause(&_CelerWallet.TransactOpts)
}

// ProposeNewOperator is a paid mutator transaction binding the contract method 0x323c4480.
//
// Solidity: function proposeNewOperator(bytes32 _walletId, address _newOperator) returns()
func (_CelerWallet *CelerWalletTransactor) ProposeNewOperator(opts *bind.TransactOpts, _walletId [32]byte, _newOperator common.Address) (*types.Transaction, error) {
	return _CelerWallet.contract.Transact(opts, "proposeNewOperator", _walletId, _newOperator)
}

// ProposeNewOperator is a paid mutator transaction binding the contract method 0x323c4480.
//
// Solidity: function proposeNewOperator(bytes32 _walletId, address _newOperator) returns()
func (_CelerWallet *CelerWalletSession) ProposeNewOperator(_walletId [32]byte, _newOperator common.Address) (*types.Transaction, error) {
	return _CelerWallet.Contract.ProposeNewOperator(&_CelerWallet.TransactOpts, _walletId, _newOperator)
}

// ProposeNewOperator is a paid mutator transaction binding the contract method 0x323c4480.
//
// Solidity: function proposeNewOperator(bytes32 _walletId, address _newOperator) returns()
func (_CelerWallet *CelerWalletTransactorSession) ProposeNewOperator(_walletId [32]byte, _newOperator common.Address) (*types.Transaction, error) {
	return _CelerWallet.Contract.ProposeNewOperator(&_CelerWallet.TransactOpts, _walletId, _newOperator)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_CelerWallet *CelerWalletTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CelerWallet.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_CelerWallet *CelerWalletSession) RenounceOwnership() (*types.Transaction, error) {
	return _CelerWallet.Contract.RenounceOwnership(&_CelerWallet.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_CelerWallet *CelerWalletTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _CelerWallet.Contract.RenounceOwnership(&_CelerWallet.TransactOpts)
}

// TransferOperatorship is a paid mutator transaction binding the contract method 0xa0c89a8c.
//
// Solidity: function transferOperatorship(bytes32 _walletId, address _newOperator) returns()
func (_CelerWallet *CelerWalletTransactor) TransferOperatorship(opts *bind.TransactOpts, _walletId [32]byte, _newOperator common.Address) (*types.Transaction, error) {
	return _CelerWallet.contract.Transact(opts, "transferOperatorship", _walletId, _newOperator)
}

// TransferOperatorship is a paid mutator transaction binding the contract method 0xa0c89a8c.
//
// Solidity: function transferOperatorship(bytes32 _walletId, address _newOperator) returns()
func (_CelerWallet *CelerWalletSession) TransferOperatorship(_walletId [32]byte, _newOperator common.Address) (*types.Transaction, error) {
	return _CelerWallet.Contract.TransferOperatorship(&_CelerWallet.TransactOpts, _walletId, _newOperator)
}

// TransferOperatorship is a paid mutator transaction binding the contract method 0xa0c89a8c.
//
// Solidity: function transferOperatorship(bytes32 _walletId, address _newOperator) returns()
func (_CelerWallet *CelerWalletTransactorSession) TransferOperatorship(_walletId [32]byte, _newOperator common.Address) (*types.Transaction, error) {
	return _CelerWallet.Contract.TransferOperatorship(&_CelerWallet.TransactOpts, _walletId, _newOperator)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_CelerWallet *CelerWalletTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _CelerWallet.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_CelerWallet *CelerWalletSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _CelerWallet.Contract.TransferOwnership(&_CelerWallet.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_CelerWallet *CelerWalletTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _CelerWallet.Contract.TransferOwnership(&_CelerWallet.TransactOpts, newOwner)
}

// TransferToWallet is a paid mutator transaction binding the contract method 0x80ba952e.
//
// Solidity: function transferToWallet(bytes32 _fromWalletId, bytes32 _toWalletId, address _tokenAddress, address _receiver, uint256 _amount) returns()
func (_CelerWallet *CelerWalletTransactor) TransferToWallet(opts *bind.TransactOpts, _fromWalletId [32]byte, _toWalletId [32]byte, _tokenAddress common.Address, _receiver common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _CelerWallet.contract.Transact(opts, "transferToWallet", _fromWalletId, _toWalletId, _tokenAddress, _receiver, _amount)
}

// TransferToWallet is a paid mutator transaction binding the contract method 0x80ba952e.
//
// Solidity: function transferToWallet(bytes32 _fromWalletId, bytes32 _toWalletId, address _tokenAddress, address _receiver, uint256 _amount) returns()
func (_CelerWallet *CelerWalletSession) TransferToWallet(_fromWalletId [32]byte, _toWalletId [32]byte, _tokenAddress common.Address, _receiver common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _CelerWallet.Contract.TransferToWallet(&_CelerWallet.TransactOpts, _fromWalletId, _toWalletId, _tokenAddress, _receiver, _amount)
}

// TransferToWallet is a paid mutator transaction binding the contract method 0x80ba952e.
//
// Solidity: function transferToWallet(bytes32 _fromWalletId, bytes32 _toWalletId, address _tokenAddress, address _receiver, uint256 _amount) returns()
func (_CelerWallet *CelerWalletTransactorSession) TransferToWallet(_fromWalletId [32]byte, _toWalletId [32]byte, _tokenAddress common.Address, _receiver common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _CelerWallet.Contract.TransferToWallet(&_CelerWallet.TransactOpts, _fromWalletId, _toWalletId, _tokenAddress, _receiver, _amount)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_CelerWallet *CelerWalletTransactor) Unpause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CelerWallet.contract.Transact(opts, "unpause")
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_CelerWallet *CelerWalletSession) Unpause() (*types.Transaction, error) {
	return _CelerWallet.Contract.Unpause(&_CelerWallet.TransactOpts)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_CelerWallet *CelerWalletTransactorSession) Unpause() (*types.Transaction, error) {
	return _CelerWallet.Contract.Unpause(&_CelerWallet.TransactOpts)
}

// Withdraw is a paid mutator transaction binding the contract method 0x8e0cc176.
//
// Solidity: function withdraw(bytes32 _walletId, address _tokenAddress, address _receiver, uint256 _amount) returns()
func (_CelerWallet *CelerWalletTransactor) Withdraw(opts *bind.TransactOpts, _walletId [32]byte, _tokenAddress common.Address, _receiver common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _CelerWallet.contract.Transact(opts, "withdraw", _walletId, _tokenAddress, _receiver, _amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0x8e0cc176.
//
// Solidity: function withdraw(bytes32 _walletId, address _tokenAddress, address _receiver, uint256 _amount) returns()
func (_CelerWallet *CelerWalletSession) Withdraw(_walletId [32]byte, _tokenAddress common.Address, _receiver common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _CelerWallet.Contract.Withdraw(&_CelerWallet.TransactOpts, _walletId, _tokenAddress, _receiver, _amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0x8e0cc176.
//
// Solidity: function withdraw(bytes32 _walletId, address _tokenAddress, address _receiver, uint256 _amount) returns()
func (_CelerWallet *CelerWalletTransactorSession) Withdraw(_walletId [32]byte, _tokenAddress common.Address, _receiver common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _CelerWallet.Contract.Withdraw(&_CelerWallet.TransactOpts, _walletId, _tokenAddress, _receiver, _amount)
}

// CelerWalletChangeOperatorIterator is returned from FilterChangeOperator and is used to iterate over the raw logs and unpacked data for ChangeOperator events raised by the CelerWallet contract.
type CelerWalletChangeOperatorIterator struct {
	Event *CelerWalletChangeOperator // Event containing the contract specifics and raw log

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
func (it *CelerWalletChangeOperatorIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CelerWalletChangeOperator)
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
		it.Event = new(CelerWalletChangeOperator)
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
func (it *CelerWalletChangeOperatorIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CelerWalletChangeOperatorIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CelerWalletChangeOperator represents a ChangeOperator event raised by the CelerWallet contract.
type CelerWalletChangeOperator struct {
	WalletId    [32]byte
	OldOperator common.Address
	NewOperator common.Address
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterChangeOperator is a free log retrieval operation binding the contract event 0x118c3f8030bc3c8254e737a0bd0584403c33646afbcbee8321c3bd5b26543cda.
//
// Solidity: event ChangeOperator(bytes32 indexed walletId, address indexed oldOperator, address indexed newOperator)
func (_CelerWallet *CelerWalletFilterer) FilterChangeOperator(opts *bind.FilterOpts, walletId [][32]byte, oldOperator []common.Address, newOperator []common.Address) (*CelerWalletChangeOperatorIterator, error) {

	var walletIdRule []interface{}
	for _, walletIdItem := range walletId {
		walletIdRule = append(walletIdRule, walletIdItem)
	}
	var oldOperatorRule []interface{}
	for _, oldOperatorItem := range oldOperator {
		oldOperatorRule = append(oldOperatorRule, oldOperatorItem)
	}
	var newOperatorRule []interface{}
	for _, newOperatorItem := range newOperator {
		newOperatorRule = append(newOperatorRule, newOperatorItem)
	}

	logs, sub, err := _CelerWallet.contract.FilterLogs(opts, "ChangeOperator", walletIdRule, oldOperatorRule, newOperatorRule)
	if err != nil {
		return nil, err
	}
	return &CelerWalletChangeOperatorIterator{contract: _CelerWallet.contract, event: "ChangeOperator", logs: logs, sub: sub}, nil
}

// WatchChangeOperator is a free log subscription operation binding the contract event 0x118c3f8030bc3c8254e737a0bd0584403c33646afbcbee8321c3bd5b26543cda.
//
// Solidity: event ChangeOperator(bytes32 indexed walletId, address indexed oldOperator, address indexed newOperator)
func (_CelerWallet *CelerWalletFilterer) WatchChangeOperator(opts *bind.WatchOpts, sink chan<- *CelerWalletChangeOperator, walletId [][32]byte, oldOperator []common.Address, newOperator []common.Address) (event.Subscription, error) {

	var walletIdRule []interface{}
	for _, walletIdItem := range walletId {
		walletIdRule = append(walletIdRule, walletIdItem)
	}
	var oldOperatorRule []interface{}
	for _, oldOperatorItem := range oldOperator {
		oldOperatorRule = append(oldOperatorRule, oldOperatorItem)
	}
	var newOperatorRule []interface{}
	for _, newOperatorItem := range newOperator {
		newOperatorRule = append(newOperatorRule, newOperatorItem)
	}

	logs, sub, err := _CelerWallet.contract.WatchLogs(opts, "ChangeOperator", walletIdRule, oldOperatorRule, newOperatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CelerWalletChangeOperator)
				if err := _CelerWallet.contract.UnpackLog(event, "ChangeOperator", log); err != nil {
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

// ParseChangeOperator is a log parse operation binding the contract event 0x118c3f8030bc3c8254e737a0bd0584403c33646afbcbee8321c3bd5b26543cda.
//
// Solidity: event ChangeOperator(bytes32 indexed walletId, address indexed oldOperator, address indexed newOperator)
func (_CelerWallet *CelerWalletFilterer) ParseChangeOperator(log types.Log) (*CelerWalletChangeOperator, error) {
	event := new(CelerWalletChangeOperator)
	if err := _CelerWallet.contract.UnpackLog(event, "ChangeOperator", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CelerWalletCreateWalletIterator is returned from FilterCreateWallet and is used to iterate over the raw logs and unpacked data for CreateWallet events raised by the CelerWallet contract.
type CelerWalletCreateWalletIterator struct {
	Event *CelerWalletCreateWallet // Event containing the contract specifics and raw log

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
func (it *CelerWalletCreateWalletIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CelerWalletCreateWallet)
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
		it.Event = new(CelerWalletCreateWallet)
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
func (it *CelerWalletCreateWalletIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CelerWalletCreateWalletIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CelerWalletCreateWallet represents a CreateWallet event raised by the CelerWallet contract.
type CelerWalletCreateWallet struct {
	WalletId [32]byte
	Owners   []common.Address
	Operator common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterCreateWallet is a free log retrieval operation binding the contract event 0xe778e91533ef049a5fc99752bc4efb2b50ca4c967dfc0d4bb4782fb128070c34.
//
// Solidity: event CreateWallet(bytes32 indexed walletId, address[] indexed owners, address indexed operator)
func (_CelerWallet *CelerWalletFilterer) FilterCreateWallet(opts *bind.FilterOpts, walletId [][32]byte, owners [][]common.Address, operator []common.Address) (*CelerWalletCreateWalletIterator, error) {

	var walletIdRule []interface{}
	for _, walletIdItem := range walletId {
		walletIdRule = append(walletIdRule, walletIdItem)
	}
	var ownersRule []interface{}
	for _, ownersItem := range owners {
		ownersRule = append(ownersRule, ownersItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _CelerWallet.contract.FilterLogs(opts, "CreateWallet", walletIdRule, ownersRule, operatorRule)
	if err != nil {
		return nil, err
	}
	return &CelerWalletCreateWalletIterator{contract: _CelerWallet.contract, event: "CreateWallet", logs: logs, sub: sub}, nil
}

// WatchCreateWallet is a free log subscription operation binding the contract event 0xe778e91533ef049a5fc99752bc4efb2b50ca4c967dfc0d4bb4782fb128070c34.
//
// Solidity: event CreateWallet(bytes32 indexed walletId, address[] indexed owners, address indexed operator)
func (_CelerWallet *CelerWalletFilterer) WatchCreateWallet(opts *bind.WatchOpts, sink chan<- *CelerWalletCreateWallet, walletId [][32]byte, owners [][]common.Address, operator []common.Address) (event.Subscription, error) {

	var walletIdRule []interface{}
	for _, walletIdItem := range walletId {
		walletIdRule = append(walletIdRule, walletIdItem)
	}
	var ownersRule []interface{}
	for _, ownersItem := range owners {
		ownersRule = append(ownersRule, ownersItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _CelerWallet.contract.WatchLogs(opts, "CreateWallet", walletIdRule, ownersRule, operatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CelerWalletCreateWallet)
				if err := _CelerWallet.contract.UnpackLog(event, "CreateWallet", log); err != nil {
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

// ParseCreateWallet is a log parse operation binding the contract event 0xe778e91533ef049a5fc99752bc4efb2b50ca4c967dfc0d4bb4782fb128070c34.
//
// Solidity: event CreateWallet(bytes32 indexed walletId, address[] indexed owners, address indexed operator)
func (_CelerWallet *CelerWalletFilterer) ParseCreateWallet(log types.Log) (*CelerWalletCreateWallet, error) {
	event := new(CelerWalletCreateWallet)
	if err := _CelerWallet.contract.UnpackLog(event, "CreateWallet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CelerWalletDepositToWalletIterator is returned from FilterDepositToWallet and is used to iterate over the raw logs and unpacked data for DepositToWallet events raised by the CelerWallet contract.
type CelerWalletDepositToWalletIterator struct {
	Event *CelerWalletDepositToWallet // Event containing the contract specifics and raw log

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
func (it *CelerWalletDepositToWalletIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CelerWalletDepositToWallet)
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
		it.Event = new(CelerWalletDepositToWallet)
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
func (it *CelerWalletDepositToWalletIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CelerWalletDepositToWalletIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CelerWalletDepositToWallet represents a DepositToWallet event raised by the CelerWallet contract.
type CelerWalletDepositToWallet struct {
	WalletId     [32]byte
	TokenAddress common.Address
	Amount       *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterDepositToWallet is a free log retrieval operation binding the contract event 0xbc8e388b96ba8b9f627cb6d72d3513182f763c33c6107ecd31191de1f71abc1a.
//
// Solidity: event DepositToWallet(bytes32 indexed walletId, address indexed tokenAddress, uint256 amount)
func (_CelerWallet *CelerWalletFilterer) FilterDepositToWallet(opts *bind.FilterOpts, walletId [][32]byte, tokenAddress []common.Address) (*CelerWalletDepositToWalletIterator, error) {

	var walletIdRule []interface{}
	for _, walletIdItem := range walletId {
		walletIdRule = append(walletIdRule, walletIdItem)
	}
	var tokenAddressRule []interface{}
	for _, tokenAddressItem := range tokenAddress {
		tokenAddressRule = append(tokenAddressRule, tokenAddressItem)
	}

	logs, sub, err := _CelerWallet.contract.FilterLogs(opts, "DepositToWallet", walletIdRule, tokenAddressRule)
	if err != nil {
		return nil, err
	}
	return &CelerWalletDepositToWalletIterator{contract: _CelerWallet.contract, event: "DepositToWallet", logs: logs, sub: sub}, nil
}

// WatchDepositToWallet is a free log subscription operation binding the contract event 0xbc8e388b96ba8b9f627cb6d72d3513182f763c33c6107ecd31191de1f71abc1a.
//
// Solidity: event DepositToWallet(bytes32 indexed walletId, address indexed tokenAddress, uint256 amount)
func (_CelerWallet *CelerWalletFilterer) WatchDepositToWallet(opts *bind.WatchOpts, sink chan<- *CelerWalletDepositToWallet, walletId [][32]byte, tokenAddress []common.Address) (event.Subscription, error) {

	var walletIdRule []interface{}
	for _, walletIdItem := range walletId {
		walletIdRule = append(walletIdRule, walletIdItem)
	}
	var tokenAddressRule []interface{}
	for _, tokenAddressItem := range tokenAddress {
		tokenAddressRule = append(tokenAddressRule, tokenAddressItem)
	}

	logs, sub, err := _CelerWallet.contract.WatchLogs(opts, "DepositToWallet", walletIdRule, tokenAddressRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CelerWalletDepositToWallet)
				if err := _CelerWallet.contract.UnpackLog(event, "DepositToWallet", log); err != nil {
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

// ParseDepositToWallet is a log parse operation binding the contract event 0xbc8e388b96ba8b9f627cb6d72d3513182f763c33c6107ecd31191de1f71abc1a.
//
// Solidity: event DepositToWallet(bytes32 indexed walletId, address indexed tokenAddress, uint256 amount)
func (_CelerWallet *CelerWalletFilterer) ParseDepositToWallet(log types.Log) (*CelerWalletDepositToWallet, error) {
	event := new(CelerWalletDepositToWallet)
	if err := _CelerWallet.contract.UnpackLog(event, "DepositToWallet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CelerWalletDrainTokenIterator is returned from FilterDrainToken and is used to iterate over the raw logs and unpacked data for DrainToken events raised by the CelerWallet contract.
type CelerWalletDrainTokenIterator struct {
	Event *CelerWalletDrainToken // Event containing the contract specifics and raw log

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
func (it *CelerWalletDrainTokenIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CelerWalletDrainToken)
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
		it.Event = new(CelerWalletDrainToken)
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
func (it *CelerWalletDrainTokenIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CelerWalletDrainTokenIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CelerWalletDrainToken represents a DrainToken event raised by the CelerWallet contract.
type CelerWalletDrainToken struct {
	TokenAddress common.Address
	Receiver     common.Address
	Amount       *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterDrainToken is a free log retrieval operation binding the contract event 0x896ecb17b26927fb33933fc5f413873193bced3c59fe736c42968a9778bf6b58.
//
// Solidity: event DrainToken(address indexed tokenAddress, address indexed receiver, uint256 amount)
func (_CelerWallet *CelerWalletFilterer) FilterDrainToken(opts *bind.FilterOpts, tokenAddress []common.Address, receiver []common.Address) (*CelerWalletDrainTokenIterator, error) {

	var tokenAddressRule []interface{}
	for _, tokenAddressItem := range tokenAddress {
		tokenAddressRule = append(tokenAddressRule, tokenAddressItem)
	}
	var receiverRule []interface{}
	for _, receiverItem := range receiver {
		receiverRule = append(receiverRule, receiverItem)
	}

	logs, sub, err := _CelerWallet.contract.FilterLogs(opts, "DrainToken", tokenAddressRule, receiverRule)
	if err != nil {
		return nil, err
	}
	return &CelerWalletDrainTokenIterator{contract: _CelerWallet.contract, event: "DrainToken", logs: logs, sub: sub}, nil
}

// WatchDrainToken is a free log subscription operation binding the contract event 0x896ecb17b26927fb33933fc5f413873193bced3c59fe736c42968a9778bf6b58.
//
// Solidity: event DrainToken(address indexed tokenAddress, address indexed receiver, uint256 amount)
func (_CelerWallet *CelerWalletFilterer) WatchDrainToken(opts *bind.WatchOpts, sink chan<- *CelerWalletDrainToken, tokenAddress []common.Address, receiver []common.Address) (event.Subscription, error) {

	var tokenAddressRule []interface{}
	for _, tokenAddressItem := range tokenAddress {
		tokenAddressRule = append(tokenAddressRule, tokenAddressItem)
	}
	var receiverRule []interface{}
	for _, receiverItem := range receiver {
		receiverRule = append(receiverRule, receiverItem)
	}

	logs, sub, err := _CelerWallet.contract.WatchLogs(opts, "DrainToken", tokenAddressRule, receiverRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CelerWalletDrainToken)
				if err := _CelerWallet.contract.UnpackLog(event, "DrainToken", log); err != nil {
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

// ParseDrainToken is a log parse operation binding the contract event 0x896ecb17b26927fb33933fc5f413873193bced3c59fe736c42968a9778bf6b58.
//
// Solidity: event DrainToken(address indexed tokenAddress, address indexed receiver, uint256 amount)
func (_CelerWallet *CelerWalletFilterer) ParseDrainToken(log types.Log) (*CelerWalletDrainToken, error) {
	event := new(CelerWalletDrainToken)
	if err := _CelerWallet.contract.UnpackLog(event, "DrainToken", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CelerWalletOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the CelerWallet contract.
type CelerWalletOwnershipTransferredIterator struct {
	Event *CelerWalletOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *CelerWalletOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CelerWalletOwnershipTransferred)
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
		it.Event = new(CelerWalletOwnershipTransferred)
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
func (it *CelerWalletOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CelerWalletOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CelerWalletOwnershipTransferred represents a OwnershipTransferred event raised by the CelerWallet contract.
type CelerWalletOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_CelerWallet *CelerWalletFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*CelerWalletOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _CelerWallet.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &CelerWalletOwnershipTransferredIterator{contract: _CelerWallet.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_CelerWallet *CelerWalletFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *CelerWalletOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _CelerWallet.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CelerWalletOwnershipTransferred)
				if err := _CelerWallet.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_CelerWallet *CelerWalletFilterer) ParseOwnershipTransferred(log types.Log) (*CelerWalletOwnershipTransferred, error) {
	event := new(CelerWalletOwnershipTransferred)
	if err := _CelerWallet.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CelerWalletPausedIterator is returned from FilterPaused and is used to iterate over the raw logs and unpacked data for Paused events raised by the CelerWallet contract.
type CelerWalletPausedIterator struct {
	Event *CelerWalletPaused // Event containing the contract specifics and raw log

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
func (it *CelerWalletPausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CelerWalletPaused)
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
		it.Event = new(CelerWalletPaused)
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
func (it *CelerWalletPausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CelerWalletPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CelerWalletPaused represents a Paused event raised by the CelerWallet contract.
type CelerWalletPaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterPaused is a free log retrieval operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_CelerWallet *CelerWalletFilterer) FilterPaused(opts *bind.FilterOpts) (*CelerWalletPausedIterator, error) {

	logs, sub, err := _CelerWallet.contract.FilterLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return &CelerWalletPausedIterator{contract: _CelerWallet.contract, event: "Paused", logs: logs, sub: sub}, nil
}

// WatchPaused is a free log subscription operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_CelerWallet *CelerWalletFilterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *CelerWalletPaused) (event.Subscription, error) {

	logs, sub, err := _CelerWallet.contract.WatchLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CelerWalletPaused)
				if err := _CelerWallet.contract.UnpackLog(event, "Paused", log); err != nil {
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

// ParsePaused is a log parse operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_CelerWallet *CelerWalletFilterer) ParsePaused(log types.Log) (*CelerWalletPaused, error) {
	event := new(CelerWalletPaused)
	if err := _CelerWallet.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CelerWalletProposeNewOperatorIterator is returned from FilterProposeNewOperator and is used to iterate over the raw logs and unpacked data for ProposeNewOperator events raised by the CelerWallet contract.
type CelerWalletProposeNewOperatorIterator struct {
	Event *CelerWalletProposeNewOperator // Event containing the contract specifics and raw log

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
func (it *CelerWalletProposeNewOperatorIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CelerWalletProposeNewOperator)
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
		it.Event = new(CelerWalletProposeNewOperator)
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
func (it *CelerWalletProposeNewOperatorIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CelerWalletProposeNewOperatorIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CelerWalletProposeNewOperator represents a ProposeNewOperator event raised by the CelerWallet contract.
type CelerWalletProposeNewOperator struct {
	WalletId    [32]byte
	NewOperator common.Address
	Proposer    common.Address
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterProposeNewOperator is a free log retrieval operation binding the contract event 0x71f9e7796b33cb192d1670169ee7f4af7c5364f8f01bab4b95466787593745c3.
//
// Solidity: event ProposeNewOperator(bytes32 indexed walletId, address indexed newOperator, address indexed proposer)
func (_CelerWallet *CelerWalletFilterer) FilterProposeNewOperator(opts *bind.FilterOpts, walletId [][32]byte, newOperator []common.Address, proposer []common.Address) (*CelerWalletProposeNewOperatorIterator, error) {

	var walletIdRule []interface{}
	for _, walletIdItem := range walletId {
		walletIdRule = append(walletIdRule, walletIdItem)
	}
	var newOperatorRule []interface{}
	for _, newOperatorItem := range newOperator {
		newOperatorRule = append(newOperatorRule, newOperatorItem)
	}
	var proposerRule []interface{}
	for _, proposerItem := range proposer {
		proposerRule = append(proposerRule, proposerItem)
	}

	logs, sub, err := _CelerWallet.contract.FilterLogs(opts, "ProposeNewOperator", walletIdRule, newOperatorRule, proposerRule)
	if err != nil {
		return nil, err
	}
	return &CelerWalletProposeNewOperatorIterator{contract: _CelerWallet.contract, event: "ProposeNewOperator", logs: logs, sub: sub}, nil
}

// WatchProposeNewOperator is a free log subscription operation binding the contract event 0x71f9e7796b33cb192d1670169ee7f4af7c5364f8f01bab4b95466787593745c3.
//
// Solidity: event ProposeNewOperator(bytes32 indexed walletId, address indexed newOperator, address indexed proposer)
func (_CelerWallet *CelerWalletFilterer) WatchProposeNewOperator(opts *bind.WatchOpts, sink chan<- *CelerWalletProposeNewOperator, walletId [][32]byte, newOperator []common.Address, proposer []common.Address) (event.Subscription, error) {

	var walletIdRule []interface{}
	for _, walletIdItem := range walletId {
		walletIdRule = append(walletIdRule, walletIdItem)
	}
	var newOperatorRule []interface{}
	for _, newOperatorItem := range newOperator {
		newOperatorRule = append(newOperatorRule, newOperatorItem)
	}
	var proposerRule []interface{}
	for _, proposerItem := range proposer {
		proposerRule = append(proposerRule, proposerItem)
	}

	logs, sub, err := _CelerWallet.contract.WatchLogs(opts, "ProposeNewOperator", walletIdRule, newOperatorRule, proposerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CelerWalletProposeNewOperator)
				if err := _CelerWallet.contract.UnpackLog(event, "ProposeNewOperator", log); err != nil {
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

// ParseProposeNewOperator is a log parse operation binding the contract event 0x71f9e7796b33cb192d1670169ee7f4af7c5364f8f01bab4b95466787593745c3.
//
// Solidity: event ProposeNewOperator(bytes32 indexed walletId, address indexed newOperator, address indexed proposer)
func (_CelerWallet *CelerWalletFilterer) ParseProposeNewOperator(log types.Log) (*CelerWalletProposeNewOperator, error) {
	event := new(CelerWalletProposeNewOperator)
	if err := _CelerWallet.contract.UnpackLog(event, "ProposeNewOperator", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CelerWalletTransferToWalletIterator is returned from FilterTransferToWallet and is used to iterate over the raw logs and unpacked data for TransferToWallet events raised by the CelerWallet contract.
type CelerWalletTransferToWalletIterator struct {
	Event *CelerWalletTransferToWallet // Event containing the contract specifics and raw log

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
func (it *CelerWalletTransferToWalletIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CelerWalletTransferToWallet)
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
		it.Event = new(CelerWalletTransferToWallet)
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
func (it *CelerWalletTransferToWalletIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CelerWalletTransferToWalletIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CelerWalletTransferToWallet represents a TransferToWallet event raised by the CelerWallet contract.
type CelerWalletTransferToWallet struct {
	FromWalletId [32]byte
	ToWalletId   [32]byte
	TokenAddress common.Address
	Receiver     common.Address
	Amount       *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterTransferToWallet is a free log retrieval operation binding the contract event 0x1b56f805e5edb1e61b0d3f46feffdcbab5e591aa0e70e978ada9fc22093601c8.
//
// Solidity: event TransferToWallet(bytes32 indexed fromWalletId, bytes32 indexed toWalletId, address indexed tokenAddress, address receiver, uint256 amount)
func (_CelerWallet *CelerWalletFilterer) FilterTransferToWallet(opts *bind.FilterOpts, fromWalletId [][32]byte, toWalletId [][32]byte, tokenAddress []common.Address) (*CelerWalletTransferToWalletIterator, error) {

	var fromWalletIdRule []interface{}
	for _, fromWalletIdItem := range fromWalletId {
		fromWalletIdRule = append(fromWalletIdRule, fromWalletIdItem)
	}
	var toWalletIdRule []interface{}
	for _, toWalletIdItem := range toWalletId {
		toWalletIdRule = append(toWalletIdRule, toWalletIdItem)
	}
	var tokenAddressRule []interface{}
	for _, tokenAddressItem := range tokenAddress {
		tokenAddressRule = append(tokenAddressRule, tokenAddressItem)
	}

	logs, sub, err := _CelerWallet.contract.FilterLogs(opts, "TransferToWallet", fromWalletIdRule, toWalletIdRule, tokenAddressRule)
	if err != nil {
		return nil, err
	}
	return &CelerWalletTransferToWalletIterator{contract: _CelerWallet.contract, event: "TransferToWallet", logs: logs, sub: sub}, nil
}

// WatchTransferToWallet is a free log subscription operation binding the contract event 0x1b56f805e5edb1e61b0d3f46feffdcbab5e591aa0e70e978ada9fc22093601c8.
//
// Solidity: event TransferToWallet(bytes32 indexed fromWalletId, bytes32 indexed toWalletId, address indexed tokenAddress, address receiver, uint256 amount)
func (_CelerWallet *CelerWalletFilterer) WatchTransferToWallet(opts *bind.WatchOpts, sink chan<- *CelerWalletTransferToWallet, fromWalletId [][32]byte, toWalletId [][32]byte, tokenAddress []common.Address) (event.Subscription, error) {

	var fromWalletIdRule []interface{}
	for _, fromWalletIdItem := range fromWalletId {
		fromWalletIdRule = append(fromWalletIdRule, fromWalletIdItem)
	}
	var toWalletIdRule []interface{}
	for _, toWalletIdItem := range toWalletId {
		toWalletIdRule = append(toWalletIdRule, toWalletIdItem)
	}
	var tokenAddressRule []interface{}
	for _, tokenAddressItem := range tokenAddress {
		tokenAddressRule = append(tokenAddressRule, tokenAddressItem)
	}

	logs, sub, err := _CelerWallet.contract.WatchLogs(opts, "TransferToWallet", fromWalletIdRule, toWalletIdRule, tokenAddressRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CelerWalletTransferToWallet)
				if err := _CelerWallet.contract.UnpackLog(event, "TransferToWallet", log); err != nil {
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

// ParseTransferToWallet is a log parse operation binding the contract event 0x1b56f805e5edb1e61b0d3f46feffdcbab5e591aa0e70e978ada9fc22093601c8.
//
// Solidity: event TransferToWallet(bytes32 indexed fromWalletId, bytes32 indexed toWalletId, address indexed tokenAddress, address receiver, uint256 amount)
func (_CelerWallet *CelerWalletFilterer) ParseTransferToWallet(log types.Log) (*CelerWalletTransferToWallet, error) {
	event := new(CelerWalletTransferToWallet)
	if err := _CelerWallet.contract.UnpackLog(event, "TransferToWallet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CelerWalletUnpausedIterator is returned from FilterUnpaused and is used to iterate over the raw logs and unpacked data for Unpaused events raised by the CelerWallet contract.
type CelerWalletUnpausedIterator struct {
	Event *CelerWalletUnpaused // Event containing the contract specifics and raw log

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
func (it *CelerWalletUnpausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CelerWalletUnpaused)
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
		it.Event = new(CelerWalletUnpaused)
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
func (it *CelerWalletUnpausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CelerWalletUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CelerWalletUnpaused represents a Unpaused event raised by the CelerWallet contract.
type CelerWalletUnpaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterUnpaused is a free log retrieval operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_CelerWallet *CelerWalletFilterer) FilterUnpaused(opts *bind.FilterOpts) (*CelerWalletUnpausedIterator, error) {

	logs, sub, err := _CelerWallet.contract.FilterLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return &CelerWalletUnpausedIterator{contract: _CelerWallet.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

// WatchUnpaused is a free log subscription operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_CelerWallet *CelerWalletFilterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *CelerWalletUnpaused) (event.Subscription, error) {

	logs, sub, err := _CelerWallet.contract.WatchLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CelerWalletUnpaused)
				if err := _CelerWallet.contract.UnpackLog(event, "Unpaused", log); err != nil {
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

// ParseUnpaused is a log parse operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_CelerWallet *CelerWalletFilterer) ParseUnpaused(log types.Log) (*CelerWalletUnpaused, error) {
	event := new(CelerWalletUnpaused)
	if err := _CelerWallet.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CelerWalletWithdrawFromWalletIterator is returned from FilterWithdrawFromWallet and is used to iterate over the raw logs and unpacked data for WithdrawFromWallet events raised by the CelerWallet contract.
type CelerWalletWithdrawFromWalletIterator struct {
	Event *CelerWalletWithdrawFromWallet // Event containing the contract specifics and raw log

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
func (it *CelerWalletWithdrawFromWalletIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CelerWalletWithdrawFromWallet)
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
		it.Event = new(CelerWalletWithdrawFromWallet)
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
func (it *CelerWalletWithdrawFromWalletIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CelerWalletWithdrawFromWalletIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CelerWalletWithdrawFromWallet represents a WithdrawFromWallet event raised by the CelerWallet contract.
type CelerWalletWithdrawFromWallet struct {
	WalletId     [32]byte
	TokenAddress common.Address
	Receiver     common.Address
	Amount       *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterWithdrawFromWallet is a free log retrieval operation binding the contract event 0xd897e862036b62a0f770979fbd2227f3210565bba2eb4d9acd1dc8ccc00c928b.
//
// Solidity: event WithdrawFromWallet(bytes32 indexed walletId, address indexed tokenAddress, address indexed receiver, uint256 amount)
func (_CelerWallet *CelerWalletFilterer) FilterWithdrawFromWallet(opts *bind.FilterOpts, walletId [][32]byte, tokenAddress []common.Address, receiver []common.Address) (*CelerWalletWithdrawFromWalletIterator, error) {

	var walletIdRule []interface{}
	for _, walletIdItem := range walletId {
		walletIdRule = append(walletIdRule, walletIdItem)
	}
	var tokenAddressRule []interface{}
	for _, tokenAddressItem := range tokenAddress {
		tokenAddressRule = append(tokenAddressRule, tokenAddressItem)
	}
	var receiverRule []interface{}
	for _, receiverItem := range receiver {
		receiverRule = append(receiverRule, receiverItem)
	}

	logs, sub, err := _CelerWallet.contract.FilterLogs(opts, "WithdrawFromWallet", walletIdRule, tokenAddressRule, receiverRule)
	if err != nil {
		return nil, err
	}
	return &CelerWalletWithdrawFromWalletIterator{contract: _CelerWallet.contract, event: "WithdrawFromWallet", logs: logs, sub: sub}, nil
}

// WatchWithdrawFromWallet is a free log subscription operation binding the contract event 0xd897e862036b62a0f770979fbd2227f3210565bba2eb4d9acd1dc8ccc00c928b.
//
// Solidity: event WithdrawFromWallet(bytes32 indexed walletId, address indexed tokenAddress, address indexed receiver, uint256 amount)
func (_CelerWallet *CelerWalletFilterer) WatchWithdrawFromWallet(opts *bind.WatchOpts, sink chan<- *CelerWalletWithdrawFromWallet, walletId [][32]byte, tokenAddress []common.Address, receiver []common.Address) (event.Subscription, error) {

	var walletIdRule []interface{}
	for _, walletIdItem := range walletId {
		walletIdRule = append(walletIdRule, walletIdItem)
	}
	var tokenAddressRule []interface{}
	for _, tokenAddressItem := range tokenAddress {
		tokenAddressRule = append(tokenAddressRule, tokenAddressItem)
	}
	var receiverRule []interface{}
	for _, receiverItem := range receiver {
		receiverRule = append(receiverRule, receiverItem)
	}

	logs, sub, err := _CelerWallet.contract.WatchLogs(opts, "WithdrawFromWallet", walletIdRule, tokenAddressRule, receiverRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CelerWalletWithdrawFromWallet)
				if err := _CelerWallet.contract.UnpackLog(event, "WithdrawFromWallet", log); err != nil {
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

// ParseWithdrawFromWallet is a log parse operation binding the contract event 0xd897e862036b62a0f770979fbd2227f3210565bba2eb4d9acd1dc8ccc00c928b.
//
// Solidity: event WithdrawFromWallet(bytes32 indexed walletId, address indexed tokenAddress, address indexed receiver, uint256 amount)
func (_CelerWallet *CelerWalletFilterer) ParseWithdrawFromWallet(log types.Log) (*CelerWalletWithdrawFromWallet, error) {
	event := new(CelerWalletWithdrawFromWallet)
	if err := _CelerWallet.contract.UnpackLog(event, "WithdrawFromWallet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
