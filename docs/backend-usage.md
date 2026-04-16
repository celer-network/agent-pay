# AgentPay Backend Usage

## What This Guide Covers

This guide is for developers and operators who need to build, run, test, and inspect the backend nodes in this repo.

It covers three common workflows:

1. Run a focused end-to-end test to validate code changes.
2. Bring up a local multi-OSP environment manually.
3. Start a backend node directly with your own profile and storage settings.

## Prerequisites

You need the following on a development machine:

- Go
- `geth` for local Ethereum-based tests
- CockroachDB only if you want the shared-SQL mode used in the manual examples

For the manual scripts, set:

```bash
export AGENTPAY=$PWD
export AGENTPAY_MANUAL_ROOT=${AGENTPAY_MANUAL_ROOT:-/tmp/celer_manual_test}
```

Useful assets already in the repo:

- Example profile schema: [test/manual/sample_profile.json](../test/manual/sample_profile.json)
- Test keystores: [testing/env/keystore](../testing/env/keystore)
- Runtime config examples: [testing/profile](../testing/profile) and [test/manual/rt_config.json](../test/manual/rt_config.json)
- Operational troubleshooting guide: [docs/backend-troubleshooting.md](./backend-troubleshooting.md)

## Build the Binaries

From the repo root:

```bash
mkdir -p ./bin
go build -o ./bin/server ./server
go build -o ./bin/osp-cli ./tools/osp-cli
```

Optional entry points you may also care about:

- `go build ./webapi/cmd`
- `go build ./webproxy/cmd`

## Fastest Validation: Focused E2E Test

If you are changing the backend and want the shortest realistic validation loop, start with the end-to-end tests in [test/e2e](../test/e2e).

From the repo root:

```bash
go test ./test/e2e -run '^TestE2E$/^e2e-grp2$/^sendCondPayWithErc20$'
```

If you are already inside `test/e2e`, the shorter form also works:

```bash
go test -run '^TestE2E$/^e2e-grp2$/^sendCondPayWithErc20$'
```

What this test setup does for you automatically:

- Starts a local geth-based chain
- Builds the required binaries into a temp output directory
- Deploys contracts and funds test accounts
- Generates runtime profiles under `/tmp/celer_e2e_*`
- Registers the OSP router used by the tests
- Starts the default OSP/backend process

Important files behind that workflow:

- [test/e2e/e2e_setup_test.go](../test/e2e/e2e_setup_test.go)
- [test/e2e/e2e_test.go](../test/e2e/e2e_test.go)
- [test/e2e/constants.go](../test/e2e/constants.go)

The default single-network e2e flow does not provision the extra networks required by the cross-net suite. Run cross-net explicitly with:

```bash
go test ./test/e2e -run '^TestE2ECrossNet$' -args -multinet
```

Useful debugging behavior:

- Successful runs delete the temp directory.
- Failed runs keep it and print a `-reuse` path so you can rerun without rebuilding or redeploying.

Example:

```bash
go test ./test/e2e -reuse /tmp/celer_e2e_1712960000/ -run '^TestE2E$/^e2e-grp2$/^sendCondPayWithErc20$'
```

## Broader Test Matrix

For a wider validation sweep, the old CI flow maps reasonably well to the following current-package commands.

Prerequisites beyond Go:

- `geth` for e2e suites
- `sqlite3` CLI for storage-related test helpers and inspection flows

Legacy CI-style unit/package sweep:

```bash
go test ./storage ./celersdk ./common/cobj ./dispatchers ./lrucache ./rpc ./rtconfig ./metrics ./route ./utils/bar
```

Recommended validation tiers:

- Payment-path regression check:

```bash
go test ./test/e2e -run '^TestE2E$/^e2e-grp2$/^sendCondPayWithErc20$'
```

- Cross-net routing check:

```bash
go test ./test/e2e -run '^TestE2ECrossNet$' -args -multinet
```

- Manual multi-OSP smoke flow: use [test/manual/README.md](../test/manual/README.md) or run `AGENTPAY=$PWD ./test/manual/smoke.sh`

The full `go test ./test/e2e` package includes broader multi-OSP and specialized integration suites in addition to the core single-network flow. Use the targeted commands above when you want predictable validation for a specific area.

## Manual Multi-OSP Workflow

The best operator-oriented walkthrough already in the repo is [test/manual/README.md](../test/manual/README.md). The steps below summarize it and point to the files that matter.

### 1. Prepare the environment

