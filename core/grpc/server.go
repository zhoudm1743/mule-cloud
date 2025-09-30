package grpc

import (
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// ServerConfig gRPC服务器配置
type ServerConfig struct {
	Port int
}

// NewServer 创建gRPC服务器
func NewServer(config *ServerConfig) *grpc.Server {
	server := grpc.NewServer(
		// 可以添加拦截器
		grpc.UnaryInterceptor(loggingInterceptor),
	)

	// 注册反射服务（用于grpcurl等工具）
	reflection.Register(server)

	return server
}

// StartServer 启动gRPC服务器
func StartServer(server *grpc.Server, port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return fmt.Errorf("监听端口失败: %v", err)
	}

	log.Printf("[gRPC] 服务器启动: :%d", port)
	return server.Serve(lis)
}

// loggingInterceptor 日志拦截器
func loggingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	log.Printf("[gRPC] 请求: %s", info.FullMethod)
	resp, err := handler(ctx, req)
	if err != nil {
		log.Printf("[gRPC] 错误: %s - %v", info.FullMethod, err)
	}
	return resp, err
}
