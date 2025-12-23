# 配置架构

## 使用场景
- 修改 `config` 包或 `config.yaml` 时必须阅读，确保与 Viper 加载流程一致。
- 引入新的运行参数或敏感配置时，遵循统一的命名与校验规则。

## 当前约定
1. 通过 `viper` 读取 `config.yaml`，并支持 `HTTP_SERVICES_<SECTION>_<KEY>` 形式的环境变量覆盖。
2. `config.LoadConfig` 负责设置默认值、读取文件、应用配置，并暴露给其它模块的全局变量（JWT、Redis、Postgres、限流、租户、PID 文件、TLS/ACME 等）。
3. 变更配置后可调用 `config.WatchConfig` 自动热加载，回调中仅允许执行轻量级逻辑（日志、缓存配置值等），重型操作需另行协调。
4. `config.CheckConfig` 必须在 `zap` logger 初始化后执行，确保关键项（JWT、数据库、Redis、初始用户等）已就绪，缺失会触发 `Fatal`。
5. TLS 启用约定：`server.enable_acme` 与 `server.enable_tls` 互斥；启用 ACME 时必须配置 `server.acme_domain`；启用本地证书模式时必须配置 `server.tls_cert_file/server.tls_key_file`。
6. PID 文件约定：`server.pid_file` 默认 `api-server.pid`，支持相对路径（相对 `config.AbsPath`）；服务启动后写入当前进程 PID，服务退出时自动删除；如不需要可配置为空字符串禁用。
7. 数据库 DSN、Redis 地址等敏感数据仅存于本地 `config.yaml`（已被 `.gitignore` 忽略），仓库仅提交 `config.yaml.example` 作为字段示例。
8. 文档及代码中面向用户的注释保持中文表达，必要的英文术语首次出现时附带中文说明。
