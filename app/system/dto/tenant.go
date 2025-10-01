package dto

import "mule-cloud/internal/models"

// TenantListRequest 租户请求
type TenantListRequest struct {
	ID    string `uri:"id" query:"id"`
	Code  string `query:"code"`
	Name  string `query:"name"`
	Value string `query:"value"`

	Page     int64 `query:"page"`
	PageSize int64 `query:"page_size"`
}

type TenantCreateRequest struct {
	Code   string `json:"code" binding:"required"`
	Name   string `json:"name" binding:"required"`
	Value  string `json:"value"`
	Remark string `json:"remark"`
}

type TenantUpdateRequest struct {
	ID     string `uri:"id"`
	Code   string `json:"code"`
	Name   string `json:"name"`
	Value  string `json:"value"`
	Remark string `json:"remark"`
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

