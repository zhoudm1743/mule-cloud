# MongoDB、Redis、Logger 使用指南

## 📖 概述

本项目已集成：
- ✅ **MongoDB**: 非关系型数据库，用于数据持久化
- ✅ **Redis**: 键值存储，用于缓存和会话管理
- ✅ **Zap Logger**: 高性能结构化日志系统

---

## ⚙️ 配置说明

### MongoDB 配置

在 `config.yaml` 中配置：

```yaml
mongodb:
  enabled: true  # 是否启用MongoDB
  # 方式1: 使用URI连接（推荐）
  uri: "mongodb://localhost:27017"
  # 方式2: 使用独立配置
  host: "127.0.0.1"
  port: 27017
  username: ""
  password: ""
  database: "mule_cloud"
  auth_source: "admin"      # 认证数据库
  replica_set: ""           # 副本集名称（可选）
  max_pool_size: 100        # 最大连接池大小
  min_pool_size: 10         # 最小连接池大小
  timeout: 10               # 连接超时（秒）
```

### Redis 配置

```yaml
redis:
  enabled: true
  host: "127.0.0.1"
  port: 6379
  password: ""
  db: 0
  pool_size: 10
```

### 日志配置

```yaml
log:
  level: "info"                    # debug, info, warn, error
  format: "text"                   # json, text
  output: "stdout"                 # stdout, file
  file_path: "logs/app.log"        # 日志文件路径
  max_size: 100                    # 单个文件最大MB
  max_backups: 3                   # 保留文件数量
  max_age: 7                       # 保留天数
  compress: true                   # 是否压缩旧文件
```

---

## 🚀 MongoDB 使用

### 1. 初始化

在 `main.go` 中初始化（已集成）：

```go
import (
    dbPkg "mule-cloud/core/database"
    cfgPkg "mule-cloud/core/config"
)

func main() {
    cfg, _ := cfgPkg.Load("config.yaml")
    
    // 初始化MongoDB
    if cfg.MongoDB.Enabled {
        if _, err := dbPkg.InitMongoDB(&cfg.MongoDB); err != nil {
            log.Fatal(err)
        }
        defer dbPkg.CloseMongoDB()
    }
}
```

### 2. 获取数据库实例

```go
import dbPkg "mule-cloud/core/database"

// 方式1: 获取数据库实例
db := dbPkg.GetMongoDB()

// 方式2: 获取集合
collection := dbPkg.GetCollection("users")
```

### 3. 基本CRUD操作

#### 插入文档

```go
package services

import (
    "context"
    "time"
    dbPkg "mule-cloud/core/database"
    "go.mongodb.org/mongo-driver/bson"
)

type User struct {
    ID        string    `bson:"_id,omitempty"`
    Name      string    `bson:"name"`
    Email     string    `bson:"email"`
    CreatedAt time.Time `bson:"created_at"`
}

// 插入单个文档
func CreateUser(user *User) error {
    collection := dbPkg.GetCollection("users")
    user.CreatedAt = time.Now()
    
    _, err := collection.InsertOne(context.Background(), user)
    return err
}

// 插入多个文档
func CreateUsers(users []User) error {
    collection := dbPkg.GetCollection("users")
    
    docs := make([]interface{}, len(users))
    for i, u := range users {
        u.CreatedAt = time.Now()
        docs[i] = u
    }
    
    _, err := collection.InsertMany(context.Background(), docs)
    return err
}
```

#### 查询文档

```go
// 查询单个文档
func GetUserByID(id string) (*User, error) {
    collection := dbPkg.GetCollection("users")
    
    var user User
    filter := bson.M{"_id": id}
    err := collection.FindOne(context.Background(), filter).Decode(&user)
    
    return &user, err
}

// 查询多个文档
func GetUsersByEmail(email string) ([]User, error) {
    collection := dbPkg.GetCollection("users")
    
    filter := bson.M{"email": email}
    cursor, err := collection.Find(context.Background(), filter)
    if err != nil {
        return nil, err
    }
    defer cursor.Close(context.Background())
    
    var users []User
    if err = cursor.All(context.Background(), &users); err != nil {
        return nil, err
    }
    
    return users, nil
}

// 查询所有文档（分页）
func GetAllUsers(page, pageSize int) ([]User, error) {
    collection := dbPkg.GetCollection("users")
    
    skip := (page - 1) * pageSize
    opts := options.Find().
        SetSkip(int64(skip)).
        SetLimit(int64(pageSize)).
        SetSort(bson.D{{"created_at", -1}})
    
    cursor, err := collection.Find(context.Background(), bson.M{}, opts)
    if err != nil {
        return nil, err
    }
    defer cursor.Close(context.Background())
    
    var users []User
    if err = cursor.All(context.Background(), &users); err != nil {
        return nil, err
    }
    
    return users, nil
}
```

