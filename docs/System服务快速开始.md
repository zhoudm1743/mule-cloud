# System 服务快速开始

## 🎯 概述

System 服务是 Mule-Cloud 的系统管理微服务，提供操作日志管理、系统监控等核心功能。

## 📦 已实现功能

### ✅ 操作日志管理

- **操作日志列表**：支持多条件筛选、分页查询
- **操作日志详情**：查看完整的请求响应信息
- **操作日志统计**：TOP用户、TOP操作、成功率等统计

## 🚀 快速启动

### 方法一：使用启动脚本（推荐）

```powershell
# 启动所有服务（包括 system）
./start.ps1

# 在运行的脚本中管理 system 服务
start system      # 启动
restart system    # 重启
stop system       # 停止
status           # 查看状态
```

### 方法二：手动启动

```bash
# 1. 编译服务
go build -o bin/system.exe cmd/system/main.go

# 2. 运行服务
./bin/system.exe

# 3. 使用自定义配置
./bin/system.exe -config config/system.yaml
```

## 📡 服务信息

- **服务名称**：`systemservice`
- **端口**：`8089`
- **健康检查**：`http://localhost:8089/health`
- **Consul注册**：自动注册到 Consul
- **网关路由**：`/admin/system/*`

## 🌐 API 端点

### 操作日志

| 方法 | 端点 | 说明 |
|------|------|------|
| GET | `/admin/system/operation-logs` | 获取操作日志列表 |
| GET | `/admin/system/operation-logs/:id` | 获取操作日志详情 |
| GET | `/admin/system/operation-logs/stats` | 获取操作日志统计 |

### 示例请求

```bash
# 1. 获取操作日志列表
curl -X GET "http://localhost:8080/admin/system/operation-logs?page=1&page_size=10" \
  -H "Authorization: Bearer YOUR_TOKEN"

# 2. 按用户名筛选
curl -X GET "http://localhost:8080/admin/system/operation-logs?page=1&page_size=10&username=张三" \
  -H "Authorization: Bearer YOUR_TOKEN"

# 3. 按时间范围筛选
curl -X GET "http://localhost:8080/admin/system/operation-logs?page=1&page_size=10&start_time=1704067200&end_time=1704153600" \
  -H "Authorization: Bearer YOUR_TOKEN"

# 4. 获取统计数据
curl -X GET "http://localhost:8080/admin/system/operation-logs/stats?start_time=1704067200&end_time=1704153600" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

## 🎨 前端页面

### 访问地址

前端页面通过 Nova-admin 访问，需要配置路由。

**推荐路由配置：**

```typescript
// 在菜单配置中添加
{
  name: 'system',
  path: '/system',
  component: 'layout.base',
  meta: {
    title: '系统管理',
    icon: 'carbon:settings',
    order: 90,
  },
  children: [
    {
      name: 'system_operation-log',
      path: '/system/operation-log',
      component: 'view.system_operation-log',
      meta: {
        title: '操作日志',
        icon: 'carbon:document',
        requiresAuth: true,
      },
    },
    {
      name: 'system_operation-log-stats',
      path: '/system/operation-log/stats',
      component: 'view.system_operation-log-stats',
      meta: {
        title: '日志统计',
        icon: 'carbon:analytics',
        requiresAuth: true,
      },
    },
  ],
}
```

### 页面功能

#### 1. 操作日志列表页

**路径**：`/system/operation-log`

**功能**：
- ✅ 多条件筛选（用户名、资源、方法、操作类型、时间范围）
- ✅ 分页展示
- ✅ 查看详情
- ✅ 实时刷新

**截图示意**：
```
┌─────────────────────────────────────────────────┐
│ 操作日志                                         │
├─────────────────────────────────────────────────┤
│ 用户名 [____] 资源 [____] 方法 [▼] 时间 [____]  │
│ [搜索] [重置]                                    │
├─────────────────────────────────────────────────┤
│ 时间 | 用户 | 方法 | 路径 | 资源 | 操作 | 状态... │
│ 2024-01-01 10:00 | 张三 | POST | /admin/... ... │
│ 2024-01-01 09:55 | 李四 | GET  | /admin/... ... │
└─────────────────────────────────────────────────┘
```

#### 2. 操作日志统计页

**路径**：`/system/operation-log/stats`

**功能**：
- ✅ 总览统计（总操作数、成功率、失败率、平均耗时）
- ✅ TOP 10 操作用户
- ✅ TOP 10 操作类型
- ✅ 自定义时间范围

**统计指标**：
```
┌──────────────────────┬──────────────────────┐
│ 总操作数: 1,000      │ 成功操作: 950 (95%) │
├──────────────────────┼──────────────────────┤
│ 失败操作: 50 (5%)    │ 平均耗时: 123ms     │
└──────────────────────┴──────────────────────┘

