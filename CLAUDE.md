## AgentPay Backend (Claude Code notes)

**Read [AGENTS.md](AGENTS.md) first.** It is the shared coding-agent contract: reading order, architecture map, protocol invariants, conventions, and the build/test commands. This file only adds Claude-Code-specific tips on top — keep the substantive guidance in `AGENTS.md` so non-Claude agents (Codex, etc.) get it too.

## Quick Reference

- **What:** Off-chain Go/gRPC implementation of the Celer AgentPay state-channel network (OSP/service nodes, client SDK, routing, persistence, admin surface).
- **Language:** Go 1.24 (toolchain notes below).
- **Process entry:** [server/server.go](server/server.go) → [cnode.CNode](cnode/cnode.go).

## Key Docs

- [README.md](README.md) — repo landing page with the docs index.
- [AGENTS.md](AGENTS.md) — coding-agent entry point (reading order + invariants + conventions). **Start here.**
- [docs/backend-implementation.md](docs/backend-implementation.md) — runtime model, package map, protocol-to-code mapping.
- [docs/backend-usage.md](docs/backend-usage.md) — build, test, runtime, profile/rtconfig formats.
- [docs/backend-troubleshooting.md](docs/backend-troubleshooting.md) — failure diagnosis and recovery playbook.
- [test/manual/README.md](test/manual/README.md) — operator-oriented multi-OSP walkthrough.
- [tools/osp-cli/README.md](tools/osp-cli/README.md) — admin / on-chain / DB inspection commands.
- [tools/scripts/README.md](tools/scripts/README.md) — code-generation and CockroachDB helpers.

The protocol-and-contract design lives in the companion `agent-pay-docs` repo (`agentpay-architecture/`). AGENTS.md §"Read This First" lists the load-bearing files; read them before any protocol-sensitive change.

## External Repositories (via MCP)

`.mcp.json` is gitignored per-developer — each contributor wires their own absolute paths. The current wiring exposes three filesystem MCP roots:

- **agent-pay-docs** — protocol architecture + SGN docs. Read `agentpay-architecture/`. Skip `state-guardian-network/` unless the task is explicitly about SGN behavior or SGN-related profile wiring.
- **agent-pay-contracts** — Solidity contracts (AgentPayLedger, PayResolver, PayRegistry, RouterRegistry, Wallet, plus the chain-canonical wrapped-native / WETH-style contract that AgentPayLedger references for native-token funding flows). Foundry project. Read only when the task touches on-chain contract logic, event semantics, or generated bindings under [chain/channel-eth-go/](chain/channel-eth-go/).
- **agent-pay-x402** — downstream Go integration that layers x402 HTTP payment over AgentPay state channels. Useful as an "external consumer" reference: it talks to this repo's WebAPI gRPC (clients) and Admin HTTP (OSP) only, never `CelerStream` directly. Friction observed there is logged in that repo's `docs/agent-pay-feedback.md`.

If a path is unavailable in your MCP roots, ask the user before guessing — do not silently proceed on protocol-sensitive work.

## Proto sources of truth

Wire and admin contracts live in [proto/](proto) and [webapi/proto/](webapi/proto). Read these before changing message behavior — `messager/`, `dispatchers/`, and `handlers/msghdl/` only realize what `.proto` says.

- `proto/message.proto` — `CelerMsg` envelope, `CondPayRequest` / `CondPayResponse`, `PaymentSettleRequest` / `PaymentSettleResponse`.
- `proto/entity.proto` — `SimplexPaymentChannel`, `ConditionalPay`, `Condition`, `CooperativeSettleInfo`.
- `proto/rpc.proto` — `Rpc` service with `CelerStream` (bidirectional streaming) + the public WebApi service.
- `proto/osp_admin.proto` — Admin gRPC (stream registration, `OpenChannel`, `Deposit`, `SendToken`, `CooperativeSettle`).
- `webapi/proto/web_api.proto` — pay-centric WebAPI used by client nodes and the optional `-webapigrpc` listener.

