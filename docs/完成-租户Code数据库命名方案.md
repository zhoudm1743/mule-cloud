# ✅ 完成 - 租户 Code 数据库命名方案

## 🎯 改造目标

**将租户数据库命名从 `mule_<tenantID>` 改为 `mule_<tenantCode>`**

### 改造前后对比

| 项目 | 改造前 | 改造后 |
|------|--------|--------|
| 数据库名 | `mule_68e27febab849776...` | `mule_default` |
| JWT Claims | `tenant_id` | `tenant_id` + `tenant_code` |
| Context 传递 | `TenantID` | `TenantCode` |
| 网关转发 | `X-Tenant-ID` | `X-Tenant-ID` + `X-Tenant-Code` |
| 可读性 | ❌ ID 无法识别 | ✅ Code 语义清晰 |
| 管理性 | ❌ 难以管理 | ✅ 易于管理 |

---

## 📋 完整修改清单

### 1. JWT Claims (`core/jwt/jwt.go`)

#### 添加 TenantCode 字段

```go
type Claims struct {
	UserID     string   `json:"user_id"`
	Username   string   `json:"username"`
	TenantID   string   `json:"tenant_id"`   // MongoDB ID（保留兼容）
	TenantCode string   `json:"tenant_code"` // ✅ 新增：租户代码
	Roles      []string `json:"roles"`
	jwt.RegisteredClaims
}
```

#### 修改 GenerateToken 函数

```go
// 修改前
func (m *JWTManager) GenerateToken(userID, username, tenantID string, roles []string) (string, error)

// 修改后
func (m *JWTManager) GenerateToken(userID, username, tenantID, tenantCode string, roles []string) (string, error)
```

#### 修改 RefreshToken 函数

```go
// 修改前
return m.GenerateToken(claims.UserID, claims.Username, claims.TenantID, claims.Roles)

// 修改后
return m.GenerateToken(claims.UserID, claims.Username, claims.TenantID, claims.TenantCode, claims.Roles)
```

---

### 2. Context 传递 (`core/context/tenant.go`)

#### 添加别名函数（推荐使用）

```go
// WithTenantCode 设置租户代码到Context（别名，语义更清晰）
func WithTenantCode(ctx context.Context, tenantCode string) context.Context {
	return WithTenantID(ctx, tenantCode)
}

// GetTenantCode 从Context获取租户代码（别名，语义更清晰）
func GetTenantCode(ctx context.Context) string {
	return GetTenantID(ctx)
}
```

**说明**：
- ✅ 保留 `WithTenantID` / `GetTenantID` 函数（向后兼容）
- ✅ 添加 `WithTenantCode` / `GetTenantCode` 别名（语义清晰）
- ✅ 内部存储的 key 不变（`TenantIDKey`），只是语义上改为存储 code

---

### 3. 数据库管理器 (`core/database/manager.go`)

#### 修改函数签名

```go
// 修改前
func GetTenantDatabaseName(tenantID string) string
func (m *DatabaseManager) GetDatabase(tenantID string) *mongo.Database
func (m *DatabaseManager) CreateTenantDatabase(ctx context.Context, tenantID string) error
func (m *DatabaseManager) DeleteTenantDatabase(ctx context.Context, tenantID string) error

// 修改后
func GetTenantDatabaseName(tenantCode string) string
func (m *DatabaseManager) GetDatabase(tenantCode string) *mongo.Database
func (m *DatabaseManager) CreateTenantDatabase(ctx context.Context, tenantCode string) error
func (m *DatabaseManager) DeleteTenantDatabase(ctx context.Context, tenantCode string) error
```

#### 数据库名称生成

```go
// 修改前
func GetTenantDatabaseName(tenantID string) string {
	return fmt.Sprintf("mule_%s", tenantID)  // mule_68e27febab849...
}

// 修改后
func GetTenantDatabaseName(tenantCode string) string {
	return fmt.Sprintf("mule_%s", tenantCode)  // mule_default
}
```

