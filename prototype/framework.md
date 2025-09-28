# 信芙云服装生产管理系统架构框架

## 一、技术栈选型

### 1.1 后端技术栈

#### 核心框架
- **Go 1.21+**：主要编程语言
- **Gin 1.9+**：HTTP Web框架
- **GORM 1.25+**：ORM框架（用于关系型数据的处理）
- **MongoDB Driver 1.12+**：MongoDB官方驱动
- **Go-Redis 9.0+**：Redis客户端

#### 微服务框架
- **go-micro/go-kit**：微服务开发框架
- **Consul**：服务发现与配置管理
- **NATS**：消息队列和事件总线
- **OpenTelemetry**：分布式追踪
- **Prometheus**：监控指标收集

#### 安全认证
- **JWT**：身份认证
- **Casbin**：权限控制
- **bcrypt**：密码加密
- **AES-256**：敏感数据加密

### 1.2 数据存储

#### 主数据库
- **MongoDB 6.0+**：文档型数据库，主要业务数据
- **MongoDB Atlas**：云服务版本（生产环境推荐）

#### 缓存系统
- **Redis 7.0+**：内存缓存和消息队列
- **Redis Cluster**：高可用集群模式

#### 文件存储
- **MinIO**：对象存储服务
- **七牛云/阿里云OSS**：生产环境文件存储

### 1.3 DevOps技术栈

#### 容器化
- **Docker 20.10+**：容器化部署
- **Docker Compose**：本地开发环境
- **Kubernetes 1.25+**：生产环境容器编排

#### CI/CD
- **GitLab CI/GitHub Actions**：持续集成
- **Harbor**：容器镜像仓库
- **Helm**：Kubernetes应用包管理

#### 监控运维
- **Grafana**：监控面板
- **Prometheus**：指标收集
- **Jaeger**：分布式追踪
- **ELK Stack**：日志收集分析

### 1.4 前端技术栈

#### Web前端
- **Vue.js 3.x**：前端框架
- **TypeScript**：类型安全
- **Element Plus**：UI组件库
- **Pinia**：状态管理
- **Axios**：HTTP客户端

#### 移动端
- **React Native**：跨平台移动应用
- **Flutter**：备选方案

## 二、微服务架构设计

### 2.1 整体架构图

```
                                Internet
                                    |
                            ┌───────────────┐
                            │  Load Balancer │
                            │   (Nginx/F5)   │
                            └───────────────┘
                                    |
                    ┌───────────────────────────────────────┐
                    │              API Gateway               │
                    │  (Rate Limiting, Auth, Routing)       │
                    └───────────────────────────────────────┘
                                    |
        ┌───────────────────────────┼───────────────────────────┐
        │                          │                           │
┌───────▼─────────┐    ┌───────────▼────────────┐    ┌─────────▼──────────┐
│   User Service  │    │    Order Service       │    │ Production Service │
│   Port: 8001    │    │    Port: 8002          │    │   Port: 8003       │
└─────────────────┘    └────────────────────────┘    └────────────────────┘
        │                          │                           │
┌───────▼─────────┐    ┌───────────▼────────────┐    ┌─────────▼──────────┐
│Timesheet Service│    │   Payroll Service      │    │  Report Service    │
│   Port: 8004    │    │   Port: 8005           │    │   Port: 8006       │
└─────────────────┘    └────────────────────────┘    └────────────────────┘
        │                          │                           │
┌───────▼─────────┐    ┌───────────▼────────────┐    ┌─────────▼──────────┐
│Master Data Svc  │    │ Notification Service   │    │   File Service     │
│   Port: 8007    │    │   Port: 8008           │    │   Port: 8009       │
└─────────────────┘    └────────────────────────┘    └────────────────────┘
        │                          │                           │
        └──────────────────────────┼───────────────────────────┘
                                   │
                        ┌─────────────────────┐
                        │   Service Mesh      │
                        │ (Consul + NATS)     │
                        └─────────────────────┘
                                   │
        ┌──────────────────────────┼──────────────────────────┐
        │                          │                          │
┌───────▼────────┐    ┌────────────▼──────────┐    ┌─────────▼────────┐
│   MongoDB       │    │        Redis         │    │   File Storage   │
│  Primary: 27017 │    │   Master: 6379       │    │   MinIO: 9000    │
│  Replica: 27018 │    │   Slave: 6380        │    │                  │
└────────────────┘    └───────────────────────┘    └──────────────────┘
```

