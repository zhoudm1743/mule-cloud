package cousul

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/hashicorp/consul/api"
)

// ConsulClient Consul客户端封装
type ConsulClient struct {
	client *api.Client
	config *ServiceConfig
}

// ServiceConfig 服务配置
type ServiceConfig struct {
	ServiceID      string       // 服务ID (唯一标识)
	ServiceName    string       // 服务名称
	ServiceAddress string       // 服务地址
	ServicePort    int          // 服务端口
	Tags           []string     // 服务标签
	HealthCheck    *HealthCheck // 健康检查配置（可选）
}

// HealthCheck 健康检查配置
type HealthCheck struct {
	HTTP                           string // HTTP健康检查地址
	Interval                       string // 检查间隔（如 "5s", "10s"）
	Timeout                        string // 超时时间（如 "3s", "5s"）
	DeregisterCriticalServiceAfter string // 失败后多久注销服务（如 "30s"）
}

// NewConsulClient 创建Consul客户端
func NewConsulClient(consulAddress string) (*ConsulClient, error) {
	config := api.DefaultConfig()
	config.Address = consulAddress

	client, err := api.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("创建Consul客户端失败: %v", err)
	}

	return &ConsulClient{
		client: client,
	}, nil
}

// RegisterService 注册服务到Consul
func (c *ConsulClient) RegisterService(cfg *ServiceConfig) error {
	// 获取本机IP
	if cfg.ServiceAddress == "" {
		ip, err := getLocalIP()
		if err != nil {
			return fmt.Errorf("获取本机IP失败: %v", err)
		}
		cfg.ServiceAddress = ip
	}

	// 默认服务ID
	if cfg.ServiceID == "" {
		cfg.ServiceID = fmt.Sprintf("%s-%s-%d", cfg.ServiceName, cfg.ServiceAddress, cfg.ServicePort)
	}

	// 默认健康检查配置
	if cfg.HealthCheck == nil {
		cfg.HealthCheck = &HealthCheck{
			HTTP:                           fmt.Sprintf("http://%s:%d/health", cfg.ServiceAddress, cfg.ServicePort),
			Interval:                       "5s",
			Timeout:                        "3s",
			DeregisterCriticalServiceAfter: "30s",
		}
	} else {
		// 填充默认值
		if cfg.HealthCheck.HTTP == "" {
			cfg.HealthCheck.HTTP = fmt.Sprintf("http://%s:%d/health", cfg.ServiceAddress, cfg.ServicePort)
		}
		if cfg.HealthCheck.Interval == "" {
			cfg.HealthCheck.Interval = "5s"
		}
		if cfg.HealthCheck.Timeout == "" {
			cfg.HealthCheck.Timeout = "3s"
		}
		if cfg.HealthCheck.DeregisterCriticalServiceAfter == "" {
			cfg.HealthCheck.DeregisterCriticalServiceAfter = "30s"
		}
	}

	c.config = cfg

	// 构建服务注册信息
	registration := &api.AgentServiceRegistration{
		ID:      cfg.ServiceID,
		Name:    cfg.ServiceName,
		Address: cfg.ServiceAddress,
		Port:    cfg.ServicePort,
		Tags:    cfg.Tags,
		Check: &api.AgentServiceCheck{
			HTTP:                           cfg.HealthCheck.HTTP,
			Timeout:                        cfg.HealthCheck.Timeout,
			Interval:                       cfg.HealthCheck.Interval,
			DeregisterCriticalServiceAfter: cfg.HealthCheck.DeregisterCriticalServiceAfter,
		},
	}

	// 注册服务
	err := c.client.Agent().ServiceRegister(registration)
	if err != nil {
		return fmt.Errorf("服务注册失败: %v", err)
	}

	log.Printf("[Consul] 服务注册成功 -> ID: %s, Name: %s, Address: %s:%d",
		cfg.ServiceID, cfg.ServiceName, cfg.ServiceAddress, cfg.ServicePort)
	return nil
}

// DeregisterService 注销服务
func (c *ConsulClient) DeregisterService() error {
	if c.config == nil {
		return fmt.Errorf("服务未注册")
	}

	err := c.client.Agent().ServiceDeregister(c.config.ServiceID)
	if err != nil {
		return fmt.Errorf("服务注销失败: %v", err)
	}

	log.Printf("[Consul] 服务注销成功 -> ID: %s", c.config.ServiceID)
	return nil
}

// DiscoverService 服务发现 - 根据服务名查找健康的服务实例
func (c *ConsulClient) DiscoverService(serviceName string) ([]*api.ServiceEntry, error) {
	services, _, err := c.client.Health().Service(serviceName, "", true, nil)
	if err != nil {
		return nil, fmt.Errorf("服务发现失败: %v", err)
	}

	if len(services) == 0 {
		return nil, fmt.Errorf("未找到服务: %s", serviceName)
	}

	return services, nil
}

// GetServiceAddress 获取服务地址（简化版，返回第一个健康实例）
func (c *ConsulClient) GetServiceAddress(serviceName string) (string, error) {
	services, err := c.DiscoverService(serviceName)
	if err != nil {
		return "", err
	}

	service := services[0].Service
	return fmt.Sprintf("%s:%d", service.Address, service.Port), nil
}

// getLocalIP 获取本机非回环IP地址
func getLocalIP() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				return ipNet.IP.String(), nil
			}
		}
	}

	return "", fmt.Errorf("未找到有效的本机IP")
}

// RegisterAndRun 注册服务到Consul并启动HTTP服务
// 包含优雅关闭功能（监听 SIGINT 和 SIGTERM 信号）
// 示例用法：
//
//	cousul.RegisterAndRun(r, &cousul.ServiceConfig{
//	    ServiceName:    "my-service",
//	    ServicePort:    8080,
//	    Tags:           []string{"api", "v1"},
//	}, "127.0.0.1:8500")
func RegisterAndRun(router *gin.Engine, config *ServiceConfig, consulAddress string) error {
	// 创建Consul客户端
	consulClient, err := NewConsulClient(consulAddress)
	if err != nil {
		log.Fatalf("连接Consul失败: %v", err)
		return err
	}

	// 注册服务
	err = consulClient.RegisterService(config)
	if err != nil {
		log.Fatalf("服务注册失败: %v", err)
		return err
	}

	// 监听退出信号，优雅关闭
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit

		log.Println("正在关闭服务...")
		if err := consulClient.DeregisterService(); err != nil {
			log.Printf("服务注销失败: %v", err)
		}
		os.Exit(0)
	}()

	// 启动HTTP服务
	addr := fmt.Sprintf(":%d", config.ServicePort)
	return router.Run(addr)
}
