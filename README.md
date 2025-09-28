# 信芙云服装生产管理系统

## 项目简介

信芙云是一个基于Go微服务架构的中小型服装生产管理系统，旨在帮助服装制造企业实现数字化生产管理。

**🎉 当前状态**：基础架构+基础数据服务已完成！用户认证、基础数据管理、API网关、微服务通信等核心组件已实现并测试通过。

### 核心功能

- **订单管理**：款式管理、订单创建、订单跟踪
- **生产管理**：生产计划、裁剪任务、进度监控
- **工时管理**：工作上报、工时统计、进度跟踪
- **工资管理**：工资计算、工资统计、工资发放
- **基础数据**：客户管理、业务员管理、工序管理
- **用户管理**：用户认证、角色权限、用户管理

### 技术栈

- **后端**：Go 1.21+, Gin, MongoDB, Redis
- **架构**：微服务架构
- **服务发现**：Consul
- **消息队列**：NATS
- **监控**：Prometheus + Grafana
- **容器化**：Docker + Docker Compose

## 快速开始

### 前置要求

- Docker 20.10+
- Docker Compose 3.8+
- Go 1.21+ (如需本地开发)

### 使用Docker Compose启动

1. **克隆项目**
```bash
git clone <repository-url>
cd mule-cloud
```

2. **启动所有服务**
```bash
docker-compose up -d
```

3. **查看服务状态**
```bash
docker-compose ps
```

4. **查看日志**
```bash
# 查看所有服务日志
docker-compose logs -f

# 查看特定服务日志
docker-compose logs -f user-service
```

### 服务访问地址

- **API网关**: http://localhost:8080
- **用户服务**: http://localhost:8001
- **Consul UI**: http://localhost:8500
- **Prometheus**: http://localhost:9090
- **Grafana**: http://localhost:3000
- **MongoDB**: localhost:27017
- **Redis**: localhost:6379
- **NATS**: localhost:4222

### 默认账号

- **系统管理员**: `admin` / `password`
- **Grafana**: `admin` / `admin123`
- **MongoDB**: `admin` / `password123`
- **Redis**: 密码 `redis123`

## API文档

### 认证相关

#### 用户注册
```bash
curl -X POST http://localhost:8001/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com", 
    "password": "password123",
    "real_name": "测试用户"
  }'
```

#### 用户登录
```bash
curl -X POST http://localhost:8001/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "password"
  }'
```

#### 获取用户资料
```bash
curl -X GET http://localhost:8001/api/v1/users/profile \
  -H "Authorization: Bearer <access_token>"
```

### 健康检查

```bash
curl http://localhost:8001/health
```

## 开发指南

### 本地开发环境

1. **安装依赖**
```bash
go mod download
```

2. **启动基础设施**
```bash
# 只启动数据库等基础服务
docker-compose up -d mongodb redis consul nats
```

3. **运行用户服务**
```bash
cd cmd/user-service
go run main.go
```

### 项目结构

```
mule-cloud/
├── cmd/                    # 主程序入口
│   ├── user-service/      # 用户服务
│   ├── order-service/     # 订单服务
│   └── gateway/           # API网关
├── internal/               # 私有代码
│   ├── handler/           # HTTP处理器
│   ├── service/           # 业务逻辑
│   ├── repository/        # 数据访问
│   ├── models/            # 数据模型
│   └── middleware/        # 中间件
├── pkg/                   # 公共库
│   ├── auth/              # 认证工具
│   ├── cache/             # 缓存工具
│   ├── config/            # 配置管理
│   ├── database/          # 数据库工具
│   └── logger/            # 日志工具
├── configs/               # 配置文件
├── deployments/           # 部署文件
├── scripts/               # 脚本文件
├── prototype/             # 原型设计
└── docs/                  # 文档
```

### 代码规范

