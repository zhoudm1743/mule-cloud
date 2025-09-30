package cache

import (
	"context"
	"fmt"
	"log"
	"mule-cloud/core/config"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	redisClient *redis.Client
	redisOnce   sync.Once
	redisErr    error
)

// Redis 全局Redis实例（懒加载）
var Redis = &RedisInstance{}

// InitRedis 初始化Redis连接
func InitRedis(cfg *config.RedisConfig) (*redis.Client, error) {
	if !cfg.Enabled {
		log.Println("⚠️  Redis未启用")
		return nil, nil
	}

	// 创建Redis客户端
	client := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password:     cfg.Password,
		DB:           cfg.DB,
		PoolSize:     cfg.PoolSize,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		PoolTimeout:  4 * time.Second,
	})

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("连接Redis失败: %v", err)
	}

	log.Printf("✅ Redis连接成功: %s:%d (DB:%d)", cfg.Host, cfg.Port, cfg.DB)

	// 保存全局实例
	redisClient = client

	return client, nil
}

// GetRedis 获取Redis客户端实例
func GetRedis() *redis.Client {
	if redisClient == nil {
		log.Fatal("Redis未初始化，请先调用 InitRedis()")
	}
	return redisClient
}

// CloseRedis 关闭Redis连接
func CloseRedis() error {
	if redisClient == nil {
		return nil
	}

	if err := redisClient.Close(); err != nil {
		return fmt.Errorf("关闭Redis连接失败: %v", err)
	}

	log.Println("✅ Redis连接已关闭")
	return nil
}

// PingRedis 检查Redis连接状态
func PingRedis() error {
	if redisClient == nil {
		return fmt.Errorf("Redis未初始化")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return redisClient.Ping(ctx).Err()
}

// RedisHealth 获取Redis健康状态
func RedisHealth() map[string]interface{} {
	status := map[string]interface{}{
		"status": "unknown",
	}

	if redisClient == nil {
		status["status"] = "not_initialized"
		return status
	}

	if err := PingRedis(); err != nil {
		status["status"] = "unhealthy"
		status["error"] = err.Error()
	} else {
		status["status"] = "healthy"

		// 获取Redis信息
		ctx := context.Background()
		if info, err := redisClient.Info(ctx, "server").Result(); err == nil {
			status["info"] = "connected"
		} else {
			status["info"] = info
		}
	}

	return status
}

// ===== 常用Redis操作封装 =====

// Set 设置键值
func Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return GetRedis().Set(ctx, key, value, expiration).Err()
}

// Get 获取键值
func Get(ctx context.Context, key string) (string, error) {
	return GetRedis().Get(ctx, key).Result()
}

// Del 删除键
func Del(ctx context.Context, keys ...string) error {
	return GetRedis().Del(ctx, keys...).Err()
}

// Exists 检查键是否存在
func Exists(ctx context.Context, keys ...string) (int64, error) {
	return GetRedis().Exists(ctx, keys...).Result()
}

// Expire 设置键过期时间
func Expire(ctx context.Context, key string, expiration time.Duration) error {
	return GetRedis().Expire(ctx, key, expiration).Err()
}

// Incr 自增
func Incr(ctx context.Context, key string) (int64, error) {
	return GetRedis().Incr(ctx, key).Result()
}

// Decr 自减
func Decr(ctx context.Context, key string) (int64, error) {
	return GetRedis().Decr(ctx, key).Result()
}

// HSet 设置哈希字段
func HSet(ctx context.Context, key string, values ...interface{}) error {
	return GetRedis().HSet(ctx, key, values...).Err()
}

// HGet 获取哈希字段
func HGet(ctx context.Context, key, field string) (string, error) {
	return GetRedis().HGet(ctx, key, field).Result()
}

// HGetAll 获取哈希所有字段
func HGetAll(ctx context.Context, key string) (map[string]string, error) {
	return GetRedis().HGetAll(ctx, key).Result()
}

