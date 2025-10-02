package repository

import (
	"context"
	"mule-cloud/core/database"
	"mule-cloud/internal/models"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// RoleRepository Role数据仓库接口
type RoleRepository interface {
	// Get 根据ID获取单条记录
	Get(ctx context.Context, id string) (*models.Role, error)

	// GetByCode 根据角色代码获取单条记录
	GetByCode(ctx context.Context, code string, tenantID string) (*models.Role, error)

	// GetByName 根据角色名称获取单条记录
	GetByName(ctx context.Context, name string, tenantID string) (*models.Role, error)

	// Find 查询记录列表
	Find(ctx context.Context, filter bson.M) ([]*models.Role, error)

	// FindOne 查询单条记录
	FindOne(ctx context.Context, filter bson.M) (*models.Role, error)

	// FindWithPage 分页查询记录列表
	FindWithPage(ctx context.Context, filter bson.M, page, pageSize int64) ([]*models.Role, error)

	// Count 统计记录数
	Count(ctx context.Context, filter bson.M) (int64, error)

	// Create 创建记录
	Create(ctx context.Context, role *models.Role) error

	// Update 更新记录
	Update(ctx context.Context, id string, update bson.M) error

	// UpdateOne 按条件更新单条记录
	UpdateOne(ctx context.Context, filter bson.M, update bson.M) error

	// Delete 删除记录
	Delete(ctx context.Context, id string) error

	// DeleteMany 批量删除
	DeleteMany(ctx context.Context, filter bson.M) (int64, error)

	// FindDeletedWithPage 分页查询已删除记录列表
	FindDeletedWithPage(ctx context.Context, filter bson.M, page, pageSize int64) ([]*models.Role, error)

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

	// GetRolesByIDs 根据角色ID数组获取角色列表
	GetRolesByIDs(ctx context.Context, ids []string) ([]*models.Role, error)

	// GetRolesByTenant 获取租户下的所有角色
	GetRolesByTenant(ctx context.Context, tenantID string) ([]*models.Role, error)
}

// roleRepository Role数据仓库实现
type roleRepository struct {
	collection *mongo.Collection
}

// NewRoleRepository 创建Role数据仓库实例
func NewRoleRepository() RoleRepository {
	collection := database.MongoDB.Collection("role")
	return &roleRepository{
		collection: collection,
	}
}

// Get 根据ID获取单条记录（排除软删除）
func (r *roleRepository) Get(ctx context.Context, id string) (*models.Role, error) {
	// 将字符串 ID 转换为 ObjectID 进行查询
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		// 如果转换失败，尝试直接用字符串查询（兼容旧数据）
		filter := bson.M{"_id": id, "is_deleted": 0}
		role := &models.Role{}
		err = r.collection.FindOne(ctx, filter).Decode(role)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return nil, nil
			}
			return nil, err
		}
		return role, nil
	}

	// 使用 ObjectID 查询（排除软删除）
	filter := bson.M{"_id": objectID, "is_deleted": 0}
	role := &models.Role{}
	err = r.collection.FindOne(ctx, filter).Decode(role)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return role, nil
}

// GetByCode 根据角色代码获取单条记录（排除软删除）
func (r *roleRepository) GetByCode(ctx context.Context, code string, tenantID string) (*models.Role, error) {
	filter := bson.M{"code": code, "tenant_id": tenantID, "is_deleted": 0}
	role := &models.Role{}
	err := r.collection.FindOne(ctx, filter).Decode(role)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return role, nil
}

// GetByName 根据角色名称获取单条记录（排除软删除）
func (r *roleRepository) GetByName(ctx context.Context, name string, tenantID string) (*models.Role, error) {
	filter := bson.M{"name": name, "tenant_id": tenantID, "is_deleted": 0}
	role := &models.Role{}
	err := r.collection.FindOne(ctx, filter).Decode(role)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return role, nil
}

// Find 查询记录列表
func (r *roleRepository) Find(ctx context.Context, filter bson.M) ([]*models.Role, error) {
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	roles := []*models.Role{}
	err = cursor.All(ctx, &roles)
	if err != nil {
		return nil, err
	}
	return roles, nil
}

// FindOne 查询单条记录
func (r *roleRepository) FindOne(ctx context.Context, filter bson.M) (*models.Role, error) {
	role := &models.Role{}
	err := r.collection.FindOne(ctx, filter).Decode(role)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return role, nil
}

