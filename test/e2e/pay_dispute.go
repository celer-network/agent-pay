// Copyright 2018-2025 Celer Network

// Conditional-pay dispute coverage. After AS-B trimmed the gaming
// state-machine machinery, the surviving "dispute" path is just:
//
//   1. send a conditional pay whose Condition references an IBooleanCond
//      contract (either VIRTUAL_CONTRACT bytecode registered off-chain or
//      a DEPLOYED_CONTRACT address);
//   2. invoke `PayResolver.resolvePaymentByConditions` (which deploys the
//      virtual contract on demand and calls IBooleanCond.{isFinalized,
//      getOutcome});
//   3. assert that the on-chain registry reflects the outcome (full amount
//      when getOutcome → true, zero when getOutcome → false).
//
// `BooleanCondMock` (deployed from `agent-pay-contracts/src/helper/`) is the
// fixture: a single byte argsQueryOutcome where any non-zero value → true
// and 0x00 / empty → false; isFinalized similarly returns true unless given
// 0x00.

package e2e

import (
	"fmt"
	"testing"

	"github.com/celer-network/agent-pay/ctype"
	"github.com/celer-network/agent-pay/entity"
	tf "github.com/celer-network/agent-pay/testing"
	"github.com/celer-network/agent-pay/testing/testapp"
	"github.com/celer-network/goutils/log"
)

func disputeEthPayWithVirtualContract(t *testing.T) {
	log.Info("============== start test disputeEthPayWithVirtualContract ==============")
	defer log.Info("============== end test disputeEthPayWithVirtualContract ==============")
	t.Parallel()
	disputePayWithVirtualContract(t, entity.TokenType_ETH, tokenAddrEth)
}

func disputeEthPayWithDeployedContract(t *testing.T) {
	log.Info("============== start test disputeEthPayWithDeployedContract ==============")
	defer log.Info("============== end test disputeEthPayWithDeployedContract ==============")
	t.Parallel()
	disputePayWithDeployedContract(t, entity.TokenType_ETH, tokenAddrEth)
}

// disputePayWithVirtualContract drives the VIRTUAL_CONTRACT path: register
// `BooleanCondMock` bytecode off-chain, send a conditional pay against the
// deterministic virtual address, then resolve on-chain. The first scenario
// uses argsQueryOutcome=0x01 (getOutcome → true) and asserts the receiver
// pulls the full amount; the second uses 0x00 (getOutcome → false) and
// asserts the registry stays at zero.
func disputePayWithVirtualContract(t *testing.T, tokenType entity.TokenType, tokenAddr string) {
	c1, c2, c1EthAddr, c2EthAddr, cleanup, err := setupTwoClientChannels(tokenType, tokenAddr)
	if err != nil {
		t.Error(err)
		return
	}
	defer cleanup()

	bytecode := testapp.BooleanCondMockBin
	constructor := []byte{}

	// Two distinct nonces so the two scenarios get different virtual addresses
	// and don't collide on chain.
	if err := runVirtualContractScenario(t, c1, c2, c1EthAddr, c2EthAddr, tokenType, tokenAddr,
		ctype.Hex2Bytes(bytecode), constructor, 1, []byte{0x01}, true); err != nil {
		t.Error(err)
		return
	}
	if err := runVirtualContractScenario(t, c1, c2, c1EthAddr, c2EthAddr, tokenType, tokenAddr,
		ctype.Hex2Bytes(bytecode), constructor, 2, []byte{0x00}, false); err != nil {
		t.Error(err)
		return
	}
}

