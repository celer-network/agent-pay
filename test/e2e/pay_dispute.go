// Copyright 2018-2025 Celer Network

// Conditional-pay dispute coverage. After AS-B trimmed the gaming
// state-machine machinery, the surviving "dispute" path is just:
//
//   1. send a conditional pay whose Condition references an IBooleanCond
//      contract (either VIRTUAL_CONTRACT bytecode registered off-chain or
//      a DEPLOYED_CONTRACT address);
//   2. for VIRTUAL_CONTRACT only, ensure the contract is on-chain by
//      calling `GetBooleanOutcomeForAppSession` (the surviving deploy-on-
//      query path through `AppClient.deployIfNeeded`);
//   3. invoke `PayResolver.resolvePaymentByConditions`, which calls
//      `IBooleanCond.isFinalized(argsQueryFinalization)` and (only if true)
//      `IBooleanCond.getOutcome(argsQueryOutcome)` on the deployed
//      contract, and reverts with:
//        - `"Nonexistent virtual address"` if step 2 was skipped for a
//          VIRTUAL_CONTRACT condition;
//        - `"Condition is not finalized"` if isFinalized returns false.
//   4. assert the resulting on-chain pay amount (full when both bytes are
//      non-zero) or that the resolve reverted as expected.
//
// `BooleanCondMock` (deployed from `agent-pay-contracts/src/helper/`) is the
// fixture: a single byte argsQuery where any non-zero value → true,
// 0x00 → false, empty → true for isFinalized but false for getOutcome.
//
// Three on-chain protocol paths are covered (one per dispute call):
//
//   - **Symmetric pass** — argsQueryFinalization == argsQueryOutcome == 0x01
//     (`runVirtualContractScenario` / `runDeployedContractScenario` with
//     expectPaid=true): isFinalized=true, getOutcome=true → registry
//     resolves to full amount.
//   - **Symmetric not-finalized** — both bytes 0x00 (same helpers with
//     expectPaid=false): isFinalized=false → PayResolver reverts with
//     "Condition is not finalized" before getOutcome is ever called.
//   - **BOOLEAN_AND short-circuit** — finalize=0x01 + outcome=0x00
//     (`run{Virtual,Deployed}ContractFalseOutcomeScenario`): isFinalized=true
//     but getOutcome=false → registry resolves to amount=0 with no revert.
//     This is the path PayResolver takes when at least one BOOLEAN_AND
//     condition is finalized-but-false.
//
// Two negative scenarios pin down on-chain prerequisites for VIRTUAL_CONTRACT:
//   - `runVirtualContractResolveBeforeDeploy`: skipping step 2 reverts with
//     "Nonexistent virtual address";
//   - `runVirtualContractParallelDeploy`: concurrent deploy-on-query calls
//     converge on a single deploy tx (the `AppChannel.mu` mutex prevents
//     duplicate submissions that would revert with VirtContractResolver's
//     "Current real address is not 0" guard); also asserts the
//     VirtContractResolver `Deploy` event log has exactly one entry for the
//     virtual address.

package e2e

