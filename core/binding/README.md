# Binding 包使用指南

统一的参数绑定和验证解决方案，集成了 validator v10，提供友好的中文错误提示。

## 核心功能

✅ 统一绑定 URI、Query、Body 参数  
✅ 自动格式化验证错误为中文提示  
✅ 内置常用自定义验证规则（手机号、身份证等）  
✅ 支持注册自定义验证规则  
✅ 避免参数覆盖问题  

---

## 快速开始

### 1. 基础绑定

```go
import "mule-cloud/core/binding"

func UpdateOrderHandler(c *gin.Context) {
    var req dto.OrderUpdateRequest
    
    // 一行代码搞定 URI + Body 参数绑定
    if err := binding.BindAll(c, &req); err != nil {
        response.Error(c, err.Error()) // 自动返回中文错误
        return
    }
    
    // 业务逻辑...
}
```

### 2. 带额外验证

```go
func CreateUserHandler(c *gin.Context) {
    var req dto.CreateUserRequest
    
    // 绑定 + 额外的自定义验证
    if err := binding.BindAndValidate(c, &req); err != nil {
        response.Error(c, err.Error())
        return
    }
    
    // 业务逻辑...
}
```

---

## 验证标签使用

### 内置验证规则

```go
type OrderCreateRequest struct {
    ContractNo   string  `json:"contract_no" binding"required"`           // 必填
    CustomerName string  `json:"customer_name" binding"required,min=2"`   // 必填，最少2个字符
    Email        string  `json:"email" binding"email"`                     // 邮箱格式
    Age          int     `json:"age" binding"gte=18,lte=100"`             // 年龄范围 18-100
    Status       string  `json:"status" binding"oneof=active inactive"`   // 枚举值
    Price        float64 `json:"price" binding"gt=0"`                      // 大于0
}
```

### 自定义验证规则（内置）

```go
type UserRequest struct {
    Mobile string `json:"mobile" binding"mobile"`   // 中国大陆手机号
    IDCard string `json:"idcard" binding"idcard"`   // 身份证号
}
```

### 错误提示示例

**请求：**
```json
{
    "contract_no": "",
    "email": "invalid-email",
    "age": 15
}
```

**返回错误（自动格式化为中文）：**
```
字段 'ContractNo' 为必填项; 字段 'Email' 必须是有效的邮箱地址; 字段 'Age' 必须大于等于 18
```

---

## API 参考

### 绑定方法

#### `BindAll(c *gin.Context, req interface{}) error`

统一绑定所有参数（推荐使用）。

**适用场景：** 绝大多数情况

```go
var req dto.OrderUpdateRequest
if err := binding.BindAll(c, &req); err != nil {
    return err
}
```

---

#### `BindAndValidate(c *gin.Context, req interface{}) error`

绑定参数 + 执行额外验证。

**适用场景：** 需要复杂业务验证逻辑时

```go
type ComplexRequest struct {
    StartDate string `json:"start_date" binding"required"`
    EndDate   string `json:"end_date" binding"required"`
}

func (r *ComplexRequest) Validate() error {
    // 自定义验证：结束日期必须晚于开始日期
    start, _ := time.Parse("2006-01-02", r.StartDate)
    end, _ := time.Parse("2006-01-02", r.EndDate)
    if end.Before(start) {
        return fmt.Errorf("结束日期不能早于开始日期")
    }
    return nil
}

func Handler(c *gin.Context) {
    var req ComplexRequest
    if err := binding.BindAndValidate(c, &req); err != nil {
        response.Error(c, err.Error())
        return
    }
    
    // 额外的自定义验证
    if err := req.Validate(); err != nil {
        response.Error(c, err.Error())
        return
    }
}
```

---

#### 单独绑定方法

```go
// 仅绑定 JSON
binding.BindJSON(c, &req)

// 仅绑定 URI
binding.BindUri(c, &req)

// 仅绑定 Query
binding.BindQuery(c, &req)
```

---

### 验证方法

#### `Validate(req interface{}) error`

手动验证结构体（不绑定参数）。

```go
user := &models.User{
    Name:   "张三",
    Mobile: "13800138000",
}

if err := binding.Validate(user); err != nil {
    return err
}
```

---

#### `FormatValidationError(err error) error`

格式化验证错误为中文。

```go
if err := someValidation(); err != nil {
    return binding.FormatValidationError(err)
}
```

