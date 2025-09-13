# Art Design Pro Edge Go Server - API 文档

## 多租户系统概述

本系统采用多租户架构，支持多个企业/组织在同一系统中独立运营，数据完全隔离。

### 基础信息
- **Base URL**: `/api/v1`
- **认证方式**: Bearer Token (JWT)
- **内容类型**: `application/json`
- **编码**: UTF-8

### 多租户特性
- ✅ **完善的租户管理**: 租户增删改查已完整实现并注册路由
- ✅ **严格权限控制**: 租户管理仅超级管理员可访问
- ✅ **数据完全隔离**: 租户级数据隔离，自动过滤
- ✅ **安全认证**: JWT Token包含租户和用户信息
- ✅ **账号独立性**: 用户账号在同一租户内唯一

---

## 1. 认证相关 API

### 1.1 获取验证码

**接口**: `GET /admin/system/user/login/captcha`

**说明**: 获取登录验证码

**请求参数**:
```json
{
  "width": 150,   // 验证码宽度，必填
  "height": 40    // 验证码高度，必填
}
```

**响应示例**:
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "id": "captcha_id_here",
    "image": "data:image/png;base64,iVBORw0KGgoAAAANSU..."
  }
}
```

### 1.2 用户登录 (多租户)

**接口**: `POST /admin/system/user/login`

**说明**: 多租户用户登录，需要企业编号

**请求参数**:
```json
{
  "tenant_code": "default",         // 企业编号，必填
  "account": "admin",              // 登录账号，必填
  "password": "your_password",     // 密码，必填
  "captcha": "123456",             // 验证码，必填
  "captcha_id": "captcha_id_here"  // 验证码ID，必填
}
```

**响应示例**:
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
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
}
```

### 1.3 获取用户信息

**接口**: `GET /admin/system/user/info`

**说明**: 获取当前登录用户信息

**请求头**:
```
Authorization: Bearer <access_token>
```

**响应示例**:
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "user_id": 1,
    "tenant_id": 1,
    "name": "超级管理员",
    "username": "admin",
    "account": "admin",
    "phone": "13800138000",
    "gender": 1,
    "status": 1,
    "role_id": 1,
    "department_id": 1
  }
}
```

### 1.4 更新用户信息

**接口**: `PUT /admin/system/user/info`

**说明**: 更新当前用户个人信息

**请求头**:
```
Authorization: Bearer <access_token>
```

**请求参数**:
```json
{
  "name": "新的姓名",
  "phone": "13800138001",
  "gender": 1
}
```

### 1.5 查询登录日志

**接口**: `GET /admin/system/login/log`

**说明**: 查询登录日志列表

**请求头**:
```
Authorization: Bearer <access_token>
```

**请求参数**:
```json
{
  "ip": "127.0.0.1",        // 可选，IP地址模糊查询
  "username": "admin",      // 可选，用户名模糊查询
  "page": 1,               // 页码，默认1
  "page_size": 10          // 每页数量，默认10
}
```

**响应示例**:
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "total": 100,
    "list": [
      {
        "id": 1,
        "tenant_code": "default",
        "user_name": "admin",
        "ip": "127.0.0.1",
        "login_status": "success",
        "created_at": "2024-01-01T10:00:00Z"
      }
    ]
  }
}
```

---

## 2. 用户管理 API

### 2.1 查询用户列表

**接口**: `GET /admin/system/user`

**说明**: 查询用户列表，自动按当前租户过滤

**请求头**:
```
Authorization: Bearer <access_token>
```

**请求参数**:
```json
{
  "username": "admin",      // 可选，昵称模糊查询
  "name": "管理员",         // 可选，姓名模糊查询
  "phone": "138",          // 可选，手机号模糊查询
  "role_id": 1,            // 可选，角色ID
  "department_id": 1,      // 可选，部门ID
  "page": 1,               // 页码，默认1
  "page_size": 10          // 每页数量，默认10
}
```

