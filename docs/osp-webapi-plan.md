# OSP WebAPI Rollout Plan

This document is the concrete rollout plan for exposing an OSP WebAPI gRPC listener from `bin/server` without creating a second runtime or breaking existing OSP callback behavior.

The plan is intentionally phased.

- [x] Phase 1 is intended to be production-ready for the pay-centric seller-OSP use cases needed today.
- [x] Later phases are only needed if this repo wants full OSP WebAPI parity for channel-scoped operations, public-interface hardening, or richer subscription semantics.
- [x] In other words: phase 1 is not a prototype. It is the first production slice.

Audience: agent-pay maintainers implementing the change in this repo, and agent-pay-x402 maintainers reviewing the contract needed for seller-OSP topologies.

Status legend: `[x]` means reviewed and locked by current design review. `[ ]` means pending implementation or validation in this repo.

## Rollout Model

- [x] The current `WebApi` surface mixes two kinds of operations: pay-centric operations and channel-scoped operations.
- [x] Pay-centric operations are already well-defined on an OSP because they are keyed by payment ID or explicit destination.
- [x] Channel-scoped operations are not yet well-defined on an OSP because the existing API shape does not name which peer/channel the OSP should act on.
- [x] Phase 1 therefore means: production-ready pay-centric OSP WebAPI support, with all ambiguous channel-scoped methods returning `codes.Unimplemented` instead of guessing.
- [x] A later parity phase is required only if we want those currently ambiguous channel-scoped methods to be supported on OSP.

## Phase 1 Acceptance Criteria

- [ ] `bin/server` accepts a new optional `-webapigrpc <host:port>` flag.
- [ ] When `-webapigrpc` is unset, OSP startup and behavior are unchanged.
- [ ] When `-webapigrpc` is set, OSP starts a separate gRPC listener without TLS transport credentials (plaintext gRPC) dedicated to `rpc.WebApiServer`.
- [ ] The OSP WebAPI listener uses the existing OSP runtime and does not initialize a second SDK client or second cnode lifecycle.
- [ ] Existing OSP fee/delegate send and receive side effects still fire after the listener is enabled.
- [ ] The phase-1 pay-centric RPC subset works through a normal `rpc.WebApiClient`.
- [ ] All deferred channel-scoped WebAPI RPCs return `codes.Unimplemented`.
- [ ] Focused unit and e2e validation cover the callback mux, the new listener, and one unsupported RPC, with the focused listener path covered by the proposed `test/e2e/osp_webapi.go` test.

## Locked Scope Decisions

- [x] Phase 1 is pay-centric only.
- [x] Phase 1 exposes `rpc.WebApiServer` only.
- [x] Phase 1 does not expose `rpc.InternalWebApiServer`.
- [x] Phase 1 does not add an HTTP or grpc-web listener.
- [x] Phase 1 does not change protobuf definitions.
- [x] Phase 1 does not claim full WebAPI parity on OSP. Channel-scoped methods stay `Unimplemented` until they have an explicit peer/channel selection model.
- [x] Phase 1 uses `-webapigrpc <host:port>` rather than a bare integer port.
- [x] Phase 1 uses gRPC without TLS transport credentials on this optional listener.
- [x] That plaintext gRPC choice is acceptable only because the intended phase-1 deployment is a colocated localhost/private caller, typically bound to `127.0.0.1:<port>` and used by a same-host client process.
- [x] Phase 1 uses single-subscriber semantics for payment subscriptions.
- [x] Phase 1 returns `codes.Unimplemented` for all deferred channel-scoped methods.

## Existing Constraints This Plan Must Respect

- [x] [webapi/api_server.go](../webapi/api_server.go) `NewApiServer(...)` initializes a fresh SDK client and waits on `callbackImpl.clientReady`, so it cannot be started inside `bin/server` without creating a second runtime.
- [x] [dispatchers/celer_msg_dispatcher.go](../dispatchers/celer_msg_dispatcher.go) stores one receive callback and one send callback, so OSP callback ownership must be preserved explicitly.
- [x] [server/server.go](../server/server.go) already registers the OSP itself as the cnode send/receive callback sink.
- [x] Channel-scoped SDK/WebAPI methods resolve token to channel through a single implicit peer, which is ambiguous on an OSP.

