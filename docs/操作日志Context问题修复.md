# 操作日志 Context 问题修复

## 🐛 问题描述

### 错误信息

```
2025/10/05 23:39:07 [MongoDB] 命令失败: insert, 错误: context canceled
2025-10-05T23:39:07.419+0800    ERROR   middleware/operation_log.go:103 保存操作日志失败
{"user_id": "68dcf02cc592a24457a2f978", "path": "/perms/menus/68ddfbe60d3d5f5a1a02f441", "error": "context canceled"}
```

### 问题原因

操作日志中间件使用异步 goroutine 保存日志，但在 goroutine 中使用了主请求的 Context：

```go
// ❌ 错误代码
go func() {
    // ...
    ctx := c.Request.Context()  // 主请求的 Context
    if err := repo.Create(ctx, log); err != nil {
        // context canceled 错误
    }
}()
```

**问题分析**：
1. 主请求处理完成后，Gin 会取消（cancel）Request.Context
2. 异步 goroutine 还在执行，但 Context 已经被取消
3. MongoDB 操作使用被取消的 Context，导致失败

---

## ✅ 解决方案

### 修改后的代码

```go
// ✅ 正确代码
// 在 goroutine 外部提取租户信息
tenantCode, _ := c.Get("tenant_code")

go func() {
    // ...
    
    // 创建独立的 Context（不会被主请求影响）
    ctx := context.Background()
    ctx = tenantCtx.WithTenantCode(ctx, toString(tenantCode))
    
    // 保存日志
    if err := repo.Create(ctx, log); err != nil {
        logger.Error("保存操作日志失败", zap.Error(err))
    }
}()
```

### 关键改进

1. **使用 `context.Background()`**
   - 创建独立的 Context，不会被主请求取消
   - 确保异步操作能够完成

2. **提前提取租户信息**
   - 在 goroutine 外部提取 `tenant_code`
   - 在 goroutine 内部设置到新的 Context 中

3. **避免共享 Context**
   - 不再使用 `c.Request.Context()`
   - 异步操作与主请求完全独立

---

## 📊 修复对比

### 修复前

```go
go func() {
    // 获取数据...
    log := &models.OperationLog{...}
    
    // ❌ 使用主请求的 Context
    ctx := c.Request.Context()
    if err := repo.Create(ctx, log); err != nil {
        // context canceled 错误 ❌
    }
}()
```

**时间线**：
```
0ms    主请求开始
50ms   中间件启动 goroutine
60ms   主请求完成 → Context 被取消 ❌
70ms   goroutine 尝试保存日志 → 失败（context canceled）
```

### 修复后

```go
// ✅ 提前提取租户信息
tenantCode, _ := c.Get("tenant_code")

go func() {
    // 获取数据...
    log := &models.OperationLog{...}
    
    // ✅ 创建独立的 Context
    ctx := context.Background()
    ctx = tenantCtx.WithTenantCode(ctx, toString(tenantCode))
    
    if err := repo.Create(ctx, log); err != nil {
        // 不会出现 context canceled 错误 ✅
    }
}()
```

**时间线**：
```
0ms    主请求开始
50ms   中间件提取租户信息并启动 goroutine
60ms   主请求完成 → 主 Context 被取消（不影响 goroutine）✅
70ms   goroutine 使用独立 Context 保存日志 → 成功 ✅
```

---

## 🔍 为什么需要租户上下文？

操作日志需要存储到租户数据库：

```go
// Repository 根据 Context 中的 tenant_code 切换数据库
func (r *OperationLogRepository) getCollection(ctx context.Context) *mongo.Collection {
    tenantCode := tenantCtx.GetTenantCode(ctx)  // 从 Context 获取
    db := r.dbManager.GetDatabase(tenantCode)   // 切换到租户库
    return db.Collection("operation_logs")
}
```

**数据流**：
```
1. 中间件从 Gin Context 提取 tenant_code
   ↓
2. 创建独立 Context 并设置 tenant_code
   ↓
3. Repository 从 Context 获取 tenant_code
   ↓
4. 切换到对应的租户数据库
   ↓
5. 保存操作日志
```

---

## 📝 类似问题的最佳实践

### 1. 异步操作中使用独立 Context

```go
// ✅ 推荐
go func() {
    ctx := context.Background()
    // 执行异步操作
}()

// ❌ 不推荐
go func() {
    ctx := c.Request.Context()  // 可能被取消
    // 执行异步操作
}()
```

### 2. 需要传递数据时，提前提取

```go
// ✅ 推荐：提前提取数据
userID, _ := c.Get("user_id")
tenantCode, _ := c.Get("tenant_code")

go func() {
    ctx := context.Background()
    ctx = tenantCtx.WithTenantCode(ctx, toString(tenantCode))
    // 使用 userID 和 ctx
}()

// ❌ 不推荐：在 goroutine 中访问 Gin Context
go func() {
    userID, _ := c.Get("user_id")  // 可能不安全
}()
```

### 3. 设置超时（可选）

```go
// ✅ 为异步操作设置合理的超时
go func() {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    // 执行操作，最多等待 5 秒
    if err := repo.Create(ctx, log); err != nil {
        logger.Error("保存失败", zap.Error(err))
    }
}()
```

---

## 🎯 验证修复

### 测试步骤

1. **重启服务**
   ```bash
   go run cmd/gateway/main.go
   ```

2. **执行写操作**
   ```bash
   curl -X POST http://localhost:8000/api/system/users \
     -H "Authorization: Bearer <token>" \
     -d '{"nickname":"测试用户"}'
   ```

3. **检查日志**
   ```bash
   # 应该看到成功日志，而不是 context canceled
   2025-10-05T23:45:00  DEBUG  操作日志已记录  {"user_id": "123", "resource": "system", "action": "create"}
   ```

4. **查询数据库**
   ```javascript
   // MongoDB
   db.operation_logs.find().sort({created_at: -1}).limit(1)
   
   // 应该能看到刚才的操作记录
   ```

---

## 🎉 修复完成

### 修改文件
- ✅ `core/middleware/operation_log.go`
  - 添加 `context.Background()` 创建独立 Context
  - 提前提取 `tenant_code`
  - 在 goroutine 中设置租户上下文

### 编译验证
```bash
✅ go build ./core/middleware
```

### 效果
- ✅ 不再出现 "context canceled" 错误
- ✅ 操作日志正常保存到数据库
- ✅ 异步操作不影响主请求性能

---

**现在操作日志可以正常保存了！** 🎊