**响应示例**:
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "total": 50,
    "list": [
      {
        "User": {
          "id": 1,
          "tenant_id": 1,
          "name": "超级管理员",
          "username": "admin",
          "account": "admin",
          "phone": "13800138000",
          "gender": 1,
          "status": 1,
          "role_id": 1,
          "department_id": 1,
          "created_at": "2024-01-01T10:00:00Z"
        },
        "role_name": "超级管理员",
        "role_desc": "系统超级管理员",
        "department_name": "技术部"
      }
    ]
  }
}
```

### 2.2 查询用户缓存列表

**接口**: `GET /admin/system/user/cache`

**说明**: 从缓存中查询用户列表（性能更好）

**请求头**:
```
Authorization: Bearer <access_token>
```

**请求参数**:
```json
{
  "id": 1,                 // 可选，指定用户ID获取单个用户
  "username": "admin",     // 可选，昵称模糊查询
  "name": "管理员",        // 可选，姓名模糊查询
  "page": 1,               // 页码，默认1
  "page_size": 10          // 每页数量，默认10
}
```

### 2.3 添加用户

**接口**: `POST /admin/system/user`

**说明**: 添加新用户，自动关联到当前租户

**请求头**:
```
Authorization: Bearer <access_token>
```

**请求参数**:
```json
{
  "name": "张三",              // 姓名，必填
  "username": "zhangsan",     // 昵称，必填
  "account": "zhangsan001",   // 登录账号，必填，同租户内唯一
  "password": "123456",       // 密码，必填
  "phone": "13800138001",     // 手机号，必填
  "gender": 1,                // 性别，必填(1:男 2:女)
  "status": 1,                // 状态，必填(1:启用 2:禁用)
  "role_id": 2,               // 角色ID，必填
  "department_id": 1          // 部门ID，必填
}
```

### 2.4 更新用户

**接口**: `PUT /admin/system/user`

**说明**: 更新用户信息

**请求头**:
```
Authorization: Bearer <access_token>
```

**请求参数**:
```json
{
  "id": 2,                    // 用户ID，必填
  "name": "张三",              // 姓名，必填
  "username": "zhangsan",     // 昵称，必填
  "account": "zhangsan001",   // 登录账号，必填
  "password": "new_password", // 密码，可选，不填则不修改
  "phone": "13800138001",     // 手机号，必填
  "gender": 1,                // 性别，必填
  "status": 1,                // 状态，必填
  "role_id": 2,               // 角色ID，必填
  "department_id": 1          // 部门ID，必填
}
```

### 2.5 删除用户

**接口**: `DELETE /admin/system/user`

**说明**: 删除用户（不能删除超级管理员）

**请求头**:
```
Authorization: Bearer <access_token>
```

**请求参数**:
```json
{
  "id": 2  // 用户ID，必填
}
```

### 2.6 获取用户菜单

**接口**: `GET /admin/system/user/menu`

**说明**: 获取当前用户的菜单权限

**请求头**:
```
Authorization: Bearer <access_token>
```

**响应示例**:
```json
{
  "code": 0,
  "msg": "success",
  "data": [
    {
      "id": 1,
      "path": "/dashboard",
      "name": "Dashboard",
      "component": "Dashboard",
      "title": "仪表盘",
      "icon": "dashboard",
      "level": 1,
      "parent_id": 0,
      "sort": 100,
      "status": 1
    }
  ]
}
```

---

## 3. 租户管理 API

> **🔐 重要安全提示**: 租户管理API仅限超级管理员（用户ID=1）访问，所有操作都受到严格权限控制

### 3.1 查询租户列表

**接口**: `GET /admin/system/tenant`

**说明**: 查询租户列表（仅超级管理员）

**权限要求**: 超级管理员（用户ID=1）

**请求头**:
```
Authorization: Bearer <access_token>
```

**请求参数**:
```json
{
  "code": "default",       // 可选，企业编号模糊查询
  "name": "默认企业",      // 可选，企业名称模糊查询
  "status": 1,             // 可选，状态筛选(1:启用 2:禁用)
  "page": 1,               // 页码，默认1
  "page_size": 10          // 每页数量，默认10
}
```

**响应示例**:
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "total": 5,
    "list": [
      {
        "id": 1,
        "code": "default",
        "name": "默认企业",
        "contact": "管理员",
        "phone": "400-000-0000",
        "email": "admin@example.com",
        "address": "北京市朝阳区",
        "status": 1,
        "expired_at": "2025-12-31T23:59:59Z",
        "max_users": 100,
        "created_at": "2024-01-01T10:00:00Z"
      }
    ]
  }
}
```

### 3.2 添加租户

**接口**: `POST /admin/system/tenant`

**说明**: 添加新租户（仅超级管理员）

**权限要求**: 超级管理员（用户ID=1）

**请求头**:
```
Authorization: Bearer <access_token>
```

**请求参数**:
```json
{
  "code": "company001",              // 企业编号，必填，全局唯一
  "name": "某某科技有限公司",        // 企业名称，必填
  "contact": "张经理",               // 联系人，可选
  "phone": "400-123-4567",          // 联系电话，可选
  "email": "contact@company.com",   // 邮箱，可选
  "address": "上海市浦东新区",       // 地址，可选
  "status": 1,                      // 状态，必填(1:启用 2:禁用)
  "expired_at": "2025-12-31T23:59:59Z", // 过期时间，可选
  "max_users": 50                   // 最大用户数，可选，默认100
}
```

