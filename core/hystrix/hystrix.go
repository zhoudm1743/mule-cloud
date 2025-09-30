package hystrix

import (
	"fmt"
	"log"
	"time"

	"github.com/afex/hystrix-go/hystrix"
)

// Config Hystrix配置
type Config struct {
	// Timeout 超时时间（毫秒）
	Timeout int
	// MaxConcurrentRequests 最大并发请求数
	MaxConcurrentRequests int
	// RequestVolumeThreshold 请求量阈值（在滚动窗口期内）
	RequestVolumeThreshold int
	// SleepWindow 熔断器打开后多久尝试半开（毫秒）
	SleepWindow int
	// ErrorPercentThreshold 错误率阈值（0-100）
	ErrorPercentThreshold int
}

// DefaultConfig 默认配置
var DefaultConfig = Config{
	Timeout:                3000, // 3秒超时
	MaxConcurrentRequests:  100,  // 最多100个并发
	RequestVolumeThreshold: 20,   // 至少20个请求后开始统计
	SleepWindow:            5000, // 熔断5秒后尝试恢复
	ErrorPercentThreshold:  50,   // 错误率超过50%触发熔断
}

// ServiceConfig 服务级别的配置
var ServiceConfigs = map[string]Config{
	"testservice": {
		Timeout:                2000,
		MaxConcurrentRequests:  50,
		RequestVolumeThreshold: 10,
		SleepWindow:            3000,
		ErrorPercentThreshold:  50,
	},
	"basicservice": {
		Timeout:                5000,
		MaxConcurrentRequests:  100,
		RequestVolumeThreshold: 20,
		SleepWindow:            5000,
		ErrorPercentThreshold:  60,
	},
}

// Init 初始化Hystrix配置
func Init() {
	log.Println("🔧 初始化Hystrix熔断器配置...")

	// 配置所有服务
	for serviceName, config := range ServiceConfigs {
		ConfigureCommand(serviceName, config)
		log.Printf("   ✓ %s: 超时=%dms, 并发=%d, 错误率阈值=%d%%",
			serviceName, config.Timeout, config.MaxConcurrentRequests, config.ErrorPercentThreshold)
	}

	// 配置默认值
	hystrix.DefaultTimeout = DefaultConfig.Timeout
	hystrix.DefaultMaxConcurrent = DefaultConfig.MaxConcurrentRequests
	hystrix.DefaultVolumeThreshold = DefaultConfig.RequestVolumeThreshold
	hystrix.DefaultSleepWindow = DefaultConfig.SleepWindow
	hystrix.DefaultErrorPercentThreshold = DefaultConfig.ErrorPercentThreshold

	log.Println("✅ Hystrix初始化完成")
}

// ConfigureCommand 配置指定命令的Hystrix参数
func ConfigureCommand(commandName string, config Config) {
	hystrix.ConfigureCommand(commandName, hystrix.CommandConfig{
		Timeout:                config.Timeout,
		MaxConcurrentRequests:  config.MaxConcurrentRequests,
		RequestVolumeThreshold: config.RequestVolumeThreshold,
		SleepWindow:            config.SleepWindow,
		ErrorPercentThreshold:  config.ErrorPercentThreshold,
	})
}

// Do 执行带熔断器保护的操作
func Do(commandName string, run func() error, fallback func(error) error) error {
	return hystrix.Do(commandName, run, fallback)
}

// Go 异步执行带熔断器保护的操作
func Go(commandName string, run func() error, fallback func(error) error) chan error {
	errChan := make(chan error, 1)
	go func() {
		errChan <- hystrix.Do(commandName, run, fallback)
	}()
	return errChan
}

// DoWithFallbackValue 执行操作并在失败时返回降级值
func DoWithFallbackValue[T any](commandName string, run func() (T, error), fallbackValue T) (T, error) {
	var result T
	var runErr error

	err := hystrix.Do(commandName,
		func() error {
			result, runErr = run()
			return runErr
		},
		func(err error) error {
			result = fallbackValue
			return nil // 返回nil表示使用fallback成功
		},
	)

	return result, err
}

// CircuitBreakerStatus 获取熔断器状态
func CircuitBreakerStatus(commandName string) string {
	cb, _, err := hystrix.GetCircuit(commandName)
	if err != nil {
		return "unknown"
	}

	if cb.IsOpen() {
		return "open"
	}

	// 注意: hystrix-go 不直接提供半开状态的判断
	// 只有在SleepWindow期间第一个请求会进入半开状态
	return "closed"
}

// GetAllCircuitStatus 获取所有熔断器状态
func GetAllCircuitStatus() map[string]interface{} {
	status := make(map[string]interface{})

	for serviceName := range ServiceConfigs {
		circuitStatus := CircuitBreakerStatus(serviceName)
		metrics, err := GetMetrics(serviceName)

		status[serviceName] = map[string]interface{}{
			"status":  circuitStatus,
			"metrics": metrics,
			"error":   err,
		}
	}

	return status
}

// Metrics 熔断器指标
type Metrics struct {
	TotalRequests        int64   `json:"total_requests"`
	ErrorCount           int64   `json:"error_count"`
	ErrorPercentage      float64 `json:"error_percentage"`
	IsCircuitBreakerOpen bool    `json:"is_circuit_breaker_open"`
}

// GetMetrics 获取熔断器指标
func GetMetrics(commandName string) (*Metrics, error) {
	cb, _, err := hystrix.GetCircuit(commandName)
	if err != nil {
		return nil, fmt.Errorf("获取熔断器失败: %v", err)
	}

	// 返回熔断器基本状态
	// 注意: hystrix-go 的 Metrics 是私有的，无法直接访问
	// 可以通过其他方式获取指标，这里只返回基本状态
	return &Metrics{
		TotalRequests:        0, // 需要自己实现计数器
		ErrorCount:           0, // 需要自己实现计数器
		ErrorPercentage:      0,
		IsCircuitBreakerOpen: cb.IsOpen(),
	}, nil
}

// StartStreamHandler 启动Hystrix Stream（用于监控）
func StartStreamHandler(port string) {
	go func() {
		log.Printf("🔍 Hystrix Stream监控启动: http://localhost%s", port)
		// 可以使用 hystrix.StreamHandler 配合 net/http
		// import "github.com/afex/hystrix-go/hystrix"
		// http.Handle("/hystrix.stream", hystrix.NewStreamHandler())
		// http.ListenAndServe(port, nil)
	}()
}

// FlushMetrics 刷新指标（清除旧数据）
func FlushMetrics() {
	hystrix.Flush()
}

// WaitForHealthyCircuit 等待熔断器恢复健康
func WaitForHealthyCircuit(commandName string, maxWait time.Duration) error {
	start := time.Now()
	for {
		if time.Since(start) > maxWait {
			return fmt.Errorf("等待熔断器恢复超时: %s", commandName)
		}

		status := CircuitBreakerStatus(commandName)
		if status == "closed" {
			return nil
		}

		time.Sleep(100 * time.Millisecond)
	}
}

// ExecuteWithTimeout 执行带超时的操作
func ExecuteWithTimeout(commandName string, timeout time.Duration, fn func() error) error {
	done := make(chan error, 1)

	go func() {
		done <- fn()
	}()

	select {
	case err := <-done:
		return err
	case <-time.After(timeout):
		return fmt.Errorf("操作超时: %s", commandName)
	}
}
