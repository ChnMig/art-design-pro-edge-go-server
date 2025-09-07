package middleware

import (
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"

	"api-server/api/response"
	"api-server/config"
)

// IPRateLimiter 基于IP的速率限制器
type IPRateLimiter struct {
	ips map[string]*rate.Limiter
	mu  *sync.RWMutex
	r   rate.Limit
	b   int
}

// NewIPRateLimiter 创建新的IP速率限制器
// r: 限制速率 (每秒允许的请求数)
// b: 令牌桶容量 (突发请求数)
func NewIPRateLimiter(r rate.Limit, b int) *IPRateLimiter {
	i := &IPRateLimiter{
		ips: make(map[string]*rate.Limiter),
		mu:  &sync.RWMutex{},
		r:   r,
		b:   b,
	}

	// 启动清理goroutine，定期清理不活跃的IP限制器
	go i.cleanupLoop()

	return i
}

// AddIP 为指定IP添加速率限制器
func (i *IPRateLimiter) AddIP(ip string) *rate.Limiter {
	i.mu.Lock()
	defer i.mu.Unlock()

	limiter := rate.NewLimiter(i.r, i.b)
	i.ips[ip] = limiter

	return limiter
}

// GetLimiter 获取指定IP的速率限制器
func (i *IPRateLimiter) GetLimiter(ip string) *rate.Limiter {
	i.mu.Lock()
	limiter, exists := i.ips[ip]

	if !exists {
		i.mu.Unlock()
		return i.AddIP(ip)
	}

	i.mu.Unlock()
	return limiter
}

// cleanupLoop 定期清理不活跃的IP限制器，防止内存泄漏
func (i *IPRateLimiter) cleanupLoop() {
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		i.mu.Lock()
		// 清理超过1小时未使用的限制器
		cutoff := time.Now().Add(-time.Hour)
		for ip, limiter := range i.ips {
			// 如果限制器的令牌桶已满且很久没有请求，则删除
			if limiter.TokensAt(cutoff) >= float64(i.b) {
				delete(i.ips, ip)
			}
		}
		i.mu.Unlock()
	}
}

// RateLimitMiddleware 速率限制中间件
func RateLimitMiddleware(limiter *IPRateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		l := limiter.GetLimiter(ip)

		if !l.Allow() {
			response.ReturnError(c, response.RESOURCE_EXHAUSTED, "请求过于频繁，请稍后再试")
			c.Abort()
			return
		}

		c.Next()
	}
}

// LoginRateLimitMiddleware 专门用于登录的速率限制中间件
// 使用配置中的登录速率限制参数
func LoginRateLimitMiddleware() gin.HandlerFunc {
	// 计算每分钟的请求频率转换为时间间隔
	interval := time.Minute / time.Duration(config.LoginRatePerMinute)
	limiter := NewIPRateLimiter(rate.Every(interval), config.LoginBurstSize)
	return RateLimitMiddleware(limiter)
}

// GeneralRateLimitMiddleware 一般API的速率限制中间件
// 使用配置中的一般API速率限制参数
func GeneralRateLimitMiddleware() gin.HandlerFunc {
	limiter := NewIPRateLimiter(rate.Limit(config.GeneralRatePerSec), config.GeneralBurstSize)
	return RateLimitMiddleware(limiter)
}