// disputePayWithDeployedContract drives the DEPLOYED_CONTRACT path against
// the `BooleanCondMock` instance deployed in `setup_onchain.go`.
func disputePayWithDeployedContract(t *testing.T, tokenType entity.TokenType, tokenAddr string) {
	c1, c2, c1EthAddr, c2EthAddr, cleanup, err := setupTwoClientChannels(tokenType, tokenAddr)
	if err != nil {
		t.Error(err)
		return
	}
	defer cleanup()

	mockAddr, ok := appAddrMap["BooleanCondMock"]
	if !ok {
		t.Errorf("BooleanCondMock address not found in appAddrMap")
		return
	}

	if err := runDeployedContractScenario(t, c1, c2, c1EthAddr, c2EthAddr, tokenType, tokenAddr,
		mockAddr, []byte{0x01}, true); err != nil {
		t.Error(err)
		return
	}
	if err := runDeployedContractScenario(t, c1, c2, c1EthAddr, c2EthAddr, tokenType, tokenAddr,
		mockAddr, []byte{0x00}, false); err != nil {
		t.Error(err)
		return
	}
}

func runVirtualContractScenario(
	t *testing.T,
	c1, c2 *tf.ClientController,
	c1EthAddr, c2EthAddr string,
	tokenType entity.TokenType,
	tokenAddr string,
	bytecode []byte,
	constructor []byte,
	nonce uint64,
	queryBytes []byte,
	expectPaid bool,
) error {
	log.Infof("virtual-contract scenario: nonce=%d query=%x expectPaid=%v", nonce, queryBytes, expectPaid)

	appChanID, err := c1.NewAppChannelOnVirtualContract(bytecode, constructor, nonce, 100)
	if err != nil {
		return fmt.Errorf("c1 NewAppChannelOnVirtualContract: %w", err)
	}
	appChanID2, err := c2.NewAppChannelOnVirtualContract(bytecode, constructor, nonce, 100)
	if err != nil {
		return fmt.Errorf("c2 NewAppChannelOnVirtualContract: %w", err)
	}
	if appChanID != appChanID2 {
		return fmt.Errorf("virtual-contract address mismatch: c1=%s c2=%s", appChanID, appChanID2)
	}

	cond := &entity.Condition{
		ConditionType:          entity.ConditionType_VIRTUAL_CONTRACT,
		VirtualContractAddress: ctype.Hex2Bytes(appChanID),
		ArgsQueryOutcome:       queryBytes,
	}

	// Trigger the deploy-on-query path before on-chain pay resolution.
	// PayResolver.resolvePaymentByConditions calls
	// VirtContractResolver.resolve(virtAddr), which reverts with "Nonexistent
	// virtual address" if the virtual contract has never been deployed. The
	// off-chain GetBooleanOutcomeForAppSession path runs `deployIfNeeded`
	// before querying IBooleanCond, so calling it here is enough to ensure
	// the virtual address has bytecode by the time the resolve tx lands.
	finalized, outcome, err := c2.GetAppChannelBooleanOutcome(appChanID, queryBytes)
	if err != nil {
		return fmt.Errorf("GetAppChannelBooleanOutcome: %w", err)
	}
	wantOutcome := expectPaid
	wantFinalized := true
	if !expectPaid {
		// BooleanCondMock.isFinalized returns false for query=0x00.
		wantFinalized = false
	}
	if finalized != wantFinalized || outcome != wantOutcome {
		return fmt.Errorf("BooleanCondMock query %x: finalized=%v outcome=%v, want finalized=%v outcome=%v",
			queryBytes, finalized, outcome, wantFinalized, wantOutcome)
	}

	return runDisputeAndAssert(c1, c2, c1EthAddr, c2EthAddr, tokenType, tokenAddr, cond, expectPaid)
}

func runDeployedContractScenario(
	t *testing.T,
	c1, c2 *tf.ClientController,
	c1EthAddr, c2EthAddr string,
	tokenType entity.TokenType,
	tokenAddr string,
	mockAddr ctype.Addr,
	queryBytes []byte,
	expectPaid bool,
) error {
	log.Infof("deployed-contract scenario: addr=%x query=%x expectPaid=%v", mockAddr, queryBytes, expectPaid)

	cond := &entity.Condition{
		ConditionType:           entity.ConditionType_DEPLOYED_CONTRACT,
		DeployedContractAddress: mockAddr.Bytes(),
		ArgsQueryOutcome:        queryBytes,
	}
	return runDisputeAndAssert(c1, c2, c1EthAddr, c2EthAddr, tokenType, tokenAddr, cond, expectPaid)
}

