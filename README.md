# ZKDAA / ZK-HTLC Deployment and Running Guide

This project is a zero-knowledge HTLC prototype for verifiable data migration. It combines data integrity auditing, Merkle commitments, Groth16 proofs, verifier contracts, and on-chain HTLC-style locking/unlocking into one end-to-end workflow.

The full GitHub project includes the `blockchain-contracts/` Hardhat project. This README assumes that directory is present.

## What This Project Does

The system implements the following flow:

1. A user generates secret materials and distributes them to DSPA and DSPB.
2. DSPA uses the normal unlock circuit to prove `H1 = MiMC(Pre_I, Sn_I)`.
3. DSPB builds a Merkle tree for the migrated data and obtains the data commitment `CIDF`.
4. DSPB uses the combined audit-unlock circuit to prove `H2 = MiMC(CIDF, Sn_II)`.
5. The on-chain `DataMigration` contract verifies the submitted proofs through the corresponding verifier contracts.
6. After both sides complete lock, audit-unlock, and unlock operations, the program exports runtime metrics as JSON.

Core circuits:

```text
circuit/unlock_circuit.go
circuit/audit_unlock_circuit.go
```

Main runtime files:

```text
main.go
batch.go
latency.go
tps.go
utils.go
```

## Project Structure

```text
zk-htlc/
  actors/                    User, DSPA, DSPB, Operator, and related roles
  audit/                     Audit workflow helpers
  circuit/                   gnark zero-knowledge circuits
  cmd/                       Standalone command entry points
  config/                    Local chain configuration examples
  contracts/                 Go contract bindings
  data/                      Data structures
  keys/                      Legacy circuit key examples
  merkle/                    Merkle tree implementation
  std/hash/poseidon/         Local replacement dependency
  tools/                     Verifier generation and debugging tools
  zkp/                       ZKP handlers
  blockchain-contracts/      Hardhat contract project
  Datamigration.sol          DataMigration contract source backup
  DataMigration.abi          DataMigration ABI
```

## Requirements

Install the following tools first:

- Go 1.24.5 or later
- Node.js 18 or later
- npm
- Git Bash, WSL, Linux, or macOS terminal

The project requires:

```text
go 1.24.5
```

If your local Go version is lower than 1.24.5, upgrade Go first. Otherwise Go may try to download the required toolchain automatically, and that can fail if your network proxy is unavailable.

## Install Dependencies

From the project root, install Go dependencies:

```bash
go mod download
```

If dependency download fails, configure a Go proxy and retry:

```bash
go env -w GOPROXY=https://proxy.golang.org,direct
go mod download
```

For mainland China networks, you can also use:

```bash
go env -w GOPROXY=https://goproxy.cn,direct
go mod download
```

Install contract project dependencies:

```bash
cd blockchain-contracts
npm install
```

Return to the project root:

```bash
cd ..
```

## Generate ZKP Setup and Verifier Contracts

The end-to-end test needs two proof systems:

- `UnlockVerifier`: normal DSPA unlock proof.
- `AuditVerifier_d<depth>`: combined DSPB audit-unlock proof.

The Merkle depth is determined by file size and chunk size:

```text
chunk_count = file_size_bytes / chunk_size_bytes
depth = ceil(log2(chunk_count))
```

Common configurations with the default chunk size of `1024` bytes:

| File size | filesize | chunksize | depth |
|---|---:|---:|---:|
| 8MB | 8388608 | 1024 | 13 |
| 16MB | 16777216 | 1024 | 14 |
| 32MB | 33554432 | 1024 | 15 |
| 64MB | 67108864 | 1024 | 16 |
| 128MB | 134217728 | 1024 | 17 |

For example, to test an 8MB file, generate setup files and verifier contracts for `depth=13`:

```bash
go run tools/gen_verifier.go -depth 13
```

This generates:

```text
build/unlock.pk
build/unlock.vk
build/audit_d13.pk
build/audit_d13.vk
blockchain-contracts/contracts/UnlockVerifier.sol
blockchain-contracts/contracts/AuditVerifier_d13.sol
```

When switching to another file size, regenerate the verifier with the corresponding depth, then recompile and redeploy the contracts.

## Start the Local Blockchain

Open terminal 1:

```bash
cd blockchain-contracts
npx hardhat node
```

Keep this terminal running.

The default local RPC endpoint is:

```text
http://127.0.0.1:8545
```

The Go code uses the first default Hardhat account private key for local testing:

```text
0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80
```

Use this private key only on the local Hardhat test chain.

## Compile and Deploy Contracts

Open terminal 2:

```bash
cd blockchain-contracts
npx hardhat compile
```

Deploy contracts with the Merkle depth used by the current test. For an 8MB file, use `depth=13`:

```bash
DEPTH=13 npx hardhat run scripts/deploy.js --network localhost
```

Windows PowerShell:

```powershell
$env:DEPTH="13"
npx hardhat run scripts/deploy.js --network localhost
```

The deploy script should deploy the verifier contracts and the `DataMigration` contracts. Record the two deployed `DataMigration` addresses and pass them to the Go program as:

```text
-addrA
-addrB
```

If the deploy output prints two addresses on one line, for example:

```text
DataMigration deployed at: 0xAddressA 0xAddressB
```

