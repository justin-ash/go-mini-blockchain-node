# Go Mini Blockchain Node

An educational, single-node blockchain prototype written in Go. It demonstrates block hashing, ordered transaction execution, account balances, nonce checks, and chain validation.

> **Status:** Learning project / prototype. This is not a production blockchain and must not be used to hold real funds.

## Features

- SHA-256 block hashes with deterministic JSON encoding
- Account balances and sequential transaction processing
- Rejects zero-value transfers, missing accounts, insufficient balances, invalid nonces, and integer overflow
- Blocks linked by previous hashes
- Chain validation that replays transactions from the configured initial balances
- Unit tests for valid transfers, invalid transfers, and tampering
- Standard library only

## Requirements

- Go 1.22 or newer

## Run

1. Create a GitHub repository named `go-mini-blockchain-node`.
2. Replace `YOUR_GITHUB_USERNAME` in `go.mod` with your GitHub username (or change the module path to your preferred module path).
3. From this directory, run:

```bash
go test ./...
go run ./cmd/node
```

## Example behavior

The demo starts Alice with 100 tokens, transfers 80 to Bob, and then attempts to transfer 30 to Carol. The second transfer is rejected because Alice only has 20 tokens left.

## Project structure

```text
.
├── blockchain/
│   ├── block.go
│   ├── chain.go
│   ├── state.go
│   ├── transaction.go
│   └── blockchain_test.go
├── cmd/node/main.go
├── go.mod
└── README.md
```

## Current limitations

- Single node; no peer-to-peer networking or consensus.
- No public-key signatures yet; callers are trusted to identify the sender.
- In-memory state only; no durable database.
- Timestamp is included in block data, but there is no proof-of-work or validator-based consensus.
- JSON is used as a deterministic encoding for this specific struct layout; production protocols should define and version a canonical serialization format.
- The chain can be recomputed by an attacker who can rewrite the whole chain. Hash linking alone does not provide consensus or immutability.

## Roadmap

1. Add transaction signatures and bind sender addresses to public keys.
2. Add an HTTP JSON-RPC/API layer.
3. Add persistent storage.
4. Add P2P transaction and block propagation.
5. Design and test a consensus/fork-choice mechanism.
6. Add benchmarks, fuzz tests, and security review.

## Interview discussion points

- Why transaction validation must happen before state mutation.
- Why each account needs a nonce.
- Why balances and nonces must be updated in transaction order.
- Why a valid block hash does not imply valid transactions.
- Why a single-node prototype is not yet a decentralized blockchain.
