# Casbin 自定义业务权限方案

## 问题场景

除了基础的 CRUD（增删改查）权限外，实际业务中还有很多特殊权限：

| 业务场景 | 自定义权限 |
|---------|-----------|
| 财务管理 | `挂账`, `核销`, `冲账`, `审核` |
| 订单管理 | `取消订单`, `修改价格`, `强制关闭` |
| 用户管理 | `重置密码`, `锁定账号`, `解锁账号` |
| 审批流程 | `提交审批`, `审批通过`, `驳回` |
| 数据操作 | `导出`, `打印`, `批量导入` |
| 状态管理 | `启用`, `禁用`, `归档` |

## 解决方案

### 方案A：扩展权限动作（推荐）

#### 核心思路
- ✅ **基础权限**：`create`, `read`, `update`, `delete`（对应 HTTP 方法）
- ✅ **业务权限**：任意自定义字符串（不对应 HTTP 方法，需要在代码中手动检查）

#### 1. 定义资源的可用权限

```go
// internal/models/menu.go

type Menu struct {
    ID              primitive.ObjectID `bson:"_id,omitempty"`
    Name            string            `bson:"name"`
    Path            string            `bson:"path"`
    Title           string            `bson:"title"`
    
    // 新增：该菜单/资源支持的权限列表
    AvailablePermissions []Permission `bson:"available_permissions" json:"available_permissions"`
}

type Permission struct {
    Action      string `json:"action"`       // read, create, update, delete, 挂账, 核销...
    Label       string `json:"label"`        // 显示名称
    Description string `json:"description"`  // 描述
    IsBasic     bool   `json:"is_basic"`     // 是否基础权限（CRUD）
}
```

#### 2. 权限配置示例

```json
// 普通菜单（只有基础权限）
{
  "name": "admin",
  "path": "/system/admin",
  "title": "管理员管理",
  "available_permissions": [
    { "action": "read",   "label": "查看", "is_basic": true },
    { "action": "create", "label": "创建", "is_basic": true },
    { "action": "update", "label": "修改", "is_basic": true },
    { "action": "delete", "label": "删除", "is_basic": true }
  ]
}

// 订单管理（有业务权限）
{
  "name": "order",
  "path": "/business/order",
  "title": "订单管理",
  "available_permissions": [
    { "action": "read",          "label": "查看",     "is_basic": true },
    { "action": "create",        "label": "创建",     "is_basic": true },
    { "action": "update",        "label": "修改",     "is_basic": true },
    { "action": "delete",        "label": "删除",     "is_basic": true },
    { "action": "cancel",        "label": "取消订单",  "is_basic": false },
    { "action": "price_adjust",  "label": "调整价格",  "is_basic": false },
    { "action": "force_close",   "label": "强制关闭",  "is_basic": false }
  ]
}

// 财务管理（有挂账等权限）
{
  "name": "finance",
  "path": "/business/finance",
  "title": "财务管理",
  "available_permissions": [
    { "action": "read",          "label": "查看",     "is_basic": true },
    { "action": "create",        "label": "创建",     "is_basic": true },
    { "action": "pending",       "label": "挂账",     "is_basic": false, "description": "将订单挂账延期支付" },
    { "action": "verify",        "label": "核销",     "is_basic": false, "description": "核销已挂账订单" },
    { "action": "reverse",       "label": "冲账",     "is_basic": false, "description": "冲销错误账目" },
    { "action": "audit",         "label": "审核",     "is_basic": false, "description": "财务审核" },
    { "action": "export",        "label": "导出",     "is_basic": false }
  ]
}
```

#### 3. 前端动态权限配置界面

