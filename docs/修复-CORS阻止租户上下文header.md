# 修复 - CORS 阻止租户上下文 header ✅

## 🐛 问题

切换租户后出现 CORS 错误：

```
Access to fetch at 'http://localhost:8080/admin/auth/getUserRoutes' from origin 
'http://localhost:9980' has been blocked by CORS policy: 
Request header field x-tenant-context is not allowed by 
Access-Control-Allow-Headers in preflight response.
```

## 🔍 原因

**CORS 配置没有允许 `X-Tenant-Context` header**

当系统管理员切换租户时：
1. 前端添加 `X-Tenant-Context` header 到请求中
2. 浏览器发送 OPTIONS 预检请求
3. **网关 CORS 中间件拒绝了该 header** ❌
4. 浏览器阻止实际请求
5. 前端报错：Failed to fetch

---

## ✅ 解决方案

**在 CORS 配置中添加 `X-Tenant-Context` 到允许的 headers**

### 修改文件

**`app/gateway/middleware/cors.go`**

#### 之前（❌ 缺少 X-Tenant-Context）

```go
c.Writer.Header().Set("Access-Control-Allow-Headers", 
    "Content-Type, Content-Length, Authorization, Accept, X-Requested-With")
```

#### 现在（✅ 添加 X-Tenant-Context）

```go
c.Writer.Header().Set("Access-Control-Allow-Headers", 
    "Content-Type, Content-Length, Authorization, Accept, X-Requested-With, X-Tenant-Context")
```

---

## 🚀 重启网关

```powershell
# 停止网关（Ctrl+C）

# 重启
go run cmd/gateway/main.go
```

---

## 🎯 验证

### 1. 检查浏览器 Network

**OPTIONS 预检请求** 应该返回：
```
Access-Control-Allow-Headers: Content-Type, Content-Length, Authorization, 
                              Accept, X-Requested-With, X-Tenant-Context
```

### 2. 测试流程

```
1. 系统管理员登录 ✅

2. 选择一个租户 ✅

3. 页面刷新 ✅

4. 查看 Network:
   - OPTIONS 请求成功 ✅
   - GET getUserRoutes 请求成功 ✅
   - 没有 CORS 错误 ✅

5. 能正常看到该租户的数据 ✅
```

---

## 📝 CORS 完整配置

```go
// app/gateway/middleware/cors.go
func CORS() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
        c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
        c.Writer.Header().Set("Access-Control-Allow-Headers", 
            "Content-Type, Content-Length, Authorization, Accept, X-Requested-With, X-Tenant-Context")
        c.Writer.Header().Set("Access-Control-Allow-Methods", 
            "GET, POST, PUT, DELETE, OPTIONS, PATCH")
        c.Writer.Header().Set("Access-Control-Max-Age", "86400")

        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }

        c.Next()
    }
}
```

---

## 🎉 完成

修复后，系统管理员应该能够：
- ✅ 切换租户（不再 CORS 错误）
- ✅ 获取菜单
- ✅ 查看租户数据
- ✅ 所有功能正常

**重启网关后测试！** 🚀
