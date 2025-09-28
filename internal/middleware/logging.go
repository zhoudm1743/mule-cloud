package middleware

import (
	"bytes"
	"io"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zhoudm1743/mule-cloud/pkg/logger"
)

// LoggingConfig 日志配置
type LoggingConfig struct {
	EnableRequestBody  bool
	EnableResponseBody bool
	MaxBodySize        int64
	SkipPaths          []string
}

// DefaultLoggingConfig 默认日志配置
func DefaultLoggingConfig() LoggingConfig {
	return LoggingConfig{
		EnableRequestBody:  false,
		EnableResponseBody: false,
		MaxBodySize:        1024 * 1024, // 1MB
		SkipPaths: []string{
			"/health",
			"/metrics",
			"/favicon.ico",
		},
	}
}

// responseBodyWriter 响应体写入器
type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w responseBodyWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// LoggingMiddleware 日志中间件
func LoggingMiddleware(log logger.Logger, config ...LoggingConfig) gin.HandlerFunc {
	var cfg LoggingConfig
	if len(config) > 0 {
		cfg = config[0]
	} else {
		cfg = DefaultLoggingConfig()
	}

	return func(c *gin.Context) {
		// 检查是否跳过日志记录
		for _, skipPath := range cfg.SkipPaths {
			if c.Request.URL.Path == skipPath {
				c.Next()
				return
			}
		}

		startTime := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// 记录请求体
		var requestBody []byte
		if cfg.EnableRequestBody && c.Request.Body != nil {
			requestBody, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		// 记录响应体
		var responseWriter responseBodyWriter
		if cfg.EnableResponseBody {
			responseWriter = responseBodyWriter{
				ResponseWriter: c.Writer,
				body:           bytes.NewBufferString(""),
			}
			c.Writer = responseWriter
		}

		// 处理请求
		c.Next()

		// 计算响应时间
		latency := time.Since(startTime)

		// 构建日志字段
		fields := map[string]interface{}{
			"type":       "access",
			"method":     c.Request.Method,
			"path":       path,
			"status":     c.Writer.Status(),
			"latency":    latency.Milliseconds(),
			"client_ip":  c.ClientIP(),
			"user_agent": c.Request.UserAgent(),
		}

		if raw != "" {
			fields["query"] = raw
		}

		// 添加用户信息
		if userID := GetUserID(c); userID != "" {
			fields["user_id"] = userID
		}

		// 添加请求体
		if cfg.EnableRequestBody && len(requestBody) > 0 && len(requestBody) <= int(cfg.MaxBodySize) {
			fields["request_body"] = string(requestBody)
		}

		// 添加响应体
		if cfg.EnableResponseBody && responseWriter.body != nil {
			bodyStr := responseWriter.body.String()
			if len(bodyStr) <= int(cfg.MaxBodySize) {
				fields["response_body"] = bodyStr
			}
		}

		// 记录错误信息
		if len(c.Errors) > 0 {
			fields["errors"] = c.Errors.String()
		}

		// 根据状态码决定日志级别
		switch {
		case c.Writer.Status() >= 500:
			log.WithFields(fields).Error("Server error")
		case c.Writer.Status() >= 400:
			log.WithFields(fields).Warn("Client error")
		case c.Writer.Status() >= 300:
			log.WithFields(fields).Info("Redirection")
		default:
			log.WithFields(fields).Info("Request completed")
		}
	}
}

// RequestIDMiddleware 请求ID中间件
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.Request.Header.Get("X-Request-ID")
		if requestID == "" {
			requestID = generateRequestID()
		}

		c.Header("X-Request-ID", requestID)
		c.Set("request_id", requestID)
		c.Next()
	}
}

// generateRequestID 生成请求ID
func generateRequestID() string {
	// 这里可以使用UUID或其他方式生成唯一ID
	return time.Now().Format("20060102150405") + randomString(6)
}

// randomString 生成随机字符串
func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(b)
}
