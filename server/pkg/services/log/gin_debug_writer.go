package log

import (
	"regexp"
	"strings"
)

// GinDebugWriter 自定义Gin调试输出
type GinDebugWriter struct{}

func NewGinDebugWriter() *GinDebugWriter {
	return &GinDebugWriter{}
}

func (w *GinDebugWriter) Write(p []byte) (n int, err error) {
	message := string(p)

	// 处理不同类型的gin调试信息
	switch {
	case strings.Contains(message, "[GIN-debug]") && strings.Contains(message, "Running in \"debug\" mode"):
		// 启动警告信息
		Logger.Warn("🔧 Gin运行在调试模式，生产环境请切换到release模式")

	case strings.Contains(message, "using env:") || strings.Contains(message, "using code:"):
		// 模式切换提示，忽略这些信息
		return len(p), nil

	case strings.Contains(message, "[GIN-debug]") && (strings.Contains(message, "GET") ||
		strings.Contains(message, "POST") || strings.Contains(message, "PUT") ||
		strings.Contains(message, "DELETE") || strings.Contains(message, "PATCH")):
		// 路由注册信息
		w.handleRouteRegistration(message)

	default:
		// 其他gin调试信息
		if strings.TrimSpace(message) != "" && strings.Contains(message, "[GIN-debug]") {
			Logger.Debug(strings.TrimSpace(message))
		}
	}

	return len(p), nil
}

// handleRouteRegistration 处理路由注册信息
func (w *GinDebugWriter) handleRouteRegistration(message string) {
	// 使用正则表达式解析路由信息
	// 格式: [GIN-debug] GET    /api/admin/test/test      --> handler (3 handlers)
	re := regexp.MustCompile(`\[GIN-debug\]\s+(\w+)\s+([^\s]+)\s+-->\s+([^\s]+(?:\s+\([^)]+\))?)`)
	matches := re.FindStringSubmatch(message)

	if len(matches) >= 4 {
		method := matches[1]
		path := matches[2]
		handler := matches[3]

		// 获取方法对应的emoji
		methodEmoji := getMethodEmoji(method)

		// 美化输出
		Logger.Infof("%s 路由注册: %s %s → %s",
			methodEmoji, method, path, cleanHandlerName(handler))
	} else {
		// 如果解析失败，输出原始信息（去掉[GIN-debug]前缀）
		cleaned := strings.Replace(message, "[GIN-debug]", "", 1)
		cleaned = strings.TrimSpace(cleaned)
		if cleaned != "" {
			Logger.Debug("🔧 " + cleaned)
		}
	}
}

// getMethodEmoji 获取HTTP方法对应的emoji
func getMethodEmoji(method string) string {
	switch strings.ToUpper(method) {
	case "GET":
		return "🔍"
	case "POST":
		return "📝"
	case "PUT":
		return "✏️"
	case "DELETE":
		return "🗑️"
	case "PATCH":
		return "🔧"
	case "OPTIONS":
		return "⚙️"
	case "HEAD":
		return "📋"
	default:
		return "🔗"
	}
}

// cleanHandlerName 清理处理器名称
func cleanHandlerName(handler string) string {
	// 移除不必要的包路径前缀
	handler = strings.ReplaceAll(handler, "mule-cloud/", "")

	// 处理函数指针格式
	if strings.Contains(handler, ".func") {
		// 提取最后的函数名部分
		parts := strings.Split(handler, ".")
		if len(parts) > 0 {
			lastPart := parts[len(parts)-1]
			// 如果是匿名函数，保留前面的结构体信息
			if strings.HasPrefix(lastPart, "func") {
				if len(parts) >= 2 {
					return parts[len(parts)-2] + "." + lastPart
				}
			} else {
				return lastPart
			}
		}
	}

	// 移除包路径，只保留相对路径和函数名
	if strings.Contains(handler, "/") {
		parts := strings.Split(handler, "/")
		if len(parts) > 0 {
			return parts[len(parts)-1]
		}
	}

	return handler
}