---

## 注册自定义验证规则

### 方式 1：在 binding 包初始化时注册

编辑 `core/binding/binding.go`：

```go
func registerCustomValidators() {
    validate.RegisterValidation("mobile", validateMobile)
    validate.RegisterValidation("idcard", validateIDCard)
    
    // 添加新的验证规则
    validate.RegisterValidation("username", validateUsername)
}

// validateUsername 验证用户名格式
func validateUsername(fl validator.FieldLevel) bool {
    username := fl.Field().String()
    // 只允许字母、数字、下划线，长度3-20
    matched, _ := regexp.MatchString(`^[a-zA-Z0-9_]{3,20}$`, username)
    return matched
}
```

### 方式 2：在应用启动时动态注册

在 `main.go` 中：

```go
import "mule-cloud/core/binding"

func main() {
    // 注册自定义验证规则
    binding.RegisterValidation("contract_no", func(fl validator.FieldLevel) bool {
        contractNo := fl.Field().String()
        // 验证合同号格式：YYYYMMDDXXXX
        matched, _ := regexp.MatchString(`^\d{12}$`, contractNo)
        return matched
    })
    
    // 启动服务...
}
```

### 方式 3：获取验证器实例注册

```go
import "mule-cloud/core/binding"

func init() {
    v := binding.GetValidator()
    
    v.RegisterValidation("custom_rule", func(fl validator.FieldLevel) bool {
        // 自定义验证逻辑
        return true
    })
}
```

---

## 完整示例

### DTO 定义

```go
// app/order/dto/order.go
package dto

type OrderCreateRequest struct {
    // URI 参数
    TenantID string `uri:"tenant_id" binding"required"`
    
    // Query 参数
    Source string `query:"source" binding"oneof=web app h5"`
    
    // JSON Body 参数
    ContractNo   string  `json:"contract_no" binding"required,len=12"`
    CustomerName string  `json:"customer_name" binding"required,min=2,max=50"`
    Mobile       string  `json:"mobile" binding"required,mobile"`      // 使用自定义验证
    Email        string  `json:"email" binding"omitempty,email"`        // 可选，但格式必须正确
    Quantity     int     `json:"quantity" binding"required,gt=0"`
    UnitPrice    float64 `json:"unit_price" binding"required,gte=0"`
    DeliveryDate string  `json:"delivery_date" binding"required"`
}
```

### Handler 实现

```go
// app/order/transport/order.go
package transport

import (
    "mule-cloud/app/order/dto"
    "mule-cloud/app/order/endpoint"
    "mule-cloud/app/order/services"
    "mule-cloud/core/binding"
    "mule-cloud/core/response"

    "github.com/gin-gonic/gin"
)

func CreateOrderHandler(svc services.IOrderService) gin.HandlerFunc {
    return func(c *gin.Context) {
        var req dto.OrderCreateRequest
        
        // 统一绑定并验证参数
        if err := binding.BindAll(c, &req); err != nil {
            response.Error(c, err.Error())
            return
        }

        ep := endpoint.CreateOrderEndpoint(svc)
        resp, err := ep(c.Request.Context(), req)
        if err != nil {
            response.Error(c, err.Error())
            return
        }

        response.Success(c, resp)
    }
}
```

### 路由注册

```go
r := gin.Default()
r.POST("/tenants/:tenant_id/orders", CreateOrderHandler(orderSvc))
```

### 测试请求

```bash
curl -X POST "http://localhost:8080/tenants/T001/orders?source=web" \
  -H "Content-Type: application/json" \
  -d '{
    "contract_no": "202510070001",
    "customer_name": "张三",
    "mobile": "13800138000",
    "email": "zhangsan@example.com",
    "quantity": 100,
    "unit_price": 29.99,
    "delivery_date": "2025-11-01"
  }'
```

---

## 常用验证标签速查

