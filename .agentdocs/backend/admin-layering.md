# Admin 分层规范（api/domain）

## 使用场景
- 修改或新增 `admin` 相关接口（`api/app/v1/private/admin/**`）时必读，确保与上游模板一致的分层与可测试性。

## 目录与职责
- `api/app/v1/private/admin/**`：仅做参数绑定、鉴权/租户上下文读取、DTO 转换、调用 domain、错误映射（`ReturnDomainError`）。
- `domain/admin/<module>`：承载业务用例与规则（校验、权限判断、范围限制等），并定义领域错误（`errors.go`）。
- `db/pgdb/system/**`：仅保留 Gorm Model 与数据库读写/事务封装（Repository 的具体实现）。
- `common/**`：仅放纯共享结构与纯函数（DTO、树构建等），禁止直接写 DB。

## 错误处理约定
- domain 返回领域错误（`errors.Is` 可识别的哨兵错误），api 通过 `ReturnDomainError(c, err, fallback)` 映射为 `api/response` 统一响应。
- api 侧记录错误日志使用 `api-server/util/log` 的 `log.WithRequest(c)`（基于请求上下文的 logger）。

## 测试约定
- domain：对纯逻辑（如菜单树解析、范围校验）补齐单元测试。
- api：补齐 gin 集成测试，优先覆盖参数校验/鉴权分支，避免依赖外部 Postgres/Redis；如需数据层能力，应通过可注入的 repo/fake 来替代真实连接。

## 说明（过渡策略）
当前实现为了先完成分层，对部分模块采取“domain 直接复用 `db/pgdb/system` 查询/事务函数”的过渡方案；后续如需要更强的可测试性与边界收敛，再逐步补齐 `repo.go` 接口并在 domain 侧通过依赖注入使用。