```vue
<!-- frontend/src/views/system/role/components/TableModal.vue -->
<template>
  <NModal v-model:show="modalVisible" title="分配权限">
    <div v-for="menu in menus" :key="menu.name">
      <NCheckbox v-model:checked="selectedMenus[menu.name]">
        {{ menu.title }}
      </NCheckbox>
      
      <!-- 动态显示该菜单的可用权限 -->
      <div v-if="selectedMenus[menu.name] && menu.available_permissions" class="ml-8 mt-2">
        <NSpace>
          <!-- 基础权限 -->
          <NTag type="info">基础权限</NTag>
          <NCheckbox
            v-for="perm in menu.available_permissions.filter(p => p.is_basic)"
            :key="perm.action"
            v-model:checked="permissions[menu.name][perm.action]"
          >
            {{ perm.label }}
          </NCheckbox>
        </NSpace>
        
        <!-- 业务权限 -->
        <NSpace v-if="menu.available_permissions.some(p => !p.is_basic)" class="mt-2">
          <NTag type="warning">业务权限</NTag>
          <NCheckbox
            v-for="perm in menu.available_permissions.filter(p => !p.is_basic)"
            :key="perm.action"
            v-model:checked="permissions[menu.name][perm.action]"
            :title="perm.description"
          >
            {{ perm.label }}
            <NTooltip v-if="perm.description">
              <template #trigger>
                <NIcon :component="QuestionCircleOutlined" class="ml-1" />
              </template>
              {{ perm.description }}
            </NTooltip>
          </NCheckbox>
        </NSpace>
      </div>
    </div>
  </NModal>
</template>

<script setup lang="ts">
// 权限数据结构
const permissions = ref<Record<string, Record<string, boolean>>>({
  'finance': {
    'read': true,
    'create': true,
    'update': true,
    'delete': false,
    'pending': true,    // 挂账
    'verify': true,     // 核销
    'reverse': false,   // 冲账
    'audit': true,      // 审核
    'export': true      // 导出
  }
})

// 转换为后端格式
function getPermissionsForSubmit() {
  const result: Record<string, string[]> = {}
  
  for (const [menuName, perms] of Object.entries(permissions.value)) {
    result[menuName] = Object.entries(perms)
      .filter(([_, enabled]) => enabled)
      .map(([action, _]) => action)
  }
  
  return result
  // 结果：{ "finance": ["read", "create", "update", "pending", "verify", "audit", "export"] }
}
</script>
```

界面效果：
```
☑️ 财务管理
   ├─ [基础权限]
   │   ☑️ 查看  ☑️ 创建  ☑️ 修改  ☐ 删除
   └─ [业务权限]
       ☑️ 挂账 ℹ️  ☑️ 核销 ℹ️  ☐ 冲账 ℹ️  ☑️ 审核 ℹ️  ☑️ 导出
```

#### 4. 后端权限检查

##### A. Gateway 自动检查（基础权限）

```go
// app/gateway/middleware/casbin_auth.go

func CasbinAuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 只检查基础 CRUD 权限
        method := c.Request.Method
        path := c.Request.URL.Path
        
        action := getActionFromMethod(method)
        
        // 检查基础权限
        allowed, _ := casbin.CheckPermission(userSub, path, action)
        if !allowed {
            response.Error(c, "权限不足")
            c.Abort()
            return
        }
        
        c.Next()
    }
}
```

##### B. 业务代码手动检查（业务权限）

```go
// app/business/services/order.go

// CancelOrder 取消订单（需要 cancel 权限）
func (s *OrderService) CancelOrder(c *gin.Context, orderID string) error {
    userID := c.GetString("user_id")
    tenantID := c.GetString("tenant_id")
    
    // 检查业务权限
    allowed, err := casbin.CheckUserPermission(
        tenantID,
        userID,
        "/business/order",  // 资源路径
        "cancel",           // 业务动作
    )
    
    if err != nil || !allowed {
        return fmt.Errorf("无权限执行此操作")
    }
    
    // 执行取消逻辑
    return s.repo.CancelOrder(orderID)
}

// AdjustPrice 调整价格（需要 price_adjust 权限）
func (s *OrderService) AdjustPrice(c *gin.Context, orderID string, newPrice float64) error {
    userID := c.GetString("user_id")
    tenantID := c.GetString("tenant_id")
    
    // 检查业务权限
    allowed, err := casbin.CheckUserPermission(
        tenantID,
        userID,
        "/business/order",
        "price_adjust",
    )
    
    if !allowed {
        return fmt.Errorf("无权限调整价格")
    }
    
    return s.repo.UpdatePrice(orderID, newPrice)
}
```

##### C. 财务挂账示例

