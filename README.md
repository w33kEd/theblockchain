# TheBlockchain

A compact, educational Go implementation of a blockchain node and libraries. This project is intended for learning and experimentation — it implements a simple P2P network, in-memory blockchain store, a tiny VM, native NFT primitives, and an HTTP API for submitting transactions and querying chain data.

This README describes exactly how the repository is organized and how to run the included example (the `main` program), with concrete commands and examples based on the code.

**Highlights**

- Minimal in-memory blockchain with block validation and state handling.
- TCP-based P2P transport with simple RPC messages (transactions, blocks, status).
- Small stack-based VM that executes byte instructions embedded inside transactions.
- Native NFT primitives: `CollectionTx` and `MintTx`.
- HTTP API (Echo) exposing `POST /tx`, `GET /tx/:hash` and `GET /block/:hashOrHeight`.

## Quick facts (code-driven)

- Module path: `github.com/w33ked/theblockchain` (see `go.mod`).
- Go version recorded in `go.mod`: `go 1.25.5`.
- Default block time used by the validator loop: 5s (`network.defaultBlockTime`).
- Genesis block funds: the genesis transaction sets `Value = 10_000_000` for the coinbase.
- Storage is in-memory (`core.NewMemorystore()`), so state is not persisted across runs.

## Build & run (concrete)

Prerequisites: Go toolchain matching your environment (the code uses Go modules).

Build the project binary:

```bash
go build -o bin/theblockchain ./...
```

Run the example `main.go` which boots up three nodes and an HTTP API used by the example sender:

```bash
go run main.go
```

What `main.go` does (useful to know):

- Spins up three network servers bound to `:3000`, `:4000`, `:5000` (and a late-joining node on `:6000`).
- The first node in `main.go` is started with a generated private key and acts as a validator — it will create blocks on a ticker.
- The code also starts an HTTP API on `:9000` (when provided) and sends a sample transaction to `http://localhost:9000/tx` using a small helper `sendTransaction()`.

## HTTP API (implemented endpoints)

The API server (in `api/server.go`) exposes the following endpoints via Echo:

- `POST /tx` — accepts a gob-encoded `core.Transaction` in the request body. The server decodes the transaction and pushes it into the node's transaction channel.
- `GET /tx/:hash` — returns a JSON representation of the transaction identified by a hex-encoded 32-byte hash.
- `GET /block/:hashOrId` — accepts either a block height (decimal) or a hex-encoded hash; returns a JSON block representation.

Notes:

- The `POST /tx` endpoint expects Go-gob binary encoding (the example client in `main.go` demonstrates how to post such data).

## Networking model

- Transport: simple TCP transport (`network/tcp_transport.go`) with peers represented by `TCPPeer`.
- Messages: encoded messages include transaction, block, get-status, status, get-blocks and blocks (see `network/message.go`).
- Bootstrap: nodes connect to addresses listed in `ServerOpts.SeedNodes` and exchange status to decide syncing.
- Mempool: `network.TxPool` keeps an ordered pending pool and a global lookup map; pending transactions are used by validators when creating blocks.

## Consensus / validator behavior

- A server with a non-nil `PrivateKey` is considered a validator (`isValidator` flag).
- Validators run `validatorLoop()` and call `createNewBlock()` on a ticker (default 5s). The created block includes all pending transactions.
- Blocks are signed using the node's private key and validated before being appended (`core.Block.Verify()` and `core.Blockchain.AddBlock`).

## Transaction model

- `core.Transaction` fields of interest:
  - `Data []byte`: arbitrary data executed by the VM when included in a block.
  - `TxInner`: used for native NFT operations (`CollectionTx` and `MintTx`).
  - `To`, `From`: public keys used for native token transfers.
  - `Value uint64`: native token transfer amount.
  - `Signature`: ECDSA P-256 signature; signing uses `crypto.PrivateKey.Sign` and verification uses `Signature.Verify`.

- The blockchain handles:
  - Native transfers via `accountState.Transfer`.
  - NFT collection creation (`CollectionTx`) and minting (`MintTx`) stored in `collectionState` and `mintState` maps.
  - Arbitrary contract code execution via the VM when `tx.Data` is present.

## VM (tiny instruction set)

- The VM (`core/vm.go`) is a very small stack-based engine supporting instructions such as `PushInt`, `Add`, `Sub`, `Pack`, and `Store`.
- `InstrStore` writes into an in-memory contract state `core.State`.

## Persistence

- Current storage implementation (`core.MemoryStore`) is a no-op in-memory store. Blocks are stored in-memory on the `Blockchain` struct (`headers`, `blocks`, `blockStore`, `txStore`).

## Tests

There are unit tests in several packages (examples):

- `core/*_test.go` — block, blockchain, vm and state unit tests.
- `network/*_test.go` and `network/local_transport_test.go` — transport and messaging tests.
- `crypto/keypair_test.go` and `types/*_test.go`.

Run tests for the whole repository:

```bash
go test ./... -v
```

Run a single package's tests, e.g. core:

```bash
go test ./core -v
```

## Example: submitting a transaction

The included `main.go` contains `sendTransaction()` which demonstrates how to build, sign and gob-encode a `core.Transaction` and POST it to the API:

```go
tx := core.NewTransaction(nil)
tx.To = recipientPub
tx.Value = 777
tx.Sign(senderPriv)
var buf bytes.Buffer
tx.Encode(core.NewGobTxEncoder(&buf))
http.Post("http://localhost:9000/tx", "application/octet-stream", &buf)
```

This is the simplest way to exercise the API; the project intentionally uses gob encoding for the examples.

## Importing this project as a library

- Module path is in `go.mod`: `github.com/w33ked/theblockchain`.
- To import packages from this repo in your own module, add a dependency on the repo and then `go mod tidy`:

```bash
go get github.com/w33ked/theblockchain@main
go mod tidy
```

If you are developing both projects locally, use `replace` in your `go.mod`:

```text
replace github.com/w33ked/theblockchain => ../path/to/theblockchain
```

## Contributing

- Submit issues for bugs and feature requests.
- Open a pull request with tests and a short explanation of the change.
- Keep changes small and focused; add tests for new behavior.

## License

See the `LICENSE` file in the repository root for licensing details.
