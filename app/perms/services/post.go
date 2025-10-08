package services

import (
	"context"
	"fmt"
	"mule-cloud/app/perms/dto"
	"mule-cloud/internal/models"
	"mule-cloud/internal/repository"
	"time"
)

type PostService struct {
	postRepo repository.PostRepository
}

func NewPostService() *PostService {
	return &PostService{
		postRepo: repository.NewPostRepository(),
	}
}

// Create 创建岗位
func (s *PostService) Create(ctx context.Context, req *dto.CreatePostRequest, createdBy string) (*models.Post, error) {
	// 检查岗位编码是否已存在
	existingPost, err := s.postRepo.FindOne(ctx, map[string]interface{}{
		"code":       req.Code,
		"is_deleted": 0,
	})
	if err != nil {
		return nil, fmt.Errorf("检查岗位编码失败: %w", err)
	}
	if existingPost != nil {
		return nil, fmt.Errorf("岗位编码 %s 已存在", req.Code)
	}

	// 检查岗位名称是否已存在
	existingPost, err = s.postRepo.FindOne(ctx, map[string]interface{}{
		"name":       req.Name,
		"is_deleted": 0,
	})
	if err != nil {
		return nil, fmt.Errorf("检查岗位名称失败: %w", err)
	}
	if existingPost != nil {
		return nil, fmt.Errorf("岗位名称 %s 已存在", req.Name)
	}

	post := &models.Post{
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
	if post.Status == 0 {
		post.Status = 1
	}

	err = s.postRepo.Create(ctx, post)
	if err != nil {
		return nil, fmt.Errorf("创建岗位失败: %w", err)
	}

	return post, nil
}

// GetByID 根据ID获取岗位
func (s *PostService) GetByID(ctx context.Context, id string) (*models.Post, error) {
	post, err := s.postRepo.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("获取岗位失败: %w", err)
	}
	if post == nil {
		return nil, fmt.Errorf("岗位不存在")
	}
	return post, nil
}

// List 查询岗位列表（分页）
func (s *PostService) List(ctx context.Context, req *dto.ListPostRequest) ([]*models.Post, int64, error) {
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
	total, err := s.postRepo.Count(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("统计岗位数量失败: %w", err)
	}

	// 查询数据
	posts, err := s.postRepo.FindWithPage(ctx, filter, int64(page), int64(pageSize))
	if err != nil {
		return nil, 0, fmt.Errorf("查询岗位列表失败: %w", err)
	}

	return posts, total, nil
}

// GetAll 获取所有岗位（不分页）
func (s *PostService) GetAll(ctx context.Context) ([]*models.Post, error) {
	filter := map[string]interface{}{
		"is_deleted": 0,
	}

	posts, err := s.postRepo.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("获取所有岗位失败: %w", err)
	}

	return posts, nil
}

// Update 更新岗位
func (s *PostService) Update(ctx context.Context, id string, req *dto.UpdatePostRequest, updatedBy string) error {
	// 检查岗位是否存在
	post, err := s.postRepo.Get(ctx, id)
	if err != nil {
		return fmt.Errorf("获取岗位失败: %w", err)
	}
	if post == nil {
		return fmt.Errorf("岗位不存在")
	}

	// 构建更新字段
	updates := map[string]interface{}{
		"updated_by": updatedBy,
		"updated_at": time.Now().Unix(),
	}

	if req.Name != "" {
		// 检查岗位名称是否已被其他岗位使用
		existingPost, err := s.postRepo.FindOne(ctx, map[string]interface{}{
			"name":       req.Name,
			"is_deleted": 0,
		})
		if err != nil {
			return fmt.Errorf("检查岗位名称失败: %w", err)
		}
		if existingPost != nil && existingPost.ID != id {
			return fmt.Errorf("岗位名称 %s 已存在", req.Name)
		}
		updates["name"] = req.Name
	}

	if req.Code != "" {
		// 检查岗位编码是否已被其他岗位使用
		existingPost, err := s.postRepo.FindOne(ctx, map[string]interface{}{
			"code":       req.Code,
			"is_deleted": 0,
		})
		if err != nil {
			return fmt.Errorf("检查岗位编码失败: %w", err)
		}
		if existingPost != nil && existingPost.ID != id {
			return fmt.Errorf("岗位编码 %s 已存在", req.Code)
		}
		updates["code"] = req.Code
	}

	if req.ParentID != "" {
		updates["parent_id"] = req.ParentID
	}

	if req.Status != nil {
		updates["status"] = *req.Status
	}

	err = s.postRepo.Update(ctx, id, updates)
	if err != nil {
		return fmt.Errorf("更新岗位失败: %w", err)
	}

	return nil
}

// Delete 删除岗位（软删除）
func (s *PostService) Delete(ctx context.Context, id string) error {
	// 检查岗位是否存在
	post, err := s.postRepo.Get(ctx, id)
	if err != nil {
		return fmt.Errorf("获取岗位失败: %w", err)
	}
	if post == nil {
		return fmt.Errorf("岗位不存在")
	}

	// 检查是否有子岗位
	children, err := s.postRepo.Find(ctx, map[string]interface{}{
		"parent_id":  id,
		"is_deleted": 0,
	})
	if err != nil {
		return fmt.Errorf("检查子岗位失败: %w", err)
	}
	if len(children) > 0 {
		return fmt.Errorf("该岗位下存在子岗位，无法删除")
	}

	err = s.postRepo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("删除岗位失败: %w", err)
	}

	return nil
}

// BatchDelete 批量删除岗位
func (s *PostService) BatchDelete(ctx context.Context, ids []string) error {
	if len(ids) == 0 {
		return fmt.Errorf("请选择要删除的岗位")
	}

	for _, id := range ids {
		err := s.Delete(ctx, id)
		if err != nil {
			return fmt.Errorf("删除岗位 %s 失败: %w", id, err)
		}
	}

	return nil
}

