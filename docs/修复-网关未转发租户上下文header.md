# 修复：网关未转发租户上下文 Header

## 问题描述

前端发送了 `X-Tenant-Context: ace` header（浏览器开发者工具可见），但后端服务收到的是空字符串，导致超管无法切换租户。

## 问题根源

**网关转发请求时，只转发了从 JWT 解析出来的用户信息 header，没有转发前端直接发送的 `X-Tenant-Context` header。**

### 数据流分析

1. **前端发送请求：**
   ```
   GET /admin/perms/admins
   Headers:
     Authorization: Bearer xxx...
     X-Tenant-Context: ace  ← 前端发送
   ```

2. **网关接收请求：**
   - CORS 中间件：允许 `X-Tenant-Context` 通过 ✅
   - JWT 中间件：解析 token，设置 `c.Set("tenant_code", "system")`
   - 代理中间件：转发请求到后端

3. **网关转发（修复前）：**
   ```go
   // ❌ 只转发从 JWT 解析的信息，没有转发前端的 header
   if tenantCode, exists := c.Get("tenant_code"); exists {
       c.Request.Header.Set("X-Tenant-Code", tenantCode.(string))  // = "system"
   }
   // ❌ X-Tenant-Context 没有被转发！
   ```

4. **后端接收（修复前）：**
   ```
   Headers:
     X-Tenant-Code: system  ← 网关从 JWT 解析的
     X-Tenant-Context: (空)  ← 前端发送的但网关没转发
   ```

## 修复方案

在网关的代理转发逻辑中，添加对 `X-Tenant-Context` header 的直接转发：

```go
// cmd/gateway/main.go

// ✅ 重要：转发前端发送的 X-Tenant-Context header（用于超管切换租户）
// 这个 header 是前端直接发送的，不在 JWT token 中，需要单独转发
if contextTenant := c.GetHeader("X-Tenant-Context"); contextTenant != "" {
    c.Request.Header.Set("X-Tenant-Context", contextTenant)
    log.Printf("[网关转发] 转发租户上下文: %s", contextTenant)
}
```

### 为什么要单独转发？

- `X-Tenant-Code`：从 JWT token 中解析，表示用户所属的租户（存储在 token 中）
- `X-Tenant-Context`：前端临时发送，表示超管**想要切换到的租户**（不在 token 中）

这两个是不同的概念：
- 系统管理员的 `X-Tenant-Code` = `"system"`（来自 JWT）
- 系统管理员想切换到的租户 `X-Tenant-Context` = `"ace"`（来自前端）

后端的 `TenantContextMiddleware` 会检查：
1. 用户的 `tenant_code` 是否为 `"system"`（有权限切换）
2. 用户是否有 `super` 角色（是超管）
3. `X-Tenant-Context` header 的值（想切换到哪个租户）

## 修改的文件

1. `cmd/gateway/main.go` - 添加 `X-Tenant-Context` header 转发逻辑

## 完整的数据流（修复后）

### 1. 前端发送
```http
GET /admin/perms/admins
Headers:
  Authorization: Bearer xxx...
  X-Tenant-Context: ace
```

### 2. 网关处理
```
[CORS 中间件] ✅ 允许 X-Tenant-Context
[JWT 中间件] ✅ 解析 token → c.Set("tenant_code", "system")
[代理中间件] ✅ 转发 headers:
  - X-Tenant-Code: system  (从 JWT 解析)
  - X-Tenant-Context: ace  (从前端直接转发) ← 新增
```

### 3. 后端接收
```
[GatewayOrJWTAuth] 从 header 读取:
  - X-Tenant-Code: system
  - ctx.WithTenantCode(ctx, "system")

[TenantContextMiddleware] 检查切换:
  - currentTenantCode = "system" ✅
  - X-Tenant-Context header = "ace" ✅
  - isSuperAdmin = true ✅
  - → 切换到租户: ace
  - ctx.WithTenantCode(ctx, "ace")
```

### 4. Repository 查询
```go
tenantCode := GetTenantCode(ctx)  // = "ace"
db := GetDatabase("ace")          // = mule_ace
```

### 5. MongoDB 查询
```
[MongoDB] 执行命令: ... "$db": "mule_ace"}  ✅
```

