package dto

import "mule-cloud/internal/models"

// AdminListRequest 管理员请求
type AdminListRequest struct {
	ID       string `uri:"id" form:"id"`
	Phone    string `form:"phone"`
	Email    string `form:"email"`
	Nickname string `form:"nickname"`
	TenantID string `form:"tenant_id"` // 租户ID过滤
	Status   *int   `form:"status"`

	Page     int64 `form:"page"`
	PageSize int64 `form:"page_size"`
}

type AdminCreateRequest struct {
	Phone    string   `json:"phone" binding"required"`
	Password string   `json:"password" binding"required"`
	Nickname string   `json:"nickname"`
	Email    string   `json:"email"`
	TenantID string   `json:"tenant_id"` // 租户ID（空表示系统级用户）
	Roles    []string `json:"roles"`
	Avatar   string   `json:"avatar"`
	Status   int      `json:"status"`
}

type AdminUpdateRequest struct {
	ID       string   `uri:"id"`
	Phone    string   `json:"phone"`
	Password string   `json:"password"`
	Nickname string   `json:"nickname"`
	Email    string   `json:"email"`
	TenantID *string  `json:"tenant_id"` // 租户ID（使用指针以支持更新为空）
	Roles    []string `json:"roles"`
	Avatar   string   `json:"avatar"`
	Status   *int     `json:"status"`
}

// AdminResponse 管理员响应
type AdminResponse struct {
	Admin *models.Admin `json:"admin"`
}

// AdminListResponse 管理员列表响应
type AdminListResponse struct {
	Admins []models.Admin `json:"admins"`
	Total  int64          `json:"total"`
}
