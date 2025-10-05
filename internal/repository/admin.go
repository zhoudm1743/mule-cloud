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

// AdminRepository Admin数据仓库接口
type AdminRepository interface {
	// Get 根据ID获取单条记录
	Get(ctx context.Context, id string) (*models.Admin, error)

	// GetByPhone 根据手机号获取单条记录
	GetByPhone(ctx context.Context, phone string) (*models.Admin, error)

	// GetByEmail 根据邮箱获取单条记录
	GetByEmail(ctx context.Context, email string) (*models.Admin, error)

	// Find 查询记录列表
	Find(ctx context.Context, filter bson.M) ([]*models.Admin, error)

	// FindOne 查询单条记录
	FindOne(ctx context.Context, filter bson.M) (*models.Admin, error)

	// FindWithPage 分页查询记录列表
	FindWithPage(ctx context.Context, filter bson.M, page, pageSize int64) ([]*models.Admin, error)

	// Count 统计记录数
	Count(ctx context.Context, filter bson.M) (int64, error)

	// Create 创建记录
	Create(ctx context.Context, admin *models.Admin) error

	// Update 更新记录
	Update(ctx context.Context, id string, update bson.M) error

	// UpdateOne 按条件更新单条记录
	UpdateOne(ctx context.Context, filter bson.M, update bson.M) error

	// Delete 删除记录
	Delete(ctx context.Context, id string) error

	// DeleteMany 批量删除
	DeleteMany(ctx context.Context, filter bson.M) (int64, error)

	// FindDeletedWithPage 分页查询已删除记录列表
	FindDeletedWithPage(ctx context.Context, filter bson.M, page, pageSize int64) ([]*models.Admin, error)

	// CountDeleted 统计已删除记录数
	CountDeleted(ctx context.Context, filter bson.M) (int64, error)

	// RestoreOne 恢复单条记录
	RestoreOne(ctx context.Context, id string) error

	// RestoreMany 批量恢复记录
	RestoreMany(ctx context.Context, ids []string) (int64, error)

	// HardDelete 物理删除记录
	HardDelete(ctx context.Context, id string) error

	// HardDeleteMany 批量物理删除记录
	HardDeleteMany(ctx context.Context, ids []string) (int64, error)

	// GetCollection 获取MongoDB集合（供高级用法使用）
	GetCollection() *mongo.Collection
}

// adminRepository Admin数据仓库实现
type adminRepository struct {
	dbManager *database.DatabaseManager
}

// NewAdminRepository 创建Admin数据仓库实例
func NewAdminRepository() AdminRepository {
	return &adminRepository{
		dbManager: database.GetDatabaseManager(),
	}
}

// getCollection 获取集合（自动根据Context中的租户ID切换数据库）
func (r *adminRepository) getCollection(ctx context.Context) *mongo.Collection {
	tenantID := tenantCtx.GetTenantID(ctx)
	db := r.dbManager.GetDatabase(tenantID)
	return db.Collection("admin")
}

// Get 根据ID获取单条记录（排除软删除）
func (r *adminRepository) Get(ctx context.Context, id string) (*models.Admin, error) {
	collection := r.getCollection(ctx)

	// 将字符串 ID 转换为 ObjectID 进行查询
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		// 如果转换失败，尝试直接用字符串查询（兼容旧数据）
		filter := bson.M{"_id": id, "is_deleted": 0}
		admin := &models.Admin{}
		err = collection.FindOne(ctx, filter).Decode(admin)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return nil, nil
			}
			return nil, err
		}
		return admin, nil
	}

	// 使用 ObjectID 查询（排除软删除）
	filter := bson.M{"_id": objectID, "is_deleted": 0}
	admin := &models.Admin{}
	err = collection.FindOne(ctx, filter).Decode(admin)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return admin, nil
}

// GetByPhone 根据手机号获取单条记录（排除软删除）
func (r *adminRepository) GetByPhone(ctx context.Context, phone string) (*models.Admin, error) {
	collection := r.getCollection(ctx)
	filter := bson.M{"phone": phone, "is_deleted": 0}
	admin := &models.Admin{}
	err := collection.FindOne(ctx, filter).Decode(admin)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return admin, nil
}

// GetByEmail 根据邮箱获取单条记录（排除软删除）
func (r *adminRepository) GetByEmail(ctx context.Context, email string) (*models.Admin, error) {
	collection := r.getCollection(ctx)
	filter := bson.M{"email": email, "is_deleted": 0}
	admin := &models.Admin{}
	err := collection.FindOne(ctx, filter).Decode(admin)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return admin, nil
}

// Find 查询记录列表
func (r *adminRepository) Find(ctx context.Context, filter bson.M) ([]*models.Admin, error) {
	collection := r.getCollection(ctx)
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	admins := []*models.Admin{}
	err = cursor.All(ctx, &admins)
	if err != nil {
		return nil, err
	}
	return admins, nil
}