```bash
export AGENTPAY=$PWD
go build -o ./osp-cli ./tools/osp-cli
cd test/manual
```

### 2. Start the local chain and generate profiles

```bash
./setup.sh
```

This does more than just start geth. It also:

- deploys the ledger, resolver, registry, wallet, and ERC20 contracts
- funds test accounts
- writes OSP profiles under `$AGENTPAY_MANUAL_ROOT/profile/`

See [test/manual/setup.go](../test/manual/setup.go) and [test/manual/sample_profile.json](../test/manual/sample_profile.json).

### 3. Fund and register OSPs

From the repo root or from `test/manual` with the built CLI available:

```bash
./osp-cli -profile $AGENTPAY_MANUAL_ROOT/profile/o1_profile.json \
  -ks $AGENTPAY/testing/env/keystore/osp1.json \
  -ethpooldeposit -amount 10000 -register -nopassword

./osp-cli -profile $AGENTPAY_MANUAL_ROOT/profile/o2_profile.json \
  -ks $AGENTPAY/testing/env/keystore/osp2.json \
  -ethpooldeposit -amount 10000 -register -nopassword
```

This is required if you want route-controller behavior that depends on on-chain router registration.

### 4. Start OSP nodes

SQLite-backed example:

```bash
./run_osp.sh 1
./run_osp.sh 2
```

For localhost manual runs, `test/manual/run_osp.sh` defaults `CELER_INSECURE_TLS=1` so inter-OSP dials work with the built-in self-signed localhost certificate.

CockroachDB-backed example:

```bash
./cockroachdb.sh start
./cockroachdb.sh 1
./cockroachdb.sh 2
./run_osp.sh 1_crdb
./run_osp.sh 2_crdb
```

See [test/manual/run_osp.sh](../test/manual/run_osp.sh) for the exact flags passed to the server.

### 5. Connect OSPs and exercise the payment path

Register an inter-OSP stream:

```bash
./osp-cli -adminhostport localhost:8190 \
  -registerstream \
  -peer 00290a43e5b2b151d530845b2d5a818240bc7c70 \
  -peerhostport localhost:10002
```

Open an OSP-to-OSP channel:

```bash
./osp-cli -adminhostport localhost:8190 \
  -openchannel \
  -peer 00290a43e5b2b151d530845b2d5a818240bc7c70 \
  -selfdeposit 10 \
  -peerdeposit 10
```

Send an off-chain payment:

```bash
./osp-cli -adminhostport localhost:8190 \
  -sendtoken \
  -receiver 00290a43e5b2b151d530845b2d5a818240bc7c70 \
  -amount 0.01
```

Inspect state:

- off-chain DB queries with [tools/osp-cli/README.md](../tools/osp-cli/README.md)
- on-chain queries with the CLI's `-onchainview` options

## Running a Backend Node Directly

You do not need the helper scripts if you already have a profile and keys.

Example command from the repo root:

```bash
go run ./server/server.go \
  -profile $AGENTPAY_MANUAL_ROOT/profile/o1_profile.json \
  -ks ./testing/env/keystore/osp1.json \
  -port 10001 \
  -adminrpc localhost:11001 \
  -adminweb localhost:8190 \
  -svrname o1 \
  -storedir $AGENTPAY_MANUAL_ROOT/store \
  -rtc ./test/manual/rt_config.json \
  -nopassword
```

If this process will dial localhost peers using the built-in localhost certificate, prefix the command with `CELER_INSECURE_TLS=1` unless you are using `test/manual/run_osp.sh`, which already does that for local manual runs.

For a CockroachDB-backed node, replace `-storedir` with `-storesql`:

```bash
-storesql 'postgresql://celer_test_o1@localhost:26257/celer_test_o1?sslmode=disable'
```

## Configuration Files

### Profile JSON

The profile schema is defined in [common/profile.go](../common/profile.go). The main sections are:

- `Ethereum`: RPC gateway, chain id, block timing, and contract addresses
- `Osp`: this node's gRPC host and ETH address
- `Sgn`: SGN-related endpoints and contract address

Example: [test/manual/sample_profile.json](../test/manual/sample_profile.json)

### Runtime config JSON

The runtime config file passed by `-rtc` is separate from the profile. It controls operational values such as:

- min/max payment timeouts
- refill thresholds and refill amounts
- deposit polling and batching
- OSP-to-OSP open-channel limits

Examples:

