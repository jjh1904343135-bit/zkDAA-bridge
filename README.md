# ZKdaa 性能测试与可视化指南

## 概述

本指南介绍如何使用项目的性能测试系统，测试不同文件大小下的系统性能。

---

## 功能特性

### 1. 参数化测试
- **可配置文件大小**：支持从 KB 到 GB 级别的文件
- **可配置分块大小**：默认 128KB，可调整
- **JSON 数据导出**：自动导出详细性能指标
- **静默模式**：适合批量测试

### 2. 性能指标 (30+ 项)

#### ⏱️ 时间指标
- Setup 时间（解锁审计电路）
- 数据准备与传输时间
- 审计证明生成与验证时间
- 解锁证明生成时间
- 交易确认时间
- **端到端总时间**

#### Gas 指标
- Lock 操作 Gas (DSPA + DSPB)
- Unlock 操作 Gas (DSPA + DSPB)
- **总 Gas 消耗**

#### 数据指标
- 文件大小、分块数量、Merkle 树深度

#### 密码学指标
- 电路约束数量
- 证明大小

#### 审计指标
- 审计成功率、挑战次数

## 快速开始

### 步骤 1: 环境准备

#### 安装 Python 依赖
```bash

#### 确保 Go 环境正常
```bash
go version  # 应该显示 Go 1.18+

### 步骤 2: 启动区块链节点

在**终端 1**中运行：
```bash
cd blockchain-contracts
npx hardhat node
```

保持此终端运行，不要关闭！

### 步骤 3: 部署合约

在**终端 2**中运行：
```bash
cd blockchain-contracts
npx hardhat compile
npx hardhat run scripts/deploy.js --network localhost
```

**记录输出的合约地址**，并更新 `main.go` 中的地址（约第 186-187 行）。

---

## 运行性能测试

### 方法 1: 单次测试（手动）

测试特定文件大小：
```bash
# 测试 10KB 文件
go run main.go -size=10240 -chunk=1024 -output=results/test_10KB.json

# 测试 1MB 文件
go run main.go -size=1048576 -chunk=1024 -output=results/test_1MB.json

# 测试 128MB 文件（静默模式）
go run main.go -size=134217728 -chunk=1024 -output=results/test_128MB.json -silent
```

**命令行参数说明**：
- `-size`: 文件大小（bytes）
- `-chunk`: 分块大小（bytes），默认 1024
- `-output`: JSON 输出路径
- `-silent`: 静默模式，只输出关键信息

### 方法 2: 批量测试（自动）
```

#### Linux/Mac 系统
```bash
chmod +x scripts/run_benchmark.sh
./scripts/run_benchmark.sh
```

批量测试会自动运行以下文件大小：
- 1MB, 2MB, 4MB, 8MB, 16MB, 32MB, 64MB, 128MB

所有结果将保存到 `results/` 目录。


## 🔧 自定义测试配置

### 修改测试文件大小范围

编辑 scripts/run_benchmark.sh：

```bash
# 示例：测试更多文件大小
set sizes=512000 1048576 5242880 10485760 52428800 104857600
set labels=500KB 1MB 5MB 10MB 50MB 100MB
```

## ❓ 常见问题

### Q1: 测试失败，提示合约地址错误
**A**: 确保已经：
1. 运行 `npx hardhat node`
2. 部署合约并更新 `main.go` 中的地址

### Q2: Python 绘图失败
**A**:
```bash
# 检查依赖
pip list | grep matplotlib
pip list | grep numpy

# 重新安装
pip install --upgrade matplotlib numpy
```

### Q3: 大文件测试时间过长
**A**: 这是正常的！128MB 文件可能需要 30-60 分钟。建议：
- 先测试小文件（1-16MB）
- 使用 `-silent` 模式减少输出
- 后台运行：`nohup ./scripts/run_benchmark.sh &`

### Q4: 想要测试更小的文件
**A**: 修改批量测试脚本，添加 100KB、500KB 等大小：
```bash
set sizes=102400 512000 1048576 ...
set labels=100KB 500KB 1MB ...
```
