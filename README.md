# ZKDAA  部署与运行指南

本项目是一个基于零知识证明和 HTLC 的数据迁移原型系统，用于验证“数据完整性审计”和“链上解锁”绑定在同一套可验证流程中的可行性。系统中包含用户、DSPA、DSPB、Merkle 数据承诺、Groth16 证明、Verifier 合约和 DataMigration 合约等模块。

## 项目功能

系统主要完成以下流程：

1. 用户生成两侧秘密材料，分别交给 DSPA 和 DSPB。
2. DSPA 侧使用普通解锁电路证明 `H1 = MiMC(Pre_I, Sn_I)`。
3. DSPB 侧先对迁移数据构建 Merkle 树，得到数据承诺 `CIDF`。
4. DSPB 侧使用联合审计解锁电路证明 `H2 = MiMC(CIDF, Sn_II)`。
5. 链上 `DataMigration` 合约调用对应 Verifier 合约验证证明。
6. 两侧锁定、审计解锁和普通解锁完成后，输出运行指标 JSON。

核心电路在：

```text
circuit/unlock_circuit.go
circuit/audit_unlock_circuit.go
```

核心运行入口在：

```text
main.go
batch.go
latency.go
tps.go
utils.go
```

## 目录结构

```text
zk-htlc/
  actors/                    用户、DSPA、DSPB、Operator 等角色
  audit/                     审计流程辅助逻辑
  circuit/                   gnark 零知识电路
  cmd/                       独立命令入口
  config/                    本地链配置示例
  contracts/                 Go 合约绑定
  data/                      数据结构
  keys/                      旧版电路密钥示例
  merkle/                    Merkle 树实现
  std/hash/poseidon/         本地替换依赖
  tools/                     verifier 生成、调试等工具
  zkp/                       ZKP handler
  blockchain-contracts/      Hardhat 合约工程
  Datamigration.sol          DataMigration 合约源码备份
  DataMigration.abi          DataMigration ABI
```

## 环境要求

请提前安装：

- Go 1.24.5 或更高版本
- Node.js 18 或更高版本
- npm
- Git Bash、WSL、Linux 或 macOS 终端

`go.mod` 中要求：

```text
go 1.24.5
```

如果本机 Go 版本低于 1.24.5，建议先升级 Go。否则运行时 Go 会尝试自动下载 toolchain，网络代理异常时可能失败。

## 安装依赖

在项目根目录安装 Go 依赖：

```bash
go mod download
```

如果下载失败，可以切换 Go 代理后重试：

```bash
go env -w GOPROXY=https://proxy.golang.org,direct
go mod download
```

国内网络可使用：

```bash
go env -w GOPROXY=https://goproxy.cn,direct
go mod download
```

安装合约工程依赖：

```bash
cd blockchain-contracts
npm install
```

回到项目根目录：

```bash
cd ..
```

## 生成 ZKP Setup 和 Verifier 合约

端到端测试需要两个证明系统：

- `UnlockVerifier`：DSPA 普通解锁证明。
- `AuditVerifier_d<depth>`：DSPB 联合审计解锁证明。

`depth` 由文件大小和分块大小决定：

```text
chunk_count = file_size_bytes / chunk_size_bytes
depth = ceil(log2(chunk_count))
```

常用配置如下，默认分块大小为 `1024` bytes：

| 文件大小 | filesize | chunksize | depth |
|---|---:|---:|---:|
| 8MB | 8388608 | 1024 | 13 |
| 16MB | 16777216 | 1024 | 14 |
| 32MB | 33554432 | 1024 | 15 |
| 64MB | 67108864 | 1024 | 16 |
| 128MB | 134217728 | 1024 | 17 |

例如测试 8MB 文件，先生成 depth=13 的 setup 和 verifier：

```bash
go run tools/gen_verifier.go -depth 13
```

该命令会生成：

```text
build/unlock.pk
build/unlock.vk
build/audit_d13.pk
build/audit_d13.vk
blockchain-contracts/contracts/UnlockVerifier.sol
blockchain-contracts/contracts/AuditVerifier_d13.sol
```

如果切换到 16MB、32MB、64MB 或 128MB，需要重新用对应 depth 生成 verifier，并重新编译部署合约。

## 启动本地区块链

打开终端 1：

```bash
cd blockchain-contracts
npx hardhat node
```

保持该终端运行，不要关闭。

默认本地 RPC 地址为：

```text
http://127.0.0.1:8545
```

项目代码默认使用 Hardhat 第一个测试账户私钥：

```text
0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80
```

该私钥只允许用于本地测试链。

## 编译和部署合约

打开终端 2：

```bash
cd blockchain-contracts
npx hardhat compile
```

部署合约时传入当前测试使用的 Merkle depth。以 8MB、depth=13 为例：

```bash
DEPTH=13 npx hardhat run scripts/deploy.js --network localhost
```

Windows PowerShell：

```powershell
$env:DEPTH="13"
npx hardhat run scripts/deploy.js --network localhost
```

