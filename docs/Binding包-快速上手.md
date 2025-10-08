# Binding 包 - 快速上手

## 🎯 核心特性

已在 `core/binding` 包中实现了集成 validator v10 的统一参数绑定方案：

✅ **一行代码** - 统一绑定 URI、Query、Body 参数  
✅ **中文错误** - 自动格式化验证错误为友好的中文提示  
✅ **自定义验证** - 内置手机号、身份证验证，支持扩展  
✅ **避免 Bug** - 解决了 Gin 原生绑定的参数覆盖问题  

---

## 🚀 快速开始

### 1. 导入包

```go
import "mule-cloud/core/binding"
```

### 2. 在 Handler 中使用

```go
func UpdateOrderHandler(svc services.IOrderService) gin.HandlerFunc {
    return func(c *gin.Context) {
        var req dto.OrderUpdateRequest
        
        // 一行代码搞定所有参数绑定和验证
        if err := binding.BindAll(c, &req); err != nil {
            response.Error(c, err.Error())
            return
        }
        
        // 业务逻辑...
    }
}
```

### 3. 定义 DTO（支持更多验证规则）

```go
type OrderCreateRequest struct {
    // URI 参数
    ID string `uri:"id" binding:"required"`
    
    // JSON 参数（带验证规则）
    ContractNo   string  `json:"contract_no" binding:"required,len=12"`
    CustomerName string  `json:"customer_name" binding:"required,min=2,max=50"`
    Mobile       string  `json:"mobile" binding:"required,mobile"`      // 自定义验证：手机号
    Email        string  `json:"email" binding:"omitempty,email"`        // 可选，但必须是邮箱格式
    Quantity     int     `json:"quantity" binding:"required,gt=0"`       // 必须大于0
    UnitPrice    float64 `json:"unit_price" binding:"required,gte=0"`    // 必须大于等于0
    Status       string  `json:"status" binding:"oneof=active inactive"` // 枚举值
}
```

---

## 📊 对比旧方法

### ❌ 旧方法（有问题）

```go
// 8 行代码，容易出错
if err := c.ShouldBindUri(&req); err != nil {
    response.Error(c, "参数错误: "+err.Error())
    return
}
if err := c.ShouldBindJSON(&req); err != nil {  // 会导致验证失败！
    response.Error(c, "参数错误: "+err.Error())
    return
}

// 错误信息不友好：
// "Key: 'OrderStyleRequest.StyleID' Error:Field validation for 'StyleID' failed on the 'required' tag"
```

### ✅ 新方法（推荐）

```go
// 4 行代码，自动处理
if err := binding.BindAll(c, &req); err != nil {
    response.Error(c, err.Error())  // 自动中文错误
    return
}

// 友好的中文错误：
// "字段 'StyleID' 为必填项"
```

---

## 🎨 内置验证规则

### 常用规则

```go
binding:"required"              // 必填
binding:"min=2"                 // 最小值/长度
binding:"max=100"               // 最大值/长度
binding:"len=11"                // 固定长度
binding:"gt=0"                  // 大于
binding:"gte=0"                 // 大于等于
binding:"email"                 // 邮箱格式
binding:"oneof=active inactive" // 枚举值
binding:"omitempty"             // 可选
```

### 自定义规则（已内置）

```go
binding:"mobile"  // 中国大陆手机号（11位，1开头）
binding:"idcard"  // 身份证号（15位或18位）
```

---

## 🔧 添加自定义验证规则

### 方法 1：在 `binding.go` 中添加

编辑 `core/binding/binding.go`：

```go
func registerCustomValidators() {
    validate.RegisterValidation("mobile", validateMobile)
    validate.RegisterValidation("idcard", validateIDCard)
    
    // 添加新规则
    validate.RegisterValidation("username", validateUsername)
}

func validateUsername(fl validator.FieldLevel) bool {
    username := fl.Field().String()
    // 只允许字母数字下划线，3-20位
    matched, _ := regexp.MatchString(`^[a-zA-Z0-9_]{3,20}$`, username)
    return matched
}
```

### 方法 2：在应用启动时注册

在 `main.go` 或 `init()` 中：

```go
import "mule-cloud/core/binding"

func init() {
    binding.RegisterValidation("contract_no", func(fl validator.FieldLevel) bool {
        contractNo := fl.Field().String()
        // 验证合同号格式
        return len(contractNo) == 12 && contractNo[0] == '2'
    })
}
```

---

