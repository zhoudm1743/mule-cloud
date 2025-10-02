# Casbin 增删改查权限实现说明

## 当前实现分析

### 1. Casbin 模型 (model.conf)

```ini
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub) && keyMatch2(r.obj, p.obj) && regexMatch(r.act, p.act)
```

**说明：**
- `sub`: 主体（用户/角色）
- `obj`: 对象（资源路径）
- `act`: 动作（read/write）
- `g`: 角色继承关系

### 2. 当前权限映射

```go
// getActionFromMethod - 当前实现
func getActionFromMethod(method string) string {
    switch strings.ToUpper(method) {
    case "GET", "HEAD", "OPTIONS":
        return "read"      // 查询
    case "POST", "PUT", "PATCH", "DELETE":
        return "write"     // 增、改、删 都是 write
    default:
        return "read"
    }
}
```

**问题：**
- ❌ 只有 read 和 write 两种权限
- ❌ 无法区分 增(POST)、改(PUT)、删(DELETE)
- ❌ 无法细粒度控制（如：允许创建但不允许删除）

## 改进方案

### 方案A：细化动作权限（推荐）

#### 1. 修改权限映射

```go
// getActionFromMethod - 改进版
func getActionFromMethod(method string) string {
    switch strings.ToUpper(method) {
    case "GET", "HEAD":
        return "read"      // 查询
    case "POST":
        return "create"    // 创建
    case "PUT", "PATCH":
        return "update"    // 更新
    case "DELETE":
        return "delete"    // 删除
    case "OPTIONS":
        return "*"         // 预检请求放行
    default:
        return "read"
    }
}
```

#### 2. 权限策略示例

```bash
# 角色 admin 对 /system/users 的权限
p, role:admin, /system/users, read      # 可以查看用户列表
p, role:admin, /system/users, create    # 可以创建用户
p, role:admin, /system/users, update    # 可以更新用户
p, role:admin, /system/users, delete    # 可以删除用户

# 角色 viewer 对 /system/users 的权限（只读）
p, role:viewer, /system/users, read     # 只能查看

# 角色 editor 对 /system/users 的权限（可改不可删）
p, role:editor, /system/users, read
p, role:editor, /system/users, create
p, role:editor, /system/users, update
# 没有 delete 权限
```

#### 3. 修改权限同步函数

```go
// SyncRoleMenus - 改进版
func SyncRoleMenus(tenantID, roleID string, menuPaths []string, permissions map[string][]string) error {
    roleSub := fmt.Sprintf("tenant:%s:role:%s", tenantID, roleID)
    
    // 删除角色的所有现有权限
    Enforcer.RemoveFilteredPolicy(0, roleSub)
    
    // 批量添加新权限
    for _, menuPath := range menuPaths {
        // 从 permissions 获取该菜单的具体权限
        actions := permissions[menuPath]
        if len(actions) == 0 {
            // 默认只给 read 权限
            actions = []string{"read"}
        }
        
        for _, action := range actions {
            Enforcer.AddPolicy(roleSub, menuPath, action)
        }
    }
    
    return Enforcer.SavePolicy()
}
```

#### 4. 数据结构调整

```go
// Role 模型
type Role struct {
    Code  string                   `json:"code"`
    Name  string                   `json:"name"`
    Menus []string                 `json:"menus"`  // 菜单列表
    
    // 新增：菜单权限映射
    MenuPermissions map[string][]string `json:"menu_permissions"`
    // 示例：
    // {
    //   "/system/admin": ["read", "create", "update", "delete"],
    //   "/system/role": ["read", "create", "update"],
    //   "/system/tenant": ["read"]
    // }
}
```

### 方案B：REST 风格细粒度权限

#### 1. 更详细的权限映射

