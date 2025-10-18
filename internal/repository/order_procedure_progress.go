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

// OrderProcedureProgressRepository 订单工序进度仓储接口
type OrderProcedureProgressRepository interface {
	Create(ctx context.Context, progress *models.OrderProcedureProgress) error
	GetByOrderAndProcedure(ctx context.Context, orderID string, procedureSeq int) (*models.OrderProcedureProgress, error)
	ListByOrder(ctx context.Context, orderID string) ([]*models.OrderProcedureProgress, error)
	UpdateReportedQty(ctx context.Context, orderID string, procedureSeq int, quantity int) error
	InitOrderProgress(ctx context.Context, orderID, contractNo string, totalQty int, procedures []models.OrderProcedure) error
	GetOrderOverallProgress(ctx context.Context, orderID string) (float64, error)
}

type orderProcedureProgressRepository struct {
	dbManager *database.DatabaseManager
}

// NewOrderProcedureProgressRepository 创建订单工序进度仓储
func NewOrderProcedureProgressRepository() OrderProcedureProgressRepository {
	return &orderProcedureProgressRepository{
		dbManager: database.GetDatabaseManager(),
	}
}

// GetCollectionWithContext 获取集合（支持租户上下文）
func (r *orderProcedureProgressRepository) GetCollectionWithContext(ctx context.Context) *mongo.Collection {
	tenantCode := tenantCtx.GetTenantCode(ctx)
	db := r.dbManager.GetDatabase(tenantCode)
	return db.Collection(models.OrderProcedureProgress{}.TableName())
}

// Create 创建订单工序进度
func (r *orderProcedureProgressRepository) Create(ctx context.Context, progress *models.OrderProcedureProgress) error {
	collection := r.GetCollectionWithContext(ctx)
	_, err := collection.InsertOne(ctx, progress)
	return err
}

// GetByOrderAndProcedure 根据订单ID和工序序号获取进度
func (r *orderProcedureProgressRepository) GetByOrderAndProcedure(ctx context.Context, orderID string, procedureSeq int) (*models.OrderProcedureProgress, error) {
	collection := r.GetCollectionWithContext(ctx)

	var progress models.OrderProcedureProgress
	err := collection.FindOne(ctx, bson.M{
		"order_id":      orderID,
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

// ListByOrder 获取订单的所有工序进度
func (r *orderProcedureProgressRepository) ListByOrder(ctx context.Context, orderID string) ([]*models.OrderProcedureProgress, error) {
	collection := r.GetCollectionWithContext(ctx)

	cursor, err := collection.Find(ctx, bson.M{"order_id": orderID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var progressList []*models.OrderProcedureProgress
	if err = cursor.All(ctx, &progressList); err != nil {
		return nil, err
	}

	return progressList, nil
}

// UpdateReportedQty 更新已上报数量
func (r *orderProcedureProgressRepository) UpdateReportedQty(ctx context.Context, orderID string, procedureSeq int, quantity int) error {
	collection := r.GetCollectionWithContext(ctx)

	filter := bson.M{
		"order_id":      orderID,
		"procedure_seq": procedureSeq,
	}

	// 先查询当前进度
	progress, err := r.GetByOrderAndProcedure(ctx, orderID, procedureSeq)
	if err != nil {
		return err
	}

	// 更新已上报数量和进度百分比
	newReportedQty := progress.ReportedQty + quantity
	newProgress := float64(newReportedQty) / float64(progress.TotalQty) * 100
	if newProgress > 100 {
		newProgress = 100 // 防止超过100%
	}

	update := bson.M{
		"$set": bson.M{
			"reported_qty": newReportedQty,
			"progress":     newProgress,
			"updated_at":   time.Now().Unix(),
		},
	}

	_, err = collection.UpdateOne(ctx, filter, update)
	return err
}

// InitOrderProgress 初始化订单的所有工序进度
func (r *orderProcedureProgressRepository) InitOrderProgress(ctx context.Context, orderID, contractNo string, totalQty int, procedures []models.OrderProcedure) error {
	collection := r.GetCollectionWithContext(ctx)

	// 检查是否已经初始化
	existing, _ := r.ListByOrder(ctx, orderID)
	if len(existing) > 0 {
		return nil // 已经初始化过了
	}

	// 批量插入所有工序的进度记录
	var documents []interface{}
	now := time.Now().Unix()

	for _, proc := range procedures {
		progress := &models.OrderProcedureProgress{
			ID:            bson.NewObjectID().Hex(),
			OrderID:       orderID,
			ContractNo:    contractNo,
			ProcedureSeq:  proc.Sequence,
			ProcedureName: proc.ProcedureName,
			TotalQty:      totalQty,
			ReportedQty:   0,
			Progress:      0,
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

// GetOrderOverallProgress 获取订单的总体进度（所有工序的平均进度）
func (r *orderProcedureProgressRepository) GetOrderOverallProgress(ctx context.Context, orderID string) (float64, error) {
	collection := r.GetCollectionWithContext(ctx)

	pipeline := []bson.D{
		{{Key: "$match", Value: bson.M{"order_id": orderID}}},
		{{Key: "$group", Value: bson.M{
			"_id":              nil,
			"average_progress": bson.M{"$avg": "$progress"},
		}}},
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return 0, err
	}
	defer cursor.Close(ctx)

	var result []struct {
		AverageProgress float64 `bson:"average_progress"`
	}
	if err = cursor.All(ctx, &result); err != nil {
		return 0, err
	}

	if len(result) > 0 {
		return result[0].AverageProgress, nil
	}

	return 0, nil
}
