package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zhoudm1743/mule-cloud/pkg/logger"
)

// RecoveryMiddleware 恢复中间件
func RecoveryMiddleware(logger logger.Logger) gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		logger.Error("Panic recovered", "error", recovered, "path", c.Request.URL.Path)

		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Internal server error",
		})
	})
}

// HealthMiddleware 健康检查中间件
func HealthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.Path == "/health" {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"time":   time.Now(),
		})
			c.Abort()
			return
		}
		c.Next()
	}
}

// SecurityMiddleware 安全中间件
func SecurityMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 设置安全头
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		c.Header("Content-Security-Policy", "default-src 'self'")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")

		c.Next()
	}
}

// NoCache 禁用缓存中间件
func NoCache() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Cache-Control", "no-cache, no-store, max-age=0, must-revalidate, value")
		c.Header("Expires", "Thu, 01 Jan 1970 00:00:00 GMT")
		c.Header("Last-Modified", time.Now().Format(http.TimeFormat))
		c.Next()
	}
}

// ContentTypeMiddleware 内容类型中间件
func ContentTypeMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 如果没有设置Content-Type，默认为JSON
		if c.GetHeader("Content-Type") == "" {
			c.Header("Content-Type", "application/json; charset=utf-8")
		}
		c.Next()
	}
}

// TimeoutMiddleware 超时中间件
func TimeoutMiddleware(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 设置超时context
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()
		
		c.Request = c.Request.WithContext(ctx)
		c.Next()
		
		// 检查是否超时
		if ctx.Err() == context.DeadlineExceeded {
			c.JSON(http.StatusRequestTimeout, gin.H{
				"code":    408,
				"message": "Request timeout",
			})
			c.Abort()
		}
	}
}

// APIVersionMiddleware API版本中间件
func APIVersionMiddleware(version string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("API-Version", version)
		c.Next()
	}
}

// TrustedProxyMiddleware 信任代理中间件
func TrustedProxyMiddleware(trustedProxies []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 验证请求是否来自信任的代理
		clientIP := c.ClientIP()

		// 如果配置了信任的代理，验证IP
		if len(trustedProxies) > 0 {
			trusted := false
			for _, proxy := range trustedProxies {
				if proxy == clientIP || proxy == "*" {
					trusted = true
					break
				}
			}

			if !trusted {
				c.JSON(http.StatusForbidden, gin.H{
					"code":    403,
					"message": "Request from untrusted proxy",
				})
				c.Abort()
				return
			}
		}

		c.Next()
	}
}

// ResponseFormatMiddleware 响应格式中间件
func ResponseFormatMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// 如果响应没有被处理，返回404
		if c.Writer.Status() == http.StatusNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"code":    404,
				"message": "Resource not found",
				"path":    c.Request.URL.Path,
			})
		}
	}
}
