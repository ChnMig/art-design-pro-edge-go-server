# 文档索引

## 产品文档
- 暂无

## 后端文档
- `backend/config.md`：记录配置加载架构、Viper/环境变量约定及安全校验要求，修改 `config` 包或 `config.yaml` 前必读。
- `backend/admin-layering.md`：记录 admin 模块的 api/domain 分层规范、错误映射与测试约定，修改 `admin` 接口前必读。

## 当前任务文档
- 暂无

## 历史任务文档
- `workflow/done/251109-sync-go-template.md`：记录已完成的 go-template(http-services) 框架对齐过程及结论。
- `workflow/done/251220-sync-go-template.md`：记录本次模板依赖/日志/TLS/health 对齐与验证结果。
- `workflow/done/251220-refactor-admin-domain.md`：将 admin 模块按上游模板分层（api/domain）拆分，并补齐菜单相关测试。

## 全局重要记忆
- Go 代码在提交前必须通过 `gofmt` 与 `go test ./...`，禁止跳过必要的本地校验。
- 日志、配置、对外交互等注释应使用中文描述，首次出现的英文专业名词补充中文解释。
