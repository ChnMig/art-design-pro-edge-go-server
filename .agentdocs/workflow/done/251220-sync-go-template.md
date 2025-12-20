# 251220-sync-go-template

## 背景
上游后端模板仓库 https://github.com/ChnMig/go-template（目录：`http-services/`）已更新，需要将本项目的基础设施层与其最新实现保持一致，同时保留现有业务逻辑（DB/Cron/多租户等）。

## 目标
- 对齐上游近期变更：Go 版本与依赖升级、日志能力增强、可选 TLS（ACME/本地证书）能力与相关配置项。
- 保持现有 API 行为与业务模块不被破坏（响应格式、多租户鉴权、数据库等）。
- 补齐必要的单元测试/集成测试，并通过 `gofmt` 与 `go test ./...` 校验。

## 同步范围（原则）
- 以“基础设施层”为主：`config/`、`main.go`、`util/log`、`api/router.go`、通用中间件与公共工具。
- 业务差异保持：多租户 JWT、数据库/Redis 配置与业务路由不回退。
- 新增能力默认关闭：TLS/ACME 等仅在配置启用时生效，避免影响现网。

## TODO
- [x] 生成对比清单：定位上游新增/变更点并映射到本项目目录结构。
- [x] 同步依赖：对齐 Go 版本与关键依赖版本，更新 `go.mod/go.sum`（必要时更新 `vendor/`）。
- [x] 同步日志：引入 `NewZapWriter/WithRequest/BoundParamsKey` 等能力，并在 gin 初始化时重定向框架日志到 zap。
- [x] 同步 TLS：补齐 `config` 配置项与 `util/acme`、`util/tlsfile` 实现，并在 `main.go` 中按配置启用。
- [x] 同步健康检查：对齐上游 health 模块结构（DTO/领域层/错误映射）并补充测试。
- [x] 更新示例配置与文档：补齐 `config.yaml.example` 新字段与使用说明。
- [x] 本地校验：执行 `gofmt -w` 与 `go test ./...`。
