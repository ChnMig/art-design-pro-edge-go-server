package middleware

import (
	"fmt"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"

	"api-server/api/response"
	"api-server/config"
)

type limiterEntry struct {
	limiter    *rate.Limiter
	lastAccess time.Time
}

const (
	MaxLimiters     = 10000
	DefaultTTL      = 5 * time.Minute
	CleanupInterval = 2 * time.Minute
)

// RateLimiter 带自动清理的限流管理器
type RateLimiter struct {
	mu       sync.RWMutex
	limiters map[string]*limiterEntry
	rate     rate.Limit
	burst    int
	ttl      time.Duration
	ticker   *time.Ticker
	stopChan chan struct{}
}

func NewRateLimiter(r rate.Limit, b int) *RateLimiter {
	rl := &RateLimiter{
		limiters: make(map[string]*limiterEntry),
		rate:     r,
		burst:    b,
		ttl:      DefaultTTL,
		stopChan: make(chan struct{}),
	}
	rl.startCleanup(CleanupInterval)
	return rl
}

func (rl *RateLimiter) getLimiter(key string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	entry, exists := rl.limiters[key]
	if !exists {
		if len(rl.limiters) >= MaxLimiters {
			rl.cleanupLocked()
			if len(rl.limiters) >= MaxLimiters {
				rl.removeOldestLocked()
			}
		}
		limiter := rate.NewLimiter(rl.rate, rl.burst)
		rl.limiters[key] = &limiterEntry{
			limiter:    limiter,
			lastAccess: time.Now(),
		}
		return limiter
	}

	entry.lastAccess = time.Now()
	return entry.limiter
}

func (rl *RateLimiter) allow(key string) bool {
	limiter := rl.getLimiter(key)
	return limiter.Allow()
}

func (rl *RateLimiter) cleanup() {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	rl.cleanupLocked()
}

func (rl *RateLimiter) cleanupLocked() {
	now := time.Now()
	for key, entry := range rl.limiters {
		if now.Sub(entry.lastAccess) > rl.ttl {
			delete(rl.limiters, key)
		}
	}
}

func (rl *RateLimiter) removeOldestLocked() {
	if len(rl.limiters) == 0 {
		return
	}
	var oldestKey string
	var oldestTime time.Time
	first := true
	for key, entry := range rl.limiters {
		if first || entry.lastAccess.Before(oldestTime) {
			oldestKey = key
			oldestTime = entry.lastAccess
			first = false
		}
	}
	if oldestKey != "" {
		delete(rl.limiters, oldestKey)
	}
}

func (rl *RateLimiter) startCleanup(interval time.Duration) {
	rl.ticker = time.NewTicker(interval)
	go func() {
		for {
			select {
			case <-rl.ticker.C:
				rl.cleanup()
			case <-rl.stopChan:
				rl.ticker.Stop()
				return
			}
		}
	}()
}

func (rl *RateLimiter) Stop() {
	close(rl.stopChan)
}

type Stats struct {
	TotalLimiters int
	Rate          rate.Limit
	Burst         int
	TTL           time.Duration
}

func (rl *RateLimiter) GetStats() Stats {
	rl.mu.RLock()
	defer rl.mu.RUnlock()
	return Stats{
		TotalLimiters: len(rl.limiters),
		Rate:          rl.rate,
		Burst:         rl.burst,
		TTL:           rl.ttl,
	}
}

var (
	limiterCache = make(map[string]*RateLimiter)
	cacheMu      sync.RWMutex
)

// CleanupAllLimiters 关闭所有限流器（程序退出时调用）
func CleanupAllLimiters() {
	cacheMu.Lock()
	defer cacheMu.Unlock()
	for _, limiter := range limiterCache {
		limiter.Stop()
	}
	limiterCache = make(map[string]*RateLimiter)
}

func GetAllStats() map[string]Stats {
	cacheMu.RLock()
	defer cacheMu.RUnlock()
	result := make(map[string]Stats)
	for key, limiter := range limiterCache {
		result[key] = limiter.GetStats()
	}
	return result
}

