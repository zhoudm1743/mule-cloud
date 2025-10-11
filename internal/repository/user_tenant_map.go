package repository

import (
	"context"
	"mule-cloud/core/database"
	"mule-cloud/internal/models"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// UserTenantMapRepository 用户-租户映射数据仓库接口
type UserTenantMapRepository interface {
	// Get 根据ID获取单条记录
	Get(ctx context.Context, id string) (*models.UserTenantMap, error)

	// GetByUserAndTenant 根据用户ID和租户ID获取关系
	GetByUserAndTenant(ctx context.Context, userID, tenantID string) (*models.UserTenantMap, error)

	// GetByUnionIDAndTenant 根据UnionID和租户ID获取关系
	GetByUnionIDAndTenant(ctx context.Context, unionID, tenantID string) (*models.UserTenantMap, error)

	// GetUserTenants 获取用户的所有租户关联（包括活跃和非活跃）
	GetUserTenants(ctx context.Context, userID string) ([]*models.UserTenantMap, error)

	// GetUserActiveTenants 获取用户的活跃租户关联
	GetUserActiveTenants(ctx context.Context, userID string) ([]*models.UserTenantMap, error)

	// GetTenantMembers 获取租户的所有成员关联
	GetTenantMembers(ctx context.Context, tenantID string) ([]*models.UserTenantMap, error)

	// GetTenantActiveMembers 获取租户的活跃成员关联
	GetTenantActiveMembers(ctx context.Context, tenantID string) ([]*models.UserTenantMap, error)

	// Find 查询记录列表
	Find(ctx context.Context, filter bson.M) ([]*models.UserTenantMap, error)

	// FindOne 查询单条记录
	FindOne(ctx context.Context, filter bson.M) (*models.UserTenantMap, error)

	// Count 统计记录数
	Count(ctx context.Context, filter bson.M) (int64, error)

	// Create 创建记录
	Create(ctx context.Context, mapping *models.UserTenantMap) error

	// Update 更新记录
	Update(ctx context.Context, id string, update bson.M) error

	// UpdateOne 按条件更新单条记录
	UpdateOne(ctx context.Context, filter bson.M, update bson.M) error

	// SetInactive 设置关系为非活跃状态（离职）
	SetInactive(ctx context.Context, id string) error

	// SetActive 设置关系为活跃状态（入职）
	SetActive(ctx context.Context, id string) error

	// Delete 软删除记录
	Delete(ctx context.Context, id string) error

	// HardDelete 物理删除记录
	HardDelete(ctx context.Context, id string) error
}

// userTenantMapRepository 用户-租户映射数据仓库实现
type userTenantMapRepository struct {
	dbManager *database.DatabaseManager
}

// NewUserTenantMapRepository 创建用户-租户映射数据仓库实例
func NewUserTenantMapRepository() UserTenantMapRepository {
	return &userTenantMapRepository{
		dbManager: database.GetDatabaseManager(),
	}
}

// getCollection 获取集合（用户-租户映射固定存储在系统库）
func (r *userTenantMapRepository) getCollection() *mongo.Collection {
	db := r.dbManager.GetSystemDatabase()
	return db.Collection("user_tenant_map")
}

// Get 根据ID获取单条记录（排除软删除）
func (r *userTenantMapRepository) Get(ctx context.Context, id string) (*models.UserTenantMap, error) {
	collection := r.getCollection()

	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		filter := bson.M{"_id": id, "is_deleted": 0}
		mapping := &models.UserTenantMap{}
		err = collection.FindOne(ctx, filter).Decode(mapping)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return nil, nil
			}
			return nil, err
		}
		return mapping, nil
	}

	filter := bson.M{"_id": objectID, "is_deleted": 0}
	mapping := &models.UserTenantMap{}
	err = collection.FindOne(ctx, filter).Decode(mapping)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return mapping, nil
}

// GetByUserAndTenant 根据用户ID和租户ID获取关系
func (r *userTenantMapRepository) GetByUserAndTenant(ctx context.Context, userID, tenantID string) (*models.UserTenantMap, error) {
	collection := r.getCollection()
	filter := bson.M{
		"user_id":    userID,
		"tenant_id":  tenantID,
		"is_deleted": 0,
	}
	mapping := &models.UserTenantMap{}
	err := collection.FindOne(ctx, filter).Decode(mapping)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return mapping, nil
}

// GetByUnionIDAndTenant 根据UnionID和租户ID获取关系
func (r *userTenantMapRepository) GetByUnionIDAndTenant(ctx context.Context, unionID, tenantID string) (*models.UserTenantMap, error) {
	collection := r.getCollection()
	filter := bson.M{
		"union_id":   unionID,
		"tenant_id":  tenantID,
		"is_deleted": 0,
	}
	mapping := &models.UserTenantMap{}
	err := collection.FindOne(ctx, filter).Decode(mapping)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return mapping, nil
}

// GetUserTenants 获取用户的所有租户关联
func (r *userTenantMapRepository) GetUserTenants(ctx context.Context, userID string) ([]*models.UserTenantMap, error) {
	collection := r.getCollection()
	filter := bson.M{
		"user_id":    userID,
		"is_deleted": 0,
	}
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	mappings := []*models.UserTenantMap{}
	err = cursor.All(ctx, &mappings)
	if err != nil {
		return nil, err
	}
	return mappings, nil
}

