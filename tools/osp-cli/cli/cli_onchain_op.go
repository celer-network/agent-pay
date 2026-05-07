// Copyright 2018-2025 Celer Network

package cli

import (
	"fmt"
	"math/big"

	"github.com/celer-network/agent-pay/chain/channel-eth-go/nativewrap"
	"github.com/celer-network/agent-pay/config"
	"github.com/celer-network/agent-pay/ctype"
	"github.com/celer-network/agent-pay/route/routerregistry"
	"github.com/celer-network/agent-pay/utils"
	"github.com/celer-network/goutils/eth"
	"github.com/celer-network/goutils/log"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
)

// NativeWrapDeposit wraps native (e.g., ETH) into the chain's
// canonical wrapped-native (WETH-style) contract under the OSP's own
// balance, then approves CelerLedger to transferFrom that wrapped balance
// — the funding-flow shape CelerLedger expects when the OSP is the
// non-msgValueReceiver peer of an open-channel call.
func (p *Processor) NativeWrapDeposit() {
	if err := p.depositNativeWrap(); err != nil {
		return
	}
	if err := p.approveNativeWrapToLedger(); err != nil {
		return
	}
	p.queryNativeWrapLedgerAllowance()
}

// NativeWrapWithdraw unwraps the OSP's wrapped-native balance back to native.
func (p *Processor) NativeWrapWithdraw() {
	if err := p.withdrawNativeWrap(); err != nil {
		return
	}
}

func (p *Processor) RegisterRouter() {
	// check router registration
	ts, err := p.queryRouterRegistry()
	if err != nil {
		return
	}
	// registry router
	if ts == 0 {
		err = p.registerRouter()
		if err != nil {
			return
		}
		p.queryRouterRegistry()
	}
	log.Infoln("Welcome to Celer Network!")
}

func (p *Processor) DeregisterRouter() {
	// check router registration
	ts, err := p.queryRouterRegistry()
	if err != nil {
		return
	}
	// registry router
	if ts == 0 {
		log.Info("OSP not registered as a network router")
		return
	}
	p.deregisterRouter()
}

func (p *Processor) depositNativeWrap() error {
	log.Infof("wrap %f native into NativeWrap and wait transaction to be mined...", *amount)
	amtWei := utils.Float2Wei(*amount)
	nativeWrapAddr := ctype.Hex2Addr(p.profile.NativeWrapAddr)

	receipt, err := p.transactor.TransactWaitMined(
		"native-wrap deposit",
		func(transactor bind.ContractTransactor, opts *bind.TransactOpts) (*types.Transaction, error) {
			contract, err2 := nativewrap.NewNativeWrapTransactor(nativeWrapAddr, transactor)
			if err2 != nil {
				return nil, err2
			}
			// WETH.deposit() credits msg.sender; the OSP self-wraps.
			return contract.Deposit(opts)
		},
		config.TransactOptions(eth.WithEthValue(amtWei))...)
	if err != nil {
		log.Error(err)
		return err
	}
	if receipt.Status != types.ReceiptStatusSuccessful {
		return fmt.Errorf("native-wrap deposit transaction %x failed", receipt.TxHash)
	}
	return nil
}

func (p *Processor) approveNativeWrapToLedger() error {
	log.Info("approve NativeWrap balance to CelerLedger and wait transaction to be mined...")
	balance, err := p.queryNativeWrapBalance()
	if err != nil {
		return err
	}
	nativeWrapAddr := ctype.Hex2Addr(p.profile.NativeWrapAddr)
	ledgerAddr := ctype.Hex2Addr(p.profile.LedgerAddr)

	receipt, err := p.transactor.TransactWaitMined(
		"native-wrap approve",
		func(transactor bind.ContractTransactor, opts *bind.TransactOpts) (*types.Transaction, error) {
			contract, err2 := nativewrap.NewNativeWrapTransactor(nativeWrapAddr, transactor)
			if err2 != nil {
				return nil, err2
			}
			return contract.Approve(opts, ledgerAddr, balance)
		},
		config.TransactOptions()...)
	if err != nil {
		log.Error(err)
		return err
	}
	if receipt.Status != types.ReceiptStatusSuccessful {
		return fmt.Errorf("native-wrap approve transaction %x failed", receipt.TxHash)
	}
	return nil
}

