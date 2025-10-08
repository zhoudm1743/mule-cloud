package dto

import "mule-cloud/internal/models"

// TenantListRequest 租户列表请求
type TenantListRequest struct {
	ID       string `uri:"id" query:"id"`
	Code     string `query:"code"`
	Name     string `query:"name"`
	Page     int64  `query:"page"`
	PageSize int64  `query:"page_size"`
}

// TenantCreateRequest 创建租户请求
type TenantCreateRequest struct {
	Code    string `json:"code" binding"required"`
	Name    string `json:"name" binding"required"`
	Contact string `json:"contact"`
	Phone   string `json:"phone"`
	Email   string `json:"email"`
	Status  int    `json:"status"`
	// 租户管理员信息（可选）
	AdminPhone    string `json:"admin_phone"`    // 管理员手机号
	AdminPassword string `json:"admin_password"` // 管理员密码
	AdminName     string `json:"admin_name"`     // 管理员名称
}

// TenantUpdateRequest 更新租户请求
type TenantUpdateRequest struct {
	ID      string `uri:"id"`
	Code    string `json:"code"`
	Name    string `json:"name"`
	Contact string `json:"contact"`
	Phone   string `json:"phone"`
	Email   string `json:"email"`
	Status  *int   `json:"status"`
}

// AssignTenantMenusRequest 分配菜单权限给租户请求（超管使用）
type AssignTenantMenusRequest struct {
	Menus []string `json:"menus" binding"required"` // 菜单ID数组
}

// GetTenantMenusResponse 获取租户菜单权限响应
type GetTenantMenusResponse struct {
	Menus []string `json:"menus"` // 菜单ID数组
}

// TenantResponse 租户响应
type TenantResponse struct {
	Tenant *models.Tenant `json:"tenant"`
}

// TenantListResponse 租户列表响应
type TenantListResponse struct {
	Tenants []models.Tenant `json:"tenants"`
	Total   int64           `json:"total"`
}
