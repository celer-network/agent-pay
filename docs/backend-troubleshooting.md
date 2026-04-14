# AgentPay Backend Troubleshooting and Operations

## Scope

This guide is for diagnosing and operating the backend nodes in this repo after they have been built or deployed. It focuses on the failures that show up most often in this codebase:

- startup and configuration failures
- stream-registration failures
- routing failures
- deposit and refill issues
- payment-path failures

It complements the usage guide in [backend-usage.md](./backend-usage.md) and the implementation guide in [backend-implementation.md](./backend-implementation.md).

## First Checks

Before going deep into a specific failure, verify the basics in this order:

1. The server process is running and listening on the expected ports.
2. The node profile points at the correct chain RPC and contract addresses.
3. The storage backend is reachable and writable.
4. Peer streams are connected.
5. Routing exists for the intended token and destination.
6. The channel exists and has enough usable balance.

Useful references:

- [server/server.go](../server/server.go)
- [tools/osp-cli/README.md](../tools/osp-cli/README.md)
- [test/manual/README.md](../test/manual/README.md)

## Operational Surfaces

The backend exposes three main operator surfaces:

- Main gRPC server on `-port`
- Admin gRPC server on `-adminrpc`
- Admin HTTP gateway and Prometheus metrics on `-adminweb`

The normal operator tool is [tools/osp-cli](../tools/osp-cli). Prefer it over ad hoc RPC calls because the repo already documents the stable command patterns there.

Useful operator commands:

```bash
./osp-cli -adminhostport localhost:8190 -querypeerosps
./osp-cli -adminhostport localhost:8190 -querydeposit -depositid <deposit-id>
./osp-cli -profile <profile.json> -storedir <store> -dbview channel -peer <peer-addr>
./osp-cli -profile <profile.json> -storedir <store> -dbview pay -payid <pay-id>
./osp-cli -profile <profile.json> -storedir <store> -dbview route -dest <dest-addr> -token <token-addr>
./osp-cli -profile <profile.json> -onchainview channel -cid <cid>
./osp-cli -profile <profile.json> -onchainview pay -payid <pay-id>
```

## Startup and Configuration Failures

### Symptom: server exits immediately on startup

Start with [server/server.go](../server/server.go), where the process validates its flags and storage selection.

Common causes:

- both `-storedir` and `-storesql` were set
- `-selfrpc` is malformed
- keystore cannot be read or decrypted
- chain RPC endpoint is unreachable
- storage backend cannot be opened

Representative log messages from the code:

- `specify only one of -storedir, -storesql`
- `invalid self-RPC`
- `Cannot setup SQL store`
- `Cannot setup local store`
- `DialETH failed`

Checks:

1. Ensure exactly one storage mode is selected.
2. Confirm the profile file is valid against the schema in [common/profile.go](../common/profile.go).
3. Verify the keystore path and password handling. For local tests, `-nopassword` is commonly used.
4. Verify the chain RPC endpoint in the profile's `Ethereum.Gateway` field.

Known-good local example:

```bash
go run ./server/server.go \
  -profile /tmp/celer_manual_test/profile/o1_profile.json \
  -ks ./testing/env/keystore/osp1.json \
  -port 10001 \
  -adminrpc localhost:11001 \
  -adminweb localhost:8190 \
  -storedir /tmp/celer_manual_test/store \
  -rtc ./test/manual/rt_config.json \
  -nopassword
```

### Symptom: server starts, but storage state is missing or unexpected

The storage path is derived differently for local and SQL modes.

- In `-storedir` mode, the actual SQLite path becomes `<storedir>/<ethaddr>/sqlite/celer.db`.
- In `-storesql` mode, the node uses the configured SQL database directly.

See `setupKVStore(...)` in [cnode/cnode.go](../cnode/cnode.go).

Checks:

1. Make sure you are querying the store for the node's actual ETH address, not only the parent directory.
2. When using SQLite, confirm you are pointing `osp-cli -storedir` at the node-specific store directory when running DB views.
3. When using SQL, verify the exact database name and credentials used by the server match the CLI invocation.

