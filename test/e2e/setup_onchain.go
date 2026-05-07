// Copyright 2018-2025 Celer Network

package e2e

import (
	"context"
	"flag"
	"math/big"
	"time"

	"github.com/celer-network/agent-pay/chain"
	"github.com/celer-network/agent-pay/chain/channel-eth-go/deploy"
	"github.com/celer-network/agent-pay/chain/channel-eth-go/ledger"
	"github.com/celer-network/agent-pay/chain/channel-eth-go/nativewrap"
	"github.com/celer-network/agent-pay/common"
	"github.com/celer-network/agent-pay/ctype"
	tf "github.com/celer-network/agent-pay/testing"
	"github.com/celer-network/agent-pay/testing/testapp"
	"github.com/celer-network/goutils/eth"
	"github.com/celer-network/goutils/log"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

var conclient *ethclient.Client
var etherBaseAuth *bind.TransactOpts
var channelAddrBundle deploy.CelerChannelAddrBundle
var nativeWrapContract *nativewrap.NativeWrap
var erc20Contract *chain.ERC20
var autoFund bool
var onchainChainID *big.Int

var grpAddrs = [][]string{
	[]string{ospEthAddr, depositorEthAddr, osp2EthAddr, osp3EthAddr, osp4EthAddr, osp5EthAddr},
	[]string{osp6EthAddr, osp7EthAddr},
	[]string{osp8EthAddr, osp9EthAddr},
}
var grpPrivs = [][]string{
	[]string{osp1Priv, depositorPriv, osp2Priv, osp3Priv, osp4Priv, osp5Priv},
	[]string{osp6Priv, osp7Priv},
	[]string{osp8Priv, osp9Priv},
}

// SetupOnChain deploy contracts, and set limit etc
// return profile, tokenAddrErc20 and set testapp related addr
func SetupOnChain(appMap map[string]ctype.Addr, groupId uint64, autofund bool) (*common.ProfileJSON, string) {
	flag.Parse()
	autoFund = autofund
	var err error
	conclient, err = ethclient.Dial(outRootDir + "chaindata/geth.ipc")
	if err != nil {
		log.Fatalf("Failed to conclientect to the Ethereum: %v", err)
	}
	ethbasePrivKey, _ := crypto.HexToECDSA(etherBasePriv)
	etherBaseAuth = bind.NewKeyedTransactor(ethbasePrivKey)
	price := big.NewInt(2e9) // 2Gwei
	etherBaseAuth.GasPrice = price

	ctx := context.Background()
	onchainChainID, err = conclient.NetworkID(ctx)
	if err != nil {
		log.Fatalf("Failed to get Ethereum network id: %v", err)
	}
	etherBaseAuth, err = bind.NewKeyedTransactorWithChainID(ethbasePrivKey, onchainChainID)
	if err != nil {
		log.Fatalf("Failed to create keyed transactor: %v", err)
	}
	etherBaseAuth.GasPrice = price
	// deploy celer channel contracts
	channelAddrBundle = deploy.DeployAll(etherBaseAuth, conclient, ctx, 0)
	// deploy router registry
	routerRegistryAddr := deploy.DeployRouterRegistry(ctx, etherBaseAuth, conclient, 0)

	// NativeWrap (WETH-style) is used later when each OSP wraps its own
	// native balance and pre-approves CelerLedger so the open-channel /
	// deposit funding-flow path can pull pre-wrapped native via
	// `WETH.transferFrom` + `WETH.withdraw` for the non-msgValueReceiver peer.
	nativeWrapContract, err = nativewrap.NewNativeWrap(channelAddrBundle.NativeWrapAddr, conclient)
	if err != nil {
		log.Fatal(err)
	}

	// Disable channel deposit limit
	ledgerContract, err := ledger.NewCelerLedger(channelAddrBundle.CelerLedgerAddr, conclient)
	if err != nil {
		log.Fatal(err)
	}
	tx1, err := ledgerContract.DisableBalanceLimits(etherBaseAuth)
	if err != nil {
		log.Fatalf("Failed disable channel deposit limits: %v", err)
	}
	// Wait for tx1 to be mined to avoid nonce race on subsequent transactions
	receipt, err := eth.WaitMined(ctx, conclient, tx1, eth.WithPollingInterval(time.Second))
	if err != nil {
		log.Fatal(err)
	}
	chkTxStatus(receipt.Status, "Disable balance limit")

	// Deploy sample ERC20 contract (MOON)
	var erc20Addr ctype.Addr
	var tx2 *ethtypes.Transaction
	erc20Addr, tx2, erc20Contract, err = chain.DeployERC20(etherBaseAuth, conclient)
	if err != nil {
		log.Fatalf("Failed to deploy ERC20: %v", err)
	}
	receipt, err = eth.WaitMined(ctx, conclient, tx2, eth.WithPollingInterval(time.Second))
	if err != nil {
		log.Fatal(err)
	}
	chkTxStatus(receipt.Status, "Deploy ERC20 "+ctype.Addr2Hex(erc20Addr))

	// Deploy BooleanCondMock — the on-chain IBooleanCond used by both
	// VIRTUAL_CONTRACT and DEPLOYED_CONTRACT dispute scenarios.
	appAddr1, tx3, _, err := testapp.DeployBooleanCondMock(etherBaseAuth, conclient)
	if err != nil {
		log.Fatalf("Failed to deploy BooleanCondMock contract: %v", err)
	}
	receipt, err = eth.WaitMined(ctx, conclient, tx3, eth.WithPollingInterval(time.Second))
	if err != nil {
		log.Fatal(err)
	}
	chkTxStatus(receipt.Status, "Deploy BooleanCondMock "+ctype.Addr2Hex(appAddr1))
	appMap["BooleanCondMock"] = appAddr1

	// Deploy a new Celer Ledger for channel migration test
	log.Infoln("Deploying new CelerLedger contract...")
	newLedgerAddr, tx6, _, err := deploy.DeployContractWithLinks(
		etherBaseAuth,
		conclient,
		ledger.CelerLedgerABI,
		ledger.CelerLedgerBin,
		map[string]ctype.Addr{
			"LedgerStruct":       channelAddrBundle.LedgerStructAddr,
			"LedgerOperation":    channelAddrBundle.OperationAddr,
			"LedgerChannel":      channelAddrBundle.LedgerChannelAddr,
			"LedgerBalanceLimit": channelAddrBundle.BalanceLimitAddr,
			"LedgerMigrate":      channelAddrBundle.MigrateAddr,
		},
		channelAddrBundle.NativeWrapAddr,
		channelAddrBundle.PayRegistryAddr,
		channelAddrBundle.CelerWalletAddr,
	)
	if err != nil {
		log.Fatalf("Failed to deploy new CelerLedger contract: %v", err)
	}
	receipt, err = eth.WaitMined(ctx, conclient, tx6, eth.WithPollingInterval(time.Second))
	if err != nil {
		log.Fatalf("Failed to WaitMined v2 CelerLedger: %v", err)
	}
	chkTxStatus(receipt.Status, "Deploy v2 Ledger contract at "+ctype.Addr2Hex(newLedgerAddr))

	log.Infoln("Add fund to OSP accounts ...")
	if groupId < 3 {
		fundEthAddrs(grpAddrs[groupId], grpPrivs[groupId])
	}

	// contruct ledger map
	ledgers := map[string]string{
		ctype.Addr2Hex(channelAddrBundle.CelerLedgerAddr): "ledger1",
		ctype.Addr2Hex(newLedgerAddr):                     "ledger2",
	}

	profileContracts := common.ProfileContracts{
		Wallet:         ctype.Addr2Hex(channelAddrBundle.CelerWalletAddr),
		Ledger:         ctype.Addr2Hex(channelAddrBundle.CelerLedgerAddr),
		VirtResolver:   ctype.Addr2Hex(channelAddrBundle.VirtResolverAddr),
		NativeWrap:     ctype.Addr2Hex(channelAddrBundle.NativeWrapAddr),
		PayResolver:    ctype.Addr2Hex(channelAddrBundle.PayResolverAddr),
		PayRegistry:    ctype.Addr2Hex(channelAddrBundle.PayRegistryAddr),
		RouterRegistry: ctype.Addr2Hex(routerRegistryAddr),
		Ledgers:        ledgers,
	}

	profileEth := common.ProfileEthereum{
		Gateway:          ethGateway,
		ChainId:          1337,
		BlockIntervalSec: 1,
		BlockDelayNum:    0,
		DisputeTimeout:   10,
		Contracts:        profileContracts,
		CheckInterval: map[string]uint64{
			"CooperativeWithdraw": 2,
			"Deploy":              2,
			"Deposit":             2,
			"IntendSettle":        2,
			"OpenChannel":         2,
			"ConfirmSettle":       2,
			"IntendWithdraw":      2,
			"ConfirmWithdraw":     2,
			"RouterUpdated":       2,
			"MigrateChannelTo":    2,
		},
	}

	profileOsp := common.ProfileOsp{
		Host:    "localhost:10000",
		Address: ospEthAddr,
	}

	// output json file
	p := &common.ProfileJSON{
		Version:  "0.1",
		Ethereum: profileEth,
		Osp:      profileOsp,
	}
	return p, ctype.Addr2Hex(erc20Addr)
}

func fundEthAddr(addrStr, privKeyStr string) {
	addr := ctype.Hex2Addr(addrStr)
	err := tf.FundAddr("100000000000000000000", []*ctype.Addr{&addr})
	if err != nil {
		log.Fatalln("failed to fund addr", addrStr, err)
	}
	tx := fundEthAddrStep1(addrStr)
	fundEthAddrStep1Check(addrStr, tx)
	tx1, tx2, tx3 := fundEthAddrStep2(addrStr, privKeyStr)
	fundEthAddrStep2Check(addrStr, tx1, tx2, tx3)
}

func fundEthAddrs(addrStrs, privKeyStr []string) {
	var addrs []*ctype.Addr
	for _, addrStr := range addrStrs {
		addr := ctype.Hex2Addr(addrStr)
		addrs = append(addrs, &addr)
	}
	err := tf.FundAddr("1000000000000000000000000", addrs) // 1 million native
	if err != nil {
		log.Fatalln("failed to fund", err)
	}

	var step1Txs []*ethtypes.Transaction
	var step2Tx1s, step2Tx2s, step2Tx3s []*ethtypes.Transaction
	for i := range addrStrs {
		step1Txs = append(step1Txs, fundEthAddrStep1(addrStrs[i]))
	}
	for i := range addrStrs {
		fundEthAddrStep1Check(addrStrs[i], step1Txs[i])
	}
	if autoFund {
		for i := range addrStrs {
			tx1, tx2, tx3 := fundEthAddrStep2(addrStrs[i], privKeyStr[i])
			step2Tx1s = append(step2Tx1s, tx1)
			step2Tx2s = append(step2Tx2s, tx2)
			step2Tx3s = append(step2Tx3s, tx3)
		}
		for i := range addrStrs {
			fundEthAddrStep2Check(addrStrs[i], step2Tx1s[i], step2Tx2s[i], step2Tx3s[i])
		}
	}
}

func fundEthAddrStep1(addrStr string) *ethtypes.Transaction {
	addr := ctype.Hex2Addr(addrStr)
	moonAmt := new(big.Int)
	moonAmt.SetString("1000000000000000000000000000", 10) // 1 billion Moon
	tx, err := erc20Contract.Transfer(etherBaseAuth, addr, moonAmt)
	if err != nil {
		log.Fatalln("failed to send MOON token for", addrStr, err)
	}
	return tx
}

// fundEthAddrStep2 prepares each OSP account so it can act as either peer in
// an open-channel call:
//
//   - wraps the OSP's native balance into WETH (WETH.deposit credits
//     msg.sender, so the OSP must self-wrap — the contract's funding-flow
//     path requires WETH already in place when the OSP is the
//     non-msgValueReceiver peer);
//   - approves CelerLedger to transferFrom the OSP's WETH balance;
//   - approves CelerLedger to transferFrom the OSP's MOON ERC20 balance.
func fundEthAddrStep2(addrStr, privKeyStr string) (*ethtypes.Transaction, *ethtypes.Transaction, *ethtypes.Transaction) {
	privKey, err := crypto.HexToECDSA(privKeyStr)
	if err != nil {
		log.Fatalln("failed to get private key", addrStr, err)
	}
	auth, err := bind.NewKeyedTransactorWithChainID(privKey, onchainChainID)
	if err != nil {
		log.Fatalln("failed to create keyed transactor", addrStr, err)
	}
	auth.GasPrice = etherBaseAuth.GasPrice

	var tx1 *ethtypes.Transaction
	if autoFund {
		// Wrap most of the OSP's native balance (leave a comfortable
		// gas buffer; fundEthAddrs funded the OSP with 1M native).
		nativeAmt := new(big.Int)
		nativeAmt.SetString("999000000000000000000000", 10) // 999_000 native
		auth.Value = nativeAmt
		tx1, err = nativeWrapContract.Deposit(auth)
		if err != nil {
			log.Fatalln("failed to wrap native into WETH for", addrStr, err)
		}
		auth.Value = big.NewInt(0)
	}

	wrapAmt := new(big.Int)
	wrapAmt.SetString("999000000000000000000000", 10) // 999_000 wrapped native
	tx2, err := nativeWrapContract.Approve(auth, channelAddrBundle.CelerLedgerAddr, wrapAmt)
	if err != nil {
		log.Fatalln("failed to approve native-wrap to celerLedger for", addrStr, err)
	}

	moonAmt := new(big.Int)
	moonAmt.SetString("1000000000000000000000000000", 10) // 1 billion Moon
	tx3, err := erc20Contract.Approve(auth, channelAddrBundle.CelerLedgerAddr, moonAmt)
	if err != nil {
		log.Fatalln("failed to approve MOON to celerLedger for", addrStr, err)
	}

	return tx1, tx2, tx3
}

func fundEthAddrStep1Check(addrStr string, tx *ethtypes.Transaction) {
	ctx := context.Background()
	receipt, err := eth.WaitMined(ctx, conclient, tx, eth.WithPollingInterval(time.Second))
	if err != nil {
		log.Fatalln("wait mined failed", addrStr, err)
	}
	chkTxStatus(receipt.Status, "transfer moon token to "+addrStr)
}

func fundEthAddrStep2Check(addrStr string, tx1, tx2, tx3 *ethtypes.Transaction) {
	ctx := context.Background()
	if autoFund {
		receipt, err := eth.WaitMined(ctx, conclient, tx1, eth.WithPollingInterval(time.Second))
		if err != nil {
			log.Fatalln("wait mined failed", addrStr, err)
		}
		chkTxStatus(receipt.Status, addrStr+" wrap native into WETH")
	}

	receipt, err := eth.WaitMined(ctx, conclient, tx2, eth.WithPollingInterval(time.Second))
	if err != nil {
		log.Fatalln("wait mined failed", addrStr, err)
	}
	chkTxStatus(receipt.Status, addrStr+" approve native-wrap to ledger")

	receipt, err = eth.WaitMined(ctx, conclient, tx3, eth.WithPollingInterval(time.Second))
	if err != nil {
		log.Fatalln("wait mined failed", addrStr, err)
	}
	chkTxStatus(receipt.Status, addrStr+" approve moon token")
	log.Infoln("finish funding for", addrStr)
}

// if status isn't 1 (sucess), log.Fatal
func chkTxStatus(s uint64, txname string) {
	if s != 1 {
		log.Fatal(txname + " tx failed")
	}
	log.Info(txname + " tx success")
}
