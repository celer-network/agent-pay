// Copyright 2018-2025 Celer Network

package app

import (
	"fmt"
	"math/big"
	"sync"

	"github.com/celer-network/agent-pay/chain/channel-eth-go/virtresolver"
	"github.com/celer-network/agent-pay/common"
	"github.com/celer-network/agent-pay/common/intfs"
	"github.com/celer-network/agent-pay/config"
	"github.com/celer-network/agent-pay/ctype"
	"github.com/celer-network/agent-pay/entity"
	"github.com/celer-network/agent-pay/storage"
	"github.com/celer-network/agent-pay/utils"
	"github.com/celer-network/goutils/eth"
	"github.com/celer-network/goutils/log"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

// AppChannel tracks a registered VIRTUAL_CONTRACT condition contract: the
// bytecode + constructor + nonce that determine its deterministic address, and
// (after the first deploy-on-query) the on-chain deployed address.
//
// Post-trim the legacy gaming surface (turn-based session state machine,
// `Callback` / `OnDispute` notifications, `Players` / `Session` for the
// deployed-multisession variant, `IntendSettle` / dispute-window / action loop)
// is gone. What remains is the bytecode-and-deploy-on-demand bookkeeping for
// stateless `IBooleanCond` virtual condition contracts.
type AppChannel struct {
	Type           entity.ConditionType
	Nonce          uint64
	ByteCode       []byte
	Constructor   []byte
	DeployedAddr   ctype.Addr
	OnChainTimeout uint64
	mu             sync.Mutex
	client         *AppClient
	cid            string
}

type AppClient struct {
	nodeConfig     common.GlobalNodeConfig
	transactor     *eth.Transactor
	transactorPool *eth.TransactorPool
	monitorService intfs.MonitorService
	dal            *storage.DAL
	signer         eth.Signer
	appChannels    map[string]*AppChannel
	cLock          sync.RWMutex
}

func NewAppClient(
	nodeConfig common.GlobalNodeConfig,
	transactor *eth.Transactor,
	transactorPool *eth.TransactorPool,
	monitorService intfs.MonitorService,
	dal *storage.DAL,
	signer eth.Signer,
) *AppClient {
	return &AppClient{
		nodeConfig:     nodeConfig,
		transactor:     transactor,
		transactorPool: transactorPool,
		monitorService: monitorService,
		dal:            dal,
		signer:         signer,
		appChannels:    make(map[string]*AppChannel),
	}
}

func (a *AppChannel) setDeployedAddr(addr ctype.Addr) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.DeployedAddr = addr
}

func (a *AppChannel) getDeployedAddr() ctype.Addr {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.DeployedAddr
}

func (c *AppClient) PutAppChannel(cid string, appChannel *AppChannel) {
	c.cLock.Lock()
	defer c.cLock.Unlock()
	c.appChannels[cid] = appChannel
}

func (c *AppClient) GetAppChannel(cid string) *AppChannel {
	c.cLock.RLock()
	defer c.cLock.RUnlock()
	return c.appChannels[cid]
}

// DeleteAppChannel removes the in-memory bookkeeping for a registered virtual
// condition contract. It does not touch any on-chain state.
func (c *AppClient) DeleteAppChannel(cid string) {
	c.cLock.Lock()
	delete(c.appChannels, cid)
	c.cLock.Unlock()
}

// NewAppChannelOnVirtualContract registers a VIRTUAL_CONTRACT condition
// contract. The cnode stores the bytecode + constructor + nonce so that, when
// dispute resolution requires it, the contract can be deployed on-chain and
// queried via `IBooleanCond.{isFinalized,getOutcome}`. Returns the deterministic
// virtual-contract address (hex) used as the session id / `OnChainAddress` in
// `Condition` payloads.
func (c *AppClient) NewAppChannelOnVirtualContract(
	byteCode []byte,
	constructor []byte,
	nonce uint64,
	onchainTimeout uint64) (string, error) {

	cid := ctype.Bytes2Hex(GetVirtualAddress(byteCode, constructor, nonce))
	appChannel := &AppChannel{
		Type:           entity.ConditionType_VIRTUAL_CONTRACT,
		Nonce:          nonce,
		ByteCode:       byteCode,
		Constructor:    constructor,
		DeployedAddr:   ctype.ZeroAddr,
		OnChainTimeout: onchainTimeout,
		client:         c,
		cid:            cid,
	}
	c.PutAppChannel(cid, appChannel)
	return cid, nil
}

// GetAppChannelDeployedAddr returns the on-chain deployed address of a
// registered virtual condition contract. If the contract has not been deployed
// yet, it probes the virt-resolver registry; returns an error if not deployed.
func (c *AppClient) GetAppChannelDeployedAddr(cid string) (ctype.Addr, error) {
	appChannel := c.GetAppChannel(cid)
	if appChannel == nil {
		return ctype.ZeroAddr, fmt.Errorf("app channel not found")
	}
	addr := appChannel.getDeployedAddr()
	if addr != (ctype.ZeroAddr) || appChannel.Type == entity.ConditionType_DEPLOYED_CONTRACT {
		return addr, nil
	}
	virtAddr := GetVirtualAddress(appChannel.ByteCode, appChannel.Constructor, appChannel.Nonce)
	deployed, addr, err := c.isDeployed(virtAddr)
	if err != nil {
		return addr, err
	}
	if !deployed {
		return ctype.ZeroAddr, fmt.Errorf("virtual contract not deployed")
	}
	appChannel.setDeployedAddr(addr)
	return addr, nil
}

