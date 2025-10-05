package services

import (
	"context"
	"fmt"
	"mule-cloud/app/perms/dto"
	"mule-cloud/core/casbin"
	"mule-cloud/core/logger"
	"mule-cloud/internal/models"
	"mule-cloud/internal/repository"
	"time"

	"go.uber.org/zap"
)

type RoleService struct {
	roleRepo   repository.RoleRepository
	tenantRepo repository.TenantRepository
	menuRepo   *repository.MenuRepository
}

func NewRoleService() *RoleService {
	return &RoleService{
		roleRepo:   repository.NewRoleRepository(),
		tenantRepo: repository.NewTenantRepository(),
		menuRepo:   repository.NewMenuRepository(),
	}
}

// Create 创建角色
func (s *RoleService) Create(ctx context.Context, req *dto.CreateRoleRequest, createdBy string) (*models.Role, error) {
	// 检查角色代码是否已存在
	existingRole, err := s.roleRepo.GetByCode(ctx, req.Code)
	if err != nil {
		return nil, fmt.Errorf("检查角色代码失败: %w", err)
	}
	if existingRole != nil {
		return nil, fmt.Errorf("角色代码 %s 已存在", req.Code)
	}

	// 检查角色名称是否已存在
	existingRole, err = s.roleRepo.GetByName(ctx, req.Name)
	if err != nil {
		return nil, fmt.Errorf("检查角色名称失败: %w", err)
	}
	if existingRole != nil {
		return nil, fmt.Errorf("角色名称 %s 已存在", req.Name)
	}

	role := &models.Role{
		Name:            req.Name,
		Code:            req.Code,
		Description:     req.Description,
		Menus:           req.Menus,
		MenuPermissions: req.MenuPermissions,
		Status:          1,
		IsDeleted:       0,
		CreatedBy:       createdBy,
		CreatedAt:       time.Now().Unix(),
		UpdatedAt:       time.Now().Unix(),
	}

	if role.Menus == nil {
		role.Menus = []string{}
	}
	if role.MenuPermissions == nil {
		role.MenuPermissions = make(map[string][]string)
	}

	err = s.roleRepo.Create(ctx, role)
	if err != nil {
		return nil, fmt.Errorf("创建角色失败: %w", err)
	}

	return role, nil
}

// GetByID 根据ID获取角色
func (s *RoleService) GetByID(ctx context.Context, id string) (*models.Role, error) {
	role, err := s.roleRepo.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("获取角色失败: %w", err)
	}
	if role == nil {
		return nil, fmt.Errorf("角色不存在")
	}
	return role, nil
}

// List 查询角色列表（分页）
func (s *RoleService) List(ctx context.Context, req *dto.ListRoleRequest) ([]*models.Role, int64, error) {
	// 构建查询条件
	filter := map[string]interface{}{
		"is_deleted": 0,
	}
	// 数据库隔离后不需要tenant_id过滤
	if req.Name != "" {
		filter["name"] = map[string]interface{}{"$regex": req.Name, "$options": "i"}
	}
	if req.Code != "" {
		filter["code"] = map[string]interface{}{"$regex": req.Code, "$options": "i"}
	}
	if req.Status != nil {
		filter["status"] = *req.Status
	}

	// 设置默认分页参数
	page := req.Page
	if page < 1 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize < 1 {
		pageSize = 10
	}

	// 查询总数
	total, err := s.roleRepo.Count(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("统计角色数量失败: %w", err)
	}

	// 查询数据
	roles, err := s.roleRepo.FindWithPage(ctx, filter, int64(page), int64(pageSize))
	if err != nil {
		return nil, 0, fmt.Errorf("查询角色列表失败: %w", err)
	}

	return roles, total, nil
}

// Update 更新角色
func (s *RoleService) Update(ctx context.Context, id string, req *dto.UpdateRoleRequest, updatedBy string) error {
	// 检查角色是否存在
	role, err := s.roleRepo.Get(ctx, id)
	if err != nil {
		return fmt.Errorf("获取角色失败: %w", err)
	}
	if role == nil {
		return fmt.Errorf("角色不存在")
	}

	// 构建更新字段
	updates := map[string]interface{}{
		"updated_by": updatedBy,
		"updated_at": time.Now().Unix(),
	}

	if req.Name != "" {
		// 检查角色名称是否已被其他角色使用
		existingRole, err := s.roleRepo.GetByName(ctx, req.Name)
		if err != nil {
			return fmt.Errorf("检查角色名称失败: %w", err)
		}
		if existingRole != nil && existingRole.ID != id {
			return fmt.Errorf("角色名称 %s 已存在", req.Name)
		}
		updates["name"] = req.Name
	}

	if req.Description != "" {
		updates["description"] = req.Description
	}

	if req.Menus != nil {
		updates["menus"] = req.Menus
	}

	if req.MenuPermissions != nil {
		updates["menu_permissions"] = req.MenuPermissions
	}

	if req.Status != nil {
		updates["status"] = *req.Status
	}

	err = s.roleRepo.Update(ctx, id, updates)
	if err != nil {
		return fmt.Errorf("更新角色失败: %w", err)
	}

	return nil
}