## Why Channel-Scoped WebAPI Calls Are Not Yet Safe On OSP

- [x] The existing client-side channel lookup resolves `token -> cid` through a single implicit peer in [client/api.go](../client/api.go#L208).
- [x] That implicit peer is `c.svrEth`, which comes from the profile fields in [common/profile.go](../common/profile.go#L65) and [common/profile.go](../common/profile.go#L66).
- [x] That model is correct for a normal client node because the client has one upstream OSP configured in its profile.
- [x] That model is not correct for an OSP because an OSP can have many peers and many channels for the same token.
- [x] As a result, existing channel-scoped methods like [Deposit](../webapi/api_server.go#L231), [CooperativeWithdraw](../webapi/api_server.go#L318), and [GetBalance](../webapi/api_server.go#L370) do not carry enough information to choose the right OSP-side channel safely.
- [x] Returning `codes.Unimplemented` is therefore the production-safe behavior in phase 1. It avoids silently acting on the wrong peer/channel.
- [x] Supporting those methods on OSP requires a later API-design phase that adds an explicit peer/channel selector or an equivalent unambiguous contract.

## External Contract After Phase 1 Lands

- [x] agent-pay-x402 will dial the OSP listener with a normal `rpc.WebApiClient`.
- [x] The intended caller in phase 1 is a colocated local or private same-host process, not an arbitrary remote client on a public network.
- [x] No x402-side protobuf regeneration is required for phase 1.
- [x] `CooperativeWithdraw` remains on the Admin RPC path for seller-OSP mode.
- [x] `GetBalance` is intentionally `Unimplemented` on OSP WebAPI in phase 1.
- [x] OSP balance observation remains available through the Admin surface, so `GetBalance` being `Unimplemented` is operationally acceptable for phase 1.
- [x] A likely phase-1b candidate is a peer-scoped `GetPeerFreeBalance(peer, tokenInfo)` WebAPI call, because it is unambiguous on an OSP.

## Phase 1 RPC Scope

- [ ] Implement `SendToken` on OSP WebAPI.
- [ ] Implement `SendConditionalPayment` on OSP WebAPI.
- [ ] Implement `GetIncomingPaymentStatus` on OSP WebAPI.
- [ ] Implement `GetIncomingPaymentInfo` on OSP WebAPI.
- [ ] Implement `GetOutgoingPaymentStatus` on OSP WebAPI.
- [ ] Implement `ConfirmOutgoingPayment` on OSP WebAPI.
- [ ] Implement `RejectIncomingPayment` on OSP WebAPI.
- [ ] Implement `SubscribeIncomingPayments` on OSP WebAPI.
- [ ] Implement `SubscribeOutgoingPayments` on OSP WebAPI.
- [x] Do not add `GetOutgoingPaymentInfo` in phase 1, because there is no existing `WebApi` RPC with that name. Adding it would require a proto change and is explicitly out of scope for this cut.
- [ ] Return `codes.Unimplemented` for every other `WebApi` RPC on the OSP listener.

## Proposed Code Changes

### 1. OSP WebAPI server type

- [ ] Add `webapi/osp_pay_api_server.go`.
- [ ] Define a narrow `OspPayBackend` interface for the phase-1 RPC subset.
- [ ] Define `OspPayApiServer` embedding `rpc.UnimplementedWebApiServer`.
- [ ] Keep `OspPayApiServer` additive; do not refactor the existing client-node `ApiServer` into a generic shared abstraction in phase 1.
- [ ] Keep `PaymentInfo` mapping behavior consistent with the existing `paymentInfoFromClientPayment(...)` logic.
- [ ] If direct helper reuse is awkward, move the minimum shared `PaymentInfo` mapping logic into `webapi/payment_convert.go` rather than maintaining two divergent mappings.

### 2. OSP backend wrapper in package main

- [ ] Add `server/osp_webapi_backend.go` in `package main`.
- [ ] Keep the backend in `package main` so it can wrap the already-running OSP state without import-cycle churn around `server.server`.
- [ ] Back the implementation with the existing `*cnode.CNode`, `*storage.DAL`, and local OSP address.
- [ ] Implement `GetIncomingPaymentStatus` and `GetOutgoingPaymentStatus` by reusing the existing payment-state-to-SDK-status mapping logic.
- [ ] Implement `GetIncomingPaymentInfo` as incoming-only: DAL lookup plus `payment.Receiver == myAddr`, otherwise `common.ErrPayNotFound`.
- [ ] Implement `ConfirmOutgoingPayment` and `RejectIncomingPayment` by calling the existing payment confirmation/rejection paths, not by adding new protocol behavior.
- [ ] Implement `SendToken` and `SendConditionalPayment` by building `entity.ConditionalPay` and calling `cnode.AddBooleanPay(...)`.
- [ ] Do not manufacture a synthetic `celersdk.Client` or `client.CelerClient` around the OSP's cnode just to reuse client wrappers.

### 3. Callback fanout

- [ ] Add `server/payment_callbacks_mux.go`.
- [ ] Add `server/payment_callbacks_mux_test.go`.
- [ ] Implement both `event.OnReceivingTokenCallback` and `event.OnSendingTokenCallback` from [common/event/event.go](../common/event/event.go).
- [ ] Fan out to exactly two logical sinks in phase 1: the existing OSP `server` callback handler and the OSP WebAPI payment event feed.
- [ ] Install the callback mux unconditionally, even when `-webapigrpc` is unset, so the cnode callback topology is identical in both startup modes.
- [ ] When `-webapigrpc` is unset, keep the WebAPI-side sink as a no-op feed target rather than removing the mux from the callback chain.
- [ ] Add an implementation comment explaining this tradeoff so a future cleanup does not "optimize away" the unconditional mux installation and accidentally break the fanout invariant.
- [ ] Replace direct registrations in [server/server.go](../server/server.go) with registrations of the mux object.
- [ ] Prove by unit test that `HandleReceivingStart`, `HandleReceivingDone`, `HandleSendComplete`, `HandleDestinationUnreachable`, and `HandleSendFail` each reach both sinks.

### 4. OSP payment event feed

- [ ] Add `webapi/payment_event_feed.go`.
- [ ] Do not reuse `callbackImpl` for OSP subscriptions.
- [ ] Implement single active subscriber per direction in phase 1: one incoming-payments subscriber and one outgoing-payments subscriber.
- [ ] If a second subscriber for the same direction appears, return `codes.FailedPrecondition` rather than replacing the first subscriber silently.
- [ ] Use bounded buffering and best-effort non-blocking publish semantics.
- [ ] Document in code comments and docs that a slow subscriber may drop events.
- [ ] Treat polling RPCs like `GetIncomingPaymentStatus` and `GetIncomingPaymentInfo` as the source of truth in tests; subscriptions are observability, not correctness authority.

### 5. Payment conversion helpers

- [ ] Add `webapi/payment_convert.go`.
- [ ] Move or duplicate only the minimum conversion logic needed for OSP pay backend and event feed.
- [ ] Convert `ConditionalPay + note + payment state/reason` into `celersdkintf.Payment` with the same field formatting used by the current client WebAPI path.
- [ ] Keep these helpers local to `webapi` rather than widening exported surfaces in `client` or `celersdk` during phase 1.
- [ ] Document in code comments or commit notes whether `paymentInfoFromClientPayment(...)` was moved, wrapped, or duplicated, so future refactors know the chosen source of truth.

### 6. OSP startup wiring

- [ ] Add `-webapigrpc <host:port>` to [server/server.go](../server/server.go).
- [ ] Add a `setUpOspWebApiService(...)` helper near `setUpAdminService(...)`.
- [ ] Construct the OSP pay backend, callback mux, event feed, and `webapi.OspPayApiServer` during startup.
- [ ] Register only `rpc.RegisterWebApiServer(...)` on the new listener.
- [ ] Do not register `rpc.RegisterInternalWebApiServer(...)` on the OSP listener.
- [ ] Fail startup if `-webapigrpc` is set and the listener cannot bind.
- [ ] Leave startup unchanged when `-webapigrpc` is unset.

## Routing-Behavior Verification Gate

- [ ] Verify that the OSP send path through `cnode.AddBooleanPay(...)` still prepends the hash-lock condition automatically for multi-hop routing when the destination is not a direct peer.
- [ ] Verify that the OSP send path still uses the `direct_pay` fast path where applicable for direct peers.
- [ ] Treat this verification as a required gate before calling phase 1 complete.
- [ ] Until item 13 lands, accept that `direct_pay` verification may rely on `pem` log inspection or other targeted debug instrumentation, because WebAPI does not currently expose `direct_pay` state.

## Deferred Methods

- [x] `GetBalance` remains deferred in phase 1.
- [x] `Deposit` remains deferred in phase 1.
- [x] `DepositNonBlocking` remains deferred in phase 1.
- [x] `MonitorDepositJob` remains deferred in phase 1.
- [x] `CooperativeWithdraw` remains deferred in phase 1.
- [x] `CooperativeWithdrawNonBlocking` remains deferred in phase 1.
- [x] `MonitorCooperativeWithdrawJob` remains deferred in phase 1.
- [x] `OpenPaymentChannel` remains deferred in phase 1.
- [x] `IntendWithdraw` remains deferred in phase 1.
- [x] `ConfirmWithdraw` remains deferred in phase 1.
- [x] `IntendSettlePaymentChannel` remains deferred in phase 1.
- [x] `ConfirmSettlePaymentChannel` remains deferred in phase 1.
- [x] `GetSettleFinalizedTimeForPaymentChannel` remains deferred in phase 1.
- [x] App-session lifecycle RPCs remain deferred in phase 1.
- [x] Debug/testing helpers like `SetMsgDropper` remain deferred in phase 1.

## Future Production Work Beyond Phase 1

- [ ] Phase 2: add explicit peer-scoped or channel-scoped OSP WebAPI APIs for the currently ambiguous methods.
- [ ] Phase 2: define the exact selector model for OSP channel operations. Candidates are explicit `peer_address`, explicit `channel_id`, or both depending on the RPC.
- [ ] Phase 2: add a production-safe OSP balance query, likely a peer-scoped balance/free-capacity RPC such as `GetPeerFreeBalance(peer, tokenInfo)`.
- [ ] Phase 2: add production-safe OSP variants for deposit and cooperative withdraw only after the selector model is locked.
- [ ] Phase 2: decide whether any existing channel-scoped `WebApi` RPCs should be extended for OSP use, or whether OSP-specific peer-scoped RPCs are cleaner.
- [ ] Phase 3: decide whether this listener ever needs to be exposed beyond loopback/private deployment.
- [ ] Phase 3: if public exposure is required, add the appropriate transport/auth story instead of relying on plaintext localhost/private gRPC.
- [ ] Phase 4: revisit subscription semantics if OSP needs multi-subscriber or lossless event delivery rather than the phase-1 single-subscriber best-effort model.
- [ ] Phase 4: revisit observability for `direct_pay` once item 13 lands, so OSP send-path verification can move from log-based checks to API-visible state.

## Testing Plan

### Unit tests

- [ ] Add `server/payment_callbacks_mux_test.go`.
- [ ] Test that receive-start hits both OSP sink and WebAPI sink.
- [ ] Test that receive-done hits both sinks.
- [ ] Test that send-complete hits both sinks.
- [ ] Test that destination-unreachable hits both sinks.
- [ ] Test that send-fail hits both sinks.

### Focused OSP WebAPI e2e

- [ ] Add a focused OSP WebAPI test file, proposed path: `test/e2e/osp_webapi.go`.
- [ ] Add a small dial helper in [testing/clientcontroller.go](../testing/clientcontroller.go) so tests can connect to an arbitrary WebAPI gRPC address without spawning `testing/testclient`.
- [ ] Start an OSP with `-webapigrpc 127.0.0.1:<test-port>`.
- [ ] Dial that port with a normal `rpc.WebApiClient`.
- [ ] Send a pay involving the OSP as the local sender or receiver.
- [ ] Assert `GetIncomingPaymentStatus` works.
- [ ] Assert `GetIncomingPaymentInfo` works and preserves current incoming-only semantics.
- [ ] Assert `RejectIncomingPayment` or `ConfirmOutgoingPayment` works on the OSP listener.
- [ ] Assert `SubscribeIncomingPayments` or `SubscribeOutgoingPayments` can observe the corresponding event.
- [ ] Assert one deferred channel-scoped RPC such as `GetBalance` returns `codes.Unimplemented`.
- [ ] Do not make the subscription event the only correctness signal in tests; polling state must remain the source of truth because slow subscribers may drop events.

### Suggested validation commands after code lands

- [ ] `go test ./server/... ./webapi/... ./testing/... -count=1`
- [ ] `go test ./test/e2e -run '^TestE2E$/^e2e-grp2$/^ospWebApiPaySubset$' -count=1`

## agent-pay-x402 Coordination After Branch Is Ready

- [ ] Start seller OSP with `-webapigrpc 127.0.0.1:<port>` in the x402 seller-OSP topology.
- [ ] Dial that port through the existing `rpc.WebApiClient` path.
- [ ] Keep seller-OSP `CooperativeWithdraw` on the Admin RPC path.
- [ ] Adjust x402 seller-OSP tests so they do not rely on OSP WebAPI `GetBalance` in phase 1.
- [ ] Use Admin-side balance visibility or drop the balance-delta assertion for the OSP-topology variant, as appropriate on the x402 side.
- [ ] Once the branch is ready, ask the x402 side to run `TestTopology_SellerOSP_Unconditional` and `TestTopology_SellerOSP_Conditional` as cross-repo acceptance.

## Rollout Order

- [ ] Add the `-webapigrpc` flag and empty listener scaffolding.
- [ ] Land the callback mux and its unit test.
- [ ] Land the payment event feed and OSP pay API server.
- [ ] Land the OSP pay backend in `server/`.
- [ ] Wire the listener into startup.
- [ ] Add the focused e2e path.
- [ ] Update operator docs.
- [ ] Re-run focused validation after each non-trivial checkpoint.

## Docs To Update After Landing

- [ ] Update [docs/backend-implementation.md](./backend-implementation.md) with the new optional OSP WebAPI gRPC listener.
- [ ] Update [docs/backend-usage.md](./backend-usage.md) with the `-webapigrpc` flag and example startup command.
- [ ] Update [docs/backend-usage.md](./backend-usage.md) to state explicitly that phase-1 OSP WebAPI is intended for localhost/private same-host callers and is not a public network-facing API.
- [ ] Update [docs/backend-usage.md](./backend-usage.md) to state explicitly that OSP WebAPI `GetBalance` is unsupported in phase 1 and balance observation remains available through Admin.
- [ ] Update [docs/backend-troubleshooting.md](./backend-troubleshooting.md) with guidance on binding the listener to loopback/private addresses only.
- [ ] Update [docs/backend-troubleshooting.md](./backend-troubleshooting.md) with a note that slow OSP WebAPI subscribers may drop events and polling status/info RPCs should be used for correctness checks.

## Review Decisions Closed By Current Feedback

- [x] `-webapigrpc <host:port>` is the accepted flag shape.
- [x] Single-subscriber semantics are acceptable for phase 1.
- [x] Plaintext gRPC is acceptable for phase 1 on a loopback/private listener.
- [x] Returning `codes.Unimplemented` for deferred methods is acceptable for the first upstreamable cut.
- [x] The x402 side is ready to use Admin RPC as the fallback path for seller-OSP `CooperativeWithdraw`.
- [x] The x402 side is willing to run seller-OSP integration acceptance against the branch once the phase-1 listener is ready.

If every unchecked item in this document is complete, phase 1 is complete and implementation can proceed to review.
