package config

import (
	"fmt"
	"sync"

	"mule-cloud/pkg/services/log"
)

var (
	once    sync.Once
	initErr error
)

// Init 初始化配置服务
func Init(configPath string) error {
	once.Do(func() {
		initErr = LoadConfig(configPath)
		if initErr != nil {
			return
		}
	})

	return initErr
}

// MustInit 初始化配置服务，失败时panic
func MustInit(configPath string) {
	if err := Init(configPath); err != nil {
		log.Logger.Fatalf("配置服务初始化失败: %v", err)
	}
}

// IsInitialized 检查配置是否已初始化
func IsInitialized() bool {
	return cfg != nil
}

// GetConfigSummary 获取配置摘要信息
func GetConfigSummary() map[string]interface{} {
	if !IsInitialized() {
		return map[string]interface{}{
			"initialized": false,
		}
	}

	return map[string]interface{}{
		"initialized": true,
		"app": map[string]interface{}{
			"name":        cfg.App.Name,
			"version":     cfg.App.Version,
			"environment": cfg.App.Environment,
			"debug":       cfg.App.Debug,
		},
		"server": map[string]interface{}{
			"address": cfg.Server.GetServerAddr(),
		},
		"database": map[string]interface{}{
			"driver": cfg.Database.Driver,
			"host":   fmt.Sprintf("%s:%d", cfg.Database.Host, cfg.Database.Port),
		},
		"redis": map[string]interface{}{
			"enabled": cfg.Redis.Host != "",
			"address": cfg.Redis.GetRedisAddr(),
		},
		"log": map[string]interface{}{
			"level":  cfg.Log.Level,
			"format": cfg.Log.Format,
			"output": cfg.Log.Output,
		},
		"jwt": map[string]interface{}{
			"secret": cfg.JWT.Secret,
			"expire": cfg.JWT.Expire,
			"issuer": cfg.JWT.Issuer,
		},
	}
}

// ValidateAndReload 验证并重新加载配置
func ValidateAndReload() error {
	if !IsInitialized() {
		return fmt.Errorf("配置服务未初始化")
	}

	oldConfig := *cfg

	if err := ReloadConfig(); err != nil {
		return fmt.Errorf("重新加载配置失败: %v", err)
	}

	log.Logger.Info("🔄 配置已重新加载")

	// 检查关键配置是否发生变化
	if oldConfig.Server.Port != cfg.Server.Port {
		log.Logger.Warnf("⚠️  服务器端口已更改: %d -> %d (需要重启服务器)",
			oldConfig.Server.Port, cfg.Server.Port)
	}

	if oldConfig.Log.Level != cfg.Log.Level {
		log.Logger.Infof("📝 日志级别已更改: %s -> %s",
			oldConfig.Log.Level, cfg.Log.Level)
	}

	return nil
}

// StartConfigWatcher 启动配置文件监听
func StartConfigWatcher() {
	if !IsInitialized() {
		log.Logger.Error("配置服务未初始化，无法启动配置监听")
		return
	}

	WatchConfig(func() {
		log.Logger.Info("🔍 检测到配置文件变化")
		if err := ValidateAndReload(); err != nil {
			log.Logger.Errorf("重新加载配置失败: %v", err)
		}
	})

}

// SetLogLevel 动态设置日志级别
func SetLogLevel(level string) error {
	if !IsInitialized() {
		return fmt.Errorf("配置服务未初始化")
	}

	// 验证日志级别
	validLevels := map[string]bool{
		"debug": true, "info": true, "warn": true,
		"error": true, "fatal": true, "panic": true,
	}

	if !validLevels[level] {
		return fmt.Errorf("无效的日志级别: %s", level)
	}

	cfg.Log.Level = level
	log.Logger.Infof("📝 日志级别已更新为: %s", level)

	return nil
}

// ToggleDebugMode 切换调试模式
func ToggleDebugMode() bool {
	if !IsInitialized() {
		return false
	}

	cfg.App.Debug = !cfg.App.Debug
	log.Logger.Infof("🐛 调试模式已%s",
		map[bool]string{true: "开启", false: "关闭"}[cfg.App.Debug])

	return cfg.App.Debug
}

// GetEnvironmentInfo 获取环境信息
func GetEnvironmentInfo() map[string]interface{} {
	if !IsInitialized() {
		return map[string]interface{}{
			"initialized": false,
		}
	}

	return map[string]interface{}{
		"environment":    cfg.App.Environment,
		"debug":          cfg.App.Debug,
		"is_development": cfg.App.IsDevelopment(),
		"is_production":  cfg.App.IsProduction(),
		"is_testing":     cfg.App.IsTesting(),
	}
}
