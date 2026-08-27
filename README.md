# 关于

本项目是 [art-design-pro-edge](https://github.com/ChnMig/art-design-pro-edge) 的后端服务。
配合前端可以做到开箱即用, 但是具体的业务功能需要自己开发.

## 项目特点

- 项目的95%代码由 `github copilot` 辅助编写

## TODO

- 业务接口细粒度权限扩展
- 接口文档
- 单元测试
- 持续的代码优化

## 部署配套服务

PostgreSQL 和 Redis 的 Docker Compose（容器编排）文件在 `docker` 目录下。密码必须通过环境变量提供，端口默认仅绑定 `127.0.0.1`。

```bash
cp .env.example .env
# 修改 .env 后启动；应用本身不会自动读取 .env
docker compose --env-file .env -f docker/docker-compose.yml up -d
```

## 技术栈

`Golang` `Gin` `Gorm` `PostgreSQL` `Redis`

## build

`CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server`

## dev

`go run main.go --dev`

## 初次启动

### 修改配置文件

> 务必修改配置文件, 尤其是密码相关

复制配置示例并修改：

```bash
cp config.yaml.example config.yaml
```

然后修改 `./config.yaml` 中的配置（JWT/Redis/Postgres/Admin 等敏感项务必替换）。

也可使用 `HTTP_SERVICES_<分组>_<字段>` 环境变量覆盖配置。配置在启动时读取并校验，修改后必须重启服务；不支持运行中热加载。

## 执行数据库初始化

`go run main.go --migrate`

迁移会幂等补齐默认租户、菜单、根角色、根部门和根用户。新增租户时，接口会在单个事务中同时创建默认部门、租户管理员角色、首个管理员账号及初始菜单范围。

升级既有数据库时，迁移不会按旧角色 ID 猜测租户管理员。若曾人工建立非默认租户，请由平台超级管理员通过平台角色接口把对应角色的 `is_tenant_admin` 明确设为 `true`。

## start

`nohup ./server &`

## QA

TODO