## 测试验证

### 1. 重启网关和 perms 服务

```bash
# 重启网关
cd cmd/gateway
./gateway.exe

# 重启 perms 服务
cd cmd/perms
./perms.exe
```

### 2. 前端测试

1. 以系统管理员身份登录
2. 选择租户 "ace"
3. 查看浏览器控制台 → Network → 请求 headers：
   ```
   X-Tenant-Context: ace  ✅
   ```

### 3. 查看网关日志

```
[网关转发] 转发租户上下文: ace
[网关转发] GET /admin/perms/admins → http://localhost:8085/perms/admins (服务: perms-service, 用户: admin)
```

### 4. 查看 perms 服务日志

```
[租户上下文] 当前 tenantCode: 'system'
[租户上下文] X-Tenant-Context header: 'ace'  ✅ 不再是空
[租户上下文] isSuperAdmin: true
[租户上下文切换] ✅ 系统管理员切换到租户: ace
[MongoDB] 执行命令: ... "$db": "mule_ace"}  ✅
```

## 为什么之前 CORS 配置没问题？

你可能会问：CORS 中间件已经配置了允许 `X-Tenant-Context`，为什么还是收不到？

```go
// app/gateway/middleware/cors.go
c.Writer.Header.Set("Access-Control-Allow-Headers", "..., X-Tenant-Context, ...")
```

**答案：**
- CORS 配置只是告诉浏览器"允许前端发送这个 header"
- 但网关在转发请求给后端时，需要**主动复制这个 header**
- 网关默认只转发自己设置的 headers（从 JWT 解析的），不会自动转发所有前端 headers

这是两个不同的层面：
1. **浏览器 → 网关**：CORS 控制（已配置 ✅）
2. **网关 → 后端**：代理转发逻辑控制（需要手动添加 ✅）

## 安全性说明

有人可能担心：如果网关直接转发 `X-Tenant-Context`，是否存在安全风险？

**答案：安全，因为后端有严格验证！**

后端的 `TenantContextMiddleware` 会检查：
```go
// 只有系统管理员可以切换
if currentTenantCode == "" || currentTenantCode == "system" {
    // 只有 super 角色可以切换
    if isSuperAdmin {
        // ✅ 允许切换
    } else {
        // ❌ 拒绝
    }
}
```

所以即使普通租户用户伪造 `X-Tenant-Context` header：
1. 他的 `X-Tenant-Code` = `"ace"`（从 JWT，无法伪造）
2. 后端检查 `currentTenantCode != "system"` → 不允许切换
3. 请求被拒绝 ❌

## 经验教训

### 1. 网关是无状态的转发器

网关不知道哪些 headers 是重要的，只会转发明确配置的 headers。

**错误假设：**
> "CORS 配置了，网关应该会自动转发所有 headers"

**正确理解：**
> "CORS 只控制浏览器能发送什么，网关转发需要单独配置"

### 2. 调试时要跟踪完整链路

问题排查应该检查：
1. 前端是否发送 ✅
2. 网关是否接收 ✅
3. 网关是否转发 ❌ ← 问题点
4. 后端是否接收
5. 后端是否处理

不要只检查一端，要跟踪完整的数据流。

### 3. 日志的重要性

添加详细的日志帮助快速定位问题：
```go
log.Printf("[网关转发] 转发租户上下文: %s", contextTenant)
```

在中间件中打印 header 也很有帮助：
```go
log.Printf("[租户上下文调试] 所有请求头: %v", c.Request.Header)
```

## 相关问题修复

这个问题修复后，之前的两个问题都解决了：

1. ✅ **数据库未正常获取**（后端问题）
   - 修复：Service 层传递 context
   - 文件：`app/perms/services/admin.go` 等

2. ✅ **页面跳转登录页**（前端问题）  
   - 修复：简化 HTTP 拦截器逻辑
   - 文件：`frontend/src/service/http/alova.ts`

3. ✅ **网关未转发 header**（网关问题） ← 本次修复
   - 修复：添加 `X-Tenant-Context` 转发
   - 文件：`cmd/gateway/main.go`

现在整个租户切换功能应该完全正常了！

