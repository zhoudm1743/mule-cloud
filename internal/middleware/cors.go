package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// CORSConfig CORS配置
type CORSConfig struct {
	AllowOrigins     []string
	AllowMethods     []string
	AllowHeaders     []string
	ExposeHeaders    []string
	AllowCredentials bool
	MaxAge           int
}

// DefaultCORSConfig 默认CORS配置
func DefaultCORSConfig() CORSConfig {
	return CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
			http.MethodOptions,
		},
		AllowHeaders: []string{
			"Origin",
			"Content-Length",
			"Content-Type",
			"Authorization",
			"X-Requested-With",
			"Accept",
			"Accept-Language",
			"Accept-Encoding",
		},
		ExposeHeaders: []string{
			"Content-Length",
			"Access-Control-Allow-Origin",
			"Access-Control-Allow-Headers",
			"Content-Type",
		},
		AllowCredentials: true,
		MaxAge:           12 * 3600, // 12 hours
	}
}

// CORSMiddleware CORS中间件
func CORSMiddleware(config ...CORSConfig) gin.HandlerFunc {
	var corsConfig CORSConfig
	if len(config) > 0 {
		corsConfig = config[0]
	} else {
		corsConfig = DefaultCORSConfig()
	}

	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// 检查是否允许该源
		allowOrigin := false
		for _, allowedOrigin := range corsConfig.AllowOrigins {
			if allowedOrigin == "*" || allowedOrigin == origin {
				allowOrigin = true
				break
			}
		}

		if allowOrigin {
			c.Header("Access-Control-Allow-Origin", origin)
		}

		// 设置允许的方法
		if len(corsConfig.AllowMethods) > 0 {
			methods := ""
			for i, method := range corsConfig.AllowMethods {
				if i > 0 {
					methods += ", "
				}
				methods += method
			}
			c.Header("Access-Control-Allow-Methods", methods)
		}

		// 设置允许的头部
		if len(corsConfig.AllowHeaders) > 0 {
			headers := ""
			for i, header := range corsConfig.AllowHeaders {
				if i > 0 {
					headers += ", "
				}
				headers += header
			}
			c.Header("Access-Control-Allow-Headers", headers)
		}

		// 设置暴露的头部
		if len(corsConfig.ExposeHeaders) > 0 {
			headers := ""
			for i, header := range corsConfig.ExposeHeaders {
				if i > 0 {
					headers += ", "
				}
				headers += header
			}
			c.Header("Access-Control-Expose-Headers", headers)
		}

		// 设置是否允许凭证
		if corsConfig.AllowCredentials {
			c.Header("Access-Control-Allow-Credentials", "true")
		}

		// 设置预检请求的缓存时间
		if corsConfig.MaxAge > 0 {
			c.Header("Access-Control-Max-Age", string(rune(corsConfig.MaxAge)))
		}

		// 处理预检请求
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
