# Consul 服务注册与发现集成指南

## 📖 概述

本项目已集成 Consul 服务注册与发现功能，支持：
- ✅ 自动服务注册/注销
- ✅ 健康检查（HTTP）
- ✅ 服务发现
- ✅ 优雅关闭
- ✅ 多实例支持

## 🚀 快速开始

### 1. 启动 Consul 服务器

**使用 Docker（推荐）:**

```bash
docker run -d \
  --name consul \
  -p 8500:8500 \
  consul:latest agent -dev -ui -client=0.0.0.0
```

**访问 Consul UI:**
```
http://localhost:8500/ui
```

### 2. 运行示例项目

```bash
# 进入项目目录
cd k:\Git\mule-cloud

# 运行测试服务
go run test/cmd/main.go
```

### 3. 验证服务注册

**方式1：通过 Consul UI**
```
访问: http://localhost:8500/ui/dc1/services
查看是否有 "testserver" 服务
```

**方式2：通过 API**
```bash
curl http://localhost:8500/v1/agent/services | json_pp
```

**方式3：测试健康检查**
```bash
curl http://localhost:8080/common/health
```

## 📝 使用说明

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
	
	// 必须有健康检查路由
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
	
	// 一键启动
	err := cousul.RegisterAndRun(r, &cousul.ServiceConfig{
		ServiceName: "my-service",
		ServicePort: 8080,
		Tags:        []string{"api", "v1"},
	}, "127.0.0.1:8500")
	
	if err != nil {
		log.Fatalf("启动失败: %v", err)
	}
}
```

### 方式2：自定义健康检查

```go
err := cousul.RegisterAndRun(r, &cousul.ServiceConfig{
	ServiceName: "my-service",
	ServicePort: 8080,
	Tags:        []string{"api", "v1"},
	HealthCheck: &cousul.HealthCheck{
		HTTP:     "http://localhost:8080/custom/health",
		Interval: "10s",  // 每10秒检查
		Timeout:  "5s",   // 5秒超时
		DeregisterCriticalServiceAfter: "30s",  // 30秒后注销
	},
}, "127.0.0.1:8500")
```

### 方式3：手动控制

```go
// 创建客户端
consulClient, err := cousul.NewConsulClient("127.0.0.1:8500")
if err != nil {
	log.Fatal(err)
}

// 注册服务
err = consulClient.RegisterService(&cousul.ServiceConfig{
	ServiceName: "my-service",
	ServicePort: 8080,
	Tags:        []string{"api", "v1"},
})

// 程序退出时注销
defer consulClient.DeregisterService()

// 启动服务
r.Run(":8080")
```

## 🔍 服务发现

### 获取服务地址

```go
consulClient, _ := cousul.NewConsulClient("127.0.0.1:8500")

// 获取第一个可用实例
address, err := consulClient.GetServiceAddress("user-service")
if err != nil {
	log.Fatal(err)
}

// 使用服务地址
resp, _ := http.Get(fmt.Sprintf("http://%s/api/users", address))
```

### 获取所有实例

```go
services, err := consulClient.DiscoverService("user-service")
if err != nil {
	log.Fatal(err)
}

for _, svc := range services {
	fmt.Printf("实例: %s:%d\n", svc.Service.Address, svc.Service.Port)
}
```

## ⚙️ 配置说明

### ServiceConfig 参数

| 参数 | 类型 | 必填 | 说明 | 默认值 |
|-----|------|------|------|--------|
| ServiceName | string | ✅ | 服务名称 | - |
| ServicePort | int | ✅ | 服务端口 | - |
| ServiceID | string | ❌ | 服务ID | 自动生成 |
| ServiceAddress | string | ❌ | 服务地址 | 自动获取本机IP |
| Tags | []string | ❌ | 服务标签 | nil |
| HealthCheck | *HealthCheck | ❌ | 健康检查配置 | 默认配置 |

### HealthCheck 参数

| 参数 | 类型 | 说明 | 默认值 |
|-----|------|------|--------|
| HTTP | string | 健康检查URL | `http://{ip}:{port}/health` |
| Interval | string | 检查间隔 | "5s" |
| Timeout | string | 超时时间 | "3s" |
| DeregisterCriticalServiceAfter | string | 失败注销时间 | "30s" |

## 🎯 实际应用场景

### 场景1：微服务架构

