package services

import (
	"context"
	"fmt"
	"mule-cloud/internal/models"
	"mule-cloud/internal/repository"
)

type MenuService struct {
	repo *repository.MenuRepository
}

func NewMenuService() *MenuService {
	return &MenuService{
		repo: repository.NewMenuRepository(),
	}
}

// GetAll 获取所有菜单（扁平结构）
func (s *MenuService) GetAll(ctx context.Context) ([]*models.Menu, error) {
	return s.repo.GetAll(ctx)
}

// GetByID 根据ID获取菜单
func (s *MenuService) GetByID(ctx context.Context, id string) (*models.Menu, error) {
	return s.repo.GetByID(ctx, id)
}

// Create 创建菜单
func (s *MenuService) Create(ctx context.Context, menu *models.Menu) error {
	// 检查name是否重复
	existing, err := s.repo.GetByName(ctx, menu.Name)
	if err != nil {
		return err
	}
	if existing != nil {
		return fmt.Errorf("路由名称已存在: %s", menu.Name)
	}

	return s.repo.Create(ctx, menu)
}

// Update 更新菜单
func (s *MenuService) Update(ctx context.Context, id string, updates map[string]interface{}) error {
	return s.repo.Update(ctx, id, updates)
}

// Delete 删除菜单
func (s *MenuService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

// List 分页查询
func (s *MenuService) List(ctx context.Context, page, pageSize int, filters map[string]interface{}) ([]*models.Menu, int64, error) {
	return s.repo.List(ctx, page, pageSize, filters)
}

// BatchDelete 批量删除
func (s *MenuService) BatchDelete(ctx context.Context, ids []string) error {
	return s.repo.BatchDelete(ctx, ids)
}
