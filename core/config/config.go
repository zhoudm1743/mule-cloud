package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// Config 全局配置
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Consul   ConsulConfig   `mapstructure:"consul"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	Hystrix  HystrixConfig  `mapstructure:"hystrix"`
	Gateway  GatewayConfig  `mapstructure:"gateway"`
	Database DatabaseConfig `mapstructure:"database"` // 已废弃
	MongoDB  MongoDBConfig  `mapstructure:"mongodb"`
	Redis    RedisConfig    `mapstructure:"redis"`
	Log      LogConfig      `mapstructure:"log"`
	Storage  StorageConfig  `mapstructure:"storage"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Name string `mapstructure:"name"`
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
	Mode string `mapstructure:"mode"` // debug, release, test
}

// ConsulConfig Consul配置
type ConsulConfig struct {
	Enabled bool   `mapstructure:"enabled"`
	Address string `mapstructure:"address"`
	Scheme  string `mapstructure:"scheme"`
	// 服务注册配置
	ServiceID   string   `mapstructure:"service_id"`
	ServiceName string   `mapstructure:"service_name"`
	ServiceIP   string   `mapstructure:"service_ip"`
	ServicePort int      `mapstructure:"service_port"`
	Tags        []string `mapstructure:"tags"`
	// 健康检查配置
	HealthCheckInterval string `mapstructure:"health_check_interval"`
	HealthCheckTimeout  string `mapstructure:"health_check_timeout"`
	DeregisterAfter     string `mapstructure:"deregister_after"`
}

// JWTConfig JWT配置
type JWTConfig struct {
	SecretKey  string `mapstructure:"secret_key"`
	ExpireTime int    `mapstructure:"expire_time"` // 小时
	Issuer     string `mapstructure:"issuer"`
}

// HystrixConfig Hystrix配置
type HystrixConfig struct {
	Enabled bool                            `mapstructure:"enabled"`
	Default HystrixCommandConfig            `mapstructure:"default"`
	Command map[string]HystrixCommandConfig `mapstructure:"commands"`
}

// HystrixCommandConfig Hystrix命令配置
type HystrixCommandConfig struct {
	Timeout                int `mapstructure:"timeout"`
	MaxConcurrentRequests  int `mapstructure:"max_concurrent_requests"`
	RequestVolumeThreshold int `mapstructure:"request_volume_threshold"`
	SleepWindow            int `mapstructure:"sleep_window"`
	ErrorPercentThreshold  int `mapstructure:"error_percent_threshold"`
}

// GatewayConfig 网关配置
type GatewayConfig struct {
	// 限流配置
	RateLimit RateLimitConfig `mapstructure:"rate_limit"`
	// 路由配置
	Routes map[string]RouteConfig `mapstructure:"routes"`
	// 超时配置
	Timeout TimeoutConfig `mapstructure:"timeout"`
}

// RateLimitConfig 限流配置
type RateLimitConfig struct {
	Enabled bool `mapstructure:"enabled"`
	Rate    int  `mapstructure:"rate"` // 每秒请求数
}

// RouteConfig 路由配置
type RouteConfig struct {
	ServiceName string   `mapstructure:"service_name"`
	RequireAuth bool     `mapstructure:"require_auth"`
	RequireRole []string `mapstructure:"require_role"`
}

// TimeoutConfig 超时配置
type TimeoutConfig struct {
	Read  int `mapstructure:"read"`  // 秒
	Write int `mapstructure:"write"` // 秒
}

// DatabaseConfig 数据库配置（已废弃，使用MongoDBConfig）
type DatabaseConfig struct {
	Type     string `mapstructure:"type"` // mysql, postgres, sqlite
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
	Charset  string `mapstructure:"charset"`
	MaxIdle  int    `mapstructure:"max_idle"`
	MaxOpen  int    `mapstructure:"max_open"`
}

