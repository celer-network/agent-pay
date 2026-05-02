// Copyright 2019-2025 Celer Network

package route

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"time"

	rt "github.com/celer-network/agent-pay/chain/channel-eth-go/routerregistry"
	"github.com/celer-network/agent-pay/common"
	"github.com/celer-network/agent-pay/common/event"
	"github.com/celer-network/agent-pay/common/intfs"
	"github.com/celer-network/agent-pay/config"
	"github.com/celer-network/agent-pay/ctype"
	"github.com/celer-network/agent-pay/ledgerview"
	"github.com/celer-network/agent-pay/route/ospreport"
	"github.com/celer-network/agent-pay/rpc"
	"github.com/celer-network/agent-pay/rtconfig"
	"github.com/celer-network/agent-pay/storage"
	"github.com/celer-network/agent-pay/utils"
	"github.com/celer-network/goutils/eth"
	"github.com/celer-network/goutils/eth/monitor"
	"github.com/celer-network/goutils/log"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"google.golang.org/protobuf/proto"
)

type BcastSendCallback func(info *rpc.RoutingRequest, ospAddrs []string)

// Controller configs to handle onchain router-related event
type Controller struct {
	nodeConfig        common.GlobalNodeConfig
	transactor        *eth.Transactor
	monitorService    intfs.MonitorService
	ethclient         *ethclient.Client
	dal               *storage.DAL
	signer            eth.Signer
	bcastSendCallback BcastSendCallback
	rtBuilder         *routingTableBuilder
	explorerReport    *ospreport.OspInfo
	explorerUrl       string // explorer url

	// Dynamic routing updates from OSPs are gathered here then
	// used for recomputing the routing table.
	routingBatch     map[ctype.Addr]*rpc.RoutingUpdate
	routingBatchLock sync.Mutex
}

// Enum corrsponding to the onchain router operation
const (
	routerAdded uint8 = iota
	routerRemoved
	routerRefreshed
)

const (
	// A router OSP checks the router registry at startup and then every checkRegistryInterval to see
	// whether its registry time is older than refreshIntervalSec, and refreshes itself onchain if so.
	// It also scans the local rtBuilder.getAllOsps() each interval (first scan at
	// startupTime + checkRegistryInterval) and removes routers whose stored timestamp is older than
	// expireTimeoutSec. The contract stores `block.timestamp` (unix seconds) per router, so these
	// thresholds are likewise in seconds.
	checkRegistryInterval = 6 * time.Hour  // time interval to check for self-refresh and OSP timeouts
	refreshIntervalSec    = uint64(432000) // 5 days, seconds — refresh self if stored ts is older
	expireTimeoutSec      = uint64(604800) // 7 days, seconds — drop a router whose stored ts is older
	// backtrackSafetyMargin is multiplied with `expireTimeoutSec` when computing
	// the on-chain event-monitor backtrack window, so a chain whose actual block
	// time briefly accelerates can't push a still-live RouterUpdated event out
	// of the replay window.
	backtrackSafetyMargin = 2
	// minBacktrackBlocks lower-bounds the computed backtrack so we always scan
	// at least a sane number of blocks, even on very slow chains.
	minBacktrackBlocks = uint64(50000)
	// blockTimeSampleSpan is how many blocks back we sample to estimate the
	// chain's effective block time at startup. Big enough to absorb per-block
	// jitter, small enough to be one HeaderByNumber call.
	blockTimeSampleSpan = uint64(1000)

	routeTTL = 15
)

// NewController creates a new process for router controller
func NewController(
	nodeConfig common.GlobalNodeConfig,
	transactor *eth.Transactor,
	monitorService intfs.MonitorService,
	ethclient *ethclient.Client,
	dal *storage.DAL,
	signer eth.Signer,
	bcastSendCallback BcastSendCallback,
	routingData []byte,
	rpcHost string,
	explorerUrl string) (*Controller, error) {
	c := &Controller{
		nodeConfig:        nodeConfig,
		transactor:        transactor,
		monitorService:    monitorService,
		ethclient:         ethclient,
		dal:               dal,
		signer:            signer,
		bcastSendCallback: bcastSendCallback,
		explorerUrl:       explorerUrl,
	}
	c.rtBuilder = newRoutingTableBuilder(nodeConfig.GetOnChainAddr(), dal)
	if c.rtBuilder == nil {
		return c, fmt.Errorf("fail to initialize routing table builder")
	}
	c.explorerReport = &ospreport.OspInfo{
		EthAddr:    nodeConfig.GetOnChainAddr().Hex(), // format required by explorer
		RpcHost:    rpcHost,
		OpenAccept: true,
	}
	err := c.startRoutingRecoverProcess(monitorService.GetCurrentBlockNumber(), routingData, nodeConfig)
	return c, err
}

