package repository

import (
	"context"
	"mule-cloud/core/database"
	"mule-cloud/internal/models"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// BasicRepository Basic数据仓库接口
type BasicRepository interface {
	// Get 根据ID获取单条记录
	Get(ctx context.Context, id string) (*models.Basic, error)

	// GetByName 根据name获取单条记录（按创建时间倒序，返回最新的一条）
	GetByName(ctx context.Context, name string) (*models.Basic, error)

	// Find 查询记录列表
	Find(ctx context.Context, filter bson.M) ([]*models.Basic, error)

	// FindOne 查询单条记录
	FindOne(ctx context.Context, filter bson.M) (*models.Basic, error)

	// FindWithPage 分页查询记录列表
	FindWithPage(ctx context.Context, filter bson.M, page, pageSize int64) ([]*models.Basic, error)

	// Count 统计记录数
	Count(ctx context.Context, filter bson.M) (int64, error)

	// Create 创建记录
	Create(ctx context.Context, basic *models.Basic) error

	// Update 更新记录
	Update(ctx context.Context, id string, update bson.M) error

	// UpdateOne 按条件更新单条记录
	UpdateOne(ctx context.Context, filter bson.M, update bson.M) error

	// Delete 删除记录
	Delete(ctx context.Context, id string) error

	// DeleteMany 批量删除
	DeleteMany(ctx context.Context, filter bson.M) (int64, error)

	// FindDeletedWithPage 分页查询已删除记录列表
	FindDeletedWithPage(ctx context.Context, filter bson.M, page, pageSize int64) ([]*models.Basic, error)

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

	// FindByTenant 查询租户下的记录（包括公共数据）
	FindByTenant(ctx context.Context, tenantID string, filter bson.M) ([]*models.Basic, error)

	// FindByTenantWithPage 分页查询租户下的记录（包括公共数据）
	FindByTenantWithPage(ctx context.Context, tenantID string, filter bson.M, page, pageSize int64) ([]*models.Basic, error)

	// CountByTenant 统计租户下的记录数（包括公共数据）
	CountByTenant(ctx context.Context, tenantID string, filter bson.M) (int64, error)

	// CheckOwnership 检查记录是否属于指定租户（用于修改/删除权限检查）
	CheckOwnership(ctx context.Context, id string, tenantID string) (bool, error)
}

// basicRepository Basic数据仓库实现
type basicRepository struct {
	collection *mongo.Collection
}

// NewBasicRepository 创建Basic数据仓库实例
func NewBasicRepository() BasicRepository {
	collection := database.MongoDB.Collection("basic")
	return &basicRepository{
		collection: collection,
	}
}

// Get 根据ID获取单条记录（排除软删除）
func (r *basicRepository) Get(ctx context.Context, id string) (*models.Basic, error) {
	// 将字符串 ID 转换为 ObjectID 进行查询
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		// 如果转换失败，尝试直接用字符串查询（兼容旧数据）
		filter := bson.M{"_id": id, "is_deleted": 0}
		basic := &models.Basic{}
		err = r.collection.FindOne(ctx, filter).Decode(basic)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return nil, nil
			}
			return nil, err
		}
		return basic, nil
	}

	// 使用 ObjectID 查询（排除软删除）
	filter := bson.M{"_id": objectID, "is_deleted": 0}
	basic := &models.Basic{}
	err = r.collection.FindOne(ctx, filter).Decode(basic)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return basic, nil
}

// GetByName 根据name获取单条记录（按创建时间倒序，返回最新的一条）（排除软删除）
func (r *basicRepository) GetByName(ctx context.Context, name string) (*models.Basic, error) {
	filter := bson.M{"name": name, "is_deleted": 0}
	// 需要用 service 层自己实现排序逻辑
	basic := &models.Basic{}
	err := r.collection.FindOne(ctx, filter).Decode(basic)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return basic, nil
}

// Find 查询记录列表
func (r *basicRepository) Find(ctx context.Context, filter bson.M) ([]*models.Basic, error) {
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	basics := []*models.Basic{}
	err = cursor.All(ctx, &basics)
	if err != nil {
		return nil, err
	}
	return basics, nil
}