### 2.2 服务详细设计

#### User Service (用户服务)
```go
// 服务职责
- 用户注册、登录、注销
- JWT Token生成和验证
- 用户信息管理
- 角色权限管理
- 密码重置和修改

// 主要接口
POST /api/v1/users/register
POST /api/v1/users/login
POST /api/v1/users/logout
GET  /api/v1/users/profile
PUT  /api/v1/users/profile
POST /api/v1/users/change-password
GET  /api/v1/users/permissions

// 数据模型
type User struct {
    ID          primitive.ObjectID `bson:"_id,omitempty"`
    Username    string            `bson:"username"`
    Email       string            `bson:"email"`
    Password    string            `bson:"password"`
    RoleID      string            `bson:"role_id"`
    Status      int               `bson:"status"`
    CreatedAt   time.Time         `bson:"created_at"`
    UpdatedAt   time.Time         `bson:"updated_at"`
}
```

#### Order Service (订单服务)
```go
// 服务职责
- 款式管理
- 订单CRUD操作
- 订单状态流转
- 订单模板和复制
- 生产计划关联

// 主要接口
POST /api/v1/styles
GET  /api/v1/styles
PUT  /api/v1/styles/:id
DELETE /api/v1/styles/:id

POST /api/v1/orders
GET  /api/v1/orders
PUT  /api/v1/orders/:id
DELETE /api/v1/orders/:id
POST /api/v1/orders/:id/copy
PUT  /api/v1/orders/:id/status

// 数据模型
type Order struct {
    ID              primitive.ObjectID `bson:"_id,omitempty"`
    OrderNo         string            `bson:"order_no"`
    CustomerID      string            `bson:"customer_id"`
    StyleID         string            `bson:"style_id"`
    Quantity        int               `bson:"quantity"`
    Status          string            `bson:"status"`
    DeliveryDate    time.Time         `bson:"delivery_date"`
    Items           []OrderItem       `bson:"items"`
    CreatedAt       time.Time         `bson:"created_at"`
    UpdatedAt       time.Time         `bson:"updated_at"`
}
```

#### Production Service (生产服务)
```go
// 服务职责
- 生产计划制定
- 裁剪任务管理
- 生产进度跟踪
- 工序流转控制
- 质量检查记录

// 主要接口
POST /api/v1/production/plans
GET  /api/v1/production/plans
PUT  /api/v1/production/plans/:id

POST /api/v1/production/cutting-tasks
GET  /api/v1/production/cutting-tasks
PUT  /api/v1/production/cutting-tasks/:id/status

GET  /api/v1/production/progress/:order_id
PUT  /api/v1/production/progress

// 数据模型
type ProductionPlan struct {
    ID          primitive.ObjectID `bson:"_id,omitempty"`
    OrderID     string            `bson:"order_id"`
    ProcessFlow []ProcessStep     `bson:"process_flow"`
    StartDate   time.Time         `bson:"start_date"`
    EndDate     time.Time         `bson:"end_date"`
    Status      string            `bson:"status"`
    CreatedAt   time.Time         `bson:"created_at"`
}
```

### 2.3 服务间通信

#### 同步通信 (HTTP/gRPC)
```go
// HTTP REST API (外部接口)
type HTTPClient struct {
    baseURL string
    client  *http.Client
    timeout time.Duration
}

// gRPC (内部服务调用)
type GRPCClient struct {
    conn   *grpc.ClientConn
    client proto.ServiceClient
}

// 服务发现
type ServiceDiscovery struct {
    consul *api.Client
    cache  map[string]string
    mutex  sync.RWMutex
}
```

#### 异步通信 (NATS)
```go
// 事件发布
type EventPublisher struct {
    nats *nats.Conn
    js   nats.JetStreamContext
}

// 事件订阅
type EventSubscriber struct {
    nats *nats.Conn
    subs []*nats.Subscription
}

// 事件定义
type OrderCreatedEvent struct {
    OrderID   string    `json:"order_id"`
    Customer  string    `json:"customer"`
    Timestamp time.Time `json:"timestamp"`
}
```

## 三、数据架构设计

### 3.1 MongoDB 数据库设计

