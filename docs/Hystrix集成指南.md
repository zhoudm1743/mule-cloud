# Hystrix 熔断器集成指南

## 📖 概述

本项目已集成 **hystrix-go** 熔断器，为微服务提供：
- ✅ **服务降级**: 服务失败时返回默认响应
- ✅ **熔断保护**: 自动隔离故障服务
- ✅ **超时控制**: 防止请求长时间等待
- ✅ **并发限制**: 控制最大并发请求数
- ✅ **实时监控**: 查看熔断器状态和指标

---

## 🏗️ 架构

```
┌─────────────┐
│   客户端     │
└──────┬──────┘
       │ HTTP请求
       ↓
┌─────────────────────┐
│    API网关          │
│  • 限流中间件       │
│  • 认证中间件       │
│  • Hystrix中间件 ⭐ │
└──────┬──────────────┘
       │ (熔断保护)
       ↓
┌─────────────────────┐
│   后端服务          │
│  • testservice      │
│  • basicservice     │
└─────────────────────┘
```

**核心流程**:
1. 请求到达网关 → Hystrix中间件拦截
2. 执行服务调用（受熔断器保护）
3. 如果服务响应正常 → 返回结果
4. 如果服务失败/超时 → 触发熔断，返回降级响应

---

## 📁 项目文件结构

```
mule-cloud/
├── core/
│   └── hystrix/
│       └── hystrix.go          # Hystrix核心封装
├── gateway/
│   └── middleware/
│       └── hystrix.go          # Hystrix中间件
└── docs/
    └── Hystrix集成指南.md      # 本文档
```

---

## ⚙️ 配置说明

### 默认配置

在 `core/hystrix/hystrix.go` 中定义：

```go
DefaultConfig = Config{
    Timeout:                3000,  // 3秒超时
    MaxConcurrentRequests:  100,   // 最多100个并发
    RequestVolumeThreshold: 20,    // 至少20个请求后开始统计
    SleepWindow:            5000,  // 熔断5秒后尝试恢复
    ErrorPercentThreshold:  50,    // 错误率超过50%触发熔断
}
```

### 服务级别配置

```go
ServiceConfigs = map[string]Config{
    "testservice": {
        Timeout:                2000,  // 2秒超时
        MaxConcurrentRequests:  50,    // 最多50个并发
        RequestVolumeThreshold: 10,    // 至少10个请求
        SleepWindow:            3000,  // 熔断3秒后尝试恢复
        ErrorPercentThreshold:  50,    // 错误率阈值50%
    },
    "basicservice": {
        Timeout:                5000,  // 5秒超时
        MaxConcurrentRequests:  100,   // 最多100个并发
        RequestVolumeThreshold: 20,    // 至少20个请求
        SleepWindow:            5000,  // 熔断5秒后尝试恢复
        ErrorPercentThreshold:  60,    // 错误率阈值60%
    },
}
```

---

## 🚀 使用方法

### 1. 网关中使用（已集成）

在 `gateway/main.go` 中已自动集成：

```go
func main() {
    // 初始化Hystrix
    hystrixPkg.Init()

    // 创建网关实例...
    gateway, _ := NewGateway("127.0.0.1:8500")

    // 创建路由
    r := gin.New()

    // 业务接口（自动应用Hystrix）
    api := r.Group("")
    api.Use(middleware.HystrixMiddleware()) // ⭐ Hystrix中间件
    {
        api.Any("/test/*path", gateway.proxyHandler())
        api.Any("/basic/*path", gateway.proxyHandler())
    }

    r.Run(":8080")
}
```

### 2. 在服务层使用（Endpoint层）

在你的服务代码中直接使用：

