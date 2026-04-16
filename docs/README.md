# AgentPay Backend Docs

This repo contains the off-chain backend implementation for AgentPay: OSP/service nodes, client-side node helpers, admin surfaces, routing, and persistence. These docs explain how that backend is built and how to run it.

Use these documents in this order:

1. [Backend Implementation](./backend-implementation.md) for the runtime model, package map, and protocol-to-code mapping.
2. [Backend Usage](./backend-usage.md) for build, test, startup, configuration, and operator workflows.
3. [Backend Troubleshooting](./backend-troubleshooting.md) for failure diagnosis, operational checks, and recovery steps.

These docs complement the companion `agentpay-architecture` documents. They do not repeat the full protocol or contract design; instead, they show how this repo realizes those ideas.

Useful repo references:

- [Project overview](../README.md)
- [Manual multi-OSP walkthrough](../test/manual/README.md)
- [OSP CLI reference](../tools/osp-cli/README.md)
- [Maintenance scripts](../tools/scripts/README.md)
- [Runtime config notes](../rtconfig/README.md)