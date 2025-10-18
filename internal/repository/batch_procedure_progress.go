package repository

import (
	"context"
	"time"

	tenantCtx "mule-cloud/core/context"
	"mule-cloud/core/database"
	"mule-cloud/internal/models"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// BatchProcedureProgressRepository 批次工序进度仓储接口
type BatchProcedureProgressRepository interface {
	Create(ctx context.Context, progress *models.BatchProcedureProgress) error
	GetByBatchAndProcedure(ctx context.Context, batchID string, procedureSeq int) (*models.BatchProcedureProgress, error)
	ListByBatch(ctx context.Context, batchID string) ([]*models.BatchProcedureProgress, error)
	UpdateReportedQty(ctx context.Context, batchID string, procedureSeq int, quantity int) error
	InitBatchProgress(ctx context.Context, batchID, bundleNo, orderID string, quantity int, procedures []models.OrderProcedure) error
}

type batchProcedureProgressRepository struct {
	dbManager *database.DatabaseManager
}

// NewBatchProcedureProgressRepository 创建批次工序进度仓储
func NewBatchProcedureProgressRepository() BatchProcedureProgressRepository {
	return &batchProcedureProgressRepository{
		dbManager: database.GetDatabaseManager(),
	}
}

// GetCollectionWithContext 获取集合（支持租户上下文）
func (r *batchProcedureProgressRepository) GetCollectionWithContext(ctx context.Context) *mongo.Collection {
	tenantCode := tenantCtx.GetTenantCode(ctx)
	db := r.dbManager.GetDatabase(tenantCode)
	return db.Collection(models.BatchProcedureProgress{}.TableName())
}

// Create 创建批次工序进度
func (r *batchProcedureProgressRepository) Create(ctx context.Context, progress *models.BatchProcedureProgress) error {
	collection := r.GetCollectionWithContext(ctx)
	_, err := collection.InsertOne(ctx, progress)
	return err
}

// GetByBatchAndProcedure 根据批次ID和工序序号获取进度
func (r *batchProcedureProgressRepository) GetByBatchAndProcedure(ctx context.Context, batchID string, procedureSeq int) (*models.BatchProcedureProgress, error) {
	collection := r.GetCollectionWithContext(ctx)

	var progress models.BatchProcedureProgress
	err := collection.FindOne(ctx, bson.M{
		"batch_id":      batchID,
		"procedure_seq": procedureSeq,
	}).Decode(&progress)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &progress, nil
}

// ListByBatch 获取批次的所有工序进度
func (r *batchProcedureProgressRepository) ListByBatch(ctx context.Context, batchID string) ([]*models.BatchProcedureProgress, error) {
	collection := r.GetCollectionWithContext(ctx)

	cursor, err := collection.Find(ctx, bson.M{"batch_id": batchID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var progressList []*models.BatchProcedureProgress
	if err = cursor.All(ctx, &progressList); err != nil {
		return nil, err
	}

	return progressList, nil
}

// UpdateReportedQty 更新已上报数量
func (r *batchProcedureProgressRepository) UpdateReportedQty(ctx context.Context, batchID string, procedureSeq int, quantity int) error {
	collection := r.GetCollectionWithContext(ctx)

	filter := bson.M{
		"batch_id":      batchID,
		"procedure_seq": procedureSeq,
	}

	// 先查询当前进度
	progress, err := r.GetByBatchAndProcedure(ctx, batchID, procedureSeq)
	if err != nil {
		return err
	}

	// 更新已上报数量
	newReportedQty := progress.ReportedQty + quantity
	isCompleted := newReportedQty >= progress.Quantity

	update := bson.M{
		"$set": bson.M{
			"reported_qty": newReportedQty,
			"is_completed": isCompleted,
			"updated_at":   time.Now().Unix(),
		},
	}

	if isCompleted {
		update["$set"].(bson.M)["completed_at"] = time.Now().Unix()
	}

	_, err = collection.UpdateOne(ctx, filter, update)
	return err
}

// InitBatchProgress 初始化批次的所有工序进度
func (r *batchProcedureProgressRepository) InitBatchProgress(ctx context.Context, batchID, bundleNo, orderID string, quantity int, procedures []models.OrderProcedure) error {
	collection := r.GetCollectionWithContext(ctx)

	// 批量插入所有工序的进度记录
	var documents []interface{}
	now := time.Now().Unix()

	for _, proc := range procedures {
		progress := &models.BatchProcedureProgress{
			ID:            bson.NewObjectID().Hex(),
			BatchID:       batchID,
			BundleNo:      bundleNo,
			OrderID:       orderID,
			ProcedureSeq:  proc.Sequence,
			ProcedureName: proc.ProcedureName,
			Quantity:      quantity,
			ReportedQty:   0,
			IsCompleted:   false,
			CreatedAt:     now,
			UpdatedAt:     now,
		}
		documents = append(documents, progress)
	}

	if len(documents) > 0 {
		_, err := collection.InsertMany(ctx, documents)
		return err
	}

	return nil
}
