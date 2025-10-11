package repository

import (
	"context"
	"mule-cloud/core/database"
	"mule-cloud/internal/models"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// WechatUserRepository 微信用户数据仓库接口
type WechatUserRepository interface {
	// Get 根据ID获取单条记录
	Get(ctx context.Context, id string) (*models.WechatUser, error)

	// GetByUnionID 根据UnionID获取单条记录
	GetByUnionID(ctx context.Context, unionID string) (*models.WechatUser, error)

	// GetByOpenID 根据OpenID获取单条记录
	GetByOpenID(ctx context.Context, openID string) (*models.WechatUser, error)

	// GetByPhone 根据手机号获取单条记录
	GetByPhone(ctx context.Context, phone string) (*models.WechatUser, error)

	// Find 查询记录列表
	Find(ctx context.Context, filter bson.M) ([]*models.WechatUser, error)

	// FindOne 查询单条记录
	FindOne(ctx context.Context, filter bson.M) (*models.WechatUser, error)

	// Count 统计记录数
	Count(ctx context.Context, filter bson.M) (int64, error)

	// Create 创建记录
	Create(ctx context.Context, user *models.WechatUser) error

	// Update 更新记录
	Update(ctx context.Context, id string, update bson.M) error

	// UpdateOne 按条件更新单条记录
	UpdateOne(ctx context.Context, filter bson.M, update bson.M) error

	// Delete 软删除记录
	Delete(ctx context.Context, id string) error

	// HardDelete 物理删除记录
	HardDelete(ctx context.Context, id string) error

	// AddTenant 添加租户ID到用户的租户列表
	AddTenant(ctx context.Context, userID string, tenantID string) error

	// RemoveTenant 从用户的租户列表中移除租户ID
	RemoveTenant(ctx context.Context, userID string, tenantID string) error
}

// wechatUserRepository 微信用户数据仓库实现
type wechatUserRepository struct {
	dbManager *database.DatabaseManager
}

// NewWechatUserRepository 创建微信用户数据仓库实例
func NewWechatUserRepository() WechatUserRepository {
	return &wechatUserRepository{
		dbManager: database.GetDatabaseManager(),
	}
}

// getCollection 获取集合（微信用户固定存储在系统库）
func (r *wechatUserRepository) getCollection() *mongo.Collection {
	db := r.dbManager.GetSystemDatabase()
	return db.Collection("wechat_user")
}

// Get 根据ID获取单条记录（排除软删除）
func (r *wechatUserRepository) Get(ctx context.Context, id string) (*models.WechatUser, error) {
	collection := r.getCollection()

	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		filter := bson.M{"_id": id, "is_deleted": 0}
		user := &models.WechatUser{}
		err = collection.FindOne(ctx, filter).Decode(user)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return nil, nil
			}
			return nil, err
		}
		return user, nil
	}

	filter := bson.M{"_id": objectID, "is_deleted": 0}
	user := &models.WechatUser{}
	err = collection.FindOne(ctx, filter).Decode(user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

// GetByUnionID 根据UnionID获取单条记录（排除软删除）
func (r *wechatUserRepository) GetByUnionID(ctx context.Context, unionID string) (*models.WechatUser, error) {
	collection := r.getCollection()
	filter := bson.M{"union_id": unionID, "is_deleted": 0}
	user := &models.WechatUser{}
	err := collection.FindOne(ctx, filter).Decode(user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

// GetByOpenID 根据OpenID获取单条记录（排除软删除）
func (r *wechatUserRepository) GetByOpenID(ctx context.Context, openID string) (*models.WechatUser, error) {
	collection := r.getCollection()
	filter := bson.M{"open_id": openID, "is_deleted": 0}
	user := &models.WechatUser{}
	err := collection.FindOne(ctx, filter).Decode(user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

// GetByPhone 根据手机号获取单条记录（排除软删除）
func (r *wechatUserRepository) GetByPhone(ctx context.Context, phone string) (*models.WechatUser, error) {
	collection := r.getCollection()
	filter := bson.M{"phone": phone, "is_deleted": 0}
	user := &models.WechatUser{}
	err := collection.FindOne(ctx, filter).Decode(user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

// Find 查询记录列表
func (r *wechatUserRepository) Find(ctx context.Context, filter bson.M) ([]*models.WechatUser, error) {
	collection := r.getCollection()
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	users := []*models.WechatUser{}
	err = cursor.All(ctx, &users)
	if err != nil {
		return nil, err
	}
	return users, nil
}

// FindOne 查询单条记录
func (r *wechatUserRepository) FindOne(ctx context.Context, filter bson.M) (*models.WechatUser, error) {
	collection := r.getCollection()
	user := &models.WechatUser{}
	err := collection.FindOne(ctx, filter).Decode(user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

// Count 统计记录数
func (r *wechatUserRepository) Count(ctx context.Context, filter bson.M) (int64, error) {
	collection := r.getCollection()
	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// Create 创建记录
func (r *wechatUserRepository) Create(ctx context.Context, user *models.WechatUser) error {
	collection := r.getCollection()
	result, err := collection.InsertOne(ctx, user)
	if err != nil {
		return err
	}

	if oid, ok := result.InsertedID.(bson.ObjectID); ok {
		user.ID = oid.Hex()
	}

	return nil
}

// Update 更新记录
func (r *wechatUserRepository) Update(ctx context.Context, id string, update bson.M) error {
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
func (r *wechatUserRepository) UpdateOne(ctx context.Context, filter bson.M, update bson.M) error {
	collection := r.getCollection()
	updateDoc := bson.M{"$set": update}
	_, err := collection.UpdateOne(ctx, filter, updateDoc)
	return err
}

// Delete 软删除记录
func (r *wechatUserRepository) Delete(ctx context.Context, id string) error {
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
func (r *wechatUserRepository) HardDelete(ctx context.Context, id string) error {
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

// AddTenant 添加租户ID到用户的租户列表
func (r *wechatUserRepository) AddTenant(ctx context.Context, userID string, tenantID string) error {
	collection := r.getCollection()

	objectID, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		filter := bson.M{"_id": userID}
		updateDoc := bson.M{
			"$addToSet": bson.M{"tenant_ids": tenantID},
			"$set":      bson.M{"updated_at": time.Now().Unix()},
		}
		_, err := collection.UpdateOne(ctx, filter, updateDoc)
		return err
	}

	filter := bson.M{"_id": objectID}
	updateDoc := bson.M{
		"$addToSet": bson.M{"tenant_ids": tenantID},
		"$set":      bson.M{"updated_at": time.Now().Unix()},
	}
	_, err = collection.UpdateOne(ctx, filter, updateDoc)
	return err
}

// RemoveTenant 从用户的租户列表中移除租户ID
func (r *wechatUserRepository) RemoveTenant(ctx context.Context, userID string, tenantID string) error {
	collection := r.getCollection()

	objectID, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		filter := bson.M{"_id": userID}
		updateDoc := bson.M{
			"$pull": bson.M{"tenant_ids": tenantID},
			"$set":  bson.M{"updated_at": time.Now().Unix()},
		}
		_, err := collection.UpdateOne(ctx, filter, updateDoc)
		return err
	}

	filter := bson.M{"_id": objectID}
	updateDoc := bson.M{
		"$pull": bson.M{"tenant_ids": tenantID},
		"$set":  bson.M{"updated_at": time.Now().Unix()},
	}
	_, err = collection.UpdateOne(ctx, filter, updateDoc)
	return err
}

