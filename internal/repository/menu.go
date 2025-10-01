package repository

import (
	"context"
	"fmt"
	"mule-cloud/core/database"
	"mule-cloud/internal/models"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MenuRepository struct {
	collection *mongo.Collection
}

func NewMenuRepository() *MenuRepository {
	return &MenuRepository{
		collection: database.MongoDB.Collection(models.Menu{}.TableName()),
	}
}

// Create 创建菜单
func (r *MenuRepository) Create(ctx context.Context, menu *models.Menu) error {
	menu.CreatedAt = time.Now().Unix()
	menu.UpdatedAt = time.Now().Unix()
	menu.Status = 1
	menu.IsDeleted = 0

	result, err := r.collection.InsertOne(ctx, menu)
	if err != nil {
		return fmt.Errorf("创建菜单失败: %w", err)
	}

	if oid, ok := result.InsertedID.(bson.ObjectID); ok {
		menu.ID = oid.Hex()
	}
	return nil
}

// GetByID 根据ID获取菜单
func (r *MenuRepository) GetByID(ctx context.Context, id string) (*models.Menu, error) {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("无效的ID格式: %w", err)
	}

	var menu models.Menu
	filter := bson.M{
		"_id":        objectID,
		"is_deleted": 0,
	}

	err = r.collection.FindOne(ctx, filter).Decode(&menu)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("菜单不存在")
		}
		return nil, fmt.Errorf("查询菜单失败: %w", err)
	}

	menu.ID = objectID.Hex()
	return &menu, nil
}

// GetAll 获取所有菜单（不分页）
func (r *MenuRepository) GetAll(ctx context.Context) ([]*models.Menu, error) {
	filter := bson.M{
		"is_deleted": 0,
		"status":     1,
	}

	// 按 order 排序
	opts := options.Find().SetSort(bson.D{
		bson.E{Key: "order", Value: 1},
		bson.E{Key: "created_at", Value: 1},
	})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("查询菜单列表失败: %w", err)
	}
	defer cursor.Close(ctx)

	var menus []*models.Menu
	for cursor.Next(ctx) {
		var menu models.Menu
		if err := cursor.Decode(&menu); err != nil {
			continue
		}
		// ObjectID已经在模型中
		menus = append(menus, &menu)
	}

	return menus, nil
}

// List 分页查询菜单
func (r *MenuRepository) List(ctx context.Context, page, pageSize int, filters map[string]interface{}) ([]*models.Menu, int64, error) {
	filter := bson.M{"is_deleted": 0}

	// 添加筛选条件
	if name, ok := filters["name"].(string); ok && name != "" {
		filter["name"] = bson.M{"$regex": name, "$options": "i"}
	}
	if title, ok := filters["title"].(string); ok && title != "" {
		filter["title"] = bson.M{"$regex": title, "$options": "i"}
	}
	if menuType, ok := filters["menuType"].(string); ok && menuType != "" {
		filter["menuType"] = menuType
	}
	if status, ok := filters["status"].(int); ok {
		filter["status"] = status
	}

	// 计算总数
	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("统计菜单数量失败: %w", err)
	}

	// 查询数据
	skip := (page - 1) * pageSize
	opts := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(pageSize)).
		SetSort(bson.D{
			bson.E{Key: "order", Value: 1},
			bson.E{Key: "created_at", Value: -1},
		})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, fmt.Errorf("查询菜单列表失败: %w", err)
	}
	defer cursor.Close(ctx)

	var menus []*models.Menu
	for cursor.Next(ctx) {
		var menu models.Menu
		if err := cursor.Decode(&menu); err != nil {
			continue
		}
		menus = append(menus, &menu)
	}

	return menus, total, nil
}

// Update 更新菜单
func (r *MenuRepository) Update(ctx context.Context, id string, updates map[string]interface{}) error {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("无效的ID格式: %w", err)
	}

	updates["updated_at"] = time.Now().Unix()

	filter := bson.M{
		"_id":        objectID,
		"is_deleted": 0,
	}

	update := bson.M{"$set": updates}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("更新菜单失败: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("菜单不存在")
	}

	return nil
}

// Delete 删除菜单（软删除）
func (r *MenuRepository) Delete(ctx context.Context, id string) error {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("无效的ID格式: %w", err)
	}

	filter := bson.M{
		"_id":        objectID,
		"is_deleted": 0,
	}

	update := bson.M{
		"$set": bson.M{
			"is_deleted": 1,
			"deleted_at": time.Now().Unix(),
			"updated_at": time.Now().Unix(),
		},
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("删除菜单失败: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("菜单不存在")
	}

	return nil
}

// BatchDelete 批量删除菜单
func (r *MenuRepository) BatchDelete(ctx context.Context, ids []string) error {
	var objectIDs []bson.ObjectID
	for _, id := range ids {
		objectID, err := bson.ObjectIDFromHex(id)
		if err != nil {
			continue
		}
		objectIDs = append(objectIDs, objectID)
	}

	if len(objectIDs) == 0 {
		return fmt.Errorf("没有有效的ID")
	}

	filter := bson.M{
		"_id":        bson.M{"$in": objectIDs},
		"is_deleted": 0,
	}

	update := bson.M{
		"$set": bson.M{
			"is_deleted": 1,
			"deleted_at": time.Now().Unix(),
			"updated_at": time.Now().Unix(),
		},
	}

	_, err := r.collection.UpdateMany(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("批量删除菜单失败: %w", err)
	}

	return nil
}

// GetByName 根据名称获取菜单
func (r *MenuRepository) GetByName(ctx context.Context, name string) (*models.Menu, error) {
	var menu models.Menu
	filter := bson.M{
		"name":       name,
		"is_deleted": 0,
	}

	err := r.collection.FindOne(ctx, filter).Decode(&menu)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("查询菜单失败: %w", err)
	}

	return &menu, nil
}

// GetByPID 根据父级ID获取子菜单
func (r *MenuRepository) GetByPID(ctx context.Context, pid *string) ([]*models.Menu, error) {
	filter := bson.M{
		"is_deleted": 0,
		"status":     1,
	}

	if pid == nil {
		filter["pid"] = nil
	} else {
		filter["pid"] = *pid
	}

	opts := options.Find().SetSort(bson.D{
		bson.E{Key: "order", Value: 1},
	})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("查询子菜单失败: %w", err)
	}
	defer cursor.Close(ctx)

	var menus []*models.Menu
	for cursor.Next(ctx) {
		var menu models.Menu
		if err := cursor.Decode(&menu); err != nil {
			continue
		}
		menus = append(menus, &menu)
	}

	return menus, nil
}