TOP 10 操作用户:
1. 🥇 张三 - 150次
2. 🥈 李四 - 120次
3. 🥉 王五 - 100次

TOP 10 操作类型:
1. 🥇 create - 300次
2. 🥈 update - 250次
3. 🥉 delete - 150次
```

## 🔧 配置说明

### 修改端口

编辑 `config/system.yaml`：

```yaml
server:
  port: 8089  # 修改为其他端口

consul:
  service_port: 8089  # 同步修改
```

### 禁用 Consul

```yaml
consul:
  enabled: false  # 禁用 Consul 注册
```

### 修改数据库

```yaml
mongodb:
  host: "127.0.0.1"
  port: 27015
  database: "mule"  # 系统数据库
```

## 📊 数据存储

### 操作日志存储策略

操作日志采用**租户隔离**存储：

1. **系统管理员操作**：存储在 `system` 数据库
2. **租户用户操作**：存储在对应租户数据库

**存储路径**：
```
MongoDB
├── system (系统库)
│   └── operation_logs (系统管理员的日志)
├── ace (租户ace)
│   └── operation_logs (租户ace的日志)
└── test (租户test)
    └── operation_logs (租户test的日志)
```

### 自动清理

建议配置 TTL 索引自动清理旧日志：

```javascript
// MongoDB Shell
use system
db.operation_logs.createIndex(
  { "created_at": 1 },
  { expireAfterSeconds: 7776000 }  // 90天后自动删除
)
```

## 🔍 故障排查

### 1. 服务无法启动

```bash
# 检查端口占用
netstat -ano | findstr ":8089"

# 查看服务日志
./bin/system.exe 2>&1 | tee system.log
```

### 2. 无法查询日志

**可能原因**：
- 权限不足（需要登录）
- 租户上下文错误
- 数据库连接失败

**解决方法**：
```bash
# 检查 JWT Token
curl -H "Authorization: Bearer YOUR_TOKEN" \
  http://localhost:8080/auth/profile

# 检查数据库连接
mongosh --host 127.0.0.1 --port 27015
```

### 3. 统计数据不准确

**可能原因**：
- 时间范围设置错误
- 索引未创建

**解决方法**：
```javascript
// 创建索引
db.operation_logs.createIndex({ "user_id": 1 })
db.operation_logs.createIndex({ "created_at": -1 })
db.operation_logs.createIndex({ "user_id": 1, "created_at": -1 })
```

## 📝 开发说明

### 添加新的系统功能

1. **创建 DTO**：`app/system/dto/your_feature.go`
2. **创建 Service**：`app/system/services/your_feature.go`
3. **创建 Endpoint**：`app/system/endpoint/your_feature.go`
4. **创建 Transport**：`app/system/transport/your_feature.go`
5. **注册路由**：在 `cmd/system/main.go` 中注册

### 前端开发

1. **创建类型**：`frontend/src/typings/api/your-feature.d.ts`
2. **创建 API**：`frontend/src/service/api/your-feature.ts`
3. **创建页面**：`frontend/src/views/system/your-feature/index.vue`

## 🎓 最佳实践

### 1. 日志查询优化

- 使用索引字段进行筛选（user_id, created_at）
- 限制时间范围（避免全表扫描）
- 使用分页（避免一次加载过多数据）

### 2. 统计分析优化

- 选择合适的时间范围（推荐7天内）
- 使用缓存减少重复查询
- 定期清理过期数据

### 3. 安全建议

- 敏感数据脱敏（密码、Token等）
- 限制查询权限（普通用户只能查看自己的日志）
- 定期备份日志数据

## 📚 相关文档

- [System 服务 README](../app/system/README.md)
- [操作日志中间件使用指南](./操作日志中间件使用指南.md)
- [租户数据库隔离方案](./数据库级别租户隔离改造方案.md)

## 🆘 获取帮助

如有问题，请查看：
1. 服务日志：`./bin/system.exe` 的输出
2. MongoDB 日志：检查数据库连接
3. Consul 日志：检查服务注册状态
4. 网关日志：检查路由配置

---

**版本**：v1.0.0  
**更新时间**：2024-01-01  
**维护者**：Mule-Cloud 开发团队