// GetUserActiveTenants 获取用户的活跃租户关联
func (r *userTenantMapRepository) GetUserActiveTenants(ctx context.Context, userID string) ([]*models.UserTenantMap, error) {
	collection := r.getCollection()
	filter := bson.M{
		"user_id":    userID,
		"status":     "active",
		"is_deleted": 0,
	}
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	mappings := []*models.UserTenantMap{}
	err = cursor.All(ctx, &mappings)
	if err != nil {
		return nil, err
	}
	return mappings, nil
}

// GetTenantMembers 获取租户的所有成员关联
func (r *userTenantMapRepository) GetTenantMembers(ctx context.Context, tenantID string) ([]*models.UserTenantMap, error) {
	collection := r.getCollection()
	filter := bson.M{
		"tenant_id":  tenantID,
		"is_deleted": 0,
	}
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	mappings := []*models.UserTenantMap{}
	err = cursor.All(ctx, &mappings)
	if err != nil {
		return nil, err
	}
	return mappings, nil
}

// GetTenantActiveMembers 获取租户的活跃成员关联
func (r *userTenantMapRepository) GetTenantActiveMembers(ctx context.Context, tenantID string) ([]*models.UserTenantMap, error) {
	collection := r.getCollection()
	filter := bson.M{
		"tenant_id":  tenantID,
		"status":     "active",
		"is_deleted": 0,
	}
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	mappings := []*models.UserTenantMap{}
	err = cursor.All(ctx, &mappings)
	if err != nil {
		return nil, err
	}
	return mappings, nil
}

// Find 查询记录列表
func (r *userTenantMapRepository) Find(ctx context.Context, filter bson.M) ([]*models.UserTenantMap, error) {
	collection := r.getCollection()
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	mappings := []*models.UserTenantMap{}
	err = cursor.All(ctx, &mappings)
	if err != nil {
		return nil, err
	}
	return mappings, nil
}

// FindOne 查询单条记录
func (r *userTenantMapRepository) FindOne(ctx context.Context, filter bson.M) (*models.UserTenantMap, error) {
	collection := r.getCollection()
	mapping := &models.UserTenantMap{}
	err := collection.FindOne(ctx, filter).Decode(mapping)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return mapping, nil
}

// Count 统计记录数
func (r *userTenantMapRepository) Count(ctx context.Context, filter bson.M) (int64, error) {
	collection := r.getCollection()
	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// Create 创建记录
func (r *userTenantMapRepository) Create(ctx context.Context, mapping *models.UserTenantMap) error {
	collection := r.getCollection()
	result, err := collection.InsertOne(ctx, mapping)
	if err != nil {
		return err
	}

	if oid, ok := result.InsertedID.(bson.ObjectID); ok {
		mapping.ID = oid.Hex()
	}

	return nil
}

// Update 更新记录
func (r *userTenantMapRepository) Update(ctx context.Context, id string, update bson.M) error {
	collection := r.getCollection()

	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
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

	filter := bson.M{"_id": objectID}
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

// UpdateOne 按条件更新单条记录
func (r *userTenantMapRepository) UpdateOne(ctx context.Context, filter bson.M, update bson.M) error {
	collection := r.getCollection()
	updateDoc := bson.M{"$set": update}
	_, err := collection.UpdateOne(ctx, filter, updateDoc)
	return err
}

// SetInactive 设置关系为非活跃状态（离职）
func (r *userTenantMapRepository) SetInactive(ctx context.Context, id string) error {
	update := bson.M{
		"status":     "inactive",
		"left_at":    time.Now().Unix(),
		"updated_at": time.Now().Unix(),
	}
	return r.Update(ctx, id, update)
}

// SetActive 设置关系为活跃状态（入职）
func (r *userTenantMapRepository) SetActive(ctx context.Context, id string) error {
	update := bson.M{
		"status":     "active",
		"joined_at":  time.Now().Unix(),
		"left_at":    0,
		"updated_at": time.Now().Unix(),
	}
	return r.Update(ctx, id, update)
}

// Delete 软删除记录
func (r *userTenantMapRepository) Delete(ctx context.Context, id string) error {
	collection := r.getCollection()

	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		filter := bson.M{"_id": id, "is_deleted": 0}
		updateDoc := bson.M{"$set": bson.M{"is_deleted": 1, "updated_at": time.Now().Unix()}}
		result, err := collection.UpdateOne(ctx, filter, updateDoc)
		if err != nil {
			return err
		}
		if result.MatchedCount == 0 {
			return mongo.ErrNoDocuments
		}
		return nil
	}

	filter := bson.M{"_id": objectID, "is_deleted": 0}
	updateDoc := bson.M{"$set": bson.M{"is_deleted": 1, "updated_at": time.Now().Unix()}}
	result, err := collection.UpdateOne(ctx, filter, updateDoc)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}

// HardDelete 物理删除记录
func (r *userTenantMapRepository) HardDelete(ctx context.Context, id string) error {
	collection := r.getCollection()

	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
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

	filter := bson.M{"_id": objectID}
	result, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}

