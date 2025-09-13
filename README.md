# Art Design Pro Edge Go Server

基于Go语言的艺术设计管理系统后端服务，采用标准的Go项目布局结构。

## 🏗️ 项目结构

本项目遵循 [golang-standards/project-layout](https://github.com/golang-standards/project-layout) 标准布局：

```
art-design-pro-edge-go-server/
├── cmd/                    # 主应用程序
│   └── api-server/        # API服务器主程序
│       └── main.go
├── internal/              # 私有应用程序和库代码
│   ├── app/              # 应用程序代码
│   │   ├── api/          # API相关代码
│   │   │   ├── router/   # 路由配置
│   │   │   ├── middleware/ # 中间件
│   │   │   └── response/ # 响应处理
│   │   ├── app/          # 业务处理器
│   │   │   └── system/   # 系统模块
│   │   └── model/        # 数据模型
│   │       └── system/   # 系统模型
│   ├── pkg/              # 私有库代码
│   │   ├── config/       # 配置管理
│   │   ├── crypto/       # 加密工具
│   │   ├── logger/       # 日志管理
│   │   └── scheduler/    # 定时任务
│   ├── repository/       # 数据访问层
│   │   ├── postgres/     # PostgreSQL数据库
│   │   └── redis/        # Redis缓存
│   └── common/           # 公共工具
├── configs/              # 配置文件
├── scripts/              # 构建和部署脚本
├── docs/                 # 文档
└── logs/                 # 日志文件
```

## 🚀 快速开始

### 环境要求

- Go 1.19+
- PostgreSQL 12+
- Redis 6+

### 安装和运行

1. **克隆项目**
   ```bash
   git clone <repository-url>
   cd art-design-pro-edge-go-server
   ```

2. **安装依赖**
   ```bash
   go mod download
   ```

3. **配置环境变量**
   ```bash
   # 复制配置文件
   cp configs/config.example.yaml configs/config.yaml

   # 设置环境变量
   export JWT_KEY="your-secret-key"
   export PGSQL_DSN="postgres://user:password@localhost/dbname?sslmode=disable"
   export REDIS_HOST="localhost:6379"
   export ADMIN_PASSWORD="admin-password"
   export PWD_SALT="password-salt"
   ```

4. **构建和运行**
   ```bash
   # 使用构建脚本
   ./scripts/build.sh

   # 数据库迁移
   ./api-server --migrate

   # 启动服务
   ./api-server
   ```

### 开发模式

```bash
# 开发模式运行
./api-server --dev

# 或使用go run
cd cmd/api-server && go run . --dev
```

## 📋 API文档

### 多租户系统

系统支持多租户架构，每个API请求都需要在JWT token中包含租户信息。

### 主要功能模块

1. **认证系统**
   - 用户登录/登出
   - JWT token管理
   - 多租户支持

2. **用户管理**
   - 用户CRUD操作
   - 角色权限管理
   - 部门管理

3. **菜单权限**
   - 菜单管理
   - 权限控制
   - 角色菜单关联

4. **租户管理** (超级管理员功能)
   - 租户CRUD操作
   - 租户状态管理
   - 用户数量限制

### API端点

详细的API文档请参考 [API_DOCUMENTATION.md](./API_DOCUMENTATION.md)

## 🛠️ 开发工具

### 构建脚本

- `scripts/build.sh` - 构建项目
- `scripts/deploy.sh` - 部署脚本
- `scripts/update-imports.sh` - 更新导入路径（重构时使用）

### 配置文件

- `configs/config.example.yaml` - 配置示例
- `configs/config.dev.yaml` - 开发环境配置
- `configs/config.prod.yaml` - 生产环境配置

## 🏛️ 架构设计

### 分层架构

1. **展示层 (Presentation Layer)**
   - HTTP路由和中间件
   - 请求/响应处理
   - 参数验证

2. **业务逻辑层 (Business Logic Layer)**
   - 业务规则处理
   - 服务编排
   - 权限控制

3. **数据访问层 (Data Access Layer)**
   - 数据库操作
   - 缓存管理
   - 事务处理

### 多租户支持

- 基于JWT token的租户识别
- 数据级别的租户隔离
- 租户配置和限制管理

### 安全特性

- bcrypt密码加密
- JWT token认证
- 角色权限控制
- 请求限流
- SQL注入防护

## 🧪 测试

```bash
# 运行测试
go test ./...

# 运行基准测试
go test -bench=. ./...

# 代码覆盖率
go test -cover ./...
```

## 📊 监控和日志

### 日志管理

使用zap日志库，支持：
- 结构化日志
- 日志级别控制
- 文件轮转
- 性能优化

### 监控指标

- 响应时间
- 错误率
- 并发用户数
- 数据库连接数

## 🚀 部署

### 开发环境部署

```bash
./scripts/deploy.sh dev
```

### 生产环境部署

```bash
./scripts/deploy.sh prod
```

### Docker部署 (待实现)

```bash
# 构建镜像
docker build -t api-server .

# 运行容器
docker run -p 8080:8080 -e JWT_KEY="your-key" api-server
```

## 🤝 贡献指南

1. Fork项目
2. 创建特性分支
3. 提交更改
4. 推送到分支
5. 创建Pull Request

## 📄 许可证

本项目采用MIT许可证 - 查看 [LICENSE](LICENSE) 文件了解详情

## 🆘 支持

如果您遇到问题或有疑问，请：

1. 查看 [文档](./docs/)
2. 搜索 [Issues](../../issues)
3. 创建新的 Issue
4. 联系维护者

## 🗺️ 路线图

- [ ] 完善API文档
- [ ] 添加单元测试
- [ ] Docker支持
- [ ] CI/CD管道
- [ ] 性能优化
- [ ] 监控面板
- [ ] 多语言支持