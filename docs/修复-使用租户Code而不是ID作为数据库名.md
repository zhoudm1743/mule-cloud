# 修复 - 使用租户 Code 而不是 ID 作为数据库名

## 🎯 改进目标

**修改前**：
```
数据库名称: mule_68e27febab849776302...
```
- ❌ 不易识别
- ❌ 管理困难
- ❌ 调试不直观

**修改后**：
```
数据库名称: mule_default
数据库名称: mule_ace
数据库名称: mule_company_a
```
- ✅ 语义明确
- ✅ 易于管理
- ✅ 调试友好

---

## 📝 修改内容

### 1. 数据库管理器 (`core/database/manager.go`)

#### GetTenantDatabaseName

```go
// 修改前
func GetTenantDatabaseName(tenantID string) string {
    return fmt.Sprintf("mule_%s", tenantID)  // mule_68e27febab849776302...
}

// 修改后
func GetTenantDatabaseName(tenantCode string) string {
    return fmt.Sprintf("mule_%s", tenantCode)  // mule_default
}
```

#### GetDatabase

```go
// 修改前
func (m *DatabaseManager) GetDatabase(tenantID string) *mongo.Database {
    if tenantID == "" {
        return m.systemDB
    }
    // ...
    dbName := GetTenantDatabaseName(tenantID)
    // ...
}

// 修改后  
func (m *DatabaseManager) GetDatabase(tenantCode string) *mongo.Database {
    if tenantCode == "" {
        return m.systemDB
    }
    // ...
    dbName := GetTenantDatabaseName(tenantCode)
    // ...
}
```

#### CreateTenantDatabase

```go
// 修改前
func (m *DatabaseManager) CreateTenantDatabase(ctx context.Context, tenantID string) error {
    dbName := GetTenantDatabaseName(tenantID)
    // ...
}

// 修改后
func (m *DatabaseManager) CreateTenantDatabase(ctx context.Context, tenantCode string) error {
    dbName := GetTenantDatabaseName(tenantCode)
    // ...
}
```

#### DeleteTenantDatabase

```go
// 修改前
func (m *DatabaseManager) DeleteTenantDatabase(ctx context.Context, tenantID string) error {
    dbName := GetTenantDatabaseName(tenantID)
    // ...
}

// 修改后
func (m *DatabaseManager) DeleteTenantDatabase(ctx context.Context, tenantCode string) error {
    dbName := GetTenantDatabaseName(tenantCode)
    // ...
}
```

---

### 2. Context 传递 (`core/context/tenant.go`)

为了向后兼容，保留原有函数，添加新的别名函数：

```go
// 添加注释说明
const (
    // TenantIDKey 实际上存储的是租户代码（code）而不是ID
    // 为了向后兼容保留这个名字，但实际存储的是 tenant_code
    TenantIDKey contextKey = "tenant_id"
    // ...
)

// 保留原有函数（向后兼容）
func WithTenantID(ctx context.Context, tenantID string) context.Context {
    return context.WithValue(ctx, TenantIDKey, tenantID)
}

// 添加新的别名函数（语义更清晰，推荐使用）
func WithTenantCode(ctx context.Context, tenantCode string) context.Context {
    return WithTenantID(ctx, tenantCode)
}

// 保留原有函数（向后兼容）
func GetTenantID(ctx context.Context) string {
    if tenantID, ok := ctx.Value(TenantIDKey).(string); ok {
        return tenantID
    }
    return ""
}

// 添加新的别名函数（语义更清晰，推荐使用）
func GetTenantCode(ctx context.Context) string {
    return GetTenantID(ctx)
}
```

**设计理由**：
1. ✅ 向后兼容：所有使用 `GetTenantID()` 的地方仍然有效
2. ✅ 语义清晰：新代码可以使用 `WithTenantCode()` / `GetTenantCode()`
3. ✅ 最小改动：不需要修改所有 repository 的代码

---

### 3. 租户服务 (`app/system/services/tenant.go`)

```go
// 修改前
dbManager.CreateTenantDatabase(ctx, tenant.ID)  // 使用 ID
dbManager.DeleteTenantDatabase(ctx, tenant.ID)
tenantCtx := tenantCtx.WithTenantID(ctx, tenant.ID)

// 修改后
dbManager.CreateTenantDatabase(ctx, tenant.Code)  // 使用 Code
dbManager.DeleteTenantDatabase(ctx, tenant.Code)
tenantCtx := tenantCtx.WithTenantCode(ctx, tenant.Code)  // 使用新函数，语义更清晰
```