import (
	"fmt"
	"strings"
	"testing"

	"github.com/celer-network/agent-pay/app"
	"github.com/celer-network/agent-pay/chain/channel-eth-go/virtresolver"
	"github.com/celer-network/agent-pay/ctype"
	"github.com/celer-network/agent-pay/entity"
	tf "github.com/celer-network/agent-pay/testing"
	"github.com/celer-network/agent-pay/testing/testapp"
	"github.com/celer-network/goutils/log"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
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
	c1, c2, _, c2EthAddr, cleanup, err := setupTwoClientChannels(tokenType, tokenAddr)
	if err != nil {
		t.Error(err)
		return
	}
	defer cleanup()

	bytecode := testapp.BooleanCondMockBin
	constructor := []byte{}

	// Three distinct nonces so the three scenarios get different virtual
	// addresses and don't collide on chain.
	if err := runVirtualContractScenario(c1, c2, c2EthAddr, tokenType, tokenAddr,
		ctype.Hex2Bytes(bytecode), constructor, 1, []byte{0x01}, true); err != nil {
		t.Error(err)
		return
	}
	if err := runVirtualContractScenario(c1, c2, c2EthAddr, tokenType, tokenAddr,
		ctype.Hex2Bytes(bytecode), constructor, 2, []byte{0x00}, false); err != nil {
		t.Error(err)
		return
	}
	// BOOLEAN_AND short-circuit: isFinalized=true, getOutcome=false. The pay
	// resolves on-chain to amount=0 (no revert) — the third protocol-supported
	// outcome that the symmetric `runVirtualContractScenario` calls above don't
	// reach. argsQueryFinalization=0x01 → isFinalized=true; argsQueryOutcome=0x00
	// → getOutcome=false → BOOLEAN_AND short-circuits to amount=0.
	if err := runVirtualContractFalseOutcomeScenario(c1, c2, c2EthAddr, tokenType, tokenAddr,
		ctype.Hex2Bytes(bytecode), constructor, 3); err != nil {
		t.Error(err)
		return
	}
	// Negative scenario: resolving a VIRTUAL_CONTRACT pay before the virtual
	// contract is deployed must fail at PayResolver. This documents the
	// deploy-before-resolve contract enforced on-chain by VirtContractResolver.
	if err := runVirtualContractResolveBeforeDeploy(c1, c2, c2EthAddr, tokenType, tokenAddr,
		ctype.Hex2Bytes(bytecode), constructor, 4); err != nil {
		t.Error(err)
		return
	}
	// Concurrency scenario: N parallel deploy-on-query calls for the same
	// virtual contract must converge on a single on-chain deploy tx. The
	// AppChannel.mu mutex serializes deployIfNeeded; a regression that
	// removes that locking would either waste a deploy tx or revert the
	// second submission with VirtContractResolver's "Current real address
	// is not 0" guard. The test also asserts exactly one Deploy event was
	// emitted by VirtContractResolver, which strengthens "no caller fails"
	// to "exactly one tx submitted".
	if err := runVirtualContractParallelDeploy(c2, ctype.Hex2Bytes(bytecode), constructor, 5); err != nil {
		t.Error(err)
		return
	}
}

// disputePayWithDeployedContract drives the DEPLOYED_CONTRACT path against
// the `BooleanCondMock` instance deployed in `setup_onchain.go`.
func disputePayWithDeployedContract(t *testing.T, tokenType entity.TokenType, tokenAddr string) {
	c1, c2, _, c2EthAddr, cleanup, err := setupTwoClientChannels(tokenType, tokenAddr)
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

	if err := runDeployedContractScenario(c1, c2, c2EthAddr, tokenType, tokenAddr,
		mockAddr, []byte{0x01}, true); err != nil {
		t.Error(err)
		return
	}
	if err := runDeployedContractScenario(c1, c2, c2EthAddr, tokenType, tokenAddr,
		mockAddr, []byte{0x00}, false); err != nil {
		t.Error(err)
		return
	}
	// BOOLEAN_AND short-circuit: argsQueryFinalization=0x01 (isFinalized=true)
	// + argsQueryOutcome=0x00 (getOutcome=false) → registry resolves to
	// amount=0 with no revert. Same protocol path as the virtual-contract
	// false-outcome scenario above.
	if err := runDeployedContractFalseOutcomeScenario(c1, c2, c2EthAddr, tokenType, tokenAddr, mockAddr); err != nil {
		t.Error(err)
		return
	}
}

func runVirtualContractScenario(
	c1, c2 *tf.ClientController,
	c2EthAddr string,
	tokenType entity.TokenType,
	tokenAddr string,
	bytecode []byte,
	constructor []byte,
	nonce uint64,
	queryBytes []byte,
	expectPaid bool,
) error {
	log.Infof("virtual-contract scenario: nonce=%d query=%x expectPaid=%v", nonce, queryBytes, expectPaid)

	appChanID, err := c1.NewAppChannelOnVirtualContract(bytecode, constructor, nonce)
	if err != nil {
		return fmt.Errorf("c1 NewAppChannelOnVirtualContract: %w", err)
	}
	appChanID2, err := c2.NewAppChannelOnVirtualContract(bytecode, constructor, nonce)
	if err != nil {
		return fmt.Errorf("c2 NewAppChannelOnVirtualContract: %w", err)
	}
	if appChanID != appChanID2 {
		return fmt.Errorf("virtual-contract address mismatch: c1=%s c2=%s", appChanID, appChanID2)
	}

	// Set ArgsQueryFinalization symmetrically with ArgsQueryOutcome so the
	// on-chain `isFinalized` is called with the same non-empty bytes as the
	// off-chain assertion below — i.e. for queryBytes=0x00 the negative case
	// reaches PayResolver via isFinalized=false (the protocol-correct rejection)
	// rather than implicitly via an empty-bytes quirk.
	cond := &entity.Condition{
		ConditionType:          entity.ConditionType_VIRTUAL_CONTRACT,
		VirtualContractAddress: ctype.Hex2Bytes(appChanID),
		ArgsQueryFinalization:  queryBytes,
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
	wantFinalized := expectPaid
	if finalized != wantFinalized || outcome != wantOutcome {
		return fmt.Errorf("BooleanCondMock query %x: finalized=%v outcome=%v, want finalized=%v outcome=%v",
			queryBytes, finalized, outcome, wantFinalized, wantOutcome)
	}

	return runDisputeAndAssert(c1, c2, c2EthAddr, tokenType, tokenAddr, cond, expectPaid)
}

func runDeployedContractScenario(
	c1, c2 *tf.ClientController,
	c2EthAddr string,
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
		ArgsQueryFinalization:   queryBytes,
		ArgsQueryOutcome:        queryBytes,
	}
	return runDisputeAndAssert(c1, c2, c2EthAddr, tokenType, tokenAddr, cond, expectPaid)
}