// MongoDBConfig MongoDB配置
type MongoDBConfig struct {
	Enabled     bool   `mapstructure:"enabled"`
	URI         string `mapstructure:"uri"`           // MongoDB连接URI
	Host        string `mapstructure:"host"`          // 主机地址
	Port        int    `mapstructure:"port"`          // 端口
	Username    string `mapstructure:"username"`      // 用户名
	Password    string `mapstructure:"password"`      // 密码
	Database    string `mapstructure:"database"`      // 数据库名
	AuthSource  string `mapstructure:"auth_source"`   // 认证数据库
	ReplicaSet  string `mapstructure:"replica_set"`   // 副本集名称
	MaxPoolSize uint64 `mapstructure:"max_pool_size"` // 最大连接池大小
	MinPoolSize uint64 `mapstructure:"min_pool_size"` // 最小连接池大小
	Timeout     int    `mapstructure:"timeout"`       // 连接超时（秒）
}

// RedisConfig Redis配置
type RedisConfig struct {
	Enabled  bool   `mapstructure:"enabled"`
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
	PoolSize int    `mapstructure:"pool_size"`
}

// LogConfig 日志配置
type LogConfig struct {
	Level      string `mapstructure:"level"`       // debug, info, warn, error
	Format     string `mapstructure:"format"`      // json, text
	Output     string `mapstructure:"output"`      // stdout, file
	FilePath   string `mapstructure:"file_path"`   // 日志文件路径
	MaxSize    int    `mapstructure:"max_size"`    // MB
	MaxBackups int    `mapstructure:"max_backups"` // 保留旧文件数量
	MaxAge     int    `mapstructure:"max_age"`     // 天
	Compress   bool   `mapstructure:"compress"`    // 是否压缩
}

// StorageConfig 存储配置
type StorageConfig struct {
	Type  string      `mapstructure:"type"`  // 存储类型: local/minio/s3/oss
	Local LocalConfig `mapstructure:"local"` // 本地存储配置
	MinIO MinIOConfig `mapstructure:"minio"` // MinIO配置
	S3    S3Config    `mapstructure:"s3"`    // S3配置
	OSS   OSSConfig   `mapstructure:"oss"`   // OSS配置
}

// LocalConfig 本地存储配置
type LocalConfig struct {
	UploadDir string `mapstructure:"upload_dir"` // 上传目录
	BaseURL   string `mapstructure:"base_url"`   // 访问基础URL
}

// MinIOConfig MinIO配置
type MinIOConfig struct {
	Endpoint        string `mapstructure:"endpoint"`          // 端点
	AccessKeyID     string `mapstructure:"access_key_id"`     // 访问密钥ID
	SecretAccessKey string `mapstructure:"secret_access_key"` // 访问密钥
	BucketName      string `mapstructure:"bucket_name"`       // 存储桶名称
	UseSSL          bool   `mapstructure:"use_ssl"`           // 是否使用SSL
	Region          string `mapstructure:"region"`            // 区域
}

// S3Config S3配置
type S3Config struct {
	Endpoint        string `mapstructure:"endpoint"`          // 端点（可选，默认AWS S3）
	AccessKeyID     string `mapstructure:"access_key_id"`     // 访问密钥ID
	SecretAccessKey string `mapstructure:"secret_access_key"` // 访问密钥
	BucketName      string `mapstructure:"bucket_name"`       // 存储桶名称
	Region          string `mapstructure:"region"`            // 区域
}

// OSSConfig 阿里云OSS配置
type OSSConfig struct {
	Endpoint        string `mapstructure:"endpoint"`          // 端点
	AccessKeyID     string `mapstructure:"access_key_id"`     // 访问密钥ID
	AccessKeySecret string `mapstructure:"access_key_secret"` // 访问密钥
	BucketName      string `mapstructure:"bucket_name"`       // 存储桶名称
}

var (
	globalConfig *Config
	configOnce   sync.Once
	configErr    error
)

// Cfg 全局配置实例（懒加载）
var Cfg = &ConfigInstance{}