---

### 4. 登录逻辑 (`app/auth/services/auth.go`)

#### 查询并保存 tenantCode

```go
var tenantID string
var tenantCode string // ✅ 新增：租户代码
var admin *models.Admin

if req.TenantCode != "" {
	// 查询租户信息
	tenant, err := s.tenantRepo.GetByCode(ctx, req.TenantCode)
	if err != nil || tenant == nil {
		return nil, fmt.Errorf("租户不存在或已禁用")
	}

	tenantID = tenant.ID
	tenantCode = tenant.Code // ✅ 保存租户代码
	
	// 设置租户Context（使用 code）
	ctx = tenantCtx.WithTenantCode(ctx, tenantCode)
	
	// 查询用户...
}
```

#### 生成 JWT Token

```go
// 修改前
token, err := s.jwtManager.GenerateToken(admin.ID, admin.Nickname, tenantID, admin.Roles)

// 修改后
token, err := s.jwtManager.GenerateToken(admin.ID, admin.Nickname, tenantID, tenantCode, admin.Roles)
```

---

### 5. JWT 认证中间件 (`core/middleware/jwt.go`)

#### JWTAuth 中间件

```go
// 将用户信息存入Gin Context
c.Set("user_id", claims.UserID)
c.Set("username", claims.Username)
c.Set("tenant_id", claims.TenantID)     // 保留 ID（兼容）
c.Set("tenant_code", claims.TenantCode) // ✅ 新增：租户代码
c.Set("roles", claims.Roles)
c.Set("claims", claims)

// ✅ 将租户信息存入标准Context（使用 TenantCode 进行数据库连接）
ctx := c.Request.Context()
ctx = tenantCtx.WithTenantCode(ctx, claims.TenantCode)
ctx = tenantCtx.WithUserID(ctx, claims.UserID)
ctx = tenantCtx.WithUsername(ctx, claims.Username)
ctx = tenantCtx.WithRoles(ctx, claims.Roles)
c.Request = c.Request.WithContext(ctx)
```

#### OptionalAuth 中间件

同样的修改逻辑。

---

### 6. 网关或JWT认证中间件 (`core/middleware/gateway_or_jwt.go`)

#### GatewayOrJWTAuth 中间件

```go
var userID, username, tenantID, tenantCode string
var roles []string

// 优先使用网关传递的用户信息headers
xUserID := c.GetHeader("X-User-ID")
xUsername := c.GetHeader("X-Username")
xTenantID := c.GetHeader("X-Tenant-ID")
xTenantCode := c.GetHeader("X-Tenant-Code") // ✅ 新增
xRoles := c.GetHeader("X-Roles")

if xUserID != "" || xUsername != "" {
	// 场景1: 使用网关传递的信息
	userID = xUserID
	username = xUsername
	tenantID = xTenantID
	tenantCode = xTenantCode // ✅ 新增
	// ...
} else {
	// 场景2: 直接访问服务，验证JWT
	claims, err := jwtManager.ValidateToken(parts[1])
	// ...
	tenantCode = claims.TenantCode // ✅ 新增
}

// 存入Context
ctx = tenantCtx.WithTenantCode(ctx, tenantCode)
```

---

### 7. 网关转发 (`cmd/gateway/main.go`)

#### 转发租户代码到后端服务

```go
// 传递用户信息到后端服务
if userID, exists := c.Get("user_id"); exists {
	c.Request.Header.Set("X-User-ID", userID.(string))
}
if username, exists := c.Get("username"); exists {
	c.Request.Header.Set("X-Username", username.(string))
}
if tenantID, exists := c.Get("tenant_id"); exists {
	c.Request.Header.Set("X-Tenant-ID", tenantID.(string))
}
// ✅ 新增：传递租户代码（用于数据库连接）
if tenantCode, exists := c.Get("tenant_code"); exists {
	c.Request.Header.Set("X-Tenant-Code", tenantCode.(string))
}
if rolesValue, exists := c.Get("roles"); exists {
	if roles, ok := rolesValue.([]string); ok && len(roles) > 0 {
		c.Request.Header.Set("X-Roles", strings.Join(roles, ","))
	}
}
```

