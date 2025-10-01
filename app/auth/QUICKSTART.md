# 认证服务快速开始

## 📦 已创建的文件

```
app/auth/
├── dto/
│   └── auth.go              # 请求响应数据结构
├── services/
│   └── auth.go              # 业务逻辑（登录、注册、JWT）
├── endpoint/
│   └── auth.go              # 业务端点
├── transport/
│   └── auth.go              # HTTP 处理器
├── README.md                # 详细文档
└── QUICKSTART.md            # 本文件

cmd/auth/
└── main.go                  # 服务启动入口

config/
└── auth.yaml                # 服务配置

scripts/
├── init_auth_users.js       # MongoDB 初始化脚本
└── test_auth_api.sh         # API 测试脚本

docs/
└── 认证服务使用指南.md       # 完整使用指南
```

## 🚀 3步启动服务

### 步骤 1: 初始化测试用户

```bash
# 连接到 MongoDB 并运行初始化脚本
mongosh mongodb://root:bgg8384495@localhost:27015/mule?authSource=admin < scripts/init_auth_users.js
```

这会创建 3 个测试账号：
- `13800138000` / `123456` (普通用户)
- `13900139000` / `123456` (管理员)
- `13700137000` / `123456` (编辑员)

### 步骤 2: 启动认证服务

```bash
go run cmd/auth/main.go
```

服务将在 `http://localhost:8002` 启动

### 步骤 3: 测试接口

#### 方式 1: 使用测试脚本（推荐）

```bash
bash scripts/test_auth_api.sh
```

#### 方式 2: 手动测试

```bash
# 1. 登录获取 token
curl -X POST http://localhost:8002/auth/login \
  -H "Content-Type: application/json" \
  -d '{"phone":"13800138000","password":"123456"}'

# 2. 使用返回的 token 访问其他接口
# 将上面返回的 token 替换到下面的 YOUR_TOKEN
curl -X GET http://localhost:8002/auth/profile \
  -H "Authorization: Bearer YOUR_TOKEN"
```

## 📝 API 列表

### 公开接口（无需认证）

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /auth/register | 用户注册 |
| POST | /auth/login | 用户登录 |
| POST | /auth/refresh | 刷新 Token |
| GET  | /health | 健康检查 |

### 需要认证的接口

| 方法 | 路径 | 说明 |
|------|------|------|
| GET  | /auth/profile | 获取个人信息 |
| PUT  | /auth/profile | 更新个人信息 |
| POST | /auth/password | 修改密码 |

## 🔑 核心功能

### 1. 用户注册

```json
POST /auth/register
{
  "phone": "13812345678",
  "password": "123456",
  "nickname": "新用户",
  "email": "user@example.com"
}
```

### 2. 用户登录

```json
POST /auth/login
{
  "phone": "13800138000",
  "password": "123456"
}
```

返回包含 JWT Token：
```json
{
  "code": 0,
  "data": {
    "token": "eyJhbGc...",
    "user_id": "13800138000",
    "nickname": "测试用户",
    "role": ["user"],
    "expires_at": 1696147200
  }
}
```

### 3. 访问受保护接口

在请求头中携带 Token：

```
Authorization: Bearer eyJhbGc...
```

## 🔧 配置说明

编辑 `config/auth.yaml`：

```yaml
server:
  port: 8002              # 服务端口

jwt:
  secret_key: "..."       # JWT 密钥（生产环境必须修改！）
  expire_time: 24         # Token 过期时间（小时）

mongodb:
  enabled: true
  host: "127.0.0.1"
  port: 27015
  username: "root"
  password: "bgg8384495"
  database: "mule"
```

## 📚 更多文档

- [app/auth/README.md](./README.md) - 详细 API 文档
- [docs/认证服务使用指南.md](../../docs/认证服务使用指南.md) - 完整使用指南
- [docs/JWT和gRPC集成指南.md](../../docs/JWT和gRPC集成指南.md) - JWT 集成指南

## ⚠️ 安全提示

**开发环境：**
- ✅ 使用 MD5 密码加密（快速）
- ✅ 简单的 JWT 配置

**生产环境务必：**
- 🔒 使用 bcrypt 替代 MD5
- 🔒 修改 JWT secret_key
- 🔒 启用 HTTPS
- 🔒 添加登录限流
- 🔒 实现 Token 黑名单

## 🐛 故障排除

### MongoDB 连接失败

```bash
# 检查 MongoDB 是否运行
mongosh mongodb://root:bgg8384495@localhost:27015/admin

# 如果连接失败，检查配置文件中的连接信息
```

### 登录失败

```bash
# 确保已运行初始化脚本创建测试用户
mongosh mongodb://root:bgg8384495@localhost:27015/mule?authSource=admin < scripts/init_auth_users.js

# 查看 admins 集合
mongosh mongodb://root:bgg8384495@localhost:27015/mule?authSource=admin
> db.admins.find({}, {password: 0}).pretty()
```

### Token 验证失败

确保：
1. Token 没有过期
2. 请求头格式正确：`Authorization: Bearer <token>`
3. JWT secret_key 配置正确

## 🎯 下一步

1. ✅ 基础认证服务已就绪
2. 📝 阅读详细文档了解更多功能
3. 🔌 集成到其他服务或网关
4. 🛡️ 配置生产环境安全策略
5. 📊 添加监控和日志

## 💡 示例代码

### 在其他服务中使用认证中间件

```go
import (
    jwtPkg "mule-cloud/core/jwt"
    "mule-cloud/app/gateway/middleware"
)

// 初始化 JWT 管理器
jwtManager := jwtPkg.NewJWTManager(
    []byte("your-secret-key"),
    24 * time.Hour,
)

// 使用认证中间件
r := gin.New()
protected := r.Group("/api")
protected.Use(middleware.JWTAuth(jwtManager))
{
    protected.GET("/data", yourHandler)
}
```

### 获取当前用户信息

```go
func YourHandler(c *gin.Context) {
    userID, _ := c.Get("user_id")
    username, _ := c.Get("username")
    roles, _ := c.Get("roles")
    
    log.Printf("用户: %s, ID: %s, 角色: %v", username, userID, roles)
}
```

---

✨ **认证服务已就绪！开始构建安全的应用吧！** ✨