// Start starts router process to instantiate OSP as a router.
func (c *Controller) Start() {
	// check if OSP is registered on-chain as a router. Stored value is the unix
	// timestamp (seconds) of the most recent register/refresh — see RouterRegistry.sol.
	registeredAt, err := c.queryRouterRegistry()
	if err != nil {
		log.Errorf("query router registry failed: %s", err)
		return
	}
	if registeredAt != 0 {
		log.Infoln("router registered / refreshed at unix ts", registeredAt)
		// check if OSP needs to send refresh transaction
		nowTs := uint64(time.Now().Unix())
		if nowTs-registeredAt > refreshIntervalSec {
			c.refreshRouterRegistry()
		}
		// start onchain events monitor
		c.monitorRouterUpdatedEvent()
		// start routine job
		go c.runRoutersRoutineJob()
	} else {
		log.Warn("NOT able to join the OSP network because this node is not registered on-chain as a router")
	}
}

// monitors the RouterUpdated event onchain
// backtrack from one interval before the current block
func (c *Controller) monitorRouterUpdatedEvent() {
	monitorCfg := &monitor.Config{
		ChainId:       config.ChainId.Uint64(),
		EventName:     event.RouterUpdated,
		Contract:      c.nodeConfig.GetRouterRegistryContract(),
		StartBlock:    c.calculateStartBlockNumber(),
		Reset:         true,
		CheckInterval: c.nodeConfig.GetCheckInterval(event.RouterUpdated),
	}
	_, err := c.monitorService.Monitor(monitorCfg,
		func(id monitor.CallbackID, eLog types.Log) bool {
			e := &rt.RouterRegistryRouterUpdated{} // event RouterUpdated
			if err := c.nodeConfig.GetRouterRegistryContract().ParseEvent(event.RouterUpdated, eLog, e); err != nil {
				log.Error(err)
				return false
			}

			// Only used in log
			routerAddr := ctype.Addr2Hex(e.RouterAddress)
			txHash := fmt.Sprintf("%x", eLog.TxHash)
			log.Infoln("Seeing RouterUpdated event, router addr:", routerAddr, "tx hash:", txHash, "callback id:", id, "blkNum:", eLog.BlockNumber)

			// Use the event's block timestamp — that's the canonical
			// `block.timestamp` the contract stamped into RouterRegistry. On
			// startup replay this can be days old, so falling back to
			// time.Now() would silently revive stale routers.
			registeredAt, err := c.blockTimestamp(eLog.BlockNumber)
			if err != nil {
				log.Errorf("fetch block %d timestamp for RouterUpdated event: %s", eLog.BlockNumber, err)
				return false
			}
			c.processRouterUpdatedEvent(e, registeredAt)
			return false
		},
	)

	if err != nil {
		log.Error(err)
	}
}

// processes the RouterUpdated event according to various router opeartion
func (c *Controller) processRouterUpdatedEvent(e *rt.RouterRegistryRouterUpdated, registeredAt uint64) {
	switch e.Op {
	case routerAdded:
		c.addRouter(e.RouterAddress, registeredAt)
	case routerRemoved:
		c.removeRouter(e.RouterAddress)
	case routerRefreshed:
		c.refreshRouter(e.RouterAddress, registeredAt)
	default:
		log.Warnf("Unknown router operation from router registry contract: %v", e.Op)
	}
}

// adds router node and record the unix timestamp of register/refresh
func (c *Controller) addRouter(routerAddr ctype.Addr, registeredAt uint64) {
	c.rtBuilder.markOsp(routerAddr, registeredAt)
}

// removes router node and delete it from the map
func (c *Controller) removeRouter(routerAddr ctype.Addr) {
	if !c.rtBuilder.hasOsp(routerAddr) {
		return
	}
	c.rtBuilder.unmarkOsp(routerAddr)
}

