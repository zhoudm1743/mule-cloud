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

// ISalesmanService 业务员服务接口
type ISalesmanService interface {
	Get(ctx context.Context, id string) (*models.Basic, error)
	GetAll(ctx context.Context, req dto.SalesmanListRequest) ([]models.Basic, error)
	List(ctx context.Context, req dto.SalesmanListRequest) ([]models.Basic, int64, error)
	Create(ctx context.Context, req dto.SalesmanCreateRequest) (*models.Basic, error)
	Update(ctx context.Context, req dto.SalesmanUpdateRequest) (*models.Basic, error)
	Delete(ctx context.Context, id string) error
}

// SalesmanService 业务员服务实现
type SalesmanService struct {
	repo repository.BasicRepository
}

// NewSalesmanService 创建业务员服务
func NewSalesmanService() ISalesmanService {
	repo := repository.NewBasicRepository()
	return &SalesmanService{repo: repo}
}

// Get 获取业务员
func (s *SalesmanService) Get(ctx context.Context, id string) (*models.Basic, error) {
	return s.repo.Get(ctx, id)
}

// List 列表（分页查询）
func (s *SalesmanService) List(ctx context.Context, req dto.SalesmanListRequest) ([]models.Basic, int64, error) {
	// 构建过滤条件
	filter := bson.M{"name": "salesman", "is_deleted": 0}
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

	salesmans := []models.Basic{}
	err = cursor.All(ctx, &salesmans)
	if err != nil {
		return nil, 0, err
	}

	return salesmans, total, nil
}

// GetAll 获取所有业务员（不分页）
func (s *SalesmanService) GetAll(ctx context.Context, req dto.SalesmanListRequest) ([]models.Basic, error) {
	// 构建过滤条件
	filter := bson.M{"name": "salesman", "is_deleted": 0}
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

	salesmans := []models.Basic{}
	err = cursor.All(ctx, &salesmans)
	if err != nil {
		return nil, err
	}

	return salesmans, nil
}

// Create 创建业务员
func (s *SalesmanService) Create(ctx context.Context, req dto.SalesmanCreateRequest) (*models.Basic, error) {
	now := time.Now().Unix()

	basic := &models.Basic{
		Name:      "salesman",
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

// Update 更新业务员
func (s *SalesmanService) Update(ctx context.Context, req dto.SalesmanUpdateRequest) (*models.Basic, error) {
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

// Delete 删除业务员
func (s *SalesmanService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
