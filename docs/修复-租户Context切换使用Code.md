# 修复 - 租户 Context 切换使用 Code

## 🐛 问题描述

### 现象

系统管理员在前端切换租户后（例如切换到 "ace"），页面数据无法加载：
- ✅ 前端正确发送了 `X-Tenant-Context` header
- ❌ 但发送的是租户 **ID**（`68e2928443a9eb9db480ed6b`）而不是租户 **Code**（`ace`）
- ❌ 后端 Repository 使用 `GetTenantCode` 获取上下文，但收到的是 ID
- ❌ 导致数据库切换失败，查询不到数据

### 截图证据

**请求头**：
```
X-Tenant-Context: 68e2928443a9eb9db480ed6b
```

**期望**：
```
X-Tenant-Context: ace
```

---

## 🔍 问题原因

### 数据流分析

```
前端 TenantSelector:
  选择租户 → value: tenant.id (❌ 使用了 ID)
    ↓
  保存到 localStorage: selected_tenant_id = "68e2928443a9eb9db480ed6b"
    ↓
  HTTP 请求 → X-Tenant-Context: "68e2928443a9eb9db480ed6b"
    ↓
后端中间件:
  TenantContextMiddleware → WithTenantID(ctx, "68e2928443a9eb9db480ed6b")
    ↓
Repository:
  GetTenantCode(ctx) → 返回 "68e2928443a9eb9db480ed6b" (❌ 这是 ID 不是 Code)
    ↓
DatabaseManager:
  GetDatabase("68e2928443a9eb9db480ed6b") → 查询数据库 "mule_68e2928443a9eb9db480ed6b" (❌ 错误)
    ↓
结果: 数据查询失败
```

### 根本原因

1. **前端保存和发送的是 ID 而不是 Code**
   - `TenantSelector` 组件中 `value: tenant.id`
   - 应该使用 `value: tenant.code`

2. **中间件使用旧的 API**
   - `TenantContextMiddleware` 使用 `WithTenantID`
   - 应该使用 `WithTenantCode`

3. **Repository 期望 Code**
   - 所有 Repository 都改成了 `GetTenantCode(ctx)`
   - 但 Context 中存储的仍然是 ID

---

## ✅ 解决方案

### 1. 修改前端 - 使用 Code 而不是 ID

#### `frontend/src/components/TenantSelector.vue`

```typescript
// ❌ 错误：使用 ID
...data.tenants.map(tenant => ({
  label: `${tenant.name} (${tenant.code})`,
  value: tenant.id, // ❌
}))

// ✅ 正确：使用 Code
...data.tenants.map(tenant => ({
  label: `${tenant.name} (${tenant.code})`,
  value: tenant.code, // ✅
}))

// ❌ 错误：保存 ID
function onTenantChange(value: string) {
  local.set('selected_tenant_id', value) // ❌
  window.location.reload()
}

// ✅ 正确：保存 Code
function onTenantChange(value: string) {
  local.set('selected_tenant_code', value) // ✅
  window.location.reload()
}

// ❌ 错误：恢复 ID
function restoreSelection() {
  const savedTenantId = local.get('selected_tenant_id') // ❌
  if (savedTenantId) {
    selectedTenantId.value = savedTenantId
  }
}

// ✅ 正确：恢复 Code
function restoreSelection() {
  const savedTenantCode = local.get('selected_tenant_code') // ✅
  if (savedTenantCode) {
    selectedTenantId.value = savedTenantCode
  }
}
```

---

#### `frontend/src/service/http/alova.ts`

```typescript
// ❌ 错误：发送 ID
const userInfo = local.get('userInfo')
const selectedTenantId = local.get('selected_tenant_id') // ❌

if (userInfo && !userInfo.tenant_id && selectedTenantId) {
  method.config.headers['X-Tenant-Context'] = selectedTenantId // ❌
}

// ✅ 正确：发送 Code
const userInfo = local.get('userInfo')
const selectedTenantCode = local.get('selected_tenant_code') // ✅

if (userInfo && !userInfo.tenant_id && selectedTenantCode) {
  method.config.headers['X-Tenant-Context'] = selectedTenantCode // ✅
}
```

---

#### `frontend/src/store/auth.ts`