#### 数据库分片策略
```javascript
// 用户数据库 (mule_cloud_users)
db.users.createIndex({ "username": 1 }, { unique: true })
db.users.createIndex({ "email": 1 }, { unique: true })
db.roles.createIndex({ "name": 1 }, { unique: true })

// 业务数据库 (mule_cloud_business)
db.orders.createIndex({ "order_no": 1 }, { unique: true })
db.orders.createIndex({ "customer_id": 1, "created_at": -1 })
db.orders.createIndex({ "status": 1, "created_at": -1 })

// 生产数据库 (mule_cloud_production)
db.work_reports.createIndex({ "worker_id": 1, "date": -1 })
db.work_reports.createIndex({ "order_id": 1, "process_id": 1 })
db.production_progress.createIndex({ "order_id": 1, "updated_at": -1 })

// 工资数据库 (mule_cloud_payroll)
db.timesheets.createIndex({ "worker_id": 1, "date": -1 })
db.payrolls.createIndex({ "worker_id": 1, "pay_period": -1 })
```

#### 数据模型详细设计
```go
// 订单模型
type Order struct {
    ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    OrderNo      string            `bson:"order_no" json:"order_no"`
    CustomerID   primitive.ObjectID `bson:"customer_id" json:"customer_id"`
    StyleID      primitive.ObjectID `bson:"style_id" json:"style_id"`
    SalespersonID primitive.ObjectID `bson:"salesperson_id" json:"salesperson_id"`
    
    // 订单基本信息
    OrderType    string    `bson:"order_type" json:"order_type"`
    TotalQty     int       `bson:"total_qty" json:"total_qty"`
    UnitPrice    float64   `bson:"unit_price" json:"unit_price"`
    TotalAmount  float64   `bson:"total_amount" json:"total_amount"`
    Currency     string    `bson:"currency" json:"currency"`
    
    // 时间信息
    OrderDate    time.Time `bson:"order_date" json:"order_date"`
    DeliveryDate time.Time `bson:"delivery_date" json:"delivery_date"`
    
    // 状态信息
    Status       string    `bson:"status" json:"status"`
    Priority     int       `bson:"priority" json:"priority"`
    
    // 明细信息
    Items        []OrderItem `bson:"items" json:"items"`
    
    // 审计信息
    CreatedBy    primitive.ObjectID `bson:"created_by" json:"created_by"`
    UpdatedBy    primitive.ObjectID `bson:"updated_by,omitempty" json:"updated_by"`
    CreatedAt    time.Time         `bson:"created_at" json:"created_at"`
    UpdatedAt    time.Time         `bson:"updated_at" json:"updated_at"`
    Version      int               `bson:"version" json:"version"`
}

// 订单明细
type OrderItem struct {
    StyleID     primitive.ObjectID `bson:"style_id" json:"style_id"`
    ColorID     primitive.ObjectID `bson:"color_id" json:"color_id"`
    SizeID      primitive.ObjectID `bson:"size_id" json:"size_id"`
    Quantity    int               `bson:"quantity" json:"quantity"`
    UnitPrice   float64           `bson:"unit_price" json:"unit_price"`
    Amount      float64           `bson:"amount" json:"amount"`
    Remark      string            `bson:"remark,omitempty" json:"remark"`
}

// 生产进度
type ProductionProgress struct {
    ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    OrderID         primitive.ObjectID `bson:"order_id" json:"order_id"`
    ProcessID       primitive.ObjectID `bson:"process_id" json:"process_id"`
    
    // 进度信息
    PlannedQty      int               `bson:"planned_qty" json:"planned_qty"`
    CompletedQty    int               `bson:"completed_qty" json:"completed_qty"`
    DefectQty       int               `bson:"defect_qty" json:"defect_qty"`
    
    // 时间信息
    PlannedStart    time.Time         `bson:"planned_start" json:"planned_start"`
    PlannedEnd      time.Time         `bson:"planned_end" json:"planned_end"`
    ActualStart     *time.Time        `bson:"actual_start,omitempty" json:"actual_start"`
    ActualEnd       *time.Time        `bson:"actual_end,omitempty" json:"actual_end"`
    
    // 状态信息
    Status          string            `bson:"status" json:"status"`
    
    // 审计信息
    CreatedAt       time.Time         `bson:"created_at" json:"created_at"`
    UpdatedAt       time.Time         `bson:"updated_at" json:"updated_at"`
}

// 工时记录
type WorkReport struct {
    ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    WorkerID    primitive.ObjectID `bson:"worker_id" json:"worker_id"`
    OrderID     primitive.ObjectID `bson:"order_id" json:"order_id"`
    ProcessID   primitive.ObjectID `bson:"process_id" json:"process_id"`
    
    // 工作信息
    Date        time.Time          `bson:"date" json:"date"`
    StartTime   time.Time          `bson:"start_time" json:"start_time"`
    EndTime     time.Time          `bson:"end_time" json:"end_time"`
    WorkHours   float64            `bson:"work_hours" json:"work_hours"`
    Quantity    int                `bson:"quantity" json:"quantity"`
    
    // 工价信息
    UnitPrice   float64            `bson:"unit_price" json:"unit_price"`
    Amount      float64            `bson:"amount" json:"amount"`
    
    // 质量信息
    QualityGrade string            `bson:"quality_grade" json:"quality_grade"`
    DefectQty    int               `bson:"defect_qty" json:"defect_qty"`
    
    // 备注
    Remark      string             `bson:"remark,omitempty" json:"remark"`
    
    // 审计信息
    ReportedBy  primitive.ObjectID `bson:"reported_by" json:"reported_by"`
    CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
    UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
}
```