#### 更新文档

```go
// 更新单个文档
func UpdateUser(id string, name, email string) error {
    collection := dbPkg.GetCollection("users")
    
    filter := bson.M{"_id": id}
    update := bson.M{
        "$set": bson.M{
            "name":  name,
            "email": email,
            "updated_at": time.Now(),
        },
    }
    
    _, err := collection.UpdateOne(context.Background(), filter, update)
    return err
}

// 更新多个文档
func UpdateUsersByRole(oldRole, newRole string) error {
    collection := dbPkg.GetCollection("users")
    
    filter := bson.M{"role": oldRole}
    update := bson.M{"$set": bson.M{"role": newRole}}
    
    _, err := collection.UpdateMany(context.Background(), filter, update)
    return err
}
```

#### 删除文档

```go
// 删除单个文档
func DeleteUser(id string) error {
    collection := dbPkg.GetCollection("users")
    
    filter := bson.M{"_id": id}
    _, err := collection.DeleteOne(context.Background(), filter)
    return err
}

// 删除多个文档
func DeleteInactiveUsers(days int) error {
    collection := dbPkg.GetCollection("users")
    
    cutoffDate := time.Now().AddDate(0, 0, -days)
    filter := bson.M{"last_login": bson.M{"$lt": cutoffDate}}
    
    _, err := collection.DeleteMany(context.Background(), filter)
    return err
}
```

---

## 🔥 Redis 使用

### 1. 初始化

在 `main.go` 中初始化（已集成）：

```go
import cachePkg "mule-cloud/core/cache"

func main() {
    // 初始化Redis
    if cfg.Redis.Enabled {
        if _, err := cachePkg.InitRedis(&cfg.Redis); err != nil {
            log.Fatal(err)
        }
        defer cachePkg.CloseRedis()
    }
}
```

### 2. 基本操作

#### 字符串操作

```go
import (
    "context"
    "time"
    cachePkg "mule-cloud/core/cache"
)

// 设置键值（带过期时间）
func CacheUserToken(userID, token string) error {
    ctx := context.Background()
    key := fmt.Sprintf("user:token:%s", userID)
    return cachePkg.Set(ctx, key, token, 24*time.Hour)
}

// 获取键值
func GetUserToken(userID string) (string, error) {
    ctx := context.Background()
    key := fmt.Sprintf("user:token:%s", userID)
    return cachePkg.Get(ctx, key)
}

// 删除键
func DeleteUserToken(userID string) error {
    ctx := context.Background()
    key := fmt.Sprintf("user:token:%s", userID)
    return cachePkg.Del(ctx, key)
}
```

#### 哈希操作

```go
// 缓存用户信息
func CacheUserInfo(userID string, user map[string]interface{}) error {
    ctx := context.Background()
    key := fmt.Sprintf("user:info:%s", userID)
    
    values := make([]interface{}, 0, len(user)*2)
    for k, v := range user {
        values = append(values, k, v)
    }
    
    return cachePkg.HSet(ctx, key, values...)
}

// 获取用户信息
func GetUserInfo(userID string) (map[string]string, error) {
    ctx := context.Background()
    key := fmt.Sprintf("user:info:%s", userID)
    return cachePkg.HGetAll(ctx, key)
}
```

#### 计数器

```go
// 增加访问计数
func IncrementViewCount(articleID string) (int64, error) {
    ctx := context.Background()
    key := fmt.Sprintf("article:views:%s", articleID)
    return cachePkg.Incr(ctx, key)
}

// 获取访问计数
func GetViewCount(articleID string) (string, error) {
    ctx := context.Background()
    key := fmt.Sprintf("article:views:%s", articleID)
    return cachePkg.Get(ctx, key)
}
```

#### 列表操作