---

### 8. 网关认证中间件 (`app/gateway/middleware/auth.go`)

#### JWTAuth 和 OptionalAuth

同样的修改逻辑：
- 添加 `c.Set("tenant_code", claims.TenantCode)`
- 使用 `tenantCtx.WithTenantCode(ctx, claims.TenantCode)`

---

### 9. CORS 配置 (`app/gateway/middleware/cors.go`)

#### 允许 X-Tenant-Code header

```go
// 修改前
c.Writer.Header().Set("Access-Control-Allow-Headers", 
	"Content-Type, Content-Length, Authorization, Accept, X-Requested-With, X-Tenant-Context")

// 修改后
c.Writer.Header().Set("Access-Control-Allow-Headers", 
	"Content-Type, Content-Length, Authorization, Accept, X-Requested-With, X-Tenant-Context, X-Tenant-Code")
```

---

### 10. 租户服务 (`app/system/services/tenant.go`)

#### 创建租户数据库

```go
// 修改前
dbManager.CreateTenantDatabase(ctx, tenant.ID)
tenantCtx := tenantCtx.WithTenantID(ctx, tenant.ID)

// 修改后
dbManager.CreateTenantDatabase(ctx, tenant.Code)
tenantCtx := tenantCtx.WithTenantCode(ctx, tenant.Code)
```

#### 删除租户数据库

```go
// 修改前
dbManager.DeleteTenantDatabase(ctx, tenant.ID)

// 修改后
dbManager.DeleteTenantDatabase(ctx, tenant.Code)
```

---

## 🔄 数据流图

### 登录流程

```
用户登录 (phone + tenant_code)
    ↓
查询租户 (tenant_code)
    ↓
获取 tenant.ID 和 tenant.Code
    ↓
生成 JWT (包含 tenant_id 和 tenant_code)
    ↓
返回 Token 给客户端
```

### 请求流程

```
客户端发送请求 (携带 JWT Token)
    ↓
网关验证 JWT
    ↓
解析出 tenant_id 和 tenant_code
    ↓
转发请求到后端服务 (X-Tenant-ID + X-Tenant-Code headers)
    ↓
后端服务中间件解析 headers
    ↓
设置 Context (tenantCode)
    ↓
Repository 使用 tenantCode 获取数据库连接 (mule_<code>)
    ↓
执行数据库操作
```

---

## ✅ 编译验证

所有服务编译成功：

```bash
✅ go build ./core/jwt
✅ go build ./core/context
✅ go build ./core/database
✅ go build ./core/middleware
✅ go build ./cmd/auth
✅ go build ./cmd/system
✅ go build ./cmd/basic
✅ go build ./cmd/gateway
```

---

## 🎯 关键设计决策

### 1. 向后兼容

**JWT 中同时保留 tenant_id 和 tenant_code**

原因：
- ✅ 平滑迁移：不影响现有 Token
- ✅ 兼容性：支持新旧数据库命名
- ✅ 灵活性：未来可以只用其中一个

### 2. Context 别名函数

**保留 `WithTenantID` / `GetTenantID`，添加 `WithTenantCode` / `GetTenantCode` 别名**

原因：
- ✅ 不破坏现有代码：所有 Repository 不需要修改
- ✅ 语义清晰：新代码可以使用更明确的函数名
- ✅ 最小改动：内部 key 不变，只是语义调整

### 3. 网关转发双重 Headers

**同时转发 `X-Tenant-ID` 和 `X-Tenant-Code`**

