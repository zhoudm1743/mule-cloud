package middleware

import (
	"bytes"
	"fmt"
	"io"
	"log"
	hystrixPkg "mule-cloud/core/hystrix"
	"net/http"
	"time"

	"github.com/afex/hystrix-go/hystrix"
	"github.com/gin-gonic/gin"
)

// HystrixMiddleware Hystrix熔断器中间件
func HystrixMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取服务名称（从context中获取，由proxyHandler设置）
		serviceName, exists := c.Get("service_name")
		if !exists {
			// 没有服务名称，直接通过
			c.Next()
			return
		}

		commandName := serviceName.(string)
		startTime := time.Now()

		// 使用Hystrix保护
		err := hystrix.Do(commandName,
			// 正常执行函数
			func() error {
				// 继续处理请求
				c.Next()

				// 检查响应状态码
				status := c.Writer.Status()
				if status >= 500 {
					return fmt.Errorf("服务返回错误: %d", status)
				}

				return nil
			},
			// 降级函数（熔断时执行）
			func(err error) error {
				duration := time.Since(startTime)
				log.Printf("[熔断触发] 服务: %s, 原因: %v, 耗时: %v", commandName, err, duration)

				// 返回降级响应
				c.JSON(http.StatusServiceUnavailable, gin.H{
					"code":    503,
					"msg":     fmt.Sprintf("服务暂时不可用: %s", commandName),
					"error":   err.Error(),
					"service": commandName,
					"fallback": true,
				})

				// 阻止继续执行
				c.Abort()
				return nil
			},
		)

		if err != nil {
			log.Printf("[熔断错误] 服务: %s, 错误: %v", commandName, err)
		}
	}
}

// HystrixConfig Hystrix配置中间件
func HystrixConfig(serviceName string, config hystrixPkg.Config) gin.HandlerFunc {
	hystrixPkg.ConfigureCommand(serviceName, config)
	
	return func(c *gin.Context) {
		c.Set("service_name", serviceName)
		c.Next()
	}
}

// CircuitBreakerCheck 熔断器状态检查（用于健康检查）
func CircuitBreakerCheck(serviceName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		status := hystrixPkg.CircuitBreakerStatus(serviceName)
		
		if status == "open" {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"code":    503,
				"msg":     "熔断器已打开",
				"service": serviceName,
				"status":  status,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// HystrixMetricsHandler 熔断器指标接口
func HystrixMetricsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		serviceName := c.Param("service")
		
		if serviceName == "" {
			// 返回所有服务的状态
			allStatus := hystrixPkg.GetAllCircuitStatus()
			c.JSON(http.StatusOK, gin.H{
				"code": 0,
				"msg":  "获取成功",
				"data": allStatus,
			})
			return
		}

		// 返回指定服务的状态
		metrics, err := hystrixPkg.GetMetrics(serviceName)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"code": 404,
				"msg":  fmt.Sprintf("服务不存在: %s", serviceName),
			})
			return
		}

		status := hystrixPkg.CircuitBreakerStatus(serviceName)
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  "获取成功",
			"data": gin.H{
				"service": serviceName,
				"status":  status,
				"metrics": metrics,
			},
		})
	}
}

// ResponseWriter 包装的响应写入器（用于捕获响应）
type responseWriter struct {
	gin.ResponseWriter
	body       *bytes.Buffer
	statusCode int
}

func newResponseWriter(w gin.ResponseWriter) *responseWriter {
	return &responseWriter{
		ResponseWriter: w,
		body:          &bytes.Buffer{},
		statusCode:    http.StatusOK,
	}
}

func (w *responseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w *responseWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

func (w *responseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

// HystrixWrapper 包装HTTP处理器（用于非Gin场景）
func HystrixWrapper(commandName string, handler http.Handler, fallbackHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := hystrix.Do(commandName,
			func() error {
				handler.ServeHTTP(w, r)
				return nil
			},
			func(err error) error {
				if fallbackHandler != nil {
					fallbackHandler.ServeHTTP(w, r)
				} else {
					http.Error(w, fmt.Sprintf("Service Unavailable: %v", err), http.StatusServiceUnavailable)
				}
				return nil
			},
		)

		if err != nil {
			log.Printf("[Hystrix错误] 命令: %s, 错误: %v", commandName, err)
		}
	})
}

// StreamBodyWrapper 包装请求体以支持重试
type streamBodyWrapper struct {
	body   []byte
	reader io.ReadCloser
}

func newStreamBodyWrapper(r io.ReadCloser) (*streamBodyWrapper, error) {
	body, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return &streamBodyWrapper{
		body:   body,
		reader: r,
	}, nil
}

func (s *streamBodyWrapper) Read(p []byte) (n int, err error) {
	return bytes.NewReader(s.body).Read(p)
}

func (s *streamBodyWrapper) Close() error {
	return s.reader.Close()
}

// Reset 重置读取器
func (s *streamBodyWrapper) Reset() {
	s.reader = io.NopCloser(bytes.NewReader(s.body))
}
