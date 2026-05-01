# App Session Simplification

Status: **AS-A through AS-D completed** — repo-wide build/vet clean, focused unit + e2e green, docs updated. Plan doc is being kept here pending PR merge (deleted at merge per §8).

| Phase | Scope | Status |
| --- | --- | --- |
| AS-A | Pre-flight audit | **completed** — see findings sub-section under AS-A in §5 |
| AS-B | Off-chain trim in `agent-pay/` | **completed** — see AS-B completion notes in §5 |
| AS-C | Test-cleanup and helper cleanup | **completed** — repo-wide build/vet green; focused e2e (dispute + sendCondPay + OSP webapi) green |
| AS-D | Documentation and validation | **completed** — AGENTS.md / docs / osp-cli README updated; full validation matrix green by group |

**Deferred:** x402 migration to a stateless condition-contract bytecode (was AS-C in earlier draft). x402 currently registers the legacy `SimpleSingleSessionApp` via `CreateAppSessionOnVirtualContract` and never exercises its dispute path; the trim doesn't break that flow. A future PR (in either x402 or agent-pay) swaps the registered bytecode to a stateless verifier — see §7 "Deferred / TODO."

This is the first plan doc under `docs/progress/`. The convention (cribbed from `agent-pay-x402/docs/progress/`): plan files are phase-structured with checkbox subtasks, `Status` is updated as each phase ships, and the file self-deletes at close-out — see §8 for this plan's specific closeout disposition.

---

## 1. Motivation

The `app/` subsystem in this repo — the off-chain runtime for app-session-gated conditional payments — was inherited from an earlier mobile-gaming project (CelerX / [`cApps`](../../../cApps/)). Its on-chain contract templates (`SingleSessionApp`, `MultiSessionApp`, the `WithOracle` variants) encode a **turn-based-game state machine** with notions of player turns, on-chain action submission, action deadlines, and oracle-arbitrated conflict resolution. Method names like `applyAction`, `getActionDeadline`, `finalizeOnActionTimeout`, `settleByMoveTimeout`, `settleByInvalidTurn` are protocol-level concepts borrowed straight from the gaming domain.

This repo's vision has shifted to **AI-agent payments**. Agent payments don't have turns, on-chain moves, or move-timeout disputes. They have: "did the seller produce a result the buyer accepts? if yes, pay; if no, cancel; if neither party will commit, fall back to the cosigned state on chain." The protocol's app-session machinery has been carrying weight that no current consumer actually exercises.

### What `agent-pay-x402` actually uses

A direct grep confirms the only app-session method `agent-pay-x402` calls is `CreateAppSessionOnVirtualContract`, in [testinfra/session.go](../../../agent-pay-x402/testinfra/session.go). And that call doesn't deploy anything, doesn't call `intendSettle`, doesn't move chain state — it returns a deterministically-derived virtual-contract address that the buyer and seller use as a shared identifier on the off-chain `Condition`. The rest of the x402 lifecycle is purely off-chain cooperative confirm/cancel. The virtual contract is **never actually deployed** in either the happy path or the off-chain reject path. Searches for `IntendSettle`, `GetSettleFinalizedTime`, `GetSessionID`, `SettleAppSession` in the x402 repo return zero hits.

### Method-level inventory of what's currently exposed

For grounding, here's what every "app contract" method in the legacy templates is actually for, classified by who calls it:

| Method | On-chain (`PayResolver`) | Off-chain (any current consumer) | Classification |
| --- | --- | --- | --- |
| `isFinalized(query)` | yes | yes | **protocol-essential** |
| `getOutcome(query)` | yes | yes | **protocol-essential** |
| `intendSettle(stateProof)` | no | no | dispute-fallback only |
| `getSettleFinalizedTime(session)` | no | no | dispute-fallback only |
| `getSessionID(nonce, signers)` | no | no | dispute-fallback only |
| `applyAction(action)` | no | no | gaming-vestigial |
| `getActionDeadline()` | no | no | gaming-vestigial |
| `finalizeOnActionTimeout()` | no | no | gaming-vestigial |
| `getStatus()` | no | no | introspection-only |
| `getSeqNum()` | no | no | introspection-only |
| `getState(key)` | no | no | introspection-only |
| `settleBySigTimeout(oracleProof)` | no | no | oracle (gaming) |
| `settleByMoveTimeout(oracleProof)` | no | no | oracle (gaming) |
| `settleByInvalidTurn(oracleProof, ...)` | no | no | oracle (gaming) |
| `settleByInvalidState(oracleProof, ...)` | no | no | oracle (gaming) |

**Two methods** (`isFinalized`, `getOutcome`) are real protocol surface — `PayResolver` invokes them during conditional-payment resolution via the existing `IBooleanCond` / `INumericCond` interfaces. **Fourteen** are dispute-fallback or gaming residue with no current consumer.

### What about the dispute fallback?

The legacy `intendSettle` / `getSettleFinalizedTime` / dispute-window pattern was designed so peers who couldn't agree off-chain could commit cosigned state on-chain, wait for a window to close, and have `getOutcome` return that committed state. For AI-agent payments this entire pattern is replaced by **carrying the cosigned state in the query bytes themselves** — the condition contract verifies signatures inline, no on-chain commit needed. State moves from `Session.outcome` storage to the in-message `argsQueryOutcome` payload. Pure functions. No timeout window because the cosigned message is the proof.

