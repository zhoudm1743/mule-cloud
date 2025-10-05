package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"strings"
	"time"

	tenantCtx "mule-cloud/core/context"
	"mule-cloud/core/logger"
	"mule-cloud/internal/models"
	"mule-cloud/internal/repository"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// OperationLogMiddleware 操作日志中间件
// 记录用户的 CRUD 操作（POST、PUT、DELETE、PATCH）
// 不记录 GET、HEAD、OPTIONS 等读操作
func OperationLogMiddleware() gin.HandlerFunc {
	repo := repository.NewOperationLogRepository()

	return func(c *gin.Context) {
		// 1. 跳过读操作（GET、HEAD、OPTIONS）
		method := c.Request.Method
		if method == "GET" || method == "HEAD" || method == "OPTIONS" {
			c.Next()
			return
		}

		// 2. 跳过健康检查、静态资源等路径
		path := c.Request.URL.Path
		if shouldSkipPath(path) {
			c.Next()
			return
		}

		// 3. 记录开始时间
		startTime := time.Now()

		// 4. 读取请求体（需要缓存，因为 Body 只能读一次）
		var requestBody string
		if c.Request.Body != nil {
			bodyBytes, err := io.ReadAll(c.Request.Body)
			if err == nil {
				requestBody = string(bodyBytes)
				// 重新设置 Body，让后续处理器可以读取
				c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			}
		}

		// 5. 创建自定义 ResponseWriter 以捕获响应状态码
		writer := &responseWriter{
			ResponseWriter: c.Writer,
			statusCode:     200, // 默认 200
		}
		c.Writer = writer

		// 6. 执行后续处理器
		c.Next()

		// 7. 记录操作日志（异步，不阻塞请求）
		// 注意：需要在 goroutine 外部提取租户信息，避免 Context 被取消
		tenantCode, _ := c.Get("tenant_code")

		go func() {
			duration := time.Since(startTime).Milliseconds()

			// 获取用户信息
			userID, _ := c.Get("user_id")
			username, _ := c.Get("username")

			// 解析资源和操作类型
			resource, action := parseResourceAndAction(method, path)

			// 获取错误信息（如果有）
			var errorMsg string
			if len(c.Errors) > 0 {
				errorMsg = c.Errors.String()
			}

			// 脱敏请求体（移除密码等敏感信息）
			sanitizedBody := sanitizeRequestBody(requestBody)

			// 创建操作日志
			log := &models.OperationLog{
				UserID:       toString(userID),
				Username:     toString(username),
				Method:       method,
				Path:         path,
				Resource:     resource,
				Action:       action,
				RequestBody:  sanitizedBody,
				ResponseCode: writer.statusCode,
				Duration:     duration,
				IP:           c.ClientIP(),
				UserAgent:    c.Request.UserAgent(),
				Error:        errorMsg,
				CreatedAt:    startTime,
			}

			// ✅ 使用独立的 Context（避免主请求 Context 被取消）
			ctx := context.Background()
			ctx = tenantCtx.WithTenantCode(ctx, toString(tenantCode))

			// 保存到数据库（租户日志存到租户库，系统管理员日志存到系统库）
			if err := repo.Create(ctx, log); err != nil {
				logger.Error("保存操作日志失败",
					zap.String("user_id", log.UserID),
					zap.String("path", log.Path),
					zap.Error(err))
			} else {
				logger.Debug("操作日志已记录",
					zap.String("user_id", log.UserID),
					zap.String("resource", log.Resource),
					zap.String("action", log.Action),
					zap.Int("status", log.ResponseCode),
					zap.Int64("duration_ms", log.Duration))
			}
		}()
	}
}

// shouldSkipPath 判断是否跳过记录
func shouldSkipPath(path string) bool {
	skipPrefixes := []string{
		"/health",
		"/metrics",
		"/favicon.ico",
		"/static/",
		"/assets/",
		"/api/auth/refresh", // Token 刷新不记录
	}

	for _, prefix := range skipPrefixes {
		if strings.HasPrefix(path, prefix) {
			return true
		}
	}

	return false
}

// parseResourceAndAction 从路径和方法解析资源和操作类型
func parseResourceAndAction(method, path string) (resource, action string) {
	// 移除前缀
	path = strings.TrimPrefix(path, "/api/")
	path = strings.TrimPrefix(path, "/")

	// 分割路径
	parts := strings.Split(path, "/")
	if len(parts) > 0 {
		resource = parts[0] // 第一部分作为资源名称（如 users, roles, tenants）
	}

	// 根据 HTTP 方法确定操作类型
	switch method {
	case "POST":
		action = "create"
	case "PUT", "PATCH":
		action = "update"
	case "DELETE":
		action = "delete"
	default:
		action = "unknown"
	}

	return
}

// sanitizeRequestBody 脱敏请求体（移除密码等敏感信息）
func sanitizeRequestBody(body string) string {
	if body == "" {
		return ""
	}

	// 限制长度（避免存储过大的请求体）
	maxLen := 10000 // 10KB
	if len(body) > maxLen {
		body = body[:maxLen] + "...(truncated)"
	}

	// 尝试解析 JSON 并移除敏感字段
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(body), &data); err == nil {
		// 移除敏感字段
		sensitiveFields := []string{"password", "passwd", "token", "secret", "api_key", "apiKey"}
		for _, field := range sensitiveFields {
			if _, exists := data[field]; exists {
				data[field] = "***"
			}
		}

		// 转回 JSON
		if sanitized, err := json.Marshal(data); err == nil {
			return string(sanitized)
		}
	}

	return body
}

// toString 安全转换为字符串
func toString(v interface{}) string {
	if v == nil {
		return ""
	}
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}

// responseWriter 自定义 ResponseWriter 以捕获状态码
type responseWriter struct {
	gin.ResponseWriter
	statusCode int
}

func (w *responseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *responseWriter) Write(data []byte) (int, error) {
	return w.ResponseWriter.Write(data)
}