部署脚本会部署 Verifier 合约和 DataMigration 合约。记录输出中的两个 `DataMigration` 地址，后续分别作为：

```text
-addrA
-addrB
```

如果部署脚本输出的是一行两个地址，例如：

```text
DataMigration deployed at: 0xAddressA 0xAddressB
```

则第一个填入 `-addrA`，第二个填入 `-addrB`。

## 运行单次端到端测试

在项目根目录执行。

由于根目录可能包含多个独立 `main` 文件，推荐显式列出主程序所需文件：

```bash
go run main.go utils.go batch.go latency.go tps.go \
  -single \
  -filesize 8388608 \
  -chunksize 1024 \
  -addrA 0xYourDataMigrationA \
  -addrB 0xYourDataMigrationB \
  -output results/metrics_8MB.json
```

参数说明：

| 参数 | 说明 |
|---|---|
| `-single` | 单次端到端测试模式 |
| `-filesize` | 文件大小，单位 bytes |
| `-chunksize` | 分块大小，单位 bytes |
| `-addrA` | A 侧 DataMigration 合约地址 |
| `-addrB` | B 侧 DataMigration 合约地址 |
| `-output` | 指标 JSON 输出路径 |

不同文件大小示例：

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

注意：切换文件大小时，要先按对应 depth 重新执行：

```bash
go run tools/gen_verifier.go -depth <depth>
```

然后重新编译和部署合约。

## 运行 Latency 测试

Latency 模式用于测试不同节点规模下的交易确认和协议延迟。

```bash
go run main.go utils.go batch.go latency.go tps.go \
  -latency \
  -nodes 20 \
  -addrA 0xYourDataMigrationA \
  -addrB 0xYourDataMigrationB
```

常用节点规模：

```bash
-nodes 20
-nodes 30
-nodes 40
-nodes 50
-nodes 60
```

## 运行 TPS 测试

TPS 模式用于测试批量 Lock / Unlock 交易吞吐能力。

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

参数说明：

| 参数 | 说明 |
|---|---|
| `-nodes` | 模拟节点规模 |
| `-repeat` | 重复测试次数 |
| `-lock-ms` | Lock 交易发送间隔 |
| `-unlock-ms` | Unlock 交易发送间隔 |

## 常见运行顺序

推荐第一次部署时按这个顺序执行：

```bash
# 1. 安装 Go 依赖
go mod download

# 2. 安装合约依赖
cd blockchain-contracts
npm install
cd ..

# 3. 生成 8MB 测试所需 verifier
go run tools/gen_verifier.go -depth 13

# 4. 终端 1 启动本地链
cd blockchain-contracts
npx hardhat node
```

另开终端：

```bash
# 5. 编译部署合约
cd blockchain-contracts
npx hardhat compile
DEPTH=13 npx hardhat run scripts/deploy.js --network localhost
```

再回到项目根目录：

```bash
# 6. 运行端到端测试
go run main.go utils.go batch.go latency.go tps.go \
  -single \
  -filesize 8388608 \
  -chunksize 1024 \
  -addrA 0xYourDataMigrationA \
  -addrB 0xYourDataMigrationB \
  -output results/metrics_8MB.json
```

## 常见问题

### 1. `go run .` 报 `main redeclared`

根目录里可能有多个实验入口文件，它们都包含 `func main()`。运行主流程时请使用：

```bash
go run main.go utils.go batch.go latency.go tps.go ...
```

### 2. `go run main.go` 报函数未定义

`main.go` 依赖 `utils.go`、`batch.go`、`latency.go`、`tps.go` 中的函数。不要只运行单个 `main.go`。

### 3. 找不到 `build/unlock.pk` 或 `build/audit_d<depth>.pk`

说明还没有生成 setup。执行：

```bash
go run tools/gen_verifier.go -depth <depth>
```

### 4. 合约部署成功，但 Go 程序交易失败

检查：

1. Hardhat 节点是否仍在运行。
2. `-addrA`、`-addrB` 是否来自当前这次部署。
3. 测试文件大小对应的 depth 是否和部署的 `AuditVerifier_d<depth>` 一致。
4. 重新生成 verifier 后，是否重新执行 `npx hardhat compile` 和部署。

### 5. Go 自动下载 toolchain 失败

项目需要 Go 1.24.5 或更高版本。建议直接安装对应版本，或者配置可用代理：

```bash
go env -w GOPROXY=https://goproxy.cn,direct
go mod download
```

### 6. 大文件测试很慢

这是正常现象。文件越大，chunk 数越多，Merkle 深度越高，证明生成和数据准备时间都会增加。建议先用 8MB 跑通完整流程，再测试 16MB、32MB、64MB 和 128MB。

## 安全提醒

- README 中出现的私钥是 Hardhat 本地测试账户私钥，只能用于本地链。
- 不要把本地测试私钥用于公网链或任何真实资产。
- 每次重启 Hardhat node 后，之前部署的合约地址都会失效，需要重新部署并更新运行命令中的地址。