use the first address for `-addrA` and the second address for `-addrB`.

## Run a Single End-to-End Test

Run the command from the project root.

Because the root directory may contain multiple standalone files with `func main()`, use an explicit file list for the main workflow:

```bash
mkdir -p results

go run main.go utils.go batch.go latency.go tps.go \
  -single \
  -filesize 8388608 \
  -chunksize 1024 \
  -addrA 0xYourDataMigrationA \
  -addrB 0xYourDataMigrationB \
  -output results/metrics_8MB.json
```

Arguments:

| Argument | Description |
|---|---|
| `-single` | Run one end-to-end test |
| `-filesize` | File size in bytes |
| `-chunksize` | Chunk size in bytes |
| `-addrA` | DataMigration contract address for side A |
| `-addrB` | DataMigration contract address for side B |
| `-output` | JSON metrics output path |

Examples for other file sizes:

```bash
# 16MB
go run main.go utils.go batch.go latency.go tps.go \
  -single \
  -filesize 16777216 \
  -chunksize 1024 \
  -addrA 0xYourDataMigrationA \
  -addrB 0xYourDataMigrationB \
  -output results/metrics_16MB.json
```

```bash
# 128MB
go run main.go utils.go batch.go latency.go tps.go \
  -single \
  -filesize 134217728 \
  -chunksize 1024 \
  -addrA 0xYourDataMigrationA \
  -addrB 0xYourDataMigrationB \
  -output results/metrics_128MB.json
```

Before changing file size, regenerate setup and verifier contracts for the corresponding depth:

```bash
go run tools/gen_verifier.go -depth <depth>
```

Then compile and deploy contracts again.

## Run Latency Tests

Latency mode measures protocol and transaction confirmation latency under different node-scale settings.

```bash
go run main.go utils.go batch.go latency.go tps.go \
  -latency \
  -nodes 20 \
  -addrA 0xYourDataMigrationA \
  -addrB 0xYourDataMigrationB
```

Common node-scale values:

```bash
-nodes 20
-nodes 30
-nodes 40
-nodes 50
-nodes 60
```

## Run TPS Tests

TPS mode measures Lock / Unlock transaction throughput.

```bash
go run main.go utils.go batch.go latency.go tps.go \
  -tps \
  -nodes 20 \
  -repeat 5 \
  -lock-ms 12 \
  -unlock-ms 18 \
  -addrA 0xYourDataMigrationA \
  -addrB 0xYourDataMigrationB
```

Arguments:

| Argument | Description |
|---|---|
| `-nodes` | Simulated node scale |
| `-repeat` | Number of repeated runs |
| `-lock-ms` | Interval between Lock transactions |
| `-unlock-ms` | Interval between Unlock transactions |

## Recommended First Run

For a clean local run, follow this order:

```bash
# 1. Install Go dependencies
go mod download

# 2. Install contract dependencies
cd blockchain-contracts
npm install
cd ..

# 3. Generate verifier contracts for the 8MB test
go run tools/gen_verifier.go -depth 13

# 4. Start the local chain in terminal 1
cd blockchain-contracts
npx hardhat node
```

Open another terminal:

```bash
# 5. Compile and deploy contracts
cd blockchain-contracts
npx hardhat compile
DEPTH=13 npx hardhat run scripts/deploy.js --network localhost
```

Return to the project root:

```bash
# 6. Run the end-to-end test
go run main.go utils.go batch.go latency.go tps.go \
  -single \
  -filesize 8388608 \
  -chunksize 1024 \
  -addrA 0xYourDataMigrationA \
  -addrB 0xYourDataMigrationB \
  -output results/metrics_8MB.json
```

## Common Issues

### 1. `go run .` reports `main redeclared`

The root directory may contain multiple standalone entry files, and each one has its own `func main()`. Use this command for the main workflow:

```bash
go run main.go utils.go batch.go latency.go tps.go ...
```

### 2. `go run main.go` reports undefined functions

`main.go` depends on functions defined in `utils.go`, `batch.go`, `latency.go`, and `tps.go`. Do not run `main.go` alone.

### 3. `build/unlock.pk` or `build/audit_d<depth>.pk` is missing

Setup files have not been generated yet. Run:

```bash
go run tools/gen_verifier.go -depth <depth>
```

### 4. Contracts deploy successfully, but Go transactions fail

Check the following:

1. The Hardhat node is still running.
2. `-addrA` and `-addrB` come from the current deployment.
3. The test file size matches the deployed `AuditVerifier_d<depth>`.
4. After regenerating verifier contracts, you recompiled and redeployed the contracts.

### 5. Go toolchain download fails

The project requires Go 1.24.5 or later. The simplest fix is to install the required Go version directly. You can also configure a working proxy:

```bash
go env -w GOPROXY=https://goproxy.cn,direct
go mod download
```

### 6. Large-file tests are slow

This is expected. Larger files produce more chunks, deeper Merkle paths, and higher data preparation and proof generation costs. Run the 8MB case first, then test 16MB, 32MB, 64MB, and 128MB.

## Security Notes

- The private key shown in this README is the default Hardhat local test account.
- Do not use it on public chains or with real assets.
- After restarting the Hardhat node, previous contract addresses become invalid. Redeploy contracts and update `-addrA` / `-addrB`.
