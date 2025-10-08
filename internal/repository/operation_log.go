package repository

import (
	"context"
	tenantCtx "mule-cloud/core/context"
	"mule-cloud/core/database"
	"mule-cloud/internal/models"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type OperationLogRepository struct {
	dbManager *database.DatabaseManager
}

func NewOperationLogRepository() *OperationLogRepository {
	return &OperationLogRepository{
		dbManager: database.GetDatabaseManager(),
	}
}

// getCollection 获取操作日志集合
// 租户的操作日志存储在租户数据库，系统管理员的操作日志存储在系统数据库
func (r *OperationLogRepository) getCollection(ctx context.Context) *mongo.Collection {
	tenantCode := tenantCtx.GetTenantCode(ctx)
	db := r.dbManager.GetDatabase(tenantCode)
	return db.Collection("operation_logs")
}

// Create 创建操作日志
func (r *OperationLogRepository) Create(ctx context.Context, log *models.OperationLog) error {
	collection := r.getCollection(ctx)

	result, err := collection.InsertOne(ctx, log)
	if err != nil {
		return err
	}

	// 将 bson.ObjectID 转换为字符串
	if oid, ok := result.InsertedID.(bson.ObjectID); ok {
		log.ID = oid.Hex()
	}

	return nil
}

// List 查询操作日志（支持分页和筛选）
func (r *OperationLogRepository) List(ctx context.Context, filter bson.M, page, pageSize int) ([]*models.OperationLog, int64, error) {
	collection := r.getCollection(ctx)

	// 计算总数
	total, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	// 分页查询
	opts := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}}). // 按时间倒序
		SetSkip(int64((page - 1) * pageSize)).
		SetLimit(int64(pageSize))

	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var logs []*models.OperationLog
	if err = cursor.All(ctx, &logs); err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}

// GetByID 根据ID获取操作日志
func (r *OperationLogRepository) GetByID(ctx context.Context, id string) (*models.OperationLog, error) {
	collection := r.getCollection(ctx)

	// 将字符串 ID 转换为 ObjectID
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var log models.OperationLog
	err = collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&log)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &log, nil
}

// CreateIndexes 创建索引（需要在各个数据库中分别创建）
func (r *OperationLogRepository) CreateIndexes(ctx context.Context) error {
	collection := r.getCollection(ctx)

	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "user_id", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "created_at", Value: -1}},
		},
		{
			Keys: bson.D{
				{Key: "user_id", Value: 1},
				{Key: "created_at", Value: -1},
			},
		},
	}

	_, err := collection.Indexes().CreateMany(ctx, indexes)
	return err
}
