# Mule-Cloud API网关

## 功能特性

- ✅ 统一入口：所有微服务通过一个端口访问
- ✅ 服务发现：自动从Consul获取服务地址
- ✅ 反向代理：动态转发请求到后端服务
- ✅ 路由管理：配置化的路由规则
- ✅ 健康检查：监控网关和后端服务状态
- ✅ 请求日志：记录所有转发请求和响应时间

## 架构图

```
客户端请求
    ↓
http://localhost:8080/test/admin/123
    ↓
┌─────────────────────┐
│   API网关 (:8080)   │
│  - 路由匹配         │
│  - 服务发现         │
│  - 反向代理         │
└──────────┬──────────┘
           │
           ├─→ Consul服务发现
           │
           ├─→ testservice (:8000)
           └─→ basicservice (:8001)
```

## 快速启动

### 前置条件

1. Consul已启动（默认: `127.0.0.1:8500`）
2. 后端服务已注册到Consul
   - `testservice` (端口: 8000)
   - `basicservice` (端口: 8001)

### 启动网关

```bash
cd gateway
go run main.go
```

输出示例：
```
========================================
🚀 Mule-Cloud API网关启动成功
📍 监听端口: :8080
🔗 Consul地址: 127.0.0.1:8500
📋 路由配置:
   /test/* → testservice (Consul)
   /basic/* → basicservice (Consul)
========================================
```

## 使用示例

### 通过网关访问服务

```bash
# 访问 test 服务
curl http://localhost:8080/test/admin/123

# 访问 basic 服务
curl http://localhost:8080/basic/color/1
```

### 网关管理接口

```bash
# 健康检查
curl http://localhost:8080/gateway/health

# 查看路由配置
curl http://localhost:8080/gateway/routes
```

## 路由配置

修改 `main.go` 中的路由映射：

```go
routes: map[string]string{
    "/test":  "testservice",   // http://localhost:8080/test/* → testservice
    "/basic": "basicservice",  // http://localhost:8080/basic/* → basicservice
    "/order": "orderservice",  // 添加新服务
}
```

## 对比：有无网关的区别

### 没有网关（原来的方式）

```bash
# 需要记住每个服务的地址和端口
curl http://192.168.31.78:8000/admin/123      # test服务
curl http://192.168.31.78:8001/basic/color/1  # basic服务
```

❌ 问题：
- 客户端需要知道所有服务的地址
- 端口变化需要修改客户端代码
- 难以统一管理认证、日志

### 有网关（推荐）

```bash
# 统一入口，只需要知道网关地址
curl http://localhost:8080/test/admin/123
curl http://localhost:8080/basic/color/1
```

✅ 优势：
- 统一入口（单一域名/IP）
- 服务地址透明（客户端不需要知道）
- 便于添加认证、限流等功能

## 扩展功能

### 1. 添加JWT认证

```go
// 在 proxyHandler 中添加
func (gw *Gateway) authMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        token := c.GetHeader("Authorization")
        if token == "" {
            c.JSON(401, gin.H{"error": "未授权"})
            c.Abort()
            return
        }
        // 验证JWT token
        c.Next()
    }
}

// 使用
r.Use(gateway.authMiddleware())
```

### 2. 添加限流

```go
import "golang.org/x/time/rate"

type Gateway struct {
    limiter *rate.Limiter
    // ...
}

func (gw *Gateway) rateLimitMiddleware() gin.HandlerFunc {
    gw.limiter = rate.NewLimiter(100, 200) // 每秒100个请求，突发200
    
    return func(c *gin.Context) {
        if !gw.limiter.Allow() {
            c.JSON(429, gin.H{"error": "请求过于频繁"})
            c.Abort()
            return
        }
        c.Next()
    }
}
```

### 3. 添加CORS支持

```go
func corsMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
        c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
        
        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }
        
        c.Next()
    }
}

r.Use(corsMiddleware())
```

## 生产部署建议

### 1. 使用环境变量

```go
consulAddr := os.Getenv("CONSUL_ADDR")
if consulAddr == "" {
    consulAddr = "127.0.0.1:8500"
}

port := os.Getenv("GATEWAY_PORT")
if port == "" {
    port = ":8080"
}
```

### 2. 启动多个实例（高可用）

```bash
# 实例1
GATEWAY_PORT=:8080 go run main.go

# 实例2
GATEWAY_PORT=:8081 go run main.go

# 前面再加Nginx负载均衡
```

### 3. 添加监控指标

```go
import "github.com/prometheus/client_golang/prometheus"

var (
    requestCounter = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "gateway_requests_total",
        },
        []string{"service", "method", "status"},
    )
)

// 在代理处理器中记录
requestCounter.WithLabelValues(serviceName, c.Request.Method, fmt.Sprint(c.Writer.Status())).Inc()
```

## 故障排查

### 问题1: 找不到服务

```
错误: 未找到可用的服务实例: testservice
```

**解决**:
1. 检查Consul是否启动：`curl http://127.0.0.1:8500/v1/status/leader`
2. 检查服务是否注册：`curl http://127.0.0.1:8500/v1/catalog/services`
3. 检查后端服务是否运行

### 问题2: 连接Consul失败

```
错误: 连接Consul失败
```

**解决**:
1. 检查Consul地址配置
2. 检查网络连接
3. 检查防火墙规则

### 问题3: 代理转发失败

```
错误: 服务不可用
```

**解决**:
1. 检查后端服务健康状态
2. 检查后端服务端口是否正确
3. 查看网关日志获取详细错误

## 性能优化

1. **连接池**: 使用HTTP连接池减少连接开销
2. **缓存服务地址**: 避免每次请求都查询Consul
3. **负载均衡**: 实现轮询、随机等负载均衡策略
4. **超时控制**: 设置合理的代理超时时间

## 下一步

1. ✅ 基础反向代理（已完成）
2. 🔲 添加JWT认证
3. 🔲 添加限流保护
4. 🔲 添加负载均衡策略
5. 🔲 添加监控指标（Prometheus）
6. 🔲 添加分布式追踪（Jaeger）

## 参考资料

- [完整架构文档](../docs/架构说明.md)
- [API网关指南](../docs/API网关指南.md)
- [Consul集成指南](../docs/Consul集成指南.md)
