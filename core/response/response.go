package response

import (
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// Response 统一响应结构
type Response struct {
	Code      int         `json:"code"`
	Msg       string      `json:"msg"`
	Data      interface{} `json:"data,omitempty"`
	Timestamp int64       `json:"timestamp,omitempty"`
	RequestID string      `json:"request_id,omitempty"`
}

// ErrorCode 错误码定义
const (
	CodeSuccess            = 0
	CodeBadRequest         = 400
	CodeUnauthorized       = 401
	CodeForbidden          = 403
	CodeNotFound           = 404
	CodeMethodNotAllowed   = 405
	CodeConflict           = 409
	CodeTooManyRequests    = 429
	CodeInternalError      = 500
	CodeServiceUnavailable = 503
	CodeGatewayTimeout     = 504
)

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	resp := Response{
		Code:      CodeSuccess,
		Msg:       "success",
		Data:      data,
		Timestamp: time.Now().Unix(),
		RequestID: getRequestID(c),
	}
	c.JSON(http.StatusOK, resp)
}

// SuccessWithMsg 带消息的成功响应
func SuccessWithMsg(c *gin.Context, msg string, data interface{}) {
	resp := Response{
		Code:      CodeSuccess,
		Msg:       msg,
		Data:      data,
		Timestamp: time.Now().Unix(),
		RequestID: getRequestID(c),
	}
	c.JSON(http.StatusOK, resp)
}

// Error 错误响应（默认错误码 -1）
func Error(c *gin.Context, msg string) {
	ErrorWithCode(c, -1, msg)
}

// ErrorWithCode 带错误码的错误响应
func ErrorWithCode(c *gin.Context, code int, msg string) {
	resp := Response{
		Code:      code,
		Msg:       msg,
		Data:      nil,
		Timestamp: time.Now().Unix(),
		RequestID: getRequestID(c),
	}
	
	// 根据错误码返回对应的HTTP状态码
	httpStatus := getHTTPStatus(code)
	c.JSON(httpStatus, resp)
}

// BadRequest 400 错误
func BadRequest(c *gin.Context, msg string) {
	ErrorWithCode(c, CodeBadRequest, msg)
}

// Unauthorized 401 未授权
func Unauthorized(c *gin.Context, msg string) {
	ErrorWithCode(c, CodeUnauthorized, msg)
}

// Forbidden 403 禁止访问
func Forbidden(c *gin.Context, msg string) {
	ErrorWithCode(c, CodeForbidden, msg)
}

// NotFound 404 未找到
func NotFound(c *gin.Context, msg string) {
	ErrorWithCode(c, CodeNotFound, msg)
}

// InternalError 500 内部错误
func InternalError(c *gin.Context, msg string) {
	ErrorWithCode(c, CodeInternalError, msg)
}

// ServiceUnavailable 503 服务不可用
func ServiceUnavailable(c *gin.Context, msg string) {
	ErrorWithCode(c, CodeServiceUnavailable, msg)
}

// getRequestID 获取请求ID
func getRequestID(c *gin.Context) string {
	if requestID, exists := c.Get("request_id"); exists {
		return requestID.(string)
	}
	return c.GetString("X-Request-ID")
}

// getHTTPStatus 根据业务错误码映射HTTP状态码
func getHTTPStatus(code int) int {
	switch code {
	case CodeSuccess:
		return http.StatusOK
	case CodeBadRequest:
		return http.StatusBadRequest
	case CodeUnauthorized:
		return http.StatusUnauthorized
	case CodeForbidden:
		return http.StatusForbidden
	case CodeNotFound:
		return http.StatusNotFound
	case CodeMethodNotAllowed:
		return http.StatusMethodNotAllowed
	case CodeConflict:
		return http.StatusConflict
	case CodeTooManyRequests:
		return http.StatusTooManyRequests
	case CodeInternalError:
		return http.StatusInternalServerError
	case CodeServiceUnavailable:
		return http.StatusServiceUnavailable
	case CodeGatewayTimeout:
		return http.StatusGatewayTimeout
	default:
		return http.StatusOK
	}
}

// UnifiedResponseMiddleware 统一响应中间件（排除健康检查）
func UnifiedResponseMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 跳过健康检查等特殊路径
		if shouldSkipUnifiedResponse(c.Request.URL.Path) {
			c.Next()
			return
		}

		// 继续处理
		c.Next()
	}
}

// RecoveryMiddleware 统一错误恢复中间件
func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 跳过健康检查等特殊路径
				if shouldSkipUnifiedResponse(c.Request.URL.Path) {
					c.Next()
					return
				}

				// 记录错误日志
				log.Printf("[Panic恢复] 错误: %v\n堆栈:\n%s", err, string(debug.Stack()))

				// 返回统一错误响应
				InternalError(c, fmt.Sprintf("服务器内部错误: %v", err))
				c.Abort()
			}
		}()
		c.Next()
	}
}

// shouldSkipUnifiedResponse 判断是否跳过统一响应
func shouldSkipUnifiedResponse(path string) bool {
	// 排除的路径列表
	skipPaths := []string{
		"/api/health",
		"/gateway/health",
		"/health",
		"/ping",
		"/metrics",
		"/hystrix.stream",
	}

	for _, skipPath := range skipPaths {
		if strings.HasPrefix(path, skipPath) || path == skipPath {
			return true
		}
	}
	return false
}