#!/bin/bash

# 更新Go导入路径脚本
# 将旧的导入路径替换为新的标准化路径

echo "正在更新Go导入路径..."

# 定义导入路径映射
declare -A imports=(
    ["api-server/api"]="api-server/internal/app/api/router"
    ["api-server/config"]="api-server/internal/pkg/config"
    ["api-server/db/pgdb"]="api-server/internal/repository/postgres"
    ["api-server/db/pgdb/system"]="api-server/internal/app/model/system"
    ["api-server/db/rdb"]="api-server/internal/repository/redis"
    ["api-server/util/log"]="api-server/internal/pkg/logger"
    ["api-server/util/encryption"]="api-server/internal/pkg/crypto"
    ["api-server/util/id"]="api-server/internal/pkg/id"
    ["api-server/util/authentication"]="api-server/internal/pkg/auth"
    ["api-server/common"]="api-server/internal/common"
    ["api-server/cron"]="api-server/internal/pkg/scheduler"
    ["api-server/api/middleware"]="api-server/internal/app/api/middleware"
    ["api-server/api/response"]="api-server/internal/app/api/response"
    ["api-server/api/auth"]="api-server/internal/app/api/auth"
)

# 获取所有Go文件
go_files=$(find . -name "*.go" -not -path "./vendor/*" -not -path "./.git/*")

# 对每个导入路径进行替换
for old_import in "${!imports[@]}"; do
    new_import="${imports[$old_import]}"
    echo "替换 $old_import -> $new_import"

    # 使用sed替换导入路径
    for file in $go_files; do
        if [[ "$OSTYPE" == "darwin"* ]]; then
            # macOS
            sed -i '' "s|\"$old_import\"|\"$new_import\"|g" "$file"
        else
            # Linux
            sed -i "s|\"$old_import\"|\"$new_import\"|g" "$file"
        fi
    done
done

echo "导入路径更新完成！"