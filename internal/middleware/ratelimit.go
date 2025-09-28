package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zhoudm1743/mule-cloud/pkg/cache"
	"github.com/zhoudm1743/mule-cloud/pkg/logger"
)

// RateLimitConfig 限流配置
type RateLimitConfig struct {
	// 每秒允许的请求数
	RequestsPerSecond int
	// 突发请求数
	Burst int
	// 窗口大小（秒）
	WindowSize int
	// 键生成函数
	KeyGenerator func(c *gin.Context) string
}

// DefaultRateLimitConfig 默认限流配置
func DefaultRateLimitConfig() RateLimitConfig {
	return RateLimitConfig{
		RequestsPerSecond: 100,
		Burst:             20,
		WindowSize:        60,
		KeyGenerator: func(c *gin.Context) string {
			return c.ClientIP()
		},
	}
}

// RateLimitMiddleware 限流中间件
func RateLimitMiddleware(cache cache.Cache, logger logger.Logger, config ...RateLimitConfig) gin.HandlerFunc {
	var cfg RateLimitConfig
	if len(config) > 0 {
		cfg = config[0]
	} else {
		cfg = DefaultRateLimitConfig()
	}

	return func(c *gin.Context) {
		key := cfg.KeyGenerator(c)
		limiterKey := fmt.Sprintf("rate_limit:%s", key)

		// 使用滑动窗口算法
		allowed, remaining, resetTime, err := slidingWindowLimiter(
			cache,
			limiterKey,
			cfg.RequestsPerSecond,
			cfg.WindowSize,
		)

		if err != nil {
			logger.Error("Rate limit check failed", "error", err.Error(), "key", key)
			// 限流检查失败时，允许请求通过
			c.Next()
			return
		}

		// 设置响应头
		c.Header("X-RateLimit-Limit", strconv.Itoa(cfg.RequestsPerSecond))
		c.Header("X-RateLimit-Remaining", strconv.Itoa(remaining))
		c.Header("X-RateLimit-Reset", strconv.FormatInt(resetTime, 10))

		if !allowed {
			logger.Warn("Rate limit exceeded", "key", key, "limit", cfg.RequestsPerSecond)
			c.JSON(http.StatusTooManyRequests, gin.H{
				"code":        429,
				"message":     "Rate limit exceeded",
				"retry_after": resetTime - time.Now().Unix(),
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// slidingWindowLimiter 滑动窗口限流器
func slidingWindowLimiter(cache cache.Cache, key string, limit int, windowSize int) (bool, int, int64, error) {
	// 这里简化实现，实际应该使用Redis的EVAL命令
	// 由于我们的cache接口不支持脚本，使用简化版本
	return simpleLimiter(cache, key, limit, windowSize)
}

// simpleLimiter 简化的限流器
func simpleLimiter(cache cache.Cache, key string, limit int, windowSize int) (bool, int, int64, error) {

	// 获取当前计数
	currentStr, err := cache.Get(context.Background(), key)
	var current int
	if err != nil {
		current = 0
	} else {
		current, _ = strconv.Atoi(currentStr)
	}

	if current >= limit {
		resetTime := time.Now().Add(time.Duration(windowSize) * time.Second).Unix()
		return false, 0, resetTime, nil
	}

	// 增加计数
	newCount, err := cache.Incr(context.Background(), key)
	if err != nil {
		return false, 0, 0, err
	}

	// 设置过期时间
	if newCount == 1 {
		cache.Expire(context.Background(), key, time.Duration(windowSize)*time.Second)
	}

	remaining := limit - int(newCount)
	if remaining < 0 {
		remaining = 0
	}

	resetTime := time.Now().Add(time.Duration(windowSize) * time.Second).Unix()
	return true, remaining, resetTime, nil
}

// IPRateLimitMiddleware IP限流中间件
func IPRateLimitMiddleware(cache cache.Cache, logger logger.Logger, requestsPerSecond int) gin.HandlerFunc {
	config := RateLimitConfig{
		RequestsPerSecond: requestsPerSecond,
		WindowSize:        60,
		KeyGenerator: func(c *gin.Context) string {
			return fmt.Sprintf("ip:%s", c.ClientIP())
		},
	}
	return RateLimitMiddleware(cache, logger, config)
}

// UserRateLimitMiddleware 用户限流中间件
func UserRateLimitMiddleware(cache cache.Cache, logger logger.Logger, requestsPerSecond int) gin.HandlerFunc {
	config := RateLimitConfig{
		RequestsPerSecond: requestsPerSecond,
		WindowSize:        60,
		KeyGenerator: func(c *gin.Context) string {
			userID := GetUserID(c)
			if userID == "" {
				return fmt.Sprintf("ip:%s", c.ClientIP())
			}
			return fmt.Sprintf("user:%s", userID)
		},
	}
	return RateLimitMiddleware(cache, logger, config)
}

// APIKeyRateLimitMiddleware API密钥限流中间件
func APIKeyRateLimitMiddleware(cache cache.Cache, logger logger.Logger, requestsPerSecond int) gin.HandlerFunc {
	config := RateLimitConfig{
		RequestsPerSecond: requestsPerSecond,
		WindowSize:        60,
		KeyGenerator: func(c *gin.Context) string {
			apiKey := c.GetHeader("X-API-Key")
			if apiKey == "" {
				return fmt.Sprintf("ip:%s", c.ClientIP())
			}
			return fmt.Sprintf("api_key:%s", apiKey)
		},
	}
	return RateLimitMiddleware(cache, logger, config)
}
