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

// ProcedureReportRepository 工序上报记录仓储接口
type ProcedureReportRepository interface {
	Create(ctx context.Context, report *models.ProcedureReport) error
	GetByID(ctx context.Context, id string) (*models.ProcedureReport, error)
	List(ctx context.Context, page, pageSize int, workerID, contractNo, startDate, endDate string) ([]*models.ProcedureReport, int64, error)
	GetStatistics(ctx context.Context, workerID, startDate, endDate string) (totalQuantity int, totalAmount float64, err error)
	GetSalaryDetails(ctx context.Context, workerID, startDate, endDate string) ([]map[string]interface{}, error)
	Delete(ctx context.Context, id string) error
}

type procedureReportRepository struct {
	dbManager *database.DatabaseManager
}

// NewProcedureReportRepository 创建工序上报记录仓储
func NewProcedureReportRepository() ProcedureReportRepository {
	return &procedureReportRepository{
		dbManager: database.GetDatabaseManager(),
	}
}

// GetCollectionWithContext 获取集合（支持租户上下文）
func (r *procedureReportRepository) GetCollectionWithContext(ctx context.Context) *mongo.Collection {
	tenantCode := tenantCtx.GetTenantCode(ctx)
	db := r.dbManager.GetDatabase(tenantCode)
	return db.Collection(models.ProcedureReport{}.TableName())
}

// Create 创建工序上报记录
func (r *procedureReportRepository) Create(ctx context.Context, report *models.ProcedureReport) error {
	collection := r.GetCollectionWithContext(ctx)
	_, err := collection.InsertOne(ctx, report)
	return err
}

// GetByID 根据ID获取工序上报记录
func (r *procedureReportRepository) GetByID(ctx context.Context, id string) (*models.ProcedureReport, error) {
	collection := r.GetCollectionWithContext(ctx)

	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var report models.ProcedureReport
	err = collection.FindOne(ctx, bson.M{"_id": objectID, "is_deleted": 0}).Decode(&report)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &report, nil
}

// List 获取工序上报记录列表
func (r *procedureReportRepository) List(ctx context.Context, page, pageSize int, workerID, contractNo, startDate, endDate string) ([]*models.ProcedureReport, int64, error) {
	collection := r.GetCollectionWithContext(ctx)

	// 构建查询条件
	filter := bson.M{"is_deleted": 0}

	if workerID != "" {
		filter["worker_id"] = workerID
	}

	if contractNo != "" {
		filter["contract_no"] = bson.M{"$regex": contractNo, "$options": "i"}
	}

	// 日期范围查询
	if startDate != "" || endDate != "" {
		dateFilter := bson.M{}
		if startDate != "" {
			startTime, _ := time.Parse("2006-01-02", startDate)
			dateFilter["$gte"] = startTime.Unix()
		}
		if endDate != "" {
			endTime, _ := time.Parse("2006-01-02", endDate)
			endTime = endTime.Add(24*time.Hour - time.Second) // 包含当天结束
			dateFilter["$lte"] = endTime.Unix()
		}
		filter["report_time"] = dateFilter
	}

	// 计算总数
	total, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	// 查询数据
	skip := int64((page - 1) * pageSize)
	limit := int64(pageSize)
	opts := options.Find().
		SetSkip(skip).
		SetLimit(limit).
		SetSort(bson.D{{Key: "report_time", Value: -1}}) // 按上报时间倒序

	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var reports []*models.ProcedureReport
	if err = cursor.All(ctx, &reports); err != nil {
		return nil, 0, err
	}

	return reports, total, nil
}

// GetStatistics 获取统计数据
func (r *procedureReportRepository) GetStatistics(ctx context.Context, workerID, startDate, endDate string) (totalQuantity int, totalAmount float64, err error) {
	collection := r.GetCollectionWithContext(ctx)

	// 构建查询条件
	filter := bson.M{"is_deleted": 0}

	if workerID != "" {
		filter["worker_id"] = workerID
	}

	// 日期范围查询
	if startDate != "" || endDate != "" {
		dateFilter := bson.M{}
		if startDate != "" {
			startTime, _ := time.Parse("2006-01-02", startDate)
			dateFilter["$gte"] = startTime.Unix()
		}
		if endDate != "" {
			endTime, _ := time.Parse("2006-01-02", endDate)
			endTime = endTime.Add(24*time.Hour - time.Second)
			dateFilter["$lte"] = endTime.Unix()
		}
		filter["report_time"] = dateFilter
	}

	// 聚合统计
	pipeline := []bson.D{
		{{Key: "$match", Value: filter}},
		{{Key: "$group", Value: bson.M{
			"_id":            nil,
			"total_quantity": bson.M{"$sum": "$quantity"},
			"total_amount":   bson.M{"$sum": "$total_price"},
		}}},
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return 0, 0, err
	}
	defer cursor.Close(ctx)

	var result []struct {
		TotalQuantity int     `bson:"total_quantity"`
		TotalAmount   float64 `bson:"total_amount"`
	}
	if err = cursor.All(ctx, &result); err != nil {
		return 0, 0, err
	}

	if len(result) > 0 {
		return result[0].TotalQuantity, result[0].TotalAmount, nil
	}

	return 0, 0, nil
}

// GetSalaryDetails 获取工资明细（按工序分组）
func (r *procedureReportRepository) GetSalaryDetails(ctx context.Context, workerID, startDate, endDate string) ([]map[string]interface{}, error) {
	collection := r.GetCollectionWithContext(ctx)

	// 构建查询条件
	filter := bson.M{"is_deleted": 0}

	if workerID != "" {
		filter["worker_id"] = workerID
	}

	// 日期范围查询
	if startDate != "" || endDate != "" {
		dateFilter := bson.M{}
		if startDate != "" {
			startTime, _ := time.Parse("2006-01-02", startDate)
			dateFilter["$gte"] = startTime.Unix()
		}
		if endDate != "" {
			endTime, _ := time.Parse("2006-01-02", endDate)
			endTime = endTime.Add(24*time.Hour - time.Second)
			dateFilter["$lte"] = endTime.Unix()
		}
		filter["report_time"] = dateFilter
	}

	// 按工序分组统计
	pipeline := []bson.D{
		{{Key: "$match", Value: filter}},
		{{Key: "$group", Value: bson.M{
			"_id": bson.M{
				"procedure_name": "$procedure_name",
				"unit_price":     "$unit_price",
			},
			"quantity": bson.M{"$sum": "$quantity"},
			"amount":   bson.M{"$sum": "$total_price"},
		}}},
		{{Key: "$project", Value: bson.M{
			"_id":            0,
			"procedure_name": "$_id.procedure_name",
			"unit_price":     "$_id.unit_price",
			"quantity":       1,
			"amount":         1,
		}}},
		{{Key: "$sort", Value: bson.D{{Key: "amount", Value: -1}}}}, // 按金额倒序
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var details []map[string]interface{}
	if err = cursor.All(ctx, &details); err != nil {
		return nil, err
	}

	return details, nil
}

// Delete 删除工序上报记录（软删除）
func (r *procedureReportRepository) Delete(ctx context.Context, id string) error {
	collection := r.GetCollectionWithContext(ctx)

	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	update := bson.M{
		"$set": bson.M{
			"is_deleted": 1,
			"updated_at": time.Now().Unix(),
		},
	}
	_, err = collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	return err
}
