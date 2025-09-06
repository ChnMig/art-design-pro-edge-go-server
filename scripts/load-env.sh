#!/bin/bash

# 环境变量加载脚本
# 此脚本从 .env 文件加载环境变量
# 使用方法: source ./scripts/load-env.sh

set -e

# 输出颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # 无颜色

# 脚本目录和项目根目录
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
ENV_FILE="$PROJECT_ROOT/.env"

print_info() {
    echo -e "${BLUE}[信息]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[成功]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[警告]${NC} $1"
}

print_error() {
    echo -e "${RED}[错误]${NC} $1"
}

# 检查 .env 文件是否存在
if [[ ! -f "$ENV_FILE" ]]; then
    print_error "在 $ENV_FILE 找不到 .env 文件"
    print_info "要创建一个："
    print_info "1. 复制示例文件: cp .env.example .env"
    print_info "2. 编辑 .env 文件，填入您的值"
    print_info "3. 运行: source ./scripts/load-env.sh"
    return 1 2>/dev/null || exit 1
fi

# 加载 .env 文件
print_info "正在从 .env 文件加载环境变量..."

# 导出 .env 文件中每一行（不是注释或空行）
while IFS= read -r line || [[ -n "$line" ]]; do
    # 跳过空行和注释
    [[ -z "$line" ]] && continue
    [[ "$line" =~ ^[[:space:]]*# ]] && continue
    
    # 导出变量
    if [[ "$line" =~ ^[[:space:]]*([A-Za-z_][A-Za-z0-9_]*)=(.*)$ ]]; then
        var_name="${BASH_REMATCH[1]}"
        var_value="${BASH_REMATCH[2]}"
        export "$var_name"="$var_value"
        print_info "已加载: $var_name"
    fi
done < "$ENV_FILE"

print_success "从 .env 文件加载环境变量成功"
print_info "现在可以运行: go run main.go"

# 显示关键变量状态（已掩码）
echo
print_info "关键环境变量状态:"
[[ -n "$JWT_KEY" ]] && echo "  JWT_KEY: ✓ 已设置 (${JWT_KEY:0:8}***)" || echo "  JWT_KEY: ✗ 缺失"
[[ -n "$POSTGRES_PASSWORD" ]] && echo "  POSTGRES_PASSWORD: ✓ 已设置 (***已掩码***)" || echo "  POSTGRES_PASSWORD: ✗ 缺失"
[[ -n "$ADMIN_PASSWORD" ]] && echo "  ADMIN_PASSWORD: ✓ 已设置 (***已掩码***)" || echo "  ADMIN_PASSWORD: ✗ 缺失"
[[ -n "$PASSWORD_SALT" ]] && echo "  PASSWORD_SALT: ✓ 已设置 (${PASSWORD_SALT:0:8}***)" || echo "  PASSWORD_SALT: ✗ 缺失"