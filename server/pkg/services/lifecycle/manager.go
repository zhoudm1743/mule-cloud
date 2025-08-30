package lifecycle

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"mule-cloud/pkg/services/log"
)

// Service 定义生命周期服务接口
type Service interface {
	// Name 返回服务名称
	Name() string
	// Start 启动服务
	Start(ctx context.Context) error
	// Stop 停止服务
	Stop(ctx context.Context) error
	// HealthCheck 健康检查
	HealthCheck(ctx context.Context) error
}

// Manager 应用生命周期管理器
type Manager struct {
	services     []Service
	shutdownChan chan os.Signal
	ctx          context.Context
	cancel       context.CancelFunc
	wg           sync.WaitGroup
	mu           sync.RWMutex
	started      bool
}

// NewManager 创建新的生命周期管理器
func NewManager() *Manager {
	ctx, cancel := context.WithCancel(context.Background())
	return &Manager{
		services:     make([]Service, 0),
		shutdownChan: make(chan os.Signal, 1),
		ctx:          ctx,
		cancel:       cancel,
	}
}

// Register 注册服务
func (m *Manager) Register(service Service) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.started {
		log.Logger.Warnf("⚠️  应用已启动，无法注册新服务: %s", service.Name())
		return
	}

	m.services = append(m.services, service)
}

// Start 启动所有服务
func (m *Manager) Start() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.started {
		return fmt.Errorf("应用已经启动")
	}

	log.Logger.Info("🚀 开始启动应用...")

	// 设置信号监听
	signal.Notify(m.shutdownChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	// 启动信号处理协程
	m.wg.Add(1)
	go m.handleShutdown()

	// 按顺序启动所有服务
	startedServices := make([]Service, 0)
	for _, service := range m.services {

		startCtx, cancel := context.WithTimeout(m.ctx, 30*time.Second)
		err := service.Start(startCtx)
		cancel()

		if err != nil {
			log.Logger.Errorf("❌ 服务启动失败: %s, 错误: %v", service.Name(), err)

			// 回滚已启动的服务
			m.rollbackServices(startedServices)
			return fmt.Errorf("服务 %s 启动失败: %v", service.Name(), err)
		}

		startedServices = append(startedServices, service)
		log.Logger.Infof("✅ 服务启动成功: %s", service.Name())
	}

	m.started = true
	log.Logger.Info("🎉 所有服务启动完成")
	// 启动健康检查
	m.wg.Add(1)
	go m.healthCheckLoop()

	return nil
}

// Stop 停止所有服务
func (m *Manager) Stop() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.started {
		return fmt.Errorf("应用未启动")
	}

	// 取消上下文
	m.cancel()

	// 按相反顺序停止服务
	for i := len(m.services) - 1; i >= 0; i-- {
		service := m.services[i]

		stopCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		err := service.Stop(stopCtx)
		cancel()

		if err != nil {
			log.Logger.Errorf("⚠️  服务停止失败: %s, 错误: %v", service.Name(), err)
		}
	}

	m.started = false
	log.Logger.Info("👋 应用关闭完成")

	return nil
}

// Wait 等待应用关闭
func (m *Manager) Wait() {
	m.wg.Wait()
}

// rollbackServices 回滚已启动的服务
func (m *Manager) rollbackServices(services []Service) {
	for i := len(services) - 1; i >= 0; i-- {
		service := services[i]

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		if err := service.Stop(ctx); err != nil {
			log.Logger.Errorf("❌ 回滚服务失败: %s, 错误: %v", service.Name(), err)
		}
		cancel()
	}
}

// handleShutdown 处理关闭信号
func (m *Manager) handleShutdown() {
	defer m.wg.Done()

	select {
	case <-m.shutdownChan:
		if err := m.Stop(); err != nil {
			log.Logger.Errorf("❌ 应用关闭失败: %v", err)
		}
	case <-m.ctx.Done():
		return
	}
}

// healthCheckLoop 健康检查循环
func (m *Manager) healthCheckLoop() {
	defer m.wg.Done()

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			m.performHealthCheck()
		case <-m.ctx.Done():
			return
		}
	}
}

// performHealthCheck 执行健康检查
func (m *Manager) performHealthCheck() {
	m.mu.RLock()
	services := make([]Service, len(m.services))
	copy(services, m.services)
	m.mu.RUnlock()

	for _, service := range services {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		err := service.HealthCheck(ctx)
		cancel()

		if err != nil {
			log.Logger.Warnf("⚕️  服务健康检查失败: %s, 错误: %v", service.Name(), err)
		}
	}
}

// GetStatus 获取应用状态
func (m *Manager) GetStatus() map[string]interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()

	serviceStatus := make([]map[string]interface{}, len(m.services))
	for i, service := range m.services {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		err := service.HealthCheck(ctx)
		cancel()

		status := "healthy"
		if err != nil {
			status = "unhealthy"
		}

		serviceStatus[i] = map[string]interface{}{
			"name":   service.Name(),
			"status": status,
			"error":  err,
		}
	}

	return map[string]interface{}{
		"started":  m.started,
		"services": serviceStatus,
	}
}
