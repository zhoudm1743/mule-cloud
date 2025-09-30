package grpc

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/consul/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// ClientManager gRPC客户端管理器
type ClientManager struct {
	consulClient *api.Client
	connections  map[string]*grpc.ClientConn
}

// NewClientManager 创建客户端管理器
func NewClientManager(consulAddr string) (*ClientManager, error) {
	config := api.DefaultConfig()
	config.Address = consulAddr

	client, err := api.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("连接Consul失败: %v", err)
	}

	return &ClientManager{
		consulClient: client,
		connections:  make(map[string]*grpc.ClientConn),
	}, nil
}

// GetConnection 获取服务的gRPC连接
func (cm *ClientManager) GetConnection(serviceName string) (*grpc.ClientConn, error) {
	// 检查是否已有连接
	if conn, exists := cm.connections[serviceName]; exists {
		if conn.GetState().String() != "SHUTDOWN" {
			return conn, nil
		}
		// 连接已关闭，删除
		delete(cm.connections, serviceName)
	}

	// 从Consul获取服务地址
	addr, err := cm.getServiceAddress(serviceName)
	if err != nil {
		return nil, err
	}

	// 创建新连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, fmt.Errorf("连接gRPC服务失败: %v", err)
	}

	cm.connections[serviceName] = conn
	log.Printf("[gRPC] 已连接到服务: %s (%s)", serviceName, addr)

	return conn, nil
}

// getServiceAddress 从Consul获取服务地址
func (cm *ClientManager) getServiceAddress(serviceName string) (string, error) {
	services, _, err := cm.consulClient.Health().Service(serviceName, "", true, nil)
	if err != nil {
		return "", fmt.Errorf("查询服务失败: %v", err)
	}

	if len(services) == 0 {
		return "", fmt.Errorf("未找到服务: %s", serviceName)
	}

	// 简单负载均衡：返回第一个健康实例
	service := services[0].Service
	return fmt.Sprintf("%s:%d", service.Address, service.Port), nil
}

// Close 关闭所有连接
func (cm *ClientManager) Close() {
	for name, conn := range cm.connections {
		if err := conn.Close(); err != nil {
			log.Printf("[gRPC] 关闭连接失败: %s, 错误: %v", name, err)
		} else {
			log.Printf("[gRPC] 已关闭连接: %s", name)
		}
	}
}

// CallWithRetry 带重试的调用（通用函数）
func CallWithRetry(ctx context.Context, maxRetries int, fn func() error) error {
	var err error
	for i := 0; i < maxRetries; i++ {
		err = fn()
		if err == nil {
			return nil
		}

		if i < maxRetries-1 {
			time.Sleep(time.Duration(i+1) * 100 * time.Millisecond)
			log.Printf("[gRPC] 重试 %d/%d: %v", i+1, maxRetries, err)
		}
	}
	return fmt.Errorf("重试%d次后失败: %v", maxRetries, err)
}
