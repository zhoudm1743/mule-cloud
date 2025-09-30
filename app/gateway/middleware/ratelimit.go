package middleware

import (
	"mule-cloud/core/response"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// RateLimiter 限流器
type RateLimiter struct {
	rate      int                    // 每秒允许的请求数
	visitors  map[string]*visitor    // 访问者映射
	mu        sync.Mutex             // 互斥锁
	cleanupCh chan struct{}          // 清理信号
}

type visitor struct {
	limiter  chan struct{}
	lastSeen time.Time
}

// NewRateLimiter 创建限流器
func NewRateLimiter(requestsPerSecond int) *RateLimiter {
	rl := &RateLimiter{
		rate:      requestsPerSecond,
		visitors:  make(map[string]*visitor),
		cleanupCh: make(chan struct{}),
	}

	// 启动清理goroutine
	go rl.cleanup()

	return rl
}

// getVisitor 获取或创建访问者
func (rl *RateLimiter) getVisitor(ip string) *visitor {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	v, exists := rl.visitors[ip]
	if !exists {
		limiter := make(chan struct{}, rl.rate)
		v = &visitor{
			limiter:  limiter,
			lastSeen: time.Now(),
		}
		rl.visitors[ip] = v

		// 启动令牌生成
		go rl.generateTokens(v)
	}

	v.lastSeen = time.Now()
	return v
}

// generateTokens 生成令牌
func (rl *RateLimiter) generateTokens(v *visitor) {
	ticker := time.NewTicker(time.Second / time.Duration(rl.rate))
	defer ticker.Stop()

	for range ticker.C {
		select {
		case v.limiter <- struct{}{}:
		default:
			// 令牌桶满了
		}
	}
}

// cleanup 清理过期的访问者
func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			rl.mu.Lock()
			for ip, v := range rl.visitors {
				if time.Since(v.lastSeen) > 10*time.Minute {
					delete(rl.visitors, ip)
				}
			}
			rl.mu.Unlock()
		case <-rl.cleanupCh:
			return
		}
	}
}

// Middleware 限流中间件
func (rl *RateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		v := rl.getVisitor(ip)

		select {
		case <-v.limiter:
			// 获取到令牌，允许请求
			c.Next()
		case <-time.After(time.Second):
			// 超时，拒绝请求
			response.ErrorWithCode(c, 429, "请求过于频繁，请稍后再试")
			c.Abort()
		}
	}
}

// Close 关闭限流器
func (rl *RateLimiter) Close() {
	close(rl.cleanupCh)
}