// FindWithPage 分页查询记录列表（按创建时间倒序）
func (r *roleRepository) FindWithPage(ctx context.Context, filter bson.M, page, pageSize int64) ([]*models.Role, error) {
	// 简化实现，不使用 options，让 service 层处理分页逻辑
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	roles := []*models.Role{}
	err = cursor.All(ctx, &roles)
	if err != nil {
		return nil, err
	}
	return roles, nil
}

// Count 统计记录数
func (r *roleRepository) Count(ctx context.Context, filter bson.M) (int64, error) {
	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// Create 创建记录
func (r *roleRepository) Create(ctx context.Context, role *models.Role) error {
	result, err := r.collection.InsertOne(ctx, role)
	if err != nil {
		return err
	}

	// 将生成的 ObjectID 转换为字符串并设置到 role.ID
	if oid, ok := result.InsertedID.(bson.ObjectID); ok {
		role.ID = oid.Hex()
	}

	return nil
}

// Update 更新记录
func (r *roleRepository) Update(ctx context.Context, id string, update bson.M) error {
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
func (r *roleRepository) UpdateOne(ctx context.Context, filter bson.M, update bson.M) error {
	updateDoc := bson.M{"$set": update}
	_, err := r.collection.UpdateOne(ctx, filter, updateDoc)
	return err
}

// Delete 软删除记录（设置 is_deleted = 1）
func (r *roleRepository) Delete(ctx context.Context, id string) error {
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

// HardDelete 物理删除记录（真正从数据库删除）
func (r *roleRepository) HardDelete(ctx context.Context, id string) error {
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

// DeleteMany 批量删除
func (r *roleRepository) DeleteMany(ctx context.Context, filter bson.M) (int64, error) {
	result, err := r.collection.DeleteMany(ctx, filter)
	if err != nil {
		return 0, err
	}
	return result.DeletedCount, nil
}

// FindDeletedWithPage 分页查询已删除记录列表（按删除时间倒序）
func (r *roleRepository) FindDeletedWithPage(ctx context.Context, filter bson.M, page, pageSize int64) ([]*models.Role, error) {
	// 添加已删除条件
	filter["is_deleted"] = 1

	// 简化实现，不使用 options，让 service 层处理分页逻辑
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	roles := []*models.Role{}
	err = cursor.All(ctx, &roles)
	if err != nil {
		return nil, err
	}
	return roles, nil
}

// CountDeleted 统计已删除记录数
func (r *roleRepository) CountDeleted(ctx context.Context, filter bson.M) (int64, error) {
	filter["is_deleted"] = 1
	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// RestoreOne 恢复单条记录（清除 is_deleted 和 deleted_at）
func (r *roleRepository) RestoreOne(ctx context.Context, id string) error {
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
func (r *roleRepository) RestoreMany(ctx context.Context, ids []string) (int64, error) {
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

// HardDeleteMany 批量物理删除记录
func (r *roleRepository) HardDeleteMany(ctx context.Context, ids []string) (int64, error) {
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
func (r *roleRepository) GetCollection() *mongo.Collection {
	return r.collection
}

// GetRolesByIDs 根据角色ID数组获取角色列表
func (r *roleRepository) GetRolesByIDs(ctx context.Context, ids []string) ([]*models.Role, error) {
	if len(ids) == 0 {
		return []*models.Role{}, nil
	}

	// 转换 ID 列表
	var objectIDs []interface{}
	for _, id := range ids {
		if objectID, err := bson.ObjectIDFromHex(id); err == nil {
			objectIDs = append(objectIDs, objectID)
		} else {
			objectIDs = append(objectIDs, id)
		}
	}

	filter := bson.M{
		"_id":        bson.M{"$in": objectIDs},
		"is_deleted": 0,
		"status":     1,
	}

	return r.Find(ctx, filter)
}

// GetRolesByTenant 获取租户下的所有角色
func (r *roleRepository) GetRolesByTenant(ctx context.Context, tenantID string) ([]*models.Role, error) {
	filter := bson.M{
		"tenant_id":  tenantID,
		"is_deleted": 0,
		"status":     1,
	}
	return r.Find(ctx, filter)
}