```go
func getActionFromMethod(method, path string) string {
    switch strings.ToUpper(method) {
    case "GET":
        // 区分列表查询和单个查询
        if strings.Contains(path, "/:id") || strings.HasSuffix(path, "/"+extractID(path)) {
            return "read:one"   // 查询单个
        }
        return "read:list"      // 查询列表
    case "POST":
        if strings.HasSuffix(path, "/batch") {
            return "create:batch"  // 批量创建
        }
        return "create"
    case "PUT":
        return "update"
    case "PATCH":
        return "update:partial"  // 部分更新
    case "DELETE":
        if strings.HasSuffix(path, "/batch") {
            return "delete:batch"  // 批量删除
        }
        return "delete"
    default:
        return "read"
    }
}
```

#### 2. 策略示例

```bash
# 完整权限
p, role:admin, /system/users, read:list
p, role:admin, /system/users, read:one
p, role:admin, /system/users, create
p, role:admin, /system/users, create:batch
p, role:admin, /system/users, update
p, role:admin, /system/users, update:partial
p, role:admin, /system/users, delete
p, role:admin, /system/users, delete:batch

# 受限权限（只能查看和单个更新）
p, role:limited, /system/users, read:list
p, role:limited, /system/users, read:one
p, role:limited, /system/users, update:partial
```

### 方案C：通配符权限（简化版）

#### 1. 使用通配符简化配置

```bash
# 使用 * 表示所有操作
p, role:admin, /system/*, *                    # 所有权限

# 使用前缀匹配
p, role:viewer, /system/*, read*               # 所有读权限（read, read:list, read:one）

# 组合权限
p, role:editor, /system/users, (read|create|update)  # 读、创建、更新
```

#### 2. 修改 Matcher

```ini
[matchers]
# 支持通配符和正则表达式
m = g(r.sub, p.sub) && keyMatch2(r.obj, p.obj) && (p.act == "*" || regexMatch(r.act, p.act))
```

## 实施步骤

### 阶段1：基础四权限（推荐先实施）

1. **修改权限映射**
```go
// app/gateway/middleware/casbin_auth.go
func getActionFromMethod(method string) string {
    switch strings.ToUpper(method) {
    case "GET", "HEAD":
        return "read"
    case "POST":
        return "create"
    case "PUT", "PATCH":
        return "update"
    case "DELETE":
        return "delete"
    default:
        return "read"
    }
}
```

2. **更新权限同步逻辑**
```go
// core/casbin/casbin.go
func SyncRoleMenusWithActions(tenantID, roleID string, menuPermissions map[string][]string) error {
    roleSub := fmt.Sprintf("tenant:%s:role:%s", tenantID, roleID)
    
    // 删除旧权限
    Enforcer.RemoveFilteredPolicy(0, roleSub)
    
    // 添加新权限
    for menuPath, actions := range menuPermissions {
        for _, action := range actions {
            Enforcer.AddPolicy(roleSub, menuPath, action)
        }
    }
    
    return Enforcer.SavePolicy()
}
```

3. **前端界面调整**

角色分配菜单时，不仅选择菜单，还选择权限：

```typescript
interface MenuPermission {
  menuName: string
  permissions: {
    read: boolean
    create: boolean
    update: boolean
    delete: boolean
  }
}
```

界面示例：
```
菜单选择：
☑️ 管理员管理
   ☑️ 查看  ☑️ 创建  ☑️ 修改  ☑️ 删除
   
☑️ 角色管理
   ☑️ 查看  ☑️ 创建  ☑️ 修改  ☐ 删除
   
☑️ 租户管理
   ☑️ 查看  ☐ 创建  ☐ 修改  ☐ 删除
```

### 阶段2：路径级权限控制

```go
// 支持更细粒度的路径控制
func SyncRoleAPIs(tenantID, roleID string, apiPermissions []APIPermission) error {
    roleSub := fmt.Sprintf("tenant:%s:role:%s", tenantID, roleID)
    
    for _, api := range apiPermissions {
        Enforcer.AddPolicy(roleSub, api.Path, api.Action)
    }
    
    return Enforcer.SavePolicy()
}

type APIPermission struct {
    Path   string   // /system/users/:id
    Action string   // read, create, update, delete
}
```

