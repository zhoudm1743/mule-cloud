package repository

import (
	"context"
	tenantCtx "mule-cloud/core/context"
	"mule-cloud/core/database"
	"mule-cloud/internal/models"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// CuttingTaskRepository 裁剪任务仓库接口
type CuttingTaskRepository interface {
	Create(ctx context.Context, task *models.CuttingTask) error
	Update(ctx context.Context, id string, task *models.CuttingTask) error
	Delete(ctx context.Context, id string) error
	GetByID(ctx context.Context, id string) (*models.CuttingTask, error)
	GetByOrderID(ctx context.Context, orderID string) (*models.CuttingTask, error)
	List(ctx context.Context, page, pageSize int, contractNo, styleNo string, status *int) ([]*models.CuttingTask, int64, error)
}

// CuttingBatchRepository 裁剪批次仓库接口
type CuttingBatchRepository interface {
	Create(ctx context.Context, batch *models.CuttingBatch) error
	Update(ctx context.Context, id string, batch *models.CuttingBatch) error
	Delete(ctx context.Context, id string) error
	GetByID(ctx context.Context, id string) (*models.CuttingBatch, error)
	List(ctx context.Context, page, pageSize int, taskID, contractNo, bedNo, bundleNo string) ([]*models.CuttingBatch, int64, error)
}

// CuttingPieceRepository 裁片监控仓库接口
type CuttingPieceRepository interface {
	Create(ctx context.Context, piece *models.CuttingPiece) error
	Update(ctx context.Context, id string, piece *models.CuttingPiece) error
	GetByID(ctx context.Context, id string) (*models.CuttingPiece, error)
	List(ctx context.Context, page, pageSize int, orderID, contractNo, bedNo, bundleNo string) ([]*models.CuttingPiece, int64, error)
}

// ==================== 裁剪任务仓库实现 ====================

type cuttingTaskRepository struct {
	dbManager *database.DatabaseManager
}

func NewCuttingTaskRepository() CuttingTaskRepository {
	return &cuttingTaskRepository{
		dbManager: database.GetDatabaseManager(),
	}
}

func (r *cuttingTaskRepository) GetCollectionWithContext(ctx context.Context) *mongo.Collection {
	tenantID := tenantCtx.GetTenantID(ctx)
	db := r.dbManager.GetDatabase(tenantID)
	return db.Collection(models.CuttingTask{}.TableName())
}

func (r *cuttingTaskRepository) Create(ctx context.Context, task *models.CuttingTask) error {
	collection := r.GetCollectionWithContext(ctx)
	_, err := collection.InsertOne(ctx, task)
	return err
}

func (r *cuttingTaskRepository) Update(ctx context.Context, id string, task *models.CuttingTask) error {
	collection := r.GetCollectionWithContext(ctx)
	_, err := collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": task})
	return err
}

func (r *cuttingTaskRepository) Delete(ctx context.Context, id string) error {
	collection := r.GetCollectionWithContext(ctx)
	now := time.Now().Unix()
	_, err := collection.UpdateOne(ctx,
		bson.M{"_id": id},
		bson.M{"$set": bson.M{"is_deleted": 1, "deleted_at": now}},
	)
	return err
}

func (r *cuttingTaskRepository) GetByID(ctx context.Context, id string) (*models.CuttingTask, error) {
	collection := r.GetCollectionWithContext(ctx)
	var task models.CuttingTask
	err := collection.FindOne(ctx, bson.M{"_id": id, "is_deleted": 0}).Decode(&task)
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *cuttingTaskRepository) GetByOrderID(ctx context.Context, orderID string) (*models.CuttingTask, error) {
	collection := r.GetCollectionWithContext(ctx)
	var task models.CuttingTask
	err := collection.FindOne(ctx, bson.M{"order_id": orderID, "is_deleted": 0}).Decode(&task)
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *cuttingTaskRepository) List(ctx context.Context, page, pageSize int, contractNo, styleNo string, status *int) ([]*models.CuttingTask, int64, error) {
	collection := r.GetCollectionWithContext(ctx)

	filter := bson.M{"is_deleted": 0}
	if contractNo != "" {
		filter["contract_no"] = bson.M{"$regex": contractNo, "$options": "i"}
	}
	if styleNo != "" {
		filter["style_no"] = bson.M{"$regex": styleNo, "$options": "i"}
	}
	if status != nil {
		filter["status"] = *status
	}

	total, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	skip := int64((page - 1) * pageSize)
	limit := int64(pageSize)
	opts := options.Find().SetSkip(skip).SetLimit(limit).SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var tasks []*models.CuttingTask
	if err = cursor.All(ctx, &tasks); err != nil {
		return nil, 0, err
	}

	return tasks, total, nil
}

// ==================== 裁剪批次仓库实现 ====================

type cuttingBatchRepository struct {
	dbManager *database.DatabaseManager
}

func NewCuttingBatchRepository() CuttingBatchRepository {
	return &cuttingBatchRepository{
		dbManager: database.GetDatabaseManager(),
	}
}

func (r *cuttingBatchRepository) GetCollectionWithContext(ctx context.Context) *mongo.Collection {
	tenantID := tenantCtx.GetTenantID(ctx)
	db := r.dbManager.GetDatabase(tenantID)
	return db.Collection(models.CuttingBatch{}.TableName())
}

func (r *cuttingBatchRepository) Create(ctx context.Context, batch *models.CuttingBatch) error {
	collection := r.GetCollectionWithContext(ctx)
	_, err := collection.InsertOne(ctx, batch)
	return err
}

func (r *cuttingBatchRepository) Update(ctx context.Context, id string, batch *models.CuttingBatch) error {
	collection := r.GetCollectionWithContext(ctx)
	_, err := collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": batch})
	return err
}

func (r *cuttingBatchRepository) Delete(ctx context.Context, id string) error {
	collection := r.GetCollectionWithContext(ctx)
	_, err := collection.UpdateOne(ctx,
		bson.M{"_id": id},
		bson.M{"$set": bson.M{"is_deleted": 1}},
	)
	return err
}

func (r *cuttingBatchRepository) GetByID(ctx context.Context, id string) (*models.CuttingBatch, error) {
	collection := r.GetCollectionWithContext(ctx)
	var batch models.CuttingBatch
	err := collection.FindOne(ctx, bson.M{"_id": id, "is_deleted": 0}).Decode(&batch)
	if err != nil {
		return nil, err
	}
	return &batch, nil
}

func (r *cuttingBatchRepository) List(ctx context.Context, page, pageSize int, taskID, contractNo, bedNo, bundleNo string) ([]*models.CuttingBatch, int64, error) {
	collection := r.GetCollectionWithContext(ctx)

	filter := bson.M{"is_deleted": 0}
	if taskID != "" {
		filter["task_id"] = taskID
	}
	if contractNo != "" {
		filter["contract_no"] = bson.M{"$regex": contractNo, "$options": "i"}
	}
	if bedNo != "" {
		filter["bed_no"] = bedNo
	}
	if bundleNo != "" {
		filter["bundle_no"] = bundleNo
	}

	total, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	skip := int64((page - 1) * pageSize)
	limit := int64(pageSize)
	opts := options.Find().SetSkip(skip).SetLimit(limit).SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var batches []*models.CuttingBatch
	if err = cursor.All(ctx, &batches); err != nil {
		return nil, 0, err
	}

	return batches, total, nil
}

// ==================== 裁片监控仓库实现 ====================

type cuttingPieceRepository struct {
	dbManager *database.DatabaseManager
}

func NewCuttingPieceRepository() CuttingPieceRepository {
	return &cuttingPieceRepository{
		dbManager: database.GetDatabaseManager(),
	}
}

func (r *cuttingPieceRepository) GetCollectionWithContext(ctx context.Context) *mongo.Collection {
	tenantID := tenantCtx.GetTenantID(ctx)
	db := r.dbManager.GetDatabase(tenantID)
	return db.Collection(models.CuttingPiece{}.TableName())
}

func (r *cuttingPieceRepository) Create(ctx context.Context, piece *models.CuttingPiece) error {
	collection := r.GetCollectionWithContext(ctx)
	_, err := collection.InsertOne(ctx, piece)
	return err
}

func (r *cuttingPieceRepository) Update(ctx context.Context, id string, piece *models.CuttingPiece) error {
	collection := r.GetCollectionWithContext(ctx)
	_, err := collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": piece})
	return err
}

func (r *cuttingPieceRepository) GetByID(ctx context.Context, id string) (*models.CuttingPiece, error) {
	collection := r.GetCollectionWithContext(ctx)
	var piece models.CuttingPiece
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&piece)
	if err != nil {
		return nil, err
	}
	return &piece, nil
}

func (r *cuttingPieceRepository) List(ctx context.Context, page, pageSize int, orderID, contractNo, bedNo, bundleNo string) ([]*models.CuttingPiece, int64, error) {
	collection := r.GetCollectionWithContext(ctx)

	filter := bson.M{}
	if orderID != "" {
		filter["order_id"] = orderID
	}
	if contractNo != "" {
		filter["contract_no"] = bson.M{"$regex": contractNo, "$options": "i"}
	}
	if bedNo != "" {
		filter["bed_no"] = bedNo
	}
	if bundleNo != "" {
		filter["bundle_no"] = bundleNo
	}

	total, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	skip := int64((page - 1) * pageSize)
	limit := int64(pageSize)
	opts := options.Find().SetSkip(skip).SetLimit(limit).SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var pieces []*models.CuttingPiece
	if err = cursor.All(ctx, &pieces); err != nil {
		return nil, 0, err
	}

	return pieces, total, nil
}
