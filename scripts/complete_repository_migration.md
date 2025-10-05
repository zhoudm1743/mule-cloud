# Repository层改造 - 批量完成指南

## 🎯 目标

完成以下Repository的改造：
- ✅ `admin.go` - 已完成
- ⏳ `role.go` - 进行中（90%完成）
- ⏳ `menu.go` - 待改造
- ⏳ `basic.go` - 待改造  
- ⏳ `tenant.go` - 待改造（特殊）

## 📝 role.go 剩余工作

role.go已经完成了主要改造，剩余需要添加 `collection := r.getCollection(ctx)` 的方法：

### 需要添加的方法

```go
// FindWithPage
func (r *roleRepository) FindWithPage(ctx context.Context, filter bson.M, page, pageSize int64) ([]*models.Role, error) {
    collection := r.getCollection(ctx)  // ← 添加这一行
    cursor, err := collection.Find(ctx, filter)
    //...
}

// Count
func (r *roleRepository) Count(ctx context.Context, filter bson.M) (int64, error) {
    collection := r.getCollection(ctx)  // ← 添加这一行
    count, err := collection.CountDocuments(ctx, filter)
    //...
}

// Update
func (r *roleRepository) Update(ctx context.Context, id string, update bson.M) error {
    collection := r.getCollection(ctx)  // ← 添加这一行在方法开头
    objectID, err := bson.ObjectIDFromHex(id)
    //...
}

// UpdateOne
func (r *roleRepository) UpdateOne(ctx context.Context, filter bson.M, update bson.M) error {
    collection := r.getCollection(ctx)  // ← 添加这一行
    updateDoc := bson.M{"$set": update}
    //...
}

// Delete
func (r *roleRepository) Delete(ctx context.Context, id string) error {
    collection := r.getCollection(ctx)  // ← 添加这一行在方法开头
    objectID, err := bson.ObjectIDFromHex(id)
    //...
}

// HardDelete
func (r *roleRepository) HardDelete(ctx context.Context, id string) error {
    collection := r.getCollection(ctx)  // ← 添加这一行在方法开头
    objectID, err := bson.ObjectIDFromHex(id)
    //...
}

// DeleteMany
func (r *roleRepository) DeleteMany(ctx context.Context, filter bson.M) (int64, error) {
    collection := r.getCollection(ctx)  // ← 添加这一行
    result, err := collection.DeleteMany(ctx, filter)
    //...
}

// FindDeletedWithPage
func (r *roleRepository) FindDeletedWithPage(ctx context.Context, filter bson.M, page, pageSize int64) ([]*models.Role, error) {
    collection := r.getCollection(ctx)  // ← 添加这一行
    filter["is_deleted"] = 1
    //...
}

// CountDeleted
func (r *roleRepository) CountDeleted(ctx context.Context, filter bson.M) (int64, error) {
    collection := r.getCollection(ctx)  // ← 添加这一行
    filter["is_deleted"] = 1
    //...
}

// RestoreOne
func (r *roleRepository) RestoreOne(ctx context.Context, id string) error {
    collection := r.getCollection(ctx)  // ← 添加这一行在方法开头
    objectID, err := bson.ObjectIDFromHex(id)
    //...
}

// RestoreMany
func (r *roleRepository) RestoreMany(ctx context.Context, ids []string) (int64, error) {
    collection := r.getCollection(ctx)  // ← 添加这一行
    if len(ids) == 0 {
        return 0, nil
    }
    //...
}

// HardDeleteMany
func (r *roleRepository) HardDeleteMany(ctx context.Context, ids []string) (int64, error) {
    collection := r.getCollection(ctx)  // ← 添加这一行
    if len(ids) == 0 {
        return 0, nil
    }
    //...
}

// GetRolesByIDs
func (r *roleRepository) GetRolesByIDs(ctx context.Context, ids []string) ([]*models.Role, error) {
    collection := r.getCollection(ctx)  // ← 添加这一行
    if len(ids) == 0 {
        return []*models.Role{}, nil
    }
    //...
}

// GetAllRoles (原 GetRolesByTenant)
func (r *roleRepository) GetAllRoles(ctx context.Context) ([]*models.Role, error) {
    collection := r.getCollection(ctx)  // ← 添加这一行
    filter := bson.M{"is_deleted": 0}  // ← 删除 tenant_id 过滤
    cursor, err := collection.Find(ctx, filter)
    //...
}

// GetCollection
func (r *roleRepository) GetCollection() *mongo.Collection {
    return r.dbManager.GetSystemDatabase().Collection("role")
}
```

