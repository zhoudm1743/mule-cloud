package middleware

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/afex/hystrix-go/hystrix"
	"github.com/hashicorp/consul/api"
)

const (
	// Consul KV 路径前缀
	RouteConfigPrefix   = "gateway/routes/"
	HystrixConfigPrefix = "gateway/hystrix/"
)

// DynamicRouteManager 动态路由管理器
type DynamicRouteManager struct {
	consulClient *api.Client
	routes       map[string]*RouteConfig
	routeLock    sync.RWMutex

	hystrixConfigs map[string]*DynamicHystrixConfig
	hystrixLock    sync.RWMutex

	stopChan chan bool
}

// RouteConfig 路由配置
type RouteConfig struct {
	ServiceName   string   `json:"service_name"`
	GatewayPrefix string   `json:"gateway_prefix"` // 网关前缀（如 /admin），转发时会去掉
	RequireAuth   bool     `json:"require_auth"`
	RequireRole   []string `json:"require_role"`
}

// DynamicHystrixConfig Hystrix配置（用于动态路由管理）
type DynamicHystrixConfig struct {
	Timeout                int `json:"timeout"`
	MaxConcurrentRequests  int `json:"max_concurrent_requests"`
	RequestVolumeThreshold int `json:"request_volume_threshold"`
	SleepWindow            int `json:"sleep_window"`
	ErrorPercentThreshold  int `json:"error_percent_threshold"`
}

// NewDynamicRouteManager 创建动态路由管理器
func NewDynamicRouteManager(consulClient *api.Client) *DynamicRouteManager {
	manager := &DynamicRouteManager{
		consulClient:   consulClient,
		routes:         make(map[string]*RouteConfig),
		hystrixConfigs: make(map[string]*DynamicHystrixConfig),
		stopChan:       make(chan bool),
	}

	// 初始加载配置
	if err := manager.loadRoutes(); err != nil {
		log.Printf("⚠️  初始加载路由配置失败: %v", err)
	}
	if err := manager.loadHystrixConfigs(); err != nil {
		log.Printf("⚠️  初始加载Hystrix配置失败: %v", err)
	}

	// 启动配置监听
	go manager.watchConfigs()

	return manager
}

// GetRoute 获取路由配置
func (m *DynamicRouteManager) GetRoute(prefix string) (*RouteConfig, bool) {
	m.routeLock.RLock()
	defer m.routeLock.RUnlock()

	route, exists := m.routes[prefix]
	return route, exists
}

// GetAllRoutes 获取所有路由配置
func (m *DynamicRouteManager) GetAllRoutes() map[string]*RouteConfig {
	m.routeLock.RLock()
	defer m.routeLock.RUnlock()

	// 返回副本，避免外部修改
	routes := make(map[string]*RouteConfig, len(m.routes))
	for k, v := range m.routes {
		routes[k] = v
	}
	return routes
}

// AddRoute 添加路由配置
func (m *DynamicRouteManager) AddRoute(prefix string, config *RouteConfig) error {
	// 保存到 Consul
	key := RouteConfigPrefix + prefix
	data, err := json.Marshal(config)
	if err != nil {
		return fmt.Errorf("序列化路由配置失败: %v", err)
	}

	kv := &api.KVPair{
		Key:   key,
		Value: data,
	}

	if _, err := m.consulClient.KV().Put(kv, nil); err != nil {
		return fmt.Errorf("保存路由配置到Consul失败: %v", err)
	}

	// 更新本地缓存
	m.routeLock.Lock()
	m.routes[prefix] = config
	m.routeLock.Unlock()

	return nil
}

// UpdateRoute 更新路由配置
func (m *DynamicRouteManager) UpdateRoute(prefix string, config *RouteConfig) error {
	return m.AddRoute(prefix, config) // 实现相同
}

// DeleteRoute 删除路由配置
func (m *DynamicRouteManager) DeleteRoute(prefix string) error {
	// 从 Consul 删除
	key := RouteConfigPrefix + prefix
	if _, err := m.consulClient.KV().Delete(key, nil); err != nil {
		return fmt.Errorf("从Consul删除路由配置失败: %v", err)
	}

	// 从本地缓存删除
	m.routeLock.Lock()
	delete(m.routes, prefix)
	m.routeLock.Unlock()

	return nil
}

// GetHystrixConfig 获取Hystrix配置
func (m *DynamicRouteManager) GetHystrixConfig(serviceName string) (*DynamicHystrixConfig, bool) {
	m.hystrixLock.RLock()
	defer m.hystrixLock.RUnlock()

	config, exists := m.hystrixConfigs[serviceName]
	return config, exists
}

// GetAllHystrixConfigs 获取所有Hystrix配置
func (m *DynamicRouteManager) GetAllHystrixConfigs() map[string]*DynamicHystrixConfig {
	m.hystrixLock.RLock()
	defer m.hystrixLock.RUnlock()

	configs := make(map[string]*DynamicHystrixConfig, len(m.hystrixConfigs))
	for k, v := range m.hystrixConfigs {
		configs[k] = v
	}
	return configs
}