### 3.3 更新租户

**接口**: `PUT /admin/system/tenant`

**说明**: 更新租户信息（仅超级管理员）

**权限要求**: 超级管理员（用户ID=1）

**请求头**:
```
Authorization: Bearer <access_token>
```

**请求参数**:
```json
{
  "id": 1,                          // 租户ID，必填
  "code": "company001",             // 企业编号，必填
  "name": "某某科技有限公司",        // 企业名称，必填
  "contact": "张经理",               // 联系人，可选
  "phone": "400-123-4567",          // 联系电话，可选
  "email": "contact@company.com",   // 邮箱，可选
  "address": "上海市浦东新区",       // 地址，可选
  "status": 1,                      // 状态，必填
  "expired_at": "2025-12-31T23:59:59Z", // 过期时间，可选
  "max_users": 50                   // 最大用户数，可选
}
```

### 3.4 删除租户

**接口**: `DELETE /admin/system/tenant`

**说明**: 删除租户（仅超级管理员，谨慎操作，会影响关联数据）

**权限要求**: 超级管理员（用户ID=1）

**请求头**:
```
Authorization: Bearer <access_token>
```

**请求参数**:
```json
{
  "id": 1  // 租户ID，必填
}
```

---

## 4. 菜单管理 API

### 4.1 查询菜单列表

**接口**: `GET /admin/system/menu`

**说明**: 查询菜单列表

**请求头**:
```
Authorization: Bearer <access_token>
```

### 4.2 添加菜单

**接口**: `POST /admin/system/menu`

**说明**: 添加新菜单

**请求头**:
```
Authorization: Bearer <access_token>
```

### 4.3 更新菜单

**接口**: `PUT /admin/system/menu`

**说明**: 更新菜单信息

**请求头**:
```
Authorization: Bearer <access_token>
```

### 4.4 删除菜单

**接口**: `DELETE /admin/system/menu`

**说明**: 删除菜单

**请求头**:
```
Authorization: Bearer <access_token>
```

### 4.5 菜单权限管理

**接口**:
- `GET /admin/system/menu/auth` - 查询菜单权限
- `POST /admin/system/menu/auth` - 添加菜单权限
- `PUT /admin/system/menu/auth` - 更新菜单权限
- `DELETE /admin/system/menu/auth` - 删除菜单权限

### 4.6 角色菜单管理

**接口**:
- `GET /admin/system/menu/role` - 根据角色ID查询菜单
- `PUT /admin/system/menu/role` - 更新角色菜单关联

---

## 5. 部门管理 API

### 5.1 查询部门列表

**接口**: `GET /admin/system/department`

**说明**: 查询部门列表，自动按当前租户过滤

**请求头**:
```
Authorization: Bearer <access_token>
```

### 5.2 添加部门

**接口**: `POST /admin/system/department`

**说明**: 添加新部门，自动关联到当前租户

**请求头**:
```
Authorization: Bearer <access_token>
```

### 5.3 更新部门

**接口**: `PUT /admin/system/department`

**说明**: 更新部门信息

**请求头**:
```
Authorization: Bearer <access_token>
```

### 5.4 删除部门

**接口**: `DELETE /admin/system/department`

**说明**: 删除部门

**请求头**:
```
Authorization: Bearer <access_token>
```

---

## 6. 角色管理 API

### 6.1 查询角色列表

**接口**: `GET /admin/system/role`

**说明**: 查询角色列表，自动按当前租户过滤

**请求头**:
```
Authorization: Bearer <access_token>
```

### 6.2 添加角色

**接口**: `POST /admin/system/role`

**说明**: 添加新角色，自动关联到当前租户

**请求头**:
```
Authorization: Bearer <access_token>
```

### 6.3 更新角色

**接口**: `PUT /admin/system/role`

**说明**: 更新角色信息

**请求头**:
```
Authorization: Bearer <access_token>
```

### 6.4 删除角色

**接口**: `DELETE /admin/system/role`

**说明**: 删除角色

**请求头**:
```
Authorization: Bearer <access_token>
```

---

## 7. 错误码说明

### 通用错误码