## Stream Registration Failures

### Symptom: `registerstream` fails or peers never connect

The admin entry point is `RegisterStream(...)` in [server/server.go](../server/server.go), which calls `CNode.RegisterStream(...)` in [cnode/cnode.go](../cnode/cnode.go).

Representative failures:

- `celer stream already exists`
- `RegisterStream failed: grpcDial ... failed`
- `RegisterStream failed: CelerStream failed`
- `waitRecvWithTimeout failed`
- `no celer stream`
- `peer not online`

What these usually mean:

- `celer stream already exists`: the server already has a live or remembered stream for that peer and RPC address.
- `grpcDial ... failed`: the target host or port is wrong, the peer process is not listening, or TLS/networking is broken.
- `waitRecvWithTimeout failed`: the transport connected, but the auth handshake did not complete.
- `peer not online` or `no celer stream`: later traffic depends on a stream that was never established or was dropped.

Checks:

1. Verify the peer gRPC port, not the admin port, is being passed to `-peerhostport`.
2. Confirm the peer ETH address matches the profile and keystore used by that peer.
3. Check whether the stream already exists before retrying the same registration.
4. Use `-querypeerosps` to see what the node currently believes about peer OSPs.

Example command:

```bash
./osp-cli -adminhostport localhost:8190 \
  -registerstream \
  -peer 00290a43e5b2b151d530845b2d5a818240bc7c70 \
  -peerhostport localhost:10002
```

If stream registration succeeds once and fails later, remember that the server installs a retry callback and may reconnect automatically after transient failures.

## Routing Failures

### Symptom: payment cannot find a route

Routing lookup is implemented in [route/forwarder.go](../route/forwarder.go). The common terminal error is `no route to destination`, but routing problems also surface indirectly as send failures or unreachable peers.

Checks:

1. Confirm a channel exists either directly to the destination or to an access OSP for that token.
2. Query the route table with `osp-cli -dbview route`.
3. Confirm the token address used in the send matches the token address used in the channel and route tables.
4. If you expect OSP routing, verify the node is registered as a router on-chain.

Example:

```bash
./osp-cli -profile <profile.json> -storedir <store> \
  -dbview route -dest <destination-addr> -token <token-addr>
```

### Symptom: OSP starts but does not join the routing network

The route controller logs this warning from [route/controller.go](../route/controller.go):

`NOT able to join the OSP network because this node is not registered on-chain as a router`

That means the node process is healthy, but the on-chain `RouterRegistry` does not show it as an active router.

Recovery:

```bash
./osp-cli -profile <profile.json> -ks <keystore.json> -register -nopassword
```

Then restart the node or wait for the route-controller logic to observe the registry state.

### Symptom: routing looks stale after restart

Checks:

1. Confirm `-loc` is enabled if this process is expected to listen to on-chain events.
2. Confirm the profile points to the expected `RouterRegistry` contract.
3. In multi-OSP setups, ensure peers have exchanged streams and routing broadcasts.
4. Verify the runtime network actually has open channels for the token you are testing.

## Deposit and Refill Issues

### Symptom: deposit request never finishes

Deposit jobs are tracked by the processor in [deposit/deposit.go](../deposit/deposit.go) and queried through admin RPC in [server/server.go](../server/server.go).

Possible states include:

- `QUEUED`
- `APPROVING_ERC20`
- `TX_SUBMITTING`
- `TX_SUBMITTED`
- `SUCCEEDED`
- `FAILED`

Checks:

1. Query the deposit job explicitly.
2. If the token is ERC20, look for an approval phase before the ledger deposit.
3. Confirm the deposit signer has funds and is the expected keystore.
4. Confirm the process is running as an event listener if you expect server-side job polling to progress automatically.

Commands:

```bash
./osp-cli -adminhostport localhost:8190 -querydeposit -depositid <deposit-id>
./osp-cli -profile <profile.json> -storedir <store> -dbview deposit -depositid <deposit-id>
```

