# digital-identity

An open source digital blockchain for individual digital identity management.

## Overview

**digital-identity** is a free, open source blockchain platform that gives every individual subscriber their own personal digital identity chain. Each subscriber's chain is a smart contract blockchain that permanently and transparently records:

- **Public identity information** – verified attributes and credentials that the subscriber chooses to make public.
- **Identification events** – every act of identification (logins, verifications, authorizations, and other identity transactions) is recorded as an immutable event on the chain.

## Core Principles

| Principle | Description |
|-----------|-------------|
| **Open source** | The entire codebase is open source and free to inspect, use, modify, and distribute under the [Apache 2.0 License](LICENSE). |
| **Free to use & adopt** | There are no licensing fees. Anyone can run a node, deploy the contracts, or build on top of the platform at no cost. |
| **User-owned** | Each identity chain is owned exclusively by the subscriber. No company, organization, or third party owns or controls a user's identity data. |
| **Decentralized** | Identity records live on a distributed blockchain, not in a centralized database controlled by a single entity. |
| **Smart contracts** | Identity and event data are stored and governed by transparent, auditable smart contracts. |

## How It Works

1. **Subscribe** – A new subscriber deploys (or is assigned) a personal smart contract that acts as their identity chain.
2. **Record identity** – Public identity attributes (name, verified credentials, public keys, etc.) are written to the smart contract.
3. **Log events** – Every identification event is appended as an immutable transaction on the chain, creating a complete, tamper-proof audit trail.
4. **Verify** – Any party can verify an identity or event by reading the public smart contract on the blockchain—no central authority required.

## Getting Started

> Implementation details, smart contract specifications, and deployment guides will be added as the project develops. Contributions are welcome!

1. Fork or clone this repository.
2. Review the [LICENSE](LICENSE) to understand your rights.
3. Open an issue or pull request to contribute.

## License

Licensed under the [Apache License, Version 2.0](LICENSE).  
This software is provided **as-is**, free of charge, with no warranty. See the license for full terms.