// refreshes a router node and update its stored register/refresh unix timestamp.
func (c *Controller) refreshRouter(routerAddr ctype.Addr, registeredAt uint64) {
	c.rtBuilder.markOsp(routerAddr, registeredAt)
}

// calculateStartBlockNumber computes the on-chain event-monitor start block
// so a RouterUpdated event from anywhere in the live `expireTimeoutSec`
// window is still in range on startup. The block-count backtrack is derived
// from the chain's actual block time (sampled from headers) — a fixed block
// count would either under-cover a fast chain (e.g. 50k blocks ≈ 28h on a
// 2s chain, far less than the 7-day expiry) or waste replay work on a slow
// one.
func (c *Controller) calculateStartBlockNumber() *big.Int {
	currentBlk := c.monitorService.GetCurrentBlockNumber()
	backtrack := new(big.Int).SetUint64(c.computeBacktrackBlocks(currentBlk))
	if backtrack.Cmp(currentBlk) >= 0 {
		return big.NewInt(0)
	}
	return new(big.Int).Sub(currentBlk, backtrack)
}

// computeBacktrackBlocks samples the chain's effective block time from a pair
// of header timestamps (current vs ~`blockTimeSampleSpan` blocks back) and
// converts the live-router window into a block count. Falls back to
// `minBacktrackBlocks` if the sample fails or yields a degenerate block time.
func (c *Controller) computeBacktrackBlocks(currentBlk *big.Int) uint64 {
	window := expireTimeoutSec * backtrackSafetyMargin
	if c.ethclient == nil || currentBlk.Sign() <= 0 {
		return minBacktrackBlocks
	}
	span := new(big.Int).SetUint64(blockTimeSampleSpan)
	if span.Cmp(currentBlk) >= 0 {
		// Brand-new chain — the window-in-blocks would be hand-wavy; use the
		// floor and let the monitor scan from genesis if currentBlk is small.
		return minBacktrackBlocks
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	latest, err := c.ethclient.HeaderByNumber(ctx, currentBlk)
	if err != nil {
		log.Warnf("backtrack block-time sample (latest) failed: %s — using minBacktrackBlocks", err)
		return minBacktrackBlocks
	}
	older, err := c.ethclient.HeaderByNumber(ctx, new(big.Int).Sub(currentBlk, span))
	if err != nil {
		log.Warnf("backtrack block-time sample (older) failed: %s — using minBacktrackBlocks", err)
		return minBacktrackBlocks
	}
	if latest.Time <= older.Time {
		log.Warnf("backtrack block-time sample non-monotonic (latest %d, older %d) — using minBacktrackBlocks", latest.Time, older.Time)
		return minBacktrackBlocks
	}
	elapsed := latest.Time - older.Time
	// blocks needed = window * span / elapsed (avoids float division).
	backtrack := window * blockTimeSampleSpan / elapsed
	if backtrack < minBacktrackBlocks {
		return minBacktrackBlocks
	}
	return backtrack
}

// blockTimestamp returns the unix-second timestamp of the given block.
func (c *Controller) blockTimestamp(blockNumber uint64) (uint64, error) {
	if c.ethclient == nil {
		return 0, fmt.Errorf("ethclient not configured")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	header, err := c.ethclient.HeaderByNumber(ctx, new(big.Int).SetUint64(blockNumber))
	if err != nil {
		return 0, err
	}
	return header.Time, nil
}

// call routerInfo in router registry contract to check if Osp has been registered.
// Return value is the unix timestamp (seconds) of the most recent register/refresh
// for this Osp address — the contract stores `block.timestamp`, see RouterRegistry.sol.
func (c *Controller) queryRouterRegistry() (uint64, error) {
	routerRegistryAddr := c.nodeConfig.GetRouterRegistryContract().GetAddr()
	caller, err := rt.NewRouterRegistryCaller(routerRegistryAddr, c.transactor.ContractCaller())
	if err != nil {
		return 0, err
	}
	registeredAt, err := caller.RouterInfo(&bind.CallOpts{}, c.transactor.Address())
	if err != nil {
		return 0, err
	}
	return registeredAt.Uint64(), nil
}

func (c *Controller) checkAndRefreshIfNeeded() {
	registeredAt, err := c.queryRouterRegistry()
	if err != nil {
		log.Errorf("query router registry failed: %s", err)
		return
	}
	// `registeredAt == 0` means RouterRegistry has no entry for us — either
	// we were never registered, or we (or the operator) deregistered. Either
	// way, refreshing here would silently re-register us; respect the
	// deregistration and let an explicit register call bring us back.
	if registeredAt == 0 {
		log.Warn("OSP not registered on-chain; skipping periodic refresh")
		return
	}
	nowTs := uint64(time.Now().Unix())
	if nowTs-registeredAt > refreshIntervalSec {
		c.refreshRouterRegistry()
	}
}

// send on-chain transaction to refresh the OSP address's last-seen unix
// timestamp in the RouterRegistry.
// CAUTION: need to pay attention if it fails to refresh
func (c *Controller) refreshRouterRegistry() {
	log.Infoln("sending RefreshRouter tx")
	routerRegistryAddr := c.nodeConfig.GetRouterRegistryContract().GetAddr()
	_, err := c.transactor.Transact(
		&eth.TransactionStateHandler{
			OnMined: func(receipt *types.Receipt) {
				if receipt.Status == types.ReceiptStatusSuccessful {
					log.Infof("RefreshRouter transaction %x succeeded", receipt.TxHash)
				} else {
					log.Errorf("RefreshRouter transaction %x failed", receipt.TxHash)
				}
			},
		},
		func(transactor bind.ContractTransactor, opts *bind.TransactOpts) (*types.Transaction, error) {
			contract, err2 := rt.NewRouterRegistryTransactor(routerRegistryAddr, transactor)
			if err2 != nil {
				log.Errorln("NewRouterRegistryTransactor err:", err2)
				return nil, err2
			}
			return contract.RefreshRouter(opts)
		},
		config.TransactOptions()...)
	if err != nil {
		log.Errorf("Fail to refresh the router: %s", err)
	}
}

// starts some routine jobs
// CAUTION: This should be run in goroutine
func (c *Controller) runRoutersRoutineJob() {
	checkTicker := time.NewTicker(checkRegistryInterval)
	bcastTicker := time.NewTicker(config.RouterBcastInterval)
	buildTicker := time.NewTicker(config.RouterBuildInterval)
	reportTicker := time.NewTicker(config.OspReportInverval)
	defer func() {
		checkTicker.Stop()
		bcastTicker.Stop()
		buildTicker.Stop()
		reportTicker.Stop()
	}()

	for {
		select {
		case <-checkTicker.C:
			c.checkAndRefreshIfNeeded()
			c.removeExpiredRouters()
		case <-bcastTicker.C:
			c.bcastRouterInfo()
		case <-buildTicker.C:
			c.buildRoutingTable()
		case <-reportTicker.C:
			c.reportOspInfoToExplorer()
		}
	}
}

// Traverses the map and remove the expired routers.
func (c *Controller) removeExpiredRouters() {
	nowTs := uint64(time.Now().Unix())
	ospInfo := c.rtBuilder.getAllOsps()
	for addr := range ospInfo {
		registeredAt := ospInfo[addr].RegistryTime

		if isRouterExpired(registeredAt, nowTs) {
			c.rtBuilder.unmarkOsp(addr)
		}
	}
}

func isRouterExpired(registeredAt, nowTs uint64) bool {
	return registeredAt+expireTimeoutSec < nowTs
}

// Get my dynamic routing information and broadcast it to peer OSPs.
// Also enqueue to the routing info batch to include it in the next
// routing recomputation.
func (c *Controller) bcastRouterInfo() {
	channels := c.gatherChannelInfo()
	myAddr := ctype.Addr2Hex(c.nodeConfig.GetOnChainAddr())
	update := &rpc.RoutingUpdate{
		Origin:   myAddr,
		Ts:       uint64(now().Unix()),
		Channels: channels,
	}

	updateBytes, err := proto.Marshal(update)
	if err != nil {
		log.Errorln("proto marshal signedUpdate err", err, update)
		return
	}
	sig, err := c.signer.SignEthMessage(updateBytes)

	signedUpdate := &rpc.SignedRoutingUpdate{
		Update: updateBytes,
		Sig:    sig,
		Ttl:    routeTTL,
	}

	info := &rpc.RoutingRequest{
		Updates: []*rpc.SignedRoutingUpdate{signedUpdate},
	}

	c.enqueueRouterInfo(update, signedUpdate.GetTtl())

	c.bcast(info, []string{myAddr}, "")
}

func (c *Controller) gatherChannelInfo() []*rpc.ChannelRoutingInfo {
	var channels []*rpc.ChannelRoutingInfo
	nowTs := uint64(time.Now().Unix())
	for _, neighbor := range c.rtBuilder.getAliveNeighbors() {
		for _, cid := range neighbor.TokenCids {
			bal, err := ledgerview.GetBalance(c.dal, cid, c.nodeConfig.GetOnChainAddr(), nowTs)
			if err != nil {
				log.Error(err)
				continue
			}
			channel := &rpc.ChannelRoutingInfo{
				Cid:     ctype.Cid2Hex(cid),
				Balance: bal.MyFree.String(),
			}
			channels = append(channels, channel)
		}
	}
	return channels
}

// Enqueue the dynamic routing information and return true if it should be
// propagated to peer OSPs in the broadcast.  The information is propagated
// if it's new to this OSP and still has time-to-live (hop counter).
//
// If this is the first information in the batch, trigger a delayed action
// to recompute the routing table, giving it some time for more routing info
// to be added to the batch.
func (c *Controller) enqueueRouterInfo(update *rpc.RoutingUpdate, ttl uint64) bool {
	if update == nil {
		return false
	}

	origin := ctype.Hex2Addr(update.GetOrigin())
	ts := update.GetTs()
	if ttl <= 0 {
		return false // this should not happen
	}

	c.routingBatchLock.Lock()
	defer c.routingBatchLock.Unlock()

	if c.routingBatch == nil {
		c.routingBatch = make(map[ctype.Addr]*rpc.RoutingUpdate)
	}

	oldUpdate, ok := c.routingBatch[origin]
	if ok && oldUpdate.GetTs() >= ts {
		return false // already have newer info from this origin
	}
	// keep osp and edges alive
	timestamp := time.Unix(int64(ts), 0).UTC()
	now := now()
	if timestamp.After(now) {
		timestamp = now
	}
	c.rtBuilder.keepOspAlive(origin, timestamp)
	for _, ch := range update.GetChannels() {
		if ch != nil {
			balance := utils.Wei2BigInt(ch.GetBalance())
			if balance == nil {
				log.Errorln("invalid balance report", ch.GetBalance())
				continue
			}
			c.rtBuilder.updateOspEdge(ctype.Hex2Cid(ch.GetCid()), balance, origin, timestamp)
		}
	}

	c.routingBatch[origin] = update
	// Propagate the info if the incoming TTL was more than 1.
	return (ttl > 1)
}

func (c *Controller) buildRoutingTable() {
	for token := range c.rtBuilder.getAllTokens() {
		c.rtBuilder.buildTable(token)
	}

	// TODO: Recompute the routing table according to bcast info.
	/*
		// Grab the current batch of routing info and release the lock.
		c.routingBatchLock.Lock()
		batch := c.routingBatch
		c.routingBatch = nil
		c.routingBatchLock.Unlock()
		log.Debugf("computing routing table from %d OSP info", len(batch))
	*/
}

// New routing information arrived from another OSP. Enqueue it for
// a future route recomputation and, if needed, forward it to other
// peer OSPs in the broadcast.
func (c *Controller) RecvBcastRoutingInfo(info *rpc.RoutingRequest) error {
	// TODO: support batch updates
	if len(info.GetUpdates()) != 1 {
		return fmt.Errorf("invalid number of routing updates in one request, %d", len(info.GetUpdates()))
	}
	signedUpdate := info.Updates[0]
	var update rpc.RoutingUpdate
	err := proto.Unmarshal(signedUpdate.GetUpdate(), &update)
	if err != nil {
		return fmt.Errorf("unmarshal signed update err: %w", err)
	}
	if !eth.IsSignatureValid(ctype.Hex2Addr(update.GetOrigin()), signedUpdate.GetUpdate(), signedUpdate.GetSig()) {
		return fmt.Errorf("route update invalid sig for origin %s", update.GetOrigin())
	}

	log.Debugln("receive router update:", &update)
	c.rtBuilder.keepNeighborAlive(ctype.Hex2Addr(info.GetSender()))
	if c.enqueueRouterInfo(&update, signedUpdate.GetTtl()) {
		info.Updates[0].Ttl--
		c.bcast(info, []string{update.GetOrigin()}, info.GetSender())
	}

	return nil
}

// Send out the given routing information request to the peer OSPs
// excluding the direct sender of this message (if any).
func (c *Controller) bcast(info *rpc.RoutingRequest, origins []string, sender string) {
	// Get peer OSPs excluding me and the given direct sender.
	var ospAddrs []string
	// TODO: support batch updates
	origin := origins[0]
	myAddr := ctype.Addr2Hex(c.nodeConfig.GetOnChainAddr())
	neighborAddrs := c.rtBuilder.getNeighborAddrs()
	for _, ospAddr := range neighborAddrs {
		ospAddrStr := ctype.Addr2Hex(ospAddr)
		if ospAddrStr == myAddr || ospAddrStr == sender {
			continue
		}
		if origin == ospAddrStr {
			continue
		}
		ospAddrs = append(ospAddrs, ospAddrStr)
	}
	if len(ospAddrs) == 0 {
		return
	}
	info.Sender = myAddr
	log.Debugf("bcast router updates: origin %s, to %s", origin, ospAddrs)
	c.bcastSendCallback(info, ospAddrs)
}

func (c *Controller) reportOspInfoToExplorer() {
	if c.explorerUrl == "" {
		return
	}
	// set osp peers
	c.explorerReport.OspPeers = nil
	nowTs := uint64(time.Now().Unix())
	for addr, neighbor := range c.rtBuilder.getAliveNeighbors() {
		peerBalances := &ospreport.PeerBalances{
			Peer: addr.Hex(), // format required by explorer
		}
		for tk, cid := range neighbor.TokenCids {
			bal, err := ledgerview.GetBalance(c.dal, cid, c.nodeConfig.GetOnChainAddr(), nowTs)
			if err != nil {
				log.Error(err)
				continue
			}
			peerBalances.Balances = append(
				peerBalances.Balances,
				&ospreport.ChannelBalance{
					Cid:         cid.Hex(), // format required by explorer
					TokenAddr:   tk.Hex(),  // format required by explorer
					SelfBalance: bal.MyFree.String(),
					PeerBalance: bal.PeerFree.String(),
				})
		}
		c.explorerReport.OspPeers = append(c.explorerReport.OspPeers, peerBalances)
	}
	// set std openchan configs
	c.explorerReport.StdOpenchanConfigs = nil
	for token, cfg := range rtconfig.GetStandardConfigs().GetConfig() {
		if cfg != nil {
			cfgReport := &ospreport.StdOpenChanConfig{
				TokenAddr:  ctype.Hex2Addr(token).Hex(), // format required by explorer
				MinDeposit: cfg.MinDeposit,
				MaxDeposit: cfg.MaxDeposit,
			}
			c.explorerReport.StdOpenchanConfigs = append(c.explorerReport.StdOpenchanConfigs, cfgReport)
		}
	}
	// set pay count
	payCount, err := c.dal.CountPayments()
	if err != nil {
		log.Error("CountPayments err:", err)
	}
	c.explorerReport.Payments = int64(payCount)
	// set timestamp
	c.explorerReport.Timestamp = uint64(now().Unix())
	// marshal and sign
	reportBytes, err := proto.Marshal(c.explorerReport)
	if err != nil {
		log.Errorln("proto marshal OSP report err:", err, c.explorerReport)
		return
	}
	sig, err := c.signer.SignEthMessage(reportBytes)
	if err != nil {
		log.Error(err)
		return
	}
	// send report
	report := map[string]string{
		"ospInfo": ctype.Bytes2Hex(reportBytes),
		"sig":     ctype.Bytes2Hex(sig),
	}
	_, err = utils.HttpPost(c.explorerUrl, report)
	if err != nil {
		log.Warnln("explorer report error:", err)
	}
}

func (c *Controller) BuildTable(tokenAddr ctype.Addr) (map[ctype.Addr]ctype.CidType, error) {
	return c.rtBuilder.buildTable(tokenAddr)
}

func (c *Controller) AddEdge(p1, p2 ctype.Addr, cid ctype.CidType, tokenAddr ctype.Addr) error {
	return c.rtBuilder.addEdge(p1, p2, cid, tokenAddr)
}

func (c *Controller) RemoveEdge(cid ctype.CidType) error {
	return c.rtBuilder.removeEdge(cid)
}

func (c *Controller) GetAllNeighbors() map[ctype.Addr]*NeighborInfo {
	return c.rtBuilder.getAllNeighbors()
}

func now() time.Time {
	return time.Now().UTC()
}