// FindOne 查询单条记录
func (r *adminRepository) FindOne(ctx context.Context, filter bson.M) (*models.Admin, error) {
	collection := r.getCollection(ctx)
	admin := &models.Admin{}
	err := collection.FindOne(ctx, filter).Decode(admin)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return admin, nil
}

// FindWithPage 分页查询记录列表（按创建时间倒序）
func (r *adminRepository) FindWithPage(ctx context.Context, filter bson.M, page, pageSize int64) ([]*models.Admin, error) {
	collection := r.getCollection(ctx)
	// 简化实现，不使用 options，让 service 层处理分页逻辑
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	admins := []*models.Admin{}
	err = cursor.All(ctx, &admins)
	if err != nil {
		return nil, err
	}
	return admins, nil
}

// Count 统计记录数
func (r *adminRepository) Count(ctx context.Context, filter bson.M) (int64, error) {
	collection := r.getCollection(ctx)
	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// Create 创建记录
func (r *adminRepository) Create(ctx context.Context, admin *models.Admin) error {
	collection := r.getCollection(ctx)
	result, err := collection.InsertOne(ctx, admin)
	if err != nil {
		return err
	}

	// 将生成的 ObjectID 转换为字符串并设置到 admin.ID
	if oid, ok := result.InsertedID.(bson.ObjectID); ok {
		admin.ID = oid.Hex()
	}

	return nil
}

// Update 更新记录
func (r *adminRepository) Update(ctx context.Context, id string, update bson.M) error {
	collection := r.getCollection(ctx)
	// 将字符串 ID 转换为 ObjectID 进行更新
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		// 如果转换失败，尝试直接用字符串更新（兼容旧数据）
		filter := bson.M{"_id": id}
		updateDoc := bson.M{"$set": update}
		result, err := collection.UpdateOne(ctx, filter, updateDoc)
		if err != nil {
			return err
		}
		if result.MatchedCount == 0 {
			return mongo.ErrNoDocuments
		}
		return nil
	}

	// 使用 ObjectID 更新
	filter := bson.M{"_id": objectID}
	updateDoc := bson.M{"$set": update}
	result, err := collection.UpdateOne(ctx, filter, updateDoc)
	if err != nil {
		return err
	}

	// 检查是否匹配到记录
	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}

// UpdateOne 按条件更新单条记录
func (r *adminRepository) UpdateOne(ctx context.Context, filter bson.M, update bson.M) error {
	collection := r.getCollection(ctx)
	updateDoc := bson.M{"$set": update}
	_, err := collection.UpdateOne(ctx, filter, updateDoc)
	return err
}

// Delete 软删除记录（设置 is_deleted = 1）
func (r *adminRepository) Delete(ctx context.Context, id string) error {
	collection := r.getCollection(ctx)
	// 将字符串 ID 转换为 ObjectID 进行删除
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		// 如果转换失败，尝试直接用字符串删除（兼容旧数据）
		filter := bson.M{"_id": id, "is_deleted": 0}
		updateDoc := bson.M{"$set": bson.M{"is_deleted": 1, "deleted_at": time.Now().Unix()}}
		result, err := collection.UpdateOne(ctx, filter, updateDoc)
		if err != nil {
			return err
		}
		if result.MatchedCount == 0 {
			return mongo.ErrNoDocuments
		}
		return nil
	}

	// 使用 ObjectID 软删除
	filter := bson.M{"_id": objectID, "is_deleted": 0}
	updateDoc := bson.M{"$set": bson.M{"is_deleted": 1, "deleted_at": time.Now().Unix()}}
	result, err := collection.UpdateOne(ctx, filter, updateDoc)
	if err != nil {
		return err
	}

	// 检查是否匹配到记录
	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}

// HardDelete 物理删除记录（真正从数据库删除）
func (r *adminRepository) HardDelete(ctx context.Context, id string) error {
	collection := r.getCollection(ctx)
	// 将字符串 ID 转换为 ObjectID 进行删除
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		// 如果转换失败，尝试直接用字符串删除（兼容旧数据）
		filter := bson.M{"_id": id}
		result, err := collection.DeleteOne(ctx, filter)
		if err != nil {
			return err
		}
		if result.DeletedCount == 0 {
			return mongo.ErrNoDocuments
		}
		return nil
	}

	// 使用 ObjectID 物理删除
	filter := bson.M{"_id": objectID}
	result, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	// 检查是否匹配到记录
	if result.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}

// DeleteMany 批量删除
func (r *adminRepository) DeleteMany(ctx context.Context, filter bson.M) (int64, error) {
	collection := r.getCollection(ctx)
	result, err := collection.DeleteMany(ctx, filter)
	if err != nil {
		return 0, err
	}
	return result.DeletedCount, nil
}

