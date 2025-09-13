#!/bin/bash

# 清理旧目录结构脚本
# 删除已经迁移到新结构的旧文件和目录

set -e

echo "🧹 开始清理旧的项目结构..."

# 确认用户操作
read -p "⚠️  此操作将删除旧的目录结构，是否继续？(y/N) " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "❌ 操作已取消"
    exit 1
fi

echo "📁 当前工作目录: $(pwd)"

# 备份重要文件
echo "💾 创建备份..."
mkdir -p backup/$(date +%Y%m%d_%H%M%S)
BACKUP_DIR="backup/$(date +%Y%m%d_%H%M%S)"

# 备份可能有价值的文件
if [ -f "main.go" ]; then
    cp main.go "$BACKUP_DIR/"
    echo "  ✅ 备份 main.go"
fi

if [ -f "README.md" ]; then
    cp README.md "$BACKUP_DIR/README_old.md"
    echo "  ✅ 备份 README.md"
fi

# 删除旧的目录结构
OLD_DIRS=(
    "api"
    "config"
    "common"
    "db"
    "util"
    "cron"
)

echo ""
echo "🗑️  删除旧目录..."
for dir in "${OLD_DIRS[@]}"; do
    if [ -d "$dir" ]; then
        echo "  🗂️  删除目录: $dir"
        rm -rf "$dir"
    fi
done

# 删除旧的主文件
OLD_FILES=(
    "main.go"
    "go.sum.bak"
    "*.bak"
)

echo ""
echo "🗑️  删除旧文件..."
for file_pattern in "${OLD_FILES[@]}"; do
    if ls $file_pattern 1> /dev/null 2>&1; then
        echo "  📄 删除文件: $file_pattern"
        rm -f $file_pattern
    fi
done

# 更新README
if [ -f "README_NEW.md" ]; then
    mv README_NEW.md README.md
    echo "  ✅ 更新 README.md"
fi

# 清理空目录
echo ""
echo "🧹 清理空目录..."
find . -type d -empty -delete 2>/dev/null || true

# 显示新的项目结构
echo ""
echo "📋 新的项目结构:"
tree -I 'backup|logs|*.log|api-server|node_modules|.git' -L 3 . || ls -la

echo ""
echo "✅ 清理完成！"
echo ""
echo "🎉 项目已成功重构为标准Go项目布局"
echo ""
echo "📚 新的项目结构："
echo "  cmd/          - 主应用程序"
echo "  internal/     - 私有代码"
echo "  configs/      - 配置文件"
echo "  scripts/      - 构建脚本"
echo "  docs/         - 文档"
echo ""
echo "🚀 使用方式："
echo "  构建: ./scripts/build.sh"
echo "  部署: ./scripts/deploy.sh [dev|prod]"
echo "  运行: ./api-server [--dev|--migrate]"
echo ""
echo "💾 备份文件保存在: $BACKUP_DIR"