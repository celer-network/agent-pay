# Project Guidelines

## Read This First

This repository implements the off-chain backend of AgentPay. Before changing code, read these documents in order:

1. Protocol reference in the companion docs repo `agent-pay-docs` under `agentpay-architecture/`:
   - `system-overview.md`
   - `off-chain-protocols/README.md`
   - `off-chain-protocols/single-hop-protocols.md`
   - `off-chain-protocols/end-to-end-protocols.md`
   - `app-contracts-and-protocols.md`
2. Repo overview in `README.md`.
3. Backend design map in `docs/backend-implementation.md`.
4. Build and runtime workflow in `docs/backend-usage.md`.
5. Failure handling in `docs/backend-troubleshooting.md`.

Companion repos are wired in via the filesystem MCP server. The canonical roots are:

- `agent-pay-docs` — protocol architecture (`agentpay-architecture/`). Required reading for protocol-sensitive work. The `state-guardian-network/` subtree is **not** required background; read it only when the task is explicitly about SGN behavior or SGN-related profile wiring.
- `agent-pay-contracts` — Solidity contracts (CelerLedger, PayResolver, PayRegistry, RouterRegistry, EthPool, Wallet). Optional background; read only when the task touches on-chain contract logic, event semantics, generated bindings under `chain/channel-eth-go/`, or profile/address wiring that depends on contract behavior.
- `agent-pay-x402` — downstream integration that consumes this repo via WebAPI gRPC + Admin HTTP. Useful as an "external consumer" reference, not required reading.

`.mcp.json` is gitignored per-developer. Copy `.mcp.json.example` to `.mcp.json` and fill in absolute paths to your local sibling clones. If a path is unavailable when needed, ask the user before guessing on protocol-sensitive work.

## Architecture

- Process entry point: `server/server.go`.
- Core runtime object: `cnode.CNode` in `cnode/cnode.go`.
- Message ingress path: `rpc.CelerStream` -> `dispatchers` -> `handlers/msghdl` -> `storage.DAL` and `messager`.
- Payment egress path: API or client call -> `cnode/pay.go` and related helpers -> `messager` -> peer stream -> handler validation.
- Routing and relay behavior: `route/controller.go` and `route/forwarder.go`.
- On-chain fallback and background processors: `deposit`, `dispute`, `cnode/cooperativewithdraw`, `migrate`, `route`.
- Client-facing wrappers: `client` and `celersdk`.

Keep `server/server.go` thin. New protocol logic normally belongs in `cnode`, `handlers/msghdl`, `messager`, or the relevant processor package.

## Protocol Invariants

- A payment channel is modeled as two independent simplex directions. Preserve the duplex design and do not collapse behavior into a shared bidirectional state machine.
- Only `peer_from` initiates simplex-state updates. Sequence number and base-sequence handling must remain compatible with the sliding-window protocol.
- Conditional pay setup and settlement must keep the pending-pay list, transferred amount, and co-signed simplex state consistent across storage and outbound messages.
- Relay nodes should stay agnostic to application logic. Do not add relay behavior that interprets app outcomes unless the protocol explicitly requires it.
- Boolean end-to-end payments should not require relay-side on-chain actions. Numeric payments may require registry checks or disputes only where the protocol already allows them.
- Channel and payment mutations that belong to one protocol step should stay inside the existing `storage.DAL` transaction boundaries.
- Multi-server mode changes must preserve client ownership and forwarding behavior in `cnode/multiserver.go`.
- All on-chain deadlines / challenge windows / timeouts are unix seconds — the contracts compare against `block.timestamp`, not `block.number`. Off-chain code uses `time.Now().Unix()` to produce and check them. This applies to `disputeTimeout`, `settleFinalizedTime`, `withdrawDeadline`, `openDeadline`, `resolveDeadline`, `resolveTimeout`, `migrationDeadline`, the `RouterRegistry` register/refresh value, and per-pay deadlines in `PayRegistry`. The one surviving exception is `testing/testapp/singlesessionapp.go` — kept as an `agent-pay-x402` back-compat carry — which still uses `block.number`-based deadlines internally; it is never queried in agent-pay's own off-chain logic, so this exception is contained to that one file. See [docs/progress/app-session-simplification.md](docs/progress/app-session-simplification.md) §7 for the deferred follow-up that retires it.

## Conventions

- Prefer changing source definitions instead of generated outputs. Files ending in `.pb.go` and generated contract bindings under `chain/channel-eth-go/` should only be edited when the corresponding source change and regeneration are part of the task.
- When changing message behavior, read the sender and receiver sides together: `messager/*`, `dispatchers/*`, and `handlers/msghdl/*`.
- Keep protocol changes aligned with `proto/*.proto`, runtime structs, persistence, and tests.
- Preserve existing logging and metrics patterns around critical protocol transitions.
- Use `github.com/celer-network/goutils/log` for all logging — leveled package-level API (`log.Infof`, `log.Warnf`, `log.Errorf`, `log.Fatalf`). For structured-ish output, embed `key=value` pairs in the format string. Do not pull in `log/slog` or stdlib `log`.
- Update docs when a change affects architecture, operator workflow, or protocol-to-code mapping.

## Build And Test

From the repo root:

```bash
mkdir -p ./bin
go build -o ./bin/server ./server
go build -o ./bin/osp-cli ./tools/osp-cli
```

Toolchain note: on macOS amd64 the local `go1.25.5` toolchain has been observed to fail with duplicate `runtime/cgo` symbols at link time (reproduces on trivial `import "C"` programs too). If you hit this, export `GOTOOLCHAIN=go1.24.9` for the shell before building or testing. SQLite-backed local builds also require cgo — `-storedir` startup fails fast with `sqlite3 requires cgo` if you build with `CGO_ENABLED=0`. See `docs/backend-troubleshooting.md` for the full diagnosis.

Fastest realistic payment-path validation:

```bash
go test ./test/e2e -run '^TestE2E$/^e2e-grp2$/^sendCondPayWithErc20$'
```

Cross-net validation is opt-in and should be run explicitly when changes may affect multi-network behavior:

```bash
go test ./test/e2e -run '^TestE2ECrossNet$' -args -multinet
```

When possible, run focused package tests for touched code first, then the focused e2e flow for protocol-path changes.

## Useful References

- `docs/backend-implementation.md` explains the protocol-to-code mapping.
- `test/e2e` is the best executable reference for end-to-end behavior.
- `test/manual/README.md` documents the operator workflow for multi-OSP setups.
- `tools/osp-cli/README.md` covers admin and inspection commands.
- `tools/scripts/README.md` covers maintenance and code-generation helpers.