```
服务A (8080) ──注册到──> Consul
服务B (8081) ──注册到──> Consul
服务C (8082) ──注册到──> Consul

API网关 ──查询──> Consul ──返回──> 服务列表
```

### 场景2：多实例负载均衡

```bash
# 启动3个实例
go run test/cmd/main.go --port 8080
go run test/cmd/main.go --port 8081
go run test/cmd/main.go --port 8082

# Consul 会自动负载均衡
```

### 场景3：生产环境部署

```bash
# 设置环境变量
export CONSUL_ADDR="consul.production.com:8500"
export SERVICE_NAME="order-service"
export SERVICE_PORT="8080"
export ENVIRONMENT="production"

# 启动服务
go run main.go
```

## 📋 健康检查说明

### 要求

1. **必须实现健康检查路由**
   ```go
   r.GET("/health", func(c *gin.Context) {
       c.JSON(200, gin.H{"status": "ok"})
   })
   ```

2. **返回 HTTP 200 状态码**
   - Consul 只认为 200 是健康状态
   - 其他状态码会标记为不健康

3. **可访问性**
   - 健康检查 URL 必须能从 Consul 服务器访问
   - 注意防火墙和网络配置

### 检查机制

```
时间线:
├─ 0s:  服务启动，注册到 Consul
├─ 5s:  第一次健康检查（成功）
├─ 10s: 第二次健康检查（成功）
├─ 15s: 第三次健康检查（失败）-> 标记为 warning
├─ 45s: 连续失败30秒 -> 自动注销服务
```

## 🔧 故障排查

### 1. 服务注册失败

**错误**: `创建Consul客户端失败`

**解决**:
```bash
# 检查 Consul 是否运行
curl http://localhost:8500/v1/status/leader

# 检查地址是否正确
telnet localhost 8500
```

### 2. 健康检查失败

**错误**: Service is in warning state

**解决**:
```bash
# 检查健康检查 URL 是否可访问
curl http://localhost:8080/common/health

# 查看 Consul 日志
docker logs consul
```

### 3. 服务未注销

**问题**: 程序退出后服务仍显示在 Consul 中

**解决**:
- 使用 `RegisterAndRun` 方法（自动处理）
- 或者使用 `defer consulClient.DeregisterService()`

### 4. 获取本机IP失败

**错误**: `未找到有效的本机IP`

**解决**:
```go
// 手动指定 ServiceAddress
&cousul.ServiceConfig{
	ServiceAddress: "192.168.1.100",
	// ...
}
```

## 📊 测试命令

### 查看所有服务

```bash
curl http://localhost:8500/v1/agent/services | json_pp
```

### 查看服务健康状态

```bash
curl http://localhost:8500/v1/health/service/testserver | json_pp
```

### 测试健康检查

```bash
curl http://localhost:8080/common/health
```

### 查看服务详情

```bash
curl http://localhost:8500/v1/catalog/service/testserver | json_pp
```

## 🌟 最佳实践

### 1. 使用环境变量

```go
consulAddr := os.Getenv("CONSUL_ADDR")
if consulAddr == "" {
	consulAddr = "127.0.0.1:8500"
}
```

### 2. 规范化服务标签

```go
Tags: []string{
	"version=v1.0.0",
	"env=production",
	"region=us-west",
	"team=backend",
}
```

### 3. 合理设置健康检查

```go
HealthCheck: &cousul.HealthCheck{
	Interval: "10s",  // 开发环境: 5s, 生产环境: 10-30s
	Timeout:  "5s",   // 应小于 Interval
	DeregisterCriticalServiceAfter: "1m",  // 生产环境建议 1-5m
}
```

### 4. 使用有意义的服务名

```go
// ✅ 好的命名
ServiceName: "order-service"
ServiceName: "user-api"
ServiceName: "payment-gateway"

// ❌ 避免
ServiceName: "service1"
ServiceName: "test"
```

## 📚 参考资料

- [Consul 官方文档](https://www.consul.io/docs)
- [Consul API 文档](https://www.consul.io/api-docs)
- [项目代码示例](../core/cousul/example_usage.go)
- [完整 API 文档](../core/cousul/README.md)

## 📞 技术支持

遇到问题？
1. 查看 [故障排查](#-故障排查) 章节
2. 查看 [Consul UI](http://localhost:8500/ui) 日志
3. 查看代码示例 `core/cousul/example_usage.go`
