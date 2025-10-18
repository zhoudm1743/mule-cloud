package repository

import (
	"context"
	"time"

	tenantCtx "mule-cloud/core/context"
	"mule-cloud/core/database"
	"mule-cloud/internal/models"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// ReworkRepository 返工记录仓储接口
type ReworkRepository interface {
	Create(ctx context.Context, rework *models.ReworkRecord) error
	Get(ctx context.Context, id string) (*models.ReworkRecord, error)
	List(ctx context.Context, page, pageSize int, status *int, workerID, contractNo string) ([]*models.ReworkRecord, int64, error)
	UpdateStatus(ctx context.Context, id string, status int) error
	Complete(ctx context.Context, id string, images []string, remark string) error
	GetStatistics(ctx context.Context, workerID string) (total int, pending int, inProgress int, completed int, err error)
	Delete(ctx context.Context, id string) error
}

type reworkRepository struct {
	dbManager *database.DatabaseManager
}

// NewReworkRepository 创建返工记录仓储
func NewReworkRepository() ReworkRepository {
	return &reworkRepository{
		dbManager: database.GetDatabaseManager(),
	}
}

// GetCollectionWithContext 获取集合（支持租户上下文）
func (r *reworkRepository) GetCollectionWithContext(ctx context.Context) *mongo.Collection {
	tenantCode := tenantCtx.GetTenantCode(ctx)
	db := r.dbManager.GetDatabase(tenantCode)
	return db.Collection(models.ReworkRecord{}.TableName())
}

// Create 创建返工记录
func (r *reworkRepository) Create(ctx context.Context, rework *models.ReworkRecord) error {
	collection := r.GetCollectionWithContext(ctx)
	_, err := collection.InsertOne(ctx, rework)
	return err
}

// Get 根据ID获取返工记录
func (r *reworkRepository) Get(ctx context.Context, id string) (*models.ReworkRecord, error) {
	collection := r.GetCollectionWithContext(ctx)

	var rework models.ReworkRecord
	err := collection.FindOne(ctx, bson.M{
		"_id":        id,
		"is_deleted": 0,
	}).Decode(&rework)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &rework, nil
}

// List 获取返工记录列表
func (r *reworkRepository) List(ctx context.Context, page, pageSize int, status *int, workerID, contractNo string) ([]*models.ReworkRecord, int64, error) {
	collection := r.GetCollectionWithContext(ctx)

	// 构建过滤条件
	filter := bson.M{"is_deleted": 0}
	if status != nil {
		filter["status"] = *status
	}
	if workerID != "" {
		filter["assigned_worker"] = workerID
	}
	if contractNo != "" {
		filter["contract_no"] = contractNo
	}

	// 获取总数
	total, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	// 分页查询
	skip := (page - 1) * pageSize
	opts := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(pageSize)).
		SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var reworks []*models.ReworkRecord
	if err = cursor.All(ctx, &reworks); err != nil {
		return nil, 0, err
	}

	return reworks, total, nil
}

// UpdateStatus 更新返工状态
func (r *reworkRepository) UpdateStatus(ctx context.Context, id string, status int) error {
	collection := r.GetCollectionWithContext(ctx)

	update := bson.M{
		"$set": bson.M{
			"status":     status,
			"updated_at": time.Now().Unix(),
		},
	}

	_, err := collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	return err
}

// Complete 完成返工
func (r *reworkRepository) Complete(ctx context.Context, id string, images []string, remark string) error {
	collection := r.GetCollectionWithContext(ctx)

	update := bson.M{
		"$set": bson.M{
			"status":       2, // 已完成
			"images":       images,
			"remark":       remark,
			"completed_at": time.Now().Unix(),
			"updated_at":   time.Now().Unix(),
		},
	}

	_, err := collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	return err
}

// GetStatistics 获取返工统计
func (r *reworkRepository) GetStatistics(ctx context.Context, workerID string) (total int, pending int, inProgress int, completed int, err error) {
	collection := r.GetCollectionWithContext(ctx)

	// 构建过滤条件
	filter := bson.M{"is_deleted": 0}
	if workerID != "" {
		filter["assigned_worker"] = workerID
	}

	// 聚合统计
	pipeline := []bson.M{
		{"$match": filter},
		{
			"$group": bson.M{
				"_id":   "$status",
				"count": bson.M{"$sum": 1},
			},
		},
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return 0, 0, 0, 0, err
	}
	defer cursor.Close(ctx)

	var results []struct {
		Status int `bson:"_id"`
		Count  int `bson:"count"`
	}

	if err = cursor.All(ctx, &results); err != nil {
		return 0, 0, 0, 0, err
	}

	for _, result := range results {
		total += result.Count
		switch result.Status {
		case 0:
			pending = result.Count
		case 1:
			inProgress = result.Count
		case 2:
			completed = result.Count
		}
	}

	return total, pending, inProgress, completed, nil
}

// Delete 软删除返工记录
func (r *reworkRepository) Delete(ctx context.Context, id string) error {
	collection := r.GetCollectionWithContext(ctx)

	update := bson.M{
		"$set": bson.M{
			"is_deleted": 1,
			"updated_at": time.Now().Unix(),
		},
	}

	_, err := collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	return err
}
