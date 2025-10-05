# 重要 - 租户 Code 使用说明

## 🎯 核心改动

**数据库命名从使用 ID 改为使用 Code**

### 修改前后对比

| 项目 | 修改前 | 修改后 |
|------|-------|-------|
| 数据库名 | `mule_68e27febab849776302...` | `mule_default` |
| Context 传递 | `tenantID` (实际是ID) | `tenantCode` (实际是Code) |
| 可读性 | ❌ 不易识别 | ✅ 语义明确 |
| 管理性 | ❌ 难以管理 | ✅ 易于管理 |

---

## 📋 使用指南

### 1. 创建租户时

**必须提供唯一的租户代码**：

```javascript
POST /system/tenants
{
  "code": "default",      // ✅ 必填，唯一
  "name": "默认租户",
  // ...
}
```

**数据库名称**：`mule_default`

### 2. Context 传递

**推荐使用新的函数**：

```go
import tenantCtx "mule-cloud/core/context"

// 设置租户上下文
ctx = tenantCtx.WithTenantCode(ctx, "default")

// 获取租户上下文
tenantCode := tenantCtx.GetTenantCode(ctx)
```

**向后兼容**（仍然有效，但不推荐）：

```go
// 旧代码仍然有效
ctx = tenantCtx.WithTenantID(ctx, "default")  // 虽然叫ID，但传入code
tenantCode := tenantCtx.GetTenantID(ctx)       // 虽然叫ID，但返回code
```

### 3. 数据库操作

**Repository 自动处理**：

```go
// Repository 内部会自动使用 GetDatabase(tenantCode)
// 开发者无需关心数据库切换细节
func (r *BasicRepository) GetByID(ctx context.Context, id string) (*models.Basic, error) {
    collection := r.getCollection(ctx)  // 自动切换到 mule_<code>
    // ...
}
```

---

## ⚠️ 重要注意事项

### 1. Code 必须唯一

```go
// ✅ 系统会检查 code 是否重复
existing, _ := s.repo.GetByCode(ctx, req.Code)
if existing != nil {
    return nil, fmt.Errorf("租户代码已存在")
}
```

**数据库索引**：
```javascript
db.tenant.createIndex({ "code": 1 }, { unique: true })
```

### 2. Code 命名规范

**推荐格式**：
- ✅ `default` - 默认租户
- ✅ `company_a` - 公司A
- ✅ `test01` - 测试租户01
- ✅ `ace` - 简短代码

**不推荐**：
- ❌ `公司A` - 避免中文
- ❌ `company-a` - 避免短横线（使用下划线）
- ❌ `Company_A` - 避免大小写混合（MongoDB数据库名区分大小写）

### 3. 中间件和鉴权

**JWT 中存储的是什么？**

看登录逻辑：

```go
// app/auth/services/auth.go
if req.TenantCode != "" {
    tenant, err := s.tenantRepo.GetByCode(ctx, req.TenantCode)
    tenantID = tenant.ID  // JWT 中存储的是 ID！
}
```

**重要**：JWT Token 中存储的仍然是 `tenant_id`（实际的MongoDB ID），不是 code！

**为什么？**
1. ✅ ID 永不改变，code 可能被修改
2. ✅ ID 是主键，查询效率更高
3. ✅ 向后兼容，不影响现有 JWT

**那么 code 在哪里使用？**

在创建数据库连接时，需要通过 ID 查询 tenant 对象获取 code：

```go
// repository 的 getCollection 方法需要修改
func (r *BasicRepository) getCollection(ctx context.Context) *mongo.Collection {
    tenantID := tenantCtx.GetTenantID(ctx)  // 从 JWT 获取 ID
    
    if tenantID == "" {
        db := r.dbManager.GetDatabase("")  // 系统库
        return db.Collection("basic")
    }
    
    // ✅ 需要通过 ID 查询 tenant 获取 code
    tenant, _ := tenantRepo.Get(context.Background(), tenantID)
    if tenant == nil {
        // 降级：仍然使用 ID（兼容旧数据）
        db := r.dbManager.GetDatabase(tenantID)
    } else {
        // 使用 code
        db := r.dbManager.GetDatabase(tenant.Code)
    }
    return db.Collection("basic")
}
```

**等等，这样太复杂了！**

---

## 🔄 更好的方案

### 方案A：JWT 中同时存储 ID 和 Code

**修改 JWT Claims**：

```go
type Claims struct {
    UserID    string   `json:"user_id"`
    Username  string   `json:"username"`
    TenantID  string   `json:"tenant_id"`   // MongoDB ID
    TenantCode string  `json:"tenant_code"` // ✅ 新增
    Roles     []string `json:"roles"`
    jwt.RegisteredClaims
}
```

**Context 传递**：

```go
// 中间件解析 JWT 后
ctx = tenantCtx.WithTenantCode(ctx, claims.TenantCode)
```

**Repository 使用**：

```go
func (r *BasicRepository) getCollection(ctx context.Context) *mongo.Collection {
    tenantCode := tenantCtx.GetTenantCode(ctx)
    db := r.dbManager.GetDatabase(tenantCode)
    return db.Collection("basic")
}
```

### 方案B：Context 中传递 Code，Repository 查询时转换

保持 JWT 不变，在中间件层转换：

```go
// 中间件：从 JWT 的 tenant_id 获取 code
tenantID := claims.TenantID
if tenantID != "" {
    tenant, _ := tenantRepo.Get(ctx, tenantID)
    if tenant != nil {
        ctx = tenantCtx.WithTenantCode(ctx, tenant.Code)
    }
}
```

---

## 🎯 推荐方案

**方案A - 在JWT中同时存储ID和Code**

### 优点
1. ✅ 性能好：不需要额外查询
2. ✅ 简单：Repository 直接使用 code
3. ✅ 清晰：语义明确

### 缺点
1. ⚠️ JWT 略大：多了一个字段
2. ⚠️ Code 改变时需要重新登录

### 实施步骤

1. **修改 JWT Claims** (core/jwt/jwt.go)
2. **修改登录逻辑** (app/auth/services/auth.go)
3. **修改中间件** (core/middleware/jwt.go)
4. **Repository 保持不变**（已经使用 GetDatabase）

---

## 📝 当前状态

### 已完成 ✅
- ✅ `core/database/manager.go` - GetDatabase 使用 code
- ✅ `core/context/tenant.go` - 添加 WithTenantCode/GetTenantCode
- ✅ `app/system/services/tenant.go` - CreateTenantDatabase 使用 code

### 待完成 ⏳
- ⏳ JWT Claims 添加 tenant_code 字段
- ⏳ 登录时查询 tenant.Code 并存入 JWT
- ⏳ 中间件解析 JWT 后设置 tenantCode 到 Context
- ⏳ 测试验证

---

## 🚀 下一步

请告诉我是否需要：
1. **继续完成 JWT 改动**（推荐）
2. **保持现状**（使用 ID 作为数据库名，不做进一步修改）
3. **实施方案B**（中间件层查询转换）

**您的选择是？** 🤔