- [testing/profile/rt_config.json](../testing/profile/rt_config.json)
- [testing/profile/rt_config_multiosp.json](../testing/profile/rt_config_multiosp.json)
- [test/manual/rt_config.json](../test/manual/rt_config.json)

## Server Flags That Matter Most

| Flag | Meaning |
| --- | --- |
| `-profile` | Chain, contract, and OSP profile |
| `-ks` | Main keystore for signing and transactions |
| `-depositks` | Optional separate deposit signer |
| `-storedir` | Local SQLite storage root |
| `-storesql` | Shared SQL store URL |
| `-port` | Main gRPC endpoint for clients and peers |
| `-adminrpc` | Admin gRPC endpoint |
| `-adminweb` | Admin HTTP endpoint that serves `/admin/` and `/metrics` |
| `-selfrpc` | Second gRPC endpoint used in multi-server mode |
| `-rtc` | Runtime config file |
| `-isosp` | Whether to run with OSP/service-node behavior |
| `-loc` | Whether this process listens to on-chain logs |
| `-tlscert`, `-tlskey`, `-tlsclient` | TLS customization |

Only one of `-storedir` and `-storesql` should be set.

## Deployment Modes

### Single-server mode

This is the default and easiest setup:

- one process owns its peers directly
- storage is local or at least logically local to that process
- forwarding never leaves the process boundary

This is what most e2e tests and the simple manual SQLite workflow use.

### Multi-server mode

This mode is enabled when the profile or flags provide both shared SQL storage and `SelfRPC`.

In that mode:

- multiple server processes share storage
- the process exposes the `MultiServer` gRPC service
- a message may be forwarded to another server if the target client is connected there

The implementation lives in [cnode/multiserver.go](../cnode/multiserver.go).

## Admin and Operator Interfaces

In practice, most operational control happens through the admin surface exposed by [server/server.go](../server/server.go):

- admin gRPC server on `-adminrpc`
- HTTP gateway on `-adminweb`, mounted under `/admin/`
- Prometheus metrics on `/metrics`

The normal operator tool for that surface is [tools/osp-cli](../tools/osp-cli).

Common admin actions:

- register a peer stream
- open an OSP-to-OSP channel
- send a payment
- start and query deposits
- inspect off-chain and on-chain state

Full command reference: [tools/osp-cli/README.md](../tools/osp-cli/README.md)

## Embedding a Client

If you are using this backend from application code instead of running only OSP nodes, there are two relevant entry points:

- [client/celer_client.go](../client/celer_client.go) for direct Go integration
- [celersdk/api.go](../celersdk/api.go) for the higher-level SDK interface

The standard client flow is:

1. Create the client with a profile and keystore.
2. Register a stream to the server OSP.
3. Open or instantiate a channel.
4. Deposit, withdraw, and send payments through the SDK/client APIs.

These clients still use the same backend protocol pipeline and storage model described in the implementation guide.

## WebAPI Notes

- `WebApi.SendToken` is the explicit alias for sending a payment without caller-specified app conditions.
- `WebApi.SendConditionalPayment` remains the lower-level payment API when you want to attach app-level conditions, or when you want to pass an empty `conditions` list explicitly.
- Even with empty `conditions`, the runtime may still prepend an internal hash-lock condition for non-direct pays.
- Public `WebApi.Deposit` and `WebApi.CooperativeWithdraw` are blocking calls that return after the transaction is mined.
- `WebApi.DepositNonBlocking` and `WebApi.CooperativeWithdrawNonBlocking` start jobs that can be tracked with `MonitorDepositJob` and `MonitorCooperativeWithdrawJob` on the public surface.
- The matching `InternalWebApi` non-blocking variants remain available for internal callers.

## Practical Notes

- The e2e tests set `CELER_INSECURE_TLS=1` so localhost clients can talk to the server's built-in localhost certificate without CA setup.
- OSP routing behavior only becomes meaningful after the OSP is registered in the on-chain `RouterRegistry`.
- The server starts periodic OSP cleanup that clears expired or on-chain-resolved payments with peer OSPs.
- `rtconfig` is operationally important. Payment timeout, refill, and deposit behavior are not hardcoded solely in Go constants.

## Suggested Reading Path

If you are new to the repo, this order works well:

1. [docs/backend-implementation.md](./backend-implementation.md)
2. [test/manual/README.md](../test/manual/README.md)
3. [tools/osp-cli/README.md](../tools/osp-cli/README.md)
4. [docs/backend-troubleshooting.md](./backend-troubleshooting.md)
5. [test/e2e](../test/e2e) for executable examples