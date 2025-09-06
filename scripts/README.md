# 脚本目录

此目录包含用于本地开发和测试的实用脚本。

## 可用脚本

### 🚀 start-local.sh / start-local.bat
本地测试启动脚本，读取 config.yaml 并安全地设置环境变量。

**Linux/macOS:**
```bash
./scripts/start-local.sh [选项]
```

**Windows:**
```cmd
scripts\start-local.bat [选项]
```

**选项:**
- `--migrate` - 在启动前运行数据库迁移
- `--env-only` - 只设置环境变量，不启动应用
- `--help` - 显示帮助信息

**示例:**
```bash
# 使用 config.yaml 值作为环境变量启动服务器
./scripts/start-local.sh

# 运行迁移后启动服务器
./scripts/start-local.sh --migrate

# 设置环境变量但不启动服务器
./scripts/start-local.sh --env-only
```

### 📝 load-env.sh
用于 .env 文件的环境变量加载器。

```bash
# 复制示例文件并编辑为您的值
cp .env.example.cn .env
nano .env

# 加载环境变量
source ./scripts/load-env.sh

# 现在运行您的应用程序
go run main.go
```

## 安全说明

- ⚠️ 这些脚本**仅用于本地开发**
- 🚫 **请勿**在生产环境中使用这些脚本
- ✅ 环境变量的优先级高于 config.yaml 值
- 🔒 敏感值在脚本输出中被掩码

## 故障排除

### 脚本权限被拒绝
```bash
chmod +x ./scripts/start-local.sh
chmod +x ./scripts/load-env.sh
```

### 找不到配置文件
确保您从项目根目录运行且 config.yaml 存在。

### 缺少 yq 工具
如果没有 yq，脚本将回退到基本的 grep/sed 解析。为获得最佳结果，请安装 yq：

```bash
# macOS
brew install yq

# Linux (Ubuntu/Debian)
sudo snap install yq
```

### 环境变量未设置
检查 config.yaml 是否包含必需的值：
- jwt.key
- postgres.password  
- admin.password
- admin.salt

## 与开发工作流程的集成

1. **初始设置:**
   ```bash
   cp .env.example.cn .env
   nano .env  # 编辑为您的值
   ```

2. **日常开发:**
   ```bash
   source ./scripts/load-env.sh
   go run main.go
   ```

3. **测试配置更改:**
   ```bash
   ./scripts/start-local.sh --env-only
   # 验证环境变量正确，然后:
   go run main.go
   ```

4. **数据库迁移:**
   ```bash
   ./scripts/start-local.sh --migrate
   ```

## 快速使用指南

### 方法 1: 使用启动脚本（推荐）

保留您现有的 `config.yaml` 包含开发值，使用启动脚本安全加载配置：

```bash
# 基本启动
./scripts/start-local.sh

# 包含数据库迁移
./scripts/start-local.sh --migrate

# 只设置环境变量（不启动服务器）
./scripts/start-local.sh --env-only
```

### 方法 2: 使用 .env 文件

```bash
# 1. 复制示例文件
cp .env.example .env

# 2. 编辑 .env 文件
nano .env

# 3. 加载环境变量
source ./scripts/load-env.sh

# 4. 启动应用程序
go run main.go
```

### 方法 3: 手动设置环境变量

```bash
export JWT_KEY="您的开发jwt密钥"
export POSTGRES_PASSWORD="您的本地postgres密码"
export ADMIN_PASSWORD="您的管理员密码"
export PASSWORD_SALT="您的密码盐值"
# ... 根据需要设置其他变量

go run main.go
```

## 脚本特性

### ✅ 功能亮点
- **智能配置解析**: 自动从 config.yaml 提取所有配置
- **安全输出**: 敏感信息在显示时自动掩码
- **完整验证**: 启动前验证所有必需的密钥
- **跨平台支持**: Linux、macOS 和 Windows 版本
- **中文界面**: 完全中文化的输出和帮助信息

### 🔒 安全特性
- **环境变量优先**: 生产环境密钥覆盖本地配置文件
- **敏感值掩码**: 所有密码和密钥在输出中被掩码
- **必需值验证**: 缺少关键配置时拒绝启动
- **安全默认值**: 非敏感配置使用安全的默认值

## 配置优先级

应用程序按以下顺序加载配置：

1. 加载 config.yaml（如果存在）
2. 使用环境变量覆盖（如果设置）
3. 应用内置默认值（仅非敏感值）
4. 验证必需的密钥存在
5. 启动应用程序或退出并显示错误