// Load 加载配置文件
func Load(configPath string) (*Config, error) {
	v := viper.New()

	// 设置配置文件路径
	if configPath != "" {
		v.SetConfigFile(configPath)
	} else {
		// 默认查找路径
		v.SetConfigName("config")
		v.SetConfigType("yaml")
		v.AddConfigPath(".")
		v.AddConfigPath("./config")
		v.AddConfigPath("../config")
		v.AddConfigPath("../../config")
	}

	// 环境变量支持
	v.SetEnvPrefix("MULE")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// 读取配置文件
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %v", err)
	}

	log.Printf("✅ 配置文件加载成功: %s", v.ConfigFileUsed())

	// 解析配置
	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %v", err)
	}

	// 环境变量覆盖
	overrideFromEnv(&cfg)

	globalConfig = &cfg
	return &cfg, nil
}

// Get 获取全局配置
func Get() *Config {
	if globalConfig == nil {
		log.Fatal("配置未初始化，请先调用 Load()")
	}
	return globalConfig
}

// overrideFromEnv 从环境变量覆盖配置
func overrideFromEnv(cfg *Config) {
	// Consul地址
	if addr := os.Getenv("CONSUL_ADDR"); addr != "" {
		cfg.Consul.Address = addr
	}

	// JWT密钥
	if secret := os.Getenv("JWT_SECRET"); secret != "" {
		cfg.JWT.SecretKey = secret
	}

	// 服务IP
	if ip := os.Getenv("SERVICE_IP"); ip != "" {
		cfg.Consul.ServiceIP = ip
	}

	// 服务端口
	if port := os.Getenv("SERVICE_PORT"); port != "" {
		fmt.Sscanf(port, "%d", &cfg.Server.Port)
	}
}

// GetConfigPath 获取配置文件路径
func GetConfigPath(customPath string) string {
	if customPath != "" {
		return customPath
	}

	// 查找配置文件
	possiblePaths := []string{
		"config.yaml",
		"config/config.yaml",
		"../config/config.yaml",
		"../../config/config.yaml",
	}

	for _, path := range possiblePaths {
		if _, err := os.Stat(path); err == nil {
			abs, _ := filepath.Abs(path)
			return abs
		}
	}

	return "config.yaml"
}

// Reload 重新加载配置
func Reload() error {
	if globalConfig == nil {
		return fmt.Errorf("配置未初始化")
	}

	cfg, err := Load("")
	if err != nil {
		return err
	}

	globalConfig = cfg
	log.Println("🔄 配置文件已重新加载")
	return nil
}

// Watch 监听配置文件变化
func Watch(callback func(*Config)) {
	v := viper.New()
	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		log.Printf("🔄 配置文件已更改: %s", e.Name)
		if err := Reload(); err != nil {
			log.Printf("❌ 重新加载配置失败: %v", err)
		} else if callback != nil {
			callback(globalConfig)
		}
	})
}

// ConfigInstance 全局配置实例包装器
type ConfigInstance struct {
	defaultPath string
}

// SetDefaultPath 设置默认配置文件路径
func (c *ConfigInstance) SetDefaultPath(path string) {
	c.defaultPath = path
}

// AutoLoad 自动加载配置（只执行一次）
func (c *ConfigInstance) AutoLoad() error {
	configOnce.Do(func() {
		path := c.defaultPath
		if path == "" {
			path = GetConfigPath("")
		}
		_, configErr = Load(path)
	})
	return configErr
}

// Get 获取配置（自动加载）
func (c *ConfigInstance) Get() *Config {
	if globalConfig == nil {
		c.AutoLoad()
	}
	return globalConfig
}

// MustGet 获取配置（自动加载，失败则panic）
func (c *ConfigInstance) MustGet() *Config {
	if err := c.AutoLoad(); err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}
	return globalConfig
}

// IsLoaded 检查配置是否已加载
func (c *ConfigInstance) IsLoaded() bool {
	return globalConfig != nil
}

// Reload 重新加载配置
func (c *ConfigInstance) Reload() error {
	return Reload()
}