```go
// 添加到最近浏览列表
func AddToRecentViews(userID, articleID string) error {
    ctx := context.Background()
    key := fmt.Sprintf("user:recent:%s", userID)
    
    // 左侧推入，保持最新在前
    if err := cachePkg.LPush(ctx, key, articleID); err != nil {
        return err
    }
    
    // 只保留最近10条
    client := cachePkg.GetRedis()
    return client.LTrim(ctx, key, 0, 9).Err()
}
```

### 3. 高级用法

#### 分布式锁

```go
import "github.com/redis/go-redis/v9"

func AcquireLock(key string, expiration time.Duration) (bool, error) {
    ctx := context.Background()
    client := cachePkg.GetRedis()
    
    result, err := client.SetNX(ctx, key, "locked", expiration).Result()
    return result, err
}

func ReleaseLock(key string) error {
    ctx := context.Background()
    return cachePkg.Del(ctx, key)
}
```

---

## 📝 Logger 使用

### 1. 初始化

在 `main.go` 中初始化（已集成）：

```go
import loggerPkg "mule-cloud/core/logger"

func main() {
    // 初始化日志系统
    if err := loggerPkg.InitLogger(&cfg.Log); err != nil {
        log.Fatal(err)
    }
    defer loggerPkg.Close()
}
```

### 2. 基本日志

```go
import (
    loggerPkg "mule-cloud/core/logger"
    "go.uber.org/zap"
)

// 简单日志
loggerPkg.Info("服务启动成功")
loggerPkg.Warn("配置文件未找到，使用默认配置")
loggerPkg.Error("连接数据库失败")

// 带字段的结构化日志
loggerPkg.Info("用户登录",
    zap.String("user_id", "123"),
    zap.String("username", "admin"),
    zap.String("ip", "192.168.1.1"),
)

// 格式化日志
loggerPkg.Infof("用户 %s 登录成功", username)
loggerPkg.Errorf("连接失败: %v", err)
```

### 3. 在服务中使用

```go
package services

import (
    loggerPkg "mule-cloud/core/logger"
    "go.uber.org/zap"
)

func CreateUser(user *User) error {
    loggerPkg.Info("创建用户",
        zap.String("name", user.Name),
        zap.String("email", user.Email),
    )
    
    // 业务逻辑
    if err := saveToDatabase(user); err != nil {
        loggerPkg.Error("保存用户失败",
            zap.Error(err),
            zap.String("user_id", user.ID),
        )
        return err
    }
    
    loggerPkg.Info("用户创建成功",
        zap.String("user_id", user.ID),
    )
    return nil
}
```

### 4. 带上下文的日志

```go
// 创建带特定字段的logger
logger := loggerPkg.With(
    zap.String("request_id", requestID),
    zap.String("user_id", userID),
)

// 使用该logger记录日志
logger.Info("处理请求")
logger.Warn("请求参数不完整")
logger.Error("处理失败", zap.Error(err))
```

---

## 🧪 测试示例

### 启动服务（已配置）

```bash
# 启动 Basic 服务
cd basic/cmd
go run main.go -config=../../config/basic.yaml

# 启动 Test 服务
cd test/cmd
go run main.go -config=../../config/test.yaml

# 启动网关
cd gateway
go run main.go -config=config/gateway.yaml
```

### 日志输出示例

```
2025-09-30T10:00:00.000+0800    INFO    logger/logger.go:71    ✅ 日志系统初始化成功    {"level": "info", "format": "text", "output": "stdout"}
2025-09-30T10:00:01.000+0800    INFO    cmd/main.go:38    🚀 BasicService 启动中...    {"service": "basicservice", "port": 8001}
2025-09-30T10:00:02.000+0800    INFO    database/mongodb.go:67    ✅ MongoDB连接成功: 127.0.0.1:27017/mule_cloud
2025-09-30T10:00:03.000+0800    INFO    cache/redis.go:47    ✅ Redis连接成功: 127.0.0.1:6379 (DB:0)
```

---

## 📚 完整示例：用户管理服务

创建 `services/user_service.go`：

