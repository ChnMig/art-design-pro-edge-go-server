# 关于

本项目是 [art-design-pro-edge](https://github.com/ChnMig/art-design-pro-edge) 的后端服务。
配合前端可以做到开箱即用, 但是具体的业务功能需要自己开发.

## 项目特点

- 项目的95%代码由 `github copilot` 辅助编写

## TODO

- API层权限管制
- 接口文档
- 单元测试
- 持续的代码优化

## 部署配套服务

PostgreSQL 和 Redis 的 docker-compose 文件在 `docker` 目录下, 可以直接使用。

> 如果部署在云端, 务必修改有 TODO 标识的配置项, 防止密码泄露!!!

```bash
docker-compose -f docker/docker-compose.yml up -d
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

修改 `./config.yaml` 中的配置

## 执行数据库初始化

`go run main.go --migrate`

## start

`nohup ./server &`

## QA

TODO
