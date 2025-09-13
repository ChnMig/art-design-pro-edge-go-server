# Go 项目重构计划 - 标准化目录布局

## 当前项目结构分析

**当前问题**:
1. `main.go` 在根目录，不符合标准
2. `api/` 目录既包含路由又包含业务逻辑，混合了概念
3. `util/` 目录应该重新组织
4. 缺少标准的 `cmd/`, `internal/`, `pkg/` 目录
5. 配置和脚本文件分布不当

## 新的标准目录结构

```
art-design-pro-edge-go-server/
├── cmd/                          # 主应用程序
│   └── api-server/              # API服务器应用
│       └── main.go              # 移动main.go到这里
├── internal/                     # 私有应用和库代码
│   ├── app/                     # 应用程序代码
│   │   ├── api/                 # API处理器
│   │   │   ├── handler/         # HTTP处理器
│   │   │   │   ├── system/      # 系统相关处理器
│   │   │   │   │   ├── auth/    # 认证相关
│   │   │   │   │   ├── user/    # 用户管理
│   │   │   │   │   ├── tenant/  # 租户管理
│   │   │   │   │   ├── role/    # 角色管理
│   │   │   │   │   ├── menu/    # 菜单管理
│   │   │   │   │   └── department/ # 部门管理
│   │   │   ├── middleware/      # 中间件
│   │   │   ├── response/        # 响应处理
│   │   │   └── router/          # 路由配置
│   │   ├── service/             # 业务逻辑服务层
│   │   │   └── system/          # 系统服务
│   │   ├── repository/          # 数据访问层
│   │   │   ├── postgres/        # PostgreSQL仓库
│   │   │   └── redis/           # Redis仓库
│   │   └── model/               # 数据模型
│   │       └── system/          # 系统模型
│   ├── pkg/                     # 私有库代码
│   │   ├── auth/                # 认证包
│   │   ├── config/              # 配置管理
│   │   ├── database/            # 数据库连接
│   │   ├── logger/              # 日志工具
│   │   ├── cache/               # 缓存工具
│   │   ├── crypto/              # 加密工具
│   │   ├── id/                  # ID生成器
│   │   ├── validator/           # 验证器
│   │   └── scheduler/           # 定时任务
│   └── common/                  # 通用代码
│       ├── constants/           # 常量定义
│       ├── errors/              # 错误定义
│       └── utils/               # 工具函数
├── pkg/                         # 公共库代码
│   └── client/                  # 客户端SDK（如果需要）
├── configs/                     # 配置文件模板
│   ├── config.yaml.example     # 配置文件示例
│   └── docker/                  # Docker配置
├── scripts/                     # 脚本文件
│   ├── build.sh                 # 构建脚本
│   ├── deploy.sh                # 部署脚本
│   └── migrate.sh               # 数据迁移脚本
├── deployments/                 # 部署配置
│   ├── docker-compose.yml      # Docker compose
│   └── k8s/                     # Kubernetes配置
├── api/                         # API定义
│   ├── openapi/                 # OpenAPI/Swagger规范
│   └── protobuf/                # Protocol Buffer定义
├── docs/                        # 项目文档
│   ├── api/                     # API文档
│   ├── design/                  # 设计文档
│   └── examples/                # 示例代码
├── test/                        # 测试文件
│   ├── integration/             # 集成测试
│   └── fixtures/                # 测试数据
├── build/                       # 打包和CI配置
│   ├── ci/                      # CI配置
│   └── package/                 # 包配置
├── web/                         # Web资源
│   └── static/                  # 静态文件
├── tools/                       # 支持工具
└── vendor/                      # 依赖包（保持不变）
```

## 详细迁移计划

### 阶段1: 创建新目录结构
1. 创建标准目录
2. 保持原有结构以确保不中断

### 阶段2: 移动应用程序代码
1. `main.go` → `cmd/api-server/main.go`
2. `api/` → `internal/app/api/`
3. `config/` → `internal/pkg/config/`

### 阶段3: 重组数据库相关代码
1. `db/pgdb/` → `internal/repository/postgres/`
2. `db/rdb/` → `internal/repository/redis/`
3. 数据模型单独分离到 `internal/app/model/`

### 阶段4: 重组工具类
1. `util/` → `internal/pkg/` 下的相应包
2. `common/` → `internal/common/`
3. `cron/` → `internal/pkg/scheduler/`

### 阶段5: 更新导入路径
1. 更新所有 Go 文件中的导入路径
2. 更新配置文件引用
3. 测试确保所有功能正常

### 阶段6: 清理和优化
1. 删除旧目录
2. 更新文档
3. 更新构建脚本

## 导入路径映射

| 旧路径 | 新路径 |
|--------|--------|
| `api-server/api` | `api-server/internal/app/api` |
| `api-server/config` | `api-server/internal/pkg/config` |
| `api-server/db/pgdb` | `api-server/internal/repository/postgres` |
| `api-server/db/rdb` | `api-server/internal/repository/redis` |
| `api-server/util/log` | `api-server/internal/pkg/logger` |
| `api-server/util/encryption` | `api-server/internal/pkg/crypto` |
| `api-server/util/id` | `api-server/internal/pkg/id` |
| `api-server/util/authentication` | `api-server/internal/pkg/auth` |
| `api-server/common` | `api-server/internal/common` |
| `api-server/cron` | `api-server/internal/pkg/scheduler` |

## 预期收益

1. **符合Go标准**: 遵循golang-standards/project-layout
2. **清晰的职责分离**: API、服务、仓库、模型分层明确
3. **更好的可维护性**: 代码组织更加合理
4. **团队协作**: 标准化的结构便于团队理解
5. **扩展性**: 为未来的微服务拆分做准备

## 注意事项

1. **保持向后兼容**: 在迁移过程中确保系统正常运行
2. **分阶段执行**: 避免一次性大规模重构
3. **充分测试**: 每个阶段完成后进行功能测试
4. **文档同步**: 及时更新相关文档