```go
package services

import (
    "context"
    "fmt"
    "time"
    
    cachePkg "mule-cloud/core/cache"
    dbPkg "mule-cloud/core/database"
    loggerPkg "mule-cloud/core/logger"
    
    "go.mongodb.org/mongo-driver/bson"
    "go.uber.org/zap"
)

type User struct {
    ID        string    `bson:"_id,omitempty" json:"id"`
    Name      string    `bson:"name" json:"name"`
    Email     string    `bson:"email" json:"email"`
    Role      string    `bson:"role" json:"role"`
    CreatedAt time.Time `bson:"created_at" json:"created_at"`
    UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

type UserService struct{}

func NewUserService() *UserService {
    return &UserService{}
}

// CreateUser 创建用户
func (s *UserService) CreateUser(ctx context.Context, user *User) error {
    loggerPkg.Info("创建用户", zap.String("email", user.Email))
    
    // 设置时间
    now := time.Now()
    user.CreatedAt = now
    user.UpdatedAt = now
    
    // 保存到MongoDB
    collection := dbPkg.GetCollection("users")
    result, err := collection.InsertOne(ctx, user)
    if err != nil {
        loggerPkg.Error("创建用户失败", zap.Error(err))
        return err
    }
    
    user.ID = result.InsertedID.(string)
    
    // 缓存用户信息到Redis
    cacheKey := fmt.Sprintf("user:%s", user.ID)
    if err := cachePkg.HSet(ctx, cacheKey,
        "name", user.Name,
        "email", user.Email,
        "role", user.Role,
    ); err != nil {
        loggerPkg.Warn("缓存用户信息失败", zap.Error(err))
    }
    
    // 设置缓存过期时间
    cachePkg.Expire(ctx, cacheKey, 24*time.Hour)
    
    loggerPkg.Info("用户创建成功", zap.String("user_id", user.ID))
    return nil
}

// GetUser 获取用户（优先从缓存读取）
func (s *UserService) GetUser(ctx context.Context, id string) (*User, error) {
    loggerPkg.Debug("获取用户", zap.String("user_id", id))
    
    // 先从Redis缓存获取
    cacheKey := fmt.Sprintf("user:%s", id)
    userInfo, err := cachePkg.HGetAll(ctx, cacheKey)
    
    if err == nil && len(userInfo) > 0 {
        loggerPkg.Debug("从缓存获取用户", zap.String("user_id", id))
        return &User{
            ID:    id,
            Name:  userInfo["name"],
            Email: userInfo["email"],
            Role:  userInfo["role"],
        }, nil
    }
    
    // 缓存未命中，从MongoDB获取
    loggerPkg.Debug("缓存未命中，从数据库获取", zap.String("user_id", id))
    collection := dbPkg.GetCollection("users")
    
    var user User
    filter := bson.M{"_id": id}
    if err := collection.FindOne(ctx, filter).Decode(&user); err != nil {
        loggerPkg.Error("获取用户失败", zap.Error(err), zap.String("user_id", id))
        return nil, err
    }
    
    // 写入缓存
    cachePkg.HSet(ctx, cacheKey,
        "name", user.Name,
        "email", user.Email,
        "role", user.Role,
    )
    cachePkg.Expire(ctx, cacheKey, 24*time.Hour)
    
    return &user, nil
}

// DeleteUser 删除用户
func (s *UserService) DeleteUser(ctx context.Context, id string) error {
    loggerPkg.Info("删除用户", zap.String("user_id", id))
    
    // 从MongoDB删除
    collection := dbPkg.GetCollection("users")
    filter := bson.M{"_id": id}
    if _, err := collection.DeleteOne(ctx, filter); err != nil {
        loggerPkg.Error("删除用户失败", zap.Error(err))
        return err
    }
    
    // 删除Redis缓存
    cacheKey := fmt.Sprintf("user:%s", id)
    if err := cachePkg.Del(ctx, cacheKey); err != nil {
        loggerPkg.Warn("删除缓存失败", zap.Error(err))
    }
    
    loggerPkg.Info("用户删除成功", zap.String("user_id", id))
    return nil
}
```

---

## 🎯 最佳实践

### 1. MongoDB

- ✅ 使用索引优化查询性能
- ✅ 合理设计文档结构，避免过深嵌套
- ✅ 使用聚合管道处理复杂查询
- ✅ 定期备份数据
- ✅ 监控慢查询

### 2. Redis

- ✅ 合理设置过期时间，避免内存溢出
- ✅ 使用合适的数据结构
- ✅ 批量操作使用 Pipeline
- ✅ 避免大key，使用哈希分片
- ✅ 监控内存使用情况

### 3. Logger

- ✅ 使用结构化日志，便于检索
- ✅ 合理设置日志级别
- ✅ 敏感信息不要记录到日志
- ✅ 生产环境使用 JSON 格式输出
- ✅ 定期清理旧日志文件

---

**高效开发，稳定运行！🚀**