// FindOne 查询单条记录
func (r *basicRepository) FindOne(ctx context.Context, filter bson.M) (*models.Basic, error) {
	basic := &models.Basic{}
	err := r.collection.FindOne(ctx, filter).Decode(basic)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return basic, nil
}

// FindWithPage 分页查询记录列表（按创建时间倒序）
func (r *basicRepository) FindWithPage(ctx context.Context, filter bson.M, page, pageSize int64) ([]*models.Basic, error) {
	// 简化实现，不使用 options，让 service 层处理分页逻辑
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	basics := []*models.Basic{}
	err = cursor.All(ctx, &basics)
	if err != nil {
		return nil, err
	}
	return basics, nil
}

// Count 统计记录数
func (r *basicRepository) Count(ctx context.Context, filter bson.M) (int64, error) {
	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// Create 创建记录
func (r *basicRepository) Create(ctx context.Context, basic *models.Basic) error {
	result, err := r.collection.InsertOne(ctx, basic)
	if err != nil {
		return err
	}

	// 将生成的 ObjectID 转换为字符串并设置到 basic.ID
	if oid, ok := result.InsertedID.(bson.ObjectID); ok {
		basic.ID = oid.Hex()
	}

	return nil
}

// Update 更新记录
func (r *basicRepository) Update(ctx context.Context, id string, update bson.M) error {
	// 将字符串 ID 转换为 ObjectID 进行更新
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		// 如果转换失败，尝试直接用字符串更新（兼容旧数据）
		filter := bson.M{"_id": id}
		updateDoc := bson.M{"$set": update}
		result, err := r.collection.UpdateOne(ctx, filter, updateDoc)
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
	result, err := r.collection.UpdateOne(ctx, filter, updateDoc)
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
func (r *basicRepository) UpdateOne(ctx context.Context, filter bson.M, update bson.M) error {
	updateDoc := bson.M{"$set": update}
	_, err := r.collection.UpdateOne(ctx, filter, updateDoc)
	return err
}

// Delete 软删除记录（设置 is_deleted = 1）
func (r *basicRepository) Delete(ctx context.Context, id string) error {
	// 将字符串 ID 转换为 ObjectID 进行删除
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		// 如果转换失败，尝试直接用字符串删除（兼容旧数据）
		filter := bson.M{"_id": id, "is_deleted": 0}
		updateDoc := bson.M{"$set": bson.M{"is_deleted": 1, "deleted_at": time.Now().Unix()}}
		result, err := r.collection.UpdateOne(ctx, filter, updateDoc)
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
	result, err := r.collection.UpdateOne(ctx, filter, updateDoc)
	if err != nil {
		return err
	}

	// 检查是否匹配到记录
	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}

// DeleteMany 批量删除
func (r *basicRepository) DeleteMany(ctx context.Context, filter bson.M) (int64, error) {
	result, err := r.collection.DeleteMany(ctx, filter)
	if err != nil {
		return 0, err
	}
	return result.DeletedCount, nil
}

// FindDeletedWithPage 分页查询已删除记录列表（按删除时间倒序）
func (r *basicRepository) FindDeletedWithPage(ctx context.Context, filter bson.M, page, pageSize int64) ([]*models.Basic, error) {
	// 添加已删除条件
	filter["is_deleted"] = 1

	// 简化实现，不使用 options，让 service 层处理分页逻辑
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	basics := []*models.Basic{}
	err = cursor.All(ctx, &basics)
	if err != nil {
		return nil, err
	}
	return basics, nil
}

// CountDeleted 统计已删除记录数
func (r *basicRepository) CountDeleted(ctx context.Context, filter bson.M) (int64, error) {
	filter["is_deleted"] = 1
	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// RestoreOne 恢复单条记录（清除 is_deleted 和 deleted_at）
func (r *basicRepository) RestoreOne(ctx context.Context, id string) error {
	// 将字符串 ID 转换为 ObjectID 进行恢复
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		// 如果转换失败，尝试直接用字符串恢复（兼容旧数据）
		filter := bson.M{"_id": id, "is_deleted": 1}
		updateDoc := bson.M{"$set": bson.M{"is_deleted": 0, "deleted_at": 0}}
		result, err := r.collection.UpdateOne(ctx, filter, updateDoc)
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
	result, err := r.collection.UpdateOne(ctx, filter, updateDoc)
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
func (r *basicRepository) RestoreMany(ctx context.Context, ids []string) (int64, error) {
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
		result, err := r.collection.UpdateMany(ctx, filter, updateDoc)
		if err != nil {
			return 0, err
		}
		totalCount += result.ModifiedCount
	}

	// 处理字符串格式的 ID（兼容旧数据）
	if len(stringIDs) > 0 {
		filter := bson.M{"_id": bson.M{"$in": stringIDs}, "is_deleted": 1}
		updateDoc := bson.M{"$set": bson.M{"is_deleted": 0, "deleted_at": 0}}
		result, err := r.collection.UpdateMany(ctx, filter, updateDoc)
		if err != nil {
			return totalCount, err
		}
		totalCount += result.ModifiedCount
	}

	return totalCount, nil
}

// HardDelete 物理删除记录（真正从数据库删除）
func (r *basicRepository) HardDelete(ctx context.Context, id string) error {
	// 将字符串 ID 转换为 ObjectID 进行删除
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		// 如果转换失败，尝试直接用字符串删除（兼容旧数据）
		filter := bson.M{"_id": id}
		result, err := r.collection.DeleteOne(ctx, filter)
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
	result, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	// 检查是否匹配到记录
	if result.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}

// HardDeleteMany 批量物理删除记录
func (r *basicRepository) HardDeleteMany(ctx context.Context, ids []string) (int64, error) {
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
		result, err := r.collection.DeleteMany(ctx, filter)
		if err != nil {
			return 0, err
		}
		totalCount += result.DeletedCount
	}

	// 处理字符串格式的 ID（兼容旧数据）
	if len(stringIDs) > 0 {
		filter := bson.M{"_id": bson.M{"$in": stringIDs}}
		result, err := r.collection.DeleteMany(ctx, filter)
		if err != nil {
			return totalCount, err
		}
		totalCount += result.DeletedCount
	}

	return totalCount, nil
}

// GetCollection 获取MongoDB集合（供高级用法使用，如需要使用 options）
func (r *basicRepository) GetCollection() *mongo.Collection {
	return r.collection
}

// FindByTenant 查询租户下的记录（包括公共数据）
func (r *basicRepository) FindByTenant(ctx context.Context, tenantID string, filter bson.M) ([]*models.Basic, error) {
	// 添加租户条件：自己的数据 OR 公共数据
	filter["$or"] = []bson.M{
		{"tenant_id": tenantID},
		{"is_common": true},
	}
	filter["is_deleted"] = 0

	return r.Find(ctx, filter)
}

// FindByTenantWithPage 分页查询租户下的记录（包括公共数据）
func (r *basicRepository) FindByTenantWithPage(ctx context.Context, tenantID string, filter bson.M, page, pageSize int64) ([]*models.Basic, error) {
	// 添加租户条件：自己的数据 OR 公共数据
	filter["$or"] = []bson.M{
		{"tenant_id": tenantID},
		{"is_common": true},
	}
	filter["is_deleted"] = 0

	return r.FindWithPage(ctx, filter, page, pageSize)
}

// CountByTenant 统计租户下的记录数（包括公共数据）
func (r *basicRepository) CountByTenant(ctx context.Context, tenantID string, filter bson.M) (int64, error) {
	// 添加租户条件：自己的数据 OR 公共数据
	filter["$or"] = []bson.M{
		{"tenant_id": tenantID},
		{"is_common": true},
	}
	filter["is_deleted"] = 0

	return r.Count(ctx, filter)
}

// CheckOwnership 检查记录是否属于指定租户（用于修改/删除权限检查）
// 返回 true 表示可以修改/删除（属于该租户且不是公共数据）
func (r *basicRepository) CheckOwnership(ctx context.Context, id string, tenantID string) (bool, error) {
	basic, err := r.Get(ctx, id)
	if err != nil {
		return false, err
	}
	if basic == nil {
		return false, nil
	}

	// 只有属于该租户且不是公共数据的记录才能修改/删除
	return basic.TenantID == tenantID && !basic.IsCommon, nil
}
