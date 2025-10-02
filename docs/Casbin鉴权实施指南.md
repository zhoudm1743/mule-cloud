# Casbin 鉴权系统实施指南

## 📋 系统架构

### 一超多租户 + Casbin RBAC 架构

```
超级管理员（Super Admin）
  └─ 管理所有租户
      └─ 租户A（Tenant A）
          ├─ 管理员用户（Admin Users）
          │   └─ 拥有角色（Roles）
          ├─ 角色（Roles）
          │   └─ 拥有菜单权限（Menu Permissions）
          └─ 菜单（Menus）- 全局资源
```

## 🗂️ 数据模型

### 1. Tenant（租户）
```go
ID          string   // 租户ID
Name        string   // 租户名称
Code        string   // 租户代码
Status      int      // 状态：1-启用 0-禁用
...
```

### 2. Admin（管理员）
```go
ID          string   // 管理员ID
TenantID    string   // 租户ID（空表示超级管理员）
Phone       string   // 手机号
Password    string   // 密码（加密）
Roles       []string // 角色ID数组
IsSuper     bool     // 是否超级管理员
...
```

### 3. Role（角色）
```go
ID          string   // 角色ID
TenantID    string   // 租户ID（空表示超级管理员角色）
Name        string   // 角色名称
Code        string   // 角色代码
Menus       []string // 菜单ID数组（权限资源）
...
```

### 4. Menu（菜单/权限资源）
```go
ID            string   // 菜单ID
PID           *string  // 父菜单ID
Name          string   // 菜单名称（路由name）
Path          string   // 菜单路径
Title         string   // 菜单标题
ComponentPath *string  // 组件路径
MenuType      string   // 菜单类型：page/dir/link
...
```

## 🔐 Casbin 策略设计

### 主体标识（Subject）

#### 超级管理员
```
super:user:{user_id}           # 超级管理员用户
super:admin                    # 超级管理员角色组
```

#### 租户用户和角色
```
tenant:{tenant_id}:user:{user_id}    # 租户用户
tenant:{tenant_id}:role:{role_id}    # 租户角色
```

### 策略示例

#### 策略（Policy）- 角色对资源的权限
```
p, tenant:abc123:role:role001, /system/users, read
p, tenant:abc123:role:role001, /system/users, write
p, super:admin, *, *
```

#### 分组（Grouping）- 用户-角色关联
```
g, tenant:abc123:user:user001, tenant:abc123:role:role001
g, super:user:admin001, super:admin
```

### Casbin 模型配置
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

## 🚀 已实现的功能

### 1. 核心模块

#### ✅ Role 模型和仓库
- 文件：`internal/models/role.go`
- 文件：`internal/repository/role.go`
- 功能：角色的 CRUD 操作，支持租户隔离

#### ✅ Admin 模型增强
- 文件：`internal/models/admin.go`
- 新增字段：
  - `TenantID`：租户ID
  - `Roles`：角色ID数组
  - `IsSuper`：超级管理员标识

#### ✅ Casbin 集成
- 文件：`core/casbin/casbin.go`
- 功能：
  - 初始化 Casbin Enforcer
  - MongoDB 适配器
  - 权限检查
  - 策略管理
  - 用户-角色管理

#### ✅ 鉴权中间件
- 文件：`app/gateway/middleware/casbin_auth.go`
- 功能：
  - 请求拦截和权限验证
  - 超级管理员特权
  - 租户隔离

### 2. API 接口

#### 角色管理 API
```
GET    /admin/system/roles              - 角色列表（分页）
GET    /admin/system/roles/:id          - 获取角色详情
GET    /admin/system/roles/tenant       - 获取租户下的所有角色
POST   /admin/system/roles              - 创建角色
PUT    /admin/system/roles/:id          - 更新角色
DELETE /admin/system/roles/:id          - 删除角色
POST   /admin/system/roles/batch-delete - 批量删除

POST   /admin/system/roles/:id/menus    - 分配菜单权限
GET    /admin/system/roles/:id/menus    - 获取角色的菜单权限
```

#### 用户角色管理 API
```
POST   /admin/system/admins/:id/roles        - 分配角色给用户
GET    /admin/system/admins/:id/roles        - 获取用户的角色
DELETE /admin/system/admins/:id/roles/:roleId - 移除用户的角色
```

## 📝 使用示例

### 1. 创建租户
```bash
curl -X POST http://localhost:8080/admin/system/tenants \
  -H "Content-Type: application/json" \
  -d '{
    "name": "示例公司",
    "code": "demo_corp",
    "contact": "张三",
    "phone": "13800138000"
  }'
```

### 2. 创建角色
```bash
curl -X POST http://localhost:8080/admin/system/roles \
  -H "Content-Type: application/json" \
  -d '{
    "tenant_id": "{tenant_id}",
    "name": "系统管理员",
    "code": "sys_admin",
    "description": "拥有系统管理权限",
    "menus": []
  }'
```

### 3. 分配菜单权限给角色
```bash
curl -X POST http://localhost:8080/admin/system/roles/{role_id}/menus \
  -H "Content-Type: application/json" \
  -d '{
    "menus": [
      "{menu_id_1}",
      "{menu_id_2}",
      "{menu_id_3}"
    ]
  }'
```

### 4. 分配角色给用户
```bash
curl -X POST http://localhost:8080/admin/system/admins/{admin_id}/roles \
  -H "Content-Type: application/json" \
  -d '{
    "roles": [
      "{role_id_1}",
      "{role_id_2}"
    ]
  }'
```

