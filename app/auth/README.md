# 认证服务 (Auth Service)

基于 JWT 的用户认证服务，提供登录、注册、Token 管理和用户信息管理功能。

## 功能特性

- ✅ 用户注册
- ✅ 用户登录
- ✅ JWT Token 生成与验证
- ✅ Token 刷新
- ✅ 获取/更新个人信息
- ✅ 修改密码
- ✅ 基于角色的权限控制
- ✅ MongoDB 数据持久化
- ✅ 密码 MD5 加密（可升级为 bcrypt）

## 快速开始

### 1. 启动服务

```bash
# 使用默认配置
go run cmd/auth/main.go

# 使用自定义配置
go run cmd/auth/main.go -config=config/auth.yaml
```

### 2. 初始化测试用户

在 MongoDB 中插入测试用户：

```javascript
use mule

db.admins.insertOne({
  phone: "13800138000",
  password: "e10adc3949ba59abbe56e057f20f883e",  // 123456的MD5
  nickname: "测试用户",
  email: "test@example.com",
  status: 1,
  role: ["user"],
  avatar: "",
  created_at: NumberLong(Date.now() / 1000),
  updated_at: NumberLong(Date.now() / 1000)
})

// 创建管理员用户
db.admins.insertOne({
  phone: "13900139000",
  password: "e10adc3949ba59abbe56e057f20f883e",  // 123456的MD5
  nickname: "管理员",
  email: "admin@example.com",
  status: 1,
  role: ["admin", "user"],
  avatar: "",
  created_at: NumberLong(Date.now() / 1000),
  updated_at: NumberLong(Date.now() / 1000)
})
```

## API 接口

### 公开接口（无需认证）

#### 1. 用户注册

```http
POST /auth/register
Content-Type: application/json

{
  "phone": "13800138000",
  "password": "123456",
  "nickname": "测试用户",
  "email": "test@example.com"
}
```

**响应示例：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "user_id": "13800138000",
    "phone": "13800138000",
    "nickname": "测试用户",
    "message": "注册成功"
  }
}
```

#### 2. 用户登录

```http
POST /auth/login
Content-Type: application/json

{
  "phone": "13800138000",
  "password": "123456"
}
```

**响应示例：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user_id": "13800138000",
    "phone": "13800138000",
    "nickname": "测试用户",
    "avatar": "",
    "role": ["user"],
    "expires_at": 1696147200
  }
}
```

#### 3. 刷新 Token

```http
POST /auth/refresh
Content-Type: application/json

{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**响应示例：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expires_at": 1696233600
  }
}
```

### 需要认证的接口

所有需要认证的接口都需要在请求头中携带 Token：

```http
Authorization: Bearer <token>
```

#### 4. 获取个人信息

```http
GET /auth/profile
Authorization: Bearer <token>
```

**响应示例：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "user_id": "13800138000",
    "phone": "13800138000",
    "nickname": "测试用户",
    "avatar": "",
    "email": "test@example.com",
    "role": ["user"],
    "status": 1,
    "created_at": 1696060800,
    "updated_at": 1696060800
  }
}
```

#### 5. 更新个人信息

```http
PUT /auth/profile
Authorization: Bearer <token>
Content-Type: application/json

{
  "nickname": "新昵称",
  "avatar": "https://example.com/avatar.jpg",
  "email": "newemail@example.com"
}
```

**响应示例：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "success": true,
    "message": "更新成功"
  }
}
```

#### 6. 修改密码

```http
POST /auth/password
Authorization: Bearer <token>
Content-Type: application/json

{
  "old_password": "123456",
  "new_password": "654321"
}
```

**响应示例：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "success": true,
    "message": "密码修改成功"
  }
}
```

## 错误码说明

| 错误信息 | 说明 |
|---------|------|
| 用户不存在 | 用户未注册或已被删除 |
| 密码错误 | 登录密码不正确 |
| 用户已存在 | 注册时手机号已被使用 |
| 用户已被禁用 | 用户状态为非激活状态 |
| token无效 | Token 格式错误或已失效 |
| 未认证 | 请求头未提供有效的 Token |

## 测试示例

### 使用 curl 测试

```bash
# 1. 注册用户
curl -X POST http://localhost:8002/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "phone": "13800138000",
    "password": "123456",
    "nickname": "测试用户",
    "email": "test@example.com"
  }'

# 2. 登录
TOKEN=$(curl -X POST http://localhost:8002/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "phone": "13800138000",
    "password": "123456"
  }' | jq -r '.data.token')

echo "Token: $TOKEN"

# 3. 获取个人信息
curl -X GET http://localhost:8002/auth/profile \
  -H "Authorization: Bearer $TOKEN"

# 4. 更新个人信息
curl -X PUT http://localhost:8002/auth/profile \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "nickname": "新昵称",
    "email": "newemail@example.com"
  }'

# 5. 修改密码
curl -X POST http://localhost:8002/auth/password \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "old_password": "123456",
    "new_password": "654321"
  }'
```

## 与网关集成

在网关配置文件中添加认证服务路由：

```yaml
gateway:
  routes:
    /auth:
      service_name: "authservice"
      require_auth: false  # 登录/注册不需要认证
      require_role: []
```

网关会自动将请求转发到认证服务，并在需要时验证 JWT Token。

## 安全建议

1. **生产环境密码加密**：将 MD5 替换为 bcrypt
   ```go
   import "golang.org/x/crypto/bcrypt"
   
   func hashPassword(password string) string {
       hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
       return string(hash)
   }
   
   func checkPassword(password, hash string) bool {
       err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
       return err == nil
   }
   ```

2. **JWT Secret Key**：在生产环境中使用强密钥
   - 长度至少 32 字符
   - 使用环境变量或配置中心管理
   - 定期轮换

3. **密码策略**：
   - 最小长度 8 位
   - 包含大小写字母、数字和特殊字符
   - 密码强度验证

4. **限流保护**：对登录/注册接口增加限流

5. **Token 黑名单**：使用 Redis 存储已注销的 Token

## 数据库索引

为提高查询性能，建议创建以下索引：

```javascript
use mule

// 手机号唯一索引
db.admins.createIndex({ "phone": 1 }, { unique: true })

// 邮箱索引
db.admins.createIndex({ "email": 1 })

// 状态索引
db.admins.createIndex({ "status": 1 })
```

## 项目结构

```
app/auth/
├── dto/
│   └── auth.go           # 数据传输对象
├── services/
│   └── auth.go           # 业务逻辑层
├── endpoint/
│   └── auth.go           # 端点层
├── transport/
│   └── auth.go           # HTTP 传输层
└── README.md             # 本文档

cmd/auth/
└── main.go               # 服务启动入口

config/
└── auth.yaml             # 服务配置文件
```

## 常见问题

### Q: Token 过期时间如何配置？
A: 在 `config/auth.yaml` 中修改 `jwt.expire_time`，单位为小时。

### Q: 如何添加更多用户角色？
A: 在注册或更新用户时设置 `role` 字段，例如 `["user", "admin", "editor"]`。

### Q: 如何实现单点登录（SSO）？
A: 可以使用 Redis 存储用户 Token，实现跨服务的 Token 共享。

### Q: 如何实现登出功能？
A: 将 Token 加入 Redis 黑名单，在中间件中检查 Token 是否在黑名单中。