### 阶段3：字段级权限（高级）

```go
// 支持字段级别的权限控制
type FieldPermission struct {
    Resource string
    Field    string
    Action   string  // read, write
}

// 示例：
// 普通管理员可以查看用户信息，但不能查看敏感字段
p, role:admin, /system/users, read
p, role:admin, /system/users.password, deny:read
p, role:admin, /system/users.salary, deny:read
```

## 配置示例

### 1. 基础配置

```bash
# 超级管理员 - 所有权限
p, super:admin, /system/*, *

# 系统管理员 - 增删改查
p, role:system_admin, /system/admins, read
p, role:system_admin, /system/admins, create
p, role:system_admin, /system/admins, update
p, role:system_admin, /system/admins, delete

p, role:system_admin, /system/roles, read
p, role:system_admin, /system/roles, create
p, role:system_admin, /system/roles, update
p, role:system_admin, /system/roles, delete

# 普通管理员 - 只读 + 部分编辑
p, role:normal_admin, /system/admins, read
p, role:normal_admin, /system/admins, update

# 只读用户 - 只能查看
p, role:viewer, /system/*, read
```

### 2. 租户级权限

```bash
# 租户A的管理员角色
p, tenant:A:role:admin, /system/users, read
p, tenant:A:role:admin, /system/users, create
p, tenant:A:role:admin, /system/users, update
p, tenant:A:role:admin, /system/users, delete

# 租户A的普通用户
p, tenant:A:role:user, /system/users, read

# 用户角色绑定
g, tenant:A:user:123, tenant:A:role:admin
g, tenant:A:user:456, tenant:A:role:user
```

## 最佳实践

### 1. 权限设计原则

- **最小权限原则**：默认没有任何权限，需要明确授予
- **分层控制**：租户 → 角色 → 用户
- **可撤销**：所有权限都应该可以撤销

### 2. 权限粒度选择

| 场景 | 推荐粒度 | 示例 |
|------|---------|------|
| 后台管理系统 | 菜单 + CRUD | `/system/users` + `read/create/update/delete` |
| API 服务 | 路径 + 方法 | `/api/v1/users/:id` + `GET/POST/PUT/DELETE` |
| 多租户SaaS | 租户 + 角色 + 资源 | `tenant:A:role:admin` + `/resources` |

### 3. 性能优化

```go
// 使用缓存减少数据库查询
var permissionCache = cache.New(5*time.Minute, 10*time.Minute)

func CheckPermissionWithCache(sub, obj, act string) (bool, error) {
    key := fmt.Sprintf("%s:%s:%s", sub, obj, act)
    
    if cached, found := permissionCache.Get(key); found {
        return cached.(bool), nil
    }
    
    allowed, err := Enforcer.Enforce(sub, obj, act)
    if err == nil {
        permissionCache.Set(key, allowed, cache.DefaultExpiration)
    }
    
    return allowed, err
}
```

## 注意事项

1. **权限变更要同步清缓存**
2. **超级管理员要特殊处理**（绕过检查或赋予 `*,*,*` 权限）
3. **API 路径要标准化**（如 `/system/users` vs `/system/user`）
4. **权限检查要在 Gateway 统一进行**
5. **敏感操作建议二次验证**（如删除租户）

## 总结

### 当前实现（简单）
```
HTTP Method → read/write → Casbin 检查
```

### 推荐改进（细粒度）
```
HTTP Method → create/read/update/delete → Casbin 检查
           ↓
      菜单权限配置（Role.menu_permissions）
           ↓
      前端界面支持勾选 CRUD
```

这样就能实现：
- ✅ 允许查看用户列表但不能创建
- ✅ 允许编辑但不能删除
- ✅ 只读权限
- ✅ 完整权限

