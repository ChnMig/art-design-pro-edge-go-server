package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"

	"api-server/api/auth"
	"api-server/api/response"
	"api-server/db/pgdb/system"
)

const (
	currentUserContextKey   = "current_user"
	currentTenantContextKey = "current_tenant"
)

type tokenPrincipalLoader func(userID, tenantID uint) (system.SystemUser, system.SystemTenant, error)

// TokenVerify 多租户JWT认证中间件
func TokenVerify(c *gin.Context) {
	tokenVerify(c, loadTokenPrincipal)
}

func tokenVerify(c *gin.Context, loadPrincipal tokenPrincipalLoader) {
	c.FormFile("file") // 防止文件未发送完成就返回错误, 导致前端504而不是正确响应

	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		response.ReturnError(c, response.UNAUTHENTICATED, "未携带 token")
		c.Abort()
		return
	}

	var tokenString string
	// 支持两种格式: "Bearer token" 或直接 "token"
	if strings.HasPrefix(authHeader, "Bearer ") {
		tokenString = strings.TrimPrefix(authHeader, "Bearer ")
	} else {
		tokenString = authHeader
	}

	// 解析多租户JWT
	claims, err := auth.JWTDecrypt(tokenString)
	if err != nil {
		response.ReturnError(c, response.UNAUTHENTICATED, "token 解析失败")
		c.Abort()
		return
	}

	if claims.UserID == 0 || claims.TenantID == 0 || claims.Account == "" {
		response.ReturnError(c, response.UNAUTHENTICATED, "token 身份信息无效")
		c.Abort()
		return
	}

	user, tenant, err := loadPrincipal(claims.UserID, claims.TenantID)
	if err != nil || !validTokenPrincipal(claims, user, tenant) {
		response.ReturnError(c, response.UNAUTHENTICATED, "token 对应的用户或租户不可用")
		c.Abort()
		return
	}

	// 上下文只保存已从数据库复核过的身份，避免信任已过期的令牌载荷。
	user.Password = ""
	c.Set("tenant_id", tenant.ID)
	c.Set("user_id", user.ID)
	c.Set("account", user.Account)
	c.Set(currentUserContextKey, user)
	c.Set(currentTenantContextKey, tenant)

	c.Next()
}

func loadTokenPrincipal(userID, tenantID uint) (system.SystemUser, system.SystemTenant, error) {
	user := system.SystemUser{}
	user.ID = userID
	user.TenantID = tenantID
	if err := system.GetUser(&user); err != nil {
		return system.SystemUser{}, system.SystemTenant{}, err
	}

	tenant := system.SystemTenant{}
	tenant.ID = tenantID
	if err := system.GetTenant(&tenant); err != nil {
		return system.SystemUser{}, system.SystemTenant{}, err
	}
	return user, tenant, nil
}

func validTokenPrincipal(claims *auth.MultiTenantClaims, user system.SystemUser, tenant system.SystemTenant) bool {
	if claims == nil {
		return false
	}
	return user.ID == claims.UserID &&
		user.TenantID == claims.TenantID &&
		user.Account == claims.Account &&
		user.Status == system.StatusEnabled &&
		tenant.ID == claims.TenantID &&
		tenant.Status == system.StatusEnabled
}

// GetTenantID 从上下文获取租户ID
func GetTenantID(c *gin.Context) uint {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		return 0
	}
	value, ok := tenantID.(uint)
	if !ok {
		return 0
	}
	return value
}

// GetCurrentUserID 从上下文获取当前用户ID
func GetCurrentUserID(c *gin.Context) uint {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0
	}
	value, ok := userID.(uint)
	if !ok {
		return 0
	}
	return value
}

// GetCurrentAccount 从上下文获取当前账号
func GetCurrentAccount(c *gin.Context) string {
	account, exists := c.Get("account")
	if !exists {
		return ""
	}
	value, ok := account.(string)
	if !ok {
		return ""
	}
	return value
}

// GetCurrentUser 从上下文获取经过数据库复核的当前用户。
func GetCurrentUser(c *gin.Context) (system.SystemUser, bool) {
	value, exists := c.Get(currentUserContextKey)
	if !exists {
		return system.SystemUser{}, false
	}
	user, ok := value.(system.SystemUser)
	return user, ok
}

// GetCurrentTenant 从上下文获取经过数据库复核的当前租户。
func GetCurrentTenant(c *gin.Context) (system.SystemTenant, bool) {
	value, exists := c.Get(currentTenantContextKey)
	if !exists {
		return system.SystemTenant{}, false
	}
	tenant, ok := value.(system.SystemTenant)
	return tenant, ok
}