## 💡 实际例子

### 订单更新接口

**DTO 定义：**

```go
type OrderStyleRequest struct {
    ID        string             `uri:"id" binding:"required"`
    StyleID   string             `json:"style_id" binding:"required"`
    Colors    []string           `json:"colors" binding:"required,min=1"`
    Sizes     []string           `json:"sizes" binding:"required,min=1"`
    UnitPrice float64            `json:"unit_price" binding:"required,gt=0"`
    Quantity  int                `json:"quantity" binding:"required,gt=0"`
    Items     []models.OrderItem `json:"items" binding:"required,min=1"`
}
```

**Handler：**

```go
func UpdateOrderStyleHandler(svc services.IOrderService) gin.HandlerFunc {
    return func(c *gin.Context) {
        var req dto.OrderStyleRequest
        
        if err := binding.BindAll(c, &req); err != nil {
            response.Error(c, err.Error())
            return
        }

        ep := endpoint.UpdateOrderStyleEndpoint(svc)
        resp, err := ep(c.Request.Context(), req)
        if err != nil {
            response.Error(c, err.Error())
            return
        }

        response.Success(c, resp)
    }
}
```

**请求示例：**

```bash
PUT /admin/order/orders/68e48c19b4eb03ee2a2b8dcd/style

{
    "style_id": "68e48be9b4eb03ee2a2b8dcb",
    "colors": ["红色", "蓝色"],
    "sizes": ["L", "XL"],
    "unit_price": 29.99,
    "quantity": 100,
    "items": [
        {"color": "红色", "size": "L", "quantity": 50}
    ]
}
```

**错误响应（如果数据无效）：**

```json
{
    "code": -1,
    "msg": "字段 'StyleID' 为必填项; 字段 'UnitPrice' 必须大于 0",
    "timestamp": 1759808602
}
```

---

## 📝 迁移步骤

### 已迁移的文件

- ✅ `app/order/transport/order.go`
- ✅ `app/order/transport/style.go`

### 待迁移的文件

其他模块的 transport 文件可以逐步迁移：

- `app/basic/transport/*.go`
- `app/perms/transport/*.go`

### 迁移模板

**查找：**
```go
if err := c.ShouldBindUri(&req); err != nil {
    response.Error(c, "参数错误: "+err.Error())
    return
}
if err := c.ShouldBindJSON(&req); err != nil {
    response.Error(c, "参数错误: "+err.Error())
    return
}
```

**替换为：**
```go
if err := binding.BindAll(c, &req); err != nil {
    response.Error(c, err.Error())
    return
}
```

---

## 🎯 效果对比

### 旧错误信息
```
参数错误: Key: 'OrderStyleRequest.StyleID' Error:Field validation for 'StyleID' failed on the 'required' tag
Key: 'OrderStyleRequest.UnitPrice' Error:Field validation for 'UnitPrice' failed on the 'required' tag
```

### 新错误信息
```
字段 'StyleID' 为必填项; 字段 'UnitPrice' 必须大于 0
```

---

## ❓ 常见问题

**Q: 为什么 BindAll 能解决验证失败问题？**  
A: 因为它正确处理了绑定顺序，先绑定 URI，再绑定 Body，避免了参数覆盖。

**Q: 如何让某些字段可选？**  
A: 使用 `binding:"omitempty,xxx"` 标签，如 `binding:"omitempty,email"`。

**Q: 自定义验证规则会影响性能吗？**  
A: 不会。验证器在包初始化时创建，规则注册一次，后续调用非常快。

**Q: 如何添加复杂的业务验证逻辑？**  
A: 在绑定后额外调用自定义验证方法：

```go
if err := binding.BindAll(c, &req); err != nil {
    return err
}

// 复杂业务验证
if err := validateBusinessRules(&req); err != nil {
    return err
}
```

---

## 📚 更多文档

- [完整 API 文档](../core/binding/README.md)
- [统一参数绑定使用指南](./统一参数绑定使用指南.md)
- [validator v10 官方文档](https://pkg.go.dev/github.com/go-playground/validator/v10)

---

## 🎉 总结

使用 `core/binding` 包可以：

1. **减少代码** - 从 8 行减少到 4 行
2. **避免 Bug** - 解决参数绑定顺序问题
3. **更友好** - 中文错误提示
4. **更强大** - 集成 validator v10 的所有功能
5. **可扩展** - 轻松添加自定义验证规则

**立即开始使用吧！** 🚀