// runDisputeAndAssert sends a conditional pay with the given `cond`, drives
// it through on-chain resolution via PayResolver, and asserts the result.
//
// `expectPaid=true` means both isFinalized and getOutcome return true — the
// pay resolves to the full amount.
//
// `expectPaid=false` means both isFinalized and getOutcome return false (the
// scenario sets ArgsQueryFinalization=ArgsQueryOutcome=0x00). PayResolver
// requires isFinalized=true, so the on-chain resolve must revert with
// "Condition is not finalized" — pay never lands in the registry. The
// alternative "isFinalized=true, getOutcome=false" path (BOOLEAN_AND with a
// false-outcome condition resolving to amount=0) is implicitly covered by
// the wire format and is not exercised here to keep this scenario aligned
// with the on-chain contract's strict isFinalized requirement.
func runDisputeAndAssert(
	c1, c2 *tf.ClientController,
	c2EthAddr string,
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
	if !expectPaid {
		if err == nil {
			return fmt.Errorf("SettleConditionalPayOnChain unexpectedly succeeded for not-finalized condition (amount=%s)", amount)
		}
		if !strings.Contains(err.Error(), "Condition is not finalized") {
			return fmt.Errorf("SettleConditionalPayOnChain error = %v, want substring %q", err, "Condition is not finalized")
		}
		return nil
	}
	if err != nil {
		return fmt.Errorf("SettleConditionalPayOnChain: %w", err)
	}
	if amount != sendAmt {
		return fmt.Errorf("on-chain pay amount = %s, want %s (expectPaid=true)", amount, sendAmt)
	}
	return nil
}

// runDeployedContractFalseOutcomeScenario is the DEPLOYED_CONTRACT counterpart
// to `runVirtualContractFalseOutcomeScenario`: argsQueryFinalization=0x01 +
// argsQueryOutcome=0x00 against the on-chain BooleanCondMock instance from
// `setup_onchain.go` resolves to amount=0 (BOOLEAN_AND short-circuit, no
// revert).
func runDeployedContractFalseOutcomeScenario(
	c1, c2 *tf.ClientController,
	c2EthAddr string,
	tokenType entity.TokenType,
	tokenAddr string,
	mockAddr ctype.Addr,
) error {
	log.Infof("deployed-contract false-outcome scenario: addr=%x", mockAddr)

	cond := &entity.Condition{
		ConditionType:           entity.ConditionType_DEPLOYED_CONTRACT,
		DeployedContractAddress: mockAddr.Bytes(),
		ArgsQueryFinalization:   []byte{0x01},
		ArgsQueryOutcome:        []byte{0x00},
	}

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
	if amount != "0" {
		return fmt.Errorf("deployed-contract BOOLEAN_AND-with-false-outcome amount = %s, want 0", amount)
	}
	log.Infof("deployed-contract false-outcome scenario correctly resolved to amount=0")
	return nil
}