| 标签 | 说明 | 示例 |
|------|------|------|
| `required` | 必填 | `binding"required"` |
| `omitempty` | 可选（为空时不验证） | `binding"omitempty,email"` |
| `min=N` | 最小值/长度 | `binding"min=2"` |
| `max=N` | 最大值/长度 | `binding"max=100"` |
| `len=N` | 固定长度 | `binding"len=11"` |
| `gt=N` | 大于 | `binding"gt=0"` |
| `gte=N` | 大于等于 | `binding"gte=18"` |
| `lt=N` | 小于 | `binding"lt=100"` |
| `lte=N` | 小于等于 | `binding"lte=100"` |
| `eq=N` | 等于 | `binding"eq=1"` |
| `ne=N` | 不等于 | `binding"ne=0"` |
| `oneof=A B C` | 枚举值 | `binding"oneof=active inactive"` |
| `email` | 邮箱格式 | `binding"email"` |
| `url` | URL 格式 | `binding"url"` |
| `uuid` | UUID 格式 | `binding"uuid"` |
| `mobile` | 手机号（自定义） | `binding"mobile"` |
| `idcard` | 身份证号（自定义） | `binding"idcard"` |

---

## 错误信息中英文对照

| 验证规则 | 中文提示 | 英文提示 |
|---------|---------|---------|
| `required` | 字段 'Name' 为必填项 | Field validation for 'Name' failed on the 'required' tag |
| `email` | 字段 'Email' 必须是有效的邮箱地址 | Field validation for 'Email' failed on the 'email' tag |
| `min=2` | 字段 'Name' 最小值为 2 | Field validation for 'Name' failed on the 'min' tag |
| `gt=0` | 字段 'Quantity' 必须大于 0 | Field validation for 'Quantity' failed on the 'gt' tag |

---

## 最佳实践

### 1. 统一使用 BindAll

```go
// ✅ 推荐
if err := binding.BindAll(c, &req); err != nil {
    response.Error(c, err.Error())
    return
}

// ❌ 不推荐（容易出错）
if err := c.ShouldBindUri(&req); err != nil { ... }
if err := c.ShouldBindJSON(&req); err != nil { ... }
```

### 2. 合理使用 omitempty

```go
type UpdateRequest struct {
    Name  string `json:"name" binding"omitempty,min=2"`  // 可选，但如果提供则必须至少2个字符
    Email string `json:"email" binding"omitempty,email"` // 可选，但如果提供则必须是邮箱格式
}
```

### 3. 复杂验证单独处理

```go
if err := binding.BindAll(c, &req); err != nil {
    response.Error(c, err.Error())
    return
}

// 复杂的业务验证逻辑
if err := validateBusinessRules(&req); err != nil {
    response.Error(c, err.Error())
    return
}
```

### 4. 自定义验证规则命名规范

- 使用小写加下划线：`custom_rule`
- 见名知意：`contract_no`, `mobile`, `id_card`
- 避免与内置规则冲突

---

## 性能考虑

- ✅ 验证器实例在包初始化时创建，全局复用
- ✅ 验证规则注册一次，后续快速查找
- ✅ 错误格式化仅在验证失败时执行
- ⚠️ 避免在验证函数中执行耗时操作（如数据库查询）

---

## 迁移指南

### 从旧代码迁移

**之前：**
```go
var req dto.OrderUpdateRequest
if err := c.ShouldBindUri(&req); err != nil {
    response.Error(c, "参数错误: "+err.Error())
    return
}
if err := c.ShouldBindJSON(&req); err != nil {
    response.Error(c, "参数错误: "+err.Error())
    return
}
```

**现在：**
```go
var req dto.OrderUpdateRequest
if err := binding.BindAll(c, &req); err != nil {
    response.Error(c, err.Error()) // 已自动格式化为中文
    return
}
```

---

## FAQ

**Q: 为什么要使用 BindAll 而不是直接用 Gin 的绑定方法？**  
A: BindAll 解决了 URI 和 Body 参数绑定的顺序问题，避免验证失败，并自动格式化错误为中文。

**Q: 如何禁用某个字段的验证？**  
A: 使用 `binding"-"` 或 `binding"-"` 标签。

**Q: 自定义验证规则可以访问其他字段吗？**  
A: 可以，使用 `validator.StructLevel` 验证。参考 validator v10 文档。

**Q: 如何国际化错误信息？**  
A: 修改 `formatFieldError` 函数，根据语言环境返回不同的错误信息。

---

## 相关链接

- [validator v10 文档](https://pkg.go.dev/github.com/go-playground/validator/v10)
- [Gin 绑定文档](https://gin-gonic.com/docs/examples/binding-and-validation/)
- [内部文档：统一参数绑定使用指南](../../docs/统一参数绑定使用指南.md)