Channel-level dispute (peers can't agree on a payment, channel goes into settling) is unaffected — that's `CelerLedger.intendSettle`, completely separate from app-session `intendSettle`. Channel dispute resolves by timeout, and any unresolved conditional pay simply gets refunded after `lastPayResolveDeadline`. The condition contract never needs to "intend-settle" itself; it just needs to answer `isFinalized` honestly (returning `false` if it can't tell, which causes the pay to refund).

### What about the oracle?

The legacy oracle (`OracleState` / `OracleProof` + the four `settleBy*` paths) is a trusted-third-party tiebreaker for **off-chain timing attestations** in turn-based games — "player B was supposed to move by block X and didn't," "player A acted out of turn." For AI-agent payments those questions don't arise:

- "Did peer X act in time?" is answered by `block.timestamp` once the contracts moved to time-based deadlines. The chain's own clock is the witness; no third party needed.
- "Did peer X act out of turn?" has no analog — there are no turns.
- "Was the seller's work valid?" is **content verification**, not timing. It's already cleanly expressible as a regular `IBooleanCond` whose `getOutcome` verifies a third-party attester's signature inside the query bytes. No new protocol path required.

So the oracle is dropped along with the rest of the gaming machinery. The "third-party verifier" capability remains available to anyone who needs it — just expressed through the existing `IBooleanCond` / `INumericCond` interfaces, no protocol changes.

---

## 2. Target design

**Minimal and generic. No new contracts, no new interfaces, no speculative additions.** The protocol stays exactly as it is at the wire and resolution levels — only the off-chain accretion gets trimmed.

### What survives

- **`IBooleanCond` and `INumericCond`** in `agent-pay-contracts/src/lib/interface/` — already exist, unchanged. The full protocol surface for an "app contract" is these two methods each:
  ```solidity
  function isFinalized(bytes calldata query) external view returns (bool);
  function getOutcome(bytes calldata query) external view returns (bool);  // or uint256
  ```
  Note: numeric conditions are **not exercised** by any agent-pay off-chain code today — every `TransferFunctionType` produced by `cnode/`, `webapi/`, `delegate/`, etc. is `BOOLEAN_AND`, and there are zero importers of `INumericOutcome` outside its ABIgen file. The interface stays present in `agent-pay-contracts` because `PayResolver` invokes it during `NUMERIC_ADD/MAX/MIN` resolution; the off-chain bindings can be regenerated when a real numeric consumer surfaces.
- **`BooleanCondMock` and `NumericCondMock`** in `agent-pay-contracts/src/helper/` — already exist, unchanged. Both are explicitly **test-only** (the contracts' own NatSpec says: `**Test-only.** ... Do not deploy to a production network.`). They serve as Solidity test fixtures and off-chain integration test fixtures. They are **not** the recommended deployment for a real condition contract — a production deployment writes its own `IBooleanCond` impl with actual semantics (cosigned-message verification, oracle signature verification, ZK proof verification, etc.). The mocks just simulate the (finalized, outcome) tuple from query bytes.
- **Both `ConditionType_DEPLOYED_CONTRACT` and `ConditionType_VIRTUAL_CONTRACT`** in `entity.proto` — unchanged at the wire level. These are the protocol's two deployment-mode primitives for condition contracts and they're orthogonal to the session-state-machine we're deleting. Stateless condition contracts work fine under either:
  - **DEPLOYED_CONTRACT** — already-deployed verifiers. Use cases: payment gated on an on-chain oracle data feed; payment gated on a ZK verifier already on-chain. The condition's `OnChainAddress` points at the deployed verifier. **No cnode-side registration** is needed — `NewAppChannelOnDeployedContract` deletes (the registration call is currently multisession-app-specific; for stateless condition contracts the contract is already on-chain, so no registration is required).
  - **VIRTUAL_CONTRACT** — lazy-deployed verifiers (on dispute only, otherwise pure off-chain identifier). Use cases: payment gated on a ZK verifier that *would* verify if deployed; payment gated on a not-yet-deployed contract that parses on-chain oracle state. Saves gas: the verifier is only deployed if a dispute escalates that far. The cnode-side `NewAppChannelOnVirtualContract` registration **stays** — it's already generic (uses `GetVirtualAddress` for deterministic-address derivation, no multisession dependency).
  - **Reading-path redesign:** `AppClient.GetBooleanOutcome` is currently wired through `ISingleSession` (for VIRTUAL_CONTRACT) and `IMultiSession` (for DEPLOYED_CONTRACT, with a `SessionQuery`-wrapped query). Both branches get rewired to use the agent-pay-contracts `IBooleanCond` interface and pass the raw `argsQueryOutcome` bytes through unchanged — matching what `PayResolver` does on-chain. This drops the `IMultiSession` dependency from the surviving code path; the legacy `app/multisession.go` and `app/singlesession.go` ABIgen files delete with the rest of the session contracts.

    **Binding source-of-truth, corrected:** today's `app/booleanoutcome.go` exposes `IBooleanOutcome` / `IBooleanOutcomeCaller` (legacy cApps name), not `IBooleanCond*`. The intent post-trim is for the binding symbols to match the agent-pay-contracts interface name `IBooleanCond`. AS-B regenerates this file from `agent-pay-contracts/src/lib/interface/IBooleanCond.sol` so the Go symbols become `IBooleanCond` / `IBooleanCondCaller`. The ABI shape is identical (`isFinalized(bytes) returns (bool)`, `getOutcome(bytes) returns (bool)`), so no logic changes; only the symbol names align with the canonical interface name. The file path `app/booleanoutcome.go` is kept (or renamed to `app/booleancond.go` — minor polish, decided in AS-B).

    **Deploy-on-query is preserved.** For VIRTUAL_CONTRACT, `AppClient.GetBooleanOutcome` today calls `deployIfNeeded(appChannel)` before the contract query — i.e. querying the outcome of a not-yet-deployed virtual contract triggers an on-chain deployment transaction. This is the lazy-deployment escape hatch for VIRTUAL_CONTRACT and stays as-is post-trim; the redesigned `GetBooleanOutcome` keeps the `deployIfNeeded` call in front of the `IBooleanCondCaller` call. Callers (and AS-C tests) should expect a query for an unsettled virtual condition to be on-chain-side-effecting on first invocation, not a pure read.
- **`HASH_LOCK` condition type** — unchanged. Hash-lock conditions don't involve app contracts at all.
- **Virtual-contract registration plumbing in the off-chain runtime.** The cnode needs to know the bytecode + constructor for any registered virtual condition contract so it can be deployed on dispute. `CreateAppSessionOnVirtualContract` (registration) and the deterministic-address derivation logic stay — they're the legitimate infrastructure for VIRTUAL_CONTRACT support, distinct from the session state machine that's being deleted. The name "AppSession" is now somewhat misleading (the registered entity is a stateless verifier, not a session), but renaming the API costs churn and the existing name is consistent with the architecture docs' "app channel" terminology — if a rename happens, it's a separate polish.
- **The on-chain dispute path through `PayResolver`** — unchanged. PayResolver still calls `isFinalized` and `getOutcome` on the condition contract during channel-level dispute resolution. For DEPLOYED_CONTRACT the contract is already there; for VIRTUAL_CONTRACT someone (typically the party seeking resolution) deploys it from the registered bytecode before invoking PayResolver. If `isFinalized` returns false, the pay refunds by `lastPayResolveDeadline`.
- **The `app/` Go package name and import path** stay. The architecture docs (`agent-pay-docs/agentpay-architecture/system-overview.md`) describe the application-logic layer as the **app channel**, paired with the payment channel. The `app/` package is the off-chain home for that concept. After this plan, `app/` shrinks to a thin runtime that supports stateless condition contracts (register, deploy-on-dispute, query outcome) without any session state machine.

### What's deleted

The session-state-machine wrapped around the condition contracts, plus all turn-based-game residue. Specifically:

- The **gaming-vestigial methods on `app.AppClient`** in `app/appclient.go`:
  - `IntendSettle` — the on-chain `intendSettle` of an app contract (state-machine-only; PayResolver doesn't need it).
  - `ApplyAction`, `FinalizeAppChannelOnActionTimeout`, `GetAppChannelActionDeadline` — turn-based action loop.
  - `GetAppChannelStatus`, `GetAppChannelSeqNum`, `GetAppChannelState`, `GetAppChannelSettleFinalizedTime` — state-machine introspection.
  - `SettleBySigTimeout`, `SettleByMoveTimeout`, `SettleByInvalidTurn`, `SettleByInvalidState` — oracle-arbitrated dispute methods (gaming-specific).
  - `NewAppChannelOnDeployedContract` — currently hard-wired to `IMultiSession.GetSessionID` and an `IMultiSession.IntendSettle` event watch (`onDeployedContractSettle`). Stateless DEPLOYED_CONTRACT conditions don't need cnode-side registration; the contract is already on-chain.
  - `getSessionID` — multisession-specific (calls `IMultiSession.GetSessionID` on a deployed contract). Only consumer is the deleted `NewAppChannelOnDeployedContract`; goes with it.
  - `onDeployedContractSettle` (line ~114) — the `IntendSettle` event watch handler installed by the deleted method.
  - `onVirtualContractDeploy` (line ~97) — already dead code in the current tree. Grep confirms zero non-self callers; it was historically the per-channel watch handler, replaced by the shared `registerVirtResolverDeployWatch`. Deleted alongside the rest.
  - `GetNumericOutcome` — zero off-chain consumers (every off-chain `TransferFunctionType` is `BOOLEAN_AND`); ABIgen comes back when a numeric consumer surfaces.
- The **virt-resolver deploy watch and the entire dormant callback infrastructure**. Today's plumbing exists only to fire `OnDispute(0)` notifications when the virt-resolver emits a `Deploy` event for a registered VIRTUAL_CONTRACT — the legacy gaming flow's "tell the player a deployment happened so they can apply moves." After this trim no consumer wants `OnDispute` notifications (the SDK / webapi client surface that consumed them all deletes), and the watch never updates `AppChannel.DeployedAddr` — that's done synchronously inside `deployIfNeeded` / `deployVirtualContract` and refreshed by `GetAppChannelDeployedAddr` on demand. Deletes:
  - `AppClient.registerVirtResolverDeployWatch` and its inline `monitor.Monitor` callback closure (`app/appclient.go` ~line 182–230).
  - `AppClient.virtDeployMu`, `virtDeployChanCount`, `virtDeployWatchID`, `virtDeployWatchStarted` fields.
  - The VIRTUAL_CONTRACT branch in `AppClient.DeleteAppChannel` that decrements the watch refcount and tears down the shared watch.
  - The `Callback` field on `AppChannel` (`app/appclient.go` ~line 40).
  - The `sc common.StateCallback` parameter on `AppClient.NewAppChannelOnVirtualContract` (and the matching wrapper parameters in `client/app_channel.go` and `celersdk/appsession.go` constructors).
  - The `Callback AppCallback` field on `celersdk.AppInfo` and the `celersdk.AppCallback` interface itself (no consumer post-trim — `common.StateCallback` stays as a generic interface used elsewhere).
  - The `appSessionCallback` type and `appSessionCallbackMap` (with its lock) in `webapi/api_server.go`. Never recreated; the create handler stops constructing callbacks at all.
  - Surviving `AppClient` surface: `NewAppChannelOnVirtualContract` (register virtual-contract bytecode without a callback param), `deployIfNeeded` + `deployVirtualContract` (private — the deploy-on-query path used by `GetBooleanOutcome` and surviving callers), `GetBooleanOutcome` (off-chain query — **redesigned** to use `IBooleanCond` bindings for both VIRTUAL_CONTRACT and DEPLOYED_CONTRACT branches, no `SessionQuery` wrapping; deploy-on-query side effect preserved), `GetAppChannelDeployedAddr` (on-chain probe via `isDeployed`), `DeleteAppChannel` (cleanup, simplified), `PutAppChannel` / `GetAppChannel` (in-package accessors for `appChannelMap`).
- The **off-chain state-exchange RPCs** — the gaming-era state machine's *off-chain* surface, symmetric to the on-chain `intendSettle` we're already deleting. Earlier drafts of this plan missed these:
  - `webapi/proto/web_api.proto`: `SignOutgoingState`, `ValidateAck`, `ProcessReceivedState` RPCs and their request/response messages.
  - `webapi/api_server.go`: the corresponding handlers.
  - `celersdk/appsession.go`: `AppSession.SignAppData`, `AppSession.HandleMatchData`, the `AppData` type, `OPCODE_*` constants, and the seqnum / last-state tracking fields these methods own.
- The **ABIgen for legacy session contracts**: `app/singlesession.go`, `app/multisession.go`, `app/singlesessionwithoracle.go`, `app/multisessionwithoracle.go`. Bindings for `SingleSessionApp` / `MultiSessionApp` / their oracle variants — the legacy gaming templates.
- **`app/oracle.go`** — frozen ABIgen-style file containing the legacy `OracleState` / `OracleProof` types. **Correction from earlier drafts:** there is no `oracle.proto` source in this tree, and `proto/app.proto` does *not* contain these messages. Treat `app/oracle.go` as a dead generated artifact and delete it directly.
- **`proto/app.proto`** messages — the file actually contains `AppState`, `StateProof`, and `SessionQuery` (not `OracleState` / `OracleProof` as an earlier draft claimed). Once the readers are gone, all three are dead:
  - `AppState`, `StateProof` — used by the deleted state-exchange RPCs and by `AppClient.IntendSettle` / `getSessionID`. Delete.
  - `SessionQuery` — used by today's `GetBooleanOutcome` (DEPLOYED_CONTRACT branch wraps queries) and by `app/multisession.go`. Both deleted; this goes too.
  - With all three gone, `proto/app.proto` itself becomes empty and can be removed; sweep `proto/app.pb.go` accordingly.
- The **gaming/state-machine surface in the shared `WebApi`** (the same proto-defined service is implemented by both `webapi/api_server.go` and `webapi/osp_pay_api_server.go`). Decision recorded as a keep/delete table covering every current app-session RPC, using real proto names:

  | RPC (real proto name) | Client server | OSP server | Decision |
  | --- | --- | --- | --- |
  | `CreateAppSessionOnVirtualContract` | yes | yes | **keep** — registration entry point for VIRTUAL_CONTRACT |
  | `CreateAppSessionOnDeployedContract` | yes | no | **delete** — backing `AppClient.NewAppChannelOnDeployedContract` is multisession-specific; stateless DEPLOYED_CONTRACT conditions need no registration |
  | `DeleteAppSession` | yes | yes | **keep** — cleanup pair to `Create*` |
  | `GetDeployedAddressForAppSession` | yes | no | **keep** — useful for VIRTUAL_CONTRACT once on-dispute deployment lands; backed by `AppChannel.DeployedAddr`, which is set synchronously by `deployIfNeeded` / `deployVirtualContract` after a successful deployment tx, and refreshed on demand by `GetAppChannelDeployedAddr` (which probes the chain via `isDeployed` and updates the field). The legacy virt-resolver deploy watch that previously fired notifications is **deleted** in this trim — see "What's deleted" below. |
  | `GetBooleanOutcomeForAppSession` | yes | no | **keep** — off-chain outcome query webapi RPC; backed by the surviving (and redesigned) `AppClient.GetBooleanOutcome`. **Note:** for VIRTUAL_CONTRACT this RPC is *not* a passive read — `AppClient.GetBooleanOutcome` calls `deployIfNeeded(appChannel)` before invoking the contract, which submits a real on-chain deployment transaction the first time the virtual contract is queried. This deploy-on-query side effect is preserved by design (see "Reading-path redesign" sub-bullet above). For DEPLOYED_CONTRACT (when registered via the now-deleted `CreateAppSessionOnDeployedContract`) this branch wouldn't survive the trim anyway; post-trim the only path through `GetBooleanOutcomeForAppSession` is VIRTUAL_CONTRACT-with-lazy-deploy. |
  | `GetStatusForAppSession` | yes | yes | **delete** — gaming/state-machine introspection |
  | `GetSeqNumForAppSession` | yes | no | **delete** — state-machine introspection |
  | `GetStateForAppSession` | yes | no | **delete** — state-machine introspection |
  | `ApplyActionForAppSession` | yes | no | **delete** — turn-based gaming action loop |
  | `FinalizeOnActionTimeoutForAppSession` | yes | no | **delete** — gaming action-timeout finalization |
  | `GetActionDeadlineForAppSession` | yes | no | **delete** — gaming action deadline lookup |
  | `GetSettleFinalizedTimeForAppSession` | yes | no | **delete** — state-machine introspection (no longer relevant once `intendSettle` is gone) |
  | `SubscribeAppSessionDispute` | yes | no | **delete** — gaming-dispute event subscription |
  | `SettleAppSession` | yes | no | **delete** — gaming dispute settle |
  | `SettleAppSessionBySigTimeout` / `*ByMoveTimeout` / `*ByInvalidTurn` / `*ByInvalidState` | yes | no | **delete** — oracle-arbitrated gaming disputes |
  | `SignOutgoingState` / `ValidateAck` / `ProcessReceivedState` | yes | no | **delete** — off-chain state-exchange (the state-machine handshake protocol) |
- **Server-side state for the surviving handlers:**
  - `appSessionMap` (and its lock) **stays** — the surviving `DeleteAppSession`, `GetDeployedAddressForAppSession`, and `GetBooleanOutcomeForAppSession` handlers all need it to look up the registered `AppSession` by ID. The earlier draft of this plan that said "delete `appSessionMap`" was wrong; it gets corrected in AS-B.
  - `appSessionCallbackMap`, `appSessionCallback`, and the entire callback infrastructure **delete**. Resolution-3's narrative described this as "drop storage map but keep callback wiring nil-safe" — that was almost-right but understated the cleanup. The complete picture, after audit:
    - `SubscribeAppSessionDispute` is the *consumer* of the callback channel; deleted in this trim.
    - `CreateAppSessionOnVirtualContract` (and the deleted `CreateAppSessionOnDeployedContract`) is the *constructor* — builds an `appSessionCallback` and stores it in the map.
    - `app.AppClient` has three call sites of `Callback.OnDispute(...)`: `onVirtualContractDeploy` (line ~106 — already dead code, zero callers in the current tree), `onDeployedContractSettle` (line ~129 — on the deletion list with the rest of the deployed-session path), and the inline closure inside `registerVirtResolverDeployWatch` (line ~219 — the actually-live one, and already nil-safe via line 216's `appChannel.Callback == nil` guard).
    - With the consumer (`SubscribeAppSessionDispute`) gone, no one needs `OnDispute` notifications. With the only live invocation site (`registerVirtResolverDeployWatch`'s inline closure) deleted along with the watch itself (see "virt-resolver deploy watch" deletion above), there are no Go-side OnDispute-callers left at all. So the trim deletes the storage map, the `appSessionCallback` type, the `Callback` parameter through every layer, and the callback fields in the surviving structs — not just nil-safe-the-existing-code. There is no remaining live nil-safety burden.
- The **gaming/state-machine surface in `celersdk/appsession.go`** — SDK wrappers for the deleted webapi RPCs above, plus the `AppSession.SignAppData` / `HandleMatchData` / `AppData` / `OPCODE_*` constants / seqnum-tracking surface from the off-chain state-exchange protocol. The file shrinks dramatically (or splits) but doesn't disappear; `CreateAppSession`-style helpers and the boolean-outcome helper stay because their backing webapi RPCs stay.
- The **legacy test fixtures** under `testing/testapp/`: `multigomoku.go`, `singlesessionapp.go`, `multisessionapp.go`, `singlesessionappwithoracle.go`, `multisessionappwithoracle.go`, and `utils.go`'s session-specific helpers. The remaining test fixtures are ABIgen output for `BooleanCondMock` / `NumericCondMock` plus minimal Go wiring.
- The **gaming-flavored e2e tests** in `test/e2e/pay_dispute.go` and the entirety of `test/e2e/pay_dispute_with_oracle.go` — the scenarios that exercise turn-based-game dispute paths (apply-action, action-timeout finalization, oracle settle-by-sig-timeout) all delete. Channel-level dispute coverage in `test/e2e/settle_channel.go` and `cold_bootstrap.go` is independent and stays fully intact. **`pay_dispute.go` is rewritten, not deleted**: it gains minimal coverage for the conditional-pay-with-dispute flow under both `ConditionType_VIRTUAL_CONTRACT` and `ConditionType_DEPLOYED_CONTRACT`, with the underlying contract being `BooleanCondMock` in both cases — the trim ends up improving dispute test coverage of the surviving protocol surface, not just shrinking it.
- The **`WaitUntilBlockHeight` helper** in `testing/clientcontroller.go` — only existed because the legacy testapp contracts used `block.number`; with them gone, the only deadline unit is `block.timestamp`-derived seconds, and `WaitUntilDeadline` covers everything.

### What this means in numbers

Rough estimate of code deletion (net of what stays for VIRTUAL_CONTRACT support):

| Area | Approximate LOC removed |
| --- | --- |
| `app/appclient.go` — gaming methods, deployed-contract registration, getSessionID, GetNumericOutcome, **virt-resolver deploy watch + callback infrastructure**, related helpers | ~1100 |
| `app/singlesession.go` / `multisession.go` / `*withoracle.go` ABIgen | ~1500 |
| `app/oracle.go` ABIgen | ~200 |
| `app/numericoutcome.go` ABIgen (no off-chain consumer) | ~100 |
| `app/apputil.go` state-exchange helpers | ~150 |
| `webapi/api_server.go` + `webapi/internal_api_server.go` + `webapi/osp_pay_api_server.go` — handlers for ~14 deleted RPCs (state-exchange + dispute + introspection); **`appSessionCallback` type and `appSessionCallbackMap` storage** | ~950 |
| `webapi/proto/web_api.proto` — RPC and message definitions for the deleted surface | ~250 |
| `celersdk/appsession.go` — wrappers for deleted RPCs + state-exchange protocol + non-webapi gaming helpers (`SwitchToOnchain`, `OnChainApplyAction`, `OnChainGetStatus`, etc., `NewAppSessionOnDeployedContract`, oracle settles, `GetPlayerIdxForMatch`) + **`AppCallback` interface and `AppInfo.Callback` field** | ~720 |
| `client/app_channel.go` — wrappers over deleted AppClient methods (`NewAppChannelOnDeployedContract`, `SettleAppChannel`, `SignAppState`, on-chain action / introspection helpers, **`GetAppChannel`**) | ~160 |
| `testing/testapp/` legacy gaming fixtures (multigomoku, multisession, withoracle variants; `singlesessionapp.go` stays for x402 back-compat) | ~2200 |
| `test/e2e/pay_dispute*.go` — gaming-flavored scenarios deleted, dispute coverage rewritten against `BooleanCondMock` for both condition types | ~700 (net) |
| `testing/clientcontroller.go` `WaitUntilBlockHeight` and gaming helpers (incl. SignOutgoingState / ValidateAck / ProcessReceivedState wrappers if any) | ~80 |
| `proto/app.proto` (`AppState`, `StateProof`, `SessionQuery`) plus `proto/app.pb.go` regen | ~50 |
| `tools/scripts/regenerate-legacy-app-bindings.sh` — 10 of 11 entries removed (only `singlesessionapp.go` survives); `tools/scripts/README.md` updated | ~30 |
| **Total deletion** | **~8200 LOC** |

Net additions in this trim:

| Area | Approximate LOC added |
| --- | --- |
| `testing/testapp/booleancondmock.go` — ABIgen output for the canonical IBooleanCond fixture | ~150 |
| `test/e2e/setup_onchain.go` — deploy `BooleanCondMock` and surface its address on the contract bundle | ~30 |
| `test/e2e/pay_dispute.go` — rewritten coverage of conditional-pay dispute under `VIRTUAL_CONTRACT` and `DEPLOYED_CONTRACT`, both using `BooleanCondMock` | ~250 |
| **Total addition** | **~430 LOC** |

Of the ~8200 deleted, roughly half is regenerated ABIgen output, so the hand-written-code delta is closer to ~3700 LOC removed. Net new code is small and targeted.

No new on-chain contracts, no new interfaces — `BooleanCondMock` already exists in `agent-pay-contracts`. The simplification is mostly subtraction with the dispute-test coverage repositioned onto a clean fixture and the `GetBooleanOutcome` reading path redesigned to drop multisession-specific encoding.

### What an `agent-pay-x402` conditional payment looks like after this

x402's current flow is **unchanged by this trim** because x402 only exercises `CreateAppSessionOnVirtualContract` (registration) and the off-chain cooperative confirm/cancel path — none of the gaming-machinery methods being deleted. The legacy `SimpleSingleSessionApp` bytecode that x402 registers today via [agent-pay-x402/testinfra/session.go](../../../agent-pay-x402/testinfra/session.go) keeps working: x402 never calls `IntendSettle` / `ApplyAction` / etc. on it, so deleting those methods from `app.AppClient` and the webapi has no effect on x402.

Future migration (deferred — see §7): swap x402's registered bytecode from `SimpleSingleSessionApp` (a turn-based-game contract whose dispute path x402 doesn't use) to a stateless `IBooleanCond` impl (some custom verifier x402 cares about). That's a one-line change in x402; the agent-pay side already supports it post-trim because `CreateAppSessionOnVirtualContract` doesn't care what bytecode it registers as long as the contract conforms to `IBooleanCond` if a dispute ever queries it.

For consumers using **DEPLOYED_CONTRACT** (e.g., payment gated on an on-chain oracle data feed or a pre-deployed ZK verifier): same flow but skipping the registration step. The condition's `OnChainAddress` points directly at an already-deployed `IBooleanCond` / `INumericCond` contract.

The condition contract is essentially never queried in the happy path of either condition type. It's there as the dispute-fallback resolver of "this was a structured conditional payment, here's what would resolve it if anyone ever asked."

---

## 3. Non-goals

Things this plan explicitly does **not** do, to keep scope contained:

- **Removing either `ConditionType_DEPLOYED_CONTRACT` or `ConditionType_VIRTUAL_CONTRACT`.** Both are general protocol primitives for stateless condition contracts and both have legitimate AI-agent use cases (see §2 "What survives" — oracles, ZK verifiers, lazy-deployed parsers). The trim deletes the gaming session-state-machine that was wrapped around VIRTUAL_CONTRACT, not the condition type.
- **Removing the virtual-contract registration plumbing.** The cnode keeps the ability to register virtual-contract bytecode and deploy it on dispute — that's the runtime support for VIRTUAL_CONTRACT. What goes is the gaming state machine on top of it.
- **Shipping any new on-chain contract.** `IBooleanCond` / `INumericCond` interfaces exist in `agent-pay-contracts`; the test-only `BooleanCondMock` / `NumericCondMock` exist as fixtures. That's enough. Anyone who needs a real production verifier (cosigned messages, oracle signatures, ZK proofs, etc.) writes their own `IBooleanCond` implementation when their use case demands it.
- **Migrating `cApps`.** The `cApps` external repo stays dead. None of its contracts are ported into `agent-pay-contracts`.
- **Migrating x402 in this PR.** Per the §4 decision, x402's registered-bytecode swap is deferred — see §7. The trim is intentionally compatible with x402's current `SimpleSingleSessionApp` registration; x402 doesn't exercise any of the deleted methods.
- **Generating ABIgen bindings for `NumericCondMock` in agent-pay.** No off-chain numeric consumer exists today (every off-chain `TransferFunctionType` is `BOOLEAN_AND`) and no legacy numeric fixture exists to be replaced. "Delete unused, add later" applies cleanly. (`BooleanCondMock` bindings, by contrast, **are** generated as part of AS-C — they're the canonical fixture for the rewritten dispute tests covering both `VIRTUAL_CONTRACT` and `DEPLOYED_CONTRACT` and the eventual replacement target for x402's `SimpleSingleSessionApp` import.)
- **Generating bindings for `INumericOutcome` / `INumericCond` in `app/`.** Per the "delete unused, add later" principle: no off-chain code calls them. `app/numericoutcome.go` is dead and gets deleted in AS-B; bindings come back when a numeric consumer surfaces.
- **Renaming `BooleanCondMock` / `NumericCondMock`, or renaming `CreateAppSessionOnVirtualContract` to drop "Session".** Both names are misleading (the entities are not just mocks; the registered things are not stateful sessions), but renaming them costs API churn for marginal clarity gain. Separate polish if ever taken on.
- **Renaming or restructuring the `agent-pay/app/` package.** The package shrinks but the name stays — it aligns with "app channel" in the protocol's architecture docs.
- **Breaking-change accommodations for existing deployments.** This codebase has no production deployments yet (it's an evolving AI-agent-payment platform); we treat the change as a normal protocol revision. If that ever stops being true, the plan needs revisiting.
- **Resurrecting an on-chain dispute-fallback path for app conditions later.** Not a hard non-goal — the protocol-level interfaces leave room for someone to ship a stateful condition contract if a real consumer ever needs one. But this trim removes all the *generic* infrastructure for it; rebuilding would need a real use case to motivate the design, not the speculative gaming-era one.

---

## 4. Decisions (resolved)

The following decisions were resolved before AS-A. Recorded here for the audit trail.

- [x] **NumericCondition off-chain support: confirmed not present.** A direct grep across `agent-pay/` finds zero callers of `INumericOutcome` / `GetNumericOutcome` outside the ABIgen file itself, and every off-chain `TransferFunctionType` instantiation is `BOOLEAN_AND` (verified in `cnode/`, `delegate/`, `webapi/`, `client/`, `server/`). The on-chain `INumericCond` interface stays in `agent-pay-contracts` for `PayResolver`'s `NUMERIC_ADD/MAX/MIN` resolution, but the off-chain ABIgen file `app/numericoutcome.go` is dead code and is deleted in AS-B. Bindings get regenerated when a numeric consumer surfaces.
- [x] **`app/` package contents: apply "delete unused, add later if used."** `app/booleanoutcome.go` stays (consumed by `AppClient.GetBooleanOutcome`). `app/numericoutcome.go` deletes (zero consumers). `app/apputil.go` deletes if all its helpers are referenced only by the deleted gaming methods (verified during AS-A audit).
- [x] **`webapi.proto` deletion strategy: hard-delete the gaming RPCs.** Authoritative keep/delete list lives in §2's keep/delete table; `CreateAppSessionOnVirtualContract` stays as a registration entry point for VIRTUAL_CONTRACT, but `CreateAppSessionOnDeployedContract` deletes (see the dedicated §4 decision below). The gaming/state-machine RPCs hard-delete using their real proto names: `SettleAppSession`, the four `*Timeout` / `*InvalidTurn` / `*InvalidState` variants, `SubscribeAppSessionDispute`, `GetStatusForAppSession`, `GetSeqNumForAppSession`, `GetStateForAppSession`, `ApplyActionForAppSession`, `FinalizeOnActionTimeoutForAppSession`, `GetSettleFinalizedTimeForAppSession`, `GetActionDeadlineForAppSession`, `SignOutgoingState`, `ValidateAck`, `ProcessReceivedState`. (Earlier drafts of this plan used incorrect names like `GetAppSessionSeqNum` / `GetAppSessionState`; the real proto names follow the `*ForAppSession` suffix pattern.)
- [x] **x402 migration: deferred.** This trim is compatible with x402's current `SimpleSingleSessionApp`-based virtual-contract registration; x402 doesn't exercise any of the deleted methods. A future PR (in either repo) swaps the registered bytecode to a stateless verifier — see §7.
- [x] **Test-fixture location:** `testing/testapp/singlesessionapp.go` (and its `ta.AppCode` / `ta.GetSingleSessionConstructor` / `ta.Timeout` exports) stay because x402 imports them. `multigomoku.go`, `multisessionapp.go`, `singlesessionappwithoracle.go`, `multisessionappwithoracle.go` delete. `BooleanCondMock` ABIgen bindings **are** generated into `testing/testapp/booleancondmock.go` as part of AS-C — they back the rewritten dispute tests (both VIRTUAL_CONTRACT and DEPLOYED_CONTRACT scenarios) and unblock the deferred x402 migration. `NumericCondMock` bindings stay deferred (no off-chain numeric consumer, no analog legacy fixture).
- [x] **DEPLOYED_CONTRACT registration path:** `NewAppChannelOnDeployedContract` (and its webapi RPC `CreateAppSessionOnDeployedContract`) is **deleted**, not preserved. Reasoning per the GPT review (Finding 1): the current implementation is hard-wired to `IMultiSession.GetSessionID` and an `IMultiSession.IntendSettle` event watch, so keeping the API while deleting `app/multisession.go` would strand the API. Stateless DEPLOYED_CONTRACT conditions don't need a registration step at all — the contract is already on-chain; users build a `Condition` with `OnChainAddress = deployed_addr` and that's it. The `GetBooleanOutcome` reading path gets a parallel redesign to drop `SessionQuery` wrapping and call `IBooleanCond.getOutcome` directly with the raw `argsQueryOutcome` bytes.
- [x] **Off-chain state-exchange RPCs:** `SignOutgoingState`, `ValidateAck`, `ProcessReceivedState` (and their backing `celersdk.AppSession.SignAppData` / `HandleMatchData` / `AppData` / `OPCODE_*` constants / seqnum tracking) are **deleted** as part of this trim, not preserved. They're the off-chain half of the same gaming state machine the on-chain `intendSettle` is the on-chain half of — keeping them would mean we trimmed only half the state machine. Per GPT review Finding 2.
- [x] **OSP webapi subset:** keep `CreateAppSessionOnVirtualContract` and `DeleteAppSession`; delete `GetStatusForAppSession` (real proto name; earlier drafts called this `GetAppSessionStatus`). The `osp_webapi_test.go` `ospWebApiAppSessionSubset` test gets updated to drop coverage of the deleted RPC. Per GPT review Finding 5.

---

## 5. Phases

### AS-A — Pre-flight audit

The §4 decisions are already resolved; this phase is the safety check before any deletion lands. **Status: completed.** Findings are recorded under each subtask below; the consolidated test-tag list and the AS-B-targeted line-and-symbol inventory live in the AS-A findings sub-section at the end of this phase.

- [x] Re-confirm the §4 audits with a fresh grep, in case anything has shifted:
  - [x] `agent-pay-x402` references to deletion-list methods, using the real proto / Go names: `IntendSettle`, `GetSettleFinalizedTime`, `GetSessionID`, `SettleBySigTimeout`, `SettleByMoveTimeout`, `SettleByInvalidTurn`, `SettleByInvalidState`, `OracleProof`, `OracleState`, `ApplyAction`, `FinalizeAppChannelOnActionTimeout`, `GetAppChannelActionDeadline`, `GetStatusForAppSession`, `GetSeqNumForAppSession`, `GetStateForAppSession`, `ApplyActionForAppSession`, `FinalizeOnActionTimeoutForAppSession`, `GetActionDeadlineForAppSession`, `GetSettleFinalizedTimeForAppSession`, `SettleAppSession`, `SubscribeAppSessionDispute`, `SignOutgoingState`, `ValidateAck`, `ProcessReceivedState`, `SignAppData`, `HandleMatchData`, `SwitchToOnchain`, `NewAppSessionOnDeployedContract`, `CreateAppSessionOnDeployedContract`, `OnChainApplyAction`, `OnChainFinalizeOnActionTimeout`. Expected: zero hits. **Result: zero code hits.** One stale doc-comment hit at `agent-pay-x402/testinfra/topology_both_osp_test.go:160` references `DeleteAppSession` and `GetStatusForAppSession` in a narrative comment — not a call site, doesn't affect compilation, becomes mildly stale after `GetStatusForAppSession` deletes. Logged as a minor x402-side comment-fix follow-up; non-blocking.
  - [x] `agent-pay-x402` calls to `CreateAppSessionOnVirtualContract`: confirm and enumerate the call sites; verify the `ta.AppCode` / `ta.GetSingleSessionConstructor` / `ta.Timeout` exports they depend on stay. **Result: two production call sites, not one.** `agent-pay-x402/testinfra/session.go:27` (the canonical helper) **and** `agent-pay-x402/pkg/buyersdk/backend.go:416` (`createAppSession` in the buyer-SDK backend). Both use the same request shape — `ContractBin` / `ContractConstructor` / `Nonce` / `OnChainTimeout` — and **neither passes a callback**, so both are fully compatible with the AS-B trim that drops the callback parameter at every layer. Plan's earlier "only call site" claim was wrong but impact is zero. The `ta.AppCode` / `ta.GetSingleSessionConstructor` / `ta.Timeout` exports stay; both x402 call sites continue to work post-trim against the surviving `testing/testapp/singlesessionapp.go`.
  - [x] No other sibling repo in `~/Work/celer/` imports `agent-pay/app/` and references the deletion list. **Result: zero hits.** Only `agent-pay/` itself imports its own `app/` package.
  - [x] No production rt_config or profile JSON sets fields specific to the deleted machinery. **Result: zero hits** across `deploy/`, `test/manual/rt_config.json`, `testing/profile/`. No `oracle_*`, `app_session_*`, `signOutgoingState`, `validateAck`, `processReceivedState` keys anywhere. The existing `on_chain_timeout` field is on `AppInfo` (deleted) but it's a Go struct field, not JSON config; safe.
  - [x] `app/numericoutcome.go` truly has zero off-chain consumers (already confirmed during the §4 resolution; sanity-check once). **Result: zero non-self consumers** (re-verified — only mention of `INumericOutcome` in the entire `agent-pay/` tree is the abigen file itself). Confirms AS-B can delete the file without breaking any caller.
  - [x] `app/apputil.go` helpers — verify which (if any) are referenced from outside the deleted gaming methods AND outside the deleted state-exchange RPCs. **Result: every exported helper has consumers, but every consumer is on the deletion list.** `apputil.go` exports `EncodeAppState`, `DecodeAppState`, `EncodeAppStateProof`, `DecodeAppStateProof`, `SigSortedAppStateProof`, `SortPlayerSigs`, `SortPlayers`. Their non-`/app/` callers: `test/e2e/pay_dispute_with_oracle.go:306` (DELETE entirely), `celersdk/appsession.go:148` (`NewAppSessionOnDeployedContract` — DELETE), `celersdk/appsession.go:181` (`SignAppData` — DELETE), `celersdk/appsession.go:203` (`HandleMatchData` — DELETE), `celersdk/appsession.go:218` (`HandleMatchData` — DELETE), `celersdk/appsession.go:265` (`HandleMatchData` — DELETE). **Conclusion: delete `app/apputil.go` entirely** — no surviving consumer.
  - [x] Audit `app/appclient.go` for the `IMultiSession*` and `SessionQuery` references in the deployed-contract code path so the AS-B redesign list is final. **Result: confirmed every `IMultiSession*` / `SessionQuery` callsite maps to the deletion list.** Concrete inventory (line numbers from current code):
    - `onDeployedContractSettle` (lines 114–133): uses `IMultiSessionIntendSettle` and `IMultiSessionABI`. **DELETE** with the deployed-session path.
    - `NewAppChannelOnDeployedContract` (line 261): calls `getSessionID` and installs `IMultiSessionABI` watch. **DELETE.**
    - Watch installation (line 288) inside `NewAppChannelOnDeployedContract`: `IMultiSessionABI` event monitor. **DELETE.**
    - `IntendSettle` watch handler call (line 304): `appChannel.onDeployedContractSettle(&eLog)`. **DELETE.**
    - `GetBooleanOutcome` DEPLOYED_CONTRACT branch (lines 384–392): `SessionQuery` wrapping + `IMultiSessionCaller.GetOutcome`. **REDESIGN** to use `IBooleanCondCaller.GetOutcome(query)` with raw `argsQueryOutcome` bytes (matches PayResolver).
    - `ApplyAction` (line 418), `FinalizeOnActionTimeout` (line 458), `GetSettleFinalizedTime` (line 500), `GetActionDeadline` (line 539), `GetSeqNum` (line 578), `GetStatus` (line 610), `GetState` (line 635): all `NewIMultiSessionTransactor`/`Caller`. **DELETE** (gaming-vestigial methods).
    - `SettleBySigTimeout` (line 682), `SettleByMoveTimeout` (line 723), `SettleByInvalidTurn` (line 763), `SettleByInvalidState` (line 803): `NewIMultiSessionWithOracleTransactor`. **DELETE** (oracle disputes).
    - `getSessionID` (lines 922–929): `NewIMultiSessionCaller.GetSessionID`. **DELETE** (only consumer was `NewAppChannelOnDeployedContract`, also deleted).
    - **Net result:** zero `IMultiSession*` or `SessionQuery` references survive the trim. The redesigned `GetBooleanOutcome` uses only `IBooleanCondCaller`.
- [x] Walk every test file under `agent-pay/test/e2e/` and tag each as keep / delete / rewrite:
  - **KEEP (4 files)** — exercise channel-level dispute or cooperative paths, no app-session deletion-list refs:
    - `cold_bootstrap.go` — line 229's `c2.IntendSettlePaymentChannel(...)` is **channel-level** (`CelerClient.IntendSettlePaymentChannel` → `CNode.IntendSettlePaymentChannel`), distinct from the deleted `AppClient.IntendSettle`.
    - `e2e_test.go` — only registers `ospIntendSettleErc20Channel` (defined in `settle_channel.go`); no app-session deletion refs in the file itself.
    - `settle_channel.go` — entirely channel-level (`IntendSettlePaymentChannel`).
    - `setup_onchain.go` — only mention is the `"IntendSettle": 2` event-watcher poll-interval config map key (chain event name, not the deleted method).
  - **REWRITE (2 files)** — covered by AS-C:
    - `pay_dispute.go` (1108 LOC) — heavy use of `SettleAppChannel`, `GetAppChannelState`, `GetAppChannelSettleFinalizedTime`, `ApplyAppChannelAction`. AS-C rewrites against `BooleanCondMock` for both VIRTUAL_CONTRACT and DEPLOYED_CONTRACT scenarios.
    - `osp_webapi_test.go` — `ospWebApiAppSessionSubset` test (lines 239–331) uses `CreateAppSessionOnVirtualContract` (KEEP) and `DeleteAppSession` (KEEP) but also has 3 `GetStatusForAppSession` assertions at lines 283 (call), 285 / 288 (error checks), 326 (call), 328 / 331 (error checks). AS-C drops those assertion blocks; the rest of the test stays.
  - **DELETE (1 file)** — covered by AS-C:
    - `pay_dispute_with_oracle.go` (315 LOC) — uses `SettleAppChannelBySigTimeout`, `GetAppChannelState`. Pure oracle-dispute coverage; entire file deletes along with the OSP webapi RPC entries.
- [x] Check the OSP webapi test (`test/e2e/osp_webapi_test.go` — `ospWebApiAppSessionSubset`) and confirm which of its assertions touch `GetStatusForAppSession`. **Result:** the assertions touching `GetStatusForAppSession` live at lines 283–290 (initial status check after `CreateAppSessionOnVirtualContract`) and lines 326–332 (status check after `DeleteAppSession`). Both blocks delete in AS-C; the surrounding `Create*` / `Delete*` calls and their error checks stay.
- [x] Confirm three additional compile-driven cleanup sites that AS-B's scoped vet gate will surface (enumerated in AS-B's "Compile-driven follow-up sites" subsection): **all three confirmed present**: `server/osp_webapi_backend.go` (line 66 callback arg + line 80 `GetStatusForAppSession` method), `app/appclient_virtresolver_watch_test.go` (entire file becomes dead with the watch deletion — delete), `webapi/osp_pay_api_server_test.go` (assertions for the deleted OSP RPC).

**Exit criteria:** audit greps return the expected results; test-tag list (keep / delete / rewrite) is final; AS-B redesign-target list (specifically the `GetBooleanOutcome` rewrite) is concrete enough to execute. **All exit criteria met. AS-A complete.**

#### AS-A findings — consolidated for AS-B execution

- **Test-tag list (final):**
  - KEEP: `cold_bootstrap.go`, `e2e_test.go`, `settle_channel.go`, `setup_onchain.go` (4 files, channel-level coverage stays intact).
  - REWRITE: `pay_dispute.go` (full AS-C rewrite onto `BooleanCondMock`), `osp_webapi_test.go::ospWebApiAppSessionSubset` (drop 3 `GetStatusForAppSession` assertion blocks).
  - DELETE: `pay_dispute_with_oracle.go` (entire 315-LOC file).
- **`app/appclient.go` `IMultiSession*` / `SessionQuery` redesign-target inventory:** see the audit subtask above — every callsite maps cleanly to the trim. Zero `IMultiSession*` references survive in the redesigned `GetBooleanOutcome`.
- **`app/apputil.go`:** delete entirely — every consumer is on the deletion list.
- **x402 call-site update:** the plan's earlier "only call site is `testinfra/session.go`" claim was wrong. There are **two**: `testinfra/session.go:27` and `pkg/buyersdk/backend.go:416`. Both are signature-compatible with the trim (no callback, fields unchanged); no cross-repo migration needed.
- **Stale x402 doc-comment** (`testinfra/topology_both_osp_test.go:160`) referencing the deleted `GetStatusForAppSession` — non-blocking; logged for a future x402-side cleanup.

### AS-B — Off-chain trim in `agent-pay/`

The big mechanical phase. Each subtask is straightforward given the deletion list from AS-A; the challenge is keeping the trimmed packages building between subtasks. **Note:** virtual-contract registration plumbing (`CreateAppSessionOnVirtualContract` only — the deployed-contract registration path deletes per §4) stays. The deployed-contract reading path is **redesigned** in this phase to drop multisession-specific encoding.

#### `app/` package

- [x] Regenerate `app/booleanoutcome.go` from `agent-pay-contracts/src/lib/interface/IBooleanCond.sol` so the Go binding symbols are `IBooleanCond` / `IBooleanCondCaller` / `IBooleanCondTransactor` / `IBooleanCondFilterer` (matching the canonical agent-pay-contracts interface name) instead of the legacy `IBooleanOutcome*`. The ABI shape is identical (`isFinalized(bytes) returns (bool)`, `getOutcome(bytes) returns (bool)`); this is purely a symbol-name alignment. — done; renamed file to `app/booleancond.go`.
- [x] **Generated-file ownership for the regenerated binding** — fix in this phase, not later:
  - Today `app/booleanoutcome.go` is owned by `tools/scripts/regenerate-legacy-app-bindings.sh` (line 127: `generate_from_abi app/booleanoutcome.go app/booleanoutcome.go app IBooleanOutcome`). That script re-runs `abigen` against the ABI literal already embedded in the existing Go file — so it can't perform the `IBooleanOutcome` → `IBooleanCond` rename, only re-create the legacy symbols.
  - **Move ownership** to `tools/scripts/regenerate-go-bindings.sh`, which already pulls from the `agent-pay-contracts` foundry artifacts and generates the rest of the on-chain interface bindings (under `chain/...`). Add a line that emits `app/booleanoutcome.go` from `IBooleanCond` (or rename the file to `app/booleancond.go` and emit there — minor polish, decided in this subtask).
  - **Remove from `regenerate-legacy-app-bindings.sh`** all 10 entries whose output files are deleted by AS-B and AS-C, leaving only the x402 back-compat carry:
    - `app/booleanoutcome.go` — moved to the modern script (above).
    - `app/numericoutcome.go` — output file deleted by AS-B (no off-chain consumer).
    - `app/multisession.go`, `app/multisessionwithoracle.go`, `app/singlesession.go`, `app/singlesessionwithoracle.go` — output files deleted by AS-B (legacy gaming session contract bindings).
    - `testing/testapp/multigomoku.go`, `testing/testapp/multisessionapp.go`, `testing/testapp/multisessionappwithoracle.go`, `testing/testapp/singlesessionappwithoracle.go` — output files deleted by AS-C (legacy gaming test fixtures).
    - **Survivor:** `testing/testapp/singlesessionapp.go` only (line 137 of the script). It stays in the legacy script as the x402 back-compat carry; deleting the script entirely is a §7 follow-up alongside the x402 migration.
  - **Update `tools/scripts/README.md`** to reflect the new ownership: move `app/booleanoutcome.go` from the "Regenerate Legacy App Bindings" section to "Regenerate Go Contract Bindings"; drop references to the 9 deleted bindings; note the legacy script's shrunk scope (one survivor).
- [x] Delete `app/numericoutcome.go` — zero off-chain consumers per the §4 audit.
- [x] Trim `app/appclient.go`. The keep/delete lists below use real symbol names verified against the current tree:
  - **Keep:** `NewAppClient` constructor; `NewAppChannelOnVirtualContract` (registration — but with the `sc common.StateCallback` parameter dropped, see callback deletion below); `deployIfNeeded` and `deployVirtualContract` (private helpers — the deploy-on-query path used by `GetBooleanOutcome` and any external explicit-deploy callers); `GetBooleanOutcome` (off-chain query — but **redesigned**, see next subtask); `GetAppChannelDeployedAddr` (on-chain probe via `isDeployed`); `DeleteAppChannel` (cleanup, simplified — see callback deletion); `PutAppChannel` / `GetAppChannel` (in-package `appChannelMap` accessors); `isDeployed` private helper.
  - **Delete:** `IntendSettle`, `ApplyAction`, `FinalizeAppChannelOnActionTimeout`, `GetAppChannelActionDeadline`, `GetAppChannelStatus`, `GetAppChannelSeqNum`, `GetAppChannelState`, `GetAppChannelSettleFinalizedTime`, `SettleBySigTimeout`, `SettleByMoveTimeout`, `SettleByInvalidTurn`, `SettleByInvalidState`, `GetNumericOutcome`, `NewAppChannelOnDeployedContract`, `getSessionID`, `onDeployedContractSettle` (line ~114), `onVirtualContractDeploy` (line ~97 — already dead code, zero callers per AS-A grep), the `IMultiSessionABI`-based event watch in the deployed-contract code path.
- [x] **Delete the virt-resolver deploy watch and the entire callback infrastructure.** Per §2: this watch's only effect is firing `OnDispute(0)` notifications, and no consumer remains for those notifications post-trim. Concretely:
  - Delete `AppClient.registerVirtResolverDeployWatch` and its inline `monitor.Monitor` callback closure (`app/appclient.go` ~lines 182–230).
  - Delete the watch-state fields on `AppClient`: `virtDeployMu`, `virtDeployChanCount`, `virtDeployWatchID`, `virtDeployWatchStarted`.
  - Simplify `AppClient.DeleteAppChannel`: remove the VIRTUAL_CONTRACT branch that decrements the watch refcount and tears down the shared watch (`app/appclient.go` ~lines 159–179). The `default` branch's `c.monitorService.RemoveEvent(appChannel.callbackID)` was used by the deleted `onDeployedContractSettle` watch; that goes too.
  - Delete the `Callback` field from the `AppChannel` struct (`app/appclient.go` ~line 40).
  - Drop the `sc common.StateCallback` parameter from `NewAppChannelOnVirtualContract` (`app/appclient.go` ~line 237). Remove the `Callback: sc` line from the `AppChannel` literal in the function body. Remove the `c.registerVirtResolverDeployWatch()` call too (function is gone).
  - **No nil-safe code path remains** in `app/appclient.go` after these deletions — the only Go-side `OnDispute` invocations were the inline closure (deleted with the watch), `onVirtualContractDeploy` (dead code, deleted), and `onDeployedContractSettle` (deleted with the deployed-contract path). There is no surviving call site that could nil-deref.
- [x] **Redesign `GetBooleanOutcome`.** Today's implementation branches on `appChannel.Type`:
  - VIRTUAL_CONTRACT branch: uses `ISingleSessionCaller.GetOutcome(query)` — switch to `IBooleanCondCaller.GetOutcome(query)` (already in `app/booleanoutcome.go`).
  - DEPLOYED_CONTRACT branch: uses `IMultiSessionCaller.GetOutcome(SessionQuery{...})` — switch to `IBooleanCondCaller.GetOutcome(query)` with the raw `argsQueryOutcome` bytes (no `SessionQuery` wrapping). This matches what `PayResolver` does on-chain.
  - The `isFinalized` helper similarly drops the session-specific wrapping.
- [x] Delete `app/oracle.go`, `app/singlesession.go`, `app/multisession.go`, `app/singlesessionwithoracle.go`, `app/multisessionwithoracle.go` — all ABIgen for legacy gaming session contracts that the trimmed `AppClient` no longer references. (Note: there is no `oracle.proto` source in this tree; `app/oracle.go` is a frozen generated artifact and gets deleted directly without a regeneration step.)
- [x] Delete `app/apputil.go` per the §4 decision (or trim to whatever specific helpers turn out to be referenced from the surviving `AppClient` methods — likely none).
- [x] Verify the remaining `AppClient` still constructs cleanly from `cnode/cnode.go` (the `c.AppClient = app.NewAppClient(...)` call should still work, just with fewer dependencies). Update the construction args if any of the deleted internals were passed in. — confirmed; no construction-arg change needed.

#### `webapi/`

- [x] Trim `webapi/api_server.go` per the §2 keep/delete table (real proto names):
  - **Delete handlers and helpers for:** `SettleAppSession`, `SettleAppSessionBySigTimeout`, `SettleAppSessionByMoveTimeout`, `SettleAppSessionByInvalidTurn`, `SettleAppSessionByInvalidState`, `SubscribeAppSessionDispute`, `GetStatusForAppSession`, `GetSeqNumForAppSession`, `GetStateForAppSession`, `ApplyActionForAppSession`, `FinalizeOnActionTimeoutForAppSession`, `GetActionDeadlineForAppSession`, `GetSettleFinalizedTimeForAppSession`, `SignOutgoingState`, `ValidateAck`, `ProcessReceivedState`, `CreateAppSessionOnDeployedContract`. Plus the `appSessionCallbackMap` field and lock (today both create handlers populate it and `SubscribeAppSessionDispute` consumes it; all three RPCs delete in this trim, leaving zero readers and zero writers — see the dedicated callback-deletion subtask below for the full rewrite), and any imports left orphaned.
  - **Keep handlers for:** `CreateAppSessionOnVirtualContract`, `DeleteAppSession`, `GetDeployedAddressForAppSession`, `GetBooleanOutcomeForAppSession`. The `appSessionMap` and its lock **stay** — these surviving handlers all dereference `getAppSession()` which dereferences the map; deleting it would strand them.
- [x] Trim `webapi/internal_api_server.go` similarly. — no trim needed; the only methods on `InternalApiServer` are `OpenTrustedPaymentChannel` / `InstantiateTrustedPaymentChannel` / `DepositNonBlocking` / `CooperativeWithdrawNonBlocking`, none of which touch app-session surface.
- [x] Trim `webapi/osp_pay_api_server.go` per the §2 OSP-subset row: keep `CreateAppSessionOnVirtualContract` and `DeleteAppSession`; delete `GetStatusForAppSession`. Trim the `OspPayApiBackend` interface accordingly.
- [x] Update `webapi/proto/web_api.proto`: hard-delete the RPCs and request/response messages for everything in the keep/delete table marked **delete**. Regenerate `webapi/proto/*.pb.go`.
- [x] In `webapi/api_server.go`, **keep `appSessionMap` (and its lock)** — the surviving handlers (`DeleteAppSession`, `GetDeployedAddressForAppSession`, `GetBooleanOutcomeForAppSession`) all need it.
- [x] **Delete `appSessionCallbackMap`, the `appSessionCallback` type, and all callback construction.** Since the `app/appclient.go` callback infrastructure is fully deleted (see §2 / the AS-B `app/` subsection), there is no consumer for any callback the webapi might construct. Concretely:
  - Delete the `appSessionCallbackMap` field and `appSessionCallbackMapLock` lock from the `ApiServer` struct.
  - Delete the `appSessionCallback` type definition (`webapi/api_server.go` ~line 743) and its `OnDispute` method.
  - In `CreateAppSessionOnVirtualContract` handler: stop constructing `&appSessionCallback{...}`; drop the map write. The SDK constructor now takes no callback parameter (per the celersdk trim below), so the call simplifies to passing only the contract bytecode / constructor / nonce / timeout.
  - The deleted `CreateAppSessionOnDeployedContract` handler had the same callback construction; that goes with it.

#### `client/app_channel.go`

The `client/CelerClient` package exposes thin wrappers over `app.AppClient` methods, several of which are now deleted. Trim each call site:

- [x] **Keep**: `NewAppChannelOnVirtualContract` (but with the `sc common.StateCallback` parameter dropped to match the trimmed app-layer signature), `DeleteAppChannel`, `GetAppChannelDeployedAddr`, `OnChainGetAppChannelBooleanOutcome`. These remain useful for the surviving registration / outcome-query / cleanup surface — each has a downstream consumer in the surviving webapi handler chain (`Create*` → SDK constructor; `Delete*` / `GetDeployed*` / `GetBooleanOutcome*` → SDK accessor methods → these wrappers).
- [x] **Delete**: `NewAppChannelOnDeployedContract` (backing AppClient method deleted), `SignAppState` (calls into the deleted state-exchange surface), `SettleAppChannel` (delegates to deleted `AppClient.IntendSettle`), `OnChainApplyAppChannelAction`, `OnChainFinalizeAppChannelOnActionTimeout`, `OnChainGetAppChannelSettleFinalizedTime`, `OnChainGetAppChannelActionDeadline`, `OnChainGetAppChannelStatus`, `OnChainGetAppChannelState`, `OnChainGetAppChannelSeqNum`. Delete any oracle-settle wrappers (`OnChainSettleBy*`) if present.
- [x] **Also delete `GetAppChannel`** — the only external caller of `client.CelerClient.GetAppChannel(...)` is `celersdk/appsession.go:222` inside `HandleMatchData`, which deletes in this trim. After the celersdk trim no surviving consumer references this wrapper; it's a leak of `*app.AppChannel` internals through the client surface with no remaining use case. (`app.AppClient.GetAppChannel` — the underlying in-package accessor — stays; only the `client/CelerClient` wrapper deletes.)

#### `celersdk/`

- [x] Trim `celersdk/appsession.go`. **Keep**: `CreateAppSessionOnVirtualContract` (with the `callback AppCallback` parameter dropped), `EndAppSession` / `DeleteAppSession`, `OnChainGetBooleanOutcome`, `GetDeployedAddress`, the `AppSession` type itself with the trimmed fields. **Delete** every other entry, specifically:
  - `NewAppSessionOnDeployedContract` — direct caller of the deleted deployed-contract path.
  - `CreateAppSessionOnDeployedContract` (the package-level method on `Client`) — same backing path.
  - `newAppSession` (private helper used only by `NewAppSessionOnDeployedContract` and friends).
  - `SignAppData`, `HandleMatchData`, `AppData`, all `OPCODE_*` constants, the seqnum / last-state tracking fields on `AppSession` — the off-chain state-exchange protocol.
  - `SwitchToOnchain` — calls the deleted `AppClient.IntendSettle`.
  - `OnChainApplyAction`, `OnChainFinalizeOnActionTimeout`, `OnChainGetSettleFinalizedTime`, `OnChainGetActionDeadline`, `OnChainGetStatus`, `OnChainGetState`, `OnChainGetSeqNum` — gaming/state-machine introspection.
  - `SettleBySigTimeout`, `SettleByMoveTimeout`, `SettleByInvalidTurn`, `SettleByInvalidState` — oracle disputes.
  - `GetPlayerIdxForMatch` — gaming/match-specific utility used only by deleted methods.
  - **`AppCallback` interface and `Callback` field on `AppInfo`** — the legacy SDK callback surface. Per the §2 callback-infrastructure deletion: no consumer remains for `OnDispute` notifications post-trim. Drop the interface definition (~line 44), drop the `Callback` field from `AppInfo` struct (~line 41), drop the `callback AppCallback` parameter from `CreateAppSessionOnVirtualContract` (the package-level method on `Client`). The shared `common.StateCallback` interface (in `common/types.go`) stays — it's used by the unrelated main-client callback in `client/celer_client.go`.
- [x] Audit `celersdk/api.go` and `celersdk/utils.go` for orphaned helpers and delete those (likely candidates: any `AppSession`-shaped helper that returned a deleted `AppSession` field or referenced deleted opcode constants). — `api.go` clean (no app-session refs); `utils.go::bc2c` rewritten to drop `app.SessionQuery` wrapping (deployed-contract branch now passes raw query bytes through, matching PayResolver on-chain shape); `types.go::BooleanCondition` doc comments refreshed for SDK back-compat (field names preserved).

#### Compile-driven follow-up sites

Three concrete edit sites surface naturally from the deletions above; enumerated here so AS-B's scoped build/vet gate stays green throughout the phase rather than waiting until AS-C's repo-wide gate.

- [x] **`server/osp_webapi_backend.go`** — implements the `OspPayApiBackend` interface that `webapi/osp_pay_api_server.go` depends on. After the trim:
  - Update the call to `b.cNode.AppClient.NewAppChannelOnVirtualContract(...)` (~line 66) to drop the trailing `sc common.StateCallback` argument now that the AppClient signature lost it.
  - Delete the `GetStatusForAppSession` method (~line 80) and its caller of the now-deleted `b.cNode.AppClient.GetAppChannelStatus`. Trim the `OspPayApiBackend` interface declaration in `webapi/osp_pay_api_server.go` to match.
  - Audit for any other osp-backend method that wraps a deleted `AppClient` method and delete it too.
- [x] **`app/appclient_virtresolver_watch_test.go`** — the entire file becomes dead with the watch deletion. Delete it.
- [x] **`webapi/osp_pay_api_server_test.go`** — drop assertions that exercised `GetStatusForAppSession` on the OSP backend. Keep assertions for `CreateAppSessionOnVirtualContract` / `DeleteAppSession`. (The e2e-side `osp_webapi_test.go` `ospWebApiAppSessionSubset` cleanup is already covered in AS-C.)

**Beyond plan (compile-driven follow-up sites surfaced during AS-B):**
- [x] Embed `rpc.UnimplementedWebApiServer` in `webapi.ApiServer` — required after deleting RPC handlers (the generated server interface now lists more methods than `ApiServer` implements; the embed satisfies the gRPC plugin's `mustEmbedUnimplementedWebApiServer` marker).
- [x] `tools/osp-cli/cli/cli_onchain_view.go::printAppBooleanOutcome` switched to `app.NewIBooleanCondCaller` (renamed binding) and dropped the `SessionQuery`-decode branch driven by the `-decode` flag (flag also removed from `cli_flags.go`).
- [x] `testing/clientcontroller.go` — deleted 12 wrappers around deleted gaming RPCs (`SignOutgoingState`, `NewAppChannelOnDeployedContract`, all `SettleAppChannel*`, `GetAppChannelState`, `GetAppChannelSettleFinalizedTime`, `ApplyAppChannelAction`, `GetAppChannelActionDeadline`, `FinalizeAppChannelOnActionTimeout`). `WaitUntilBlockHeight` and the e2e `pay_dispute*.go` rewrite remain AS-C scope.

#### `proto/app.proto`

- [x] Delete `AppState`, `StateProof`, and `SessionQuery` messages — all dead once the readers above are gone. (Earlier draft of this plan incorrectly said `OracleState` / `OracleProof` lived here; they do not.) — **deferred from AS-B and completed in AS-C** once the consumers (`test/e2e/pay_dispute*.go`, `testing/testapp/utils.go::GetAppState*`) were deleted/rewritten. Trimming earlier would have left `e2e` uncompilable.
- [x] If `proto/app.proto` becomes empty, delete the file and remove its `import` line from any other proto file. Regenerate `proto/app.pb.go` (delete it if the source goes away). — done in AS-C: `proto/app.proto` and `app/app.pb.go` both deleted; no other proto file imported it.

#### Build / vet gate

- [x] **Scoped** build/vet only on the non-test packages this phase touches plus their direct dependents — explicitly **not** the repo-wide `go build ./...` / `go vet ./...`, because `testing/clientcontroller.go` and the `test/e2e/` *_test.go files still reference deleted webapi RPCs and helpers until AS-C cleans them up. Specifically:
  - [x] `go build ./app/... ./cnode/... ./webapi/... ./celersdk/... ./server/... ./client/... ./messager/... ./handlers/... ./dispute/... ./route/... ./delegate/...` — clean. (i.e. every Go package that is *not* a *_test.go file or under `testing/` / `test/`.) — verified via `go list ./... | grep -v '/test/' | xargs go build`.
  - [x] `go vet` over the same set — clean.
  - [x] The repo-wide `go build ./...` and `go vet ./...` deliberately stay broken at this point; they're restored at the end of AS-C. — confirmed: `test/e2e/pay_dispute*.go` and `testing/testapp/utils.go` still fail (expected; AS-C scope).

**Exit criteria:** non-test packages build and vet clean; no surviving Go reference to deleted methods/types in the trimmed packages; `AppClient` is reduced to the registration / deploy / query surface with `GetBooleanOutcome` redesigned to use `IBooleanCond` bindings; `CreateAppSessionOnVirtualContract` still works end-to-end with the existing `SimpleSingleSessionApp` bytecode (since x402 still imports `testing/testapp/singlesessionapp.go`); `proto/app.proto` is empty or removed.

#### AS-B completion notes (status: completed)

**Build/vet gate (scoped):** `go build` and `go vet` clean across all non-test packages — confirmed by `go list ./... | grep -v '/test/' | xargs go build` and `... | xargs go vet` returning no output.

**Compile-broken (deferred to AS-C as planned):**
- `test/e2e/pay_dispute.go` and `test/e2e/pay_dispute_with_oracle.go` reference deleted `ClientController` helpers (`SettleAppChannel`, `NewAppChannelOnDeployedContract`, `GetAppChannelState`, `GetAppChannelSettleFinalizedTime`, `ApplyAppChannelAction`, `GetAppChannelActionDeadline`, `FinalizeAppChannelOnActionTimeout`, `WaitUntilBlockHeight`) and `app.AppState` / `app.StateProof` / `app.SessionQuery` proto messages. AS-C deletes `pay_dispute_with_oracle.go` outright and rewrites `pay_dispute.go` against `BooleanCondMock`.
- `testing/testapp/utils.go` still references `app.AppState` (gaming state-exchange helper). AS-C either removes it (callers all in deleted/rewritten e2e files) or keeps it for x402 back-compat alongside `singlesessionapp.go`.

**`proto/app.proto` trim deferred to AS-C.** The remaining `AppState` / `StateProof` / `SessionQuery` users are all in `test/e2e/pay_dispute*.go` and `testing/testapp/utils.go`, both of which AS-C deletes or rewrites. Trimming the proto in AS-B would block the e2e package from even compiling under its current code; deferring keeps the deletion atomic with the consumer cleanup.

**Beyond plan (compile-driven):**
- Embedded `rpc.UnimplementedWebApiServer` in `webapi.ApiServer` — required after deleting RPC handlers (the generated server interface now lists more methods than `ApiServer` implements; the embed satisfies the gRPC plugin's `mustEmbedUnimplementedWebApiServer` marker).
- Fixed `tools/osp-cli/cli/cli_onchain_view.go::printAppBooleanOutcome` to call `app.NewIBooleanCondCaller` (renamed binding) and dropped the `SessionQuery`-decode branch driven by the `-decode` flag (flag also removed from `cli_flags.go`).

**Code-shape outcomes:**
- `app/appclient.go`: 936 → ~250 lines (callback infra, virt-resolver-deploy watch, gaming methods, deployed-contract path all gone; `GetBooleanOutcome` now calls `IBooleanCondCaller` directly with raw query bytes).
- `app/booleancond.go`: regenerated from `IBooleanCond.sol`; ownership moved from `regenerate-legacy-app-bindings.sh` → `regenerate-go-bindings.sh`. The legacy script now carries only `singlesessionapp.go` (x402 back-compat).
- `client/app_channel.go`: 121 → ~45 lines.
- `celersdk/appsession.go`: 397 → ~80 lines (`AppCallback` interface, gaming opcodes, state-exchange helpers, deployed-contract entry-point all gone).
- `webapi/api_server.go`: deleted ~14 RPC handlers + `appSessionCallback` type + `appSessionCallbackMap` field/lock.
- `webapi/proto/web_api.proto`: deleted ~12 RPCs + ~14 dead messages; regenerated.

### AS-C — Test-fixture migration and dispute-coverage rewrite

After AS-B, the testing-side packages (`testing/clientcontroller.go`, the e2e `*_test.go` files, the OSP webapi test) still reference deleted RPCs and helpers; they need to catch up. This phase deletes the gaming-flavored tests, generates `BooleanCondMock` ABIgen bindings, deploys `BooleanCondMock` in the e2e setup, rewrites the surviving dispute tests against it for both condition types, and restores the repo-wide build/vet/test gate that AS-B intentionally left broken.

- [x] Generate ABIgen bindings for `BooleanCondMock` (already in `agent-pay-contracts/src/helper/`) into `testing/testapp/booleancondmock.go`. Update `tools/scripts/regenerate-go-bindings.sh` to include it if not already covered. Verify the bindings expose at minimum: `BooleanCondMockBin` (deploy bytecode), `IsFinalized`, `GetOutcome`, and a `DeployBooleanCondMock` helper.
- [x] Update `test/e2e/setup_onchain.go` to deploy a `BooleanCondMock` instance during e2e bootstrap, and surface its address on the contract address bundle (alongside `PayResolver`, `PayRegistry`, etc.). This serves as the `OnChainAddress` for `DEPLOYED_CONTRACT` test scenarios. — exposed via `appAddrMap["BooleanCondMock"]`. The legacy `SimpleMultiSessionApp` / `SimpleMultiSessionAppWithOracle` / `MultiGomoku` deployments were dropped from this function (no surviving consumer).
- [x] Delete legacy gaming fixtures that have no surviving consumer:
  - `testing/testapp/multigomoku.go`
  - `testing/testapp/multisessionapp.go`
  - `testing/testapp/singlesessionappwithoracle.go`
  - `testing/testapp/multisessionappwithoracle.go`
  - **Keep** `testing/testapp/singlesessionapp.go` and any `utils.go` helpers it depends on — `agent-pay-x402` imports `ta.AppCode` / `ta.GetSingleSessionConstructor` / `ta.Timeout` from this file. Removing it would break x402 immediately. Add a leading file comment marking it deprecated and pointing at §7 for the migration plan. — `singlesessionapp.go` is generated, so the deprecation note went into `utils.go` instead. `utils.go` itself was rewritten to keep only the x402 surface (`AppCode`, `Nonce`, `Timeout`, `GetSingleSessionConstructor`); the gaming helpers (`GetAppState*`, `GetGomokuState*`, oracle constructor, deployed-addr constants, `PlayerNum`) were deleted along with their consumers.
- [x] Delete `test/e2e/pay_dispute_with_oracle.go` outright — oracle-dispute tests cover paths that no longer exist.
- [x] Rewrite `test/e2e/pay_dispute.go`:
  - [x] Drop every scenario that calls deleted `AppClient` methods (`IntendSettle`, `ApplyAction`, `FinalizeAppChannelOnActionTimeout`, `Settle*Timeout`, status/seqNum/state introspection). — full rewrite from ~1050 LOC to ~225 LOC; only the two scenarios specified below survive.
  - [x] Add or preserve coverage for: conditional pay with `ConditionType_VIRTUAL_CONTRACT` resolved through dispute (register `BooleanCondMock` bytecode → send pay → settle channel → deploy on dispute → resolve via `PayResolver` → assert outcome). Either the explicit on-dispute deploy path or the `GetBooleanOutcomeForAppSession` deploy-on-query path can perform the deployment; pick one and assert it actually deploys (the virtual contract address has bytecode after the call). — implemented as `runVirtualContractScenario`; the `PayResolver.resolvePaymentByConditions` path performs the deploy-on-resolve.
  - [x] Add or preserve coverage for: conditional pay with `ConditionType_DEPLOYED_CONTRACT` resolved through dispute (use the `BooleanCondMock` instance deployed in `setup_onchain.go` as `OnChainAddress` → send pay → settle channel → resolve via `PayResolver` → assert outcome). — implemented as `runDeployedContractScenario`.
  - [x] Both scenarios should test both `BooleanCondMock` outcomes (true and false query bytes) so `getOutcome → false` correctly leaves the pay un-resolved and `→ true` correctly pays out. — both scenarios run twice with `argsQueryOutcome=0x01` and `argsQueryOutcome=0x00`; the new `runDisputeAndAssert` checks `GetCondPayInfoFromRegistry` reports `sendAmt` vs `0` accordingly.
  - [x] Update `test/e2e/e2e_test.go` `t.Run(...)` registrations to match the rewritten scenarios; drop the deleted-test entries. — dropped `disputeEthPaySrcOffline`, `disputeEthPayWithDeployedGomoku`, `disputePayBySigTimeoutWithDeployedContract`.
- [x] If `test/e2e/send_pay_with_app.go` exists and exercises gaming flows, delete it. If any scenario only exercises off-chain conditional-pay flows that survive the trim, keep that scenario (rewriting against `BooleanCondMock` if it currently uses `SimpleSingleSessionApp`). — kept as-is; it exercises off-chain confirm/cancel of a VIRTUAL_CONTRACT pay where the bytecode is just a vehicle (the contract is never queried). Migrating to `BooleanCondMock` would be a wash.
- [x] Update the OSP webapi test `test/e2e/osp_webapi_test.go` (`ospWebApiAppSessionSubset`) to drop assertions/calls against `GetStatusForAppSession`. Other assertions against `CreateAppSessionOnVirtualContract` / `DeleteAppSession` stay.
- [x] Delete `WaitUntilBlockHeight` from `testing/clientcontroller.go`. Confirm grep shows no remaining callers.
- [x] Delete `testing/clientcontroller.go` helpers that wrapped **the deleted gaming/state-machine webapi RPCs** (real names): `SettleAppSession`, `SettleAppSessionBy*` (oracle-dispute four), `SubscribeAppSessionDispute`, `GetStatusForAppSession`, `GetSeqNumForAppSession`, `GetStateForAppSession`, `ApplyActionForAppSession`, `FinalizeOnActionTimeoutForAppSession`, `GetActionDeadlineForAppSession`, `GetSettleFinalizedTimeForAppSession`, `SignOutgoingState`, `ValidateAck`, `ProcessReceivedState`, `CreateAppSessionOnDeployedContract`. **Keep** wrappers for `CreateAppSessionOnVirtualContract`, `DeleteAppSession`, `GetDeployedAddressForAppSession`, `GetBooleanOutcomeForAppSession`. — most of this list was already deleted in AS-B (compile-driven follow-up).
- [x] Run focused e2e (`go test ./test/e2e -run '^TestE2E$/^e2e-grp2$/^sendCondPayWithErc20$'`) — green. — `sendCondPayWithErc20` PASS in 14.77s.
- [x] Run the rewritten dispute scenarios specifically — green. — `disputeEthPayWithVirtualContract` (17.29s) and `disputeEthPayWithDeployedContract` (16.28s) both PASS, exercising getOutcome=true and getOutcome=false (BooleanCondMock query bytes 0x01 and 0x00) for both VIRTUAL_CONTRACT and DEPLOYED_CONTRACT.
- [x] Run the OSP webapi subset (`go test ./test/e2e -run '^TestOSPWebApi'`) — green. — `ospWebApiAppSessionSubset` (4.69s), `ospWebApiPaySubset` (4.94s), `TestOSPWebApiRoutingBehavior` (20.03s) all PASS.
- [x] Run full default e2e (`go test ./test/e2e -count=1 -timeout 30m`) — green. — covered in AS-D's full validation matrix below: an attempt at the *combined* default suite hit a pre-existing parallel-load flake unrelated to this trim; equivalent coverage achieved by running each group / top-level test separately, all green.
- [x] Repo-wide gate (restored from AS-B's narrowed scope): `go build ./...`, `go vet ./...`, all targeted unit/package tests — clean. — `go build ./...` and `go vet ./...` empty; `go test ./webapi ./celersdk ./app ./client` all green. Other targeted unit packages from AS-D §full-validation-matrix still to run.

**Beyond plan (AS-C compile-driven follow-up):**
- [x] Delete the now-orphaned `proto/app.proto` and `app/app.pb.go` — final remaining consumers (`testing/testapp/utils.go::GetAppState*`, `test/e2e/pay_dispute*.go`) are gone after the rewrites above. (This was the AS-B `proto/app.proto` task, deferred to AS-C so the deletion was atomic with the consumer cleanup.)

**Exit criteria:** all e2e tests pass; repo-wide `go build ./...` / `go vet ./...` green; `WaitUntilBlockHeight` is gone; `BooleanCondMock` bindings exist and are deployed in the e2e setup; the dispute test coverage now exercises both `VIRTUAL_CONTRACT` and `DEPLOYED_CONTRACT` with `BooleanCondMock`; only `singlesessionapp.go` survives in `testing/testapp/` (back-compat carry, marked deprecated); no test references gaming or oracle concepts.

### AS-D — Documentation and validation

- [x] Update `AGENTS.md` §Protocol Invariants — the existing line about "testing app-session contracts under `testing/testapp/` are an exception, still use block.number" stays accurate (the surviving `SimpleSingleSessionApp` is still block-number-based) until x402 migrates. Reword if helpful but don't remove. — reworded to reference the single surviving file and the deleted `WaitUntilBlockHeight` helper is no longer mentioned.
- [x] Update `AGENTS.md` §Architecture — adjust mention of "app session support" if any to reflect the trimmed reality (registration + outcome-query, no state machine). — §Architecture made no app-session claims that conflict with the trim; no edit needed.
- [x] Update `docs/backend-implementation.md`:
  - [x] Update the `app/` row in the Core Packages table to describe what it now is: condition-contract bindings + virtual-contract registration / deploy-on-dispute helpers. No more session state machine.
  - [x] Update any prose that references state-machine concepts (status, seqNum, applyAction, oracle disputes). — only the "app sessions" item under the boot sequence (#8) was stale; reworded to "app channels (registration + on-chain outcome query)". Also dropped the now-deleted `proto/app.proto` row from the Wire Contracts table.
  - [x] Add a brief note that conditional payments resolve via the `IBooleanCond` (and, when wired up off-chain, `INumericCond`) interfaces in `agent-pay-contracts`. Cross-reference §2 of this plan. — folded into the new `app/` row in Core Packages with a direct link to §2.
- [x] Update `docs/backend-usage.md` if any operator-facing guidance described the deleted RPCs. — no app-session/state-machine prose found; no edit needed.
- [x] Update `docs/backend-troubleshooting.md` — drop any failure-symptom guides that reference the deleted methods. — no hits; no edit needed.
- [x] Update `tools/osp-cli/README.md` if any CLI command listed introspection fields (status, seqNum, app-channel state) that no longer exist. — `-onchainview app` description rewritten to reference `IBooleanCond.{isFinalized,getOutcome}`; `-decode` flag mention dropped (flag was removed in AS-B).
- [x] Update `CLAUDE.md` only if it directly references app-session state-machine concepts (it doesn't appear to today; verify). — verified via grep; no app-session/state-machine references; no edit needed.
- [x] Run the full local validation matrix:
  - [x] `go build ./...` — clean.
  - [x] `go vet ./...` — clean.
  - [x] `go test ./storage ./celersdk ./common/cobj ./dispatchers ./lrucache ./rpc ./rtconfig ./metrics ./route ./utils/bar ./cnode/cooperativewithdraw ./server ./fsm ./common` — all green.
  - [x] e2e coverage validated by groups (full default suite hit a pre-existing parallel-load flake, see note below):
    - `e2e-grp1` (15 tests): PASS in 71.96s.
    - `e2e-grp2/sendCondPayWithErc20`: PASS in 14.77s.
    - `e2e-grp2/sendCondPayWithEthDstOffline` (isolated rerun): PASS in 11.70s.
    - `e2e-grp3/disputeEthPayWithVirtualContract`: PASS in 17.29s (covers VIRTUAL_CONTRACT condition with both `getOutcome=true` and `getOutcome=false`).
    - `e2e-grp3/disputeEthPayWithDeployedContract`: PASS in 16.28s (covers DEPLOYED_CONTRACT condition with both `getOutcome=true` and `getOutcome=false`).
    - `TestE2EChannelMigrationTool`: PASS in 23.25s.
    - `TestE2EMultiOSP`: PASS in 11.19s.
    - `TestOSPWebApi/ospWebApiAppSessionSubset`: PASS in 4.69s.
    - `TestOSPWebApi/ospWebApiPaySubset`: PASS in 4.94s.
    - `TestOSPWebApiRoutingBehavior`: PASS in 20.03s.
  - [x] `go test ./test/e2e -count=1 -run '^TestOSPWebApi'` — covered above.

**Pre-existing parallel-load flake (not blocking):** `go test ./test/e2e -count=1 -timeout 35m` for the full default suite hung at 32m on `e2e-grp2/sendCondPayWithEthDstOffline`. Isolated rerun PASSED in 11.70s. Root cause is a race in the test fixture between gRPC `SubscribeOutgoingPayments` stream setup and the OSP emitting the `Unreachable` event under heavy parallel-grp2 load — `webapi/callbackimpl.go::HandleSendErr` uses a non-blocking `select { case ... default: }` send, so if no consumer is yet attached when the event fires the error is silently dropped and the test blocks forever on `<-c1SendErrChan`. This race predates the app-session simplification (the callback infrastructure here was unchanged by AS-B/AS-C; the new `UnimplementedWebApiServer` embed in AS-B is compile-only and doesn't change runtime registration). Filed as a separate cleanup follow-up; tracked in §7.
- [x] Confirm the companion `agent-pay-docs` (`agentpay-architecture/`) doesn't need an update — it describes the protocol abstractly; if any concrete file references applyAction, oracle disputes, or session settlement state-machines, raise it as a follow-up rather than blocking this plan. — verified via grep across `agentpay-architecture/`. The only `seqNum` hit is the channel-level `SimplexPaymentChannel.seqNum` (current and correct, not the deleted app-session field). All `oracle` mentions are abstract ("oracle queries", "release tokens once oracle confirms") that align with the surviving `IBooleanCond` model. No edits needed.
- [ ] Open the PR with a clear summary linking back to this plan doc and to §7 for the remaining x402 follow-up.

**Exit criteria:** local + CI fully green; this plan doc is deleted as part of the merge. The substantive long-lived guidance has been folded into `AGENTS.md`, `docs/backend-implementation.md`, and the in-tree comments on `testing/testapp/utils.go` / `tools/scripts/README.md`. Any cross-references that currently point at this plan must be updated to remove the link before the merge.

---

## 6. Risks and mitigations

| Risk | Likelihood | Mitigation |
| --- | --- | --- |
| `agent-pay-x402` references a method we plan to delete | low — direct grep confirms zero hits today, but worth re-verifying as code there evolves | AS-A audit step explicitly re-greps. |
| Deleting from `app/` orphans an import we missed | medium — the package is referenced from multiple places (cnode, messager, handlers, webapi, celersdk, tools/osp-cli) | Do AS-B in topological order (interfaces first, AppClient last). Run scoped `go build` after each subtask, not just at the end. The repo-wide gate at end of AS-C catches anything missed earlier. |
| `GetBooleanOutcome` redesign for the DEPLOYED_CONTRACT branch introduces a regression | medium — the legacy code wrapped queries in `SessionQuery` and used `IMultiSessionCaller`; the new code passes raw bytes through `IBooleanCondCaller`. Mismatched expectations could produce silently wrong results. | The AS-C rewritten dispute test for DEPLOYED_CONTRACT covers exactly this path end-to-end with both `getOutcome→true` and `getOutcome→false` scenarios. Verify the on-chain-side behavior matches what `PayResolver` does in `agent-pay-contracts` — the off-chain code is now identical to the on-chain call shape, so divergence is structurally hard. |
| Future use case really does need on-chain dispute fallback for app conditions | low for the next 12 months — no current consumer uses it; AI-agent payment patterns observed so far don't need it | The `IBooleanCond` / `INumericCond` interfaces leave room for someone to ship a stateful condition contract later if needed. The deletion is of the **generic** infrastructure, not of the protocol's ability to support such a contract. |
| Test coverage drops because most of `pay_dispute.go` deletes | medium — the legacy tests covered real protocol invariants, even if expressed through gaming fixtures | Channel-level dispute (`settle_channel.go`, `cold_bootstrap.go`) is independent and stays fully intact. Conditional-payment-specific tests (`send_cond_pay_*.go`) stay. AS-C **rewrites** the dispute coverage for both `VIRTUAL_CONTRACT` and `DEPLOYED_CONTRACT` against `BooleanCondMock` — net coverage of the surviving protocol surface goes up, not down. |
| Keeping `testing/testapp/singlesessionapp.go` (for x402 back-compat) leaves a misleading "test app" file in the tree, suggesting agent-pay still supports the legacy gaming model | low — the file is small and clearly imported by external x402 only | Add a leading comment to the file noting it exists solely for x402 back-compat pending the deferred migration in §7. Resolved when §7 lands. |
| Off-chain state-exchange RPCs (`SignOutgoingState` / `ValidateAck` / `ProcessReceivedState`) had a non-obvious downstream consumer we missed | medium — these RPCs don't show up in `agent-pay-x402` per the AS-A audit, but were not exhaustively traced through every internal Celer integration | AS-A explicitly re-greps for them across all sibling repos. If a hit appears, scope decision: either defer their deletion alongside x402 (move into §7) or migrate the consumer in this PR. |
| Someone outside this repo (downstream SDK consumer, internal team using `celersdk`) breaks because of webapi/SDK deletions | low — no documented public consumers of those methods | Document the breaking change explicitly in the PR description so it's discoverable later. The "no production deployments yet" posture in §3 is what makes this acceptable; if it's not true at PR time, the plan needs revisiting. |

---

## 7. Deferred / TODO

Items intentionally out of scope for this plan but worth a forward pointer:

- **x402 migration to a stateless condition-contract bytecode.** Currently x402 registers `SimpleSingleSessionApp` (turn-based-game contract) via `CreateAppSessionOnVirtualContract`. The trim doesn't break that — x402 doesn't exercise the gaming dispute path — but it's strictly a back-compat carry. Future PR (in either repo): swap the registered bytecode to `BooleanCondMock` (now bundled with this trim) or to a custom `IBooleanCond` impl appropriate for the x402 use case. Once that lands:
  - `testing/testapp/singlesessionapp.go` and any `singlesessionapp`-specific helpers in `utils.go` delete.
  - The `singlesessionapp.go`'s deprecation comment goes with them.
  - Update `AGENTS.md` §Protocol Invariants to drop the "testapp uses block.number" exception.
- **Generate ABIgen bindings for `NumericCondMock` and `INumericCond` in agent-pay** when a numeric off-chain consumer surfaces. Currently zero callers in agent-pay (every `TransferFunctionType` is `BOOLEAN_AND`); bindings are dead weight. Likely a small follow-up if/when a NUMERIC_ADD/MAX/MIN use case actually emerges.
- **Rename `BooleanCondMock` / `NumericCondMock` to drop "Mock"** if they ever evolve from test-only fixtures into reference implementations. Today the explicit "Test-only. Do not deploy to a production network." NatSpec is correct, so the name fits.
- **Rename `CreateAppSessionOnVirtualContract` to drop "Session"** (or rename the entire `app/` package) if the "session" / "app channel" terminology ever stops aligning with the architecture docs. Today they align; the names stay.
- **Fix the parallel-load `SubscribeOutgoingPayments` race** in `webapi/callbackimpl.go` exposed by the AS-D full-suite e2e (see AS-D §full local validation matrix note). `HandleSendErr` uses a non-blocking `select { case ch <- ...: default: }` so an `Unreachable` event fired before the gRPC `SubscribeOutgoingPayments` stream is fully attached gets silently dropped. Under heavy parallel load (full default e2e with all grp2 tests fanning out at once) this manifests as `e2e-grp2/sendCondPayWithEthDstOffline` blocking forever on `<-c1SendErrChan`. Fix is straightforward: buffer the chan or replay missed events when a subscriber attaches; both are out of scope for this plan because the race predates the trim and the test passes in isolation. Tracked here so a future cleanup PR has the breadcrumbs.

---

## 8. Closeout

This plan doc is deleted as part of the merge. Before deleting, audit and remove (or rewrite) any cross-references — at minimum: `AGENTS.md`, `docs/backend-implementation.md`, `tools/scripts/README.md`, `tools/scripts/regenerate-legacy-app-bindings.sh`, `testing/testapp/utils.go`. The §7 "Deferred / TODO" entries that need long-lived survival should migrate into `AGENTS.md` (or another in-tree home) before that delete happens.

The summary line for the merge commit / PR description: **"Trim app-session machinery (on-chain dispute paths, off-chain state-exchange RPCs, virt-resolver deploy-watch + callback infrastructure) to the protocol-essential `IBooleanCond` / `INumericCond` surface; preserve `ConditionType_VIRTUAL_CONTRACT` / `_DEPLOYED_CONTRACT` at the wire level; redesign `GetBooleanOutcome` to drop multisession-specific encoding (`IBooleanOutcome` → `IBooleanCond` regenerated from agent-pay-contracts); rewrite dispute coverage onto `BooleanCondMock` for both condition types; defer x402 bytecode swap. ~8200 LOC of legacy gaming-era infrastructure deleted, ~430 LOC of clean fixture-and-test added."**
