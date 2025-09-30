# Consul 服务注册与发现模块

本模块封装了 Consul 服务注册、注销和服务发现功能，支持自动健康检查和优雅关闭。

## 📦 功能特性

- ✅ 服务注册/注销
- ✅ 自动健康检查（HTTP）
- ✅ 服务发现
- ✅ 优雅关闭（自动注销服务）
- ✅ 自动获取本机IP

## 🚀 快速开始

### 方式1：一键式启动（推荐）

```go
package main

import (
    "log"
    "mule-cloud/core/cousul"
    "github.com/gin-gonic/gin"
)

func main() {
    r := gin.Default()
    
    // 注册健康检查路由（必须）
    r.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{"status": "ok"})
    })
    
    // 一键注册并启动服务
    err := cousul.RegisterAndRun(r, &cousul.ServiceConfig{
        ServiceName: "my-service",
        ServicePort: 8080,
        Tags:        []string{"api", "v1"},
    }, "127.0.0.1:8500")
    
    if err != nil {
        log.Fatalf("服务启动失败: %v", err)
    }
}
```

### 方式2：手动控制

```go
package main

import (
    "log"
    "os"
    "os/signal"
    "syscall"
    "mule-cloud/core/cousul"
    "github.com/gin-gonic/gin"
)

func main() {
    r := gin.Default()
    
    // 健康检查路由
    r.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{"status": "ok"})
    })
    
    // 创建Consul客户端
    consulClient, err := cousul.NewConsulClient("127.0.0.1:8500")
    if err != nil {
        log.Fatalf("连接Consul失败: %v", err)
    }
    
    // 注册服务
    err = consulClient.RegisterService(&cousul.ServiceConfig{
        ServiceName: "my-service",
        ServicePort: 8080,
        Tags:        []string{"api", "v1"},
    })
    if err != nil {
        log.Fatalf("服务注册失败: %v", err)
    }
    
    // 监听退出信号
    go func() {
        quit := make(chan os.Signal, 1)
        signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
        <-quit
        
        log.Println("正在注销服务...")
        consulClient.DeregisterService()
        os.Exit(0)
    }()
    
    // 启动服务
    r.Run(":8080")
}
```

## 📝 配置说明

### ServiceConfig 结构体

```go
type ServiceConfig struct {
    ServiceID      string   // 服务ID（可选，自动生成）
    ServiceName    string   // 服务名称（必填）
    ServiceAddress string   // 服务地址（可选，自动获取本机IP）
    ServicePort    int      // 服务端口（必填）
    Tags           []string // 服务标签（可选）
    HealthCheck    string   // 健康检查地址（可选，默认 /health）
}
```

### 配置示例

```go
&cousul.ServiceConfig{
    ServiceName:    "user-service",        // 服务名
    ServicePort:    8080,                  // 端口
    ServiceAddress: "192.168.1.100",       // 可选，不填自动获取
    Tags:           []string{"api", "v1"}, // 标签
    HealthCheck:    "http://192.168.1.100:8080/health", // 可选
}
```

## 🔍 服务发现

### 获取服务地址

```go
consulClient, _ := cousul.NewConsulClient("127.0.0.1:8500")

// 获取第一个健康的服务实例地址
address, err := consulClient.GetServiceAddress("user-service")
if err != nil {
    log.Printf("服务发现失败: %v", err)
}
fmt.Println(address) // 输出: 192.168.1.100:8080
```

### 获取所有服务实例

```go
consulClient, _ := cousul.NewConsulClient("127.0.0.1:8500")

// 获取所有健康的服务实例
services, err := consulClient.DiscoverService("user-service")
if err != nil {
    log.Printf("服务发现失败: %v", err)
}

for _, svc := range services {
    fmt.Printf("服务地址: %s:%d\n", 
        svc.Service.Address, 
        svc.Service.Port)
}
```

## ⚙️ Consul 配置

### 健康检查参数（默认）

- **检查间隔**: 5秒
- **超时时间**: 3秒
- **失败注销**: 30秒后自动注销不健康的服务

### Consul 连接地址

```go
// 本地开发
consulClient, _ := cousul.NewConsulClient("127.0.0.1:8500")

// 生产环境
consulClient, _ := cousul.NewConsulClient("consul.example.com:8500")
```

## 🧪 测试

### 启动本地 Consul（Docker）

```bash
docker run -d \
  --name consul \
  -p 8500:8500 \
  consul:latest agent -dev -ui -client=0.0.0.0
```

### 访问 Consul UI

```
http://localhost:8500/ui
```

### 查看注册的服务

```bash
# 命令行查看
curl http://localhost:8500/v1/agent/services

# 查看健康状态
curl http://localhost:8500/v1/health/service/my-service
```

## 📌 注意事项

1. **健康检查路由必须存在**
   - 默认路径: `/health`
   - 必须返回 HTTP 200 状态码

2. **端口冲突**
   - 确保 ServicePort 与实际启动端口一致
   - 同一台机器运行多个实例时，使用不同端口

3. **优雅关闭**
   - 使用 `RegisterAndRun` 会自动处理信号监听
   - 手动控制时需要 `defer consulClient.DeregisterService()`

4. **网络环境**
   - Consul 服务器必须可访问
   - 健康检查地址必须可从 Consul 访问

## 🔧 故障排查

### 1. 服务注册失败

```
错误: 创建Consul客户端失败
解决: 检查 Consul 是否运行，地址是否正确
```

### 2. 健康检查失败

```
错误: Service is in warning state
解决: 
- 确保 /health 路由存在
- 检查防火墙是否阻止访问
- 确认 ServiceAddress 正确
```

### 3. 获取本机IP失败

```
错误: 未找到有效的本机IP
解决: 手动指定 ServiceAddress
```

## 🌟 最佳实践

### 1. 使用环境变量

```go
consulAddr := os.Getenv("CONSUL_ADDR")
if consulAddr == "" {
    consulAddr = "127.0.0.1:8500" // 默认值
}
```

### 2. 服务标签规范

```go
Tags: []string{
    "version=v1.0.0",
    "env=production",
    "region=us-west",
}
```

### 3. 多实例部署

```go
// 使用不同端口启动多个实例
// 实例1: 8080
// 实例2: 8081
// Consul 会自动负载均衡
```

## 📚 相关文档

- [Consul 官方文档](https://www.consul.io/docs)
- [服务注册与发现原理](https://www.consul.io/docs/architecture)
- [健康检查配置](https://www.consul.io/docs/discovery/checks)
