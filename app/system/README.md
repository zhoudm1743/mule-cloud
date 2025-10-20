# System Service - 系统服务

## 📋 概述

System 微服务提供系统级功能，包括操作日志管理、系统监控等功能。

## 🏗️ 架构

遵循项目标准的四层架构：

```
app/system/
├── dto/                # 数据传输对象
│   └── operation_log.go
├── services/           # 业务逻辑层
│   └── operation_log.go
├── endpoint/           # 端点层
│   └── operation_log.go
└── transport/          # HTTP传输层
    └── operation_log.go
```

## 📦 功能模块

### 1. 操作日志管理

#### 功能列表
- ✅ 操作日志列表查询（分页、筛选）
- ✅ 操作日志详情查看
- ✅ 操作日志统计分析

#### API 端点

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/admin/system/operation-logs` | 获取操作日志列表 |
| GET | `/admin/system/operation-logs/:id` | 获取操作日志详情 |
| GET | `/admin/system/operation-logs/stats` | 获取操作日志统计 |

#### 请求参数

**列表查询参数：**
```typescript
{
  page: number            // 页码（必填）
  page_size: number       // 每页数量（必填）
  user_id?: string        // 用户ID过滤
  username?: string       // 用户名过滤（模糊查询）
  method?: string         // HTTP方法过滤
  resource?: string       // 资源名称过滤（模糊查询）
  action?: string         // 操作类型过滤
  response_code?: number  // 响应状态码过滤
  start_time?: number     // 开始时间（Unix时间戳）
  end_time?: number       // 结束时间（Unix时间戳）
}
```

**统计查询参数：**
```typescript
{
  start_time: number  // 开始时间（必填）
  end_time: number    // 结束时间（必填）
  group_by?: string   // 分组方式: user, action, resource
}
```

#### 响应示例

**列表响应：**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "list": [
      {
        "id": "507f1f77bcf86cd799439011",
        "user_id": "user123",
        "username": "张三",
        "method": "POST",
        "path": "/admin/perms/admins",
        "resource": "admin",
        "action": "create",
        "request_body": "{\"nickname\":\"测试\"}",
        "response_code": 200,
        "duration": 156,
        "ip": "192.168.1.100",
        "user_agent": "Mozilla/5.0...",
        "created_at": "2024-01-01T10:00:00Z"
      }
    ],
    "total": 100,
    "page": 1,
    "page_size": 10
  }
}
```

**统计响应：**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "total": 1000,
    "success_num": 950,
    "fail_num": 50,
    "avg_time": 123.45,
    "top_users": [
      {
        "user_id": "user123",
        "username": "张三",
        "count": 150
      }
    ],
    "top_actions": [
      {
        "action": "create",
        "count": 300
      }
    ]
  }
}
```

## 🚀 启动服务

### 1. 编译服务

```bash
go build -o bin/system.exe cmd/system/main.go
```

### 2. 运行服务

```bash
# 使用默认配置
./bin/system.exe

# 指定配置文件
./bin/system.exe -config config/system.yaml
```

### 3. 使用启动脚本

```powershell
# 启动所有服务（包括 system）
./start.ps1

# 在运行中的脚本中启动 system 服务
start system

# 重启 system 服务
restart system

# 停止 system 服务
stop system
```

## ⚙️ 配置说明

配置文件：`config/system.yaml`

```yaml
server:
  name: "systemservice"
  port: 8089

consul:
  enabled: true
  service_name: "systemservice"
  service_port: 8089

mongodb:
  enabled: true
  host: "127.0.0.1"
  port: 27015
  database: "mule"
```

## 🎨 前端集成

### 页面路由

- 操作日志列表：`/system/operation-log`
- 操作日志统计：`/system/operation-log/stats`

### 组件结构

```
frontend/src/views/system/operation-log/
├── index.vue                    # 列表页面
├── stats.vue                    # 统计页面
└── components/
    └── DetailDrawer.vue         # 详情抽屉
```

### 类型定义

```
frontend/src/typings/api/
└── operation-log.d.ts           # TypeScript 类型定义
```

### API 服务

```
frontend/src/service/api/
└── operation-log.ts             # API 调用封装
```

## 📊 数据存储

### MongoDB 集合

操作日志存储在对应的数据库中：
- 系统管理员操作：存储在 `system` 数据库的 `operation_logs` 集合
- 租户用户操作：存储在租户数据库的 `operation_logs` 集合

### 索引

```javascript
// 用户ID索引
{ user_id: 1 }

// 创建时间索引（降序）
{ created_at: -1 }

// 复合索引
{ user_id: 1, created_at: -1 }
```

## 🔒 权限控制

操作日志功能需要通过认证中间件：
- 需要登录（JWT 认证）
- 支持租户上下文切换
- 自动记录当前用户的租户信息

## 📝 日志记录

操作日志由中间件自动记录，无需手动调用。

中间件：`core/middleware/operation_log.go`

自动记录以下信息：
- 用户信息（ID、用户名）
- 请求信息（方法、路径、请求体）
- 响应信息（状态码、耗时）
- 客户端信息（IP、User Agent）

## 🐛 故障排查

### 1. 服务无法启动

检查配置文件：
```bash
cat config/system.yaml
```

检查端口占用：
```bash
netstat -ano | findstr ":8089"
```

### 2. 数据库连接失败

检查 MongoDB 连接：
```bash
mongosh --host 127.0.0.1 --port 27015 -u root -p
```

### 3. Consul 注册失败

检查 Consul 服务：
```bash
curl http://127.0.0.1:8500/v1/health/service/systemservice
```

## 📚 相关文档

- [操作日志中间件使用指南](../../docs/操作日志中间件使用指南.md)
- [租户数据库隔离方案](../../docs/数据库级别租户隔离改造方案.md)
- [中间件使用指南](../../docs/中间件极简使用指南.md)

## 🎯 下一步计划

- [ ] 添加日志导出功能
- [ ] 添加日志归档功能
- [ ] 添加实时日志监控
- [ ] 添加异常告警功能
- [ ] 添加日志分析报表

## 👥 维护者

Mule-Cloud 开发团队