// runVirtualContractFalseOutcomeScenario covers the third protocol-supported
// outcome that the symmetric `runVirtualContractScenario` paths don't reach:
// `(isFinalized=true, getOutcome=false)` resolves on-chain to amount=0 via
// PayResolver's BOOLEAN_AND short-circuit. argsQueryFinalization=0x01 makes
// `BooleanCondMock.isFinalized(0x01)` return true so the resolve doesn't
// revert; argsQueryOutcome=0x00 makes `getOutcome(0x00)` return false so the
// pay registry stays at amount=0. Pre-deploys via the deploy-on-query path so
// the resolve doesn't hit "Nonexistent virtual address".
func runVirtualContractFalseOutcomeScenario(
	c1, c2 *tf.ClientController,
	c2EthAddr string,
	tokenType entity.TokenType,
	tokenAddr string,
	bytecode []byte,
	constructor []byte,
	nonce uint64,
) error {
	log.Infof("virtual-contract false-outcome scenario: nonce=%d", nonce)

	appChanID, err := c1.NewAppChannelOnVirtualContract(bytecode, constructor, nonce)
	if err != nil {
		return fmt.Errorf("c1 NewAppChannelOnVirtualContract: %w", err)
	}
	if _, err := c2.NewAppChannelOnVirtualContract(bytecode, constructor, nonce); err != nil {
		return fmt.Errorf("c2 NewAppChannelOnVirtualContract: %w", err)
	}

	// Trigger deploy-on-query before the on-chain resolve so PayResolver
	// doesn't revert with "Nonexistent virtual address".
	if _, _, err := c2.GetAppChannelBooleanOutcome(appChanID, []byte{0x01}); err != nil {
		return fmt.Errorf("GetAppChannelBooleanOutcome (deploy trigger): %w", err)
	}

	cond := &entity.Condition{
		ConditionType:          entity.ConditionType_VIRTUAL_CONTRACT,
		VirtualContractAddress: ctype.Hex2Bytes(appChanID),
		ArgsQueryFinalization:  []byte{0x01},
		ArgsQueryOutcome:       []byte{0x00},
	}

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
	if amount != "0" {
		return fmt.Errorf("BOOLEAN_AND-with-false-outcome on-chain pay amount = %s, want 0 (registry should stay at zero, not revert)", amount)
	}
	log.Infof("false-outcome scenario correctly resolved to amount=0")
	return nil
}

// runVirtualContractResolveBeforeDeploy registers a virtual condition
// contract (so the off-chain Condition is well-formed) but deliberately skips
// the deploy-on-query step, then sends a conditional pay against that
// undeployed virtual address and asserts that PayResolver rejects the resolve
// with "Nonexistent virtual address". This pins down the on-chain
// deploy-before-resolve contract that the positive scenarios above rely on.
func runVirtualContractResolveBeforeDeploy(
	c1, c2 *tf.ClientController,
	c2EthAddr string,
	tokenType entity.TokenType,
	tokenAddr string,
	bytecode []byte,
	constructor []byte,
	nonce uint64,
) error {
	log.Infof("virtual-contract negative scenario (resolve-before-deploy): nonce=%d", nonce)

	appChanID, err := c1.NewAppChannelOnVirtualContract(bytecode, constructor, nonce)
	if err != nil {
		return fmt.Errorf("c1 NewAppChannelOnVirtualContract: %w", err)
	}
	if _, err := c2.NewAppChannelOnVirtualContract(bytecode, constructor, nonce); err != nil {
		return fmt.Errorf("c2 NewAppChannelOnVirtualContract: %w", err)
	}

	cond := &entity.Condition{
		ConditionType:          entity.ConditionType_VIRTUAL_CONTRACT,
		VirtualContractAddress: ctype.Hex2Bytes(appChanID),
		ArgsQueryOutcome:       []byte{0x01},
	}

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

	// Intentionally skip GetAppChannelBooleanOutcome — the virtual contract
	// must remain undeployed for this scenario.
	_, _, err = c2.SettleConditionalPayOnChain(payID)
	if err == nil {
		return fmt.Errorf("SettleConditionalPayOnChain unexpectedly succeeded against undeployed virtual contract")
	}
	if !strings.Contains(err.Error(), "Nonexistent virtual address") {
		return fmt.Errorf("SettleConditionalPayOnChain error = %v, want substring %q", err, "Nonexistent virtual address")
	}
	log.Infof("resolve-before-deploy correctly rejected: %v", err)
	return nil
}