func limiterKey(r rate.Limit, b int) string {
	return fmt.Sprintf("%.4f-%d", float64(r), b)
}

func getLimiterFromCache(r rate.Limit, b int) *RateLimiter {
	key := limiterKey(r, b)
	cacheMu.RLock()
	if limiter, exists := limiterCache[key]; exists {
		cacheMu.RUnlock()
		return limiter
	}
	cacheMu.RUnlock()

	cacheMu.Lock()
	defer cacheMu.Unlock()
	if limiter, exists := limiterCache[key]; exists {
		return limiter
	}
	limiter := NewRateLimiter(r, b)
	limiterCache[key] = limiter
	return limiter
}

type RateLimitOptions struct {
	Rate    rate.Limit
	Burst   int
	KeyFunc func(*gin.Context) string
	Message string
}

func IPRateLimit(r, b int) gin.HandlerFunc {
	return RateLimitWithOptions(RateLimitOptions{
		Rate:    rate.Limit(r),
		Burst:   b,
		KeyFunc: func(c *gin.Context) string { return c.ClientIP() },
		Message: "IP 请求过于频繁",
	})
}

func TokenRateLimit(r, b int) gin.HandlerFunc {
	return RateLimitWithOptions(RateLimitOptions{
		Rate:  rate.Limit(r),
		Burst: b,
		KeyFunc: func(c *gin.Context) string {
			return getTokenKey(c)
		},
		Message: "请求过于频繁",
	})
}

func RateLimitWithOptions(opts RateLimitOptions) gin.HandlerFunc {
	if opts.Rate <= 0 {
		opts.Rate = rate.Limit(1)
	}
	if opts.Burst <= 0 {
		opts.Burst = 1
	}
	if opts.KeyFunc == nil {
		opts.KeyFunc = func(c *gin.Context) string { return c.ClientIP() }
	}
	if opts.Message == "" {
		opts.Message = "请求过于频繁，请稍后再试"
	}
	limiter := getLimiterFromCache(opts.Rate, opts.Burst)

	return func(c *gin.Context) {
		key := opts.KeyFunc(c)
		if !limiter.allow(key) {
			response.ReturnError(c, response.RESOURCE_EXHAUSTED, opts.Message)
			return
		}
		c.Next()
	}
}

func getTokenKey(c *gin.Context) string {
	jwtData, exists := c.Get("jwtData")
	if !exists {
		return c.ClientIP()
	}
	switch v := jwtData.(type) {
	case string:
		if v != "" {
			return v
		}
	case map[string]interface{}:
		if id, ok := v["id"].(string); ok && id != "" {
			return id
		}
		if uid, ok := v["user_id"].(string); ok && uid != "" {
			return uid
		}
	}
	return c.ClientIP()
}

func StrictRateLimit() gin.HandlerFunc {
	return IPRateLimit(5, 10)
}

func ModerateRateLimit() gin.HandlerFunc {
	return IPRateLimit(50, 100)
}

func RelaxedRateLimit() gin.HandlerFunc {
	return IPRateLimit(100, 200)
}

// LoginRateLimitMiddleware 使用配置驱动的登录限流
func LoginRateLimitMiddleware() gin.HandlerFunc {
	ratePerSec := rate.Limit(float64(config.LoginRatePerMinute) / 60.0)
	if ratePerSec <= 0 {
		ratePerSec = rate.Limit(1)
	}
	return RateLimitWithOptions(RateLimitOptions{
		Rate:    ratePerSec,
		Burst:   config.LoginBurstSize,
		KeyFunc: func(c *gin.Context) string { return c.ClientIP() },
		Message: "登录请求过于频繁，请稍后再试",
	})
}

// GeneralRateLimitMiddleware 全局接口限流
func GeneralRateLimitMiddleware() gin.HandlerFunc {
	return RateLimitWithOptions(RateLimitOptions{
		Rate:    rate.Limit(config.GeneralRatePerSec),
		Burst:   config.GeneralBurstSize,
		KeyFunc: func(c *gin.Context) string { return c.ClientIP() },
		Message: "请求过于频繁，请稍后再试",
	})
}
