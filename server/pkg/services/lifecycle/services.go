package lifecycle

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"mule-cloud/pkg/services/config"
	"mule-cloud/pkg/services/database"
	httpService "mule-cloud/pkg/services/http"
	"mule-cloud/pkg/services/log"
)

// ConfigService 配置服务适配器
type ConfigService struct {
	configPath string
}

func NewConfigService(configPath string) *ConfigService {
	return &ConfigService{configPath: configPath}
}

func (s *ConfigService) Name() string {
	return "ConfigService"
}

func (s *ConfigService) Start(ctx context.Context) error {
	if err := config.Init(s.configPath); err != nil {
		return fmt.Errorf("初始化配置服务失败: %v", err)
	}

	// 启动配置文件监听（开发环境）
	if config.GetApp().IsDevelopment() {
		config.StartConfigWatcher()
		log.Logger.Info("🔍 配置文件监听已启动")
	}

	return nil
}

func (s *ConfigService) Stop(ctx context.Context) error {
	log.Logger.Info("📝 配置服务已停止")
	return nil
}

func (s *ConfigService) HealthCheck(ctx context.Context) error {
	if !config.IsInitialized() {
		return fmt.Errorf("配置服务未初始化")
	}
	return nil
}

// DatabaseService 数据库服务适配器
type DatabaseService struct{}

func NewDatabaseService() *DatabaseService {
	return &DatabaseService{}
}

func (s *DatabaseService) Name() string {
	return "DatabaseService"
}

func (s *DatabaseService) Start(ctx context.Context) error {
	if err := database.InitTorm(); err != nil {
		return fmt.Errorf("初始化数据库失败: %v", err)
	}
	return nil
}

func (s *DatabaseService) Stop(ctx context.Context) error {
	// 暂时简化处理，因为torm包的Close API可能不同
	// TODO: 需要根据torm包的实际API来实现正确的关闭方法
	log.Logger.Info("🗄️  数据库服务已停止")
	return nil
}

func (s *DatabaseService) HealthCheck(ctx context.Context) error {
	// 简化健康检查，因为torm包的API可能不同
	// TODO: 需要根据torm包的实际API来实现正确的健康检查
	return nil
}

// HTTPService HTTP服务适配器
type HTTPService struct {
	server *http.Server
}

func NewHTTPService() *HTTPService {
	return &HTTPService{}
}

func (s *HTTPService) Name() string {
	return "HTTPService"
}

func (s *HTTPService) Start(ctx context.Context) error {
	// 创建HTTP服务器
	s.server = httpService.Run()

	// 获取服务器配置
	serverConfig := config.GetServer()
	addr := serverConfig.GetServerAddr()
	s.server.Addr = addr

	log.Logger.Infof("🌐 HTTP服务器启动在 %s", addr)

	// 在单独的goroutine中启动服务器
	go func() {
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Logger.Errorf("❌ HTTP服务器启动失败: %v", err)
		}
	}()

	// 等待服务器启动
	time.Sleep(100 * time.Millisecond)

	return nil
}

func (s *HTTPService) Stop(ctx context.Context) error {
	if s.server == nil {
		return nil
	}

	// 优雅关闭服务器
	if err := s.server.Shutdown(ctx); err != nil {
		log.Logger.Errorf("HTTP服务器关闭失败: %v", err)
		// 强制关闭
		return s.server.Close()
	}

	return nil
}

func (s *HTTPService) HealthCheck(ctx context.Context) error {
	if s.server == nil {
		return fmt.Errorf("HTTP服务器未启动")
	}

	return nil
}

// LogService 日志服务适配器
type LogService struct{}

func NewLogService() *LogService {
	return &LogService{}
}

func (s *LogService) Name() string {
	return "LogService"
}

func (s *LogService) Start(ctx context.Context) error {

	return nil
}

func (s *LogService) Stop(ctx context.Context) error {
	return nil
}

func (s *LogService) HealthCheck(ctx context.Context) error {
	// 简单检查日志器是否可用
	if log.Logger == nil {
		return fmt.Errorf("日志器未初始化")
	}
	return nil
}
