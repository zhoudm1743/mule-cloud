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

// TenantMemberRepository 租户成员数据仓库接口
type TenantMemberRepository interface {
	// Get 根据ID获取单条记录
	Get(ctx context.Context, id string) (*models.TenantMember, error)

	// GetByUnionID 根据UnionID获取单条记录（在指定租户库中）
	GetByUnionID(ctx context.Context, unionID string) (*models.TenantMember, error)

	// GetByUserID 根据UserID获取单条记录（在指定租户库中）
	GetByUserID(ctx context.Context, userID string) (*models.TenantMember, error)

	// GetByJobNumber 根据工号获取单条记录
	GetByJobNumber(ctx context.Context, jobNumber string) (*models.TenantMember, error)

	// Find 查询记录列表
	Find(ctx context.Context, filter bson.M) ([]*models.TenantMember, error)

	// FindOne 查询单条记录
	FindOne(ctx context.Context, filter bson.M) (*models.TenantMember, error)

	// FindWithPage 分页查询记录列表
	FindWithPage(ctx context.Context, filter bson.M, page, pageSize int64) ([]*models.TenantMember, error)

	// Count 统计记录数
	Count(ctx context.Context, filter bson.M) (int64, error)

	// Create 创建记录
	Create(ctx context.Context, member *models.TenantMember) error

	// Update 更新记录
	Update(ctx context.Context, id string, update bson.M) error

	// UpdateOne 按条件更新单条记录
	UpdateOne(ctx context.Context, filter bson.M, update bson.M) error

	// SetInactive 设置成员为离职状态
	SetInactive(ctx context.Context, id string) error

	// SetActive 设置成员为在职状态
	SetActive(ctx context.Context, id string) error

	// Delete 软删除记录
	Delete(ctx context.Context, id string) error

	// HardDelete 物理删除记录
	HardDelete(ctx context.Context, id string) error
}

// tenantMemberRepository 租户成员数据仓库实现
type tenantMemberRepository struct {
	dbManager *database.DatabaseManager
}

// NewTenantMemberRepository 创建租户成员数据仓库实例
func NewTenantMemberRepository() TenantMemberRepository {
	return &tenantMemberRepository{
		dbManager: database.GetDatabaseManager(),
	}
}

// getCollection 获取集合（自动根据Context中的租户Code切换数据库）
func (r *tenantMemberRepository) getCollection(ctx context.Context) *mongo.Collection {
	tenantCode := tenantCtx.GetTenantCode(ctx)
	db := r.dbManager.GetDatabase(tenantCode)
	return db.Collection("member")
}

// Get 根据ID获取单条记录（排除软删除）
func (r *tenantMemberRepository) Get(ctx context.Context, id string) (*models.TenantMember, error) {
	collection := r.getCollection(ctx)

	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		filter := bson.M{"_id": id, "is_deleted": 0}
		member := &models.TenantMember{}
		err = collection.FindOne(ctx, filter).Decode(member)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return nil, nil
			}
			return nil, err
		}
		return member, nil
	}

	filter := bson.M{"_id": objectID, "is_deleted": 0}
	member := &models.TenantMember{}
	err = collection.FindOne(ctx, filter).Decode(member)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return member, nil
}

// GetByUnionID 根据UnionID获取单条记录（排除软删除）
func (r *tenantMemberRepository) GetByUnionID(ctx context.Context, unionID string) (*models.TenantMember, error) {
	collection := r.getCollection(ctx)
	filter := bson.M{"union_id": unionID, "is_deleted": 0}
	member := &models.TenantMember{}
	err := collection.FindOne(ctx, filter).Decode(member)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return member, nil
}

// GetByUserID 根据UserID获取单条记录（排除软删除）
func (r *tenantMemberRepository) GetByUserID(ctx context.Context, userID string) (*models.TenantMember, error) {
	collection := r.getCollection(ctx)
	filter := bson.M{"user_id": userID, "is_deleted": 0}
	member := &models.TenantMember{}
	err := collection.FindOne(ctx, filter).Decode(member)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return member, nil
}

// GetByJobNumber 根据工号获取单条记录（排除软删除）
func (r *tenantMemberRepository) GetByJobNumber(ctx context.Context, jobNumber string) (*models.TenantMember, error) {
	collection := r.getCollection(ctx)
	filter := bson.M{"job_number": jobNumber, "is_deleted": 0}
	member := &models.TenantMember{}
	err := collection.FindOne(ctx, filter).Decode(member)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return member, nil
}

// Find 查询记录列表
func (r *tenantMemberRepository) Find(ctx context.Context, filter bson.M) ([]*models.TenantMember, error) {
	collection := r.getCollection(ctx)
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	members := []*models.TenantMember{}
	err = cursor.All(ctx, &members)
	if err != nil {
		return nil, err
	}
	return members, nil
}

// FindOne 查询单条记录
func (r *tenantMemberRepository) FindOne(ctx context.Context, filter bson.M) (*models.TenantMember, error) {
	collection := r.getCollection(ctx)
	member := &models.TenantMember{}
	err := collection.FindOne(ctx, filter).Decode(member)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return member, nil
}

// FindWithPage 分页查询记录列表
func (r *tenantMemberRepository) FindWithPage(ctx context.Context, filter bson.M, page, pageSize int64) ([]*models.TenantMember, error) {
	collection := r.getCollection(ctx)
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	members := []*models.TenantMember{}
	err = cursor.All(ctx, &members)
	if err != nil {
		return nil, err
	}
	return members, nil
}

// Count 统计记录数
func (r *tenantMemberRepository) Count(ctx context.Context, filter bson.M) (int64, error) {
	collection := r.getCollection(ctx)
	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// Create 创建记录
func (r *tenantMemberRepository) Create(ctx context.Context, member *models.TenantMember) error {
	collection := r.getCollection(ctx)
	result, err := collection.InsertOne(ctx, member)
	if err != nil {
		return err
	}

	if oid, ok := result.InsertedID.(bson.ObjectID); ok {
		member.ID = oid.Hex()
	}

	return nil
}

// Update 更新记录
func (r *tenantMemberRepository) Update(ctx context.Context, id string, update bson.M) error {
	collection := r.getCollection(ctx)

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
func (r *tenantMemberRepository) UpdateOne(ctx context.Context, filter bson.M, update bson.M) error {
	collection := r.getCollection(ctx)
	updateDoc := bson.M{"$set": update}
	_, err := collection.UpdateOne(ctx, filter, updateDoc)
	return err
}

// SetInactive 设置成员为离职状态
func (r *tenantMemberRepository) SetInactive(ctx context.Context, id string) error {
	update := bson.M{
		"status":     "inactive",
		"left_at":    time.Now().Unix(),
		"updated_at": time.Now().Unix(),
	}
	return r.Update(ctx, id, update)
}

// SetActive 设置成员为在职状态
func (r *tenantMemberRepository) SetActive(ctx context.Context, id string) error {
	update := bson.M{
		"status":     "active",
		"employed_at": time.Now().Unix(),
		"left_at":    0,
		"updated_at": time.Now().Unix(),
	}
	return r.Update(ctx, id, update)
}

// Delete 软删除记录
func (r *tenantMemberRepository) Delete(ctx context.Context, id string) error {
	collection := r.getCollection(ctx)

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
func (r *tenantMemberRepository) HardDelete(ctx context.Context, id string) error {
	collection := r.getCollection(ctx)

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

