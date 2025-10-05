package services

import (
	"context"
	"mule-cloud/app/basic/dto"
	"mule-cloud/internal/models"
	"mule-cloud/internal/repository"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// IOrderTypeService 订单类型服务接口
type IOrderTypeService interface {
	Get(ctx context.Context, id string) (*models.Basic, error)
	GetAll(ctx context.Context, req dto.OrderTypeListRequest) ([]models.Basic, error)
	List(ctx context.Context, req dto.OrderTypeListRequest) ([]models.Basic, int64, error)
	Create(ctx context.Context, req dto.OrderTypeCreateRequest) (*models.Basic, error)
	Update(ctx context.Context, req dto.OrderTypeUpdateRequest) (*models.Basic, error)
	Delete(ctx context.Context, id string) error
}

// OrderTypeService 订单类型服务实现
type OrderTypeService struct {
	repo repository.BasicRepository
}

// NewOrderTypeService 创建订单类型服务
func NewOrderTypeService() IOrderTypeService {
	repo := repository.NewBasicRepository()
	return &OrderTypeService{repo: repo}
}

// Get 获取订单类型
func (s *OrderTypeService) Get(ctx context.Context, id string) (*models.Basic, error) {
	return s.repo.Get(ctx, id)
}

// List 列表（分页查询）
func (s *OrderTypeService) List(ctx context.Context, req dto.OrderTypeListRequest) ([]models.Basic, int64, error) {
	// 构建过滤条件
	filter := bson.M{"name": "order_type", "is_deleted": 0}
	if req.Value != "" {
		// 支持模糊搜索
		filter["value"] = bson.M{"$regex": req.Value, "$options": "i"}
	}
	if req.ID != "" {
		filter["_id"] = req.ID
	}
	// 获取总数
	total, err := s.repo.Count(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (req.Page - 1) * req.PageSize
	opts := options.Find().
		SetSkip(offset).
		SetLimit(req.PageSize).
		SetSort(bson.M{"created_at": -1})

	// 使用 GetCollection 获取原始集合以使用 options
	collection := s.repo.GetCollectionWithContext(ctx)
	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	order_types := []models.Basic{}
	err = cursor.All(ctx, &order_types)
	if err != nil {
		return nil, 0, err
	}

	return order_types, total, nil
}

// GetAll 获取所有订单类型（不分页）
func (s *OrderTypeService) GetAll(ctx context.Context, req dto.OrderTypeListRequest) ([]models.Basic, error) {
	// 构建过滤条件
	filter := bson.M{"name": "order_type", "is_deleted": 0}
	if req.Value != "" {
		// 支持模糊搜索
		filter["value"] = bson.M{"$regex": req.Value, "$options": "i"}
	}
	if req.ID != "" {
		filter["_id"] = req.ID
	}
	// 使用 GetCollection 获取原始集合以使用排序选项
	opts := options.Find().SetSort(bson.M{"created_at": -1})
	collection := s.repo.GetCollectionWithContext(ctx)
	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	order_types := []models.Basic{}
	err = cursor.All(ctx, &order_types)
	if err != nil {
		return nil, err
	}

	return order_types, nil
}

// Create 创建订单类型
func (s *OrderTypeService) Create(ctx context.Context, req dto.OrderTypeCreateRequest) (*models.Basic, error) {
	now := time.Now().Unix()

	basic := &models.Basic{
		Name:      "order_type",
		Value:     req.Value,
		Remark:    req.Remark,
		CreatedAt: now,
		UpdatedAt: now,
	}

	err := s.repo.Create(ctx, basic)
	if err != nil {
		return nil, err
	}

	return basic, nil
}

// Update 更新订单类型
func (s *OrderTypeService) Update(ctx context.Context, req dto.OrderTypeUpdateRequest) (*models.Basic, error) {
	// 更新字段
	update := bson.M{
		"value":      req.Value,
		"remark":     req.Remark,
		"updated_at": time.Now().Unix(),
	}

	err := s.repo.Update(ctx, req.ID, update)
	if err != nil {
		return nil, err
	}

	// 返回更新后的数据
	return s.repo.Get(ctx, req.ID)
}

// Delete 删除订单类型
func (s *OrderTypeService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