### 快速批量修改方法

#### 方法1：使用编辑器查找替换

1. 打开 `role.go`
2. 查找所有 `func (r \*roleRepository\) (\w+)\(ctx context\.Context`
3. 在函数体第一行添加 `collection := r.getCollection(ctx)`

#### 方法2：手动逐个方法添加

按照上面的列表，在每个方法的第一行添加 `collection := r.getCollection(ctx)`

---

## 📋 完整改造步骤（适用于其他Repository）

### Step 1: 修改导入

```go
import (
    "context"
    tenantCtx "mule-cloud/core/context"  // ← 新增
    "mule-cloud/core/database"
    //...
)
```

### Step 2: 修改结构体

```go
// ❌ 旧代码
type xxxRepository struct {
    collection *mongo.Collection
}

func NewXXXRepository() XXXRepository {
    collection := database.MongoDB.Collection("xxx")
    return &xxxRepository{collection: collection}
}

// ✅ 新代码
type xxxRepository struct {
    dbManager *database.DatabaseManager
}

func NewXXXRepository() XXXRepository {
    return &xxxRepository{
        dbManager: database.GetDatabaseManager(),
    }
}

func (r *xxxRepository) getCollection(ctx context.Context) *mongo.Collection {
    tenantID := tenantCtx.GetTenantID(ctx)
    db := r.dbManager.GetDatabase(tenantID)
    return db.Collection("xxx")
}
```

### Step 3: 修改所有方法

在每个方法的**第一行**添加：
```go
collection := r.getCollection(ctx)
```

然后确保所有 `r.collection.XXX` 都改为 `collection.XXX`

### Step 4: 修复GetCollection方法

```go
func (r *xxxRepository) GetCollection() *mongo.Collection {
    return r.dbManager.GetSystemDatabase().Collection("xxx")
}
```

---

## 🔧 特殊处理：tenant.go

`tenant.go` 固定使用系统数据库：

```go
type tenantRepository struct {
    dbManager *database.DatabaseManager
}

func NewTenantRepository() TenantRepository {
    return &tenantRepository{
        dbManager: database.GetDatabaseManager(),
    }
}

// ⚠️ 注意：不需要ctx参数，直接使用系统库
func (r *tenantRepository) getCollection() *mongo.Collection {
    db := r.dbManager.GetSystemDatabase()
    return db.Collection("tenant")
}

// 所有方法调用
func (r *tenantRepository) Get(ctx context.Context, id string) (*models.Tenant, error) {
    collection := r.getCollection()  // ← 不需要ctx
    //...
}
```

---

## ✅ 改造检查清单

### role.go
- [x] 修改导入
- [x] 修改结构体和构造函数
- [x] 添加 getCollection 方法
- [x] Get 方法
- [x] GetByCode 方法  
- [x] GetByName 方法
- [x] Find 方法
- [x] FindOne 方法
- [x] Create 方法
- [ ] FindWithPage 方法
- [ ] Count 方法
- [ ] Update 方法
- [ ] UpdateOne 方法
- [ ] Delete 方法
- [ ] HardDelete 方法
- [ ] DeleteMany 方法
- [ ] FindDeletedWithPage 方法
- [ ] CountDeleted 方法
- [ ] RestoreOne 方法
- [ ] RestoreMany 方法
- [ ] HardDeleteMany 方法
- [ ] GetRolesByIDs 方法
- [ ] GetAllRoles 方法（原GetRolesByTenant）
- [ ] GetCollection 方法

### menu.go
- [ ] 完整改造

### basic.go
- [ ] 完整改造

### tenant.go
- [ ] 完整改造（特殊：固定使用系统库）

---

## 🎯 预期完成时间

- role.go 剩余工作：10分钟
- menu.go：20分钟
- basic.go：20分钟
- tenant.go：15分钟

**总计：** 约1小时

---

## 🧪 验证改造

改造完成后，确保：

1. 所有方法开头都有 `collection := r.getCollection(ctx)`
2. 所有 `r.collection.XXX` 都改为 `collection.XXX`
3. 接口中的方法签名更新（删除tenantID参数）
4. 编译通过：`go build ./internal/repository/...`

```bash
# 验证编译
cd K:\Git\mule-cloud
go build ./internal/repository/...

# 如果有错误，会显示具体的问题
```

---

**完成这些改造后，整个Repository层就改造完成了！**