### 3.2 Redis 缓存架构

#### 缓存策略
```go
// 缓存键命名规范
const (
    // 用户相关
    UserCachePrefix     = "user:"
    UserSessionPrefix   = "session:"
    UserPermPrefix      = "perm:"
    
    // 业务数据
    OrderCachePrefix    = "order:"
    StyleCachePrefix    = "style:"
    CustomerCachePrefix = "customer:"
    
    // 统计数据
    StatsCachePrefix    = "stats:"
    ReportCachePrefix   = "report:"
    
    // 配置数据
    ConfigCachePrefix   = "config:"
    DictCachePrefix     = "dict:"
)

// 缓存管理器
type CacheManager struct {
    rdb *redis.Client
    ttl map[string]time.Duration
}

// 缓存操作接口
type Cache interface {
    Get(ctx context.Context, key string) (string, error)
    Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
    Del(ctx context.Context, keys ...string) error
    Exists(ctx context.Context, key string) (bool, error)
    HGet(ctx context.Context, key, field string) (string, error)
    HSet(ctx context.Context, key string, values ...interface{}) error
}
```

#### 分布式锁
```go
// 分布式锁实现
type DistributedLock struct {
    rdb    *redis.Client
    key    string
    value  string
    expiry time.Duration
}

func (l *DistributedLock) Lock(ctx context.Context) error {
    script := `
        if redis.call("GET", KEYS[1]) == ARGV[1] then
            return redis.call("PEXPIRE", KEYS[1], ARGV[2])
        else
            return redis.call("SET", KEYS[1], ARGV[1], "PX", ARGV[2], "NX")
        end
    `
    result := l.rdb.Eval(ctx, script, []string{l.key}, l.value, l.expiry.Milliseconds())
    return result.Err()
}
```

## 四、安全架构

### 4.1 认证授权

#### JWT 认证
```go
// JWT配置
type JWTConfig struct {
    SecretKey       string
    AccessTokenTTL  time.Duration
    RefreshTokenTTL time.Duration
    Issuer          string
}

// Token管理
type TokenManager struct {
    config   JWTConfig
    keyFunc  jwt.Keyfunc
}

// Claims结构
type CustomClaims struct {
    UserID      string   `json:"user_id"`
    Username    string   `json:"username"`
    Roles       []string `json:"roles"`
    Permissions []string `json:"permissions"`
    jwt.RegisteredClaims
}
```

#### RBAC权限模型
```go
// Casbin权限配置
type RBACConfig struct {
    ModelPath  string
    PolicyPath string
    Driver     string
    DSN        string
}

// 权限检查中间件
func AuthMiddleware(enforcer *casbin.Enforcer) gin.HandlerFunc {
    return func(c *gin.Context) {
        // 1. 验证JWT Token
        token := extractToken(c)
        claims, err := validateToken(token)
        if err != nil {
            c.JSON(401, gin.H{"error": "Invalid token"})
            c.Abort()
            return
        }
        
        // 2. 检查权限
        obj := c.Request.URL.Path
        act := c.Request.Method
        sub := claims.UserID
        
        allowed, err := enforcer.Enforce(sub, obj, act)
        if err != nil || !allowed {
            c.JSON(403, gin.H{"error": "Access denied"})
            c.Abort()
            return
        }
        
        c.Set("user", claims)
        c.Next()
    }
}
```

### 4.2 数据安全