```go
package endpoint

import (
    hystrixPkg "mule-cloud/core/hystrix"
    "github.com/go-kit/kit/endpoint"
)

// 方法1: 使用Do包装
func GetAdminEndpoint(svc services.IAdminService) endpoint.Endpoint {
    return func(ctx context.Context, request interface{}) (interface{}, error) {
        req := request.(AdminRequest)
        
        var result *services.Admin
        var err error
        
        // 使用Hystrix保护
        hystrixErr := hystrixPkg.Do("get-admin",
            // 正常执行
            func() error {
                result, err = svc.GetAdmin(ctx, req.ID)
                return err
            },
            // 降级处理
            func(hystrixErr error) error {
                // 返回默认值或缓存数据
                result = &services.Admin{
                    ID:   req.ID,
                    Name: "服务暂时不可用",
                }
                return nil
            },
        )
        
        if hystrixErr != nil {
            return nil, hystrixErr
        }
        
        return AdminResponse{Admin: result}, nil
    }
}

// 方法2: 使用DoWithFallbackValue
func GetColorEndpoint(svc services.IColorService) endpoint.Endpoint {
    return func(ctx context.Context, request interface{}) (interface{}, error) {
        req := request.(ColorRequest)
        
        // 带默认值的熔断器
        color, err := hystrixPkg.DoWithFallbackValue(
            "get-color",
            func() (*services.Color, error) {
                return svc.GetColor(ctx, req.ID)
            },
            &services.Color{
                ID:   req.ID,
                Name: "默认颜色",
                Hex:  "#000000",
            },
        )
        
        return ColorResponse{Color: color}, err
    }
}
```

### 3. 自定义配置

```go
// 为特定命令配置Hystrix
hystrixPkg.ConfigureCommand("my-service", hystrixPkg.Config{
    Timeout:                1000,  // 1秒超时
    MaxConcurrentRequests:  10,    // 最多10个并发
    RequestVolumeThreshold: 5,     // 至少5个请求
    SleepWindow:            2000,  // 熔断2秒后尝试恢复
    ErrorPercentThreshold:  30,    // 错误率30%触发熔断
})
```

---

## 🧪 测试熔断器

### 启动服务

```bash
# 启动Consul
consul agent -dev

# 启动所有服务
.\scripts\start_all.bat
```

### 测试场景

#### 1. 正常请求

```bash
# 访问basicservice
curl http://localhost:8080/basic/color/1

# 响应（正常）
{
  "code": 0,
  "msg": "success",
  "data": {
    "id": 1,
    "name": "红色",
    "hex": "#FF0000"
  }
}
```

#### 2. 触发熔断（停止后端服务）

```bash
# 1. 停止basicservice
# 按Ctrl+C停止 basic/cmd/main.go

# 2. 多次请求触发熔断
for i in {1..30}; do
  curl http://localhost:8080/basic/color/1
  sleep 0.1
done

# 响应（熔断降级）
{
  "code": 503,
  "msg": "服务暂时不可用: basicservice",
  "error": "hystrix: circuit open",
  "service": "basicservice",
  "fallback": true
}
```

#### 3. 超时触发熔断

```bash
# 模拟慢服务（在服务代码中添加延迟）
time.Sleep(10 * time.Second) // 超过配置的超时时间

# 响应
{
  "code": 503,
  "msg": "服务暂时不可用: basicservice",
  "error": "hystrix: timeout",
  "service": "basicservice",
  "fallback": true
}
```

---

## 📊 监控与指标

### 查看所有服务熔断器状态

```bash
curl http://localhost:8080/gateway/hystrix/metrics

# 响应
{
  "code": 0,
  "msg": "获取成功",
  "data": {
    "basicservice": {
      "status": "closed",
      "metrics": {
        "total_requests": 150,
        "error_count": 5,
        "error_percentage": 3.33,
        "is_circuit_breaker_open": false
      }
    },
    "testservice": {
      "status": "open",
      "metrics": {
        "total_requests": 50,
        "error_count": 30,
        "error_percentage": 60.0,
        "is_circuit_breaker_open": true
      }
    }
  }
}
```

### 查看指定服务状态

```bash
curl http://localhost:8080/gateway/hystrix/metrics/testservice

# 响应
{
  "code": 0,
  "msg": "获取成功",
  "data": {
    "service": "testservice",
    "status": "closed",
    "metrics": {
      "total_requests": 100,
      "error_count": 5,
      "error_percentage": 5.0,
      "is_circuit_breaker_open": false
    }
  }
}
```

---

## 🔥 熔断器状态说明

### Closed（关闭状态）

- **含义**: 正常工作，请求正常通过
- **条件**: 错误率低于阈值
- **行为**: 所有请求正常执行

### Open（打开状态）

- **含义**: 熔断器打开，直接返回降级响应
- **条件**: 错误率超过阈值
- **行为**: 
  - 所有请求立即失败（快速失败）
  - 执行fallback降级函数
  - 等待SleepWindow时间后进入Half-Open

### Half-Open（半开状态）

- **含义**: 尝试恢复中
- **条件**: SleepWindow时间后
- **行为**: 
  - 允许少量请求通过
  - 如果成功 → 进入Closed
  - 如果失败 → 重新进入Open

