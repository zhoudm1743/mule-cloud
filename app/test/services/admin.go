package services

import (
	"fmt"
	"time"
)

// Admin 管理员模型
type Admin struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

// IAdminService 管理服务接口
type IAdminService interface {
	Get(id string, userID string) (*Admin, error)
	Delete(id string, userID string) error
	GetByIds(ids []string) ([]*Admin, error)
	Create(name, email, role string) (*Admin, error)
	Update(id, name, email, role string) error
}

// AdminService 管理服务实现
type AdminService struct {
	// 实际项目中应该注入数据库连接
	data map[string]*Admin
}

// NewAdminService 创建管理服务
func NewAdminService() IAdminService {
	// 模拟数据
	data := map[string]*Admin{
		"1": {
			ID:        "1",
			Name:      "超级管理员",
			Email:     "admin@example.com",
			Role:      "admin",
			CreatedAt: time.Now().Add(-30 * 24 * time.Hour),
		},
		"2": {
			ID:        "2",
			Name:      "普通管理员",
			Email:     "manager@example.com",
			Role:      "manager",
			CreatedAt: time.Now().Add(-15 * 24 * time.Hour),
		},
		"3": {
			ID:        "3",
			Name:      "操作员",
			Email:     "operator@example.com",
			Role:      "operator",
			CreatedAt: time.Now().Add(-7 * 24 * time.Hour),
		},
	}

	return &AdminService{data: data}
}

// Get 获取管理员信息（带用户权限检查）
func (s *AdminService) Get(id string, userID string) (*Admin, error) {
	admin, exists := s.data[id]
	if !exists {
		return nil, fmt.Errorf("管理员不存在: %s", id)
	}

	// 这里可以根据 userID 做权限检查
	// 例如：普通用户只能查看自己的信息
	return admin, nil
}

// Delete 删除管理员（带权限检查）
func (s *AdminService) Delete(id string, userID string) error {
	if _, exists := s.data[id]; !exists {
		return fmt.Errorf("管理员不存在: %s", id)
	}

	// 权限检查：不能删除自己
	if id == userID {
		return fmt.Errorf("不能删除自己")
	}

	delete(s.data, id)
	return nil
}

// GetByIds 批量获取管理员（用于gRPC调用）
func (s *AdminService) GetByIds(ids []string) ([]*Admin, error) {
	admins := make([]*Admin, 0)
	for _, id := range ids {
		if admin, exists := s.data[id]; exists {
			admins = append(admins, admin)
		}
	}
	return admins, nil
}

// Create 创建管理员
func (s *AdminService) Create(name, email, role string) (*Admin, error) {
	id := fmt.Sprintf("%d", time.Now().Unix())
	admin := &Admin{
		ID:        id,
		Name:      name,
		Email:     email,
		Role:      role,
		CreatedAt: time.Now(),
	}
	s.data[id] = admin
	return admin, nil
}

// Update 更新管理员
func (s *AdminService) Update(id, name, email, role string) error {
	admin, exists := s.data[id]
	if !exists {
		return fmt.Errorf("管理员不存在: %s", id)
	}

	admin.Name = name
	admin.Email = email
	admin.Role = role
	return nil
}