#### 敏感数据加密
```go
// AES加密工具
type AESCrypto struct {
    key []byte
}

func (a *AESCrypto) Encrypt(plaintext string) (string, error) {
    block, err := aes.NewCipher(a.key)
    if err != nil {
        return "", err
    }
    
    // GCM模式加密
    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return "", err
    }
    
    nonce := make([]byte, gcm.NonceSize())
    if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
        return "", err
    }
    
    ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
    return base64.StdEncoding.EncodeToString(ciphertext), nil
}
```

## 五、部署架构

### 5.1 Docker化部署

#### Dockerfile示例
```dockerfile
# 多阶段构建
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/server

# 运行阶段
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/

COPY --from=builder /app/main .
COPY --from=builder /app/configs ./configs

EXPOSE 8080
CMD ["./main"]
```

#### Docker Compose配置
```yaml
version: '3.8'

services:
  # API网关
  gateway:
    build: ./gateway
    ports:
      - "8080:8080"
    environment:
      - CONSUL_ADDRESS=consul:8500
    depends_on:
      - consul
      - user-service
      - order-service

  # 用户服务
  user-service:
    build: ./services/user
    ports:
      - "8001:8001"
    environment:
      - MONGO_URI=mongodb://mongo:27017/mule_cloud_users
      - REDIS_URI=redis://redis:6379
      - CONSUL_ADDRESS=consul:8500
    depends_on:
      - mongo
      - redis
      - consul

  # 订单服务
  order-service:
    build: ./services/order
    ports:
      - "8002:8002"
    environment:
      - MONGO_URI=mongodb://mongo:27017/mule_cloud_business
      - REDIS_URI=redis://redis:6379
      - CONSUL_ADDRESS=consul:8500
    depends_on:
      - mongo
      - redis
      - consul

  # MongoDB
  mongo:
    image: mongo:6.0
    ports:
      - "27017:27017"
    environment:
      - MONGO_INITDB_ROOT_USERNAME=admin
      - MONGO_INITDB_ROOT_PASSWORD=password
    volumes:
      - mongo_data:/data/db

  # Redis
  redis:
    image: redis:7.0-alpine
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data

  # Consul
  consul:
    image: consul:1.15
    ports:
      - "8500:8500"
    command: consul agent -dev -client=0.0.0.0

volumes:
  mongo_data:
  redis_data:
```

### 5.2 Kubernetes部署

#### 服务部署配置
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: user-service
  namespace: mule-cloud
spec:
  replicas: 3
  selector:
    matchLabels:
      app: user-service
  template:
    metadata:
      labels:
        app: user-service
    spec:
      containers:
      - name: user-service
        image: mule-cloud/user-service:latest
        ports:
        - containerPort: 8001
        env:
        - name: MONGO_URI
          valueFrom:
            secretKeyRef:
              name: mongo-secret
              key: uri
        - name: REDIS_URI
          valueFrom:
            secretKeyRef:
              name: redis-secret
              key: uri
        resources:
          requests:
            memory: "128Mi"
            cpu: "100m"
          limits:
            memory: "512Mi"
            cpu: "500m"
        livenessProbe:
          httpGet:
            path: /health
            port: 8001
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /ready
            port: 8001
          initialDelaySeconds: 5
          periodSeconds: 5
---
apiVersion: v1
kind: Service
metadata:
  name: user-service
  namespace: mule-cloud
spec:
  selector:
    app: user-service
  ports:
  - port: 8001
    targetPort: 8001
  type: ClusterIP
```

## 六、监控与日志

### 6.1 Prometheus监控

#### 指标定义
```go
// 业务指标
var (
    // HTTP请求指标
    httpRequestsTotal = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
            Help: "Total number of HTTP requests",
        },
        []string{"method", "endpoint", "status"},
    )
    
    // 响应时间指标
    httpRequestDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name: "http_request_duration_seconds",
            Help: "HTTP request duration in seconds",
            Buckets: prometheus.DefBuckets,
        },
        []string{"method", "endpoint"},
    )
    
    // 业务指标
    ordersCreatedTotal = prometheus.NewCounter(
        prometheus.CounterOpts{
            Name: "orders_created_total",
            Help: "Total number of orders created",
        },
    )
    
    activeWorkersGauge = prometheus.NewGauge(
        prometheus.GaugeOpts{
            Name: "active_workers",
            Help: "Number of active workers",
        },
    )
)
```

### 6.2 日志系统

#### 结构化日志
```go
// 日志配置
type LogConfig struct {
    Level      string `yaml:"level"`
    Format     string `yaml:"format"`
    Output     string `yaml:"output"`
    MaxSize    int    `yaml:"max_size"`
    MaxBackups int    `yaml:"max_backups"`
    MaxAge     int    `yaml:"max_age"`
}