## 🔧 集成步骤

### 1. 初始化 Casbin（在网关启动时）

```go
import (
    casbinPkg "mule-cloud/core/casbin"
)

// 在 main.go 中初始化
func main() {
    // ... 其他初始化代码 ...
    
    // 初始化 Casbin
    casbinConfig := &casbinPkg.Config{
        MongoURI:     cfg.MongoDB.URI,
        DatabaseName: cfg.MongoDB.Database,
        ModelPath:    "core/casbin/model.conf", // 可选
    }
    
    enforcer, err := casbinPkg.InitCasbin(casbinConfig)
    if err != nil {
        log.Fatalf("初始化Casbin失败: %v", err)
    }
    
    log.Println("✅ Casbin 初始化成功")
}
```

### 2. 在网关使用鉴权中间件

```go
import (
    "mule-cloud/app/gateway/middleware"
)

// 在动态路由处理中添加鉴权中间件
func (m *DynamicRouteManager) buildHandlers(routeConfig *RouteConfig) []gin.HandlerFunc {
    handlers := []gin.HandlerFunc{}
    
    // CORS
    handlers = append(handlers, middleware.CORSMiddleware())
    
    // 认证（如果需要）
    if routeConfig.RequireAuth {
        handlers = append(handlers, middleware.JWTAuthMiddleware())
    }
    
    // Casbin 鉴权
    handlers = append(handlers, middleware.CasbinAuthMiddleware())
    
    // 反向代理
    handlers = append(handlers, m.proxyHandler())
    
    return handlers
}
```

### 3. 同步角色菜单权限到 Casbin

```go
import casbinPkg "mule-cloud/core/casbin"

// 在 RoleService.AssignMenus 中同步
func (s *RoleService) AssignMenus(ctx context.Context, roleID string, menuIDs []string, updatedBy string) error {
    // 获取角色
    role, err := s.roleRepo.Get(ctx, roleID)
    if err != nil {
        return err
    }
    
    // 更新数据库
    updates := map[string]interface{}{
        "menus":      menuIDs,
        "updated_by": updatedBy,
        "updated_at": time.Now().Unix(),
    }
    err = s.roleRepo.Update(ctx, roleID, updates)
    if err != nil {
        return err
    }
    
    // 获取菜单路径列表
    menuPaths := []string{}
    for _, menuID := range menuIDs {
        menu, _ := s.menuRepo.GetByID(ctx, menuID)
        if menu != nil {
            menuPaths = append(menuPaths, menu.Path)
        }
    }
    
    // 同步到 Casbin
    return casbinPkg.SyncRoleMenus(role.TenantID, roleID, menuPaths)
}
```

### 4. 同步用户角色到 Casbin

```go
import casbinPkg "mule-cloud/core/casbin"

// 在 AdminService.AssignRoles 中同步
func (s *AdminService) AssignRoles(ctx context.Context, adminID string, roleIDs []string, updatedBy string) error {
    // 获取管理员
    admin, err := s.repo.Get(ctx, adminID)
    if err != nil {
        return err
    }
    
    // 更新数据库
    updates := bson.M{
        "roles":      roleIDs,
        "updated_by": updatedBy,
        "updated_at": time.Now().Unix(),
    }
    err = s.repo.Update(ctx, adminID, updates)
    if err != nil {
        return err
    }
    
    // 同步到 Casbin
    return casbinPkg.SyncUserRoles(admin.TenantID, adminID, roleIDs)
}
```

## ⚠️ 待完成事项

### 1. JWT 集成
- [ ] 在 JWT token 中包含：`user_id`、`tenant_id`、`is_super`
- [ ] JWT 中间件解析 token 并设置到 Gin Context
- [ ] 更新 CasbinAuthMiddleware 从 context 获取用户信息

### 2. 网关初始化 Casbin
- [ ] 在 `cmd/gateway/main.go` 中初始化 Casbin
- [ ] 配置 MongoDB 连接信息

### 3. 权限同步
- [ ] 在角色分配菜单时同步到 Casbin
- [ ] 在用户分配角色时同步到 Casbin
- [ ] 在删除角色/用户时清理 Casbin 策略

### 4. 超级管理员初始化
- [ ] 创建初始化脚本添加超级管理员
- [ ] 设置超级管理员的特殊权限

### 5. 前端集成
- [ ] 角色管理界面
- [ ] 权限分配界面
- [ ] 用户角色绑定界面

## 📚 参考资料

- [Casbin 官方文档](https://casbin.org/zh/)
- [Casbin MongoDB 适配器](https://github.com/casbin/mongodb-adapter)
- [Go-Casbin API 文档](https://pkg.go.dev/github.com/casbin/casbin/v2)

## 🎯 核心优势

1. **租户隔离**：每个租户的数据完全隔离
2. **灵活权限**：基于 RBAC 的细粒度权限控制
3. **超级管理员**：支持跨租户的超级管理员
4. **可扩展**：易于添加新的权限规则和角色
5. **持久化**：权限策略持久化到 MongoDB
6. **高性能**：Casbin 内置高效的策略匹配引擎

## 🔒 安全建议

1. 所有密码必须加密存储
2. JWT token 设置合理的过期时间
3. 敏感操作记录审计日志
4. 定期review权限配置
5. 实施最小权限原则