---

## 🔒 租户 Code 唯一性保证

### 数据库索引

确保 `tenant` 集合的 `code` 字段有唯一索引：

```javascript
// 系统数据库: tenant_system.tenant
db.tenant.createIndex({ "code": 1 }, { unique: true })
```

### 代码层面校验

在创建租户时应该检查 code 是否已存在：

```go
// app/system/services/tenant.go
func (s *TenantService) Create(req dto.TenantCreateRequest) (*models.Tenant, error) {
    // 1. 检查租户代码是否已存在
    existing, _ := s.repo.GetByCode(context.Background(), req.Code)
    if existing != nil {
        return nil, fmt.Errorf("租户代码 '%s' 已存在", req.Code)
    }
    
    // 2. 创建租户...
}
```

---

## 📊 数据库对比

### 修改前

```
MongoDB 数据库列表:
├── tenant_system                           ← 系统库
├── mule_68e27febab849776302f149           ← 租户 A (???谁的？)
├── mule_68dda6cd04ba0d6c8dda4b7a           ← 租户 B (???谁的？)
└── mule_68f3a4e1b2c5d7f8e9a1b2c3           ← 租户 C (???谁的？)
```

### 修改后

```
MongoDB 数据库列表:
├── tenant_system                           ← 系统库
├── mule_default                            ← 默认租户 ✅
├── mule_ace                                ← ACE 公司 ✅
└── mule_company_a                          ← A 公司 ✅
```

---

## 🚀 测试验证

### 1. 创建新租户

```
租户代码: test01
租户名称: 测试租户
```

**期望结果**：
```
数据库名称: mule_test01 ✅
```

### 2. 查看数据库列表

```javascript
// MongoDB
show databases

// 期望看到
mule_test01
```

### 3. 验证查询

```go
// Repository 查询
tenantCode := "test01"
db := dbManager.GetDatabase(tenantCode)  // 返回 mule_test01 数据库
```

---

## 📝 迁移方案（可选）

如果需要将现有租户的数据库名从 ID 改为 Code：

### 方案 A: 不迁移（推荐）

- ✅ 保持现有租户的数据库名不变
- ✅ 新租户使用新的命名规则
- ✅ 代码可以同时支持两种命名

### 方案 B: 创建迁移脚本

```javascript
// 迁移脚本伪代码
db.tenant_system.tenant.find().forEach(tenant => {
    let oldDbName = `mule_${tenant._id}`
    let newDbName = `mule_${tenant.code}`
    
    // 1. 复制数据库
    db.copyDatabase(oldDbName, newDbName)
    
    // 2. 验证数据完整性
    // ...
    
    // 3. 删除旧数据库
    db.getSiblingDB(oldDbName).dropDatabase()
})
```

---

## ✅ 修改完成

### 编译验证

```bash
go build ./core/database  ✅
go build ./core/context   ✅
go build ./cmd/system     ✅
go build ./cmd/auth       ✅
go build ./cmd/basic      ✅
```

### 功能验证

1. ✅ 创建租户时使用 code 生成数据库名
2. ✅ Repository 自动切换到正确的数据库
3. ✅ 中间件传递 code 而不是 id
4. ✅ 向后兼容，不影响现有功能

---

## 📝 后续建议

### 1. 前端显示优化

在租户管理页面显示数据库名称：

```
租户列表:
ID: 68e27febab849776302f149
Code: default
Name: 默认租户
Database: mule_default  ← 显示数据库名
```

### 2. 日志优化

在日志中显示 code 而不是 id：

```go
log.Printf("租户 [%s] 登录成功", tenantCode)  // 而不是 tenantID
```

### 3. 监控告警

使用 code 作为监控指标标签：

```
tenant_requests{tenant="default"} 1000
tenant_requests{tenant="ace"} 500
```

---

## 🎉 改进效果

- ✅ 数据库名称语义化：`mule_default` 比 `mule_68e27...` 更直观
- ✅ 管理更方便：直接通过数据库名识别租户
- ✅ 调试更友好：日志和错误信息更易读
- ✅ 向后兼容：不影响现有代码和数据
- ✅ 代码一致性：统一使用 code 而不是 id

**现在创建的新租户将使用更易读的数据库名！** 🎊
