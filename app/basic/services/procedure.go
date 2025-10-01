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

// IProcedureService 工序服务接口
type IProcedureService interface {
	Get(id string) (*models.Basic, error)
	GetAll(req dto.ProcedureListRequest) ([]models.Basic, error)
	List(req dto.ProcedureListRequest) ([]models.Basic, int64, error)
	Create(req dto.ProcedureCreateRequest) (*models.Basic, error)
	Update(req dto.ProcedureUpdateRequest) (*models.Basic, error)
	Delete(id string) error
}

// ProcedureService 工序服务实现
type ProcedureService struct {
	repo repository.BasicRepository
}

// NewProcedureService 创建工序服务
func NewProcedureService() IProcedureService {
	repo := repository.NewBasicRepository()
	return &ProcedureService{repo: repo}
}

// Get 获取工序
func (s *ProcedureService) Get(id string) (*models.Basic, error) {
	ctx := context.Background()
	return s.repo.Get(ctx, id)
}

// List 列表（分页查询）
func (s *ProcedureService) List(req dto.ProcedureListRequest) ([]models.Basic, int64, error) {
	ctx := context.Background()

	// 构建过滤条件
	filter := bson.M{"name": "procedure"}
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
		SetSort(bson.M{"created_at": -1})

	// 使用 GetCollection 获取原始集合以使用 options
	collection := s.repo.GetCollection()
	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	procedures := []models.Basic{}
	err = cursor.All(ctx, &procedures)
	if err != nil {
		return nil, 0, err
	}

	return procedures, total, nil
}

// GetAll 获取所有工序（不分页）
func (s *ProcedureService) GetAll(req dto.ProcedureListRequest) ([]models.Basic, error) {
	ctx := context.Background()

	// 构建过滤条件
	filter := bson.M{"name": "procedure"}
	if req.Value != "" {
		filter["value"] = req.Value
	}
	if req.ID != "" {
		filter["_id"] = req.ID
	}
	// 使用 GetCollection 获取原始集合以使用排序选项
	opts := options.Find().SetSort(bson.M{"created_at": -1})
	collection := s.repo.GetCollection()
	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	procedures := []models.Basic{}
	err = cursor.All(ctx, &procedures)
	if err != nil {
		return nil, err
	}

	return procedures, nil
}

// Create 创建工序
func (s *ProcedureService) Create(req dto.ProcedureCreateRequest) (*models.Basic, error) {
	ctx := context.Background()
	now := time.Now().Unix()

	basic := &models.Basic{
		Name:      "procedure",
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

// Update 更新工序
func (s *ProcedureService) Update(req dto.ProcedureUpdateRequest) (*models.Basic, error) {
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

// Delete 删除工序
func (s *ProcedureService) Delete(id string) error {
	ctx := context.Background()
	return s.repo.Delete(ctx, id)
}
