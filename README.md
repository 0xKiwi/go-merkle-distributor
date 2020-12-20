# SAFE1-Claims

This repo is a golang implementation of generating the information needed for a Merkle Distributor contract.

Note: This implementation uses a sparse merkle trie instead of the standard merkle tree. Normal proof verification will not work.

Instructions to run:
1. If using new data, create a `.env` with RPC_URL as your RPC endpoint.
2. `go run .`