// Copyright 2018-2025 Celer Network

// SDK APIs dealing with app sessions backed by stateless `IBooleanCond`
// virtual condition contracts.
//
// Post-trim the legacy gaming surface (turn-based state-exchange protocol with
// `SignAppData` / `HandleMatchData` / opcodes / seqnum tracking, on-chain
// `applyAction` / introspection / oracle disputes, the `NewAppSessionOnDeployedContract`
// path and its multisession dependencies) is gone. What remains is the thin
// wrapper around `client.CelerClient`'s registration / outcome-query surface
// for VIRTUAL_CONTRACT condition contracts.

package celersdk

import (
	"github.com/celer-network/agent-pay/client"
	"github.com/celer-network/agent-pay/ctype"
)

type AppSession struct {
	ID string
	cc *client.CelerClient
}

// CreateAppSessionOnVirtualContract registers a VIRTUAL_CONTRACT condition
// contract on the cnode and returns an `AppSession` keyed by its deterministic
// virtual-contract address. The contract is deployed lazily on first outcome
// query (via `OnChainGetBooleanOutcome`).
func (mc *Client) CreateAppSessionOnVirtualContract(
	contractBin string,
	constructor string,
	nonce uint64) (*AppSession, error) {
	sessionID, err := mc.c.NewAppChannelOnVirtualContract(
		ctype.Hex2Bytes(contractBin),
		ctype.Hex2Bytes(constructor),
		nonce)
	if err != nil {
		return nil, err
	}
	return &AppSession{ID: sessionID, cc: mc.c}, nil
}

// EndAppSession removes the registered virtual condition contract from the
// cnode's in-memory bookkeeping. The current implementation always succeeds;
// it cannot fail.
func (mc *Client) EndAppSession(sessionid string) {
	mc.c.DeleteAppChannel(sessionid)
}

// GetDeployedAddress returns the on-chain deployed address of the registered
// virtual condition contract. Returns an error if the contract has not been
// deployed yet (deployment is triggered lazily by `OnChainGetBooleanOutcome`).
func (s *AppSession) GetDeployedAddress() (string, error) {
	addr, err := s.cc.GetAppChannelDeployedAddr(s.ID)
	return ctype.Addr2Hex(addr), err
}

// AppBooleanOutcome carries the result of an `IBooleanCond` query: whether the
// outcome has been finalized, and the boolean outcome itself.
type AppBooleanOutcome struct {
	Finalized bool
	Outcome   bool
}

// OnChainGetBooleanOutcome queries `IBooleanCond.{isFinalized,getOutcome}` on
// the registered condition contract. For VIRTUAL_CONTRACT this triggers
// deploy-on-query: if the virtual contract has not been deployed yet, this call
// submits a deployment transaction first.
func (s *AppSession) OnChainGetBooleanOutcome(query []byte) (*AppBooleanOutcome, error) {
	finalized, outcome, err := s.cc.OnChainGetAppChannelBooleanOutcome(s.ID, query)
	return &AppBooleanOutcome{Finalized: finalized, Outcome: outcome}, err
}
