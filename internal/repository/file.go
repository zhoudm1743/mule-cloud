package repository

import (
	"context"
	"mule-cloud/internal/models"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// FileRepository 文件仓储接口
type FileRepository interface {
	Create(ctx context.Context, tenantCode string, file *models.FileInfo) error
	Get(ctx context.Context, tenantCode string, id string) (*models.FileInfo, error)
	List(ctx context.Context, tenantCode string, page, pageSize int, businessType string) ([]*models.FileInfo, int64, error)
	Delete(ctx context.Context, tenantCode string, id string) error
	GetByKey(ctx context.Context, tenantCode string, storageKey string) (*models.FileInfo, error)
}

// fileRepository 文件仓储实现
type fileRepository struct {
	db *mongo.Database
}

// NewFileRepository 创建文件仓储实例
func NewFileRepository(db *mongo.Database) FileRepository {
	return &fileRepository{db: db}
}

// getCollection 获取租户专属集合
func (r *fileRepository) getCollection(tenantCode string) *mongo.Collection {
	dbName := "mule_" + tenantCode
	db := r.db.Client().Database(dbName)
	return db.Collection(models.FileInfo{}.TableName())
}

// Create 创建文件记录
func (r *fileRepository) Create(ctx context.Context, tenantCode string, file *models.FileInfo) error {
	collection := r.getCollection(tenantCode)

	file.TenantCode = tenantCode
	file.CreatedAt = time.Now()
	file.UpdatedAt = time.Now()

	result, err := collection.InsertOne(ctx, file)
	if err != nil {
		return err
	}

	// 将 bson.ObjectID 转换为字符串
	if oid, ok := result.InsertedID.(bson.ObjectID); ok {
		file.ID = oid.Hex()
	}

	return nil
}

// Get 根据ID获取文件记录
func (r *fileRepository) Get(ctx context.Context, tenantCode string, id string) (*models.FileInfo, error) {
	collection := r.getCollection(tenantCode)

	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var file models.FileInfo
	err = collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&file)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, err
		}
		return nil, err
	}

	return &file, nil
}

// List 获取文件列表
func (r *fileRepository) List(ctx context.Context, tenantCode string, page, pageSize int, businessType string) ([]*models.FileInfo, int64, error) {
	collection := r.getCollection(tenantCode)

	// 构建查询条件
	filter := bson.M{}
	if businessType != "" {
		filter["business_type"] = businessType
	}

	// 计算总数
	total, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	// 查询列表
	skip := int64((page - 1) * pageSize)
	opts := options.Find().
		SetSkip(skip).
		SetLimit(int64(pageSize)).
		SetSort(bson.M{"created_at": -1})

	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var files []*models.FileInfo
	if err = cursor.All(ctx, &files); err != nil {
		return nil, 0, err
	}

	return files, total, nil
}

// Delete 删除文件记录
func (r *fileRepository) Delete(ctx context.Context, tenantCode string, id string) error {
	collection := r.getCollection(tenantCode)

	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	result, err := collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}

// GetByKey 根据存储键获取文件记录
func (r *fileRepository) GetByKey(ctx context.Context, tenantCode string, storageKey string) (*models.FileInfo, error) {
	collection := r.getCollection(tenantCode)

	var file models.FileInfo
	err := collection.FindOne(ctx, bson.M{"storage_key": storageKey}).Decode(&file)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &file, nil
}