## Toolchain notes

- **macOS amd64 cgo linker bug.** The local `go1.25.5` toolchain has been observed to fail with duplicate `runtime/cgo` symbols on this codebase (reproduces on trivial `import "C"` programs too — it's a toolchain issue, not a repo regression). Before `go build` / `go test`, export `GOTOOLCHAIN=go1.24.9` for the shell. Documented in [docs/backend-troubleshooting.md](docs/backend-troubleshooting.md).
- **CGO is required for SQLite-backed builds.** `-storedir` mode fails fast with `sqlite3 requires cgo` if you accidentally build with `CGO_ENABLED=0`. CI sets `CGO_ENABLED=1` explicitly.
- **`AGENTPAY_INSECURE_TLS=1` is normal for localhost.** Inter-OSP and client→OSP localhost dials use the built-in self-signed localhost cert; `test/manual/run_osp.sh` and the e2e harness already set this. Set it yourself when launching binaries directly against `localhost`/`127.0.0.1`.
- **e2e is slow.** CI gives the default suite 30 minutes (40-minute job cap) and crossnet 30/45. For a local validation loop, prefer the focused single-test runs in [AGENTS.md §Build And Test](AGENTS.md) over `go test ./test/e2e`. On failed runs the harness keeps `/tmp/celer_e2e_*` and prints a `-reuse` path — use it instead of paying the rebuild cost again.

## Logging convention

Use [`github.com/celer-network/goutils/log`](https://github.com/celer-network/goutils) — leveled package-level API: `log.Infof` / `log.Warnf` / `log.Errorf` / `log.Fatalf`. This is what the rest of the codebase uses (e.g. [cnode/cnode.go](cnode/cnode.go), [server/server.go](server/server.go)). For structured-ish output, embed `key=value` pairs in the format string; do **not** pull in `log/slog` or stdlib `log`.

## Useful slash commands

- `/review` — pull-request review helper after pushing a branch.
- `/security-review` — focused security review of pending changes. Run before PRs that touch payment-path code, signing, or anything in [cnode/](cnode), [messager/](messager), [handlers/msghdl/](handlers/msghdl), [dispute/](dispute), or [chain/](chain).
- `/ultrareview` — multi-agent cloud review (user-triggered, billed). Claude Code cannot launch this itself; the human invokes it.

## Effort estimates

You — Claude Code — are the developer here. Do not size work in human-developer time units ("an afternoon", "1–2 days"). Those numbers don't translate to LLM bottlenecks (context window, tool latency, review rounds) and have repeatedly been wrong by an order of magnitude on adjacent Celer codebases.

Use units that actually predict your cost:

- **Edit footprint** — files / packages touched, approximate LOC delta.
- **Conceptual scope** — number of distinct seams; whether `proto/*.proto`, storage schema, or generated bindings under `chain/channel-eth-go/` change; whether `agent-pay-contracts` or `agent-pay-docs` need parallel changes.
- **Test surface** — focused package test vs focused e2e (`sendCondPayWithErc20`) vs full `go test ./test/e2e` vs manual multi-OSP harness.
- **Review risk** — likelihood of protocol-correctness corner cases (sequence-number handling, pending-pay list invariants, simplex/duplex symmetry) versus straight implementation.

Example: instead of *"~half a day"*, say *"small — one new method on `messager.Messager` (~40 LOC), reuses existing dispatcher path, validation via the existing focused `sendCondPayWithErc20` e2e (~3 min on macOS); no proto change; low review risk."* That's actionable for sequencing; "half a day" is not.

## When in doubt

If a memory or recollection conflicts with what's currently in `AGENTS.md`, `docs/`, or the source — trust the repo. AgentPay's protocol surface is touchy (duplex simplex, sliding-window sequence numbers, DAL transaction boundaries); silent drift between agent recollection and code is how regressions land.
