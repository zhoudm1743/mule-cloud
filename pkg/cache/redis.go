package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

// Config Redis配置
type Config struct {
	Host     string
	Port     int
	Password string
	DB       int
}

// Cache 缓存接口
type Cache interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	Del(ctx context.Context, keys ...string) error
	Exists(ctx context.Context, key string) (bool, error)
	HGet(ctx context.Context, key, field string) (string, error)
	HSet(ctx context.Context, key string, values ...interface{}) error
	HDel(ctx context.Context, key string, fields ...string) error
	HGetAll(ctx context.Context, key string) (map[string]string, error)
	Incr(ctx context.Context, key string) (int64, error)
	Expire(ctx context.Context, key string, expiration time.Duration) error
	Lock(ctx context.Context, key string, expiration time.Duration) (*DistributedLock, error)
}

// RedisCache Redis缓存实现
type RedisCache struct {
	client *redis.Client
}

// NewRedisCache 创建Redis缓存客户端
func NewRedisCache(config Config) (Cache, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Host, config.Port),
		Password: config.Password,
		DB:       config.DB,
	})

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &RedisCache{client: rdb}, nil
}

// Get 获取值
func (r *RedisCache) Get(ctx context.Context, key string) (string, error) {
	result, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("key %s not found", key)
	}
	return result, err
}

// Set 设置值
func (r *RedisCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	var data string
	switch v := value.(type) {
	case string:
		data = v
	default:
		bytes, err := json.Marshal(value)
		if err != nil {
			return fmt.Errorf("failed to marshal value: %w", err)
		}
		data = string(bytes)
	}

	return r.client.Set(ctx, key, data, ttl).Err()
}

// Del 删除键
func (r *RedisCache) Del(ctx context.Context, keys ...string) error {
	return r.client.Del(ctx, keys...).Err()
}

// Exists 检查键是否存在
func (r *RedisCache) Exists(ctx context.Context, key string) (bool, error) {
	result, err := r.client.Exists(ctx, key).Result()
	return result > 0, err
}

// HGet 获取哈希字段值
func (r *RedisCache) HGet(ctx context.Context, key, field string) (string, error) {
	return r.client.HGet(ctx, key, field).Result()
}

// HSet 设置哈希字段值
func (r *RedisCache) HSet(ctx context.Context, key string, values ...interface{}) error {
	return r.client.HSet(ctx, key, values...).Err()
}

// HDel 删除哈希字段
func (r *RedisCache) HDel(ctx context.Context, key string, fields ...string) error {
	return r.client.HDel(ctx, key, fields...).Err()
}

// HGetAll 获取哈希所有字段
func (r *RedisCache) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	return r.client.HGetAll(ctx, key).Result()
}

// Incr 自增
func (r *RedisCache) Incr(ctx context.Context, key string) (int64, error) {
	return r.client.Incr(ctx, key).Result()
}

// Expire 设置过期时间
func (r *RedisCache) Expire(ctx context.Context, key string, expiration time.Duration) error {
	return r.client.Expire(ctx, key, expiration).Err()
}

// Lock 获取分布式锁
func (r *RedisCache) Lock(ctx context.Context, key string, expiration time.Duration) (*DistributedLock, error) {
	lockValue := uuid.New().String()

	lock := &DistributedLock{
		client: r.client,
		key:    key,
		value:  lockValue,
		expiry: expiration,
	}

	err := lock.Lock(ctx)
	if err != nil {
		return nil, err
	}

	return lock, nil
}

// DistributedLock 分布式锁
type DistributedLock struct {
	client *redis.Client
	key    string
	value  string
	expiry time.Duration
}

// Lock 获取锁
func (l *DistributedLock) Lock(ctx context.Context) error {
	script := `
		if redis.call("GET", KEYS[1]) == ARGV[1] then
			return redis.call("PEXPIRE", KEYS[1], ARGV[2])
		else
			return redis.call("SET", KEYS[1], ARGV[1], "PX", ARGV[2], "NX")
		end
	`

	result := l.client.Eval(ctx, script, []string{l.key}, l.value, l.expiry.Milliseconds())
	if result.Err() != nil {
		return fmt.Errorf("failed to acquire lock: %w", result.Err())
	}

	return nil
}

// Unlock 释放锁
func (l *DistributedLock) Unlock(ctx context.Context) error {
	script := `
		if redis.call("GET", KEYS[1]) == ARGV[1] then
			return redis.call("DEL", KEYS[1])
		else
			return 0
		end
	`

	result := l.client.Eval(ctx, script, []string{l.key}, l.value)
	return result.Err()
}

// CacheManager 缓存管理器
type CacheManager struct {
	cache Cache
	ttl   map[string]time.Duration
}

// NewCacheManager 创建缓存管理器
func NewCacheManager(cache Cache) *CacheManager {
	return &CacheManager{
		cache: cache,
		ttl: map[string]time.Duration{
			"user:":     24 * time.Hour,   // 用户信息缓存24小时
			"session:":  24 * time.Hour,   // 会话信息缓存24小时
			"perm:":     1 * time.Hour,    // 权限信息缓存1小时
			"order:":    30 * time.Minute, // 订单信息缓存30分钟
			"style:":    1 * time.Hour,    // 款式信息缓存1小时
			"customer:": 1 * time.Hour,    // 客户信息缓存1小时
			"stats:":    5 * time.Minute,  // 统计数据缓存5分钟
			"report:":   15 * time.Minute, // 报表数据缓存15分钟
			"config:":   1 * time.Hour,    // 配置数据缓存1小时
			"dict:":     1 * time.Hour,    // 字典数据缓存1小时
		},
	}
}

// SetWithDefaultTTL 使用默认TTL设置缓存
func (cm *CacheManager) SetWithDefaultTTL(ctx context.Context, key string, value interface{}) error {
	ttl := cm.getDefaultTTL(key)
	return cm.cache.Set(ctx, key, value, ttl)
}

// getDefaultTTL 获取默认TTL
func (cm *CacheManager) getDefaultTTL(key string) time.Duration {
	for prefix, ttl := range cm.ttl {
		if len(key) >= len(prefix) && key[:len(prefix)] == prefix {
			return ttl
		}
	}
	return 15 * time.Minute // 默认15分钟
}

// Get 获取缓存
func (cm *CacheManager) Get(ctx context.Context, key string) (string, error) {
	return cm.cache.Get(ctx, key)
}

// Set 设置缓存
func (cm *CacheManager) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	return cm.cache.Set(ctx, key, value, ttl)
}

// Del 删除缓存
func (cm *CacheManager) Del(ctx context.Context, keys ...string) error {
	return cm.cache.Del(ctx, keys...)
}
