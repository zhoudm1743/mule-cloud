package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"mule-cloud/pkg/services/validator"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// LoadConfig 加载配置文件
func LoadConfig(configPath string) error {

	// 设置配置文件路径
	if configPath == "" {
		configPath = getDefaultConfigPath()
	}

	// 配置viper
	setupViper(configPath)

	// 设置环境变量支持
	setupEnvironmentVariables()

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		// 如果配置文件不存在，创建默认配置
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			if err := createDefaultConfig(configPath); err != nil {
				return fmt.Errorf("创建默认配置文件失败: %v", err)
			}
			// 重新读取配置
			if err := viper.ReadInConfig(); err != nil {
				return fmt.Errorf("读取配置文件失败: %v", err)
			}
		} else {
			return fmt.Errorf("读取配置文件失败: %v", err)
		}
	}

	// 解析配置到结构体
	cfg = &Config{}
	if err := viper.Unmarshal(cfg); err != nil {
		return fmt.Errorf("解析配置失败: %v", err)
	}

	// 验证配置
	if err := ValidateConfig(cfg); err != nil {
		return fmt.Errorf("配置验证失败: %v", err)
	}

	return nil
}

// getDefaultConfigPath 获取默认配置文件路径
func getDefaultConfigPath() string {
	// 按优先级查找配置文件
	paths := []string{
		"./configs",
		"./config",
		"./",
		"../configs",
		"../config",
	}

	names := []string{
		"app.yaml",
		"app.yml",
		"config.yaml",
		"config.yml",
		"app.json",
		"config.json",
	}

	for _, path := range paths {
		for _, name := range names {
			fullPath := filepath.Join(path, name)
			if _, err := os.Stat(fullPath); err == nil {
				return fullPath
			}
		}
	}

	// 如果没有找到，返回默认路径
	return "./configs/app.yaml"
}

// setupViper 配置viper
func setupViper(configPath string) {
	// 获取文件信息
	dir := filepath.Dir(configPath)
	filename := filepath.Base(configPath)
	ext := filepath.Ext(filename)
	name := strings.TrimSuffix(filename, ext)

	// 配置viper
	viper.AddConfigPath(dir)
	viper.SetConfigName(name)
	viper.SetConfigType(strings.TrimPrefix(ext, "."))

	// 设置默认值
	setDefaultValues()
}

// setupEnvironmentVariables 设置环境变量支持
func setupEnvironmentVariables() {
	viper.AutomaticEnv()
	viper.SetEnvPrefix("MULE")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// 绑定特定环境变量
	envBindings := map[string]string{
		"app.environment":   "MULE_ENV",
		"app.debug":         "MULE_DEBUG",
		"server.port":       "MULE_PORT",
		"database.host":     "MULE_DB_HOST",
		"database.port":     "MULE_DB_PORT",
		"database.username": "MULE_DB_USER",
		"database.password": "MULE_DB_PASSWORD",
		"database.database": "MULE_DB_NAME",
		"redis.host":        "MULE_REDIS_HOST",
		"redis.port":        "MULE_REDIS_PORT",
		"redis.password":    "MULE_REDIS_PASSWORD",
		"jwt.secret":        "MULE_JWT_SECRET",
		"jwt.expire":        "MULE_JWT_EXPIRE",
		"jwt.issuer":        "MULE_JWT_ISSUER",
	}

	for key, env := range envBindings {
		viper.BindEnv(key, env)
	}
}

// setDefaultValues 设置默认配置值
func setDefaultValues() {
	// 应用默认配置
	viper.SetDefault("app.name", "mule-cloud")
	viper.SetDefault("app.version", "1.0.0")
	viper.SetDefault("app.environment", "development")
	viper.SetDefault("app.debug", true)

	// 服务器默认配置
	viper.SetDefault("server.host", "0.0.0.0")
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.read_timeout", 10)
	viper.SetDefault("server.write_timeout", 10)
	viper.SetDefault("server.max_header_bytes", 1048576) // 1MB
	viper.SetDefault("server.shutdown_timeout", 30)

	// 数据库默认配置
	viper.SetDefault("database.driver", "mysql")
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 3306)
	viper.SetDefault("database.username", "root")
	viper.SetDefault("database.password", "")
	viper.SetDefault("database.database", "mule_cloud")
	viper.SetDefault("database.charset", "utf8mb4")
	viper.SetDefault("database.max_idle_conns", 10)
	viper.SetDefault("database.max_open_conns", 100)
	viper.SetDefault("database.conn_max_lifetime", 3600)

	// 日志默认配置
	viper.SetDefault("log.level", "info")
	viper.SetDefault("log.format", "text")
	viper.SetDefault("log.output", "stdout")
	viper.SetDefault("log.file_path", "./logs/app.log")
	viper.SetDefault("log.max_size", 100)
	viper.SetDefault("log.max_backups", 3)
	viper.SetDefault("log.max_age", 7)
	viper.SetDefault("log.compress", true)

	// Redis默认配置
	viper.SetDefault("redis.host", "localhost")
	viper.SetDefault("redis.port", 6379)
	viper.SetDefault("redis.password", "")
	viper.SetDefault("redis.database", 0)
	viper.SetDefault("redis.max_retries", 3)
	viper.SetDefault("redis.pool_size", 10)
	viper.SetDefault("redis.min_idle_conns", 5)

	// JWT默认配置
	viper.SetDefault("jwt.secret", "a72cc3325e7d9f530d2468ebfb470373")
	viper.SetDefault("jwt.expire", "36h")
	viper.SetDefault("jwt.issuer", "mule-cloud")
}

// ValidateConfig 验证配置
func ValidateConfig(cfg *Config) error {
	result := validatorInst.Validate(cfg)
	if !result.Valid {
		return formatValidationResult(result)
	}
	return nil
}

// formatValidationResult 格式化验证结果
func formatValidationResult(result *validator.ValidationResult) error {
	var errors []string

	for _, fieldError := range result.Errors {
		errors = append(errors, fieldError.Message)
	}

	for _, structError := range result.StructErrors {
		errors = append(errors, structError)
	}

	return fmt.Errorf("配置验证错误:\n%s", strings.Join(errors, "\n"))
}

// ReloadConfig 重新加载配置
func ReloadConfig() error {
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("重新读取配置文件失败: %v", err)
	}

	newCfg := &Config{}
	if err := viper.Unmarshal(newCfg); err != nil {
		return fmt.Errorf("重新解析配置失败: %v", err)
	}

	if err := ValidateConfig(newCfg); err != nil {
		return fmt.Errorf("重新验证配置失败: %v", err)
	}

	cfg = newCfg
	return nil
}

// WatchConfig 监听配置文件变化
func WatchConfig(callback func()) {
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		if callback != nil {
			callback()
		}
	})
}