func (p *Processor) queryNativeWrapBalance() (*big.Int, error) {
	nativeWrapAddr := ctype.Hex2Addr(p.profile.NativeWrapAddr)
	contract, err := nativewrap.NewNativeWrapCaller(nativeWrapAddr, p.transactor.ContractCaller())
	if err != nil {
		log.Error(err)
		return nil, err
	}
	balance, err := contract.BalanceOf(&bind.CallOpts{}, p.myAddr)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	log.Infoln("my balance at NativeWrap:", balance)
	return balance, nil
}

func (p *Processor) queryNativeWrapLedgerAllowance() (*big.Int, error) {
	nativeWrapAddr := ctype.Hex2Addr(p.profile.NativeWrapAddr)
	ledgerAddr := ctype.Hex2Addr(p.profile.LedgerAddr)
	contract, err := nativewrap.NewNativeWrapCaller(nativeWrapAddr, p.transactor.ContractCaller())
	if err != nil {
		log.Error(err)
		return nil, err
	}
	allowance, err := contract.Allowance(&bind.CallOpts{}, p.myAddr, ledgerAddr)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	log.Infoln("allowance from NativeWrap to Ledger is:", allowance)
	return allowance, nil
}

func (p *Processor) withdrawNativeWrap() error {
	log.Infof("unwrap %f from NativeWrap to native and wait transaction to be mined...", *amount)
	amtWei := utils.Float2Wei(*amount)
	nativeWrapAddr := ctype.Hex2Addr(p.profile.NativeWrapAddr)

	receipt, err := p.transactor.TransactWaitMined(
		"native-wrap withdraw",
		func(transactor bind.ContractTransactor, opts *bind.TransactOpts) (*types.Transaction, error) {
			contract, err2 := nativewrap.NewNativeWrapTransactor(nativeWrapAddr, transactor)
			if err2 != nil {
				return nil, err2
			}
			return contract.Withdraw(opts, amtWei)
		},
		config.TransactOptions()...)
	if err != nil {
		log.Error(err)
		return err
	}
	if receipt.Status != types.ReceiptStatusSuccessful {
		return fmt.Errorf("native-wrap withdraw transaction %x failed", receipt.TxHash)
	}
	return nil
}

func (p *Processor) queryRouterRegistry() (uint64, error) {
	routerRegistryAddr := ctype.Hex2Addr(p.profile.RouterRegistryAddr)
	contract, err := routerregistry.NewRouterRegistryCaller(routerRegistryAddr, p.transactor.ContractCaller())
	if err != nil {
		log.Error(err)
		return 0, err
	}
	info, err := contract.RouterInfo(&bind.CallOpts{}, p.myAddr)
	if err != nil {
		log.Error(err)
		return 0, err
	}
	ts := info.Uint64()
	if ts != 0 {
		log.Infoln("router registered / refreshed at unix-ts", ts)
	}
	return ts, nil
}

func (p *Processor) registerRouter() error {
	log.Info("register OSP as state channel router and wait transaction to be mined...")
	routerRegistryAddr := ctype.Hex2Addr(p.profile.RouterRegistryAddr)

	receipt, err := p.transactor.TransactWaitMined(
		"register router",
		func(transactor bind.ContractTransactor, opts *bind.TransactOpts) (*types.Transaction, error) {
			contract, err2 := routerregistry.NewRouterRegistryTransactor(routerRegistryAddr, transactor)
			if err2 != nil {
				return nil, err2
			}
			return contract.RegisterRouter(opts)
		},
		config.TransactOptions()...)
	if err != nil {
		log.Error(err)
		return err
	}
	if receipt.Status != types.ReceiptStatusSuccessful {
		return fmt.Errorf("register router transaction %x failed", receipt.TxHash)
	}
	return nil
}

func (p *Processor) deregisterRouter() error {
	log.Info("deregister OSP as state channel router and wait transaction to be mined...")
	routerRegistryAddr := ctype.Hex2Addr(p.profile.RouterRegistryAddr)

	receipt, err := p.transactor.TransactWaitMined(
		"deregister router",
		func(transactor bind.ContractTransactor, opts *bind.TransactOpts) (*types.Transaction, error) {
			contract, err2 := routerregistry.NewRouterRegistryTransactor(routerRegistryAddr, transactor)
			if err2 != nil {
				return nil, err2
			}
			return contract.DeregisterRouter(opts)
		},
		config.TransactOptions()...)
	if err != nil {
		log.Error(err)
		return err
	}
	if receipt.Status != types.ReceiptStatusSuccessful {
		return fmt.Errorf("deregister router transaction %x failed", receipt.TxHash)
	}
	return nil
}
