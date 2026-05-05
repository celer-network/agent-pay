// Copyright 2021 Celer Network

package cli

import (
	"encoding/json"
	"io/ioutil"

	"github.com/celer-network/agent-pay/ctype"
	"github.com/celer-network/agent-pay/utils"
	"github.com/celer-network/goutils/log"
)

// XnetConfig is the operator-supplied cross-net routing configuration.
//
// A "net" is the off-chain network identity of a (chainId, contractSet)
// pair, where the contract set is the deployment of CelerLedger /
// PayResolver / PayRegistry / EthPool / VirtResolver / Wallet that an OSP
// boots against (configured via its profile JSON). Two OSPs sharing a netId
// must boot against the same contract addresses on the same chain;
// otherwise channels and signed messages won't validate across them.
//
// netId is intentionally separate from `block.chainid` for two reasons:
//
//  1. A contract-set redeployment on the same chain (upgrades, hard-fork-
//     style migrations) is a new net even though chainId is unchanged —
//     the new PayResolver enforces `pay.payResolver == address(this)` and
//     the new Ledger enforces `initializer.ledger_address == address(this)`,
//     so signed messages don't cross over. A bridge pair operated by the
//     same business entity can carry pays across the migration window
//     without users having to coordinate the upgrade.
//  2. Test harnesses can simulate cross-net routing on a single geth
//     instance by deploying multiple contract sets and labelling them as
//     distinct nets.
//
// On-chain replay protection lives in the contracts: chainId and
// ledger_address are signed into PaymentChannelInitializer; chainId is
// signed into ConditionalPay; payResolver is bound on every pay. netId
// only drives off-chain forwarding decisions.
//
// Bridge OSPs at a net boundary are operated as a trust unit (typically by
// the same business entity). The cross-bridge link is a direct gRPC stream
// between the two bridge processes, not a payment channel — there's no
// shared on-chain state between bridges, and reconciliation between them
// is off-protocol.
type XnetConfig struct {
	NetId         uint64                       `json:"net_id"`         // local net id
	NetBridge     map[string]uint64            `json:"net_bridge"`     // bridgeAddr -> bridgeNetId
	BridgeRouting map[uint64]string            `json:"bridge_routing"` // destNetId -> nextHopBridgeAddr
	NetToken      map[string]map[uint64]string `json:"net_token"`      // localTokenAddr -> map(remoteNetId -> remoteTokenAddr)
}

func ParseXnetConfig(path string) (*XnetConfig, error) {
	raw, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	xnet := new(XnetConfig)
	json.Unmarshal(raw, xnet)
	return xnet, nil
}

func (p *Processor) ConfigXnet() {
	if *batchfile == "" {
		log.Fatal("no config file provided")
	}
	xnet, err := ParseXnetConfig(*batchfile)
	if err != nil {
		log.Fatal(err)
	}
	p.setNetId(xnet.NetId)
	for bridge, netid := range xnet.NetBridge {
		p.setNetBridge(bridge, netid)
	}
	for netid, bridge := range xnet.BridgeRouting {
		p.setBridgeRouting(netid, bridge)
	}
	for local, remote := range xnet.NetToken {
		for netid, token := range remote {
			p.setNetToken(netid, token, local)
		}
	}
}

func (p *Processor) SetNetId() {
	p.setNetId(*netid)
}

func (p *Processor) SetNetBridge() {
	p.setNetBridge(*bridgeaddr, *netid)
}

func (p *Processor) SetBridgeRouting() {
	p.setBridgeRouting(*netid, *bridgeaddr)
}

func (p *Processor) SetNetToken() {
	p.setNetToken(*netid, *tokenaddr, *localtoken)
}

func (p *Processor) DeleteNetBridge() {
	log.Infoln("Delete netbridge", *bridgeaddr)
	err := p.dal.DeleteNetBridge(ctype.Hex2Addr(*bridgeaddr))
	if err != nil {
		log.Fatal(err)
	}
}

func (p *Processor) DeleteBridgeRouting() {
	log.Infoln("Delete bridge routing for dest net id", *netid)
	err := p.dal.DeleteBridgeRouting(*netid)
	if err != nil {
		log.Fatal(err)
	}
}

func (p *Processor) DeleteNetToken() {
	log.Infof("Delete net token for net id: %d, token :%s", *netid, *tokenaddr)
	err := p.dal.DeleteNetToken(*netid, utils.GetTokenInfoFromAddress(ctype.Hex2Addr(*tokenaddr)))
	if err != nil {
		log.Fatal(err)
	}
}

func (p *Processor) setNetId(netid uint64) {
	log.Infoln("Update net id", netid)
	err := p.dal.PutNetId(netid)
	if err != nil {
		log.Fatal(err)
	}
}

func (p *Processor) setNetBridge(bridgeAddr string, netId uint64) {
	log.Infof("Update netbridge addr: %s, net id: %d", bridgeAddr, netId)
	err := p.dal.UpsertNetBridge(ctype.Hex2Addr(bridgeAddr), netId)
	if err != nil {
		log.Fatal(err)
	}
}

func (p *Processor) setBridgeRouting(netId uint64, bridgeAddr string) {
	log.Infof("Update bridge routing dest net id: %d, bridge addr: %s", netId, bridgeAddr)
	err := p.dal.UpsertBridgeRouting(netId, ctype.Hex2Addr(bridgeAddr))
	if err != nil {
		log.Fatal(err)
	}
}

func (p *Processor) setNetToken(netId uint64, tokenAddr, localToken string) {
	log.Infof("Update net token for net id: %d, net token %s, local token :%s", netId, tokenAddr, localToken)
	err := p.dal.UpsertNetToken(netId,
		utils.GetTokenInfoFromAddress(ctype.Hex2Addr(tokenAddr)),
		utils.GetTokenInfoFromAddress(ctype.Hex2Addr(localToken)))
	if err != nil {
		log.Fatal(err)
	}
}
