package repository

import (
	"context"
	"mule-cloud/core/database"
	tenantCtx "mule-cloud/core/context"
	"mule-cloud/internal/models"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// StyleRepository 款式数据仓库接口
type StyleRepository interface {
	Get(ctx context.Context, id string) (*models.Style, error)
	Create(ctx context.Context, style *models.Style) error
	Update(ctx context.Context, id string, update bson.M) error
	Delete(ctx context.Context, id string) error
	Count(ctx context.Context, filter bson.M) (int64, error)
	GetCollectionWithContext(ctx context.Context) *mongo.Collection
	GetByStyleNo(ctx context.Context, styleNo string) (*models.Style, error)
}

type styleRepository struct {
	dbManager *database.DatabaseManager
}

// NewStyleRepository 创建款式数据仓库
func NewStyleRepository() StyleRepository {
	return &styleRepository{
		dbManager: database.GetDatabaseManager(),
	}
}

// GetCollectionWithContext 获取集合（支持租户上下文）
func (r *styleRepository) GetCollectionWithContext(ctx context.Context) *mongo.Collection {
	tenantID := tenantCtx.GetTenantID(ctx)
	db := r.dbManager.GetDatabase(tenantID)
	return db.Collection(models.Style{}.TableName())
}

// Get 获取款式
func (r *styleRepository) Get(ctx context.Context, id string) (*models.Style, error) {
	collection := r.GetCollectionWithContext(ctx)
	
	var style models.Style
	err := collection.FindOne(ctx, bson.M{"_id": id, "is_deleted": 0}).Decode(&style)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrNotFound
		}
		return nil, err
	}
	
	return &style, nil
}

// GetByStyleNo 根据款式编号获取款式
func (r *styleRepository) GetByStyleNo(ctx context.Context, styleNo string) (*models.Style, error) {
	collection := r.GetCollectionWithContext(ctx)
	
	var style models.Style
	err := collection.FindOne(ctx, bson.M{"style_no": styleNo, "is_deleted": 0}).Decode(&style)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrNotFound
		}
		return nil, err
	}
	
	return &style, nil
}

// Create 创建款式
func (r *styleRepository) Create(ctx context.Context, style *models.Style) error {
	collection := r.GetCollectionWithContext(ctx)
	
	result, err := collection.InsertOne(ctx, style)
	if err != nil {
		return err
	}
	
	style.ID = result.InsertedID.(string)
	return nil
}

// Update 更新款式
func (r *styleRepository) Update(ctx context.Context, id string, update bson.M) error {
	collection := r.GetCollectionWithContext(ctx)
	
	_, err := collection.UpdateOne(
		ctx,
		bson.M{"_id": id, "is_deleted": 0},
		bson.M{"$set": update},
	)
	
	return err
}

// Delete 删除款式（软删除）
func (r *styleRepository) Delete(ctx context.Context, id string) error {
	collection := r.GetCollectionWithContext(ctx)
	
	now := time.Now().Unix()
	_, err := collection.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.M{"$set": bson.M{"is_deleted": 1, "deleted_at": now}},
	)
	
	return err
}

// Count 统计款式数量
func (r *styleRepository) Count(ctx context.Context, filter bson.M) (int64, error) {
	collection := r.GetCollectionWithContext(ctx)
	return collection.CountDocuments(ctx, filter)
}
