package dto

import "mule-cloud/internal/models"

// SalesmanRequest 业务员请求
type SalesmanListRequest struct {
	ID    string `uri:"id" query:"id"`
	Value string `query:"value"`

	Page     int64 `query:"page"`
	PageSize int64 `query:"page_size"`
}

type SalesmanCreateRequest struct {
	Value  string `json:"value"`
	Remark string `json:"remark"`
}

type SalesmanUpdateRequest struct {
	ID     string `uri:"id"`
	Value  string `json:"value"`
	Remark string `json:"remark"`
}

// SalesmanResponse 业务员响应
type SalesmanResponse struct {
	Salesman *models.Basic `json:"salesman"`
}

// SalesmanListResponse 业务员列表响应
type SalesmanListResponse struct {
	Salesmans []models.Basic `json:"salesmans"`
	Total     int64          `json:"total"`
}
