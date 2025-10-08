package database

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const (
	// SystemDatabase 系统数据库名称（存储租户元数据和系统超管）
	SystemDatabase = "mule_system"
)

// DatabaseManager 多租户数据库管理器
type DatabaseManager struct {
	client    *mongo.Client
	systemDB  *mongo.Database
	tenantDBs sync.Map // map[tenantID]*mongo.Database
	mu        sync.RWMutex
}

var (
	globalDBManager *DatabaseManager
	managerOnce     sync.Once
)

// InitDatabaseManager 初始化数据库管理器
func InitDatabaseManager(client *mongo.Client) *DatabaseManager {
	managerOnce.Do(func() {
		globalDBManager = &DatabaseManager{
			client:   client,
			systemDB: client.Database(SystemDatabase),
		}
		log.Printf("✅ 数据库管理器初始化成功，系统库: %s", SystemDatabase)
	})
	return globalDBManager
}

// GetDatabaseManager 获取全局数据库管理器
func GetDatabaseManager() *DatabaseManager {
	if globalDBManager == nil {
		log.Fatal("❌ 数据库管理器未初始化，请先调用 InitDatabaseManager()")
	}
	return globalDBManager
}

// GetDatabase 获取租户数据库（空tenantCode或"system"返回系统库）
// 注意：参数改为 tenantCode 而不是 tenantID，使数据库名更易读（如 mule_default）
func (m *DatabaseManager) GetDatabase(tenantCode string) *mongo.Database {
	// 系统超管使用系统库
	if tenantCode == "" || tenantCode == "system" {
		return m.systemDB
	}

	// 从缓存获取（使用 code 作为 key）
	if db, ok := m.tenantDBs.Load(tenantCode); ok {
		return db.(*mongo.Database)
	}

	// 创建新连接
	m.mu.Lock()
	defer m.mu.Unlock()

	// 双重检查
	if db, ok := m.tenantDBs.Load(tenantCode); ok {
		return db.(*mongo.Database)
	}

	dbName := GetTenantDatabaseName(tenantCode)
	db := m.client.Database(dbName)
	m.tenantDBs.Store(tenantCode, db)

	log.Printf("🔗 创建租户数据库连接: %s", dbName)
	return db
}

// GetSystemDatabase 获取系统数据库
func (m *DatabaseManager) GetSystemDatabase() *mongo.Database {
	return m.systemDB
}

// CreateTenantDatabase 创建租户数据库（初始化集合和索引）
// 参数使用 tenantCode 而不是 tenantID，这样数据库名更易读
func (m *DatabaseManager) CreateTenantDatabase(ctx context.Context, tenantCode string) error {
	dbName := GetTenantDatabaseName(tenantCode)
	db := m.client.Database(dbName)

	log.Printf("📦 开始创建租户数据库: %s", dbName)

	// 创建集合列表
	collections := []string{
		"admin", // 管理员
		"role",  // 角色
		// "menu",  // 菜单（租户可以自定义菜单，但通常从系统库同步）
		"basic", // 基础数据（统一存储：颜色、尺码、客户、订单类型、工序、业务员等，通过 type 字段区分）
		// ❌ 不再创建独立集合：color、size、customer、order_type、procedure、salesman
		// ✅ 所有基础数据统一存储在 basic 集合中
	}

	for _, collName := range collections {
		// 创建集合
		err := db.CreateCollection(ctx, collName)
		if err != nil {
			// 检查是否是集合已存在错误
			if !strings.Contains(err.Error(), "already exists") {
				return fmt.Errorf("创建集合 %s 失败: %v", collName, err)
			}
		}

		// 创建索引
		collection := db.Collection(collName)

		// 为所有集合创建 is_deleted 索引（软删除查询优化）
		_, err = collection.Indexes().CreateOne(ctx, mongo.IndexModel{
			Keys: bson.D{{Key: "is_deleted", Value: 1}},
		})
		if err != nil {
			log.Printf("⚠️  创建 is_deleted 索引失败 (%s): %v", collName, err)
		}

		// admin 集合特殊索引
		if collName == "admin" {
			// 手机号唯一索引
			_, err = collection.Indexes().CreateOne(ctx, mongo.IndexModel{
				Keys:    bson.D{{Key: "phone", Value: 1}},
				Options: options.Index().SetUnique(true).SetSparse(true),
			})
			if err != nil {
				log.Printf("⚠️  创建 phone 索引失败: %v", err)
			}

			// 邮箱索引
			_, err = collection.Indexes().CreateOne(ctx, mongo.IndexModel{
				Keys:    bson.D{{Key: "email", Value: 1}},
				Options: options.Index().SetSparse(true),
			})
			if err != nil {
				log.Printf("⚠️  创建 email 索引失败: %v", err)
			}
		}

		// basic 集合特殊索引
		if collName == "basic" {
			// type 索引（用于区分基础数据类型：color、size、customer 等）
			_, err = collection.Indexes().CreateOne(ctx, mongo.IndexModel{
				Keys: bson.D{{Key: "type", Value: 1}},
			})
			if err != nil {
				log.Printf("⚠️  创建 type 索引失败: %v", err)
			}

			// type + code 复合索引（确保同类型下 code 唯一）
			_, err = collection.Indexes().CreateOne(ctx, mongo.IndexModel{
				Keys: bson.D{
					{Key: "type", Value: 1},
					{Key: "code", Value: 1},
				},
				Options: options.Index().SetUnique(true).SetSparse(true),
			})
			if err != nil {
				log.Printf("⚠️  创建 type+code 复合索引失败: %v", err)
			}
		}

		log.Printf("  ✅ 集合 %s 创建成功", collName)
	}

	// 缓存数据库连接（使用 code 作为缓存 key）
	m.tenantDBs.Store(tenantCode, db)

	log.Printf("✅ 租户数据库创建完成: %s", dbName)
	return nil
}

