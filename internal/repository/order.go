package repository

import (
	"context"
	tenantCtx "mule-cloud/core/context"
	"mule-cloud/core/database"
	"mule-cloud/internal/models"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// OrderRepository 订单数据仓库接口
type OrderRepository interface {
	Get(ctx context.Context, id string) (*models.Order, error)
	Create(ctx context.Context, order *models.Order) error
	Update(ctx context.Context, id string, update bson.M) error
	Delete(ctx context.Context, id string) error
	Count(ctx context.Context, filter bson.M) (int64, error)
	GetCollectionWithContext(ctx context.Context) *mongo.Collection
	GetByContractNo(ctx context.Context, contractNo string) (*models.Order, error)
}

type orderRepository struct {
	dbManager *database.DatabaseManager
}

// NewOrderRepository 创建订单数据仓库
func NewOrderRepository() OrderRepository {
	return &orderRepository{
		dbManager: database.GetDatabaseManager(),
	}
}

// GetCollectionWithContext 获取集合（支持租户上下文）
func (r *orderRepository) GetCollectionWithContext(ctx context.Context) *mongo.Collection {
	tenantCode := tenantCtx.GetTenantCode(ctx)
	db := r.dbManager.GetDatabase(tenantCode)
	return db.Collection(models.Order{}.TableName())
}

// Get 获取订单
func (r *orderRepository) Get(ctx context.Context, id string) (*models.Order, error) {
	collection := r.GetCollectionWithContext(ctx)

	// 将字符串 ID 转换为 ObjectID
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var order models.Order
	err = collection.FindOne(ctx, bson.M{"_id": objectID, "is_deleted": 0}).Decode(&order)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &order, nil
}

// GetByContractNo 根据合同号获取订单
func (r *orderRepository) GetByContractNo(ctx context.Context, contractNo string) (*models.Order, error) {
	collection := r.GetCollectionWithContext(ctx)

	var order models.Order
	err := collection.FindOne(ctx, bson.M{"contract_no": contractNo, "is_deleted": 0}).Decode(&order)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &order, nil
}

// Create 创建订单
func (r *orderRepository) Create(ctx context.Context, order *models.Order) error {
	collection := r.GetCollectionWithContext(ctx)

	result, err := collection.InsertOne(ctx, order)
	if err != nil {
		return err
	}

	// 将 bson.ObjectID 转换为字符串
	if oid, ok := result.InsertedID.(bson.ObjectID); ok {
		order.ID = oid.Hex()
	}
	return nil
}

// Update 更新订单
func (r *orderRepository) Update(ctx context.Context, id string, update bson.M) error {
	collection := r.GetCollectionWithContext(ctx)

	// 将字符串 ID 转换为 ObjectID
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = collection.UpdateOne(
		ctx,
		bson.M{"_id": objectID, "is_deleted": 0},
		bson.M{"$set": update},
	)

	return err
}

// Delete 删除订单（软删除）
func (r *orderRepository) Delete(ctx context.Context, id string) error {
	collection := r.GetCollectionWithContext(ctx)

	// 将字符串 ID 转换为 ObjectID
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	now := time.Now().Unix()
	_, err = collection.UpdateOne(
		ctx,
		bson.M{"_id": objectID},
		bson.M{"$set": bson.M{"is_deleted": 1, "deleted_at": now}},
	)

	return err
}

// Count 统计订单数量
func (r *orderRepository) Count(ctx context.Context, filter bson.M) (int64, error) {
	collection := r.GetCollectionWithContext(ctx)
	return collection.CountDocuments(ctx, filter)
}