```go
// app/business/services/finance.go

// PendingPayment 挂账
func (s *FinanceService) PendingPayment(c *gin.Context, req *dto.PendingRequest) error {
    userID := c.GetString("user_id")
    tenantID := c.GetString("tenant_id")
    
    // 检查挂账权限
    allowed, err := casbin.CheckUserPermission(
        tenantID,
        userID,
        "/business/finance",
        "pending",  // 挂账权限
    )
    
    if !allowed {
        return fmt.Errorf("无挂账权限")
    }
    
    // 执行挂账逻辑
    return s.repo.CreatePendingPayment(req)
}

// VerifyPayment 核销
func (s *FinanceService) VerifyPayment(c *gin.Context, recordID string) error {
    userID := c.GetString("user_id")
    tenantID := c.GetString("tenant_id")
    
    // 检查核销权限
    allowed, err := casbin.CheckUserPermission(
        tenantID,
        userID,
        "/business/finance",
        "verify",  // 核销权限
    )
    
    if !allowed {
        return fmt.Errorf("无核销权限")
    }
    
    return s.repo.VerifyPendingPayment(recordID)
}
```

#### 5. Casbin 策略示例

```bash
# 财务主管角色 - 完整权限
p, tenant:A:role:finance_manager, /business/finance, read
p, tenant:A:role:finance_manager, /business/finance, create
p, tenant:A:role:finance_manager, /business/finance, update
p, tenant:A:role:finance_manager, /business/finance, delete
p, tenant:A:role:finance_manager, /business/finance, pending    # 挂账
p, tenant:A:role:finance_manager, /business/finance, verify     # 核销
p, tenant:A:role:finance_manager, /business/finance, reverse    # 冲账
p, tenant:A:role:finance_manager, /business/finance, audit      # 审核
p, tenant:A:role:finance_manager, /business/finance, export     # 导出

# 财务专员 - 部分权限（不能冲账、审核）
p, tenant:A:role:finance_staff, /business/finance, read
p, tenant:A:role:finance_staff, /business/finance, create
p, tenant:A:role:finance_staff, /business/finance, pending      # 可以挂账
p, tenant:A:role:finance_staff, /business/finance, verify       # 可以核销
p, tenant:A:role:finance_staff, /business/finance, export       # 可以导出
# 没有 reverse（冲账）和 audit（审核）权限

# 财务查看 - 只读 + 导出
p, tenant:A:role:finance_viewer, /business/finance, read
p, tenant:A:role:finance_viewer, /business/finance, export
```

---

## 方案B：按钮级权限控制

### 思路
为每个功能按钮定义独立的权限标识

#### 1. 定义按钮权限

```typescript
// frontend/src/views/business/finance/index.vue

const buttonPermissions = {
  'btn_pending': 'pending',      // 挂账按钮
  'btn_verify': 'verify',        // 核销按钮
  'btn_reverse': 'reverse',      // 冲账按钮
  'btn_audit': 'audit',          // 审核按钮
  'btn_export': 'export',        // 导出按钮
}

// 检查按钮权限
function hasButtonPermission(btnKey: string) {
  const action = buttonPermissions[btnKey]
  return hasPermission(`/business/finance:${action}`)
}
```

#### 2. 前端使用

```vue
<template>
  <NSpace>
    <!-- 基础按钮 -->
    <NButton v-if="hasPermission('/business/finance:create')" @click="handleCreate">
      新建
    </NButton>
    
    <!-- 业务按钮 -->
    <NButton 
      v-if="hasButtonPermission('btn_pending')" 
      type="warning"
      @click="handlePending"
    >
      挂账
    </NButton>
    
    <NButton 
      v-if="hasButtonPermission('btn_verify')" 
      type="success"
      @click="handleVerify"
    >
      核销
    </NButton>
    
    <NButton 
      v-if="hasButtonPermission('btn_reverse')" 
      type="error"
      @click="handleReverse"
    >
      冲账
    </NButton>
    
    <NButton 
      v-if="hasButtonPermission('btn_audit')" 
      type="info"
      @click="handleAudit"
    >
      审核
    </NButton>
    
    <NButton 
      v-if="hasButtonPermission('btn_export')" 
      @click="handleExport"
    >
      导出
    </NButton>
  </NSpace>
</template>
```

---

## 方案C：操作权限配置表

### 适用场景
权限非常复杂，需要动态配置

#### 数据结构

```typescript
// 操作权限配置
interface ActionConfig {
  code: string           // pending, verify, reverse
  name: string           // 挂账, 核销, 冲账
  resource: string       // /business/finance
  method?: string        // POST, PUT, DELETE（可选）
  endpoint?: string      // /api/finance/pending（可选）
  description: string
  category: 'basic' | 'business' | 'advanced'
}

// 示例数据
const financeActions: ActionConfig[] = [
  {
    code: 'pending',
    name: '挂账',
    resource: '/business/finance',
    method: 'POST',
    endpoint: '/api/finance/pending',
    description: '将订单挂账延期支付',
    category: 'business'
  },
  {
    code: 'verify',
    name: '核销',
    resource: '/business/finance',
    method: 'POST',
    endpoint: '/api/finance/verify',
    description: '核销已挂账订单',
    category: 'business'
  }
]
```

