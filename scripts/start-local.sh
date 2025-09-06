#!/bin/bash

# 本地开发测试启动脚本
# 此脚本从 config.yaml 读取配置并设置环境变量，用于安全的本地测试
# 使用方法: ./scripts/start-local.sh [--migrate|--help]

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
CONFIG_FILE="$PROJECT_ROOT/config.yaml"

# 彩色输出函数
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

# 显示帮助信息
show_help() {
    cat << EOF
本地开发测试启动脚本

此脚本从 config.yaml 读取配置并设置环境变量，用于安全的本地测试，然后启动应用程序。

使用方法:
    $0 [选项]

选项:
    --migrate       在启动前运行数据库迁移
    --env-only      只设置环境变量，不启动应用
    --help          显示此帮助信息

设置的环境变量:
    SERVER_PORT         服务器监听端口
    JWT_KEY            JWT 签名密钥
    JWT_EXPIRATION     JWT 令牌过期时间
    REDIS_HOST         Redis 服务器地址
    REDIS_PASSWORD     Redis 密码
    POSTGRES_HOST      PostgreSQL 主机
    POSTGRES_USER      PostgreSQL 用户名
    POSTGRES_PASSWORD  PostgreSQL 密码
    POSTGRES_DBNAME    PostgreSQL 数据库名
    POSTGRES_PORT      PostgreSQL 端口
    POSTGRES_SSLMODE   PostgreSQL SSL 模式
    POSTGRES_TIMEZONE  PostgreSQL 时区
    ADMIN_PASSWORD     管理员密码
    PASSWORD_SALT      密码哈希盐值

示例:
    ./scripts/start-local.sh              # 使用当前配置启动
    ./scripts/start-local.sh --migrate    # 运行迁移后启动
    ./scripts/start-local.sh --env-only   # 只设置环境变量

EOF
}

# 从 YAML 中提取值的函数（使用 yq 或回退到 grep/sed）
extract_yaml_value() {
    local yaml_path="$1"
    local config_file="$2"
    
    # 尝试使用 yq（最准确）
    if command -v yq &> /dev/null; then
        yq eval "$yaml_path" "$config_file" 2>/dev/null || echo ""
    else
        # 回退到使用 grep/sed 进行基本提取
        case "$yaml_path" in
            ".server.port")
                grep -A 1 "^server:" "$config_file" | grep "port:" | sed 's/.*port: *\([0-9]*\).*/\1/' | head -1
                ;;
            ".jwt.key")
                grep -A 2 "^jwt:" "$config_file" | grep "key:" | sed 's/.*key: *"\([^"]*\)".*/\1/' | head -1
                ;;
            ".jwt.expiration")
                grep -A 2 "^jwt:" "$config_file" | grep "expiration:" | sed 's/.*expiration: *"\([^"]*\)".*/\1/' | head -1
                ;;
            ".redis.host")
                grep -A 2 "^redis:" "$config_file" | grep "host:" | sed 's/.*host: *"\([^"]*\)".*/\1/' | head -1
                ;;
            ".redis.password")
                grep -A 2 "^redis:" "$config_file" | grep "password:" | sed 's/.*password: *"\([^"]*\)".*/\1/' | head -1
                ;;
            ".postgres.host")
                grep -A 6 "^postgres:" "$config_file" | grep "host:" | sed 's/.*host: *"\([^"]*\)".*/\1/' | head -1
                ;;
            ".postgres.user")
                grep -A 6 "^postgres:" "$config_file" | grep "user:" | sed 's/.*user: *"\([^"]*\)".*/\1/' | head -1
                ;;
            ".postgres.password")
                grep -A 6 "^postgres:" "$config_file" | grep "password:" | sed 's/.*password: *"\([^"]*\)".*/\1/' | head -1
                ;;
            ".postgres.dbname")
                grep -A 6 "^postgres:" "$config_file" | grep "dbname:" | sed 's/.*dbname: *"\([^"]*\)".*/\1/' | head -1
                ;;
            ".postgres.port")
                grep -A 6 "^postgres:" "$config_file" | grep "port:" | sed 's/.*port: *\([0-9]*\).*/\1/' | head -1
                ;;
            ".postgres.sslmode")
                grep -A 6 "^postgres:" "$config_file" | grep "sslmode:" | sed 's/.*sslmode: *"\([^"]*\)".*/\1/' | head -1
                ;;
            ".postgres.timezone")
                grep -A 6 "^postgres:" "$config_file" | grep "timezone:" | sed 's/.*timezone: *"\([^"]*\)".*/\1/' | head -1
                ;;
            ".admin.password")
                grep -A 2 "^admin:" "$config_file" | grep "password:" | sed 's/.*password: *"\([^"]*\)".*/\1/' | head -1
                ;;
            ".admin.salt")
                grep -A 2 "^admin:" "$config_file" | grep "salt:" | sed 's/.*salt: *"\([^"]*\)".*/\1/' | head -1
                ;;
            *)
                echo ""
                ;;
        esac
    fi
}

