# 多租户系统架构说明

## 系统概述

本系统已升级为多租户架构，支持多个企业/组织在同一系统中独立运营，数据完全隔离。

## 核心表结构

### 1. SystemTenant (租户/企业表)

```go
type SystemTenant struct {
    gorm.Model
    Code        string     // 企业编号，唯一
    Name        string     // 企业名称
    Contact     string     // 联系人
    Phone       string     // 联系电话
    Email       string     // 邮箱
    Address     string     // 地址
    Status      uint       // 状态(1:启用 2:禁用)
    ExpiredAt   *time.Time // 过期时间
    MaxUsers    uint       // 最大用户数
}
```

### 2. SystemUser (用户表)

```go
type SystemUser struct {
    gorm.Model
    TenantID     uint   // 租户ID，数据隔离关键字段
    DepartmentID uint   
    RoleID       uint   
    Name         string // 姓名
    Username     string // 昵称
    Account      string // 登录账号，同租户内唯一
    Password     string // 密码（bcrypt加密）
    Phone        string
    Gender       uint   // 性别(1:男 2:女)
    Status       uint   // 状态(1:启用 2:禁用)
}
```

### 3. 其他业务表

所有业务表都增加了 `TenantID` 字段实现数据隔离：

- `SystemDepartment` - 部门表
- `SystemRole` - 角色表  
- `SystemUserLoginLog` - 登录日志表

## 多租户登录方式

登录时需要提供三个字段：

```json
{
  "tenant_code": "default",  // 企业编号
  "account": "admin",        // 登录账号
  "password": "your_password"
}
```

## JWT Token 结构

新的JWT Token包含租户信息：

```go
type MultiTenantClaims struct {
    UserID   uint   // 用户ID
    TenantID uint   // 租户ID  
    Account  string // 账号
    jwt.RegisteredClaims
}
```

## 登录响应结构

```json
{
  "access_token": "jwt_token_here",
  "tenant_info": {
    "tenant_id": 1,
    "tenant_code": "default",
    "tenant_name": "默认企业"
  },
  "user_info": {
    "user_id": 1,
    "name": "超级管理员",
    "username": "admin", 
    "account": "admin"
  }
}
```

## 数据隔离原理

1. **租户级隔离**: 所有业务数据都包含 `TenantID` 字段
2. **用户账号隔离**: `Account` 字段在同一租户内唯一，不同租户可重复
3. **JWT验证**: 每个请求都验证租户和用户身份
4. **API自动过滤**: 所有查询自动添加租户过滤条件

## 默认初始化数据

系统会创建默认租户和管理员：

```yaml
租户:
  code: "default"
  name: "默认企业"
  
管理员:
  tenant_id: 1
  name: "超级管理员"
  username: "admin"
  account: "admin"
  password: [配置文件AdminPassword]
```

## API 认证

所有需要认证的接口都需要：

1. 在请求头添加: `Authorization: Bearer <jwt_token>`
2. JWT中包含有效的租户和用户信息
3. 系统自动进行租户数据隔离

## 中间件

- `MultiTenantAuth()`: 多租户认证中间件
- `GetTenantID(c)`: 获取当前租户ID
- `GetCurrentUserID(c)`: 获取当前用户ID

## 租户管理

提供完整的租户管理API：

- 租户增删改查
- 用户数量限制
- 租户状态管理
- 到期时间控制

## 注意事项

1. **登录变更**: 前端登录页面需要添加企业编号输入框
2. **数据隔离**: 所有查询都会自动添加租户过滤
3. **账号唯一性**: 同租户内账号唯一，跨租户可重复
4. **权限控制**: 用户只能访问自己租户的数据
5. **扩展性**: 支持无限租户，每个租户独立配置