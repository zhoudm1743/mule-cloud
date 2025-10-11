# 修复：微信小程序 UnionID 为空导致登录失败问题

## 问题描述

用户登录时报错：
```json
{
  "code": -1,
  "msg": "查询成员信息失败: %!w(<nil>)",
  "timestamp": 1760197427
}
```

## 问题原因

在 `app/miniapp/services/wechat.go` 的 `generateTokenForTenant` 函数中，查询租户成员信息时只使用了 `UnionID`：

```go
member, err := s.memberRepo.GetByUnionID(tenantContext, user.UnionID)
```

**问题关键**：
- 如果小程序**没有配置微信开放平台**，`UnionID` 会是**空字符串**
- 使用空字符串查询 MongoDB 会返回 `nil`（找不到记录）
- 导致用户无法正常登录

## 解决方案

修改 `generateTokenForTenant` 函数，增加 `UserID` 作为备用查询条件：

```go
// 查询租户成员信息（获取角色）
tenantContext := tenantCtx.WithTenantCode(ctx, tenantMap.TenantCode)

// 优先使用UnionID查询，如果没有则使用UserID
var member *models.TenantMember
if user.UnionID != "" {
    member, err = s.memberRepo.GetByUnionID(tenantContext, user.UnionID)
}

// 如果通过UnionID找不到或UnionID为空，尝试通过UserID查询
if member == nil {
    member, err = s.memberRepo.GetByUserID(tenantContext, user.ID)
}

if err != nil || member == nil {
    logger.Error("查询租户成员失败",
        zap.String("tenant_code", tenantMap.TenantCode),
        zap.String("user_id", user.ID),
        zap.String("union_id", user.UnionID),
        zap.Error(err))
    return "", nil, fmt.Errorf("查询成员信息失败: %w", err)
}
```

## 技术细节

### 1. UnionID vs OpenID

| 字段 | 说明 | 适用场景 |
|------|------|----------|
| **OpenID** | 用户在单个小程序中的唯一标识 | 未绑定开放平台的小程序 |
| **UnionID** | 用户在同一开放平台下所有应用的唯一标识 | 已绑定开放平台的小程序，用于打通多个应用 |

### 2. 数据库设计

在 `tenant_xxx.member` 集合中，有两个字段用于关联用户：

```go
type TenantMember struct {
    UnionID string `bson:"union_id" json:"union_id"` // 微信UnionID（如果有）
    UserID  string `bson:"user_id" json:"user_id"`   // 全局用户ID（必有）
    // ... 其他字段
}
```

### 3. 查询优先级

1. **优先使用 UnionID**：如果有开放平台，UnionID 更稳定
2. **备用 UserID**：兼容没有开放平台的场景

## 影响范围

此修复影响以下场景：
- ✅ **微信小程序登录**（只有一个租户）
- ✅ **微信小程序登录**（多个租户，选择后）
- ✅ **切换租户**
- ✅ **绑定租户**

## 测试验证

### 测试场景 1：无开放平台的小程序

1. 用户首次登录（UnionID 为空）
2. 绑定租户
3. 登录成功 ✅

### 测试场景 2：有开放平台的小程序

1. 用户首次登录（UnionID 有值）
2. 绑定租户
3. 登录成功 ✅

### 测试场景 3：用户有多个租户

1. 用户登录
2. 选择租户
3. 登录成功 ✅

### 测试场景 4：切换租户

1. 已登录用户
2. 切换到另一个租户
3. 切换成功 ✅

## 注意事项

1. **数据一致性**：在创建成员记录时（`BindTenant` 函数），需要同时保存 `UnionID` 和 `UserID`
2. **查询性能**：建议在 `member` 集合的 `user_id` 和 `union_id` 字段上创建索引
3. **日志记录**：增加了详细的错误日志，便于排查问题

## 相关文件

- `app/miniapp/services/wechat.go` - 微信服务实现
- `internal/repository/tenant_member.go` - 租户成员数据仓库
- `internal/models/tenant_member.go` - 租户成员模型

## 相关文档

- [微信小程序多租户用户系统设计方案](./微信小程序多租户用户系统设计方案.md)
- [微信小程序用户信息与手机号绑定功能说明](./微信小程序用户信息与手机号绑定功能说明.md)

