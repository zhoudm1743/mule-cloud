package services

import (
	"context"
	"fmt"
	"mule-cloud/app/perms/dto"
	"mule-cloud/internal/models"
	"mule-cloud/internal/repository"
	"time"
)

type DepartmentService struct {
	deptRepo repository.DepartmentRepository
}

func NewDepartmentService() *DepartmentService {
	return &DepartmentService{
		deptRepo: repository.NewDepartmentRepository(),
	}
}

// Create 创建部门
func (s *DepartmentService) Create(ctx context.Context, req *dto.CreateDepartmentRequest, createdBy string) (*models.Department, error) {
	// 检查部门编码是否已存在
	existingDept, err := s.deptRepo.FindOne(ctx, map[string]interface{}{
		"code":       req.Code,
		"is_deleted": 0,
	})
	if err != nil {
		return nil, fmt.Errorf("检查部门编码失败: %w", err)
	}
	if existingDept != nil {
		return nil, fmt.Errorf("部门编码 %s 已存在", req.Code)
	}

	// 检查部门名称是否已存在
	existingDept, err = s.deptRepo.FindOne(ctx, map[string]interface{}{
		"name":       req.Name,
		"is_deleted": 0,
	})
	if err != nil {
		return nil, fmt.Errorf("检查部门名称失败: %w", err)
	}
	if existingDept != nil {
		return nil, fmt.Errorf("部门名称 %s 已存在", req.Name)
	}

	dept := &models.Department{
		Name:      req.Name,
		Code:      req.Code,
		ParentID:  req.ParentID,
		Status:    req.Status,
		IsDeleted: 0,
		CreatedBy: createdBy,
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}

	// 默认状态为启用
	if dept.Status == 0 {
		dept.Status = 1
	}

	err = s.deptRepo.Create(ctx, dept)
	if err != nil {
		return nil, fmt.Errorf("创建部门失败: %w", err)
	}

	return dept, nil
}

// GetByID 根据ID获取部门
func (s *DepartmentService) GetByID(ctx context.Context, id string) (*models.Department, error) {
	dept, err := s.deptRepo.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("获取部门失败: %w", err)
	}
	if dept == nil {
		return nil, fmt.Errorf("部门不存在")
	}
	return dept, nil
}

// List 查询部门列表（分页）
func (s *DepartmentService) List(ctx context.Context, req *dto.ListDepartmentRequest) ([]*models.Department, int64, error) {
	// 构建查询条件
	filter := map[string]interface{}{
		"is_deleted": 0,
	}

	if req.Name != "" {
		filter["name"] = map[string]interface{}{"$regex": req.Name, "$options": "i"}
	}
	if req.Code != "" {
		filter["code"] = map[string]interface{}{"$regex": req.Code, "$options": "i"}
	}
	if req.ParentID != "" {
		filter["parent_id"] = req.ParentID
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
	total, err := s.deptRepo.Count(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("统计部门数量失败: %w", err)
	}

	// 查询数据
	depts, err := s.deptRepo.FindWithPage(ctx, filter, int64(page), int64(pageSize))
	if err != nil {
		return nil, 0, fmt.Errorf("查询部门列表失败: %w", err)
	}

	return depts, total, nil
}

// GetAll 获取所有部门（不分页）
func (s *DepartmentService) GetAll(ctx context.Context) ([]*models.Department, error) {
	filter := map[string]interface{}{
		"is_deleted": 0,
	}

	depts, err := s.deptRepo.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("获取所有部门失败: %w", err)
	}

	return depts, nil
}

// Update 更新部门
func (s *DepartmentService) Update(ctx context.Context, id string, req *dto.UpdateDepartmentRequest, updatedBy string) error {
	// 检查部门是否存在
	dept, err := s.deptRepo.Get(ctx, id)
	if err != nil {
		return fmt.Errorf("获取部门失败: %w", err)
	}
	if dept == nil {
		return fmt.Errorf("部门不存在")
	}

	// 构建更新字段
	updates := map[string]interface{}{
		"updated_by": updatedBy,
		"updated_at": time.Now().Unix(),
	}

	if req.Name != "" {
		// 检查部门名称是否已被其他部门使用
		existingDept, err := s.deptRepo.FindOne(ctx, map[string]interface{}{
			"name":       req.Name,
			"is_deleted": 0,
		})
		if err != nil {
			return fmt.Errorf("检查部门名称失败: %w", err)
		}
		if existingDept != nil && existingDept.ID != id {
			return fmt.Errorf("部门名称 %s 已存在", req.Name)
		}
		updates["name"] = req.Name
	}

	if req.Code != "" {
		// 检查部门编码是否已被其他部门使用
		existingDept, err := s.deptRepo.FindOne(ctx, map[string]interface{}{
			"code":       req.Code,
			"is_deleted": 0,
		})
		if err != nil {
			return fmt.Errorf("检查部门编码失败: %w", err)
		}
		if existingDept != nil && existingDept.ID != id {
			return fmt.Errorf("部门编码 %s 已存在", req.Code)
		}
		updates["code"] = req.Code
	}

	if req.ParentID != "" {
		updates["parent_id"] = req.ParentID
	}

	if req.Status != nil {
		updates["status"] = *req.Status
	}

	err = s.deptRepo.Update(ctx, id, updates)
	if err != nil {
		return fmt.Errorf("更新部门失败: %w", err)
	}

	return nil
}

// Delete 删除部门（软删除）
func (s *DepartmentService) Delete(ctx context.Context, id string) error {
	// 检查部门是否存在
	dept, err := s.deptRepo.Get(ctx, id)
	if err != nil {
		return fmt.Errorf("获取部门失败: %w", err)
	}
	if dept == nil {
		return fmt.Errorf("部门不存在")
	}

	// 检查是否有子部门
	children, err := s.deptRepo.Find(ctx, map[string]interface{}{
		"parent_id":  id,
		"is_deleted": 0,
	})
	if err != nil {
		return fmt.Errorf("检查子部门失败: %w", err)
	}
	if len(children) > 0 {
		return fmt.Errorf("该部门下存在子部门，无法删除")
	}

	err = s.deptRepo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("删除部门失败: %w", err)
	}

	return nil
}

// BatchDelete 批量删除部门
func (s *DepartmentService) BatchDelete(ctx context.Context, ids []string) error {
	if len(ids) == 0 {
		return fmt.Errorf("请选择要删除的部门")
	}

	for _, id := range ids {
		err := s.Delete(ctx, id)
		if err != nil {
			return fmt.Errorf("删除部门 %s 失败: %w", id, err)
		}
	}

	return nil
}

