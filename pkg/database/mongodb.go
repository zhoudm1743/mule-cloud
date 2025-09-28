package database

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// MongoConfig MongoDB配置
type MongoConfig struct {
	URI      string
	Database string
	Username string
	Password string
}

// MongoDB MongoDB客户端
type MongoDB struct {
	client   *mongo.Client
	database *mongo.Database
	config   MongoConfig
}

// NewMongoDB 创建MongoDB客户端
func NewMongoDB(config MongoConfig) (*MongoDB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 创建客户端选项
	clientOptions := options.Client().ApplyURI(config.URI)

	if config.Username != "" && config.Password != "" {
		clientOptions.SetAuth(options.Credential{
			Username: config.Username,
			Password: config.Password,
		})
	}

	// 连接MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// 测试连接
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	database := client.Database(config.Database)

	return &MongoDB{
		client:   client,
		database: database,
		config:   config,
	}, nil
}

// GetClient 获取MongoDB客户端
func (m *MongoDB) GetClient() *mongo.Client {
	return m.client
}

// GetDatabase 获取数据库
func (m *MongoDB) GetDatabase() *mongo.Database {
	return m.database
}

// GetCollection 获取集合
func (m *MongoDB) GetCollection(name string) *mongo.Collection {
	return m.database.Collection(name)
}

// Close 关闭连接
func (m *MongoDB) Close(ctx context.Context) error {
	return m.client.Disconnect(ctx)
}

// Health 健康检查
func (m *MongoDB) Health(ctx context.Context) error {
	return m.client.Ping(ctx, readpref.Primary())
}

// CreateIndexes 创建索引
func (m *MongoDB) CreateIndexes(ctx context.Context) error {
	// 用户索引
	userCollection := m.GetCollection("users")
	userIndexes := []mongo.IndexModel{
		{
			Keys:    map[string]interface{}{"username": 1},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys:    map[string]interface{}{"email": 1},
			Options: options.Index().SetUnique(true),
		},
	}
	_, err := userCollection.Indexes().CreateMany(ctx, userIndexes)
	if err != nil {
		return fmt.Errorf("failed to create user indexes: %w", err)
	}

	// 订单索引
	orderCollection := m.GetCollection("orders")
	orderIndexes := []mongo.IndexModel{
		{
			Keys:    map[string]interface{}{"order_no": 1},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: map[string]interface{}{
				"customer_id": 1,
				"created_at":  -1,
			},
		},
		{
			Keys: map[string]interface{}{
				"status":     1,
				"created_at": -1,
			},
		},
	}
	_, err = orderCollection.Indexes().CreateMany(ctx, orderIndexes)
	if err != nil {
		return fmt.Errorf("failed to create order indexes: %w", err)
	}

	// 工作报告索引
	workReportCollection := m.GetCollection("work_reports")
	workReportIndexes := []mongo.IndexModel{
		{
			Keys: map[string]interface{}{
				"worker_id": 1,
				"date":      -1,
			},
		},
		{
			Keys: map[string]interface{}{
				"order_id":   1,
				"process_id": 1,
			},
		},
	}
	_, err = workReportCollection.Indexes().CreateMany(ctx, workReportIndexes)
	if err != nil {
		return fmt.Errorf("failed to create work report indexes: %w", err)
	}

	// 生产进度索引
	progressCollection := m.GetCollection("production_progress")
	progressIndexes := []mongo.IndexModel{
		{
			Keys: map[string]interface{}{
				"order_id":   1,
				"updated_at": -1,
			},
		},
	}
	_, err = progressCollection.Indexes().CreateMany(ctx, progressIndexes)
	if err != nil {
		return fmt.Errorf("failed to create production progress indexes: %w", err)
	}

	// 工时记录索引
	timesheetCollection := m.GetCollection("timesheets")
	timesheetIndexes := []mongo.IndexModel{
		{
			Keys: map[string]interface{}{
				"worker_id": 1,
				"date":      -1,
			},
		},
	}
	_, err = timesheetCollection.Indexes().CreateMany(ctx, timesheetIndexes)
	if err != nil {
		return fmt.Errorf("failed to create timesheet indexes: %w", err)
	}

	// 工资单索引
	payrollCollection := m.GetCollection("payrolls")
	payrollIndexes := []mongo.IndexModel{
		{
			Keys: map[string]interface{}{
				"worker_id":  1,
				"pay_period": -1,
			},
		},
	}
	_, err = payrollCollection.Indexes().CreateMany(ctx, payrollIndexes)
	if err != nil {
		return fmt.Errorf("failed to create payroll indexes: %w", err)
	}

	return nil
}
