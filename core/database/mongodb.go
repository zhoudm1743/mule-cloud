package database

import (
	"context"
	"fmt"
	"log"
	"mule-cloud/core/config"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/v2/event"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

var (
	mongoClient *mongo.Client
	mongoDB     *mongo.Database
	mongoOnce   sync.Once
	mongoErr    error
)

// MongoDB 全局MongoDB实例（懒加载）
var MongoDB = &MongoDBInstance{}

// InitMongoDB 初始化MongoDB连接
func InitMongoDB(cfg *config.MongoDBConfig) (*mongo.Client, error) {
	if !cfg.Enabled {
		log.Println("⚠️  MongoDB未启用")
		return nil, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.Timeout)*time.Second)
	defer cancel()

	// 构建连接选项
	clientOpts := options.Client()

	// 方式1: 使用URI（推荐）
	if cfg.URI != "" {
		clientOpts.ApplyURI(cfg.URI)
	} else {
		// 方式2: 使用独立配置构建URI
		uri := buildMongoDBURI(cfg)
		clientOpts.ApplyURI(uri)
	}

	// 配置BSON选项：将ObjectID自动转换为十六进制字符串
	bsonOpts := &options.BSONOptions{
		ObjectIDAsHexString: true,
	}
	clientOpts.SetBSONOptions(bsonOpts)

	// 连接池配置
	if cfg.MaxPoolSize > 0 {
		clientOpts.SetMaxPoolSize(cfg.MaxPoolSize)
	}
	if cfg.MinPoolSize > 0 {
		clientOpts.SetMinPoolSize(cfg.MinPoolSize)
	}

	// 副本集配置
	if cfg.ReplicaSet != "" {
		clientOpts.SetReplicaSet(cfg.ReplicaSet)
	}

	// 命令监控（开发环境）
	cmdMonitor := &event.CommandMonitor{
		Started: func(ctx context.Context, e *event.CommandStartedEvent) {
			log.Printf("[MongoDB] 执行命令: %s, 命令内容: %v", e.CommandName, e.Command)
		},
		Succeeded: func(ctx context.Context, e *event.CommandSucceededEvent) {
			// log.Printf("[MongoDB] 命令成功: %s, 耗时: %dms", e.CommandName, e.Duration.Milliseconds())
		},
		Failed: func(ctx context.Context, e *event.CommandFailedEvent) {
			log.Printf("[MongoDB] 命令失败: %s, 错误: %v", e.CommandName, e.Failure)
		},
	}
	clientOpts.SetMonitor(cmdMonitor)

	// 连接MongoDB
	client, err := mongo.Connect(clientOpts)
	if err != nil {
		return nil, fmt.Errorf("连接MongoDB失败: %v", err)
	}

	// Ping测试
	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, fmt.Errorf("Ping MongoDB失败: %v", err)
	}

	log.Printf("✅ MongoDB连接成功: %s:%d/%s", cfg.Host, cfg.Port, cfg.Database)

	// 保存全局实例
	mongoClient = client
	mongoDB = client.Database(cfg.Database)

	return client, nil
}

// buildMongoDBURI 构建MongoDB连接URI
func buildMongoDBURI(cfg *config.MongoDBConfig) string {
	// 基本URI格式: mongodb://[username:password@]host:port[/database][?options]
	uri := "mongodb://"

	// 用户认证
	if cfg.Username != "" && cfg.Password != "" {
		uri += fmt.Sprintf("%s:%s@", cfg.Username, cfg.Password)
	}

	// 主机和端口
	uri += fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)

	// 数据库
	if cfg.Database != "" {
		uri += fmt.Sprintf("/%s", cfg.Database)
	}

	// 认证数据库
	if cfg.AuthSource != "" {
		uri += fmt.Sprintf("?authSource=%s", cfg.AuthSource)
	}

	return uri
}

// GetMongoDB 获取MongoDB数据库实例
func GetMongoDB() *mongo.Database {
	if mongoDB == nil {
		log.Fatal("MongoDB未初始化，请先调用 InitMongoDB()")
	}
	return mongoDB
}

// GetMongoClient 获取MongoDB客户端实例
func GetMongoClient() *mongo.Client {
	if mongoClient == nil {
		log.Fatal("MongoDB未初始化，请先调用 InitMongoDB()")
	}
	return mongoClient
}

// GetCollection 获取MongoDB集合
func GetCollection(name string) *mongo.Collection {
	return GetMongoDB().Collection(name)
}

// CloseMongoDB 关闭MongoDB连接
func CloseMongoDB() error {
	if mongoClient == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := mongoClient.Disconnect(ctx); err != nil {
		return fmt.Errorf("关闭MongoDB连接失败: %v", err)
	}

	log.Println("✅ MongoDB连接已关闭")
	return nil
}

// PingMongoDB 检查MongoDB连接状态
func PingMongoDB() error {
	if mongoClient == nil {
		return fmt.Errorf("MongoDB未初始化")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return mongoClient.Ping(ctx, readpref.Primary())
}

// MongoDBHealth 获取MongoDB健康状态
func MongoDBHealth() map[string]interface{} {
	status := map[string]interface{}{
		"status": "unknown",
	}

	if mongoClient == nil {
		status["status"] = "not_initialized"
		return status
	}

	if err := PingMongoDB(); err != nil {
		status["status"] = "unhealthy"
		status["error"] = err.Error()
	} else {
		status["status"] = "healthy"
		status["database"] = mongoDB.Name()
	}

	return status
}

// MongoDBInstance 全局MongoDB实例包装器
type MongoDBInstance struct{}

// Init 初始化MongoDB（从全局配置）
func (m *MongoDBInstance) Init() error {
	cfg := config.Get()
	if !cfg.MongoDB.Enabled {
		return fmt.Errorf("MongoDB未启用")
	}
	_, err := InitMongoDB(&cfg.MongoDB)
	return err
}

// AutoInit 自动初始化（只执行一次）
func (m *MongoDBInstance) AutoInit() error {
	mongoOnce.Do(func() {
		cfg := config.Get()
		if cfg.MongoDB.Enabled {
			_, mongoErr = InitMongoDB(&cfg.MongoDB)
		}
	})
	return mongoErr
}

// DB 获取数据库实例（自动初始化）
func (m *MongoDBInstance) DB() *mongo.Database {
	if mongoDB == nil {
		m.AutoInit()
	}
	return mongoDB
}

// Client 获取客户端实例（自动初始化）
func (m *MongoDBInstance) Client() *mongo.Client {
	if mongoClient == nil {
		m.AutoInit()
	}
	return mongoClient
}

// Collection 获取集合（自动初始化）
func (m *MongoDBInstance) Collection(name string) *mongo.Collection {
	return m.DB().Collection(name)
}

// IsConnected 检查是否已连接
func (m *MongoDBInstance) IsConnected() bool {
	return mongoClient != nil
}

// Close 关闭连接
func (m *MongoDBInstance) Close() error {
	return CloseMongoDB()
}
