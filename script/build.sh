#!/bin/bash

# Build script for art-design-pro-edge-go-server
# Following golang-standards/project-layout structure

set -e

echo "🔨 开始构建 art-design-pro-edge-go-server..."

# 检查Go环境
if ! command -v go &> /dev/null; then
    echo "❌ Go环境未安装"
    exit 1
fi

echo "✅ Go版本: $(go version)"

# 进入项目根目录
cd "$(dirname "$0")/.."

echo "📁 当前工作目录: $(pwd)"

# 清理之前的构建
echo "🧹 清理之前的构建..."
rm -f ./api-server
rm -f ./cmd/api-server/api-server

# 下载依赖
echo "📦 下载Go模块依赖..."
go mod download
go mod tidy

# 构建主程序
echo "🔨 构建主程序..."
cd cmd/api-server
go build -o ../../api-server .
cd ../..

# 验证构建结果
if [ -f "./api-server" ]; then
    echo "✅ 构建成功！"
    echo "📊 构建信息:"
    ls -la ./api-server
    echo ""
    echo "🚀 运行方式:"
    echo "  开发模式: ./api-server --dev"
    echo "  数据库迁移: ./api-server --migrate"
    echo "  生产模式: ./api-server"
else
    echo "❌ 构建失败！"
    exit 1
fi

echo "🎉 构建完成！"