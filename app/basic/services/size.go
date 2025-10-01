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

// ISizeService 尺寸服务接口
type ISizeService interface {
	Get(id string) (*models.Basic, error)
	GetAll(req dto.SizeListRequest) ([]*models.Basic, error)
	List(req dto.SizeListRequest) ([]*models.Basic, int64, error)
	Create(req dto.SizeCreateRequest) (*models.Basic, error)
	Update(req dto.SizeUpdateRequest) (*models.Basic, error)
	Delete(id string) error
}

// SizeService 尺寸服务实现
type SizeService struct {
	repo repository.BasicRepository
}

// NewSizeService 创建尺寸服务
func NewSizeService() ISizeService {
	repo := repository.NewBasicRepository()
	return &SizeService{repo: repo}
}

// Get 获取尺寸
func (s *SizeService) Get(id string) (*models.Basic, error) {
	ctx := context.Background()
	return s.repo.Get(ctx, id)
}

// GetAll 获取所有尺寸（不分页）
func (s *SizeService) GetAll(req dto.SizeListRequest) ([]*models.Basic, error) {
	ctx := context.Background()

	// 构建过滤条件
	filter := bson.M{"name": "size"}
	if req.Value != "" {
		filter["value"] = req.Value
	}
	if req.ID != "" {
		filter["_id"] = req.ID
	}

	// 使用 GetCollection 获取原始集合以使用排序选项
	opts := options.Find().SetSort(bson.M{"created_at": 1}) // 按创建时间正序
	collection := s.repo.GetCollection()
	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var sizes []*models.Basic
	err = cursor.All(ctx, &sizes)
	if err != nil {
		return nil, err
	}

	return sizes, nil
}

// List 列表（分页查询）
func (s *SizeService) List(req dto.SizeListRequest) ([]*models.Basic, int64, error) {
	ctx := context.Background()

	// 构建过滤条件
	filter := bson.M{"name": "size"}
	if req.Value != "" {
		filter["value"] = req.Value
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
		SetSort(bson.M{"created_at": 1})

	// 使用 GetCollection 获取原始集合以使用 options
	collection := s.repo.GetCollection()
	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var sizes []*models.Basic
	err = cursor.All(ctx, &sizes)
	if err != nil {
		return nil, 0, err
	}

	return sizes, total, nil
}

// Create 创建尺寸
func (s *SizeService) Create(req dto.SizeCreateRequest) (*models.Basic, error) {
	ctx := context.Background()
	now := time.Now().Unix()

	basic := &models.Basic{
		Name:      "size",
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

// Update 更新尺寸
func (s *SizeService) Update(req dto.SizeUpdateRequest) (*models.Basic, error) {
	ctx := context.Background()

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

// Delete 删除尺寸
func (s *SizeService) Delete(id string) error {
	ctx := context.Background()
	return s.repo.Delete(ctx, id)
}
