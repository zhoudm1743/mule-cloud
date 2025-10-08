package services

import (
	"context"
	"mule-cloud/app/perms/dto"
	"mule-cloud/internal/models"
	"mule-cloud/internal/repository"
	"mule-cloud/util"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// IAdminService 管理员服务接口
type IAdminService interface {
	Get(ctx context.Context, id string) (*models.Admin, error)
	GetAll(ctx context.Context, req dto.AdminListRequest) ([]models.Admin, error)
	List(ctx context.Context, req dto.AdminListRequest) ([]models.Admin, int64, error)
	Create(ctx context.Context, req dto.AdminCreateRequest) (*models.Admin, error)
	Update(ctx context.Context, req dto.AdminUpdateRequest) (*models.Admin, error)
	Delete(ctx context.Context, id string) error
}

// AdminService 管理员服务实现
type AdminService struct {
	repo repository.AdminRepository
}

// NewAdminService 创建管理员服务
func NewAdminService() IAdminService {
	repo := repository.NewAdminRepository()
	return &AdminService{repo: repo}
}

// Get 获取管理员
func (s *AdminService) Get(ctx context.Context, id string) (*models.Admin, error) {
	return s.repo.Get(ctx, id)
}

// List 列表（分页查询）
func (s *AdminService) List(ctx context.Context, req dto.AdminListRequest) ([]models.Admin, int64, error) {
	// 构建过滤条件（排除软删除）
	filter := bson.M{"is_deleted": 0}
	if req.Phone != "" {
		filter["phone"] = req.Phone
	}
	if req.Email != "" {
		filter["email"] = req.Email
	}
	if req.Nickname != "" {
		filter["nickname"] = bson.M{"$regex": req.Nickname, "$options": "i"} // 模糊查询
	}
	// 数据库隔离后不需要tenant_id过滤
	if req.Status != nil {
		filter["status"] = *req.Status
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

	// 使用 GetCollectionWithContext 获取正确的租户数据库集合
	collection := s.repo.GetCollectionWithContext(ctx)
	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	admins := []models.Admin{}
	err = cursor.All(ctx, &admins)
	if err != nil {
		return nil, 0, err
	}

	return admins, total, nil
}

// GetAll 获取所有管理员（不分页）
func (s *AdminService) GetAll(ctx context.Context, req dto.AdminListRequest) ([]models.Admin, error) {
	// 构建过滤条件（排除软删除）
	filter := bson.M{"is_deleted": 0}
	if req.Phone != "" {
		filter["phone"] = req.Phone
	}
	if req.Email != "" {
		filter["email"] = req.Email
	}
	if req.Nickname != "" {
		filter["nickname"] = bson.M{"$regex": req.Nickname, "$options": "i"}
	}
	// 数据库隔离后不需要tenant_id过滤
	if req.Status != nil {
		filter["status"] = *req.Status
	}
	if req.ID != "" {
		filter["_id"] = req.ID
	}

	// 使用 GetCollectionWithContext 获取正确的租户数据库集合
	opts := options.Find().SetSort(bson.M{"created_at": -1})
	collection := s.repo.GetCollectionWithContext(ctx)
	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	admins := []models.Admin{}
	err = cursor.All(ctx, &admins)
	if err != nil {
		return nil, err
	}

	return admins, nil
}

// Create 创建管理员
func (s *AdminService) Create(ctx context.Context, req dto.AdminCreateRequest) (*models.Admin, error) {
	now := time.Now().Unix()
	password := util.ToolsUtil.Md5(req.Password + "mule-zdm")

	admin := &models.Admin{
		Phone:     req.Phone,
		Password:  password,
		Nickname:  req.Nickname,
		Email:     req.Email,
		Roles:     req.Roles, // 使用请求中的角色
		Avatar:    req.Avatar,
		Status:    req.Status,
		IsDeleted: 0, // 初始化为未删除
		CreatedAt: now,
		UpdatedAt: now,
	}

	err := s.repo.Create(ctx, admin)
	if err != nil {
		return nil, err
	}

	return admin, nil
}

// Update 更新管理员
func (s *AdminService) Update(ctx context.Context, req dto.AdminUpdateRequest) (*models.Admin, error) {
	// 更新字段
	update := bson.M{
		"updated_at": time.Now().Unix(),
	}
	if req.Phone != "" {
		update["phone"] = req.Phone
	}
	if req.Password != "" {
		update["password"] = util.ToolsUtil.Md5(req.Password + "mule-zdm")
	}
	if req.Nickname != "" {
		update["nickname"] = req.Nickname
	}
	if req.Email != "" {
		update["email"] = req.Email
	}
	// 数据库隔离后不需要更新tenant_id
	if req.Roles != nil {
		update["roles"] = req.Roles
	}
	if req.Avatar != "" {
		update["avatar"] = req.Avatar
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

// Delete 删除管理员
func (s *AdminService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

// AssignRoles 分配角色给管理员
func (s *AdminService) AssignRoles(ctx context.Context, adminID string, roleIDs []string, updatedBy string) error {
	// 检查管理员是否存在
	admin, err := s.repo.Get(ctx, adminID)
	if err != nil {
		return err
	}
	if admin == nil {
		return repository.ErrNotFound
	}

	// 更新角色
	update := bson.M{
		"roles":      roleIDs,
		"updated_by": updatedBy,
		"updated_at": time.Now().Unix(),
	}

	return s.repo.Update(ctx, adminID, update)
}

// GetAdminRoles 获取管理员的角色
func (s *AdminService) GetAdminRoles(ctx context.Context, adminID string) ([]string, error) {
	admin, err := s.repo.Get(ctx, adminID)
	if err != nil {
		return nil, err
	}
	if admin == nil {
		return nil, repository.ErrNotFound
	}

	if admin.Roles == nil {
		return []string{}, nil
	}

	return admin.Roles, nil
}

// RemoveRole 移除管理员的某个角色
func (s *AdminService) RemoveRole(ctx context.Context, adminID string, roleID string, updatedBy string) error {
	// 获取管理员当前角色
	admin, err := s.repo.Get(ctx, adminID)
	if err != nil {
		return err
	}
	if admin == nil {
		return repository.ErrNotFound
	}

	// 过滤掉要移除的角色
	newRoles := []string{}
	for _, r := range admin.Roles {
		if r != roleID {
			newRoles = append(newRoles, r)
		}
	}

	// 更新角色
	update := bson.M{
		"roles":      newRoles,
		"updated_by": updatedBy,
		"updated_at": time.Now().Unix(),
	}

	return s.repo.Update(ctx, adminID, update)
}
