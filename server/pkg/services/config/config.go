package config

import (
	"fmt"

	"mule-cloud/pkg/services/validator"
)

// Config 应用程序配置结构
type Config struct {
	// 应用配置
	App AppConfig `mapstructure:"app" validate:"required"`

	// 服务器配置
	Server ServerConfig `mapstructure:"server" validate:"required"`

	// 数据库配置
	Database DatabaseConfig `mapstructure:"database" validate:"required"`

	// 日志配置
	Log LogConfig `mapstructure:"log" validate:"required"`

	// Redis配置
	Redis RedisConfig `mapstructure:"redis"`

	// JWT配置
	JWT JWTConfig `mapstructure:"jwt" validate:"required"`
}

// AppConfig 应用程序基本配置
type AppConfig struct {
	Name        string `mapstructure:"name" validate:"required"`
	Version     string `mapstructure:"version" validate:"required"`
	Environment string `mapstructure:"environment" validate:"required,oneof=development production testing"`
	Debug       bool   `mapstructure:"debug"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Host            string `mapstructure:"host" validate:"required"`
	Port            int    `mapstructure:"port" validate:"required,min=1,max=65535"`
	ReadTimeout     int    `mapstructure:"read_timeout" validate:"min=1"`
	WriteTimeout    int    `mapstructure:"write_timeout" validate:"min=1"`
	MaxHeaderBytes  int    `mapstructure:"max_header_bytes" validate:"min=1"`
	ShutdownTimeout int    `mapstructure:"shutdown_timeout" validate:"min=1"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Driver          string `mapstructure:"driver" validate:"required,oneof=mysql postgres sqlite"`
	Host            string `mapstructure:"host" validate:"required"`
	Port            int    `mapstructure:"port" validate:"required,min=1,max=65535"`
	Username        string `mapstructure:"username" validate:"required"`
	Password        string `mapstructure:"password"`
	Database        string `mapstructure:"database" validate:"required"`
	Charset         string `mapstructure:"charset"`
	MaxIdleConns    int    `mapstructure:"max_idle_conns" validate:"min=1"`
	MaxOpenConns    int    `mapstructure:"max_open_conns" validate:"min=1"`
	ConnMaxLifetime int    `mapstructure:"conn_max_lifetime" validate:"min=1"`
}

// LogConfig 日志配置
type LogConfig struct {
	Level      string `mapstructure:"level" validate:"required,oneof=debug info warn error fatal panic"`
	Format     string `mapstructure:"format" validate:"required,oneof=text json"`
	Output     string `mapstructure:"output" validate:"required,oneof=stdout stderr file"`
	FilePath   string `mapstructure:"file_path"`
	MaxSize    int    `mapstructure:"max_size" validate:"min=1"`
	MaxBackups int    `mapstructure:"max_backups" validate:"min=0"`
	MaxAge     int    `mapstructure:"max_age" validate:"min=1"`
	Compress   bool   `mapstructure:"compress"`
}

// RedisConfig Redis配置
type RedisConfig struct {
	Host         string `mapstructure:"host" validate:"required"`
	Port         int    `mapstructure:"port" validate:"required,min=1,max=65535"`
	Password     string `mapstructure:"password"`
	Database     int    `mapstructure:"database" validate:"min=0,max=15"`
	MaxRetries   int    `mapstructure:"max_retries" validate:"min=0"`
	PoolSize     int    `mapstructure:"pool_size" validate:"min=1"`
	MinIdleConns int    `mapstructure:"min_idle_conns" validate:"min=0"`
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret string `mapstructure:"secret" validate:"required"`
	Expire string `mapstructure:"expire" validate:"required"`
	Issuer string `mapstructure:"issuer" validate:"required"`
}

var (
	cfg           *Config
	validatorInst validator.Validator
)

func init() {
	// 初始化验证器
	validatorInst = validator.GetManager().GetValidator()
}

// GetConfig 获取全局配置实例
func GetConfig() *Config {
	return cfg
}

// GetApp 获取应用配置
func GetApp() *AppConfig {
	return &cfg.App
}

// GetServer 获取服务器配置
func GetServer() *ServerConfig {
	return &cfg.Server
}

// GetDatabase 获取数据库配置
func GetDatabase() *DatabaseConfig {
	return &cfg.Database
}

// GetLog 获取日志配置
func GetLog() *LogConfig {
	return &cfg.Log
}

// GetRedis 获取Redis配置
func GetRedis() *RedisConfig {
	return &cfg.Redis
}

// GetServerAddr 获取服务器监听地址
func (s *ServerConfig) GetServerAddr() string {
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}

// GetDSN 获取数据库连接字符串
func (d *DatabaseConfig) GetDSN() string {
	switch d.Driver {
	case "mysql":
		return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
			d.Username, d.Password, d.Host, d.Port, d.Database, d.Charset)
	case "postgres":
		return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			d.Host, d.Port, d.Username, d.Password, d.Database)
	case "sqlite":
		return d.Database
	default:
		return ""
	}
}

// GetRedisAddr 获取Redis连接地址
func (r *RedisConfig) GetRedisAddr() string {
	return fmt.Sprintf("%s:%d", r.Host, r.Port)
}

// IsDevelopment 判断是否为开发环境
func (a *AppConfig) IsDevelopment() bool {
	return a.Environment == "development"
}

// IsProduction 判断是否为生产环境
func (a *AppConfig) IsProduction() bool {
	return a.Environment == "production"
}

// IsTesting 判断是否为测试环境
func (a *AppConfig) IsTesting() bool {
	return a.Environment == "testing"
}