// GetBooleanOutcome queries `IBooleanCond.{isFinalized,getOutcome}` for the
// registered condition contract. For VIRTUAL_CONTRACT, this triggers
// deploy-on-query: if the virtual contract has not been deployed yet, this call
// submits a deployment transaction first. The query bytes are passed through
// unchanged (matches what `PayResolver` does on-chain) — no `SessionQuery`
// wrapping.
func (c *AppClient) GetBooleanOutcome(cid string, query []byte) (bool, bool, error) {
	appChannel := c.GetAppChannel(cid)
	if appChannel == nil {
		return false, false, fmt.Errorf("GetBooleanOutcome error: app channel not found")
	}
	if err := c.deployIfNeeded(appChannel); err != nil {
		return false, false, err
	}
	deployedAddr := appChannel.getDeployedAddr()
	contract, err := NewIBooleanCondCaller(deployedAddr, c.transactorPool.ContractCaller())
	if err != nil {
		return false, false, fmt.Errorf("GetBooleanOutcome error: %w", err)
	}
	finalized, err := contract.IsFinalized(&bind.CallOpts{}, query)
	if err != nil {
		return false, false, fmt.Errorf("contract IsFinalized error: %w", err)
	}
	result, err := contract.GetOutcome(&bind.CallOpts{}, query)
	if err != nil {
		return false, false, fmt.Errorf("contract GetOutcome error: %w", err)
	}
	return finalized, result, nil
}

// deployIfNeeded ensures the registered virtual condition contract is deployed
// on-chain. For VIRTUAL_CONTRACT entries with no deployed address yet, it
// submits the deployment transaction via the virt-resolver and caches the
// resulting address on the `AppChannel`.
func (c *AppClient) deployIfNeeded(appChannel *AppChannel) error {
	deployedAddr := appChannel.getDeployedAddr()
	if appChannel.Type == entity.ConditionType_VIRTUAL_CONTRACT && deployedAddr == (ctype.ZeroAddr) {
		deployedAddr, err :=
			c.deployVirtualContract(appChannel.Nonce, appChannel.ByteCode, appChannel.Constructor)
		if err != nil {
			log.Error("virtual contract not deployed")
			return err
		}
		appChannel.setDeployedAddr(deployedAddr)
	}
	return nil
}

func (c *AppClient) deployVirtualContract(
	nonce uint64, byteCode []byte, constructor []byte) (ctype.Addr, error) {

	virtResolverContract := c.nodeConfig.GetVirtResolverContract()
	virtAddr := GetVirtualAddress(byteCode, constructor, nonce)
	deployed, addr, err := c.isDeployed(virtAddr)
	if err != nil {
		log.Error(err)
		return ctype.ZeroAddr, err
	}
	if deployed {
		return addr, nil
	}
	log.Debugln("deploying virtual contract...")
	codeWithCons := append(byteCode, constructor...)

	receipt, err := c.transactorPool.SubmitWaitMined(
		"deploy virtual contract",
		func(transactor bind.ContractTransactor, opts *bind.TransactOpts) (*types.Transaction, error) {
			contract, err2 :=
				virtresolver.NewVirtContractResolverTransactor(virtResolverContract.GetAddr(), transactor)
			if err2 != nil {
				return nil, err2
			}
			return contract.Deploy(opts, codeWithCons, new(big.Int).SetUint64(nonce))
		},
		config.QuickTransactOptions(eth.WithGasLimit(4000000))...)
	if err != nil {
		log.Errorf("deploy virtual contract tx %x error %s", receipt.TxHash, err)
		return ctype.ZeroAddr, err
	}
	deployed, addr, err = c.isDeployed(virtAddr)
	if err != nil {
		return ctype.ZeroAddr, err
	}
	if deployed {
		log.Debugln("deployed virtual contract at", ctype.Addr2Hex(addr))
		return addr, nil
	}
	return ctype.ZeroAddr, fmt.Errorf("virtual contract not deployed")
}

// isDeployed checks if the given virtual address has been deployed on-chain;
// if yes, also returns the deployment address.
func (c *AppClient) isDeployed(virtAddr []byte) (bool, ctype.Addr, error) {
	contract, err := virtresolver.NewVirtContractResolverCaller(
		c.nodeConfig.GetVirtResolverContract().GetAddr(), c.transactorPool.ContractCaller())
	if err != nil {
		return false, ctype.ZeroAddr, err
	}
	var virt [32]byte
	copy(virt[:], virtAddr[:])
	deployedAddr, err := contract.Resolve(&bind.CallOpts{}, virt)
	if deployedAddr == (ctype.ZeroAddr) {
		return false, deployedAddr, nil
	}
	return true, deployedAddr, nil
}

// GetVirtualAddress derives the deterministic virtual-contract address from
// `(bytecode, constructor, nonce)`. Used both at registration time (to compute
// the session id) and at deploy time (to look up the eventual on-chain address
// in the virt-resolver).
func GetVirtualAddress(byteCode []byte, constructor []byte, nonce uint64) []byte {
	codeWithCons := append(byteCode, constructor...)
	return crypto.Keccak256(codeWithCons, utils.Pad(new(big.Int).SetUint64(nonce).Bytes(), 32))
}
