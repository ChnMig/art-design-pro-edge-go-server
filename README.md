# 关于

> !!!项目还在开发中, 并未正式发布

本项目是 [art-design-pro-edge](https://github.com/ChnMig/art-design-pro-edge) 的后端服务。
配合前端可以做到开箱即用, 但是具体的业务功能需要自己开发.

# 技术栈

`Golang` `Gin` `Gorm` `PostgreSQL` `Redis`

# build

`CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server`

# dev

`go run main.go --dev`

# 初次启动

## 修改配置文件

修改 `config/config.go` 中的配置(尤其是 TODO 标识的配置项)

## 执行数据库初始化

`go run main.go --migrate`

# start

`nohup ./server &`

# QA

## 为什么不做数据库地址的动态配置

动态配置意味着需要将数据库地址存储在配置文件中，如果服务器被入侵, 会导致数据库地址泄露。所以我们不提供数据库地址的动态配置。
而将数据库地址写死在代码中, 可以避免这种情况的发生。
我们考虑过对文件进行加密, 但是这样带来的好处并不是很大, 考虑到用户是有开发经验的, 所以我们认为用户可以查看文档通过修改代码的方式进行配置。