If the job is missing entirely, the admin query may return `deposit job not found`.

### Symptom: payment send fails because balance is too low even though the channel exists

During send-path execution, [messager/send_cond_pay_request.go](../messager/send_cond_pay_request.go) computes free balance from the working simplex state and on-chain balance view. The common error is `balance not enough`.

Checks:

1. Inspect the channel with `osp-cli -dbview channel` and confirm free balance, not only total deposited balance.
2. If this is an OSP, check whether refill thresholds in the runtime config are causing automatic refill behavior.
3. Confirm there are not too many unresolved pending payments consuming available capacity.

Relevant runtime config examples:

- [testing/profile/rt_config.json](../testing/profile/rt_config.json)
- [testing/profile/rt_config_multiosp.json](../testing/profile/rt_config_multiosp.json)

## Payment-Path Failures

### Symptom: send fails from admin or SDK

Admin sends use `SendToken(...)` in [server/server.go](../server/server.go). Common immediate failures are:

- `Can't parse amount.`
- `Can't parse dst.`
- `Can't parse token address.`
- `no celer stream`
- `no route to destination`
- `balance not enough`
- `invalid pay resolve deadline`

Checks:

1. Validate the receiver and token addresses before retrying.
2. Confirm the stream and route exist.
3. Confirm the current block height and `rtconfig` payment-timeout settings are consistent with the send path.
4. Confirm the destination is reachable on the intended network.

### Symptom: receive-side logs show `invalid sequence number`

This comes from the simplex sliding-window protocol in the message handlers. It typically means one of the following:

- request loss or replay
- sender and receiver disagree on the last co-signed simplex state
- a later request arrived before the expected base sequence was acknowledged

Relevant code paths:

- [handlers/msghdl/handle_cond_pay_request.go](../handlers/msghdl/handle_cond_pay_request.go)
- [handlers/msghdl/handle_pay_settle_request.go](../handlers/msghdl/handle_pay_settle_request.go)

Checks:

1. Look for earlier ACK, NACK, reconnect, or dropped-stream events for the same peer.
2. Inspect channel state with `osp-cli -dbview channel` and compare simplex sequence numbers across both peers.
3. If this happened after a restart or network interruption, re-establish the stream and retry the operation.

### Symptom: pay or channel is reported missing during settlement

Common errors include:

- `payment not found`
- `channel not found`
- `channel simplex state not found`

These generally point to one of three situations:

- wrong store or wrong node profile is being queried
- the channel/payment was never created on this node's side
- the operator is looking at ingress/egress state on the wrong hop

Checks:

1. Query the payment by pay ID in the node that should own ingress or egress state.
2. Query the channel by peer and token or by CID.
3. Confirm that the profile and storage path used by `osp-cli` match the exact node instance you are debugging.

## Recovery Playbook

For most operational issues, use this order instead of trying random retries:

1. Confirm the process, ports, and profile are correct.
2. Confirm storage access and the correct store path.
3. Re-register missing peer streams.
4. Re-check route table state.
5. Re-check channel state and free balance.
6. Re-check deposit status and on-chain state.
7. Retry the payment or channel operation only after the earlier layers look correct.

This sequence matches the way the backend is structured: transport first, then routing, then channel state, then payment state.

## Useful File Map

- Startup and admin surfaces: [server/server.go](../server/server.go)
- Core node initialization: [cnode/cnode.go](../cnode/cnode.go)
- Stream auth and registration: [cnode/auth.go](../cnode/auth.go) and [cnode/cnode.go](../cnode/cnode.go)
- Routing lookup and controller: [route/forwarder.go](../route/forwarder.go) and [route/controller.go](../route/controller.go)
- Payment send path: [messager/send_cond_pay_request.go](../messager/send_cond_pay_request.go)
- Payment receive and settlement handlers: [handlers/msghdl](../handlers/msghdl)
- Deposit job processing: [deposit](../deposit)
- Shared errors: [common/errs.go](../common/errs.go)