// runDisputeAndAssert sends a conditional pay with the given `cond`, drives
// it through on-chain resolution via PayResolver, and asserts the resulting
// on-chain pay amount matches `expectPaid`.
func runDisputeAndAssert(
	c1, c2 *tf.ClientController,
	c1EthAddr, c2EthAddr string,
	tokenType entity.TokenType,
	tokenAddr string,
	cond *entity.Condition,
	expectPaid bool,
) error {
	payID, err := c1.SendPaymentWithBooleanConditions(
		c2EthAddr, sendAmt, tokenType, tokenAddr, []*entity.Condition{cond}, 100)
	if err != nil {
		return fmt.Errorf("SendPaymentWithBooleanConditions: %w", err)
	}
	if err := waitForPaymentPending(payID, c1, c2); err != nil {
		return fmt.Errorf("waitForPaymentPending: %w", err)
	}

	done := make(chan bool)
	go tf.AdvanceBlocksUntilDone(done)
	defer func() { done <- true }()

	amount, _, err := c2.SettleConditionalPayOnChain(payID)
	if err != nil {
		return fmt.Errorf("SettleConditionalPayOnChain: %w", err)
	}
	wantAmount := "0"
	if expectPaid {
		wantAmount = sendAmt
	}
	if amount != wantAmount {
		return fmt.Errorf("on-chain pay amount = %s, want %s (expectPaid=%v)", amount, wantAmount, expectPaid)
	}
	return nil
}

// setupTwoClientChannels creates funded c1 / c2 clients, opens a channel for
// each against the OSP, and returns a cleanup func that kills both clients.
func setupTwoClientChannels(tokenType entity.TokenType, tokenAddr string) (
	*tf.ClientController, *tf.ClientController, string, string, func(), error,
) {
	ks, addrs, err := tf.CreateAccountsWithBalance(2, accountBalance)
	if err != nil {
		return nil, nil, "", "", nil, fmt.Errorf("CreateAccountsWithBalance: %w", err)
	}
	if tokenAddr != tokenAddrEth {
		if err := tf.FundAccountsWithErc20(tokenAddr, addrs, accountBalance); err != nil {
			return nil, nil, "", "", nil, fmt.Errorf("FundAccountsWithErc20: %w", err)
		}
	}

	c1, err := tf.StartC1WithoutProxy(ks[0])
	if err != nil {
		return nil, nil, "", "", nil, fmt.Errorf("StartC1: %w", err)
	}
	c2, err := tf.StartC2WithoutProxy(ks[1])
	if err != nil {
		c1.Kill()
		return nil, nil, "", "", nil, fmt.Errorf("StartC2: %w", err)
	}
	cleanup := func() {
		c1.Kill()
		c2.Kill()
	}

	if _, err := c1.OpenChannel(addrs[0], tokenType, tokenAddr, initialBalance, initialBalance); err != nil {
		cleanup()
		return nil, nil, "", "", nil, fmt.Errorf("c1 OpenChannel: %w", err)
	}
	if _, err := c2.OpenChannel(addrs[1], tokenType, tokenAddr, initialBalance, initialBalance); err != nil {
		cleanup()
		return nil, nil, "", "", nil, fmt.Errorf("c2 OpenChannel: %w", err)
	}
	if err := c1.AssertBalance(tokenAddr, initialBalance, "0", initialBalance); err != nil {
		cleanup()
		return nil, nil, "", "", nil, fmt.Errorf("c1 AssertBalance: %w", err)
	}
	if err := c2.AssertBalance(tokenAddr, initialBalance, "0", initialBalance); err != nil {
		cleanup()
		return nil, nil, "", "", nil, fmt.Errorf("c2 AssertBalance: %w", err)
	}
	return c1, c2, addrs[0], addrs[1], cleanup, nil
}