- 遵循Go官方代码规范
- 使用`gofmt`格式化代码
- 所有公共函数和结构体需要添加注释
- 错误处理使用包装错误的方式
- 数据库操作必须使用事务
- API接口需要参数验证

## 微服务架构

### 服务列表

1. **用户服务** (user-service:8001)
   - 用户注册、登录、认证
   - 用户信息管理
   - 角色权限管理

2. **订单服务** (order-service:8002)
   - 订单管理
   - 款式管理
   - 客户管理

3. **生产服务** (production-service:8003)
   - 生产计划
   - 裁剪任务
   - 生产进度

4. **工时服务** (timesheet-service:8004)
   - 工作上报
   - 工时统计
   - 进度跟踪

5. **工资服务** (payroll-service:8005)
   - 工资计算
   - 工资统计
   - 工资发放

6. **报表服务** (report-service:8006)
   - 数据统计
   - 报表生成
   - 监控面板

7. **基础数据服务** (master-data-service:8007)
   - 工序管理
   - 尺码颜色管理
   - 字典数据

8. **通知服务** (notification-service:8008)
   - 消息推送
   - 邮件通知
   - 系统公告

9. **文件服务** (file-service:8009)
   - 文件上传
   - 图片处理
   - 文档管理

10. **API网关** (gateway:8080)
    - 路由转发
    - 负载均衡
    - 限流熔断

### 服务间通信

- **同步通信**: HTTP/REST API (外部调用)
- **异步通信**: NATS消息队列 (内部事件)
- **服务发现**: Consul注册中心
- **配置管理**: Consul KV存储

## 部署指南

### Docker部署

1. **构建镜像**
```bash
# 构建用户服务
docker build -f deployments/user-service/Dockerfile -t mule-cloud/user-service:latest .

# 构建API网关
docker build -f deployments/gateway/Dockerfile -t mule-cloud/gateway:latest .
```

2. **推送镜像**
```bash
docker tag mule-cloud/user-service:latest registry.example.com/mule-cloud/user-service:latest
docker push registry.example.com/mule-cloud/user-service:latest
```

### Kubernetes部署

详见 `deployments/k8s/` 目录下的配置文件。

### 生产环境配置

1. **安全配置**
   - 修改默认密码
   - 配置HTTPS证书
   - 设置防火墙规则

2. **性能优化**
   - 调整数据库连接池
   - 配置Redis集群
   - 启用Gzip压缩

3. **监控告警**
   - 配置Prometheus监控
   - 设置Grafana告警
   - 配置日志收集

## 监控运维

### 健康检查

所有服务都提供健康检查端点：
```bash
curl http://service-host:port/health
```

### 性能监控

- **Prometheus指标**: http://localhost:9090
- **Grafana面板**: http://localhost:3000
- **服务监控**: http://localhost:8500

### 日志管理

日志输出格式为JSON，包含以下字段：
- `timestamp`: 时间戳
- `level`: 日志级别
- `message`: 日志消息
- `service`: 服务名称
- `request_id`: 请求ID
- `user_id`: 用户ID

## 常见问题

### Q: 如何重置数据库？

```bash
# 停止服务
docker-compose down

# 删除数据卷
docker volume rm mule-cloud_mongodb_data

# 重新启动
docker-compose up -d
```

### Q: 如何修改服务端口？

修改 `docker-compose.yaml` 文件中对应服务的端口映射。

### Q: 如何查看服务注册状态？

访问Consul UI: http://localhost:8500

### Q: 如何扩展服务？

```bash
# 扩展用户服务到3个实例
docker-compose up -d --scale user-service=3
```

## 贡献指南

1. Fork项目
2. 创建功能分支
3. 提交代码
4. 创建Pull Request

## 许可证

本项目采用MIT许可证，详见LICENSE文件。

## 联系我们

- 项目主页：https://github.com/mule-cloud/mule-cloud
- 问题反馈：https://github.com/mule-cloud/mule-cloud/issues
- 邮箱：support@mulecloud.com
