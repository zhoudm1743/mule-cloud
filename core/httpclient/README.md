# HTTP Client 服务间调用客户端

## 功能特性

- ✅ **服务发现**：自动从 Consul 获取服务地址
- ✅ **负载均衡**：支持多实例服务的负载均衡
- ✅ **自动鉴权**：支持 JWT token 自动传递
- ✅ **上下文传递**：支持用户ID、租户ID等信息传递
- ✅ **超时控制**：支持请求超时和上下文取消
- ✅ **错误处理**：统一的错误处理机制

## 使用方法

### 1. 创建客户端

```go
import "mule-cloud/core/httpclient"

// 创建服务客户端
client, err := httpclient.NewServiceClient("localhost:8500")
if err != nil {
    log.Fatal(err)
}

// （可选）设置默认的服务间调用 token
client.SetDefaultToken("your-service-token")
```

### 2. 基本调用

```go
ctx := context.Background()

// GET 请求
data, err := client.Get(ctx, "system-service", "/system/menus/all", nil)

// POST 请求
reqBody := map[string]interface{}{"name": "test"}
data, err := client.Post(ctx, "system-service", "/system/menus", reqBody, nil)
```

### 3. 自动鉴权调用

#### 方式1：通过 Context 传递 Token

```go
import "mule-cloud/core/httpclient"

// 从请求中获取 token
token := c.GetHeader("Authorization") // Bearer xxx

// 将 token 添加到 context
ctx := httpclient.WithToken(context.Background(), token)

// 调用服务时自动带上 token
client.Get(ctx, "system-service", "/system/menus/all", nil)
```

#### 方式2：手动传递 Headers

```go
headers := map[string]string{
    "Authorization": "Bearer " + token,
}
client.Get(ctx, "system-service", "/system/menus/all", headers)
```

### 4. 传递用户上下文

```go
import "mule-cloud/core/httpclient"

// 构建包含用户信息的 context
ctx := context.Background()
ctx = httpclient.WithToken(ctx, token)
ctx = httpclient.WithUserID(ctx, userID)
ctx = httpclient.WithTenantID(ctx, tenantID)

// 服务调用时会自动添加以下 headers：
// - Authorization: Bearer <token>
// - X-User-ID: <userID>
// - X-Tenant-ID: <tenantID>
client.Get(ctx, "system-service", "/system/roles", nil)
```

### 5. 完整示例：在 Gateway 中转发请求

```go
func ProxyHandler(client *httpclient.ServiceClient) gin.HandlerFunc {
    return func(c *gin.Context) {
        // 1. 从请求中提取认证信息
        token := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")
        userID := c.GetString("user_id")
        tenantID := c.GetString("tenant_id")
        
        // 2. 构建带认证信息的 context
        ctx := context.Background()
        ctx = httpclient.WithToken(ctx, token)
        ctx = httpclient.WithUserID(ctx, userID)
        ctx = httpclient.WithTenantID(ctx, tenantID)
        
        // 3. 调用后端服务
        var result interface{}
        err := client.CallService(
            ctx,
            c.Request.Method,
            "system-service",
            c.Request.URL.Path,
            c.Request.Body,
            &result,
            nil,
        )
        
        if err != nil {
            c.JSON(500, gin.H{"error": err.Error()})
            return
        }
        
        c.JSON(200, result)
    }
}
```

### 6. 在 Auth Service 中使用

```go
func (s *AuthService) GetUserRoutes(ctx context.Context, userID string) ([]Route, error) {
    // 构建包含用户信息的 context
    ctx = httpclient.WithUserID(ctx, userID)
    
    // 调用 system 服务获取路由
    var result struct {
        Code int     `json:"code"`
        Data []Route `json:"data"`
    }
    
    err := s.httpClient.CallService(
        ctx,
        "GET",
        "system-service",
        fmt.Sprintf("/system/users/%s/routes", userID),
        nil,
        &result,
        nil,
    )
    
    if err != nil {
        return nil, err
    }
    
    return result.Data, nil
}
```

## Token 优先级

当有多种方式提供 token 时，优先级如下：

1. **Context 中的 token** - 最高优先级
2. **手动传递的 headers** - 中等优先级
3. **默认 token (defaultToken)** - 最低优先级

## 最佳实践

### 1. Gateway 中间件

在 Gateway 中创建中间件，自动将请求的认证信息注入到 context：

```go
func AuthContextMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        token := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")
        if token != "" {
            ctx := httpclient.WithToken(c.Request.Context(), token)
            c.Request = c.Request.WithContext(ctx)
        }
        c.Next()
    }
}
```

### 2. 服务间调用使用 Service Token

对于不需要用户上下文的服务间调用（如定时任务），使用专用的 service token：

```go
client.SetDefaultToken(os.Getenv("SERVICE_TOKEN"))
```

### 3. 超时控制

始终使用带超时的 context：

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

client.Get(ctx, "system-service", "/api/data", nil)
```

## 注意事项

1. **Token 安全**：不要在日志中打印 token
2. **Consul 连接**：确保 Consul 服务可用，否则服务发现会失败
3. **超时设置**：根据业务需求调整客户端超时时间
4. **错误处理**：妥善处理网络错误和服务不可用的情况

