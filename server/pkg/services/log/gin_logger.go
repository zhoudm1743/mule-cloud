package log

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// GinLogger 返回一个gin中间件，使用logrus记录请求日志
func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 开始时间
		startTime := time.Now()

		// 处理请求
		c.Next()

		// 结束时间
		endTime := time.Now()

		// 执行时间
		latencyTime := endTime.Sub(startTime)

		// 请求方式
		reqMethod := c.Request.Method

		// 请求路由
		reqUri := c.Request.RequestURI

		// 状态码
		statusCode := c.Writer.Status()

		// 请求IP
		clientIP := c.ClientIP()

		// 状态码颜色映射
		var statusColor string
		switch {
		case statusCode >= 200 && statusCode < 300:
			statusColor = "✅"
		case statusCode >= 300 && statusCode < 400:
			statusColor = "📝"
		case statusCode >= 400 && statusCode < 500:
			statusColor = "⚠️"
		case statusCode >= 500:
			statusColor = "❌"
		default:
			statusColor = "❓"
		}

		// 格式化延迟时间
		var latencyStr string
		if latencyTime < time.Millisecond {
			latencyStr = fmt.Sprintf("%.0fμs", float64(latencyTime.Nanoseconds())/1000)
		} else if latencyTime < time.Second {
			latencyStr = fmt.Sprintf("%.2fms", float64(latencyTime.Nanoseconds())/1e6)
		} else {
			latencyStr = fmt.Sprintf("%.2fs", latencyTime.Seconds())
		}

		// 美化的日志消息
		message := fmt.Sprintf("%s %s %s %s [%s] from %s",
			statusColor, reqMethod, reqUri,
			fmt.Sprintf("(%d)", statusCode), latencyStr, clientIP)

		// 根据状态码选择日志级别
		switch {
		case statusCode >= 400 && statusCode < 500:
			Logger.Warn(message)
		case statusCode >= 500:
			Logger.Error(message)
		default:
			Logger.Info(message)
		}
	}
}

// GinRecovery 返回一个gin恢复中间件，使用logrus记录panic
func GinRecovery() gin.HandlerFunc {
	return gin.RecoveryWithWriter(&LogrusWriter{logger: Logger})
}

// LogrusWriter 包装logrus.Logger以实现io.Writer接口
type LogrusWriter struct {
	logger *logrus.Logger
}

func (l *LogrusWriter) Write(p []byte) (n int, err error) {
	l.logger.Error(string(p))
	return len(p), nil
}

// LogRoute 记录路由注册信息
func LogRoute(method, path, handler string) {
	var methodEmoji string
	switch method {
	case "GET":
		methodEmoji = "🔍"
	case "POST":
		methodEmoji = "📝"
	case "PUT":
		methodEmoji = "✏️"
	case "DELETE":
		methodEmoji = "🗑️"
	case "PATCH":
		methodEmoji = "🔧"
	default:
		methodEmoji = "🔗"
	}

	Logger.Infof("%s 注册路由: %s %s → %s", methodEmoji, method, path, handler)
}
