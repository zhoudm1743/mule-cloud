# 快速修复 - getUserRoutes CORS 错误

## 问题

切换租户后，前端疯狂调用 `/auth/getUserRoutes` 接口，但全部失败：
```
[GET] - [/auth/getUserRoutes] - Failed to fetch
CORS error
```

## 原因

**auth 和 system 服务缺少租户上下文中间件**，导致无法处理带有 `X-Tenant-Context` header 的请求。

## 已修复 ✅

### 修改的文件

1. **`cmd/auth/main.go`** ✅
   - 添加了 `coreMdw.TenantContextMiddleware()`
   
2. **`cmd/system/main.go`** ✅
   - 添加了 `coreMdw.TenantContextMiddleware()`

3. **`cmd/basic/main.go`** ✅
   - 之前已添加

## 解决步骤

### 1. 重启所有服务

```powershell
# 终止所有正在运行的服务（Ctrl+C）

# 重启 auth 服务
go run cmd/auth/main.go

# 新开一个终端，重启 system 服务
go run cmd/system/main.go

# 新开一个终端，重启 basic 服务
go run cmd/basic/main.go

# 如果使用网关，也重启网关
go run cmd/gateway/main.go
```

### 2. 刷新前端页面

```
1. 清除浏览器缓存（Ctrl+Shift+R）
2. 或者完全刷新页面
```

### 3. 测试流程

```
1. 系统管理员登录
   账号: 17858361617
   密码: 123456
   租户代码: （留空）

2. 登录成功后，顶部应该看到租户选择器

3. 选择一个租户（例如"测试租户 default"）

4. 页面刷新，应该：
   ✅ 没有 CORS 错误
   ✅ 没有 getUserRoutes 失败
   ✅ 正常显示该租户的菜单和数据

5. 查看浏览器 Network 面板：
   - getUserRoutes 请求应该成功（200）
   - 请求 headers 应该包含：
     X-Tenant-Context: {租户ID}

6. 后端日志应该显示：
   [租户上下文切换] 系统管理员切换到租户: {租户ID}
```

## 验证

### 检查 HTTP Headers

打开浏览器 F12 → Network → 找到 `getUserRoutes` 请求：

**Request Headers** 应该包含：
```
Authorization: Bearer {token}
X-Tenant-Context: 68dda6cd04ba0d6c8dda4b7a  ← 租户ID
```

**Response** 应该：
```
Status: 200 OK
```

### 检查后端日志

应该看到类似：
```
[租户上下文切换] 系统管理员切换到租户: 68dda6cd04ba0d6c8dda4b7a
🔗 创建租户数据库连接: mule_68dda6cd04ba0d6c8dda4b7a
```

## 如果还有问题

### 1. 检查中间件顺序

确保中间件顺序正确：
```go
protected.Use(middleware.JWTAuth(jwtManager))      // 1. JWT 认证
protected.Use(coreMdw.TenantContextMiddleware())   // 2. 租户上下文切换
```

### 2. 检查 CORS 配置

如果使用网关，检查网关的 CORS 配置是否允许 `X-Tenant-Context` header。

### 3. 清除 localStorage

如果租户选择器有问题，清除 localStorage：
```javascript
// 在浏览器 Console 执行
localStorage.removeItem('selected_tenant_id')
```

### 4. 查看完整日志

```powershell
# 增加日志级别
go run cmd/auth/main.go -config config/auth.yaml
```

## 成功标志

当你看到：
- ✅ 顶部有租户选择器
- ✅ 选择租户后页面正常刷新
- ✅ 没有 CORS 错误
- ✅ 没有 getUserRoutes 失败
- ✅ 能看到该租户的数据

说明问题已解决！🎉
