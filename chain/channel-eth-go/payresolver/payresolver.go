// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package payresolver

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

// PayResolverMetaData contains all meta data concerning the PayResolver contract.
var PayResolverMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"_registryAddr\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_virtResolverAddr\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"payRegistry\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIPayRegistry\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"resolvePaymentByConditions\",\"inputs\":[{\"name\":\"_resolvePayRequest\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"resolvePaymentByVouchedResult\",\"inputs\":[{\"name\":\"_vouchedPayResult\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"virtResolver\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIVirtContractResolver\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"ResolvePayment\",\"inputs\":[{\"name\":\"payId\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"resolveDeadline\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"ECDSAInvalidSignature\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ECDSAInvalidSignatureLength\",\"inputs\":[{\"name\":\"length\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ECDSAInvalidSignatureS\",\"inputs\":[{\"name\":\"s\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}]",
	Bin: "0x6080604052348015600e575f5ffd5b506040516124b93803806124b9833981016040819052602b916074565b5f80546001600160a01b039384166001600160a01b0319918216179091556001805492909316911617905560a0565b80516001600160a01b0381168114606f575f5ffd5b919050565b5f5f604083850312156084575f5ffd5b608b83605a565b9150609760208401605a565b90509250929050565b61240c806100ad5f395ff3fe608060405234801561000f575f5ffd5b506004361061004a575f3560e01c80634367e45e1461004e57806353fc513f146100635780635fff88c814610091578063ead54c1b146100a4575b5f5ffd5b61006161005c36600461217e565b6100b7565b005b5f54610075906001600160a01b031681565b6040516001600160a01b03909116815260200160405180910390f35b61006161009f36600461217e565b6101a7565b600154610075906001600160a01b031681565b5f6100f683838080601f0160208091040260200160405190810160405280939291908181526020018383808284375f9201919091525061037392505050565b90505f610105825f015161050e565b6080810151519091505f9081816005811115610123576101236121ec565b0361013d5761013683856020015161076a565b9150610189565b6001816005811115610151576101516121ec565b03610164576101368385602001516109ab565b61016d81610bfb565b156101815761013683856020015183610c4e565b610189612200565b8351805160209091012061019e8482856110b7565b50505050505050565b5f6101e683838080601f0160208091040260200160405190810160405280939291908181526020018383808284375f920191909152506114e692505050565b90505f6101f5825f01516115a2565b90505f610204825f015161050e565b905080608001516020015160200151602001518260200151111561026f5760405162461bcd60e51b815260206004820152601a60248201527f457863656564206d6178207472616e7366657220616d6f756e7400000000000060448201526064015b60405180910390fd5b825180516020918201207f19457468657265756d205369676e6564204d6573736167653a0a3332000000005f908152601c91909152603c8120918501516102b7908390611633565b90505f6102d186604001518461163390919063ffffffff16565b905083602001516001600160a01b0316826001600160a01b031614801561030d575083604001516001600160a01b0316816001600160a01b0316145b61034d5760405162461bcd60e51b815260206004820152601160248201527010da1958dac81cda59dcc819985a5b1959607a1b6044820152606401610266565b845180516020918201209086015161036890869083906110b7565b505050505050505050565b60408051808201909152606080825260208201525f6103a483604080518082019091525f8152602081019190915290565b90505f6103b282600261165b565b9050806002815181106103c7576103c7612214565b602002602001015167ffffffffffffffff8111156103e7576103e7612228565b60405190808252806020026020018201604052801561041a57816020015b60608152602001906001900390816104055790505b5083602001819052505f8160028151811061043757610437612214565b6020026020010181815250505f5f5b602084015151845110156105055761045d84611714565b90925090508160010361047a576104738461174c565b8552610446565b816002036104f65761048b8461174c565b8560200151846002815181106104a3576104a3612214565b6020026020010151815181106104bb576104bb612214565b6020026020010181905250826002815181106104d9576104d9612214565b6020026020010180518091906104ee90612250565b905250610446565b6105008482611804565b610446565b50505050919050565b610516612092565b604080518082019091525f8082526020820184905261053682600861165b565b90508060048151811061054b5761054b612214565b602002602001015167ffffffffffffffff81111561056b5761056b612228565b6040519080825280602002602001820160405280156105a457816020015b6105916120f2565b8152602001906001900390816105895790505b5083606001819052505f816004815181106105c1576105c1612214565b6020026020010181815250505f5f5b60208401515184511015610505576105e784611714565b909250905081600103610604576105fd84611879565b85526105d0565b816002036106305761061d6106188561174c565b6118e8565b6001600160a01b031660208601526105d0565b81600303610657576106446106188561174c565b6001600160a01b031660408601526105d0565b816004036106db5761067061066b8561174c565b6118f2565b85606001518460048151811061068857610688612214565b6020026020010151815181106106a0576106a0612214565b6020026020010181905250826004815181106106be576106be612214565b6020026020010180518091906106d390612250565b9052506105d0565b816005036106fe576106f46106ef8561174c565b611a2a565b60808601526105d0565b816006036107195761070f84611879565b60a08601526105d0565b816007036107345761072a84611879565b60c08601526105d0565b8160080361075b576107486106188561174c565b6001600160a01b031660e08601526105d0565b6107658482611804565b6105d0565b5f8080805b85606001515181101561097e575f8660600151828151811061079357610793612214565b602002602001015190505f60028111156107af576107af6121ec565b815160028111156107c2576107c26121ec565b0361081c5780602001518685815181106107de576107de612214565b6020026020010151805190602001201461080a5760405162461bcd60e51b815260040161026690612268565b8361081481612250565b945050610975565b600181516002811115610831576108316121ec565b148061084f575060028151600281111561084d5761084d6121ec565b145b1561096d575f61085e82611ae7565b6080830151604051632f36f6a560e21b815291925082916001600160a01b0383169163bcdbda94916108939190600401612290565b602060405180830381865afa1580156108ae573d5f5f3e3d5ffd5b505050506040513d601f19601f820116820180604052508101906108d291906122c5565b6108ee5760405162461bcd60e51b8152600401610266906122e4565b60a083015160405163ea4ba8eb60e01b81526001600160a01b0383169163ea4ba8eb9161091e9190600401612290565b602060405180830381865afa158015610939573d5f5f3e3d5ffd5b505050506040513d601f19601f8201168201806040525081019061095d91906122c5565b61096657600194505b5050610975565b610975612200565b5060010161076f565b50801561098f575f925050506109a5565b8460800151602001516020015160200151925050505b92915050565b5f808080805b866060015151811015610bc5575f876060015182815181106109d5576109d5612214565b602002602001015190505f60028111156109f1576109f16121ec565b81516002811115610a0457610a046121ec565b03610a5e578060200151878681518110610a2057610a20612214565b60200260200101518051906020012014610a4c5760405162461bcd60e51b815260040161026690612268565b84610a5681612250565b955050610bbc565b600181516002811115610a7357610a736121ec565b1480610a915750600281516002811115610a8f57610a8f6121ec565b145b15610bb4575f610aa082611ae7565b6080830151604051632f36f6a560e21b815291925082916001600160a01b0383169163bcdbda9491610ad59190600401612290565b602060405180830381865afa158015610af0573d5f5f3e3d5ffd5b505050506040513d601f19601f82011682018060405250810190610b1491906122c5565b610b305760405162461bcd60e51b8152600401610266906122e4565b60a083015160405163ea4ba8eb60e01b8152600197506001600160a01b0383169163ea4ba8eb91610b649190600401612290565b602060405180830381865afa158015610b7f573d5f5f3e3d5ffd5b505050506040513d601f19601f82011682018060405250810190610ba391906122c5565b15610bad57600194505b5050610bbc565b610bbc612200565b506001016109b1565b50811580610bd05750805b15610bf057856080015160200151602001516020015193505050506109a5565b5f93505050506109a5565b5f6003826005811115610c1057610c106121ec565b1480610c2d57506004826005811115610c2b57610c2b6121ec565b145b806109a557506005826005811115610c4757610c476121ec565b1492915050565b5f808080805b876060015151811015611027575f88606001518281518110610c7857610c78612214565b602002602001015190505f6002811115610c9457610c946121ec565b81516002811115610ca757610ca76121ec565b03610d01578060200151888581518110610cc357610cc3612214565b60200260200101518051906020012014610cef5760405162461bcd60e51b815260040161026690612268565b83610cf981612250565b94505061101e565b600181516002811115610d1657610d166121ec565b1480610d345750600281516002811115610d3257610d326121ec565b145b15611016575f610d4382611ae7565b6080830151604051632f36f6a560e21b815291925082916001600160a01b0383169163bcdbda9491610d789190600401612290565b602060405180830381865afa158015610d93573d5f5f3e3d5ffd5b505050506040513d601f19601f82011682018060405250810190610db791906122c5565b610dd35760405162461bcd60e51b8152600401610266906122e4565b6003896005811115610de757610de76121ec565b03610e6c5760a083015160405163ea4ba8eb60e01b81526001600160a01b0383169163ea4ba8eb91610e1c9190600401612290565b602060405180830381865afa158015610e37573d5f5f3e3d5ffd5b505050506040513d601f19601f82011682018060405250810190610e5b919061231b565b610e659088612332565b965061100b565b6004896005811115610e8057610e806121ec565b03610efd57610e6587826001600160a01b031663ea4ba8eb8660a001516040518263ffffffff1660e01b8152600401610eb99190612290565b602060405180830381865afa158015610ed4573d5f5f3e3d5ffd5b505050506040513d601f19601f82011682018060405250810190610ef8919061231b565b611bdb565b6005896005811115610f1157610f116121ec565b03611003578415610f9457610e6587826001600160a01b031663ea4ba8eb8660a001516040518263ffffffff1660e01b8152600401610f509190612290565b602060405180830381865afa158015610f6b573d5f5f3e3d5ffd5b505050506040513d601f19601f82011682018060405250810190610f8f919061231b565b611bea565b60a083015160405163ea4ba8eb60e01b81526001600160a01b0383169163ea4ba8eb91610fc49190600401612290565b602060405180830381865afa158015610fdf573d5f5f3e3d5ffd5b505050506040513d601f19601f82011682018060405250810190610e65919061231b565b61100b612200565b60019450505061101e565b61101e612200565b50600101610c54565b50801561109957866080015160200151602001516020015183111561108e5760405162461bcd60e51b815260206004820152601a60248201527f457863656564206d6178207472616e7366657220616d6f756e740000000000006044820152606401610266565b8293505050506110b0565b866080015160200151602001516020015193505050505b9392505050565b60a083015142908111156111205760405162461bcd60e51b815260206004820152602a60248201527f50617373656420706179207265736f6c766520646561646c696e6520696e20636044820152696f6e64506179206d736760b01b6064820152608401610266565b5f61112b8430611bf9565b5f80546040516304f61c0b60e31b815260048101849052929350909182916001600160a01b0316906327b0e058906024016040805180830381865afa158015611176573d5f5f3e3d5ffd5b505050506040513d601f19601f8201168201806040525081019061119a9190612345565b91509150805f14806111ac5750808411155b6112045760405162461bcd60e51b815260206004820152602360248201527f506173736564206f6e636861696e207265736f6c76652070617920646561646c604482015262696e6560e81b6064820152608401610266565b80156113ae578185116112595760405162461bcd60e51b815260206004820152601860248201527f4e657720616d6f756e74206973206e6f74206c617267657200000000000000006044820152606401610266565b86608001516020015160200151602001518503611317575f54604051630e1e354960e41b81526004810188905260248101879052604481018690526001600160a01b039091169063e1e35490906064015f604051808303815f87803b1580156112c0575f5ffd5b505af11580156112d2573d5f5f3e3d5ffd5b505060408051888152602081018890528693507fa87e293885636c5018108e8ee0e41d65206d1dfc0a9066f26f2a91a78b2beb179250015b60405180910390a261019e565b5f5460405163f8fb012f60e01b815260048101889052602481018790526001600160a01b039091169063f8fb012f906044015f604051808303815f87803b158015611360575f5ffd5b505af1158015611372573d5f5f3e3d5ffd5b505060408051888152602081018590528693507fa87e293885636c5018108e8ee0e41d65206d1dfc0a9066f26f2a91a78b2beb1792500161130a565b5f876080015160200151602001516020015186036113cd57508361143b565b6113ea8860c00151866113e09190612332565b8960a00151611bea565b90505f811161143b5760405162461bcd60e51b815260206004820152601960248201527f4e6577207265736f6c766520646561646c696e652069732030000000000000006044820152606401610266565b5f54604051630e1e354960e41b81526004810189905260248101889052604481018390526001600160a01b039091169063e1e35490906064015f604051808303815f87803b15801561148b575f5ffd5b505af115801561149d573d5f5f3e3d5ffd5b505060408051898152602081018590528793507fa87e293885636c5018108e8ee0e41d65206d1dfc0a9066f26f2a91a78b2beb1792500160405180910390a25050505050505050565b61150a60405180606001604052806060815260200160608152602001606081525090565b604080518082019091525f80825260208201849052805b6020830151518351101561159a5761153883611714565b9092509050816001036115555761154e8361174c565b8452611521565b81600203611570576115668361174c565b6020850152611521565b8160030361158b576115818361174c565b6040850152611521565b6115958382611804565b611521565b505050919050565b604080518082018252606081525f602080830182905283518085019094528184528301849052909190805b6020830151518351101561159a576115e483611714565b909250905081600103611601576115fa8361174c565b84526115cd565b816002036116245761161a6116158461174c565b611c46565b60208501526115cd565b61162e8382611804565b6115cd565b5f5f5f5f6116418686611c7b565b9250925092506116518282611cc4565b5090949350505050565b815160609061166b836001612332565b67ffffffffffffffff81111561168357611683612228565b6040519080825280602002602001820160405280156116ac578160200160208202803683370190505b5091505f5f5b6020860151518651101561170b576116c986611714565b809250819350505060018483815181106116e5576116e5612214565b602002602001018181516116f99190612332565b9052506117068682611804565b6116b2565b50509092525090565b5f5f5f61172084611879565b905061172d600882612367565b9250806007166005811115611744576117446121ec565b915050915091565b60605f61175883611879565b90505f81845f015161176a9190612332565b905083602001515181111561177d575f5ffd5b8167ffffffffffffffff81111561179657611796612228565b6040519080825280601f01601f1916602001820160405280156117c0576020820181803683370190505b5060208086015186519295509181860191908301015f5b858110156117f95781810151838201526117f2602082612332565b90506117d7565b505050935250919050565b5f816005811115611817576118176121ec565b0361182a5761182582611879565b505050565b600281600581111561183e5761183e6121ec565b0361004a575f61184d83611879565b905080835f018181516118609190612332565b90525060208301515183511115611825575f5ffd5b5050565b60208082015182518101909101515f9182805b600a81101561004a5783811a91506118a5816007612386565b82607f16901b85179450816080165f036118e0576118c4816001612332565b865187906118d3908390612332565b9052509395945050505050565b60010161188c565b5f6109a582611d7c565b6118fa6120f2565b604080518082019091525f80825260208201849052805b6020830151518351101561159a5761192883611714565b90925090508160010361197d5761193e83611879565b600281111561194f5761194f6121ec565b84906002811115611962576119626121ec565b90816002811115611975576119756121ec565b905250611911565b816002036119a0576119966119918461174c565b611d99565b6020850152611911565b816003036119c7576119b46106188461174c565b6001600160a01b03166040850152611911565b816004036119e5576119db6119918461174c565b6060850152611911565b81600503611a00576119f68361174c565b6080850152611911565b81600603611a1b57611a118361174c565b60a0850152611911565b611a258382611804565b611911565b611a3261212e565b604080518082019091525f80825260208201849052805b6020830151518351101561159a57611a6083611714565b909250905081600103611ab557611a7683611879565b6005811115611a8757611a876121ec565b84906005811115611a9a57611a9a6121ec565b90816005811115611aad57611aad6121ec565b905250611a49565b81600203611ad857611ace611ac98461174c565b611daf565b6020850152611a49565b611ae28382611804565b611a49565b5f600182516002811115611afd57611afd6121ec565b03611b0a57506040015190565b600282516002811115611b1f57611b1f6121ec565b03611b9a576001546060830151604051635c23bdf560e01b81526001600160a01b0390921691635c23bdf591611b5b9160040190815260200190565b602060405180830381865afa158015611b76573d5f5f3e3d5ffd5b505050506040513d601f19601f820116820180604052508101906109a5919061239d565b60405162461bcd60e51b8152602060048201526016602482015275496e76616c696420636f6e646974696f6e207479706560501b6044820152606401610266565b5f8282188284110282186110b0565b5f8282188284100282186110b0565b5f8282604051602001611c2892919091825260601b6bffffffffffffffffffffffff1916602082015260340190565b60405160208183030381529060405280519060200120905092915050565b5f602082511115611c55575f5ffd5b6020820151905081516020611c6a91906123c3565b611c75906008612386565b1c919050565b5f5f5f8351604103611cb2576020840151604085015160608601515f1a611ca488828585611e65565b955095509550505050611cbd565b505081515f91506002905b9250925092565b5f826003811115611cd757611cd76121ec565b03611ce0575050565b6001826003811115611cf457611cf46121ec565b03611d125760405163f645eedf60e01b815260040160405180910390fd5b6002826003811115611d2657611d266121ec565b03611d475760405163fce698f760e01b815260048101829052602401610266565b6003826003811115611d5b57611d5b6121ec565b03611875576040516335e2f38360e21b815260048101829052602401610266565b5f8151601414611d8a575f5ffd5b5060200151600160601b900490565b5f8151602014611da7575f5ffd5b506020015190565b604080516080810182525f8183018181526060830182905282528251808401845281815260208082018390528084019190915283518085019094528184528301849052909190805b6020830151518351101561159a57611e0e83611714565b909250905081600103611e3357611e2c611e278461174c565b611f2d565b8452611df7565b81600203611e5657611e4c611e478461174c565b611ffb565b6020850152611df7565b611e608382611804565b611df7565b5f80807f7fffffffffffffffffffffffffffffff5d576e7357a4501ddfe92f46681b20a0841115611e9e57505f91506003905082611f23565b604080515f808252602082018084528a905260ff891692820192909252606081018790526080810186905260019060a0016020604051602081039080840390855afa158015611eef573d5f5f3e3d5ffd5b5050604051601f1901519150506001600160a01b038116611f1a57505f925060019150829050611f23565b92505f91508190505b9450945094915050565b604080518082019091525f8082526020820152604080518082019091525f80825260208201849052505f5f5b6020830151518351101561159a57611f7083611714565b909250905081600103611fc557611f8683611879565b6002811115611f9757611f976121ec565b84906002811115611faa57611faa6121ec565b90816002811115611fbd57611fbd6121ec565b905250611f59565b81600203611fec57611fd96106188461174c565b6001600160a01b03166020850152611f59565b611ff68382611804565b611f59565b6040805180820182525f808252602080830182905283518085019094528184528301849052909190805b6020830151518351101561159a5761203c83611714565b909250905081600103612065576120556106188461174c565b6001600160a01b03168452612025565b81600203612083576120796116158461174c565b6020850152612025565b61208d8382611804565b612025565b6040518061010001604052805f81526020015f6001600160a01b031681526020015f6001600160a01b03168152602001606081526020016120d161212e565b81526020015f81526020015f81526020015f6001600160a01b031681525090565b6040805160c08101909152805f81526020015f81526020015f6001600160a01b031681526020015f815260200160608152602001606081525090565b60408051808201909152805f8152602001612179604080516080810182525f818301818152606083018290528252825180840190935280835260208381019190915290919082015290565b905290565b5f5f6020838503121561218f575f5ffd5b823567ffffffffffffffff8111156121a5575f5ffd5b8301601f810185136121b5575f5ffd5b803567ffffffffffffffff8111156121cb575f5ffd5b8560208284010111156121dc575f5ffd5b6020919091019590945092505050565b634e487b7160e01b5f52602160045260245ffd5b634e487b7160e01b5f52600160045260245ffd5b634e487b7160e01b5f52603260045260245ffd5b634e487b7160e01b5f52604160045260245ffd5b634e487b7160e01b5f52601160045260245ffd5b5f600182016122615761226161223c565b5060010190565b6020808252600e908201526d57726f6e6720707265696d61676560901b604082015260600190565b602081525f82518060208401528060208501604085015e5f604082850101526040601f19601f83011684010191505092915050565b5f602082840312156122d5575f5ffd5b815180151581146110b0575f5ffd5b6020808252601a908201527f436f6e646974696f6e206973206e6f742066696e616c697a6564000000000000604082015260600190565b5f6020828403121561232b575f5ffd5b5051919050565b808201808211156109a5576109a561223c565b5f5f60408385031215612356575f5ffd5b505080516020909101519092909150565b5f8261238157634e487b7160e01b5f52601260045260245ffd5b500490565b80820281158282048414176109a5576109a561223c565b5f602082840312156123ad575f5ffd5b81516001600160a01b03811681146110b0575f5ffd5b818103818111156109a5576109a561223c56fea26469706673582212208f1004bf702ea90805aa2480ee6a8286b55ec76d95408c5915a5b753208bb2e264736f6c634300081e0033",
}

