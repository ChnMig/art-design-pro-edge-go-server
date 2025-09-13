#!/bin/bash

# 部署脚本 for art-design-pro-edge-go-server
# 支持开发环境和生产环境部署

set -e

ENVIRONMENT=${1:-dev}
CONFIG_FILE=""

echo "🚀 开始部署 art-design-pro-edge-go-server (环境: $ENVIRONMENT)..."

# 检查环境参数
case $ENVIRONMENT in
    dev|development)
        echo "📝 开发环境部署"
        CONFIG_FILE="configs/config.dev.yaml"
        ;;
    prod|production)
        echo "📝 生产环境部署"
        CONFIG_FILE="configs/config.prod.yaml"
        ;;
    *)
        echo "❌ 无效的环境参数: $ENVIRONMENT"
        echo "使用方式: $0 [dev|prod]"
        exit 1
        ;;
esac

# 检查配置文件
if [ ! -f "$CONFIG_FILE" ]; then
    echo "⚠️  配置文件不存在: $CONFIG_FILE"
    echo "将使用环境变量配置"
fi

# 构建项目
echo "🔨 构建项目..."
./scripts/build.sh

# 检查必要的环境变量
echo "🔍 检查环境变量..."
REQUIRED_VARS=("JWT_KEY" "PGSQL_DSN" "REDIS_HOST" "ADMIN_PASSWORD" "PWD_SALT")
MISSING_VARS=()

for var in "${REQUIRED_VARS[@]}"; do
    if [ -z "${!var}" ]; then
        MISSING_VARS+=("$var")
    fi
done

if [ ${#MISSING_VARS[@]} -ne 0 ]; then
    echo "❌ 缺少必要的环境变量:"
    printf "  %s\n" "${MISSING_VARS[@]}"
    echo ""
    echo "请设置以下环境变量："
    echo "export JWT_KEY=\"your-secret-key\""
    echo "export PGSQL_DSN=\"postgres://user:password@localhost/dbname?sslmode=disable\""
    echo "export REDIS_HOST=\"localhost:6379\""
    echo "export ADMIN_PASSWORD=\"admin-password\""
    echo "export PWD_SALT=\"password-salt\""
    exit 1
fi

echo "✅ 环境变量检查通过"

# 数据库迁移（仅在生产环境中）
if [ "$ENVIRONMENT" = "prod" ] || [ "$ENVIRONMENT" = "production" ]; then
    echo "🗄️  执行数据库迁移..."
    ./api-server --migrate
    if [ $? -eq 0 ]; then
        echo "✅ 数据库迁移成功"
    else
        echo "❌ 数据库迁移失败"
        exit 1
    fi
fi

# 创建systemd服务文件（仅在Linux生产环境）
if [ "$ENVIRONMENT" = "prod" ] && [ "$(uname)" = "Linux" ]; then
    echo "📋 创建systemd服务文件..."

    cat > /tmp/api-server.service << EOF
[Unit]
Description=Art Design Pro Edge Go Server
After=network.target postgresql.service redis.service

[Service]
Type=simple
User=www-data
WorkingDirectory=$(pwd)
ExecStart=$(pwd)/api-server
Restart=always
RestartSec=5
StandardOutput=journal
StandardError=journal
SyslogIdentifier=api-server

Environment=JWT_KEY=$JWT_KEY
Environment=PGSQL_DSN=$PGSQL_DSN
Environment=REDIS_HOST=$REDIS_HOST
Environment=ADMIN_PASSWORD=$ADMIN_PASSWORD
Environment=PWD_SALT=$PWD_SALT

[Install]
WantedBy=multi-user.target
EOF

    echo "💾 systemd服务文件已生成: /tmp/api-server.service"
    echo "请手动复制到 /etc/systemd/system/ 并启用服务:"
    echo "  sudo cp /tmp/api-server.service /etc/systemd/system/"
    echo "  sudo systemctl daemon-reload"
    echo "  sudo systemctl enable api-server"
    echo "  sudo systemctl start api-server"
fi

echo ""
echo "🎉 部署准备完成！"
echo ""
echo "🚀 启动服务:"
if [ "$ENVIRONMENT" = "dev" ]; then
    echo "  开发模式: ./api-server --dev"
else
    echo "  生产模式: ./api-server"
fi
echo ""
echo "📊 查看日志:"
echo "  tail -f logs/api-server.log"
echo ""
echo "🔍 健康检查:"
echo "  curl http://localhost:8080/api/v1/health"