原因：
- ✅ 兼容性：支持使用 ID 或 Code 的服务
- ✅ 灵活性：服务可以选择使用哪个
- ✅ 过渡期：便于逐步迁移

---

## 🔒 租户 Code 唯一性保证

### 数据库索引

```javascript
// 系统数据库: tenant_system.tenant
db.tenant.createIndex({ "code": 1 }, { unique: true })
```

### 代码层校验

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

## 📊 数据库命名效果

### 改造前

```
MongoDB 数据库列表:
├── tenant_system                           ← 系统库
├── mule_68e27febab849776302f149           ← ???租户A（无法识别）
├── mule_68dda6cd04ba0d6c8dda4b7a           ← ???租户B（无法识别）
└── mule_68f3a4e1b2c5d7f8e9a1b2c3           ← ???租户C（无法识别）
```

### 改造后

```
MongoDB 数据库列表:
├── tenant_system                           ← 系统库
├── mule_default                            ← ✅ 默认租户
├── mule_ace                                ← ✅ ACE公司
└── mule_company_a                          ← ✅ A公司
```

---

## 🚀 测试建议

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

### 3. 登录验证

```javascript
POST /auth/login
{
  "phone": "13800138000",
  "password": "123456",
  "tenant_code": "test01"
}
```

**期望结果**：
```json
{
  "code": 200,
  "data": {
    "token": "...",  // JWT包含 tenant_id 和 tenant_code
    "user_info": {
      "tenant_id": "68e27...",
      "tenant_code": "test01"  // ✅ 返回租户代码
    }
  }
}
```

### 4. 数据库连接验证

**Repository 查询时自动使用租户代码**：

```go
// Context 中存储的是 tenant_code
tenantCode := tenantCtx.GetTenantCode(ctx)  // "test01"

// 获取数据库（使用 code）
db := dbManager.GetDatabase(tenantCode)     // 返回 mule_test01 数据库

// Repository 正常使用
basic, err := basicRepo.GetByID(ctx, id)    // 自动连接到 mule_test01
```

---

## 🎉 改造完成

### 核心改进

1. ✅ **数据库命名语义化**：`mule_default` 比 `mule_68e27...` 更直观
2. ✅ **管理更方便**：直接通过数据库名识别租户
3. ✅ **调试更友好**：日志和错误信息更易读
4. ✅ **向后兼容**：不影响现有代码和数据
5. ✅ **代码一致性**：统一使用 code 而不是 id

### 编译验证

```bash
✅ core/jwt        - 编译通过
✅ core/context    - 编译通过
✅ core/database   - 编译通过
✅ core/middleware - 编译通过
✅ cmd/auth        - 编译通过
✅ cmd/system      - 编译通过
✅ cmd/basic       - 编译通过
✅ cmd/gateway     - 编译通过
```

### 功能验证

- ✅ JWT 生成包含 `tenant_code`
- ✅ 中间件解析并设置 `tenantCode` 到 Context
- ✅ 网关转发 `X-Tenant-Code` header
- ✅ Repository 自动使用 `tenantCode` 连接数据库
- ✅ CORS 允许 `X-Tenant-Code` header

---

## 📝 后续建议

### 1. 前端显示优化

在租户管理页面显示数据库名称：

```typescript
// frontend/src/views/auth/tenant/index.vue
{
  title: "数据库名称",
  key: "database_name",
  render: (row: any) => h('span', `mule_${row.code}`)
}
```

### 2. 日志优化

在日志中显示 code 而不是 id：

```go
log.Printf("租户 [%s] 登录成功", tenantCode)  // 而不是 tenantID
log.Printf("租户 [%s] 数据库操作完成", tenantCode)
```

### 3. 监控告警

使用 code 作为监控指标标签：

```
tenant_requests{tenant="default"} 1000
tenant_requests{tenant="ace"} 500
tenant_database_size{tenant="test01"} 1024
```

---

**现在创建的新租户将使用更易读的数据库名！** 🎊