---

## ⚙️ 高级用法

### 1. 异步执行

```go
errChan := hystrixPkg.Go("my-command",
    func() error {
        // 执行异步操作
        return doSomething()
    },
    func(err error) error {
        // 降级处理
        return handleFallback(err)
    },
)

// 等待结果
if err := <-errChan; err != nil {
    log.Printf("异步操作失败: %v", err)
}
```

### 2. 等待熔断器恢复

```go
err := hystrixPkg.WaitForHealthyCircuit("my-service", 30*time.Second)
if err != nil {
    log.Printf("服务恢复超时: %v", err)
}
```

### 3. 手动控制熔断器

```go
// 检查状态
status := hystrixPkg.CircuitBreakerStatus("my-service")
if status == "open" {
    log.Println("熔断器已打开")
}

// 获取指标
metrics, err := hystrixPkg.GetMetrics("my-service")
if err == nil {
    fmt.Printf("总请求: %d, 错误: %d, 错误率: %.2f%%\n",
        metrics.TotalRequests,
        metrics.ErrorCount,
        metrics.ErrorPercentage)
}

// 刷新指标
hystrixPkg.FlushMetrics()
```

---

## 🎯 最佳实践

### 1. 合理设置超时时间

```go
// ❌ 错误：超时设置过短
Config{Timeout: 100} // 100ms可能不够

// ✅ 正确：根据实际响应时间设置
Config{Timeout: 3000} // 3秒合理
```

### 2. 降级策略

```go
// ❌ 错误：降级函数抛出错误
func(err error) error {
    return fmt.Errorf("服务不可用: %v", err)
}

// ✅ 正确：返回默认值
func(err error) error {
    result = getFromCache() // 从缓存获取
    return nil
}
```

### 3. 合理设置错误率阈值

```go
// 根据服务重要性调整
ServiceConfigs = map[string]Config{
    "critical-service": {
        ErrorPercentThreshold: 30, // 核心服务：容忍度低
    },
    "optional-service": {
        ErrorPercentThreshold: 70, // 可选服务：容忍度高
    },
}
```

### 4. 监控告警

```go
// 定期检查熔断器状态
go func() {
    ticker := time.NewTicker(30 * time.Second)
    for range ticker.C {
        allStatus := hystrixPkg.GetAllCircuitStatus()
        for svc, status := range allStatus {
            if status.(map[string]interface{})["status"] == "open" {
                log.Printf("⚠️  警告：%s 熔断器已打开", svc)
                // 发送告警通知
            }
        }
    }
}()
```

---

## 🐛 故障排查

### 熔断器一直打开

**原因**:
- 后端服务真的故障
- 超时时间设置过短
- 错误率阈值设置过低

**解决**:
```bash
# 1. 检查后端服务是否正常
curl http://localhost:8000/health

# 2. 查看熔断器指标
curl http://localhost:8080/gateway/hystrix/metrics/testservice

# 3. 调整配置
# 修改 core/hystrix/hystrix.go 中的 ServiceConfigs
```

### 降级函数未执行

**原因**: 降级函数返回错误

**解决**:
```go
// ✅ 确保降级函数返回nil
func(err error) error {
    // 处理降级逻辑
    return nil // 返回nil表示降级成功
}
```

### 请求被拒绝（并发限制）

**原因**: `MaxConcurrentRequests` 设置过小

**解决**:
```go
Config{
    MaxConcurrentRequests: 200, // 增加并发限制
}
```

---

## 📋 API接口

| 接口 | 方法 | 说明 |
|------|------|------|
| `/gateway/hystrix/metrics` | GET | 获取所有服务熔断器状态 |
| `/gateway/hystrix/metrics/:service` | GET | 获取指定服务熔断器状态 |

---

## 📚 参考资料

- [hystrix-go GitHub](https://github.com/afex/hystrix-go)
- [Netflix Hystrix Wiki](https://github.com/Netflix/Hystrix/wiki)
- [微服务熔断器模式](https://martinfowler.com/bliki/CircuitBreaker.html)

---

## 🔄 更新日志

### v1.0.0 (2025-09-30)
- ✅ 集成 hystrix-go
- ✅ 网关Hystrix中间件
- ✅ 服务级别配置
- ✅ 监控指标接口
- ✅ 完整文档

---

**享受熔断保护！🔥**
