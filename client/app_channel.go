// Copyright 2018-2025 Celer Network

package client

import (
	"github.com/celer-network/agent-pay/ctype"
)

// NewAppChannelOnVirtualContract registers a VIRTUAL_CONTRACT condition
// contract on the cnode and returns its deterministic virtual-contract address
// (used as the session id / Condition.OnChainAddress for VIRTUAL_CONTRACT pays).
// The bytecode + constructor + nonce are stored so the contract can be deployed
// on-chain on demand (e.g. on dispute or during outcome query).
func (c *CelerClient) NewAppChannelOnVirtualContract(
	byteCode []byte,
	constructor []byte,
	nonce uint64,
	onchainTimeout uint64) (string, error) {
	return c.cNode.AppClient.NewAppChannelOnVirtualContract(byteCode, constructor, nonce, onchainTimeout)
}

// DeleteAppChannel removes the registered virtual condition contract from the
// cnode's in-memory bookkeeping. Does not touch on-chain state.
func (c *CelerClient) DeleteAppChannel(cid string) error {
	c.cNode.AppClient.DeleteAppChannel(cid)
	return nil
}

// GetAppChannelDeployedAddr returns the on-chain deployed address of a
// registered virtual condition contract, probing the virt-resolver if needed.
// Returns an error if the contract has not been deployed yet.
func (c *CelerClient) GetAppChannelDeployedAddr(cid string) (ctype.Addr, error) {
	return c.cNode.AppClient.GetAppChannelDeployedAddr(cid)
}

// OnChainGetAppChannelBooleanOutcome queries IBooleanCond.{isFinalized,
// getOutcome} on the registered condition contract. For VIRTUAL_CONTRACT this
// triggers deploy-on-query: if the virtual contract has not been deployed yet,
// this call submits a deployment transaction first.
func (c *CelerClient) OnChainGetAppChannelBooleanOutcome(cid string, query []byte) (bool, bool, error) {
	return c.cNode.AppClient.GetBooleanOutcome(cid, query)
}
