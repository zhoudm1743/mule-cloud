package config

import (
	"fmt"
	"os"
	"path/filepath"
)

// createDefaultConfig 创建默认配置文件
func createDefaultConfig(configPath string) error {
	// 确保目录存在
	dir := filepath.Dir(configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("创建配置目录失败: %v", err)
	}

	// 生成默认配置内容
	content := getDefaultConfigContent()

	// 写入配置文件
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("写入配置文件失败: %v", err)
	}

	return nil
}

// getDefaultConfigContent 获取默认配置文件内容
func getDefaultConfigContent() string {
	return `# Mule Cloud 应用配置文件

# 应用基本配置
app:
  name: "mule-cloud"
  version: "1.0.0"
  environment: "development"  # development, production, testing
  debug: true

# 服务器配置
server:
  host: "0.0.0.0"
  port: 8080
  read_timeout: 10          # 秒
  write_timeout: 10         # 秒
  max_header_bytes: 1048576 # 1MB
  shutdown_timeout: 30      # 秒

# 数据库配置
database:
  driver: "mysql"           # mysql, postgres, sqlite
  host: "localhost"
  port: 3306
  username: "root"
  password: ""
  database: "mule_cloud"
  charset: "utf8mb4"
  max_idle_conns: 10
  max_open_conns: 100
  conn_max_lifetime: 3600   # 秒

# 日志配置
log:
  level: "info"             # debug, info, warn, error, fatal, panic
  format: "text"            # text, json
  output: "stdout"          # stdout, stderr, file
  file_path: "./logs/app.log"
  max_size: 100             # MB
  max_backups: 3
  max_age: 7                # 天
  compress: true

# Redis配置 (可选)
redis:
  host: "localhost"
  port: 6379
  password: ""
  database: 0
  max_retries: 3
  pool_size: 10
  min_idle_conns: 5

# JWT配置
jwt:
  secret: "a72cc3325e7d9f530d2468ebfb470373"
  expire: "36h"
  issuer: "mule-cloud"
`
}

// ExportDefaultConfig 导出默认配置到指定路径
func ExportDefaultConfig(outputPath string) error {
	return createDefaultConfig(outputPath)
}

// GenerateConfigTemplate 生成配置模板
func GenerateConfigTemplate(format string, outputPath string) error {
	var content string

	switch format {
	case "json":
		content = getDefaultConfigJSON()
	case "yaml", "yml":
		content = getDefaultConfigContent()
	default:
		return fmt.Errorf("不支持的配置格式: %s", format)
	}

	// 确保目录存在
	dir := filepath.Dir(outputPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("创建目录失败: %v", err)
	}

	// 写入文件
	if err := os.WriteFile(outputPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("写入文件失败: %v", err)
	}

	return nil
}

// getDefaultConfigJSON 获取JSON格式的默认配置
func getDefaultConfigJSON() string {
	return `{
  "app": {
    "name": "mule-cloud",
    "version": "1.0.0", 
    "environment": "development",
    "debug": true
  },
  "server": {
    "host": "0.0.0.0",
    "port": 8080,
    "read_timeout": 10,
    "write_timeout": 10,
    "max_header_bytes": 1048576,
    "shutdown_timeout": 30
  },
  "database": {
    "driver": "mysql",
    "host": "localhost",
    "port": 3306,
    "username": "root",
    "password": "",
    "database": "mule_cloud",
    "charset": "utf8mb4",
    "max_idle_conns": 10,
    "max_open_conns": 100,
    "conn_max_lifetime": 3600
  },
  "log": {
    "level": "info",
    "format": "text",
    "output": "stdout",
    "file_path": "./logs/app.log",
    "max_size": 100,
    "max_backups": 3,
    "max_age": 7,
    "compress": true
  },
  "redis": {
    "host": "localhost",
    "port": 6379,
    "password": "",
    "database": 0,
    "max_retries": 3,
    "pool_size": 10,
    "min_idle_conns": 5
  },
  "jwt": {
    "secret": "a72cc3325e7d9f530d2468ebfb470373",
    "expire": "36h",
    "issuer": "mule-cloud"
  }
}`
}