// 日志记录器
type Logger struct {
    logger *logrus.Logger
    fields logrus.Fields
}

// 操作日志
type OperationLog struct {
    UserID      string      `json:"user_id"`
    Action      string      `json:"action"`
    Resource    string      `json:"resource"`
    ResourceID  string      `json:"resource_id"`
    OldData     interface{} `json:"old_data,omitempty"`
    NewData     interface{} `json:"new_data,omitempty"`
    IP          string      `json:"ip"`
    UserAgent   string      `json:"user_agent"`
    Timestamp   time.Time   `json:"timestamp"`
    Duration    int64       `json:"duration_ms"`
    Error       string      `json:"error,omitempty"`
}
```

## 七、开发规范

### 7.1 代码规范

#### 项目结构
```
mule-cloud/
├── cmd/                    # 主程序入口
│   ├── gateway/           # API网关
│   ├── user-service/      # 用户服务
│   └── order-service/     # 订单服务
├── internal/               # 私有代码
│   ├── config/            # 配置管理
│   ├── middleware/        # 中间件
│   ├── models/            # 数据模型
│   ├── repository/        # 数据访问层
│   ├── service/           # 业务逻辑层
│   └── handler/           # 处理器层
├── pkg/                   # 公共代码库
│   ├── auth/              # 认证工具
│   ├── cache/             # 缓存工具
│   ├── database/          # 数据库工具
│   ├── logger/            # 日志工具
│   └── utils/             # 通用工具
├── api/                   # API定义
│   ├── proto/             # gRPC协议
│   └── swagger/           # API文档
├── configs/               # 配置文件
├── scripts/               # 脚本文件
├── deployments/           # 部署文件
└── docs/                  # 文档
```

#### 编码规范
```go
// 接口定义
type UserService interface {
    CreateUser(ctx context.Context, user *User) error
    GetUser(ctx context.Context, id string) (*User, error)
    UpdateUser(ctx context.Context, user *User) error
    DeleteUser(ctx context.Context, id string) error
}

// 实现结构
type userService struct {
    repo   UserRepository
    cache  Cache
    logger Logger
}

// 构造函数
func NewUserService(repo UserRepository, cache Cache, logger Logger) UserService {
    return &userService{
        repo:   repo,
        cache:  cache,
        logger: logger,
    }
}

// 错误处理
func (s *userService) CreateUser(ctx context.Context, user *User) error {
    // 参数验证
    if err := validateUser(user); err != nil {
        return fmt.Errorf("invalid user data: %w", err)
    }
    
    // 业务逻辑
    if err := s.repo.Create(ctx, user); err != nil {
        s.logger.Error("failed to create user", "error", err, "user_id", user.ID)
        return fmt.Errorf("create user failed: %w", err)
    }
    
    // 缓存更新
    if err := s.cache.Del(ctx, fmt.Sprintf("user:%s", user.ID)); err != nil {
        s.logger.Warn("failed to invalidate cache", "error", err)
    }
    
    return nil
}
```

### 7.2 API设计规范

#### RESTful API设计
```go
// 路由定义
func SetupRoutes(r *gin.Engine, userHandler *UserHandler) {
    v1 := r.Group("/api/v1")
    {
        // 用户相关路由
        users := v1.Group("/users")
        users.Use(AuthMiddleware())
        {
            users.POST("", userHandler.CreateUser)
            users.GET("", userHandler.ListUsers)
            users.GET("/:id", userHandler.GetUser)
            users.PUT("/:id", userHandler.UpdateUser)
            users.DELETE("/:id", userHandler.DeleteUser)
        }
        
        // 订单相关路由
        orders := v1.Group("/orders")
        orders.Use(AuthMiddleware())
        {
            orders.POST("", orderHandler.CreateOrder)
            orders.GET("", orderHandler.ListOrders)
            orders.GET("/:id", orderHandler.GetOrder)
            orders.PUT("/:id", orderHandler.UpdateOrder)
            orders.DELETE("/:id", orderHandler.DeleteOrder)
            orders.POST("/:id/copy", orderHandler.CopyOrder)
            orders.PUT("/:id/status", orderHandler.UpdateStatus)
        }
    }
}

