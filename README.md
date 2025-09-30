# Mule-Cloud 微服务项目

基于 **Go-Kit** + **Gin** + **Consul** + **JWT** 的微服务架构示例项目。

## 🎯 项目特性

- ✅ **三层架构**: Service → Endpoint → Transport
- ✅ **JWT认证**: 基于角色的权限控制
- ✅ **API网关**: 统一入口、路由转发、认证鉴权、限流保护
- ✅ **Consul服务发现**: 自动服务注册与发现
- ✅ **Hystrix熔断器**: 服务降级、超时控制、并发限制
- ✅ **统一响应**: 统一返回格式、统一错误处理
- ✅ **配置管理**: Viper + YAML，支持环境变量覆盖
- ✅ **CORS支持**: 跨域请求处理

## 📁 项目结构

```
mule-cloud/
├── core/                    # 核心工具库
│   ├── config/             # 配置管理
│   ├── jwt/                # JWT认证
│   ├── consul/             # Consul集成
│   ├── hystrix/            # Hystrix熔断器
│   └── response/           # 统一响应
├── config/                  # 配置文件目录
│   ├── gateway.yaml        # 网关配置
│   ├── basic.yaml          # Basic服务配置
│   └── test.yaml           # Test服务配置
├── gateway/                 # API网关
│   ├── middleware/         # 中间件（认证、限流、熔断、CORS）
│   └── main.go            # 网关启动
├── test/                    # Test服务（需要认证）
│   ├── services/          # 业务逻辑
│   ├── endpoint/          # 端点层
│   ├── transport/         # HTTP处理
│   └── cmd/               # 启动入口
├── basic/                   # Basic服务（公开访问）
│   ├── services/
│   ├── endpoint/
│   ├── transport/
│   └── cmd/
├── scripts/                 # 脚本
│   ├── start_all.bat      # 启动所有服务
│   ├── build_all.bat      # 编译所有服务
│   └── test_services.bat  # 测试服务
└── docs/                    # 文档
    ├── 架构说明.md
    ├── API网关指南.md
    └── 快速开始.md
```

## 🚀 快速开始

### 1. 前置条件