// PayResolverABI is the input ABI used to generate the binding from.
// Deprecated: Use PayResolverMetaData.ABI instead.
var PayResolverABI = PayResolverMetaData.ABI

// PayResolverBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use PayResolverMetaData.Bin instead.
var PayResolverBin = PayResolverMetaData.Bin

// DeployPayResolver deploys a new Ethereum contract, binding an instance of PayResolver to it.
func DeployPayResolver(auth *bind.TransactOpts, backend bind.ContractBackend, _registryAddr common.Address, _virtResolverAddr common.Address) (common.Address, *types.Transaction, *PayResolver, error) {
	parsed, err := PayResolverMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(PayResolverBin), backend, _registryAddr, _virtResolverAddr)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &PayResolver{PayResolverCaller: PayResolverCaller{contract: contract}, PayResolverTransactor: PayResolverTransactor{contract: contract}, PayResolverFilterer: PayResolverFilterer{contract: contract}}, nil
}

// PayResolver is an auto generated Go binding around an Ethereum contract.
type PayResolver struct {
	PayResolverCaller     // Read-only binding to the contract
	PayResolverTransactor // Write-only binding to the contract
	PayResolverFilterer   // Log filterer for contract events
}

// PayResolverCaller is an auto generated read-only Go binding around an Ethereum contract.
type PayResolverCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PayResolverTransactor is an auto generated write-only Go binding around an Ethereum contract.
type PayResolverTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PayResolverFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type PayResolverFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PayResolverSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type PayResolverSession struct {
	Contract     *PayResolver      // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// PayResolverCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type PayResolverCallerSession struct {
	Contract *PayResolverCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts      // Call options to use throughout this session
}

// PayResolverTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type PayResolverTransactorSession struct {
	Contract     *PayResolverTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// PayResolverRaw is an auto generated low-level Go binding around an Ethereum contract.
type PayResolverRaw struct {
	Contract *PayResolver // Generic contract binding to access the raw methods on
}

// PayResolverCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type PayResolverCallerRaw struct {
	Contract *PayResolverCaller // Generic read-only contract binding to access the raw methods on
}

// PayResolverTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type PayResolverTransactorRaw struct {
	Contract *PayResolverTransactor // Generic write-only contract binding to access the raw methods on
}

// NewPayResolver creates a new instance of PayResolver, bound to a specific deployed contract.
func NewPayResolver(address common.Address, backend bind.ContractBackend) (*PayResolver, error) {
	contract, err := bindPayResolver(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &PayResolver{PayResolverCaller: PayResolverCaller{contract: contract}, PayResolverTransactor: PayResolverTransactor{contract: contract}, PayResolverFilterer: PayResolverFilterer{contract: contract}}, nil
}

// NewPayResolverCaller creates a new read-only instance of PayResolver, bound to a specific deployed contract.
func NewPayResolverCaller(address common.Address, caller bind.ContractCaller) (*PayResolverCaller, error) {
	contract, err := bindPayResolver(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &PayResolverCaller{contract: contract}, nil
}

// NewPayResolverTransactor creates a new write-only instance of PayResolver, bound to a specific deployed contract.
func NewPayResolverTransactor(address common.Address, transactor bind.ContractTransactor) (*PayResolverTransactor, error) {
	contract, err := bindPayResolver(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &PayResolverTransactor{contract: contract}, nil
}

// NewPayResolverFilterer creates a new log filterer instance of PayResolver, bound to a specific deployed contract.
func NewPayResolverFilterer(address common.Address, filterer bind.ContractFilterer) (*PayResolverFilterer, error) {
	contract, err := bindPayResolver(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &PayResolverFilterer{contract: contract}, nil
}

// bindPayResolver binds a generic wrapper to an already deployed contract.
func bindPayResolver(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := PayResolverMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_PayResolver *PayResolverRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _PayResolver.Contract.PayResolverCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_PayResolver *PayResolverRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PayResolver.Contract.PayResolverTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_PayResolver *PayResolverRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _PayResolver.Contract.PayResolverTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_PayResolver *PayResolverCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _PayResolver.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_PayResolver *PayResolverTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PayResolver.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_PayResolver *PayResolverTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _PayResolver.Contract.contract.Transact(opts, method, params...)
}

// PayRegistry is a free data retrieval call binding the contract method 0x53fc513f.
//
// Solidity: function payRegistry() view returns(address)
func (_PayResolver *PayResolverCaller) PayRegistry(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _PayResolver.contract.Call(opts, &out, "payRegistry")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// PayRegistry is a free data retrieval call binding the contract method 0x53fc513f.
//
// Solidity: function payRegistry() view returns(address)
func (_PayResolver *PayResolverSession) PayRegistry() (common.Address, error) {
	return _PayResolver.Contract.PayRegistry(&_PayResolver.CallOpts)
}

// PayRegistry is a free data retrieval call binding the contract method 0x53fc513f.
//
// Solidity: function payRegistry() view returns(address)
func (_PayResolver *PayResolverCallerSession) PayRegistry() (common.Address, error) {
	return _PayResolver.Contract.PayRegistry(&_PayResolver.CallOpts)
}

// VirtResolver is a free data retrieval call binding the contract method 0xead54c1b.
//
// Solidity: function virtResolver() view returns(address)
func (_PayResolver *PayResolverCaller) VirtResolver(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _PayResolver.contract.Call(opts, &out, "virtResolver")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// VirtResolver is a free data retrieval call binding the contract method 0xead54c1b.
//
// Solidity: function virtResolver() view returns(address)
func (_PayResolver *PayResolverSession) VirtResolver() (common.Address, error) {
	return _PayResolver.Contract.VirtResolver(&_PayResolver.CallOpts)
}

// VirtResolver is a free data retrieval call binding the contract method 0xead54c1b.
//
// Solidity: function virtResolver() view returns(address)
func (_PayResolver *PayResolverCallerSession) VirtResolver() (common.Address, error) {
	return _PayResolver.Contract.VirtResolver(&_PayResolver.CallOpts)
}

// ResolvePaymentByConditions is a paid mutator transaction binding the contract method 0x4367e45e.
//
// Solidity: function resolvePaymentByConditions(bytes _resolvePayRequest) returns()
func (_PayResolver *PayResolverTransactor) ResolvePaymentByConditions(opts *bind.TransactOpts, _resolvePayRequest []byte) (*types.Transaction, error) {
	return _PayResolver.contract.Transact(opts, "resolvePaymentByConditions", _resolvePayRequest)
}

// ResolvePaymentByConditions is a paid mutator transaction binding the contract method 0x4367e45e.
//
// Solidity: function resolvePaymentByConditions(bytes _resolvePayRequest) returns()
func (_PayResolver *PayResolverSession) ResolvePaymentByConditions(_resolvePayRequest []byte) (*types.Transaction, error) {
	return _PayResolver.Contract.ResolvePaymentByConditions(&_PayResolver.TransactOpts, _resolvePayRequest)
}

// ResolvePaymentByConditions is a paid mutator transaction binding the contract method 0x4367e45e.
//
// Solidity: function resolvePaymentByConditions(bytes _resolvePayRequest) returns()
func (_PayResolver *PayResolverTransactorSession) ResolvePaymentByConditions(_resolvePayRequest []byte) (*types.Transaction, error) {
	return _PayResolver.Contract.ResolvePaymentByConditions(&_PayResolver.TransactOpts, _resolvePayRequest)
}

// ResolvePaymentByVouchedResult is a paid mutator transaction binding the contract method 0x5fff88c8.
//
// Solidity: function resolvePaymentByVouchedResult(bytes _vouchedPayResult) returns()
func (_PayResolver *PayResolverTransactor) ResolvePaymentByVouchedResult(opts *bind.TransactOpts, _vouchedPayResult []byte) (*types.Transaction, error) {
	return _PayResolver.contract.Transact(opts, "resolvePaymentByVouchedResult", _vouchedPayResult)
}

// ResolvePaymentByVouchedResult is a paid mutator transaction binding the contract method 0x5fff88c8.
//
// Solidity: function resolvePaymentByVouchedResult(bytes _vouchedPayResult) returns()
func (_PayResolver *PayResolverSession) ResolvePaymentByVouchedResult(_vouchedPayResult []byte) (*types.Transaction, error) {
	return _PayResolver.Contract.ResolvePaymentByVouchedResult(&_PayResolver.TransactOpts, _vouchedPayResult)
}

// ResolvePaymentByVouchedResult is a paid mutator transaction binding the contract method 0x5fff88c8.
//
// Solidity: function resolvePaymentByVouchedResult(bytes _vouchedPayResult) returns()
func (_PayResolver *PayResolverTransactorSession) ResolvePaymentByVouchedResult(_vouchedPayResult []byte) (*types.Transaction, error) {
	return _PayResolver.Contract.ResolvePaymentByVouchedResult(&_PayResolver.TransactOpts, _vouchedPayResult)
}

// PayResolverResolvePaymentIterator is returned from FilterResolvePayment and is used to iterate over the raw logs and unpacked data for ResolvePayment events raised by the PayResolver contract.
type PayResolverResolvePaymentIterator struct {
	Event *PayResolverResolvePayment // Event containing the contract specifics and raw log

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
func (it *PayResolverResolvePaymentIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PayResolverResolvePayment)
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
		it.Event = new(PayResolverResolvePayment)
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
func (it *PayResolverResolvePaymentIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PayResolverResolvePaymentIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PayResolverResolvePayment represents a ResolvePayment event raised by the PayResolver contract.
type PayResolverResolvePayment struct {
	PayId           [32]byte
	Amount          *big.Int
	ResolveDeadline *big.Int
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterResolvePayment is a free log retrieval operation binding the contract event 0xa87e293885636c5018108e8ee0e41d65206d1dfc0a9066f26f2a91a78b2beb17.
//
// Solidity: event ResolvePayment(bytes32 indexed payId, uint256 amount, uint256 resolveDeadline)
func (_PayResolver *PayResolverFilterer) FilterResolvePayment(opts *bind.FilterOpts, payId [][32]byte) (*PayResolverResolvePaymentIterator, error) {

	var payIdRule []interface{}
	for _, payIdItem := range payId {
		payIdRule = append(payIdRule, payIdItem)
	}

	logs, sub, err := _PayResolver.contract.FilterLogs(opts, "ResolvePayment", payIdRule)
	if err != nil {
		return nil, err
	}
	return &PayResolverResolvePaymentIterator{contract: _PayResolver.contract, event: "ResolvePayment", logs: logs, sub: sub}, nil
}

// WatchResolvePayment is a free log subscription operation binding the contract event 0xa87e293885636c5018108e8ee0e41d65206d1dfc0a9066f26f2a91a78b2beb17.
//
// Solidity: event ResolvePayment(bytes32 indexed payId, uint256 amount, uint256 resolveDeadline)
func (_PayResolver *PayResolverFilterer) WatchResolvePayment(opts *bind.WatchOpts, sink chan<- *PayResolverResolvePayment, payId [][32]byte) (event.Subscription, error) {

	var payIdRule []interface{}
	for _, payIdItem := range payId {
		payIdRule = append(payIdRule, payIdItem)
	}

	logs, sub, err := _PayResolver.contract.WatchLogs(opts, "ResolvePayment", payIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PayResolverResolvePayment)
				if err := _PayResolver.contract.UnpackLog(event, "ResolvePayment", log); err != nil {
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

// ParseResolvePayment is a log parse operation binding the contract event 0xa87e293885636c5018108e8ee0e41d65206d1dfc0a9066f26f2a91a78b2beb17.
//
// Solidity: event ResolvePayment(bytes32 indexed payId, uint256 amount, uint256 resolveDeadline)
func (_PayResolver *PayResolverFilterer) ParseResolvePayment(log types.Log) (*PayResolverResolvePayment, error) {
	event := new(PayResolverResolvePayment)
	if err := _PayResolver.contract.UnpackLog(event, "ResolvePayment", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