```typescript
// ❌ 错误：清除旧的 key
local.remove('selected_tenant_id') // ❌

// ✅ 正确：清除新的 key
local.remove('selected_tenant_code') // ✅
```

---

### 2. 修改后端中间件 - 使用 Code API

#### `core/middleware/tenant_context.go`

```go
// ❌ 错误：使用 ID API
func TenantContextMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		currentTenantID := tenantCtx.GetTenantID(ctx) // ❌

		if currentTenantID == "" { // ❌
			contextTenantID := c.GetHeader("X-Tenant-Context")

			if contextTenantID != "" {
				// ...
				ctx = tenantCtx.WithTenantID(ctx, contextTenantID) // ❌
				c.Request = c.Request.WithContext(ctx)
			}
		}

		c.Next()
	}
}

// ✅ 正确：使用 Code API
func TenantContextMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		currentTenantCode := tenantCtx.GetTenantCode(ctx) // ✅

		// ✅ 系统管理员（tenantCode为空或"system"）可以切换
		if currentTenantCode == "" || currentTenantCode == "system" {
			contextTenantCode := c.GetHeader("X-Tenant-Context") // ✅ 现在是 Code

			if contextTenantCode != "" {
				// ...
				ctx = tenantCtx.WithTenantCode(ctx, contextTenantCode) // ✅
				c.Request = c.Request.WithContext(ctx)
			}
		}

		c.Next()
	}
}
```

---

## 📊 修复后的数据流

```
前端 TenantSelector:
  选择租户 → value: tenant.code (✅ "ace")
    ↓
  保存到 localStorage: selected_tenant_code = "ace" ✅
    ↓
  HTTP 请求 → X-Tenant-Context: "ace" ✅
    ↓
后端中间件:
  TenantContextMiddleware → WithTenantCode(ctx, "ace") ✅
    ↓
Repository:
  GetTenantCode(ctx) → 返回 "ace" ✅
    ↓
DatabaseManager:
  GetDatabase("ace") → 查询数据库 "mule_ace" ✅
    ↓
结果: 数据查询成功 ✅
```

---

## 🎯 修改文件清单

### 前端
- ✅ `frontend/src/components/TenantSelector.vue` - 使用 `tenant.code` 和 `selected_tenant_code`
- ✅ `frontend/src/service/http/alova.ts` - 从 localStorage 读取 `selected_tenant_code`
- ✅ `frontend/src/store/auth.ts` - 清除和恢复使用 `selected_tenant_code`

### 后端
- ✅ `core/middleware/tenant_context.go` - 使用 `GetTenantCode` 和 `WithTenantCode`

---

## ✅ 编译验证

```bash
# 后端
go build ./core/middleware  ✅
go build ./cmd/perms        ✅
go build ./cmd/auth         ✅
go build ./cmd/basic        ✅

# 前端
npm run build               ✅
```

---

## 🔍 验证方法

### 1. 清除旧数据

```javascript
// 浏览器 Console
localStorage.removeItem('selected_tenant_id')  // 删除旧的 key
```

### 2. 重新登录

- 使用系统管理员账号登录
- 选择租户 "ACE租户 (ace)"
- 查看浏览器 Network 请求头

**期望看到**：
```
X-Tenant-Context: ace
```

### 3. 检查数据加载

- 切换到岗位管理页面
- 应该能看到租户的岗位数据
- Network 请求应该返回数据而不是空列表

### 4. 查看日志

```
[租户上下文切换] 系统管理员切换到租户: ace
[MongoDB] 使用数据库: mule_ace
```

---

## 📝 localStorage Key 变更

| 旧 Key | 新 Key | 说明 |
|--------|--------|------|
| `selected_tenant_id` | `selected_tenant_code` | 存储租户 Code 而不是 ID |

---

## 🎉 修复完成

### 效果

- ✅ 系统管理员可以切换到指定租户
- ✅ 切换后可以查看和管理租户数据
- ✅ 数据库正确切换到 `mule_<code>`
- ✅ 前后端数据流一致

---

## 🔗 相关文档

- [租户 Code 数据库命名方案](完成-租户Code数据库命名方案.md)
- [修复-租户管理员菜单获取问题](修复-租户管理员菜单获取问题.md)
- [重要-租户Code使用说明](重要-租户Code使用说明.md)

---

**现在系统管理员可以正确切换租户并查看租户数据了！** 🎊