// AddHystrixConfig 添加Hystrix配置
func (m *DynamicRouteManager) AddHystrixConfig(serviceName string, config *DynamicHystrixConfig) error {
	// 保存到 Consul
	key := HystrixConfigPrefix + serviceName
	data, err := json.Marshal(config)
	if err != nil {
		return fmt.Errorf("序列化Hystrix配置失败: %v", err)
	}

	kv := &api.KVPair{
		Key:   key,
		Value: data,
	}

	if _, err := m.consulClient.KV().Put(kv, nil); err != nil {
		return fmt.Errorf("保存Hystrix配置到Consul失败: %v", err)
	}

	// 更新本地缓存
	m.hystrixLock.Lock()
	m.hystrixConfigs[serviceName] = config
	m.hystrixLock.Unlock()

	// 动态更新Hystrix配置
	m.applyHystrixConfig(serviceName, config)

	return nil
}

// DeleteHystrixConfig 删除Hystrix配置
func (m *DynamicRouteManager) DeleteHystrixConfig(serviceName string) error {
	// 从 Consul 删除
	key := HystrixConfigPrefix + serviceName
	if _, err := m.consulClient.KV().Delete(key, nil); err != nil {
		return fmt.Errorf("从Consul删除Hystrix配置失败: %v", err)
	}

	// 从本地缓存删除
	m.hystrixLock.Lock()
	delete(m.hystrixConfigs, serviceName)
	m.hystrixLock.Unlock()

	return nil
}

// loadRoutes 从Consul加载所有路由配置
func (m *DynamicRouteManager) loadRoutes() error {
	pairs, _, err := m.consulClient.KV().List(RouteConfigPrefix, nil)
	if err != nil {
		return fmt.Errorf("从Consul获取路由配置失败: %v", err)
	}

	m.routeLock.Lock()
	defer m.routeLock.Unlock()

	// 清空现有配置
	m.routes = make(map[string]*RouteConfig)

	for _, pair := range pairs {
		var config RouteConfig
		if err := json.Unmarshal(pair.Value, &config); err != nil {
			log.Printf("⚠️  解析路由配置失败 (key=%s): %v", pair.Key, err)
			continue
		}

		// 提取路由前缀（移除前缀路径）
		prefix := pair.Key[len(RouteConfigPrefix):]
		// 确保前缀有前导斜杠
		if !strings.HasPrefix(prefix, "/") {
			prefix = "/" + prefix
		}
		m.routes[prefix] = &config

	}

	return nil
}

// loadHystrixConfigs 从Consul加载所有Hystrix配置
func (m *DynamicRouteManager) loadHystrixConfigs() error {
	pairs, _, err := m.consulClient.KV().List(HystrixConfigPrefix, nil)
	if err != nil {
		return fmt.Errorf("从Consul获取Hystrix配置失败: %v", err)
	}

	m.hystrixLock.Lock()
	defer m.hystrixLock.Unlock()

	// 清空现有配置
	m.hystrixConfigs = make(map[string]*DynamicHystrixConfig)

	for _, pair := range pairs {
		var config DynamicHystrixConfig
		if err := json.Unmarshal(pair.Value, &config); err != nil {
			log.Printf("⚠️  解析Hystrix配置失败 (key=%s): %v", pair.Key, err)
			continue
		}

		// 提取服务名（移除前缀路径）
		serviceName := pair.Key[len(HystrixConfigPrefix):]
		m.hystrixConfigs[serviceName] = &config

		// 应用Hystrix配置
		m.applyHystrixConfig(serviceName, &config)
	}

	return nil
}

// applyHystrixConfig 应用Hystrix配置
func (m *DynamicRouteManager) applyHystrixConfig(serviceName string, config *DynamicHystrixConfig) {
	// 动态配置Hystrix
	hystrix.ConfigureCommand(serviceName, hystrix.CommandConfig{
		Timeout:                config.Timeout,
		MaxConcurrentRequests:  config.MaxConcurrentRequests,
		RequestVolumeThreshold: config.RequestVolumeThreshold,
		SleepWindow:            config.SleepWindow,
		ErrorPercentThreshold:  config.ErrorPercentThreshold,
	})

}

// watchConfigs 监听配置变化
func (m *DynamicRouteManager) watchConfigs() {
	ticker := time.NewTicker(10 * time.Second) // 每10秒检查一次
	defer ticker.Stop()

	log.Println("🔍 启动配置监听器...")

	for {
		select {
		case <-ticker.C:
			// 重新加载路由配置
			if err := m.loadRoutes(); err != nil {
				log.Printf("⚠️  重新加载路由配置失败: %v", err)
			}

			// 重新加载Hystrix配置
			if err := m.loadHystrixConfigs(); err != nil {
				log.Printf("⚠️  重新加载Hystrix配置失败: %v", err)
			}

		case <-m.stopChan:
			log.Println("🛑 停止配置监听器")
			return
		}
	}
}

// Stop 停止管理器
func (m *DynamicRouteManager) Stop() {
	close(m.stopChan)
}
