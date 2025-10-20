package dto

import "mule-cloud/internal/models"

// CustomerRequest 客户请求
type CustomerListRequest struct {
	ID    string `uri:"id" form:"id"`
	Value string `form:"value"`

	Page     int64 `form:"page"`
	PageSize int64 `form:"page_size"`
}

type CustomerCreateRequest struct {
	Value  string `json:"value"`
	Remark string `json:"remark"`
}

type CustomerUpdateRequest struct {
	ID     string `uri:"id"`
	Value  string `json:"value"`
	Remark string `json:"remark"`
}

// CustomerResponse 客户响应
type CustomerResponse struct {
	Customer *models.Basic `json:"customer"`
}

// CustomerListResponse 客户列表响应
type CustomerListResponse struct {
	Customers []models.Basic `json:"customers"`
	Total     int64          `json:"total"`
}