# 从 config.yaml 加载配置到环境变量
load_config_to_env() {
    print_info "正在从 $CONFIG_FILE 加载配置"
    
    if [[ ! -f "$CONFIG_FILE" ]]; then
        print_error "找不到配置文件: $CONFIG_FILE"
        print_error "请创建 config.yaml 文件或从项目根目录运行。"
        exit 1
    fi
    
    # 从 YAML 中提取值
    SERVER_PORT=$(extract_yaml_value ".server.port" "$CONFIG_FILE")
    JWT_KEY=$(extract_yaml_value ".jwt.key" "$CONFIG_FILE")
    JWT_EXPIRATION=$(extract_yaml_value ".jwt.expiration" "$CONFIG_FILE")
    REDIS_HOST=$(extract_yaml_value ".redis.host" "$CONFIG_FILE")
    REDIS_PASSWORD=$(extract_yaml_value ".redis.password" "$CONFIG_FILE")
    POSTGRES_HOST=$(extract_yaml_value ".postgres.host" "$CONFIG_FILE")
    POSTGRES_USER=$(extract_yaml_value ".postgres.user" "$CONFIG_FILE")
    POSTGRES_PASSWORD=$(extract_yaml_value ".postgres.password" "$CONFIG_FILE")
    POSTGRES_DBNAME=$(extract_yaml_value ".postgres.dbname" "$CONFIG_FILE")
    POSTGRES_PORT=$(extract_yaml_value ".postgres.port" "$CONFIG_FILE")
    POSTGRES_SSLMODE=$(extract_yaml_value ".postgres.sslmode" "$CONFIG_FILE")
    POSTGRES_TIMEZONE=$(extract_yaml_value ".postgres.timezone" "$CONFIG_FILE")
    ADMIN_PASSWORD=$(extract_yaml_value ".admin.password" "$CONFIG_FILE")
    PASSWORD_SALT=$(extract_yaml_value ".admin.salt" "$CONFIG_FILE")
    
    # 导出环境变量
    export SERVER_PORT="$SERVER_PORT"
    export JWT_KEY="$JWT_KEY"
    export JWT_EXPIRATION="$JWT_EXPIRATION"
    export REDIS_HOST="$REDIS_HOST"
    export REDIS_PASSWORD="$REDIS_PASSWORD"
    export POSTGRES_HOST="$POSTGRES_HOST"
    export POSTGRES_USER="$POSTGRES_USER"
    export POSTGRES_PASSWORD="$POSTGRES_PASSWORD"
    export POSTGRES_DBNAME="$POSTGRES_DBNAME"
    export POSTGRES_PORT="$POSTGRES_PORT"
    export POSTGRES_SSLMODE="$POSTGRES_SSLMODE"
    export POSTGRES_TIMEZONE="$POSTGRES_TIMEZONE"
    export ADMIN_PASSWORD="$ADMIN_PASSWORD"
    export PASSWORD_SALT="$PASSWORD_SALT"
    
    print_success "从 config.yaml 设置环境变量成功"
    
    # 验证关键变量
    local missing_vars=()
    [[ -z "$JWT_KEY" ]] && missing_vars+=("JWT_KEY")
    [[ -z "$POSTGRES_PASSWORD" ]] && missing_vars+=("POSTGRES_PASSWORD")
    [[ -z "$ADMIN_PASSWORD" ]] && missing_vars+=("ADMIN_PASSWORD")
    [[ -z "$PASSWORD_SALT" ]] && missing_vars+=("PASSWORD_SALT")
    
    if [[ ${#missing_vars[@]} -gt 0 ]]; then
        print_error "缺少必需的配置值:"
        for var in "${missing_vars[@]}"; do
            print_error "  - $var"
        done
        exit 1
    fi
    
    # 显示加载的配置（敏感值加掩码）
    print_info "已加载的配置:"
    echo "  SERVER_PORT=${SERVER_PORT:-8080}"
    echo "  JWT_KEY=${JWT_KEY:0:8}*** (已掩码)"
    echo "  JWT_EXPIRATION=${JWT_EXPIRATION:-12h}"
    echo "  REDIS_HOST=${REDIS_HOST:-127.0.0.1:6379}"
    echo "  REDIS_PASSWORD=${REDIS_PASSWORD:+***已掩码***}"
    echo "  POSTGRES_HOST=${POSTGRES_HOST:-127.0.0.1}"
    echo "  POSTGRES_USER=${POSTGRES_USER:-postgres}"
    echo "  POSTGRES_PASSWORD=***已掩码***"
    echo "  POSTGRES_DBNAME=${POSTGRES_DBNAME:-server}"
    echo "  POSTGRES_PORT=${POSTGRES_PORT:-5432}"
    echo "  POSTGRES_SSLMODE=${POSTGRES_SSLMODE:-disable}"
    echo "  POSTGRES_TIMEZONE=${POSTGRES_TIMEZONE:-Asia/Shanghai}"
    echo "  ADMIN_PASSWORD=***已掩码***"
    echo "  PASSWORD_SALT=${PASSWORD_SALT:0:8}*** (已掩码)"
}

# 运行数据库迁移
run_migration() {
    print_info "正在运行数据库迁移..."
    cd "$PROJECT_ROOT"
    
    if [[ ! -f "main.go" ]]; then
        print_error "找不到 main.go 文件。请从项目根目录运行。"
        exit 1
    fi
    
    go run main.go --migrate
    print_success "数据库迁移完成"
}

# 启动应用程序
start_application() {
    print_info "正在启动应用程序..."
    cd "$PROJECT_ROOT"
    
    if [[ ! -f "main.go" ]]; then
        print_error "找不到 main.go 文件。请从项目根目录运行。"
        exit 1
    fi
    
    print_success "正在使用环境变量启动服务器 (config.yaml 值将被忽略)"
    print_info "按 Ctrl+C 停止服务器"
    go run main.go
}

# 主函数
main() {
    local run_migrate=false
    local env_only=false
    
    # 解析命令行参数
    while [[ $# -gt 0 ]]; do
        case $1 in
            --migrate)
                run_migrate=true
                shift
                ;;
            --env-only)
                env_only=true
                shift
                ;;
            --help|-h)
                show_help
                exit 0
                ;;
            *)
                print_error "未知选项: $1"
                show_help
                exit 1
                ;;
        esac
    done
    
    print_info "艺术设计专业版 Go 服务器 - 本地测试脚本"
    print_warning "此脚本仅用于本地测试"
    print_warning "请勿在生产环境中使用此脚本"
    echo
    
    # 加载配置
    load_config_to_env
    echo
    
    # 如果需要，运行迁移
    if [[ "$run_migrate" == true ]]; then
        run_migration
        echo
    fi
    
    # 除非是仅环境变量模式，否则启动应用程序
    if [[ "$env_only" == false ]]; then
        start_application
    else
        print_success "环境变量已设置。现在可以运行您的应用程序:"
        print_info "cd $PROJECT_ROOT && go run main.go"
    fi
}

# 使用所有参数运行主函数
main "$@"