---

## 完整流程示例

### 1. 系统初始化时配置菜单权限

```javascript
// scripts/init_menu_permissions.js

db.menus.updateOne(
  { name: "finance" },
  {
    $set: {
      available_permissions: [
        { action: "read",    label: "查看", is_basic: true },
        { action: "create",  label: "创建", is_basic: true },
        { action: "update",  label: "修改", is_basic: true },
        { action: "delete",  label: "删除", is_basic: true },
        { action: "pending", label: "挂账", is_basic: false, description: "将订单挂账延期支付" },
        { action: "verify",  label: "核销", is_basic: false, description: "核销已挂账订单" },
        { action: "reverse", label: "冲账", is_basic: false, description: "冲销错误账目" },
        { action: "audit",   label: "审核", is_basic: false, description: "财务审核" },
        { action: "export",  label: "导出", is_basic: false }
      ]
    }
  }
)
```

### 2. 管理员分配角色权限

```
管理员登录 → 角色管理 → 编辑"财务专员"角色 → 分配权限：

☑️ 财务管理
   [基础权限]
   ☑️ 查看  ☑️ 创建  ☑️ 修改  ☐ 删除
   
   [业务权限]
   ☑️ 挂账  ☑️ 核销  ☐ 冲账  ☐ 审核  ☑️ 导出
```

### 3. 后端保存权限到 Casbin

```go
// 保存到 MongoDB casbin_rule
{
  ptype: "p",
  v0: "tenant:A:role:finance_staff",
  v1: "/business/finance",
  v2: "pending"
}
```

### 4. 用户操作时检查权限

```go
// 用户点击"挂账"按钮 → 发送请求
POST /api/finance/pending

// 后端检查
allowed := casbin.CheckUserPermission(
    tenantID,
    userID,
    "/business/finance",
    "pending"
)

// 如果 allowed = true → 执行挂账
// 如果 allowed = false → 返回 403
```

---

## 最佳实践

### 1. 权限命名规范

```
基础权限: read, create, update, delete
状态操作: enable, disable, archive, restore
审批流程: submit, approve, reject, cancel
财务操作: pending, verify, reverse, audit, refund
数据操作: export, import, print, download
账号操作: reset_password, lock, unlock
```

### 2. 权限分层

```
第一层：HTTP 方法权限（Gateway 自动检查）
  - GET → read
  - POST → create
  - PUT/PATCH → update
  - DELETE → delete

第二层：业务操作权限（代码手动检查）
  - pending（挂账）
  - verify（核销）
  - audit（审核）
  ...
```

### 3. 前端权限控制

```typescript
// 封装权限检查 hook
import { usePermission } from '@/hooks'

const { hasPermission } = usePermission()

// 基础权限
const canCreate = hasPermission('/business/finance:create')
const canUpdate = hasPermission('/business/finance:update')

// 业务权限
const canPending = hasPermission('/business/finance:pending')
const canVerify = hasPermission('/business/finance:verify')

// 按钮显示控制
<NButton v-if="canPending" @click="handlePending">挂账</NButton>
```

---

## 总结

### 实现步骤

1. **定义权限** → Menu 模型添加 `available_permissions` 字段
2. **配置界面** → 前端角色管理支持勾选自定义权限
3. **保存策略** → 后端同步到 Casbin
4. **权限检查** → 业务代码调用 `CheckUserPermission(resource, action)`

### 优势

- ✅ 灵活扩展：可以添加任意自定义权限
- ✅ 统一管理：所有权限通过 Casbin 统一控制
- ✅ 前后端协同：前端隐藏按钮，后端拒绝请求
- ✅ 可视化配置：管理员可以通过界面配置权限

### 示例总结

```
基础权限（HTTP）→ read, create, update, delete
业务权限（自定义）→ pending(挂账), verify(核销), reverse(冲账), audit(审核)
前端显示控制 → v-if="hasPermission(...)"
后端强制检查 → CheckUserPermission(resource, action)
```

这样就能完美支持任意自定义业务权限了！🎉

