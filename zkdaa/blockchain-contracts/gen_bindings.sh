#!/bin/bash

# 检查是否安装了 jq
if ! command -v jq &> /dev/null; then
    echo "❌ 错误: 未找到 jq 工具。请先运行: sudo apt install jq"
    exit 1
fi

# 定义目录
BUILD_DIR="build_go"
OUT_DIR="contracts_go"
TEMP_DIR="temp_ignored_contracts" # 临时存放被跳过文件的目录

# 1. 创建输出目录
rm -rf $BUILD_DIR $OUT_DIR
mkdir -p $BUILD_DIR
mkdir -p $OUT_DIR

# ==========================================
# 🙈 步骤 A: 临时移走 Batch 相关文件
# ==========================================
echo "🙈 正在临时移走 Batch 相关合约以跳过错误..."
mkdir -p $TEMP_DIR

# 移动所有 Batch 开头的文件
mv contracts/Batch* $TEMP_DIR/ 2>/dev/null || true
# 如果 IVerifier.sol 还在且报错，也移走（虽然日志说它已经不在了）
mv contracts/IVerifier.sol $TEMP_DIR/ 2>/dev/null || true

# 确保在脚本退出或出错时，一定把文件移回去 (Trap)
cleanup() {
    if [ -d "$TEMP_DIR" ]; then
        echo "♻️  正在恢复被跳过的合约文件..."
        mv $TEMP_DIR/* contracts/ 2>/dev/null || true
        rmdir $TEMP_DIR
    fi
}
trap cleanup EXIT

# ==========================================
# 🔨 步骤 B: 运行 Hardhat 编译
# ==========================================
echo "🔨 正在运行 Hardhat 编译 (仅编译剩余的核心合约)..."
npx hardhat compile

if [ $? -ne 0 ]; then
    echo "❌ Hardhat 编译失败！请检查 contracts/ 下剩余文件是否有误。"
    exit 1
fi

echo "✅ 编译完成"

# ==========================================
# 📦 步骤 C: 提取并生成 Go 绑定
# ==========================================

# 定义要处理的合约名称
CONTRACTS=("UnlockVerifier" "AuditVerifier" "DataMigration")

echo "🚀 开始提取 JSON 并生成 Go 绑定..."

for CONTRACT in "${CONTRACTS[@]}"; do
    echo "----------------------------------------"
    echo "📦 处理合约: $CONTRACT"

    # Hardhat 的构建产物路径
    ARTIFACT_PATH="artifacts/contracts/${CONTRACT}.sol/${CONTRACT}.json"

    if [ ! -f "$ARTIFACT_PATH" ]; then
        echo "❌ 找不到构建文件: $ARTIFACT_PATH"
        echo "   (可能编译未包含该文件，或文件名与合约名不一致)"
        continue
    fi

    # A. 提取 ABI
    cat "$ARTIFACT_PATH" | jq .abi > "$BUILD_DIR/${CONTRACT}.abi"

    # B. 提取 Bytecode
    cat "$ARTIFACT_PATH" | jq -r .bytecode > "$BUILD_DIR/${CONTRACT}.bin"

    # C. 生成 Go 代码
    OUT_FILE="$OUT_DIR/$(echo $CONTRACT | tr '[:upper:]' '[:lower:]').go"

    abigen --abi "$BUILD_DIR/${CONTRACT}.abi" \
           --bin "$BUILD_DIR/${CONTRACT}.bin" \
           --pkg contracts \
           --type $CONTRACT \
           --out "$OUT_FILE"

    echo "✅ 生成绑定: $OUT_FILE"
done

echo "----------------------------------------"
echo "🎉 流程结束！Batch 文件已归位。"
echo "👉 请运行: cp contracts_go/*.go ../contracts/ (根据你的路径调整)"