// DeleteTenantDatabase 删除租户数据库（谨慎操作！）
func (m *DatabaseManager) DeleteTenantDatabase(ctx context.Context, tenantCode string) error {
	dbName := GetTenantDatabaseName(tenantCode)

	log.Printf("⚠️  准备删除租户数据库: %s", dbName)

	// 删除数据库
	err := m.client.Database(dbName).Drop(ctx)
	if err != nil {
		return fmt.Errorf("删除数据库失败: %v", err)
	}

	// 从缓存移除（使用 code 作为 key）
	m.tenantDBs.Delete(tenantCode)

	log.Printf("✅ 租户数据库已删除: %s", dbName)
	return nil
}

// GetTenantDatabaseName 获取租户数据库名称（使用租户代码）
func GetTenantDatabaseName(tenantCode string) string {
	return fmt.Sprintf("mule_%s", tenantCode)
}

// ListTenantDatabases 列出所有租户数据库
func (m *DatabaseManager) ListTenantDatabases(ctx context.Context) ([]string, error) {
	databases, err := m.client.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	tenantDBs := []string{}
	for _, db := range databases {
		if strings.HasPrefix(db, "tenant_") && db != SystemDatabase {
			tenantDBs = append(tenantDBs, db)
		}
	}

	return tenantDBs, nil
}

// CheckTenantDatabaseExists 检查租户数据库是否存在
func (m *DatabaseManager) CheckTenantDatabaseExists(ctx context.Context, tenantCode string) (bool, error) {
	dbName := GetTenantDatabaseName(tenantCode)
	databases, err := m.client.ListDatabaseNames(ctx, bson.M{"name": dbName})
	if err != nil {
		return false, err
	}

	for _, db := range databases {
		if db == dbName {
			return true, nil
		}
	}

	return false, nil
}

// GetDatabaseStats 获取数据库统计信息
func (m *DatabaseManager) GetDatabaseStats(ctx context.Context, tenantID string) (map[string]interface{}, error) {
	db := m.GetDatabase(tenantID)

	var result bson.M
	err := db.RunCommand(ctx, bson.D{{Key: "dbStats", Value: 1}}).Decode(&result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// CleanupInactiveDatabases 清理不活跃的数据库连接（释放内存）
func (m *DatabaseManager) CleanupInactiveDatabases() {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 清空缓存，下次访问时会重新创建连接
	m.tenantDBs = sync.Map{}
	log.Println("🧹 已清理不活跃的数据库连接缓存")
}

// HealthCheck 健康检查
func (m *DatabaseManager) HealthCheck(ctx context.Context) error {
	// 检查系统库连接
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err := m.client.Ping(ctx, nil)
	if err != nil {
		return fmt.Errorf("数据库连接失败: %v", err)
	}

	return nil
}
