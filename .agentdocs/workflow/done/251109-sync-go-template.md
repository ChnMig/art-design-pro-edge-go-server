# 251109-sync-go-template

## 背景
用户希望将项目对齐到上游 https://github.com/ChnMig/go-template/tree/main/http-services 的最新框架（Viper 配置、优雅关闭、请求中间件等），再把现有业务模块合入，以保持长期同步能力。

## 方案要点
1. 引入上游主干代码（config/log/main/middleware 等）并保持与其目录/依赖一致。
2. 保留并迁移现有业务逻辑（API、DB、Cron 等），确保新骨架只替换基础设施层。
3. 更新配置文件、文档与测试指令，确保后续二次开发遵循统一规范。

## TODO
- [x] 对齐基础设施：复制/改造 config、util/log、util/path-tool、api/middleware、main.go 等骨架文件。
- [x] 迁移业务入口：在新骨架下重新挂载原有 router、cron、db、config 校验逻辑并增加 CLI 选项。
- [x] 更新文档与验证：同步 `config.yaml` 示例说明，补充 AGENTS.md 要求并执行 gofmt/go test 验证。
