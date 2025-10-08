package services

import (
	"context"
	"mule-cloud/app/order/dto"
	"mule-cloud/internal/models"
	"mule-cloud/internal/repository"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// IStyleService 款式服务接口
type IStyleService interface {
	Get(ctx context.Context, id string) (*models.Style, error)
	GetAll(ctx context.Context, req dto.StyleListRequest) ([]models.Style, error)
	List(ctx context.Context, req dto.StyleListRequest) ([]models.Style, int64, error)
	Create(ctx context.Context, req dto.StyleCreateRequest) (*models.Style, error)
	Update(ctx context.Context, req dto.StyleUpdateRequest) (*models.Style, error)
	Delete(ctx context.Context, id string) error
}

// StyleService 款式服务实现
type StyleService struct {
	repo repository.StyleRepository
}

// NewStyleService 创建款式服务
func NewStyleService() IStyleService {
	return &StyleService{repo: repository.NewStyleRepository()}
}

// Get 获取款式
func (s *StyleService) Get(ctx context.Context, id string) (*models.Style, error) {
	return s.repo.Get(ctx, id)
}

// GetAll 获取所有款式（不分页）
func (s *StyleService) GetAll(ctx context.Context, req dto.StyleListRequest) ([]models.Style, error) {
	// 构建过滤条件
	filter := bson.M{"is_deleted": 0}

	if req.StyleNo != "" {
		filter["style_no"] = bson.M{"$regex": req.StyleNo, "$options": "i"}
	}
	if req.StyleName != "" {
		filter["style_name"] = bson.M{"$regex": req.StyleName, "$options": "i"}
	}
	if req.Category != "" {
		filter["category"] = req.Category
	}
	if req.Season != "" {
		filter["season"] = req.Season
	}
	if req.Year != "" {
		filter["year"] = req.Year
	}
	if req.Status > 0 {
		filter["status"] = req.Status
	}

	opts := options.Find().SetSort(bson.M{"created_at": -1})
	collection := s.repo.GetCollectionWithContext(ctx)
	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	styles := []models.Style{}
	err = cursor.All(ctx, &styles)
	if err != nil {
		return nil, err
	}

	return styles, nil
}

// List 列表（分页查询）
func (s *StyleService) List(ctx context.Context, req dto.StyleListRequest) ([]models.Style, int64, error) {
	// 构建过滤条件
	filter := bson.M{"is_deleted": 0}

	if req.ID != "" {
		filter["_id"] = req.ID
	}
	if req.StyleNo != "" {
		filter["style_no"] = bson.M{"$regex": req.StyleNo, "$options": "i"}
	}
	if req.StyleName != "" {
		filter["style_name"] = bson.M{"$regex": req.StyleName, "$options": "i"}
	}
	if req.Category != "" {
		filter["category"] = req.Category
	}
	if req.Season != "" {
		filter["season"] = req.Season
	}
	if req.Year != "" {
		filter["year"] = req.Year
	}
	if req.Status > 0 {
		filter["status"] = req.Status
	}

	// 获取总数
	total, err := s.repo.Count(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	// 设置分页默认值
	page := req.Page
	if page <= 0 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}

	// 分页查询
	offset := int64((page - 1) * pageSize)
	opts := options.Find().
		SetSkip(offset).
		SetLimit(int64(pageSize)).
		SetSort(bson.M{"created_at": -1})

	collection := s.repo.GetCollectionWithContext(ctx)
	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	styles := []models.Style{}
	err = cursor.All(ctx, &styles)
	if err != nil {
		return nil, 0, err
	}

	return styles, total, nil
}

// Create 创建款式
func (s *StyleService) Create(ctx context.Context, req dto.StyleCreateRequest) (*models.Style, error) {
	// 验证工序
	if len(req.Procedures) > 0 {
		if err := ValidateStyleProcedures(req.Procedures); err != nil {
			return nil, err
		}
	}

	now := time.Now().Unix()

	style := &models.Style{
		StyleNo:    req.StyleNo,
		StyleName:  req.StyleName,
		Category:   req.Category,
		Season:     req.Season,
		Year:       req.Year,
		Images:     req.Images,
		Colors:     req.Colors,
		Sizes:      req.Sizes,
		UnitPrice:  req.UnitPrice,
		Remark:     req.Remark,
		Procedures: req.Procedures,
		Status:     req.Status,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	err := s.repo.Create(ctx, style)
	if err != nil {
		return nil, err
	}

	return style, nil
}

// Update 更新款式
func (s *StyleService) Update(ctx context.Context, req dto.StyleUpdateRequest) (*models.Style, error) {
	update := bson.M{"updated_at": time.Now().Unix()}

	if req.StyleName != "" {
		update["style_name"] = req.StyleName
	}
	if req.Category != "" {
		update["category"] = req.Category
	}
	if req.Season != "" {
		update["season"] = req.Season
	}
	if req.Year != "" {
		update["year"] = req.Year
	}
	if len(req.Images) > 0 {
		update["images"] = req.Images
	}
	if len(req.Colors) > 0 {
		update["colors"] = req.Colors
	}
	if len(req.Sizes) > 0 {
		update["sizes"] = req.Sizes
	}
	if req.UnitPrice > 0 {
		update["unit_price"] = req.UnitPrice
	}
	if req.Remark != "" {
		update["remark"] = req.Remark
	}
	if len(req.Procedures) > 0 {
		// 验证工序
		if err := ValidateStyleProcedures(req.Procedures); err != nil {
			return nil, err
		}
		update["procedures"] = req.Procedures
	}
	if req.Status >= 0 {
		update["status"] = req.Status
	}

	err := s.repo.Update(ctx, req.ID, update)
	if err != nil {
		return nil, err
	}

	return s.repo.Get(ctx, req.ID)
}

// Delete 删除款式
func (s *StyleService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