// LPush 列表左侧推入
func LPush(ctx context.Context, key string, values ...interface{}) error {
	return GetRedis().LPush(ctx, key, values...).Err()
}

// RPush 列表右侧推入
func RPush(ctx context.Context, key string, values ...interface{}) error {
	return GetRedis().RPush(ctx, key, values...).Err()
}

// LPop 列表左侧弹出
func LPop(ctx context.Context, key string) (string, error) {
	return GetRedis().LPop(ctx, key).Result()
}

// RPop 列表右侧弹出
func RPop(ctx context.Context, key string) (string, error) {
	return GetRedis().RPop(ctx, key).Result()
}

// SAdd 集合添加元素
func SAdd(ctx context.Context, key string, members ...interface{}) error {
	return GetRedis().SAdd(ctx, key, members...).Err()
}

// SMembers 获取集合所有元素
func SMembers(ctx context.Context, key string) ([]string, error) {
	return GetRedis().SMembers(ctx, key).Result()
}

// ZAdd 有序集合添加元素
func ZAdd(ctx context.Context, key string, members ...redis.Z) error {
	return GetRedis().ZAdd(ctx, key, members...).Err()
}

// ZRange 有序集合范围查询
func ZRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	return GetRedis().ZRange(ctx, key, start, stop).Result()
}

// RedisInstance 全局Redis实例包装器
type RedisInstance struct{}

// Init 初始化Redis（从全局配置）
func (r *RedisInstance) Init() error {
	cfg := config.Get()
	if !cfg.Redis.Enabled {
		return fmt.Errorf("Redis未启用")
	}
	_, err := InitRedis(&cfg.Redis)
	return err
}

// AutoInit 自动初始化（只执行一次）
func (r *RedisInstance) AutoInit() error {
	redisOnce.Do(func() {
		cfg := config.Get()
		if cfg.Redis.Enabled {
			_, redisErr = InitRedis(&cfg.Redis)
		}
	})
	return redisErr
}

// Client 获取Redis客户端（自动初始化）
func (r *RedisInstance) Client() *redis.Client {
	if redisClient == nil {
		r.AutoInit()
	}
	return redisClient
}

// IsConnected 检查是否已连接
func (r *RedisInstance) IsConnected() bool {
	return redisClient != nil
}

// Close 关闭连接
func (r *RedisInstance) Close() error {
	return CloseRedis()
}

// ===== 便捷方法（自动初始化） =====

// Set 设置键值（自动初始化）
func (r *RedisInstance) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return Set(ctx, key, value, expiration)
}

// Get 获取键值（自动初始化）
func (r *RedisInstance) Get(ctx context.Context, key string) (string, error) {
	return Get(ctx, key)
}

// Del 删除键（自动初始化）
func (r *RedisInstance) Del(ctx context.Context, keys ...string) error {
	return Del(ctx, keys...)
}

// Exists 检查键是否存在（自动初始化）
func (r *RedisInstance) Exists(ctx context.Context, keys ...string) (int64, error) {
	return Exists(ctx, keys...)
}

// Expire 设置过期时间（自动初始化）
func (r *RedisInstance) Expire(ctx context.Context, key string, expiration time.Duration) error {
	return Expire(ctx, key, expiration)
}

// Incr 自增（自动初始化）
func (r *RedisInstance) Incr(ctx context.Context, key string) (int64, error) {
	return Incr(ctx, key)
}

// HSet 设置哈希字段（自动初始化）
func (r *RedisInstance) HSet(ctx context.Context, key string, values ...interface{}) error {
	return HSet(ctx, key, values...)
}

// HGet 获取哈希字段（自动初始化）
func (r *RedisInstance) HGet(ctx context.Context, key, field string) (string, error) {
	return HGet(ctx, key, field)
}

// HGetAll 获取哈希所有字段（自动初始化）
func (r *RedisInstance) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	return HGetAll(ctx, key)
}
