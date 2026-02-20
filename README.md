# digital-identity

An open source, chain-agnostic identity ledger for individual digital identity management.

## Overview

**digital-identity** is a proposed free, open source platform that gives every individual subscriber their own personal digital identity chain. Each subscriber's chain is user-owned and append-only, and permanently records:

- **Public identity information** – verified attributes and credentials that the subscriber chooses to make public.
- **Identification events** – every act of identification (logins, verifications, authorizations, and other identity transactions) is recorded as an immutable event on the chain.

Smart contract integration is optional. The core protocol does **not** depend on any single blockchain.

## Core Principles

| Principle | Description |
|-----------|-------------|
| **Open source** | The entire codebase is open source and free to inspect, use, modify, and distribute under the [Apache 2.0 License](LICENSE). |
| **Free to use & adopt** | There are no licensing fees. Anyone can run a node, deploy the contracts, or build on top of the platform at no cost. |
| **User-owned** | Each identity chain is owned exclusively by the subscriber. No company, organization, or third party owns or controls a user's identity data. |
| **Decentralized** | Identity records live on a distributed blockchain, not in a centralized database controlled by a single entity. |
| **Strong root security** | Root actions require one long phrase plus 2FA/3FA, depending on risk level. |
| **Rich auditability** | High-fidelity, immutable event logs are the primary fraud defense mechanism. |
| **Chain agnostic** | The protocol works without smart contracts and supports optional chain adapters for anchoring. |

## How It Works

1. **Create chain** – A subscriber initializes a personal identity chain rooted in their own key.
2. **Protect root** – Sensitive actions require long phrase + 2FA (or 3FA for high-risk actions).
3. **Log events** – Every identification event is appended as an immutable, signed block link.
4. **Verify** – Any party can verify chain integrity cryptographically; optional blockchain anchors strengthen public auditability.

## Assurance Model

Identity assurance supports multiple trust levels:

- **Self asserted** – user-managed identity with cryptographic accountability.
- **Verified** – additional third-party or ecosystem attestations.
- **Government verified** – optional evidence such as passport/KYC to increase trust level.

Passport/KYC is optional and acts as trust elevation rather than sole proof of identity.

## Getting Started

Initial Go implementation is now bootstrapped.

### Current code layout

- `cmd/identityd` – API daemon entrypoint.
- `pkg/chain` – append-only identity chain types, append logic, and verification.
- `pkg/auth` – long phrase + 2FA/3FA policy checks.
- `pkg/assurance` – optional assurance evidence model (including passport/KYC types).
- `pkg/anchor` – chain-agnostic anchoring adapter interface.
- `pkg/store` – persistent JSON file store for chains.
- `pkg/api` – HTTP API for create, append, list, and verify operations.

### Run

1. Ensure Go 1.22+ is installed.
2. From project root, run:
	- `go run ./cmd/identityd`
	- `go test ./...`

The API starts on `:8080` by default.

### HTTP API (initial)

- `GET /healthz`
- `POST /chains` (create chain, auto-generate root key if `root_public_key` is omitted)
- `GET /chains`
- `GET /chains/{ownerID}`
- `POST /chains/{ownerID}/events` (requires long phrase + 2FA/3FA auth context and signer private key)
- `GET /chains/{ownerID}/verify`

Example flow:

1. Create a chain:
	- `curl -s -X POST http://localhost:8080/chains -H 'Content-Type: application/json' -d '{"owner_id":"user:alice"}'`
2. Append an event (using returned `root_private_key`):
	- `curl -s -X POST http://localhost:8080/chains/user:alice/events -H 'Content-Type: application/json' -d '{"event":{"type":"verification","risk":"normal","payload":{"method":"totp+email"}},"signer_private_key":"<ROOT_PRIVATE_KEY>","auth":{"long_phrase":true,"email_otp":true,"totp":true,"risk":"normal"}}'`
3. Verify integrity:
	- `curl -s http://localhost:8080/chains/user:alice/verify`

### Run with Docker

From project root:

- Build image: `docker build -t digital-identity:latest .`
- Run container: `docker run --rm -p 8080:8080 digital-identity:latest`

Or with compose:

- `docker compose up --build`

### Next implementation milestones

1. Persistent storage for chains and events.
2. API service for create/append/verify operations.
3. Checkpointing and Merkle proof support.
4. Production-grade key derivation, recovery, and hardware key integration.
5. Optional blockchain anchor adapters.

Contributions are welcome.

## License

Licensed under the [Apache License, Version 2.0](LICENSE).  
This software is provided **as-is**, free of charge, with no warranty. See the license for full terms.