// Delete 删除角色（软删除）
func (s *RoleService) Delete(ctx context.Context, id string) error {
	// 检查角色是否存在
	role, err := s.roleRepo.Get(ctx, id)
	if err != nil {
		return fmt.Errorf("获取角色失败: %w", err)
	}
	if role == nil {
		return fmt.Errorf("角色不存在")
	}

	err = s.roleRepo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("删除角色失败: %w", err)
	}

	return nil
}

// BatchDelete 批量删除角色
func (s *RoleService) BatchDelete(ctx context.Context, ids []string) error {
	if len(ids) == 0 {
		return fmt.Errorf("请选择要删除的角色")
	}

	for _, id := range ids {
		err := s.roleRepo.Delete(ctx, id)
		if err != nil {
			return fmt.Errorf("删除角色 %s 失败: %w", id, err)
		}
	}

	return nil
}

// AssignMenus 分配菜单权限（支持细粒度权限）
func (s *RoleService) AssignMenus(ctx context.Context, roleID string, menuNames []string, menuPermissions map[string][]string, updatedBy string) error {
	// 检查角色是否存在
	role, err := s.roleRepo.Get(ctx, roleID)
	if err != nil {
		return fmt.Errorf("获取角色失败: %w", err)
	}
	if role == nil {
		return fmt.Errorf("角色不存在")
	}

	// 数据库隔离后不再需要租户验证

	// 更新菜单权限
	updates := map[string]interface{}{
		"menus":            menuNames,
		"menu_permissions": menuPermissions,
		"updated_by":       updatedBy,
		"updated_at":       time.Now().Unix(),
	}

	err = s.roleRepo.Update(ctx, roleID, updates)
	if err != nil {
		return fmt.Errorf("分配菜单权限失败: %w", err)
	}

	// 同步到 Casbin
	if menuPermissions != nil && len(menuPermissions) > 0 {
		// 构建菜单名到路径的映射
		menuPathMap := make(map[string]string)
		for _, menuName := range menuNames {
			menu, err := s.menuRepo.GetByName(ctx, menuName)
			if err != nil {
				logger.Warn("查询菜单失败", zap.String("menu_name", menuName), zap.Error(err))
				continue
			}
			if menu != nil {
				menuPathMap[menuName] = menu.Path
			}
		}

		// 同步权限到 Casbin
		// 数据库隔离后，Casbin只使用roleID作为唯一标识
		err = casbin.SyncRoleMenusWithPermissions("", roleID, menuPermissions, menuPathMap)
		if err != nil {
			logger.Warn("同步Casbin权限失败", zap.String("role_id", roleID), zap.Error(err))
			// 不中断流程，只记录日志
		} else {
			logger.Info("角色权限已同步到 Casbin", zap.String("role_id", roleID))
		}
	}

	return nil
}

// GetRoleMenus 获取角色的菜单权限
func (s *RoleService) GetRoleMenus(ctx context.Context, roleID string) ([]string, error) {
	role, err := s.roleRepo.Get(ctx, roleID)
	if err != nil {
		return nil, fmt.Errorf("获取角色失败: %w", err)
	}
	if role == nil {
		return nil, fmt.Errorf("角色不存在")
	}

	return role.Menus, nil
}

// GetRolesByTenant 获取租户下的所有角色
func (s *RoleService) GetRolesByTenant(ctx context.Context, tenantID string) ([]*models.Role, error) {
	roles, err := s.roleRepo.GetAllRoles(ctx)
	if err != nil {
		return nil, fmt.Errorf("获取租户角色失败: %w", err)
	}
	return roles, nil
}

// GetRolesByIDs 根据角色ID数组获取角色列表
func (s *RoleService) GetRolesByIDs(ctx context.Context, ids []string) ([]*models.Role, error) {
	roles, err := s.roleRepo.GetRolesByIDs(ctx, ids)
	if err != nil {
		return nil, fmt.Errorf("获取角色列表失败: %w", err)
	}
	return roles, nil
}