// FindDeletedWithPage 分页查询已删除记录列表（按删除时间倒序）
func (r *adminRepository) FindDeletedWithPage(ctx context.Context, filter bson.M, page, pageSize int64) ([]*models.Admin, error) {
	collection := r.getCollection(ctx)
	// 添加已删除条件
	filter["is_deleted"] = 1

	// 简化实现，不使用 options，让 service 层处理分页逻辑
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	admins := []*models.Admin{}
	err = cursor.All(ctx, &admins)
	if err != nil {
		return nil, err
	}
	return admins, nil
}

// CountDeleted 统计已删除记录数
func (r *adminRepository) CountDeleted(ctx context.Context, filter bson.M) (int64, error) {
	collection := r.getCollection(ctx)
	filter["is_deleted"] = 1
	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// RestoreOne 恢复单条记录（清除 is_deleted 和 deleted_at）
func (r *adminRepository) RestoreOne(ctx context.Context, id string) error {
	collection := r.getCollection(ctx)
	// 将字符串 ID 转换为 ObjectID 进行恢复
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		// 如果转换失败，尝试直接用字符串恢复（兼容旧数据）
		filter := bson.M{"_id": id, "is_deleted": 1}
		updateDoc := bson.M{"$set": bson.M{"is_deleted": 0, "deleted_at": 0}}
		result, err := collection.UpdateOne(ctx, filter, updateDoc)
		if err != nil {
			return err
		}
		if result.MatchedCount == 0 {
			return mongo.ErrNoDocuments
		}
		return nil
	}

	// 使用 ObjectID 恢复
	filter := bson.M{"_id": objectID, "is_deleted": 1}
	updateDoc := bson.M{"$set": bson.M{"is_deleted": 0, "deleted_at": 0}}
	result, err := collection.UpdateOne(ctx, filter, updateDoc)
	if err != nil {
		return err
	}

	// 检查是否匹配到记录
	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}

// RestoreMany 批量恢复记录
func (r *adminRepository) RestoreMany(ctx context.Context, ids []string) (int64, error) {
	collection := r.getCollection(ctx)
	if len(ids) == 0 {
		return 0, nil
	}

	// 转换 ID 列表
	var objectIDs []interface{}
	var stringIDs []interface{}

	for _, id := range ids {
		if objectID, err := bson.ObjectIDFromHex(id); err == nil {
			objectIDs = append(objectIDs, objectID)
		} else {
			stringIDs = append(stringIDs, id)
		}
	}

	var totalCount int64

	// 处理 ObjectID 格式的 ID
	if len(objectIDs) > 0 {
		filter := bson.M{"_id": bson.M{"$in": objectIDs}, "is_deleted": 1}
		updateDoc := bson.M{"$set": bson.M{"is_deleted": 0, "deleted_at": 0}}
		result, err := collection.UpdateMany(ctx, filter, updateDoc)
		if err != nil {
			return 0, err
		}
		totalCount += result.ModifiedCount
	}

	// 处理字符串格式的 ID（兼容旧数据）
	if len(stringIDs) > 0 {
		filter := bson.M{"_id": bson.M{"$in": stringIDs}, "is_deleted": 1}
		updateDoc := bson.M{"$set": bson.M{"is_deleted": 0, "deleted_at": 0}}
		result, err := collection.UpdateMany(ctx, filter, updateDoc)
		if err != nil {
			return totalCount, err
		}
		totalCount += result.ModifiedCount
	}

	return totalCount, nil
}

// HardDeleteMany 批量物理删除记录
func (r *adminRepository) HardDeleteMany(ctx context.Context, ids []string) (int64, error) {
	collection := r.getCollection(ctx)
	if len(ids) == 0 {
		return 0, nil
	}

	// 转换 ID 列表
	var objectIDs []interface{}
	var stringIDs []interface{}

	for _, id := range ids {
		if objectID, err := bson.ObjectIDFromHex(id); err == nil {
			objectIDs = append(objectIDs, objectID)
		} else {
			stringIDs = append(stringIDs, id)
		}
	}

	var totalCount int64

	// 处理 ObjectID 格式的 ID
	if len(objectIDs) > 0 {
		filter := bson.M{"_id": bson.M{"$in": objectIDs}}
		result, err := collection.DeleteMany(ctx, filter)
		if err != nil {
			return 0, err
		}
		totalCount += result.DeletedCount
	}

	// 处理字符串格式的 ID（兼容旧数据）
	if len(stringIDs) > 0 {
		filter := bson.M{"_id": bson.M{"$in": stringIDs}}
		result, err := collection.DeleteMany(ctx, filter)
		if err != nil {
			return totalCount, err
		}
		totalCount += result.DeletedCount
	}

	return totalCount, nil
}

// GetCollection 获取MongoDB集合（供高级用法使用，如需要使用 options）
// 注意：此方法需要从Context获取租户ID，建议传入Context
func (r *adminRepository) GetCollection() *mongo.Collection {
	// 返回系统库的collection作为默认值（向下兼容）
	// 实际使用时应该调用 getCollection(ctx)
	return r.dbManager.GetSystemDatabase().Collection("admin")
}
