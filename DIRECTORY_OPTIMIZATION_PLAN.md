# 目录结构优化方案

## 🎯 优化目标
- 消除命名歧义和重复
- 提高目录结构的语义清晰度
- 遵循Go社区最佳实践
- 建立清晰的分层架构

## 📊 当前问题分析

### 高优先级问题
1. **`internal/app/app/`** - 重复的"app"命名
2. **重复的repository目录** - rdb/redis, pgdb/postgres
3. **语义不清晰** - common目录职责模糊

### 中优先级问题
4. **configs vs config** - 单复数不统一
5. **缺少service层** - 业务逻辑层缺失

## 🏗️ 推荐的新结构

```
art-design-pro-edge-go-server/
├── cmd/
│   └── api-server/
├── internal/
│   ├── transport/              # 传输层
│   │   └── http/              # HTTP传输
│   │       ├── router/        # 路由配置
│   │       ├── middleware/    # 中间件
│   │       └── response/      # 响应处理
│   ├── handler/               # HTTP处理器 (原 app/app)
│   │   └── system/           # 系统处理器
│   ├── service/              # 业务服务层 (新增)
│   │   └── system/          # 系统服务
│   ├── domain/              # 领域模型 (原 model)
│   │   └── system/         # 系统模型
│   ├── repository/         # 数据访问层
│   │   ├── postgres/      # PostgreSQL (清理重复)
│   │   └── redis/         # Redis (清理重复)
│   ├── pkg/              # 内部工具包
│   └── shared/           # 共享代码 (原 common)
├── config/               # 配置文件 (统一单数)
├── script/              # 脚本文件 (统一单数)
└── doc/                # 文档 (统一单数)
```

## 🔄 分阶段实施计划

### 阶段1: 修复明显问题 (高优先级)
- [ ] `internal/app/app/` → `internal/handler/`
- [ ] `internal/common/` → `internal/shared/`
- [ ] 清理重复的repository目录

### 阶段2: 结构优化 (中优先级)
- [ ] `internal/app/` 重新组织为transport层
- [ ] `configs/` → `config/`
- [ ] 添加 `internal/service/` 业务服务层

### 阶段3: 语义优化 (低优先级)
- [ ] `internal/app/model/` → `internal/domain/`
- [ ] 统一其他目录命名

## 💡 架构优势

新结构遵循Clean Architecture原则：

1. **Transport层** - 处理HTTP、gRPC等传输协议
2. **Handler层** - 处理请求和响应
3. **Service层** - 业务逻辑和用例
4. **Domain层** - 业务实体和规则
5. **Repository层** - 数据访问

## 🔍 实施建议

**推荐方案**: 渐进式重构
- 先修复最明显的问题
- 每次改动后验证编译
- 更新导入路径和文档
- 保持功能完整性

**是否继续**: 用户确认后开始实施第一阶段改进