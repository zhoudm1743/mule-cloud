# MongoDB 嵌套 `$set` 错误修复记录

## 📅 修复时间
2025年10月19日

## 🐛 问题描述

### 错误信息
```
write exception: write errors: [The dollar ($) prefixed field '$set' in '$set' is not allowed in the context of an update's replacement document. Consider using an aggregation pipeline with $replaceWith.]
```

### 错误原因

在 `internal/repository/order.go` 中，`Update` 方法已经自动包装了 `$set`：

```go
func (r *orderRepository) Update(ctx context.Context, id string, update bson.M) error {
    collection := r.GetCollectionWithContext(ctx)
    
    objectID, err := bson.ObjectIDFromHex(id)
    if err != nil {
        return err
    }
    
    _, err = collection.UpdateOne(
        ctx,
        bson.M{"_id": objectID, "is_deleted": 0},
        bson.M{"$set": update},  // ← 这里已经包装了 $set
    )
    
    return err
}
```

但是在多处调用时，又嵌套了一层 `$set`，导致最终生成的 MongoDB 命令变成：

```json
{
  "$set": {
    "$set": {
      "progress": 0.05,
      "updated_at": 1760840655
    }
  }
}
```

MongoDB 不允许这种嵌套的 `$set` 操作符。

## 🔍 受影响的文件

共修复了 **3 个文件**：

### 1. `app/production/services/report.go`

**位置**：`updateOrderProgressFromPieces` 方法

**修复前**：
```go
// 3. 更新订单进度字段
err = s.orderRepo.Update(ctx, orderID, map[string]interface{}{
    "$set": map[string]interface{}{  // ❌ 嵌套的 $set
        "progress":   orderProgress,
        "updated_at": time.Now().Unix(),
    },
})
```

**修复后**：
```go
// 3. 更新订单进度字段
// 注意：orderRepo.Update 方法内部会自动包装 $set，这里直接传字段即可
err = s.orderRepo.Update(ctx, orderID, bson.M{
    "progress":   orderProgress,
    "updated_at": time.Now().Unix(),
})
```

### 2. `core/workflow/order_workflow.go`

**位置**：`TransitionTo` 方法

**修复前**：
```go
// 更新数据库
err = w.orderRepo.Update(ctx, orderID, bson.M{
    "$set": bson.M{  // ❌ 嵌套的 $set
        "status":     int(nextStatus),
        "updated_at": time.Now().Unix(),
    },
})
```

**修复后**：
```go
// 更新数据库
// 注意：orderRepo.Update 方法内部会自动包装 $set，这里直接传字段即可
err = w.orderRepo.Update(ctx, orderID, bson.M{
    "status":     int(nextStatus),
    "updated_at": time.Now().Unix(),
})
```

### 3. `core/workflow/order_workflow_advanced.go`

**位置**：`TransitionToAdvanced` 方法

**修复前**：
```go
// 更新数据库
err = w.orderRepo.Update(ctx, orderID, map[string]interface{}{
    "$set": map[string]interface{}{  // ❌ 嵌套的 $set
        "status":     int(nextStatus),
        "updated_at": time.Now().Unix(),
    },
})
```

**修复后**：
```go
// 更新数据库
// 注意：orderRepo.Update 方法内部会自动包装 $set，这里直接传字段即可
err = w.orderRepo.Update(ctx, orderID, bson.M{
    "status":     int(nextStatus),
    "updated_at": time.Now().Unix(),
})
```

**额外修复**：添加 `bson` 导入

```go
import (
    "context"
    "encoding/json"
    "fmt"
    "time"

    "mule-cloud/internal/repository"

    "go.mongodb.org/mongo-driver/v2/bson"  // ← 新增
)
```

## ✅ 修复验证

### 编译测试
```bash
# Production Service
go build -o ./bin/production.exe ./cmd/production
# ✅ 编译成功

# Order Service
go build -o ./bin/order.exe ./cmd/order
# ✅ 编译成功
```

### 日志验证

修复后，MongoDB 命令应该变成：

```json
{
  "update": "orders",
  "updates": [{
    "q": {"_id": {"$oid": "..."}},
    "u": {
      "$set": {
        "progress": 0.05,
        "updated_at": 1760840655
      }
    }
  }]
}
```

## 📚 最佳实践

### ⚠️ 规则：不要在调用时包装 `$set`

当使用 Repository 层的 `Update` 方法时：

**❌ 错误做法**：
```go
repo.Update(ctx, id, bson.M{
    "$set": bson.M{  // ← 不要手动包装
        "field": value,
    },
})
```

**✅ 正确做法**：
```go
repo.Update(ctx, id, bson.M{
    "field": value,  // ← 直接传字段
})
```

### 📝 添加注释说明

为了避免未来再次犯同样的错误，建议在所有调用 `Update` 的地方添加注释：

```go
// 注意：orderRepo.Update 方法内部会自动包装 $set，这里直接传字段即可
err = s.orderRepo.Update(ctx, orderID, bson.M{
    "progress":   orderProgress,
    "updated_at": time.Now().Unix(),
})
```

### 🔍 检查清单

在使用 Repository 的 `Update` 方法时，检查：

1. ✅ 是否阅读了 Repository 层的 `Update` 方法实现？
2. ✅ 是否了解该方法是否已经包装了 `$set`？
3. ✅ 是否避免了手动添加 `$set` 操作符？
4. ✅ 是否添加了注释说明？

## 🔄 其他 Repository 方法

### 检查是否有类似问题

建议检查其他 Repository 的 `Update` 方法，确认它们的实现方式：

```bash
# 查找所有 Update 方法
grep -r "func.*Update.*bson.M" internal/repository/
```

如果其他 Repository 也使用了相同的模式（内部包装 `$set`），需要检查所有调用处。

## 📊 影响范围

| 功能模块 | 受影响方法 | 修复状态 |
|---------|----------|---------|
| 生产上报 | `updateOrderProgressFromPieces` | ✅ 已修复 |
| 工作流 - 基础 | `TransitionTo` | ✅ 已修复 |
| 工作流 - 高级 | `TransitionToAdvanced` | ✅ 已修复 |

## 🚀 部署步骤

1. **重新编译服务**
   ```bash
   # 编译 Production Service
   go build -o ./bin/production.exe ./cmd/production
   
   # 编译 Order Service（如果工作流相关功能在此）
   go build -o ./bin/order.exe ./cmd/order
   ```

2. **重启服务**
   ```bash
   # 停止 Production Service
   # 启动 Production Service
   
   # 停止 Order Service
   # 启动 Order Service
   ```

3. **验证功能**
   - 扫码上报工序
   - 检查订单进度是否更新
   - 检查订单状态是否自动转换

4. **查看日志**
   ```
   🚀 触发订单进度更新: 订单=xxx, 租户=xxx
   📊 订单进度计算（基于裁片）: ...
   ✅ 订单 xxx 已转换到生产中状态  ← 应该看到这个
   ```

## 🎯 预期结果

修复后：
- ✅ 订单进度正常更新
- ✅ 订单状态自动转换（已下单 → 生产中）
- ✅ 不再出现 MongoDB `$set` 嵌套错误

## 📝 相关文档

1. 《工作流自动更新修复说明》 - `docs/工作流自动更新修复说明.md`
2. 《扫码打菲工序上报实施方案》 - `docs/扫码打菲工序上报实施方案.md`

---

**修复完成！** 🎉

所有 MongoDB 嵌套 `$set` 错误已修复，服务编译成功，可以重新部署测试。

