package boot

import (
	"os"

	"mule-cloud/pkg/services/lifecycle"
	"mule-cloud/pkg/services/log"
)

// Setup 设置并启动应用
func Setup() error {
	// 获取配置文件路径
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "./configs/app.yaml"
	}

	// 创建生命周期管理器
	manager := lifecycle.NewManager()

	// 注册服务（按启动顺序）
	manager.Register(lifecycle.NewLogService())              // 1. 日志服务
	manager.Register(lifecycle.NewConfigService(configPath)) // 2. 配置服务
	manager.Register(lifecycle.NewDatabaseService())         // 3. 数据库服务
	manager.Register(lifecycle.NewHTTPService())             // 4. HTTP服务

	// 启动所有服务
	if err := manager.Start(); err != nil {
		log.Logger.Fatalf("应用启动失败: %v", err)
		return err
	}

	// 等待应用关闭
	manager.Wait()

	return nil
}

// GetApplicationStatus 获取应用状态（用于健康检查等）
func GetApplicationStatus() map[string]interface{} {
	// 这个函数可以被其他包调用来获取应用状态
	// 目前返回基本信息，后续可以扩展
	return map[string]interface{}{
		"status": "running",
		"name":   "mule-cloud",
	}
}