// 响应格式统一
type ApiResponse struct {
    Code    int         `json:"code"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
    Meta    *Meta       `json:"meta,omitempty"`
}

type Meta struct {
    Page       int `json:"page"`
    PageSize   int `json:"page_size"`
    Total      int `json:"total"`
    TotalPages int `json:"total_pages"`
}
```

## 八、性能优化

### 8.1 数据库优化

#### 查询优化
```go
// 索引创建
func CreateIndexes(db *mongo.Database) error {
    collections := map[string][]mongo.IndexModel{
        "orders": {
            {Keys: bson.D{{"order_no", 1}}, Options: options.Index().SetUnique(true)},
            {Keys: bson.D{{"customer_id", 1}, {"created_at", -1}}},
            {Keys: bson.D{{"status", 1}, {"created_at", -1}}},
        },
        "work_reports": {
            {Keys: bson.D{{"worker_id", 1}, {"date", -1}}},
            {Keys: bson.D{{"order_id", 1}, {"process_id", 1}}},
        },
    }
    
    for collName, indexes := range collections {
        coll := db.Collection(collName)
        _, err := coll.Indexes().CreateMany(context.Background(), indexes)
        if err != nil {
            return fmt.Errorf("failed to create indexes for %s: %w", collName, err)
        }
    }
    
    return nil
}
```

### 8.2 缓存优化

#### 多级缓存
```go
// L1缓存：本地内存缓存
type L1Cache struct {
    cache map[string]interface{}
    mutex sync.RWMutex
    ttl   time.Duration
}

// L2缓存：Redis分布式缓存
type L2Cache struct {
    rdb *redis.Client
}

// 多级缓存管理器
type MultiLevelCache struct {
    l1 *L1Cache
    l2 *L2Cache
}

func (c *MultiLevelCache) Get(key string) (interface{}, error) {
    // 先查L1缓存
    if value, ok := c.l1.Get(key); ok {
        return value, nil
    }
    
    // 再查L2缓存
    value, err := c.l2.Get(key)
    if err != nil {
        return nil, err
    }
    
    // 回写L1缓存
    c.l1.Set(key, value)
    
    return value, nil
}
```

## 九、质量保证

### 9.1 单元测试

#### 测试框架
```go
// 使用testify框架
func TestUserService_CreateUser(t *testing.T) {
    // 准备测试数据
    mockRepo := &MockUserRepository{}
    mockCache := &MockCache{}
    mockLogger := &MockLogger{}
    
    service := NewUserService(mockRepo, mockCache, mockLogger)
    
    user := &User{
        Username: "testuser",
        Email:    "test@example.com",
        Password: "password123",
    }
    
    // 设置Mock期望
    mockRepo.On("Create", mock.Anything, user).Return(nil)
    mockCache.On("Del", mock.Anything, mock.Anything).Return(nil)
    
    // 执行测试
    err := service.CreateUser(context.Background(), user)
    
    // 断言结果
    assert.NoError(t, err)
    mockRepo.AssertExpectations(t)
    mockCache.AssertExpectations(t)
}
```

### 9.2 集成测试

#### 测试容器
```go
// 使用testcontainers进行集成测试
func TestOrderIntegration(t *testing.T) {
    // 启动测试用的MongoDB容器
    mongoContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
        ContainerRequest: testcontainers.ContainerRequest{
            Image:        "mongo:6.0",
            ExposedPorts: []string{"27017/tcp"},
            Env: map[string]string{
                "MONGO_INITDB_ROOT_USERNAME": "test",
                "MONGO_INITDB_ROOT_PASSWORD": "test",
            },
            WaitingFor: wait.ForLog("waiting for connections on port 27017"),
        },
        Started: true,
    })
    require.NoError(t, err)
    defer mongoContainer.Terminate(ctx)
    
    // 获取容器连接信息
    mongoHost, err := mongoContainer.Host(ctx)
    require.NoError(t, err)
    
    mongoPort, err := mongoContainer.MappedPort(ctx, "27017")
    require.NoError(t, err)
    
    // 连接数据库并执行测试
    mongoURI := fmt.Sprintf("mongodb://test:test@%s:%s", mongoHost, mongoPort.Port())
    // ... 测试逻辑
}
```

---

*本架构框架文档提供了系统实现的详细技术指导，可作为开发团队的技术参考手册*