| 错误码 | 状态 | 说明 | HTTP状态码 |
|--------|------|------|-----------|
| 200 | OK | 请求成功 | 200 |
| 400 | INVALID_ARGUMENT | 请求参数错误 | 400 |
| 400 | FAILED_PRECONDITION | 无法执行客户端请求 | 400 |
| 400 | OUT_OF_RANGE | 客户端越限访问 | 400 |
| 401 | UNAUTHENTICATED | 身份验证失败 | 401 |
| 403 | PERMISSION_DENIED | 客户端权限不足 | 403 |
| 404 | NOT_FOUND | 资源不存在 | 404 |
| 409 | ABORTED | 数据处理冲突 | 409 |
| 409 | ALREADY_EXISTS | 资源已存在 | 409 |
| 429 | RESOURCE_EXHAUSTED | 资源配额不足或速率限制 | 429 |
| 499 | CANCELLED | 请求被客户端取消 | 499 |
| 500 | DATA_LOSS | 处理数据发生错误 | 500 |
| 500 | UNKNOWN | 服务器未知错误 | 500 |
| 500 | INTERNAL | 服务器内部错误 | 500 |
| 501 | NOT_IMPLEMENTED | API不存在 | 501 |
| 503 | UNAVAILABLE | 服务不可用 | 503 |
| 504 | DEALINE_EXCEED | 请求超时 | 504 |

### 业务错误信息

- `"验证码错误"` - 登录验证码不正确
- `"账号或密码错误"` - 登录凭证错误
- `"账号已被禁用"` - 用户状态为禁用
- `"租户不存在或已禁用"` - 租户无效
- `"租户已过期"` - 租户过期时间已到
- `"Invalid tenant context"` - 租户上下文无效
- `"不能删除超级管理员"` - 保护超级管理员
- `"权限不足，只有超级管理员可以管理租户"` - 租户管理权限限制
- `"权限不足，需要管理员权限"` - 普通用户权限不足

---

## 8. 注意事项

### 8.1 多租户特性
1. **登录变更**: 前端需要在登录页面添加企业编号输入框
2. **数据隔离**: 所有查询自动添加租户过滤，无需手动处理
3. **账号唯一性**: 账号在同一租户内唯一，不同租户可重复
4. **权限控制**: 用户只能访问自己租户的数据

### 8.2 安全性
1. **JWT Token**: 包含租户ID、用户ID和账号信息
2. **Token验证**: 所有需要认证的接口都需要Bearer Token
3. **密码加密**: 使用bcrypt进行密码哈希
4. **登录日志**: 自动记录登录成功/失败日志

### 8.3 性能优化
1. **缓存查询**: 提供用户缓存查询接口提升性能
2. **分页查询**: 支持分页避免大量数据返回
3. **模糊查询**: 支持多字段模糊查询
4. **索引优化**: 租户ID和账号建立了复合索引

### 8.4 扩展性
1. **无限租户**: 支持无限数量租户
2. **独立配置**: 每个租户可独立配置最大用户数和到期时间
3. **灵活架构**: 所有业务表都支持多租户扩展

### 8.5 系统完善性 ✅

**已完成的重要修复**:
1. ✅ **租户管理路由**: 已在 `api/router.go` 中正确注册所有租户管理API路由
2. ✅ **权限控制机制**: 实现了`SuperAdminVerify`中间件，确保只有超级管理员能管理租户
3. ✅ **完整API文档**: 包含所有API端点的详细参数、响应格式和权限要求
4. ✅ **安全性加固**: 租户管理操作受到严格的身份验证和权限控制

**系统架构特点**:
- 🔒 **多层权限验证**: TokenVerify + SuperAdminVerify 双重验证
- 🏢 **完整租户生命周期**: 创建、查询、更新、删除全覆盖
- 🛡️ **数据安全隔离**: 租户数据完全隔离，防止跨租户访问
- ⚡ **高性能支持**: 缓存查询、分页查询、索引优化

---

## 9. 快速开始

### 9.1 前端集成步骤

1. **登录页面改造**
   - 添加企业编号输入框
   - 修改登录请求参数格式

2. **Token处理**
   - 保存登录返回的access_token
   - 在请求头中添加Authorization

3. **错误处理**
   - 401错误自动跳转登录页
   - 其他错误显示具体错误信息

4. **租户上下文**
   - 保存tenant_info用于显示
   - 用户只能看到同租户数据

### 9.2 测试账号

**默认租户**:
- 企业编号: `default`
- 企业名称: `默认企业`

**默认管理员**:
- 账号: `admin`
- 密码: 配置文件中的AdminPassword

---

这份API文档涵盖了多租户系统的核心功能，前端开发者可以根据此文档进行接口对接开发。