#!/bin/bash

# 商品上架 Agent 启动脚本

echo "🚀 启动商品上架 Agent..."
echo ""

# 检查 .env 文件
if [ ! -f ".env" ]; then
    echo "⚠️  警告: .env 文件不存在"
    echo "📝 请复制 .env.example 并配置你的 API 密钥:"
    echo "   cp .env.example .env"
    echo "   vim .env  # 编辑配置"
    echo ""
    exit 1
fi

# 检查 Python 版本
PYTHON_CMD="python3"
if ! command -v python3 &> /dev/null; then
    PYTHON_CMD="python"
fi

PYTHON_VERSION=$($PYTHON_CMD --version 2>&1 | awk '{print $2}')
echo "🐍 Python 版本: $PYTHON_VERSION"

# 激活虚拟环境（如果存在）
if [ -d "venv" ]; then
    echo "📦 激活虚拟环境..."
    source venv/bin/activate
elif [ -d ".venv" ]; then
    echo "📦 激活虚拟环境..."
    source .venv/bin/activate
fi

# 检查依赖
echo "📚 检查依赖..."
pip list | grep -q langgraph || {
    echo "⚠️  依赖未安装，正在安装..."
    pip install -e .
}

echo ""
echo "✨ 启动 Agent..."
echo ""

# 运行 Agent
$PYTHON_CMD -m product_listing_agent.main
