// Copyright 2018-2025 Celer Network

package common

import (
	"encoding/json"
	"os"
)

// Defines what new profile json looks like. Note if we need to
// output profile json keys begin w/ lowercase, we'll have to split each fields
// into its own line so tag like `json:"version"` can work. otherwise
// tag is applied to all fields defined in same line and json.Marshal fails

// ProfileJSON handles new profile json schema
type ProfileJSON struct {
	// schema version, ignored for now but will be useful
	// when need to handle incompatible schema in the future
	Version  string
	Ethereum ProfileEthereum
	Osp      ProfileOsp
	Sgn      ProfileSgn
}

type ProfileEthereum struct {
	Gateway                                                  string
	ChainId, BlockIntervalSec, BlockDelayNum, DisputeTimeout uint64
	Contracts                                                ProfileContracts
	// CheckInterval is map of eventname to its check interval for monitor service
	// if not set (ie. 0) will check every blockIntervalSec (ie. same as check new block head)
	// if specify, key must be one of event.go const string values
	// monitor will check every checkInterval * blockIntervalSec
	CheckInterval map[string]uint64
}

type ProfileContracts struct {
	Wallet, Ledger, VirtResolver, EthPool, PayResolver, PayRegistry, RouterRegistry string
	Ledgers                                                                         map[string]string
}

type ProfileOsp struct {
	Host, Address, ExplorerUrl string
}

type ProfileSgn struct {
	Gateway         string
	SgnContractAddr string
}

func (pj *ProfileJSON) ToCProfile() *CProfile {
	cp := &CProfile{
		ChainId:            int64(pj.Ethereum.ChainId),
		ETHInstance:        pj.Ethereum.Gateway,
		BlockDelayNum:      pj.Ethereum.BlockDelayNum,
		PollingInterval:    pj.Ethereum.BlockIntervalSec,
		DisputeTimeout:     pj.Ethereum.DisputeTimeout,
		WalletAddr:         pj.Ethereum.Contracts.Wallet,
		LedgerAddr:         pj.Ethereum.Contracts.Ledger,
		VirtResolverAddr:   pj.Ethereum.Contracts.VirtResolver,
		EthPoolAddr:        pj.Ethereum.Contracts.EthPool,
		PayResolverAddr:    pj.Ethereum.Contracts.PayResolver,
		PayRegistryAddr:    pj.Ethereum.Contracts.PayRegistry,
		RouterRegistryAddr: pj.Ethereum.Contracts.RouterRegistry,
		Ledgers:            pj.Ethereum.Contracts.Ledgers,
		SvrETHAddr:         pj.Osp.Address,
		SvrRPC:             pj.Osp.Host,
		ExplorerUrl:        pj.Osp.ExplorerUrl,
		CheckInterval:      pj.Ethereum.CheckInterval, // json.Unmarshal guarantee non-nil map (could be empty)
		SgnGateway:         pj.Sgn.Gateway,            // json.Unmarshal guarantee non-nil map (could be empty)
		SgnContractAddr:    pj.Sgn.SgnContractAddr,    // json.Unmarshal guarantee non-nil map (could be empty)
	}
	return cp
}

// ParseProfile parses file content at path and returns CProfile
// supports both old and new schema
func ParseProfile(path string) *CProfile {
	raw, _ := os.ReadFile(path)
	return Bytes2Profile(raw)
}

func ParseProfileJSON(path string) *ProfileJSON {
	raw, _ := os.ReadFile(path)
	pj := new(ProfileJSON)
	json.Unmarshal(raw, pj)
	return pj
}

// Bytes2Profile does json.Unmarshal and return CProfile
func Bytes2Profile(data []byte) *CProfile {
	// Try parsing as new schema first
	pj := new(ProfileJSON)
	_ = json.Unmarshal(data, pj)

	// Heuristic: consider new schema valid if key fields are populated
	newSchemaOk := false
	if pj.Ethereum.Gateway != "" || pj.Ethereum.ChainId != 0 ||
		pj.Ethereum.Contracts.Wallet != "" || pj.Ethereum.Contracts.Ledger != "" ||
		pj.Osp.Host != "" || pj.Osp.Address != "" || pj.Version != "" {
		newSchemaOk = true
	}

	if newSchemaOk {
		return pj.ToCProfile()
	}

	// Fallback to old schema
	cp := new(CProfile)
	_ = json.Unmarshal(data, cp)
	return cp
}