// runVirtualContractParallelDeploy fires N concurrent deploy-on-query calls
// against a freshly-registered virtual contract — interleaved with a few
// `GetAppChannelDeployedAddr` probes that take the same per-AppChannel mutex
// — and asserts:
//
//   - Every concurrent call returns either `(true, true, nil)` from
//     GetBooleanOutcome or a successfully-resolved address from
//     GetAppChannelDeployedAddr.
//   - Exactly one VirtContractResolver `Deploy` event was emitted for the
//     virtual address. This is the stronger invariant the mutex is supposed
//     to guarantee: "exactly one deploy tx submitted," not just "no caller
//     observed an error."
//
// Without the mutex, at least one goroutine would submit a duplicate deploy
// tx and one of two things would happen: VirtContractResolver reverts the
// second one with "Current real address is not 0" (the first goroutine sees
// an error from `SubmitWaitMined`) or the second tx lands and a Deploy event
// count of 2 would surface here. With the mutex, the first goroutine deploys,
// every other goroutine sees the cached `DeployedAddr` and returns without
// a tx.
func runVirtualContractParallelDeploy(
	c *tf.ClientController,
	bytecode []byte,
	constructor []byte,
	nonce uint64,
) error {
	log.Infof("virtual-contract parallel-deploy scenario: nonce=%d", nonce)

	appChanID, err := c.NewAppChannelOnVirtualContract(bytecode, constructor, nonce)
	if err != nil {
		return fmt.Errorf("NewAppChannelOnVirtualContract: %w", err)
	}

	const outcomeWorkers = 4
	const addrWorkers = 2
	totalWorkers := outcomeWorkers + addrWorkers

	errs := make(chan error, totalWorkers)
	for i := 0; i < outcomeWorkers; i++ {
		go func(idx int) {
			f, o, e := c.GetAppChannelBooleanOutcome(appChanID, []byte{0x01})
			if e != nil {
				errs <- fmt.Errorf("parallel GetAppChannelBooleanOutcome[%d]: %w", idx, e)
				return
			}
			if !f || !o {
				errs <- fmt.Errorf("parallel GetAppChannelBooleanOutcome[%d]: finalized=%v outcome=%v, want true/true",
					idx, f, o)
				return
			}
			errs <- nil
		}(i)
	}
	// Interleave a couple of `GetAppChannelDeployedAddr` calls. This path
	// takes the same `appChannel.mu` as `GetBooleanOutcome` so it should
	// either return the resolved address (after the deploy lands) or block
	// on the mutex until it does. Either way, no error.
	for i := 0; i < addrWorkers; i++ {
		go func(idx int) {
			addr, e := c.GetAppChannelDeployedAddr(appChanID)
			if e != nil {
				errs <- fmt.Errorf("parallel GetAppChannelDeployedAddr[%d]: %w", idx, e)
				return
			}
			if addr == "" {
				errs <- fmt.Errorf("parallel GetAppChannelDeployedAddr[%d]: empty address", idx)
				return
			}
			errs <- nil
		}(i)
	}
	for i := 0; i < totalWorkers; i++ {
		if e := <-errs; e != nil {
			return e
		}
	}

	// Strengthen the regression test: assert exactly one Deploy event was
	// emitted for the virtual address. A buggy implementation that, e.g.,
	// serialized via a coarser-grained lock that nonetheless allowed a
	// duplicate submission and swallowed the resulting revert in the second
	// goroutine would fail this check.
	virtAddr := app.GetVirtualAddress(bytecode, constructor, nonce)
	var virt32 [32]byte
	copy(virt32[:], virtAddr[:])
	resolver, err := virtresolver.NewVirtContractResolverFilterer(channelAddrBundle.VirtResolverAddr, conclient)
	if err != nil {
		return fmt.Errorf("NewVirtContractResolverFilterer: %w", err)
	}
	iter, err := resolver.FilterDeploy(&bind.FilterOpts{}, [][32]byte{virt32})
	if err != nil {
		return fmt.Errorf("FilterDeploy: %w", err)
	}
	defer iter.Close()
	deployCount := 0
	for iter.Next() {
		deployCount++
	}
	if iter.Error() != nil {
		return fmt.Errorf("Deploy event iteration: %w", iter.Error())
	}
	if deployCount != 1 {
		return fmt.Errorf("VirtContractResolver Deploy events for virtAddr = %d, want exactly 1 (mutex should serialize concurrent deploys to a single tx)", deployCount)
	}
	log.Infof("parallel-deploy scenario passed: %d concurrent calls converged on a single deploy tx (Deploy event count = %d)", totalWorkers, deployCount)
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