- ✅ **Go 1.21+**
- ✅ **Consul** ([下载](https://www.consul.io/downloads))
- ✅ **curl** 或 **Postman** (测试用)

### 2. 安装依赖

```bash
# 克隆项目
cd mule-cloud

# 安装Go依赖
go mod tidy
```

### 3. 启动Consul

```bash
# 开发模式启动
consul agent -dev
```

访问 Consul UI: http://localhost:8500

### 4. 启动所有服务

**方式1: 使用配置文件启动**

```bash
# 终端1: Test HTTP服务
cd test/cmd
go run main.go -config=../../config/test.yaml
# 监听: :8000

# 终端2: Basic HTTP服务
cd basic/cmd
go run main.go -config=../../config/basic.yaml
# 监听: :8001

# 终端3: API网关
cd gateway
go run main.go -config=config/gateway.yaml
# 监听: :8080
```

**方式2: 一键启动（推荐）**
```bash
.\scripts\start_all.bat
```

### 5. 测试服务

```bash
# 自动测试脚本
.\scripts\test_services.bat

# 或手动测试
# 健康检查
curl http://localhost:8080/gateway/health

# 公开接口（无需认证）
curl http://localhost:8080/basic/color/1

# 登录获取Token
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d "{\"username\":\"admin\",\"password\":\"admin123\"}"

# 使用Token访问受保护接口
curl -H "Authorization: Bearer {your-token}" \
  http://localhost:8080/test/admin/1
```

## 📚 核心概念

### 架构图

```
┌─────────────┐
│  前端/客户端 │
└──────┬──────┘
       │ HTTP + JWT
       ↓
┌─────────────────────┐
│    API网关 (:8080)  │
│  • JWT认证          │
│  • 路由转发         │
│  • 限流保护         │
│  • Hystrix熔断      │
│  • CORS支持         │
└──────┬──────────────┘
       │
       ├─→ test (:8000) ✅ 需要认证
       └─→ basic (:8001) 🌍 公开
                 ↑
            Consul (:8500)
```

### 路由配置

| 路径 | 服务 | 认证 | 说明 |
|------|------|------|------|
| `/api/login` | 网关 | ❌ | 用户登录 |
| `/gateway/health` | 网关 | ❌ | 健康检查 |
| `/basic/*` | basicservice | ❌ | 公开访问（颜色、尺寸） |
| `/test/*` | testservice | ✅ | 需要登录（管理员CRUD） |

### 测试账号

| 用户名 | 密码 | 角色 | 权限 |
|--------|------|------|------|
| admin | admin123 | admin, user | 所有接口 |
| user | user123 | user | 部分接口 |

## 🔐 JWT认证流程

```
1. 用户登录 → POST /api/login
2. 获得Token
3. 后续请求带上Token
   Header: Authorization: Bearer {token}
4. 网关验证Token
5. 提取用户信息传递给后端服务
   Header: X-User-ID, X-Username
```

## 🔧 开发指南

### 添加新接口

**步骤**:
1. 在 `services/` 添加业务逻辑
2. 在 `endpoint/` 添加Endpoint函数
3. 在 `transport/` 添加Handler
4. 在 `cmd/main.go` 注册路由

详见: [架构说明.md](docs/架构说明.md)

### 修改网关配置

编辑 `gateway/main.go`:

```go
routes: map[string]*RouteConfig{
    "/your-service": {
        ServiceName: "your-service",
        RequireAuth: true,
    },
}
```

## 📖 文档

- 📘 [架构说明](docs/架构说明.md) - 三层架构详解
- 📗 [API网关指南](docs/API网关指南.md) - 网关配置和使用
- 📙 [快速开始](docs/快速开始.md) - 5分钟快速体验
- 📕 [Consul集成指南](docs/Consul集成指南.md) - 服务注册发现
- 📓 [快速开发指南](docs/快速开发指南.md) - 添加新接口的模板
- 🔥 [Hystrix集成指南](docs/Hystrix集成指南.md) - 熔断器配置和使用
- ⚙️ [配置文件指南](docs/配置文件指南.md) - Viper配置管理
- 💾 [MongoDB-Redis-Logger使用指南](docs/MongoDB-Redis-Logger使用指南.md) - 数据库、缓存、日志
- 🎯 [全局实例使用指南](docs/全局实例使用指南.md) - 懒加载全局实例（推荐）

## 🧪 测试API

### 1. 公开接口

```bash
# 获取颜色
GET http://localhost:8080/basic/color/1

# 获取所有颜色
GET http://localhost:8080/basic/color

# 获取尺寸
GET http://localhost:8080/basic/size/2

# 获取所有尺寸
GET http://localhost:8080/basic/size
```

### 2. 认证接口

```bash
# 登录
POST http://localhost:8080/api/login
Content-Type: application/json

{
  "username": "admin",
  "password": "admin123"
}

# 响应
{
  "code": 0,
  "msg": "登录成功",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "username": "admin",
    "roles": ["admin", "user"]
  }
}

# 使用Token访问管理员接口
GET http://localhost:8080/test/admin/1
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...

# 创建管理员
POST http://localhost:8080/test/admin
Authorization: Bearer {token}
Content-Type: application/json

{
  "name": "新管理员",
  "email": "new@example.com",
  "role": "manager"
}

# 更新管理员
PUT http://localhost:8080/test/admin/1
Authorization: Bearer {token}
Content-Type: application/json

{
  "name": "更新的名字",
  "email": "updated@example.com",
  "role": "manager"
}

# 删除管理员
DELETE http://localhost:8080/test/admin/3
Authorization: Bearer {token}
```

## 🛠️ 常用命令

```bash
# 启动所有服务
.\scripts\start_all.bat

# 编译所有服务
.\scripts\build_all.bat

# 测试所有服务
.\scripts\test_services.bat

# 查看Consul服务
curl http://localhost:8500/v1/catalog/services

# 查看网关健康状态
curl http://localhost:8080/gateway/health

# 清理依赖
go mod tidy
```

## ⚠️ 注意事项

### 生产环境配置

1. **修改JWT密钥**
   ```bash
   export JWT_SECRET="your-super-secret-key-min-32-chars"
   ```

2. **修改服务IP**
   ```bash
   export SERVICE_IP="实际服务器IP"
   export CONSUL_ADDR="consul服务器地址:8500"
   ```

3. **使用HTTPS**
   - 配置SSL证书
   - 修改网关监听端口

4. **配置日志和监控**
   - 添加日志系统（如ELK）
   - 添加监控（如Prometheus）

## 🐛 故障排查

### 服务无法启动

```bash
# 检查端口占用
netstat -ano | findstr "8080"
netstat -ano | findstr "8000"

# 检查Consul是否启动
curl http://localhost:8500/v1/status/leader
```

### Token验证失败

- 检查Token格式: `Bearer {token}`
- 检查Token是否过期（24小时）
- 检查JWT_SECRET是否一致

## 📊 服务端口一览

| 服务 | 端口 | URL | 说明 |
|------|------|-----|------|
| Consul | 8500 | http://localhost:8500 | 服务注册中心 |
| Test服务 | 8000 | http://localhost:8000 | Admin管理服务 |
| Basic服务 | 8001 | http://localhost:8001 | 基础服务（颜色、尺寸） |
| **API网关** | **8080** | **http://localhost:8080** | **统一入口** |

## 📝 更新日志

### v1.3.0 (2025-09-30)
- ✅ 集成 MongoDB 数据库（非关系型）
- ✅ 集成 Redis 缓存
- ✅ 集成 Zap 结构化日志系统
- ✅ 完整的数据库、缓存、日志封装
- ✅ 所有服务使用配置文件启动

### v1.2.0 (2025-09-30)
- ✅ 集成 Hystrix-go 熔断器
- ✅ 统一响应格式和错误处理
- ✅ Viper + YAML 配置管理
- ✅ 环境变量支持

### v1.0.0 (2025-01-01)
- ✅ 基础三层架构
- ✅ JWT认证系统
- ✅ API网关（路由、认证、限流）
- ✅ Consul集成
- ✅ 完整文档

## 🤝 贡献

欢迎提交Issue和Pull Request！

## 📄 许可证

MIT License

---

**快速链接**:
- 📚 [完整文档](docs/)
- 🚀 [快速开始](docs/快速开始.md)
- 🏗️ [架构说明](docs/架构说明.md)

**享受编码！🎉**