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

// QualityInspectionRepository 质检记录仓储接口
type QualityInspectionRepository interface {
	Create(ctx context.Context, inspection *models.QualityInspection) error
	Get(ctx context.Context, id string) (*models.QualityInspection, error)
	List(ctx context.Context, page, pageSize int, inspectorID, contractNo string, startDate, endDate int64) ([]*models.QualityInspection, int64, error)
	UpdateReworkID(ctx context.Context, id, reworkID string) error
	GetStatistics(ctx context.Context, inspectorID string, startDate, endDate int64) (totalInspected int, totalQualified int, totalUnqualified int, err error)
	Delete(ctx context.Context, id string) error
}

type qualityInspectionRepository struct {
	dbManager *database.DatabaseManager
}

// NewQualityInspectionRepository 创建质检记录仓储
func NewQualityInspectionRepository() QualityInspectionRepository {
	return &qualityInspectionRepository{
		dbManager: database.GetDatabaseManager(),
	}
}

// GetCollectionWithContext 获取集合（支持租户上下文）
func (r *qualityInspectionRepository) GetCollectionWithContext(ctx context.Context) *mongo.Collection {
	tenantCode := tenantCtx.GetTenantCode(ctx)
	db := r.dbManager.GetDatabase(tenantCode)
	return db.Collection(models.QualityInspection{}.TableName())
}

// Create 创建质检记录
func (r *qualityInspectionRepository) Create(ctx context.Context, inspection *models.QualityInspection) error {
	collection := r.GetCollectionWithContext(ctx)
	_, err := collection.InsertOne(ctx, inspection)
	return err
}

// Get 根据ID获取质检记录
func (r *qualityInspectionRepository) Get(ctx context.Context, id string) (*models.QualityInspection, error) {
	collection := r.GetCollectionWithContext(ctx)

	var inspection models.QualityInspection
	err := collection.FindOne(ctx, bson.M{
		"_id":        id,
		"is_deleted": 0,
	}).Decode(&inspection)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &inspection, nil
}

// List 获取质检记录列表
func (r *qualityInspectionRepository) List(ctx context.Context, page, pageSize int, inspectorID, contractNo string, startDate, endDate int64) ([]*models.QualityInspection, int64, error) {
	collection := r.GetCollectionWithContext(ctx)

	// 构建过滤条件
	filter := bson.M{"is_deleted": 0}
	if inspectorID != "" {
		filter["inspector_id"] = inspectorID
	}
	if contractNo != "" {
		filter["contract_no"] = contractNo
	}
	if startDate > 0 && endDate > 0 {
		filter["inspection_time"] = bson.M{
			"$gte": startDate,
			"$lte": endDate,
		}
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
		SetSort(bson.D{{Key: "inspection_time", Value: -1}})

	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var inspections []*models.QualityInspection
	if err = cursor.All(ctx, &inspections); err != nil {
		return nil, 0, err
	}

	return inspections, total, nil
}

// UpdateReworkID 更新返工单ID
func (r *qualityInspectionRepository) UpdateReworkID(ctx context.Context, id, reworkID string) error {
	collection := r.GetCollectionWithContext(ctx)

	update := bson.M{
		"$set": bson.M{
			"rework_id":  reworkID,
			"updated_at": time.Now().Unix(),
		},
	}

	_, err := collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	return err
}

// GetStatistics 获取质检统计
func (r *qualityInspectionRepository) GetStatistics(ctx context.Context, inspectorID string, startDate, endDate int64) (totalInspected int, totalQualified int, totalUnqualified int, err error) {
	collection := r.GetCollectionWithContext(ctx)

	// 构建过滤条件
	filter := bson.M{"is_deleted": 0}
	if inspectorID != "" {
		filter["inspector_id"] = inspectorID
	}
	if startDate > 0 && endDate > 0 {
		filter["inspection_time"] = bson.M{
			"$gte": startDate,
			"$lte": endDate,
		}
	}

	// 聚合统计
	pipeline := []bson.M{
		{"$match": filter},
		{
			"$group": bson.M{
				"_id":               nil,
				"total_inspected":   bson.M{"$sum": "$inspected_qty"},
				"total_qualified":   bson.M{"$sum": "$qualified_qty"},
				"total_unqualified": bson.M{"$sum": "$unqualified_qty"},
			},
		},
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return 0, 0, 0, err
	}
	defer cursor.Close(ctx)

	var result []struct {
		TotalInspected   int `bson:"total_inspected"`
		TotalQualified   int `bson:"total_qualified"`
		TotalUnqualified int `bson:"total_unqualified"`
	}

	if err = cursor.All(ctx, &result); err != nil {
		return 0, 0, 0, err
	}

	if len(result) > 0 {
		return result[0].TotalInspected, result[0].TotalQualified, result[0].TotalUnqualified, nil
	}

	return 0, 0, 0, nil
}

// Delete 软删除质检记录
func (r *qualityInspectionRepository) Delete(ctx context.Context, id string) error {
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
