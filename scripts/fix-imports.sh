#!/bin/bash

# 修复剩余的导入路径问题

echo "🔧 修复剩余的导入路径问题..."

# 创建需要的目录和文件
echo "📁 创建缺失的目录..."
mkdir -p internal/app/api/auth
mkdir -p internal/pkg/path-tool

# 复制认证相关文件
if [ ! -f "internal/app/api/auth/multi_tenant.go" ] && [ -f "internal/app/auth/multi_tenant.go" ]; then
    cp internal/app/auth/multi_tenant.go internal/app/api/auth/
fi

# 复制路径工具
if [ ! -f "internal/pkg/path-tool/path.go" ] && [ -f "internal/pkg/util/path-tool/path.go" ]; then
    cp internal/pkg/util/path-tool/path.go internal/pkg/path-tool/
fi

# 批量替换导入路径
echo "🔄 批量更新导入路径..."

find internal -name "*.go" -type f -exec sed -i.bak '
    s|api-server/api/auth|api-server/internal/app/api/auth|g
    s|api-server/util/path-tool|api-server/internal/pkg/path-tool|g
    s|api-server/db/rdb|api-server/internal/repository/redis|g
    s|api-server/db/pgdb/system|api-server/internal/app/model/system|g
    s|api-server/config|api-server/internal/pkg/config|g
' {} \;

# 清理备份文件
find internal -name "*.bak" -delete

echo "✅ 导入路径修复完成"