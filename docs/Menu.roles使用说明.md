# Menu.roles 字段使用说明

## ⚠️ 重要说明

`Menu.roles` 字段容易与角色分配菜单（`Role.menus`）混淆，使用前请明确理解。

## 两种权限控制的区别

### 1. Role.menus（角色拥有的菜单）- **主要方式**
```yaml
含义: 这个角色可以访问哪些菜单
配置位置: 角色管理 → 分配权限菜单
使用者: 权限管理员
```

**示例：**
```json
{
  "code": "manager",
  "name": "部门经理",
  "menus": ["dashboard", "admin", "role"]
}
```
- 含义：manager 角色的用户可以访问 dashboard、admin、role 三个菜单

### 2. Menu.roles（菜单允许的角色标识）- **系统级限制**
```yaml
含义: 访问这个菜单需要具备哪些角色标识
配置位置: 系统初始化 / 超级管理员
使用者: 系统管理员
```

**示例：**
```json
{
  "name": "tenant",
  "title": "租户管理",
  "roles": ["super"]  // 只有具有 'super' 标识的用户能访问
}
```

## 问题场景

### ❌ 错误理解
```
管理员A: 在角色管理中，给 manager 角色分配了 tenant 菜单
管理员B: 在菜单管理中，设置 tenant.roles = ['super']

结果: manager 角色被分配了 tenant，但用户访问时被拒绝
原因: 用户的角色标识中没有 'super'
```

**这会造成困惑：为什么分配了菜单却不能访问？**

## 正确使用方式

### 方案A：不使用 Menu.roles（推荐）

**优点：**
- 简单直观，易于理解
- 权限管理员只需关注角色分配
- 避免混淆

**实现：**
```json
// 所有菜单的 roles 字段为空
{
  "name": "admin",
  "title": "管理员管理",
  "roles": []  // 空 = 不限制角色标识
}

// 权限完全由 Role.menus 控制
Role(manager): menus = ['dashboard', 'admin']
Role(viewer): menus = ['dashboard']
```

**界面设计：**
- 菜单管理界面：**不显示** roles 字段
- 角色管理界面：分配菜单列表（已实现）

### 方案B：Menu.roles 作为系统保护（高级）

**适用场景：**
- 需要多层安全控制
- 保护超级敏感功能
- 防止权限管理员误分配

**实现：**
```json
// 系统初始化时设置，不可通过普通界面修改
{
  "name": "tenant",
  "title": "租户管理", 
  "roles": ["super"]  // 硬限制
}

{
  "name": "system_config",
  "title": "系统配置",
  "roles": ["super"]  // 硬限制
}

// 其他菜单保持空
{
  "name": "dashboard",
  "roles": []  // 任何角色都可以，只要被分配
}
```

**权限检查流程：**
```typescript
function canAccessMenu(user, menu) {
  // 1. 用户角色是否被分配了这个菜单
  const hasMenuInRole = user.roles.some(role => 
    role.menus.includes(menu.name)
  )
  if (!hasMenuInRole) return false
  
  // 2. 如果菜单设置了 roles，检查用户是否有对应的角色标识
  if (menu.roles && menu.roles.length > 0) {
    const hasRequiredRole = menu.roles.some(roleCode =>
      user.roleIdentifiers.includes(roleCode)
    )
    if (!hasRequiredRole) return false
  }
  
  return true
}
```

**数据结构：**
```typescript
// 用户数据
{
  "user_id": "xxx",
  "roles": [
    {
      "id": "role_id_1",
      "code": "manager",    // 角色标识
      "menus": ["dashboard", "admin", "tenant"]
    }
  ],
  "roleIdentifiers": ["manager"]  // 所有角色的 code 集合
}

// 菜单数据
{
  "name": "tenant",
  "roles": ["super"]  // 需要 super 标识
}

// 检查
user.roles 包含 tenant ✅
但 user.roleIdentifiers 不包含 'super' ❌
→ 不能访问
```

### 方案C：角色标识与角色分离（最灵活）

**设计思路：**
- 角色（Role）：用于分配菜单权限
- 标签（Tag/Label）：用于细粒度访问控制

**实现：**
```typescript
// Admin 模型
{
  "roles": ["role_id_1", "role_id_2"],  // 角色ID（决定菜单）
  "tags": ["super", "auditor"]           // 标签（决定细粒度权限）
}

// Menu 模型  
{
  "name": "tenant",
  "required_tags": ["super"]  // 需要的标签
}

// 检查
hasMenu = user.roles 的 menus 包含 tenant
hasTags = user.tags 包含 "super"
canAccess = hasMenu && hasTags
```

## 推荐实施

### 阶段1：简化模式（立即可用）
1. **不使用 Menu.roles**
2. 菜单管理界面不显示 roles 字段
3. 完全基于角色分配菜单控制
4. 清晰、简单、易维护 ✅

### 阶段2：高级模式（可选）
1. **保留 Menu.roles** 作为系统级保护
2. 通过代码或数据库脚本初始化敏感菜单的 roles
3. 普通管理员不可修改
4. 仅超级管理员通过特殊界面修改

### 阶段3：标签系统（长期）
1. 引入独立的标签系统
2. 角色用于分配菜单（粗粒度）
3. 标签用于细粒度控制（Menu.roles 改为 Menu.required_tags）
4. 更灵活、更清晰

## 结论

**对于当前系统，建议：**
- ✅ **不在菜单管理界面显示 roles 字段**
- ✅ **保留字段但不通过界面配置**
- ✅ **完全通过角色分配菜单控制权限**
- ✅ **Menu.roles 留空或由系统初始化**

**如果需要细粒度控制：**
- 🔄 考虑引入标签系统
- 🔄 或者将 Menu.roles 改名为 Menu.access_tags
- 🔄 避免与角色分配菜单混淆

