package services

import (
	"context"
	"mule-cloud/app/system/dto"
	"mule-cloud/internal/models"
	"mule-cloud/internal/repository"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// ITenantService 租户服务接口
type ITenantService interface {
	Get(id string) (*models.Tenant, error)
	GetAll(req dto.TenantListRequest) ([]models.Tenant, error)
	List(req dto.TenantListRequest) ([]models.Tenant, int64, error)
	Create(req dto.TenantCreateRequest) (*models.Tenant, error)
	Update(req dto.TenantUpdateRequest) (*models.Tenant, error)
	Delete(id string) error
}

// TenantService 租户服务实现
type TenantService struct {
	repo repository.TenantRepository
}

// NewTenantService 创建租户服务
func NewTenantService() ITenantService {
	repo := repository.NewTenantRepository()
	return &TenantService{repo: repo}
}

// Get 获取租户
func (s *TenantService) Get(id string) (*models.Tenant, error) {
	ctx := context.Background()
	return s.repo.Get(ctx, id)
}

// List 列表（分页查询）
func (s *TenantService) List(req dto.TenantListRequest) ([]models.Tenant, int64, error) {
	ctx := context.Background()

	// 构建过滤条件（排除软删除）
	filter := bson.M{"is_deleted": 0}
	if req.Code != "" {
		filter["code"] = req.Code
	}
	if req.Name != "" {
		filter["name"] = req.Name
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

	tenants := []models.Tenant{}
	err = cursor.All(ctx, &tenants)
	if err != nil {
		return nil, 0, err
	}

	return tenants, total, nil
}

// GetAll 获取所有租户（不分页）
func (s *TenantService) GetAll(req dto.TenantListRequest) ([]models.Tenant, error) {
	ctx := context.Background()

	// 构建过滤条件（排除软删除）
	filter := bson.M{"is_deleted": 0}
	if req.Code != "" {
		filter["code"] = req.Code
	}
	if req.Name != "" {
		filter["name"] = req.Name
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

	tenants := []models.Tenant{}
	err = cursor.All(ctx, &tenants)
	if err != nil {
		return nil, err
	}

	return tenants, nil
}

// Create 创建租户
func (s *TenantService) Create(req dto.TenantCreateRequest) (*models.Tenant, error) {
	ctx := context.Background()
	now := time.Now().Unix()

	tenant := &models.Tenant{
		Code:      req.Code,
		Name:      req.Name,
		Contact:   req.Contact,
		Phone:     req.Phone,
		Email:     req.Email,
		Menus:     []string{}, // 初始化为空菜单数组
		Status:    req.Status,
		IsDeleted: 0, // 初始化为未删除
		CreatedAt: now,
		UpdatedAt: now,
	}

	// 如果未指定状态，默认为启用
	if tenant.Status == 0 {
		tenant.Status = 1
	}

	err := s.repo.Create(ctx, tenant)
	if err != nil {
		return nil, err
	}

	return tenant, nil
}

// Update 更新租户
func (s *TenantService) Update(req dto.TenantUpdateRequest) (*models.Tenant, error) {
	ctx := context.Background()

	// 更新字段
	update := bson.M{
		"updated_at": time.Now().Unix(),
	}
	if req.Code != "" {
		update["code"] = req.Code
	}
	if req.Name != "" {
		update["name"] = req.Name
	}
	if req.Contact != "" {
		update["contact"] = req.Contact
	}
	if req.Phone != "" {
		update["phone"] = req.Phone
	}
	if req.Email != "" {
		update["email"] = req.Email
	}
	if req.Status != nil {
		update["status"] = *req.Status
	}

	err := s.repo.Update(ctx, req.ID, update)
	if err != nil {
		return nil, err
	}

	// 返回更新后的数据
	return s.repo.Get(ctx, req.ID)
}

// Delete 删除租户
func (s *TenantService) Delete(id string) error {
	ctx := context.Background()
	return s.repo.Delete(ctx, id)
}

// AssignMenus 分配菜单权限给租户（超管使用）
func (s *TenantService) AssignMenus(ctx context.Context, tenantID string, menuIDs []string, updatedBy string) error {
	// 检查租户是否存在
	tenant, err := s.repo.Get(ctx, tenantID)
	if err != nil {
		return err
	}
	if tenant == nil {
		return repository.ErrNotFound
	}

	// 更新菜单权限到数据库
	update := bson.M{
		"menus":      menuIDs,
		"updated_by": updatedBy,
		"updated_at": time.Now().Unix(),
	}

	err = s.repo.Update(ctx, tenantID, update)
	if err != nil {
		return err
	}

	// TODO: 同步到 Casbin（可选，需要先获取菜单路径）
	// 这里可以集成 MenuRepository 来获取菜单路径
	// menuPaths := []string{}
	// for _, menuID := range menuIDs {
	//     menu, _ := menuRepo.GetByID(ctx, menuID)
	//     if menu != nil {
	//         menuPaths = append(menuPaths, menu.Path)
	//     }
	// }
	// casbinPkg.SyncTenantMenus(tenantID, menuPaths)

	return nil
}

// GetTenantMenus 获取租户的菜单权限
func (s *TenantService) GetTenantMenus(ctx context.Context, tenantID string) ([]string, error) {
	tenant, err := s.repo.Get(ctx, tenantID)
	if err != nil {
		return nil, err
	}
	if tenant == nil {
		return nil, repository.ErrNotFound
	}

	if tenant.Menus == nil {
		return []string{}, nil
	}

	return tenant.Menus, nil
}