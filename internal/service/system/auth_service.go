package system

import (
	"context"
	"errors"
	"time"

	"go.uber.org/zap"

	"api-server/internal/domain/system"
	"api-server/internal/pkg/config"
	"api-server/internal/repository/redis/captcha"
	"api-server/internal/transport/http/auth"
)

// AuthService 认证服务接口
type AuthService interface {
	// Login 用户登录
	Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error)
	// Logout 用户登出
	Logout(ctx context.Context, userID uint) error
	// RefreshToken 刷新令牌
	RefreshToken(ctx context.Context, token string) (*TokenResponse, error)
	// VerifyCaptcha 验证验证码
	VerifyCaptcha(captchaID, captcha string) bool
}

// LoginRequest 登录请求
type LoginRequest struct {
	TenantCode string `json:"tenant_code" binding:"required"` // 企业编号
	Account    string `json:"account" binding:"required"`     // 登录账号
	Password   string `json:"password" binding:"required"`    // 密码
	Captcha    string `json:"captcha" binding:"required"`     // 验证码
	CaptchaID  string `json:"captcha_id" binding:"required"`  // 验证码ID
	IP         string `json:"-"`                              // 客户端IP（由Handler层设置）
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token     string      `json:"token"`      // JWT令牌
	User      *system.SystemUser   `json:"user"`       // 用户信息
	Tenant    *system.SystemTenant `json:"tenant"`     // 租户信息
	ExpiresAt time.Time   `json:"expires_at"` // 令牌过期时间
}

// TokenResponse 令牌响应
type TokenResponse struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
}

// authService 认证服务实现
type authService struct {
	logger *zap.Logger
}

// NewAuthService 创建认证服务
func NewAuthService() AuthService {
	return &authService{
		logger: zap.L(),
	}
}

// Login 用户登录业务逻辑
func (s *authService) Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
	// 1. 验证验证码
	if !s.VerifyCaptcha(req.CaptchaID, req.Captcha) {
		s.logFailedLogin(req, "验证码错误")
		return nil, errors.New("验证码错误")
	}

	// 2. 验证用户和租户
	user, tenant, err := system.VerifyUser(req.TenantCode, req.Account, req.Password)
	if err != nil {
		s.logger.Error("用户验证失败", zap.Error(err))
		s.logFailedLogin(req, "验证失败")
		return nil, errors.New("查询用户失败")
	}

	// 3. 检查用户是否存在
	if user.ID == 0 {
		s.logFailedLogin(req, "账号或密码错误")
		return nil, errors.New("账号或密码错误")
	}

	// 4. 检查用户状态
	if user.Status != 1 {
		s.logFailedLogin(req, "账号已禁用")
		return nil, errors.New("账号已被禁用")
	}

	// 5. 验证租户状态
	if err := system.ValidateTenant(&tenant); err != nil {
		s.logFailedLogin(req, "租户无效")
		return nil, errors.New("租户状态异常")
	}

	// 6. 生成JWT令牌
	token, expiresAt, err := s.generateJWTToken(&user, &tenant)
	if err != nil {
		s.logger.Error("生成JWT令牌失败", zap.Error(err))
		return nil, errors.New("生成令牌失败")
	}

	// 7. 记录成功登录日志
	s.logSuccessLogin(&user, req.IP)

	return &LoginResponse{
		Token:     token,
		User:      &user,
		Tenant:    &tenant,
		ExpiresAt: expiresAt,
	}, nil
}

// Logout 用户登出
func (s *authService) Logout(ctx context.Context, userID uint) error {
	// 这里可以实现令牌黑名单逻辑
	s.logger.Info("用户登出", zap.Uint("user_id", userID))
	return nil
}

// RefreshToken 刷新令牌
func (s *authService) RefreshToken(ctx context.Context, token string) (*TokenResponse, error) {
	// 这里实现令牌刷新逻辑
	return nil, errors.New("暂未实现")
}

// VerifyCaptcha 验证验证码
func (s *authService) VerifyCaptcha(captchaID, captchaValue string) bool {
	return captcha.GetRedisStore().Verify(captchaID, captchaValue, true)
}

// generateJWTToken 生成JWT令牌
func (s *authService) generateJWTToken(user *system.SystemUser, tenant *system.SystemTenant) (string, time.Time, error) {
	token, err := auth.JWTIssue(user.ID, tenant.ID, user.Username)
	if err != nil {
		return "", time.Time{}, err
	}

	expiresAt := time.Now().Add(time.Duration(config.JWTExpiration) * time.Hour)
	return token, expiresAt, nil
}

// logFailedLogin 记录失败登录日志
func (s *authService) logFailedLogin(req *LoginRequest, reason string) {
	log := &system.SystemUserLoginLog{
		TenantCode:  req.TenantCode,
		UserName:    req.Account,
		IP:          req.IP,
		LoginStatus: "failed",
	}

	if err := system.CreateLoginLog(log); err != nil {
		s.logger.Error("记录登录失败日志失败", zap.Error(err))
	}

	s.logger.Warn("用户登录失败",
		zap.String("tenant_code", req.TenantCode),
		zap.String("account", req.Account),
		zap.String("ip", req.IP),
		zap.String("reason", reason),
	)
}

// logSuccessLogin 记录成功登录日志
func (s *authService) logSuccessLogin(user *system.SystemUser, ip string) {
	log := &system.SystemUserLoginLog{
		TenantCode:  "", // 需要通过租户ID获取租户编号
		UserName:    user.Username,
		IP:          ip,
		LoginStatus: "success",
	}

	if err := system.CreateLoginLog(log); err != nil {
		s.logger.Error("记录登录成功日志失败", zap.Error(err))
	}

	s.logger.Info("用户登录成功",
		zap.Uint("user_id", user.ID),
		zap.String("username", user.Username),
		zap.String("ip", ip),
	)
}