# Service层架构示例

## 🎯 为什么需要Service层？

### 当前问题
现在的Handler直接包含所有业务逻辑，导致：
- 代码复杂，难以测试
- 业务逻辑分散，难以复用
- 违反单一职责原则
- HTTP层和业务逻辑耦合

## 🏗️ Service层重构示例

### 1. 现在的代码结构
```
Handler (login.go)
├── 参数验证
├── 验证码验证
├── 用户查询和验证
├── 状态检查
├── 日志记录
├── JWT生成
└── 响应处理
```

### 2. 理想的代码结构
```
Handler (login.go)              ← 只处理HTTP相关
├── 参数验证和转换
└── 调用Service层

Service (auth_service.go)       ← 业务逻辑层
├── 验证码验证
├── 用户认证逻辑
├── 状态检查
├── 日志记录
└── JWT生成

Repository (user_repository.go) ← 数据访问层
├── 用户查询
├── 日志存储
└── 租户查询
```

## 🔧 具体实现示例

### Handler层 (简化后)
```go
func Login(c *gin.Context) {
    // 1. 参数验证和绑定
    req := &LoginRequest{}
    if !middleware.CheckParam(req, c) {
        return
    }

    // 2. 调用Service层处理业务逻辑
    authService := service.NewAuthService()
    result, err := authService.Login(c.Request.Context(), req)
    if err != nil {
        response.ReturnError(c, response.INTERNAL_ERROR, err.Error())
        return
    }

    // 3. 返回结果
    response.ReturnData(c, result)
}
```

### Service层 (业务逻辑)
```go
type AuthService interface {
    Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error)
    Logout(ctx context.Context, userID uint) error
    RefreshToken(ctx context.Context, token string) (*TokenResponse, error)
}

type authService struct {
    userRepo    repository.UserRepository
    tenantRepo  repository.TenantRepository
    captchaRepo repository.CaptchaRepository
    logger      *zap.Logger
}

func (s *authService) Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
    // 1. 验证验证码
    if !s.captchaRepo.Verify(req.CaptchaID, req.Captcha) {
        return nil, errors.New("验证码错误")
    }

    // 2. 验证用户和租户
    user, tenant, err := s.userRepo.VerifyUser(req.TenantCode, req.Account, req.Password)
    if err != nil {
        s.logFailedLogin(req, "验证失败")
        return nil, err
    }

    // 3. 检查用户状态
    if user.Status != 1 {
        s.logFailedLogin(req, "账号已禁用")
        return nil, errors.New("账号已被禁用")
    }

    // 4. 生成JWT Token
    token, err := s.generateJWTToken(user, tenant)
    if err != nil {
        return nil, err
    }

    // 5. 记录成功日志
    s.logSuccessLogin(user)

    return &LoginResponse{
        Token: token,
        User:  user,
    }, nil
}
```

### Repository层 (数据访问)
```go
type UserRepository interface {
    VerifyUser(tenantCode, account, password string) (User, Tenant, error)
    CreateLoginLog(log *LoginLog) error
    FindByID(id uint) (User, error)
}

type userRepository struct {
    db *gorm.DB
}

func (r *userRepository) VerifyUser(tenantCode, account, password string) (User, Tenant, error) {
    // 纯粹的数据库操作，不包含业务逻辑
    // ...
}
```

## 🚀 Service层的优势

### 1. 单一职责
- Handler只负责HTTP请求/响应
- Service只负责业务逻辑
- Repository只负责数据访问

### 2. 可测试性
```go
func TestAuthService_Login(t *testing.T) {
    // 可以轻松mock Repository依赖
    mockUserRepo := &MockUserRepository{}
    authService := NewAuthService(mockUserRepo, ...)

    // 测试业务逻辑，不依赖HTTP或数据库
    result, err := authService.Login(ctx, req)
    assert.NoError(t, err)
}
```

### 3. 业务逻辑复用
Service层的方法可以被多个Handler使用：
- HTTP API Handler
- gRPC Handler
- 定时任务
- 命令行工具

### 4. 依赖注入
```go
type AuthService struct {
    userRepo    UserRepository      // 接口，可替换
    logger      Logger             // 接口，可替换
    jwtService  JWTService         // 接口，可替换
}
```

## 📁 建议的目录结构

```
internal/
├── handler/           # HTTP处理器
│   └── system/
├── service/          # 业务服务层 (新增)
│   ├── auth.go       # 认证服务
│   ├── user.go       # 用户服务
│   └── tenant.go     # 租户服务
├── repository/       # 数据访问层
└── domain/          # 业务模型
```

## 🎯 总结

Service层让我们的架构：
- ✅ 更清晰 - 每层职责明确
- ✅ 更可测试 - 业务逻辑独立
- ✅ 更可复用 - 逻辑可被多处调用
- ✅ 更易维护 - 改动影响范围小

这就是为什么建议添加Service层的原因！