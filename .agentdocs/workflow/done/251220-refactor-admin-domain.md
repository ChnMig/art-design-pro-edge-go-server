# 251220-refactor-admin-domain

## 背景
当前 `admin` 相关接口（`/api/v1/private/admin/*`）的业务逻辑分散在：
- `api/app/v1/private/admin/**`（handler 里包含大量业务判断与数据库操作）
- `db/pgdb/system/**`（同时承载 Gorm Model、查询/事务、以及部分业务规则）
- `common/menu/**`（包含 DB 写操作 `SaveRoleMenu`，不符合“common 纯共享”定位）

上游模板强调分层：`api` 仅负责参数绑定/DTO/错误映射，`domain` 承载核心业务规则与用例，基础设施（DB/缓存）作为实现细节通过接口注入。

## 目标
- 将所有 admin 业务逻辑下沉到 `domain`，`api` 仅保留：参数校验、鉴权/租户上下文读取、DTO 映射、错误响应。
- `db/pgdb/**` 仅保留数据访问（Repository）与 Model；禁止在 `common/**` 中直接写 DB。
- 保持接口路径、请求参数与响应结构尽量不变（避免前端联调成本）。
- 补齐单元测试与集成测试，并通过 `gofmt` + `go test ./...`。

## 设计约定（对齐模板）
- `domain/admin/<module>`：每个模块一个用例入口（Service/UseCase），定义领域错误（`errors.go`）。
- `api/...`：只做 DTO 与调用 domain，用 `ReturnDomainError` 统一映射错误（参考 `api/app/v1/open/health`）。
说明：当前以“先完成分层”为目标，domain 层暂时直接复用 `db/pgdb/system` 的查询/事务函数；后续若需要更强的可测试性，可再逐步补齐 repo 接口与依赖注入。

## 迁移范围
- admin/system：`tenant` / `department` / `role` / `user` / `menu` / `login` / `captcha` / `tenant_search`
- admin/platform：`menu` / `role`（平台侧特权能力）

## TODO
- [x] 搭建 domain 分层骨架：模块目录、统一错误映射模式。
- [x] 迁移 tenant/department/role：从 handler 下沉用例逻辑，API 仅调用 domain。
- [x] 迁移 user/login/captcha：拆分登录用例、验证码校验、登录日志与 JWT 签发边界。
- [x] 迁移 menu：角色菜单、租户范围、平台菜单定义与权限树逻辑全部下沉到 domain；清理 `common/menu` 的 DB 写操作。
- [x] 补齐测试：domain/menu 单测 + 关键接口的 gin 集成测试（不依赖外部 Postgres/Redis）。
- [x] 本地校验：执行 `gofmt -w` 与 `go test ./...`。
