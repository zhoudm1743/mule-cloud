# JWT认证 + gRPC内部调用 - 完整指南

## 📖 目录
- [架构设计](#架构设计)
- [JWT认证](#jwt认证)
- [gRPC内部调用](#grpc内部调用)
- [完整示例](#完整示例)
- [测试指南](#测试指南)

---

## 架构设计

### 整体架构

```
┌──────────────────────────────────────────────────────┐
│                     前端/客户端                       │
└────────────────────┬─────────────────────────────────┘
                     │ HTTP + JWT Token
                     ↓
┌────────────────────────────────────────────────────────┐
│              API网关 (:8080)                           │
│  - JWT认证                                             │
│  - 路由转发                                            │
│  - 限流保护                                            │
│  - CORS处理                                            │
└──────────┬─────────────────────────────────────────────┘
           │
           ├─→ HTTP (对外API)
           │   ├─→ testservice (:8000)
           │   └─→ basicservice (:8001)
           │
           └─→ gRPC (内部服务间调用)
               ├─→ userservice-grpc (:9000)
               ├─→ orderservice-grpc (:9001)
               └─→ productservice-grpc (:9002)
                     ↑
                     └─── Consul服务发现
```

### 通信方式

| 场景 | 协议 | 说明 |
|------|------|------|
| 前端 → 网关 | HTTP + JWT | RESTful API，JWT认证 |
| 网关 → HTTP服务 | HTTP | 转发HTTP请求 |
| 服务间调用 | gRPC | 高性能内部RPC调用 |

---

## JWT认证

### 1. JWT工具类

已创建 `core/jwt/jwt.go`，提供：
- ✅ 生成Token
- ✅ 验证Token
- ✅ 刷新Token
- ✅ 角色权限检查

### 2. 认证中间件

**位置**: `gateway/middleware/auth.go`

#### JWTAuth - 必须认证
```go
// 使用示例
r.Use(middleware.JWTAuth(jwtManager))
```

#### OptionalAuth - 可选认证
```go
// 有token则验证，没有则跳过
r.Use(middleware.OptionalAuth(jwtManager))
```

#### RequireRole - 角色检查
```go
// 需要admin角色
r.Use(middleware.RequireRole("admin"))

// 需要admin或moderator角色之一
r.Use(middleware.RequireRole("admin", "moderator"))
```

### 3. 登录流程

```mermaid
sequenceDiagram
    客户端->>网关: POST /api/login {username, password}
    网关->>用户服务(gRPC): ValidateCredentials
    用户服务(gRPC)-->>网关: 返回用户信息
    网关->>网关: 生成JWT Token
    网关-->>客户端: 返回Token
    客户端->>网关: GET /test/xxx (Header: Authorization: Bearer {token})
    网关->>网关: 验证Token
    网关->>后端服务: 转发请求 (Header: X-User-ID)
    后端服务-->>网关: 返回响应
    网关-->>客户端: 返回响应
```

### 4. Token格式

**Header**:
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**Payload** (Claims):
```json
{
  "user_id": "123",
  "username": "admin",
  "roles": ["admin", "user"],
  "exp": 1735689600,
  "iat": 1735603200,
  "iss": "mule-cloud"
}
```

---

## gRPC内部调用

### 1. 为什么使用gRPC？

| 特性 | HTTP/REST | gRPC |
|------|-----------|------|
| 性能 | ⭐⭐⭐ | ⭐⭐⭐⭐⭐ |
| 类型安全 | ❌ | ✅ |
| 双向流 | ❌ | ✅ |
| 代码生成 | ❌ | ✅ |
| 传输效率 | JSON (大) | Protobuf (小) |
| 适用场景 | 对外API | 内部服务间 |

### 2. 定义Proto文件

**示例**: `proto/user.proto`

```protobuf
syntax = "proto3";

package user;

service UserService {
  rpc GetUserById(GetUserByIdRequest) returns (GetUserByIdResponse);
  rpc ValidateCredentials(ValidateCredentialsRequest) returns (ValidateCredentialsResponse);
}

message User {
  string id = 1;
  string username = 2;
  repeated string roles = 3;
}
```

### 3. 生成Go代码

```bash
# 安装protoc编译器
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# 生成代码
protoc --go_out=. --go_opt=paths=source_relative \
       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
       proto/user.proto
```

### 4. 实现gRPC服务端

```go
// services/user_grpc.go
package services

import (
    "context"
    pb "mule-cloud/proto/user"
)

type UserGRPCService struct {
    pb.UnimplementedUserServiceServer
}

func (s *UserGRPCService) GetUserById(ctx context.Context, req *pb.GetUserByIdRequest) (*pb.GetUserByIdResponse, error) {
    // 业务逻辑
    user := &pb.User{
        Id:       req.Id,
        Username: "user_" + req.Id,
        Roles:    []string{"user"},
    }
    return &pb.GetUserByIdResponse{User: user}, nil
}

func (s *UserGRPCService) ValidateCredentials(ctx context.Context, req *pb.ValidateCredentialsRequest) (*pb.ValidateCredentialsResponse, error) {
    // 验证用户名密码
    if req.Username == "admin" && req.Password == "admin123" {
        return &pb.ValidateCredentialsResponse{
            Valid: true,
            User: &pb.User{
                Id:       "1",
                Username: "admin",
                Roles:    []string{"admin", "user"},
            },
        }, nil
    }
    return &pb.ValidateCredentialsResponse{Valid: false}, nil
}
```

### 5. 启动gRPC服务器

```go
// cmd/user_grpc/main.go
package main

import (
    "log"
    "mule-cloud/core/grpc"
    pb "mule-cloud/proto/user"
    "mule-cloud/services"
)

func main() {
    // 创建gRPC服务器
    server := grpc.NewServer(&grpc.ServerConfig{Port: 9000})
    
    // 注册服务
    pb.RegisterUserServiceServer(server, &services.UserGRPCService{})
    
    // 启动服务器
    log.Fatal(grpc.StartServer(server, 9000))
}
```

### 6. 调用gRPC服务（客户端）

```go
// 在其他服务中调用UserService
package services

import (
    "context"
    "mule-cloud/core/grpc"
    pb "mule-cloud/proto/user"
)

type OrderService struct {
    grpcClient *grpc.ClientManager
}

func (s *OrderService) CreateOrder(userID string) error {
    // 获取gRPC连接
    conn, err := s.grpcClient.GetConnection("userservice-grpc")
    if err != nil {
        return err
    }
    
    // 创建客户端
    client := pb.NewUserServiceClient(conn)
    
    // 调用远程方法
    ctx := context.Background()
    resp, err := client.GetUserById(ctx, &pb.GetUserByIdRequest{Id: userID})
    if err != nil {
        return err
    }
    
    // 使用响应
    log.Printf("用户信息: %+v", resp.User)
    return nil
}
```

---

## 完整示例

### 项目结构

```
mule-cloud/
├── core/
│   ├── jwt/              # JWT工具
│   │   └── jwt.go
│   ├── grpc/             # gRPC工具
│   │   ├── client.go     # 客户端管理器
│   │   └── server.go     # 服务器帮助函数
│   └── response/         # 统一响应
├── gateway/
│   ├── middleware/       # 中间件
│   │   ├── auth.go       # JWT认证
│   │   ├── cors.go       # CORS
│   │   └── ratelimit.go  # 限流
│   ├── main.go           # 简单版网关
│   └── main_v2.go        # 增强版网关（JWT + 限流）
├── proto/
│   └── user.proto        # Protobuf定义
├── services/
│   ├── user_grpc.go      # gRPC服务实现
│   └── ...
└── cmd/
    ├── gateway/          # 网关启动
    └── user_grpc/        # gRPC服务启动
```

### 启动服务

#### 1. 启动Consul
```bash
consul agent -dev
```

#### 2. 启动gRPC服务（用户服务）
```bash
cd cmd/user_grpc
go run main.go
# 监听: :9000 (gRPC)
```

#### 3. 启动HTTP服务
```bash
# testservice
cd test/cmd
go run main.go
# 监听: :8000 (HTTP)

# basicservice
cd basic/cmd
go run main.go
# 监听: :8001 (HTTP)
```

#### 4. 启动API网关
```bash
cd gateway
go run main_v2.go
# 监听: :8080 (HTTP)
```

---

## 测试指南

### 1. 测试登录（获取Token）

```bash
# 管理员登录
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "admin123"
  }'

# 响应:
# {
#   "code": 0,
#   "msg": "登录成功",
#   "data": {
#     "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
#     "username": "admin",
#     "roles": ["admin", "user"]
#   }
# }

# 普通用户登录
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "user",
    "password": "user123"
  }'
```

### 2. 测试公开接口（无需认证）

```bash
# 访问basic服务（公开）
curl http://localhost:8080/basic/color/1

# 响应: {"code":0,"msg":"success","data":{"result":"red"}}
```

### 3. 测试需要认证的接口

```bash
# 保存token
TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."

# 访问需要认证的接口
curl http://localhost:8080/test/admin/123 \
  -H "Authorization: Bearer $TOKEN"

# 成功响应: {"code":0,"msg":"success","data":{...}}
```

### 4. 测试权限控制

```bash
# 使用admin token访问管理接口
curl http://localhost:8080/admin/users \
  -H "Authorization: Bearer $ADMIN_TOKEN"
# ✅ 成功

# 使用普通用户token访问管理接口
curl http://localhost:8080/admin/users \
  -H "Authorization: Bearer $USER_TOKEN"
# ❌ 失败: {"code":403,"msg":"权限不足"}
```

### 5. 测试限流

```bash
# 快速发送多个请求
for i in {1..150}; do
  curl http://localhost:8080/basic/color/1 &
done

# 部分请求会返回:
# {"code":429,"msg":"请求过于频繁，请稍后再试"}
```

### 6. 测试gRPC调用

```bash
# 安装grpcurl
go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest

# 查看服务
grpcurl -plaintext localhost:9000 list

# 调用gRPC方法
grpcurl -plaintext \
  -d '{"id": "123"}' \
  localhost:9000 user.UserService/GetUserById

# 响应:
# {
#   "user": {
#     "id": "123",
#     "username": "user_123",
#     "roles": ["user"]
#   }
# }
```

---

## 路由配置说明

### 网关路由表

| 路径前缀 | 服务名 | 认证要求 | 角色要求 |
|----------|--------|---------|---------|
| `/basic/*` | basicservice | ❌ 无需认证 | - |
| `/test/*` | testservice | ✅ 需要登录 | - |
| `/admin/*` | testservice | ✅ 需要登录 | admin |

### 配置方式

在 `gateway/main_v2.go` 中修改：

```go
routes: map[string]*RouteConfig{
    "/basic": {
        ServiceName: "basicservice",
        RequireAuth: false,  // 公开访问
    },
    "/test": {
        ServiceName: "testservice",
        RequireAuth: true,   // 需要登录
    },
    "/admin": {
        ServiceName: "testservice",
        RequireAuth: true,
        RequireRole: []string{"admin"},  // 需要admin角色
    },
}
```

---

## 最佳实践

### 1. JWT密钥管理

```bash
# 生产环境使用环境变量
export JWT_SECRET="your-super-secret-key-min-32-chars"

# Docker
docker run -e JWT_SECRET="..." your-gateway
```

### 2. Token刷新策略

```go
// 前端定时刷新token
setInterval(() => {
    fetch('/api/refresh', {
        headers: {
            'Authorization': 'Bearer ' + oldToken
        }
    }).then(res => {
        // 更新token
        localStorage.setItem('token', res.data.token);
    });
}, 20 * 60 * 1000); // 每20分钟刷新
```

### 3. gRPC连接管理

- ✅ 使用连接池（`ClientManager`已实现）
- ✅ 健康检查和自动重连
- ✅ 超时控制
- ✅ 错误重试

### 4. 安全建议

- ✅ HTTPS（生产环境必须）
- ✅ 强密钥（至少32字符）
- ✅ Token过期时间合理（建议1-24小时）
- ✅ 敏感操作二次验证
- ✅ 日志记录（不记录密码）

---

## 常见问题

### Q1: JWT vs Session，为什么选JWT？

| 特性 | Session | JWT |
|------|---------|-----|
| 状态 | 服务端保存 | 无状态 |
| 扩展性 | 需要共享存储 | ✅ 无需存储 |
| 性能 | 需要查询存储 | ✅ 直接验证 |
| 适合场景 | 单体应用 | ✅ 微服务 |

### Q2: HTTP vs gRPC，如何选择？

```
对外API → HTTP/REST (易用、通用)
内部调用 → gRPC (高性能、类型安全)
```

### Q3: 如何吊销JWT Token？

JWT是无状态的，无法直接吊销。解决方案：
1. **短期Token** + 刷新Token
2. **黑名单**（Redis存储已吊销的Token）
3. **版本号**（用户密码修改后增加版本号）

```go
// 黑名单示例
func (m *JWTManager) RevokeToken(tokenString string) {
    redis.Set("blacklist:"+tokenString, "1", m.tokenDuration)
}

func (m *JWTManager) IsRevoked(tokenString string) bool {
    return redis.Exists("blacklist:" + tokenString)
}
```

### Q4: gRPC服务如何注册到Consul？

```go
// 注册gRPC服务
consulClient.RegisterService(&cousul.ServiceConfig{
    ServiceName: "userservice-grpc",
    ServicePort: 9000,
    Tags:        []string{"grpc", "user"},
    HealthCheck: &cousul.HealthCheck{
        GRPC:     "192.168.1.10:9000",  // gRPC健康检查
        Interval: "10s",
    },
})
```

---

## 监控和日志

### 请求日志格式

```
[网关转发] GET /test/admin/123 → http://192.168.31.78:8000/admin/123 (服务: testservice, 用户: admin)
[网关响应] GET /test/admin/123 耗时: 15.2ms
[gRPC] 请求: /user.UserService/GetUserById
```

### 监控指标

建议监控：
- 请求QPS
- 响应时间
- 错误率
- Token验证失败率
- gRPC连接状态

---

## 下一步

- [ ] 添加Prometheus监控
- [ ] 添加分布式追踪（Jaeger）
- [ ] 实现Token刷新接口
- [ ] 添加Redis缓存层
- [ ] 实现更多gRPC服务
- [ ] 添加服务降级和熔断

---

## 总结

✅ **已实现功能**:
- JWT认证和角色权限
- API网关（路由、认证、限流）
- gRPC内部调用框架
- Consul服务发现集成
- CORS支持

🎯 **推荐架构**:
```
前端 --HTTP+JWT--> 网关 --HTTP--> HTTP服务（对外API）
                          --gRPC--> gRPC服务（内部调用）
```

这种架构兼顾了**对外易用性**（HTTP）和**内部高性能**（gRPC）！
